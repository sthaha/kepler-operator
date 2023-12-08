package modelserver

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sustainable.computing.io/kepler-operator/pkg/api/v1alpha1"
	"github.com/sustainable.computing.io/kepler-operator/pkg/components"
)

func TestConfigMap(t *testing.T) {

	tt := []struct {
		spec     *v1alpha1.InternalModelServerSpec
		data     map[string]string
		scenario string
	}{
		{
			spec: &v1alpha1.InternalModelServerSpec{},
			data: map[string]string{
				"MODEL_PATH": "/mnt/models",
			},
			scenario: "default case",
		},
		{
			spec: &v1alpha1.InternalModelServerSpec{
				URL:         "fake-url",
				Path:        "fake-model-path",
				RequestPath: "fake-request-path",
				ListPath:    "fake-model-list-path",
				PipelineURL: "fake-pipeline",
				ErrorKey:    "fake-error-key",
			},
			data: map[string]string{
				"MODEL_PATH":                   "fake-model-path",
				"MODEL_SERVER_REQ_PATH":        "fake-request-path",
				"MODEL_SERVER_MODEL_LIST_PATH": "fake-model-list-path",
				"INITIAL_PIPELINE_URL":         "fake-pipeline",
				"ERROR_KEY":                    "fake-error-key",
			},
			scenario: "user defined server-api config",
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()
			k := v1alpha1.KeplerInternal{
				Spec: v1alpha1.KeplerInternalSpec{
					ModelServer: tc.spec,
				},
			}
			actual := NewConfigMap(k.ModelServerDeploymentName(), components.Full, k.Spec.ModelServer, k.Spec.Exporter.Deployment.Namespace)
			assert.Equal(t, len(tc.data), len(actual.Data))
			for k, v := range tc.data {
				assert.Equal(t, v, actual.Data[k])
			}
		})
	}

}

func TestService(t *testing.T) {

	tt := []struct {
		spec        v1alpha1.InternalModelServerSpec
		servicePort int32
		scenario    string
	}{
		{
			spec: v1alpha1.InternalModelServerSpec{
				Port: 8080,
			},
			servicePort: 8080,
			scenario:    "user defined port",
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()
			k := v1alpha1.KeplerInternal{
				Spec: v1alpha1.KeplerInternalSpec{
					ModelServer: &tc.spec,
				},
			}
			actual := NewService(k.ModelServerDeploymentName(), k.Spec.ModelServer, k.Spec.Exporter.Deployment.Namespace)
			assert.Equal(t, actual.Spec.Ports[0].Port, tc.servicePort)
		})
	}

}

func TestServerAPIContainer(t *testing.T) {

	tt := []struct {
		spec        v1alpha1.InternalModelServerSpec
		servicePort int32
		scenario    string
	}{
		{
			spec: v1alpha1.InternalModelServerSpec{
				Port: 8080,
			},
			servicePort: 8080,
			scenario:    "user defined port",
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()
			k := v1alpha1.KeplerInternal{
				Spec: v1alpha1.KeplerInternalSpec{
					ModelServer: &tc.spec,
				},
			}
			actual := NewDeployment(k.ModelServerDeploymentName(), k.Spec.ModelServer, k.Spec.Exporter.Deployment.Namespace)
			containers := actual.Spec.Template.Spec.Containers
			assert.Equal(t, len(containers), 1)
			assert.Equal(t, containers[0].Ports[0].ContainerPort, tc.servicePort)
		})
	}

}
