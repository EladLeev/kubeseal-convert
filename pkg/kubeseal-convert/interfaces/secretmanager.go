package interfaces

type SecretsManager interface {
	GetSecret(secretName string, timeout int) map[string]interface{}
}
