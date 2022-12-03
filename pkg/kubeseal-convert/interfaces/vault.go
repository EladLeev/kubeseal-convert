package interfaces

type Vault interface {
	GetSecret(secretName string) map[string]interface{}
}
