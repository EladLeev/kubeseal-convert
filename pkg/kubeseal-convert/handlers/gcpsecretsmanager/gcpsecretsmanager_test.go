package gcpsecretsmanager

import (
	"context"
	"errors"
	"testing"

	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"github.com/googleapis/gax-go/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

// mockGCPClient is a local implementation of gcpSecretsClientAPI for testing.
type mockGCPClient struct {
	response *secretmanagerpb.AccessSecretVersionResponse
	err      error
}

func (m *mockGCPClient) AccessSecretVersion(
	_ context.Context,
	_ *secretmanagerpb.AccessSecretVersionRequest,
	_ ...gax.CallOption,
) (*secretmanagerpb.AccessSecretVersionResponse, error) {
	return m.response, m.err
}

func (m *mockGCPClient) Close() error { return nil }

func Test_getSecret(t *testing.T) {
	t.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "../../../../test/testdata/mock_gcp_creds.json")

	tests := []struct {
		name       string
		client     gcpSecretsClientAPI
		secretName string
		want       map[string]interface{}
		wantErr    bool
	}{
		{
			name: "full secret ID returns secret value",
			client: &mockGCPClient{
				response: &secretmanagerpb.AccessSecretVersionResponse{
					Payload: &secretmanagerpb.SecretPayload{
						Data: []byte("my-secret-value"),
					},
				},
			},
			secretName: "projects/my-project/secrets/my-secret/versions/1",
			want:       map[string]interface{}{"my-secret": "my-secret-value"},
		},
		{
			name: "short secret name uses default credentials",
			client: &mockGCPClient{
				response: &secretmanagerpb.AccessSecretVersionResponse{
					Payload: &secretmanagerpb.SecretPayload{
						Data: []byte("short-name-value"),
					},
				},
			},
			secretName: "my-short-secret",
			want:       map[string]interface{}{"my-short-secret": "short-name-value"},
		},
		{
			name: "client error propagates",
			client: &mockGCPClient{
				err: errors.New("permission denied"),
			},
			secretName: "projects/my-project/secrets/my-secret/versions/1",
			want:       nil,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getSecret(context.Background(), tt.client, tt.secretName)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
