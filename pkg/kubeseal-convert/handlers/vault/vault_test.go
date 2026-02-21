package vault

import (
	"context"
	"errors"
	"testing"

	vault "github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockKVv2 is a local implementation of vaultKVv2API for testing.
type mockKVv2 struct {
	secret *vault.KVSecret
	err    error
}

func (m *mockKVv2) Get(_ context.Context, _ string) (*vault.KVSecret, error) {
	return m.secret, m.err
}

func Test_getSecret(t *testing.T) {
	tests := []struct {
		name       string
		kv         vaultKVv2API
		secretName string
		want       map[string]interface{}
		wantErr    bool
	}{
		{
			name: "secret found",
			kv: &mockKVv2{
				secret: &vault.KVSecret{
					Data: map[string]interface{}{
						"username": "admin",
						"password": "s3cr3t",
					},
				},
			},
			secretName: "dev/my-secret",
			want:       map[string]interface{}{"username": "admin", "password": "s3cr3t"},
		},
		{
			name: "secret not found returns error",
			kv: &mockKVv2{
				err: errors.New("secret not found"),
			},
			secretName: "dev/missing-secret",
			want:       nil,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getSecret(context.Background(), tt.kv, tt.secretName)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
