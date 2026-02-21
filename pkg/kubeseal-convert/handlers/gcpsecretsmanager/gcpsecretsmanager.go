package gcpsecretsmanager

import (
	"context"
	"fmt"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"github.com/googleapis/gax-go/v2"
	"golang.org/x/oauth2/google"

	"github.com/eladleev/kubeseal-convert/pkg/kubeseal-convert/interfaces"
)

// gcpSecretsClientAPI is the subset of secretmanager.Client used by this handler.
type gcpSecretsClientAPI interface {
	AccessSecretVersion(
		ctx context.Context,
		req *secretmanagerpb.AccessSecretVersionRequest,
		opts ...gax.CallOption,
	) (*secretmanagerpb.AccessSecretVersionResponse, error)
	Close() error
}

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
		return fmt.Sprintf(
			"projects/%v/secrets/%v/versions/%v",
			credentials.ProjectID,
			cleanSecretName,
			"latest",
		)
	}

	cleanSecretName = secretSlice[3]
	return secretName
}

func getSecret(
	ctx context.Context,
	client gcpSecretsClientAPI,
	secretName string,
) (map[string]interface{}, error) {
	mp := make(map[string]interface{})

	accessRequest := &secretmanagerpb.AccessSecretVersionRequest{
		Name: buildSecretId(ctx, secretName),
	}
	log.Debugf("accessRequest: %v", accessRequest)

	result, err := client.AccessSecretVersion(ctx, accessRequest)
	log.Debugf("result: %v", result)
	if err != nil {
		return nil, fmt.Errorf("failed to access secret version: %v", err)
	}

	mp[cleanSecretName] = string(result.Payload.Data)
	return mp, nil
}

type GcpSecretsManagerImp struct {
	client gcpSecretsClientAPI
}

func New() interfaces.SecretsManager {
	return &GcpSecretsManagerImp{}
}

func (g *GcpSecretsManagerImp) GetSecret(secretName string, timeout int) map[string]interface{} {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	cli := g.client
	if cli == nil {
		var err error
		cli, err = secretmanager.NewClient(ctx)
		log.Debugf("client: %v", cli)
		if err != nil {
			log.Errorf("failed to setup client: %v", err)
			return nil
		}
		defer func() { _ = cli.Close() }()
	}

	secret, err := getSecret(ctx, cli, secretName)
	if err != nil {
		log.Errorf("failed to get secret: %v", err)
	}
	return secret
}
