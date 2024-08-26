package interfaces

type Vault interface {
	GetSecret(secretName string, timeout int) map[string]interface{}
}
