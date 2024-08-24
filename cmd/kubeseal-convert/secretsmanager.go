package cmd

import (
	"github.com/spf13/cobra"

	"github.com/eladleev/kubeseal-convert/pkg/kubeseal-convert/domain"
)

var secretsmanagerCmd = &cobra.Command{
	Use:     "secretsmanager",
	Aliases: []string{"sm", "secretsmanager", "aws", "aws-secret"},
	Short:   "Convert AWS Secrets-Manager secrets",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		secretVal := domain.SecretValues{
			Data:        SecretsManager.GetSecret(args[0]),
			Name:        ParseStringFlag(cmd, "name"),
			Namespace:   ParseStringFlag(cmd, "namespace"),
			Labels:      ParseLabels(cmd),
			Annotations: ParseAnnotations(cmd),
		}
		KubeSeal.BuildSecretFile(secretVal)
	},
}

func init() {
	rootCmd.AddCommand(secretsmanagerCmd)
}
