package secretsmanager

import (
	"context"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

func Test_getSecret(t *testing.T) {
	type args struct {
		svc        *secretsmanager.Client
		secretName string
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := getSecret(context.TODO(), tt.args.svc, tt.args.secretName); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("getSecret() = %v, want %v", got, tt.want)
			}
		})
	}
}
