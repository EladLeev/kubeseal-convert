package domain

// SecretValues struct defines the basic structure of a secret.
// Name => The secret name
// Namespace => The namespace name
// Data => The raw data of the secret
// Annotations => The k8s annotations provided by the user to be included
// Labels => The k8s labels provided by the user to be included
type SecretValues struct {
	Name        string
	Namespace   string
	Annotations map[string]string
	Labels      map[string]string
	Data        map[string]interface{}
}
