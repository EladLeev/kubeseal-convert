package secretsmanager

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/eladleev/kubeseal-convert/pkg/kubeseal-convert/interfaces"
)

// TODO: Implement proper context
func createConfig() aws.Config {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("Unable to load SDK config, %v", err)
	}
	return cfg
}

// getSecret wil get the secret into a map[string]interface{} as the return value may vary
func getSecret(svc *secretsmanager.Client, secretName string) map[string]interface{} {
	r, err := svc.GetSecretValue(context.TODO(), &secretsmanager.GetSecretValueInput{SecretId: &secretName})
	if err != nil {
		log.Fatal(err.Error())
	}

	mp := make(map[string]interface{})
	error := json.Unmarshal([]byte(*r.SecretString), &mp)
	if error != nil {
		log.Fatalf("Unable to parse secret with err: %v", error)
	}
	return mp
}

type SecretsManagerImp struct {
}

func New() interfaces.SecretsManager {
	return &SecretsManagerImp{}
}

func (*SecretsManagerImp) GetSecret(secretName string) map[string]interface{} {
	cfg := createConfig()
	svc := secretsmanager.NewFromConfig(cfg)
	return getSecret(svc, secretName)
}
