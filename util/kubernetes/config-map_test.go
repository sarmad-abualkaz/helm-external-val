package util

import (
	"reflect"
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestGetConfigMap(t *testing.T) {
	type args struct {
		namespace string
		name      string
	}
	tests := []struct {
		name string
		args args
		want *v1.ConfigMap
	}{
		{
			name: "Should fail and log not found",
			args: args{
				namespace: "kuuji",
				name:      "helm-values",
			},
			want: &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "helm-values",
					Namespace: "kuuji",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := Client{
				Clientset: fake.NewSimpleClientset(),
			}
			_, err := GetConfigMap(tt.args.namespace, tt.args.name, client)
			if err.Error() != "configmaps \"helm-values\" not found" {
				t.Errorf("Incorrect error message when %q not found", tt.args.name)
			}
		})
	}
}

func TestComposeValues(t *testing.T) {
	type args struct {
		configmap *v1.ConfigMap
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Should create file",
			args: args{
				configmap: &v1.ConfigMap{
					Data: map[string]string{
						"values.yaml": "replicas: \"3\"\ndeployment:\n  server:\n    replicas: \"3\"\n",
					},
				},
			},
			want: "replicas: \"3\"\ndeployment:\n  server:\n    replicas: \"3\"\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ComposeValues(tt.args.configmap); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetConfigMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
