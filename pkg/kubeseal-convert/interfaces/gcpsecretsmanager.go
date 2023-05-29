package interfaces

type GcpSecretsManager interface {
	GetSecret(secretName string) map[string]interface{}
}
