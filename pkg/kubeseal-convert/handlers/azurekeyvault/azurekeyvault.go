package azurekeyvault

import (
	"context"
	"encoding/json"
	"log"
	"os"

	identity "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	secrets "github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"
	"github.com/eladleev/kubeseal-convert/pkg/kubeseal-convert/interfaces"
)

func createClient() *secrets.Client {
	// see https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity#readme-defaultazurecredential: this allows getting credentials
	// via either environment variables, managed identity, or 'az login'
	cred, err := identity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("Failed to obtain a credential needed to login to the azure vault: %v", err)
	}

	vaultURI := os.Getenv("AZURE_KEY_VAULT_URI")
	client, err := secrets.NewClient(vaultURI, cred, nil)
	if err != nil {
		log.Fatalf("Failed to connect to vault '%s': %v", vaultURI, err)
	}
	return client
}

// retrieve secret by name with the client
func getSecret(client *secrets.Client, secretName string) map[string]interface{} {
	secret, err := client.GetSecret(context.TODO(), secretName, "", &secrets.GetSecretOptions{})
	if err != nil {
		log.Fatalf("Unable to read secret %s from the Azure Key Vault: %v", secretName, err)
	}

	mp := make(map[string]interface{})
	error := json.Unmarshal([]byte(*secret.Value), &mp)
	if error != nil {
		mp[secretName] = *secret.Value
	}
	return mp
}

type AzureKeyVaultImp struct {
}

func New() interfaces.AzureKeyVault {
	return &AzureKeyVaultImp{}
}

func (*AzureKeyVaultImp) GetSecret(secretName string) map[string]interface{} {
	cli := createClient()
	return getSecret(cli, secretName)
}
