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

// secretsManagerAPI is the subset of secretsmanager.Client used by this handler.
type secretsManagerAPI interface {
	GetSecretValue(
		ctx context.Context,
		input *secretsmanager.GetSecretValueInput,
		optFns ...func(*secretsmanager.Options),
	) (*secretsmanager.GetSecretValueOutput, error)
}

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
	svc secretsManagerAPI,
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
	cfg    aws.Config
	client secretsManagerAPI
}

func New() interfaces.SecretsManager {
	return &SecretsManagerImp{cfg: createConfig()}
}

func (s *SecretsManagerImp) GetSecret(secretName string, timeout int) map[string]interface{} {
	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	svc := s.client
	if svc == nil {
		svc = secretsmanager.NewFromConfig(s.cfg)
	}
	secret, err := getSecret(ctx, svc, secretName)
	if err != nil {
		log.Errorf("failed to get secret: %v", err)
	}
	return secret
}
