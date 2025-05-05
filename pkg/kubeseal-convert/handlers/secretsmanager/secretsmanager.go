package secretsmanager

import (
	"context"
	"encoding/json"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	log "github.com/sirupsen/logrus"

	"github.com/eladleev/kubeseal-convert/pkg/kubeseal-convert/interfaces"
)

func createConfig() aws.Config {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Errorf("unable to load SDK config, %v", err)
	}
	return cfg
}

// getSecret wil get the secret into a map[string]interface{} as the return value may vary
func getSecret(
	ctx context.Context,
	svc *secretsmanager.Client,
	secretName string,
) (map[string]interface{}, error) {
	r, err := svc.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{SecretId: &secretName})
	if err != nil {
		return nil, err
	}

	mp := make(map[string]interface{})
	err = json.Unmarshal([]byte(*r.SecretString), &mp)
	if err != nil {
		return nil, err
	}
	return mp, nil
}

type SecretsManagerImp struct {
	cfg aws.Config
}

func New() interfaces.SecretsManager {
	return &SecretsManagerImp{cfg: createConfig()}
}

func (s *SecretsManagerImp) GetSecret(secretName string, timeout int) map[string]interface{} {
	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	svc := secretsmanager.NewFromConfig(s.cfg)
	secret, err := getSecret(ctx, svc, secretName)
	if err != nil {
		log.Errorf("failed to get secret: %v", err)
	}
	return secret
}
