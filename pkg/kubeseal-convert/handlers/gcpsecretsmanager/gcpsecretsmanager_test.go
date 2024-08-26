package gcpsecretsmanager

import (
	"context"
	"reflect"
	"testing"
)

func Test_buildSecretId(t *testing.T) {
	t.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "../../../../test/testdata/mock_gcp_creds.json")

	ctx := context.TODO()
	secretName := "projects/my-project/secrets/my-secret/versions/1"
	expectedResult := secretName
	result := buildSecretId(ctx, secretName)
	if result != expectedResult {
		t.Errorf("Unexpected result. Got: %s, Want: %s", result, expectedResult)
	}

	secretNameOnly := "my-secret"
	expectedValue := "projects//secrets/my-secret/versions/latest"
	result = buildSecretId(ctx, secretNameOnly)
	if result != expectedValue {
		t.Errorf("Unexpected result. Got: %s, Want: %s", result, expectedValue)
	}

}

func Test_getSecret(t *testing.T) {
	type args struct {
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
			if got, _ := getSecret(context.TODO(), tt.args.secretName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getSecret() = %v, want %v", got, tt.want)
			}
		})
	}
}
