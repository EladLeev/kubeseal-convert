package kubesealconvert

import (
	"fmt"

	coreV1 "k8s.io/api/core/v1"

	"github.com/eladleev/kubeseal-convert/pkg/kubeseal-convert/domain"
)

// buildDataBytes merges the SecretValues.Data into the new k8s raw secret
func buildDataBytes(sv domain.SecretValues) map[string][]byte {
	data := make(map[string][]byte)

	for k, v := range sv.Data {
		str, ok := v.(string)
		if !ok {
			fmt.Println("Oops, unexpected field value. Unable to decode secret value to string.")
		}
		data[k] = []byte(str)
	}
	return data
}

// buildKubeSecret gets the SecretValues struct and returns a raw k8s secret
func buildKubeSecret(sv domain.SecretValues) coreV1.Secret {
	var secretSpec coreV1.Secret

	secretSpec.Name = sv.Name
	secretSpec.Namespace = sv.Namespace
	secretSpec.APIVersion = "v1"
	secretSpec.Kind = "Secret"
	secretSpec.Type = "Opaque"
	secretSpec.Data = buildDataBytes(sv)
	secretSpec.Annotations = sv.Annotations
	secretSpec.Labels = sv.Labels

	return secretSpec
}
