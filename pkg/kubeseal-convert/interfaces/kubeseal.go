package interfaces

import "github.com/eladleev/kubeseal-convert/pkg/kubeseal-convert/domain"

type KubeSeal interface {
	Seal(secret string)
	RawSeal(secretValues domain.SecretValues)
	BuildSecretFile(secretValues domain.SecretValues, useRaw bool)
}
