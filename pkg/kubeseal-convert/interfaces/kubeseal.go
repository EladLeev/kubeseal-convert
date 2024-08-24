package interfaces

import "github.com/eladleev/kubeseal-convert/pkg/kubeseal-convert/domain"

type KubeSeal interface {
	Seal(secret string)
	BuildSecretFile(secretValues domain.SecretValues, useRaw bool)
}
