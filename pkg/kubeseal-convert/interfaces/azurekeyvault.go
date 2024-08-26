package interfaces

type AzureKeyVault interface {
	GetSecrets(vaultName string, timeout int) map[string]interface{}
}
