package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/eladleev/kubeseal-convert/pkg/kubeseal-convert/domain"
)

var secretsmanagerCmd = &cobra.Command{
	Use:     "secretsmanager",
	Aliases: []string{"sm", "secretsmanager", "aws", "aws-secret"},
	Short:   "Convert AWS Secrets-Manager secrets",
	Args:    cobra.ExactArgs(1),
	PreRun:  toggleDebug,
	Run: func(cmd *cobra.Command, args []string) {
		secretVal := domain.SecretValues{
			Data:        SecretsManager.GetSecret(args[0]),
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
	rootCmd.AddCommand(secretsmanagerCmd)
}
