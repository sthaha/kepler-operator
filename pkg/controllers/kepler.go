package controllers

import (
	"context"
	"strings"

	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	"github.com/sustainable.computing.io/kepler-operator/pkg/api/v1alpha1"
	"github.com/sustainable.computing.io/kepler-operator/pkg/reconciler"
	"github.com/sustainable.computing.io/kepler-operator/pkg/utils/k8s"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/util/retry"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	ctrl "sigs.k8s.io/controller-runtime"
)

const (
	Finalizer                       = "kepler.system.sustainable.computing.io/finalizer"
	KeplerBpfAttachMethodAnnotation = "kepler.sustainable.computing.io/bpf-attach-method"
	KeplerBpfAttachMethodBCC        = "bcc"
	KeplerBpfAttachMethodLibbpf     = "libbpf"
)

// Config that will be set from outside
var (
	Config = struct {
		Image       string
		ImageLibbpf string
	}{}
)

// KeplerReconciler reconciles a Kepler object
type KeplerReconciler struct {
	client.Client
	Scheme *runtime.Scheme

	logger  logr.Logger
	Cluster k8s.Cluster
}

// Owned resource
//+kubebuilder:rbac:groups=kepler.system.sustainable.computing.io,resources=*,verbs=*

// SetupWithManager sets up the controller with the Manager.
func (r *KeplerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Kepler{}).
		Owns(&v1alpha1.KeplerInternal{},
			builder.WithPredicates(predicate.ResourceVersionChangedPredicate{})).
		Complete(r)
}

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Kepler object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *KeplerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// TODO: remove these keys from the log
	// "controller": "kepler", "controllerGroup": "kepler.system.sustainable.computing.io",
	// "controllerKind": "Kepler", "Kepler": {"name":"kepler"},

	logger := log.FromContext(ctx)
	r.logger = logger

	logger.Info("Start of  reconcile")
	defer logger.Info("End of reconcile")

	kepler, err := r.getKepler(ctx, req)
	if err != nil {
		// retry since some error has occurred
		logger.V(6).Info("Get Error ", "error", err)
		return ctrl.Result{}, err
	}

	if kepler == nil {
		// no kepler found , so stop here
		logger.V(6).Info("Kepler Nil")
		return ctrl.Result{}, nil
	}

	// TODO: have admission webhook decline all but `kepler` instance
	if kepler.Name != "kepler" {
		return r.setInvalidStatus(ctx, req)
	}

	logger.V(6).Info("Running sub reconcilers", "kepler", kepler.Spec)

	result, recErr := r.runKeplerReconcilers(ctx, kepler)
	updateErr := r.updateStatus(ctx, req, err)

	if recErr != nil {
		return result, recErr
	}
	return result, updateErr
}

func (r KeplerReconciler) runKeplerReconcilers(ctx context.Context, kepler *v1alpha1.Kepler) (ctrl.Result, error) {

	reconcilers := r.reconcilersForKepler(kepler)
	r.logger.V(6).Info("renconcilers ...", "count", len(reconcilers))

	return reconciler.Runner{
		Reconcilers: reconcilers,
		Client:      r.Client,
		Scheme:      r.Scheme,
		Logger:      r.logger,
	}.Run(ctx)
}

func (r KeplerReconciler) updateStatus(ctx context.Context, req ctrl.Request, recErr error) error {
	return retry.RetryOnConflict(retry.DefaultBackoff, func() error {

		k, _ := r.getKepler(ctx, req)
		// may be deleted
		if k == nil || !k.GetDeletionTimestamp().IsZero() {
			// retry since some error has occurred
			r.logger.V(6).Info("kepler has been deleted; skipping status update")
			return nil
		}

		internal, _ := r.getInternalForKepler(ctx, k)
		// may be deleted
		if internal == nil || !internal.GetDeletionTimestamp().IsZero() {
			// retry since some error has occurred
			r.logger.V(6).Info("keplerinternal has deleted; skipping status update")
			return nil
		}
		if !hasInternalStatusChanged(internal) {
			r.logger.V(6).Info("keplerinternal has not changed; skipping status update")
			return nil
		}

		// NOTE: although, this copies the internal status, the observed generation
		// should be set to kepler's current generation to indicate that the
		// current generation has been "observed"
		k.Status = internal.Status
		for i := range k.Status.Conditions {
			k.Status.Conditions[i].ObservedGeneration = k.Generation
		}
		return r.Client.Status().Update(ctx, k)
	})
}

