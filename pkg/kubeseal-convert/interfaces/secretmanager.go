package interfaces

type SecretsManager interface {
	GetSecret(secretName string) map[string]interface{}
}
