package util

import (
	"reflect"
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestGetSecret(t *testing.T) {
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
				name:      "helm-secret-values",
			},
			want: &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "helm-secret-values",
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
			// client := GetK8sClient()
			_, err := GetSecret(tt.args.namespace, tt.args.name, client)
			if err.Error() != "secrets \"helm-secret-values\" not found" {
				t.Errorf("Incorrect error message when %q not found", tt.args.name)
			}
		})
	}
}

func TestComposeSecretValues(t *testing.T) {
	type args struct {
		secret *v1.Secret
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Should get the right value out",
			args: args{
				secret: &v1.Secret{
					Data: map[string][]byte{
						"values.yaml": []byte("replicas: \"3\"\ndeployment:\n  server:\n    replicas: \"3\"\n"),
					},
				},
			},
			want: "replicas: \"3\"\ndeployment:\n  server:\n    replicas: \"3\"\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ComposeSecretValues(tt.args.secret); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetConfigMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
