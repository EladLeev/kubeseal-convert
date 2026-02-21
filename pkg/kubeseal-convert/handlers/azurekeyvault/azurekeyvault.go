package azurekeyvault

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	identity "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	secrets "github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"

	"github.com/eladleev/kubeseal-convert/pkg/kubeseal-convert/interfaces"
)

// azureSecretsClientAPI is the subset of secrets.Client used by this handler.
type azureSecretsClientAPI interface {
	NewListSecretsPager(
		opts *secrets.ListSecretsOptions,
	) *runtime.Pager[secrets.ListSecretsResponse]
	GetSecret(
		ctx context.Context,
		name, version string,
		opts *secrets.GetSecretOptions,
	) (secrets.GetSecretResponse, error)
}

func createClient(vaultName string) (*secrets.Client, error) {
	/*
		see https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity#readme-defaultazurecredential: this allows getting credentials
		via either environment variables, managed identity, or 'az login'
	*/
	cred, err := identity.NewDefaultAzureCredential(nil)
	log.Debugf("Azure identity: %v", cred)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to obtain a credential needed to login to the azure vault: %v",
			err,
		)
	}

	vaultURI := fmt.Sprintf("https://%s.vault.azure.net", vaultName)
	log.Debugf("vaultURI: %v", vaultURI)
	client, err := secrets.NewClient(vaultURI, cred, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to vault '%s': %v", vaultURI, err)
	}
	return client, nil
}

// retrieve secret by name with the client
func getSecrets(
	ctx context.Context,
	client azureSecretsClientAPI,
	vaultName string,
) (map[string]interface{}, error) {
	mp := make(map[string]interface{})

	pager := client.NewListSecretsPager(&secrets.ListSecretsOptions{})

	for pager.More() {
		log.Debugf("pager: %v", pager)
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve secrets from vault '%s': %v", vaultName, err)
		}
		for _, secret := range page.Value {
			value, err := client.GetSecret(
				ctx,
				secret.ID.Name(),
				secret.ID.Version(),
				&secrets.GetSecretOptions{},
			)
			log.Debugf("secret value: %v", value)
			if err != nil {
				return nil, fmt.Errorf(
					"failed to retrieve secret '%s' from vault '%s': %v",
					secret.ID.Name(),
					vaultName,
					err,
				)
			}
			mp[secret.ID.Name()] = *value.Value
		}
	}

	return mp, nil
}

type AzureKeyVaultImp struct {
	client azureSecretsClientAPI
}

func New() interfaces.AzureKeyVault {
	return &AzureKeyVaultImp{}
}

func (a *AzureKeyVaultImp) GetSecrets(vaultName string, timeout int) map[string]interface{} {
	log.Debugf("Getting secrets from vault %v", vaultName)
	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	cli := a.client
	if cli == nil {
		var err error
		cli, err = createClient(vaultName)
		if err != nil {
			log.Fatalf("failed to create client: %v", err)
		}
	}

	secret, err := getSecrets(ctx, cli, vaultName)
	if err != nil {
		log.Fatalf("failed to get secret: %v", err)
	}
	return secret
}
