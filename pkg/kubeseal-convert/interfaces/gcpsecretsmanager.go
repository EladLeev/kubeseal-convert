package interfaces

type GcpSecretsManager interface {
	GetSecret(secretName string, timeout int) map[string]interface{}
}
