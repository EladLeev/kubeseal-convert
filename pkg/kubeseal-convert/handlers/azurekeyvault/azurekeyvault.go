package azurekeyvault

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"

	identity "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	secrets "github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"

	"github.com/eladleev/kubeseal-convert/pkg/kubeseal-convert/interfaces"
)

func createClient(vaultName string) *secrets.Client {
	// see https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity#readme-defaultazurecredential: this allows getting credentials
	// via either environment variables, managed identity, or 'az login'
	cred, err := identity.NewDefaultAzureCredential(nil)
	log.Debugf("Azure identity: %v", cred)
	if err != nil {
		log.Fatalf("Failed to obtain a credential needed to login to the azure vault: %v", err)
	}

	vaultURI := fmt.Sprintf("https://%s.vault.azure.net", vaultName)
	log.Debugf("vaultURI: %v", vaultURI)
	client, err := secrets.NewClient(vaultURI, cred, nil)
	if err != nil {
		log.Fatalf("Failed to connect to vault '%s': %v", vaultURI, err)
	}
	return client
}

// retrieve secret by name with the client
func getSecrets(client *secrets.Client, vaultName string) map[string]interface{} {
	mp := make(map[string]interface{})

	pager := client.NewListSecretsPager(&secrets.ListSecretsOptions{})

	for pager.More() {
		log.Debugf("pager: %v", pager)
		page, err := pager.NextPage(context.TODO())
		if err != nil {
			log.Fatalf("Failed to retrieve secrets from vault '%s': %v", vaultName, err)
		}
		for _, secret := range page.Value {
			value, err := client.GetSecret(context.TODO(), secret.ID.Name(), secret.ID.Version(), &secrets.GetSecretOptions{})
			log.Debugf("secret value: %v", value)
			if err != nil {
				log.Fatalf("Failed to retrieve secret '%s' from vault '%s': %v", secret.ID.Name(), vaultName, err)
			}
			mp[secret.ID.Name()] = *value.Value
		}
	}

	return mp
}

type AzureKeyVaultImp struct {
}

func New() interfaces.AzureKeyVault {
	return &AzureKeyVaultImp{}
}

func (*AzureKeyVaultImp) GetSecrets(vaultName string) map[string]interface{} {
	log.Debugf("Getting secrets from vault %v", vaultName)
	cli := createClient(vaultName)
	return getSecrets(cli, vaultName)
}
