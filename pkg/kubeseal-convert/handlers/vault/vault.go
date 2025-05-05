package vault

import (
	"context"
	"fmt"
	"os"
	"time"

	vault "github.com/hashicorp/vault/api"
	log "github.com/sirupsen/logrus"

	"github.com/eladleev/kubeseal-convert/pkg/kubeseal-convert/interfaces"
)

// createClient creates a new Vault client with default config
// and returns context and client
func createClient() (*vault.Client, error) {
	config := vault.DefaultConfig()
	config.Address = os.Getenv("VAULT_ADDR")

	client, err := vault.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize a Vault client: %v", err)
	}

	client.SetToken(os.Getenv("VAULT_TOKEN"))

	return client, nil
}

// getSecret get the Vault context, client, and secret name and retrieve the secret from Vault
func getSecret(
	ctx context.Context,
	client *vault.Client,
	secretName string,
) (map[string]interface{}, error) {
	secret, err := client.KVv2("secret").Get(ctx, secretName)
	log.Debugf("secret: %v", secret)
	if err != nil {
		return nil, fmt.Errorf("unable to read secret from the Vault: %v", err)
	}

	return secret.Data, nil
}

type VaultImp struct{}

func New() interfaces.Vault {
	return &VaultImp{}
}

func (*VaultImp) GetSecret(secretName string, timeout int) map[string]interface{} {
	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	cli, err := createClient()
	if err != nil {
		log.Fatalf("unable to initialize a Vault client: %v", err)
	}

	log.Debugf("ctx: %v", ctx)

	secret, err := getSecret(ctx, cli, secretName)
	if err != nil {
		log.Fatalf("failed to get secret: %v", err)
	}
	return secret
}
