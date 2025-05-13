// SPDX-FileCopyrightText: 2025 The Kepler Authors
// SPDX-License-Identifier: Apache-2.0

package reconciler

import (
	"context"
	"fmt"

	"github.com/sustainable.computing.io/kepler-operator/api/v1alpha1"
	powermonitor "github.com/sustainable.computing.io/kepler-operator/pkg/components/power-monitor"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// AdditionalConfigReconciler reconciles DaemonSet with annotations from ConfigMap
type AdditionalConfigReconciler struct {
	Pmi *v1alpha1.PowerMonitorInternal
	Ds  *appsv1.DaemonSet
}

// Reconcile implements the AdditionalConfigReconciler interface
func (r AdditionalConfigReconciler) Reconcile(ctx context.Context, cli client.Client, s *runtime.Scheme) Result {
	cfm := &corev1.ConfigMap{}
	if err := cli.Get(ctx, types.NamespacedName{
		Name:      r.Pmi.Name,
		Namespace: r.Pmi.Namespace(),
	}, cfm); err != nil {
		if errors.IsNotFound(err) {
			return Updater{Owner: r.Pmi, Resource: r.Ds}.Reconcile(ctx, cli, s)
		}
		return Result{Action: Stop, Error: fmt.Errorf("error getting ConfigMap: %w", err)}
	}

	powermonitor.MountConfigMapToDaemonSet(r.Ds, cfm)

	// Update the DaemonSet
	return Updater{Owner: r.Pmi, Resource: r.Ds}.Reconcile(ctx, cli, s)
}

// PowerMonitorConfigMapReconciler reconciles ConfigMap with user-provided ConfigMaps
type PowerMonitorConfigMapReconciler struct {
	Pmi *v1alpha1.PowerMonitorInternal
	Cfm *corev1.ConfigMap
}

// Reconcile implements the PowerMonitorConfigMapReconciler interface
func (r PowerMonitorConfigMapReconciler) Reconcile(ctx context.Context, cli client.Client, s *runtime.Scheme) Result {
	configMaps := r.Pmi.Spec.Kepler.Config.AdditionalConfigMaps
	if len(configMaps) == 0 {
		return Updater{Owner: r.Pmi, Resource: r.Cfm}.Reconcile(ctx, cli, s)
	}

	// Get the merged config data
	configData, err := r.getMergedConfig(ctx, cli, configMaps)
	if err != nil {
		return Result{Action: Stop, Error: fmt.Errorf("error getting ConfigMaps: %w", err)}
	}

	// Assign the merged config data directly to the ConfigMap
	r.Cfm.Data = configData

	// Update the ConfigMap in the cluster
	return Updater{Owner: r.Pmi, Resource: r.Cfm}.Reconcile(ctx, cli, s)
}

// getMergedConfig fetches the ConfigMaps referenced in the spec, merges them, and returns the final config data
func (r PowerMonitorConfigMapReconciler) getMergedConfig(ctx context.Context, cli client.Client, refs []v1alpha1.ConfigMapRef) (map[string]string, error) {
	configMaps := make([]*corev1.ConfigMap, 0, len(refs))
	ns := r.Pmi.Namespace()

	// Fetch all referenced ConfigMaps
	for _, ref := range refs {
		cfm := &corev1.ConfigMap{}
		if err := cli.Get(ctx, types.NamespacedName{Namespace: ns, Name: ref.Name}, cfm); err != nil {
			if errors.IsNotFound(err) {
				return nil, fmt.Errorf("ConfigMap %q not found in %q namespace", ref.Name, ns)
			}
			return nil, fmt.Errorf("failed to get ConfigMap %q: %w", ref.Name, err)
		}
		configMaps = append(configMaps, cfm)
	}

	// Extract YAML configurations from custom ConfigMaps
	var additionalConfigs []string
	for _, cfm := range configMaps {
		if cfm != nil && cfm.Data != nil {
			for filename, content := range cfm.Data {
				if content != "" && filename == powermonitor.KeplerConfigFile {
					additionalConfigs = append(additionalConfigs, content)
				}
			}
		}
	}

	// Get the merged config
	mergedConfig, err := powermonitor.KeplerConfig(r.Pmi, additionalConfigs...)
	if err != nil {
		return nil, fmt.Errorf("failed to merge configurations: %w", err)
	}

	// Return the config data that should be assigned to the original ConfigMap
	configData := map[string]string{
		powermonitor.KeplerConfigFile: mergedConfig,
	}

	return configData, nil
}
