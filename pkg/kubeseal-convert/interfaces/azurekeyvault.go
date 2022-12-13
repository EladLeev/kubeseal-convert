package interfaces

type AzureKeyVault interface {
	GetSecret(secretName string) map[string]interface{}
}
