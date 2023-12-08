package e2e

import (
	"testing"

	"github.com/sustainable.computing.io/kepler-operator/pkg/api/v1alpha1"
	"github.com/sustainable.computing.io/kepler-operator/pkg/controllers"
	"github.com/sustainable.computing.io/kepler-operator/pkg/utils/test"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

func TestModelServer_Reconciliation(t *testing.T) {
	f := test.NewFramework(t)
	// pre-condition
	f.AssertNoResourceExists("model-server", "", &v1alpha1.KeplerInternal{}, test.NoWait())

	// when
	ki := f.CreateInternal("model-server", f.WithDefaultExporter("model-server-e2e-test"))

	// then
	f.AssertResourceExists(controllers.KeplerDeploymentNS, "", &corev1.Namespace{})
	ds := appsv1.Deployment{}
	f.AssertResourceExists(ki.ModelServerDeploymentName(), controllers.KeplerDeploymentNS, &ds)
}
