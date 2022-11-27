package kubesealconvert

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	internal "github.com/eladleev/kubeseal-convert/pkg/kubeseal-convert/internal"
	"github.com/spf13/cobra"
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

// Secretsmanager implements AWS SecretsManager logic to retrieve secrets and build the secretData struct
func Secretsmanager(secretName string, cmd *cobra.Command) (secretData internal.SecretValues) {
	cfg := createConfig()
	svc := secretsmanager.NewFromConfig(cfg)

	secretData.Name = ParseStringFlag(cmd, "name")           // ->
	secretData.Namespace = ParseStringFlag(cmd, "namespace") // ->
	secretData.Labels = ParseLabels(cmd)                     // ->
	secretData.Annotations = ParseAnnotations(cmd)           // ->
	secretData.Data = getSecret(svc, secretName)

	return secretData
}
