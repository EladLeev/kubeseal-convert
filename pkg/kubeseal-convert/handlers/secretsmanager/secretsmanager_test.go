package secretsmanager

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockSMClient is a local implementation of secretsManagerAPI for testing.
type mockSMClient struct {
	output *secretsmanager.GetSecretValueOutput
	err    error
}

func (m *mockSMClient) GetSecretValue(
	_ context.Context,
	_ *secretsmanager.GetSecretValueInput,
	_ ...func(*secretsmanager.Options),
) (*secretsmanager.GetSecretValueOutput, error) {
	return m.output, m.err
}

func Test_getSecret(t *testing.T) {
	tests := []struct {
		name       string
		client     secretsManagerAPI
		secretName string
		want       map[string]interface{}
		wantErr    bool
	}{
		{
			name: "json secret happy path",
			client: &mockSMClient{
				output: &secretsmanager.GetSecretValueOutput{
					SecretString: aws.String(`{"username":"admin","password":"s3cr3t"}`),
				},
			},
			secretName: "dev/my-secret",
			want:       map[string]interface{}{"username": "admin", "password": "s3cr3t"},
		},
		{
			name: "non-json secret value returns error",
			client: &mockSMClient{
				output: &secretsmanager.GetSecretValueOutput{
					SecretString: aws.String("plain-text-not-json"),
				},
			},
			secretName: "dev/plain-secret",
			want:       nil,
			wantErr:    true,
		},
		{
			name: "aws api error propagates",
			client: &mockSMClient{
				err: errors.New("ResourceNotFoundException"),
			},
			secretName: "dev/missing-secret",
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
