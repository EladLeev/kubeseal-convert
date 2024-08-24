package gcpsecretsmanager

import (
	"context"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"golang.org/x/oauth2/google"

	"github.com/eladleev/kubeseal-convert/pkg/kubeseal-convert/interfaces"
)

var cleanSecretName string

// buildSecretId will build the full secret ID based on the secret name.
// If a full secret ID is provided, it will be returned as is, if not,
// then, we need to extract the project ID from the default credentials
func buildSecretId(ctx context.Context, secretName string) string {
	secretSlice := strings.Split(secretName, "/")
	log.Debugf("secretSlice: %v", secretSlice)

	// The only supported format is: projects/<PROJECT_ID>/secrets/<SECRET_NAME>/versions/<VERSION>
	if len(secretSlice) != 6 {
		credentials, err := google.FindDefaultCredentials(ctx)
		if err != nil {
			log.Fatalf("failed to FindDefaultCredentials: %v", err)
		}
		cleanSecretName = secretName
		return fmt.Sprintf("projects/%v/secrets/%v/versions/%v", credentials.ProjectID, cleanSecretName, "latest")
	}

	cleanSecretName = secretSlice[3]
	return secretName
}

func getSecret(secretName string) map[string]interface{} {
	mp := make(map[string]interface{})
	ctx := context.Background()

	client, err := secretmanager.NewClient(ctx)
	log.Debugf("client: %v", client)
	if err != nil {
		log.Fatalf("failed to setup client: %v", err)
	}

	defer func(client *secretmanager.Client) {
		_ = client.Close()
	}(client)

	accessRequest := &secretmanagerpb.AccessSecretVersionRequest{
		Name: buildSecretId(ctx, secretName),
	}
	log.Debugf("accessRequest: %v", accessRequest)

	result, err := client.AccessSecretVersion(ctx, accessRequest)
	log.Debugf("result: %v", result)
	if err != nil {
		log.Fatalf("failed to access secret version: %v", err)
	}

	mp[cleanSecretName] = string(result.Payload.Data)
	return mp
}

type GcpSecretsManagerImp struct {
}

func New() interfaces.SecretsManager {
	return &GcpSecretsManagerImp{}
}

func (*GcpSecretsManagerImp) GetSecret(secretName string) map[string]interface{} {
	return getSecret(secretName)
}
