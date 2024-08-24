package vault

import (
	"context"
	"log"
	"os"

	vault "github.com/hashicorp/vault/api"

	"github.com/eladleev/kubeseal-convert/pkg/kubeseal-convert/interfaces"
)

// createClientContext creates a new Vault client with default config
// and returns context and client
func createClientContext() (context.Context, *vault.Client) {
	config := vault.DefaultConfig()
	config.Address = os.Getenv("VAULT_ADDR")

	client, err := vault.NewClient(config)
	if err != nil {
		log.Fatalf("Unable to initialize a Vault client: %v", err)
	}

	client.SetToken(os.Getenv("VAULT_TOKEN"))

	return context.Background(), client
}

// getSecret get the Vault context, client, and secret name and retrieve the secret from Vault
func getSecret(ctx context.Context, client *vault.Client, secretName string) map[string]interface{} {
	secret, err := client.KVv2("secret").Get(ctx, secretName)
	if err != nil {
		log.Fatalf("Unable to read secret from the Vault: %v", err)
	}

	return secret.Data
}

type VaultImp struct {
}

func New() interfaces.Vault {
	return &VaultImp{}
}

func (*VaultImp) GetSecret(secretName string) map[string]interface{} {
	ctx, cli := createClientContext()
	return getSecret(ctx, cli, secretName)
}
