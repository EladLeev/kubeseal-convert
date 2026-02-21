package azurekeyvault

import (
	"context"
	"errors"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	secrets "github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockAzureClient is a local implementation of azureSecretsClientAPI for testing.
type mockAzureClient struct {
	listPages []secrets.ListSecretsResponse
	listErr   error
	secretMap map[string]string
	getErr    error
}

func (m *mockAzureClient) NewListSecretsPager(
	_ *secrets.ListSecretsOptions,
) *runtime.Pager[secrets.ListSecretsResponse] {
	pages := m.listPages
	idx := 0
	return runtime.NewPager(runtime.PagingHandler[secrets.ListSecretsResponse]{
		More: func(_ secrets.ListSecretsResponse) bool {
			return idx < len(pages)
		},
		Fetcher: func(_ context.Context, _ *secrets.ListSecretsResponse) (secrets.ListSecretsResponse, error) {
			if m.listErr != nil {
				return secrets.ListSecretsResponse{}, m.listErr
			}
			if idx >= len(pages) {
				return secrets.ListSecretsResponse{}, errors.New("no more pages")
			}
			page := pages[idx]
			idx++
			return page, nil
		},
	})
}

func (m *mockAzureClient) GetSecret(
	_ context.Context,
	name, _ string,
	_ *secrets.GetSecretOptions,
) (secrets.GetSecretResponse, error) {
	if m.getErr != nil {
		return secrets.GetSecretResponse{}, m.getErr
	}
	val, ok := m.secretMap[name]
	if !ok {
		return secrets.GetSecretResponse{}, errors.New("secret not found")
	}
	return secrets.GetSecretResponse{SecretBundle: secrets.SecretBundle{Value: &val}}, nil
}

// secretID creates a minimal azsecrets ID string in the expected format.
func secretID(name string) secrets.ID {
	return secrets.ID("https://myvault.vault.azure.net/secrets/" + name + "/abc")
}

func Test_getSecrets(t *testing.T) {
	tests := []struct {
		name      string
		client    azureSecretsClientAPI
		vaultName string
		want      map[string]interface{}
		wantErr   bool
	}{
		{
			name: "multiple secrets retrieved successfully",
			client: &mockAzureClient{
				listPages: []secrets.ListSecretsResponse{
					{
						SecretListResult: secrets.SecretListResult{
							Value: func() []*secrets.SecretItem {
								id1 := secretID("db-password")
								id2 := secretID("api-key")
								return []*secrets.SecretItem{
									{ID: &id1},
									{ID: &id2},
								}
							}(),
						},
					},
				},
				secretMap: map[string]string{
					"db-password": "s3cr3t",
					"api-key":     "key-abc123",
				},
			},
			vaultName: "my-vault",
			want: map[string]interface{}{
				"db-password": "s3cr3t",
				"api-key":     "key-abc123",
			},
		},
		{
			name: "empty vault returns empty map",
			client: &mockAzureClient{
				listPages: []secrets.ListSecretsResponse{
					{
						SecretListResult: secrets.SecretListResult{
							Value: []*secrets.SecretItem{},
						},
					},
				},
				secretMap: map[string]string{},
			},
			vaultName: "empty-vault",
			want:      map[string]interface{}{},
		},
		{
			name: "list pager error propagates",
			client: &mockAzureClient{
				listErr: errors.New("vault unreachable"),
			},
			vaultName: "my-vault",
			want:      nil,
			wantErr:   true,
		},
		{
			name: "GetSecret error propagates",
			client: &mockAzureClient{
				listPages: []secrets.ListSecretsResponse{
					{
						SecretListResult: secrets.SecretListResult{
							Value: func() []*secrets.SecretItem {
								id := secretID("mysecret")
								return []*secrets.SecretItem{{ID: &id}}
							}(),
						},
					},
				},
				getErr: errors.New("permission denied"),
			},
			vaultName: "my-vault",
			want:      nil,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getSecrets(context.Background(), tt.client, tt.vaultName)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