// returns true (i.e. status has changed ) if any of the Conditions'
// ObservedGeneration is equal to the current generation
func hasInternalStatusChanged(internal *v1alpha1.KeplerInternal) bool {
	for i := range internal.Status.Conditions {
		if internal.Status.Conditions[i].ObservedGeneration == internal.Generation {
			return true
		}
	}
	return false
}

func (r KeplerReconciler) getKepler(ctx context.Context, req ctrl.Request) (*v1alpha1.Kepler, error) {
	logger := r.logger

	kepler := v1alpha1.Kepler{}

	if err := r.Client.Get(ctx, req.NamespacedName, &kepler); err != nil {
		if errors.IsNotFound(err) {
			logger.V(3).Info("kepler could not be found; may be marked for deletion")
			return nil, nil
		}
		logger.Error(err, "failed to get kepler")
		return nil, err
	}

	return &kepler, nil
}

func (r KeplerReconciler) getInternalForKepler(ctx context.Context, k *v1alpha1.Kepler) (*v1alpha1.KeplerInternal, error) {
	logger := r.logger.WithValues("kepler-internal", k.Name)

	internal := v1alpha1.KeplerInternal{}
	if err := r.Client.Get(ctx, client.ObjectKey{Name: k.Name}, &internal); err != nil {
		if errors.IsNotFound(err) {
			logger.V(3).Info("kepler-internal could not be found; may be marked for deletion")
			return nil, nil
		}
		logger.Error(err, "failed to get kepler-internal")
		return nil, err
	}
	return &internal, nil
}

func (r KeplerReconciler) reconcilersForKepler(k *v1alpha1.Kepler) []reconciler.Reconciler {
	op := deleteResource
	if update := k.DeletionTimestamp.IsZero(); update {
		op = newUpdaterWithOwner(k)
	}

	rs := []reconciler.Reconciler{
		op(newKeplerInternal(k)),
		reconciler.Finalizer{
			Resource: k, Finalizer: Finalizer, Logger: r.logger,
		},
	}
	return rs
}

func (r KeplerReconciler) setInvalidStatus(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	err := retry.RetryOnConflict(retry.DefaultBackoff, func() error {
		invalidKepler, _ := r.getKepler(ctx, req)
		// may be deleted
		if invalidKepler == nil || !invalidKepler.GetDeletionTimestamp().IsZero() {
			return nil
		}

		invalidKepler.Status.Conditions = []v1alpha1.Condition{{
			Type:               v1alpha1.Reconciled,
			Status:             v1alpha1.ConditionFalse,
			ObservedGeneration: invalidKepler.Generation,
			LastTransitionTime: metav1.Now(),
			Reason:             v1alpha1.InvalidKeplerResource,
			Message:            "Only a single instance of Kepler named kepler is reconciled",
		}, {
			Type:               v1alpha1.Available,
			Status:             v1alpha1.ConditionUnknown,
			ObservedGeneration: invalidKepler.Generation,
			LastTransitionTime: metav1.Now(),
			Reason:             v1alpha1.InvalidKeplerResource,
			Message:            "This instance of Kepler is invalid",
		}}
		return r.Client.Status().Update(ctx, invalidKepler)
	})

	// retry only on error
	return ctrl.Result{}, err
}

func newKeplerInternal(k *v1alpha1.Kepler) *v1alpha1.KeplerInternal {

	keplerImage := Config.Image
	if IsLibbpfAttachType(k) {
		keplerImage = Config.ImageLibbpf
	}

	return &v1alpha1.KeplerInternal{
		TypeMeta: metav1.TypeMeta{
			Kind:       "KeplerInternal",
			APIVersion: v1alpha1.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        k.Name,
			Annotations: k.Annotations,
		},
		Spec: v1alpha1.KeplerInternalSpec{
			Exporter: v1alpha1.InternalExporterSpec{
				Deployment: v1alpha1.InternalExporterDeploymentSpec{
					ExporterDeploymentSpec: k.Spec.Exporter.Deployment,
					Image:                  keplerImage,
				},
			},
		},
	}
}

func IsLibbpfAttachType(k *v1alpha1.Kepler) bool {
	bpftype, ok := k.Annotations[KeplerBpfAttachMethodAnnotation]
	return ok && strings.ToLower(bpftype) == KeplerBpfAttachMethodLibbpf
}
