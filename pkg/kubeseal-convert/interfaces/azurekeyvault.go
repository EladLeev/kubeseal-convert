package interfaces

type AzureKeyVault interface {
	GetSecrets(vaultName string) map[string]interface{}
}
