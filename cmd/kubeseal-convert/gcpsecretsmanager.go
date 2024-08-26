package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/eladleev/kubeseal-convert/pkg/kubeseal-convert/domain"
)

var gcpSecretsmanagerCmd = &cobra.Command{
	Use:     "gcpsecretsmanager",
	Aliases: []string{"gs", "gcpsecretsmanager", "gcp", "gcp-secret"},
	Short:   "Convert GCP Secrets-Manager secrets",
	Args:    cobra.ExactArgs(1),
	PreRun:  toggleDebug,
	Run: func(cmd *cobra.Command, args []string) {
		secretVal := domain.SecretValues{
			Data:        GcpSecretsManager.GetSecret(args[0], timeout),
			Name:        ParseStringFlag(cmd, "name"),
			Namespace:   ParseStringFlag(cmd, "namespace"),
			Labels:      ParseLabels(cmd),
			Annotations: ParseAnnotations(cmd),
		}
		log.Debugf("secret values: %v", secretVal)
		KubeSeal.BuildSecretFile(secretVal, rawMode)
	},
}

func init() {
	rootCmd.AddCommand(gcpSecretsmanagerCmd)
}
