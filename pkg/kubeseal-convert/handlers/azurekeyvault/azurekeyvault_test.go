package azurekeyvault

import (
	"os"
	"testing"
)

func Test_getSecret(t *testing.T) {
	os.Setenv("AZURE_KEY_VAULT_URI", "https://example.vault.azure.net")

	// TODO

}
