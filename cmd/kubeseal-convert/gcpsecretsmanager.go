package cmd

import (
	"github.com/spf13/cobra"

	"github.com/eladleev/kubeseal-convert/pkg/kubeseal-convert/domain"
)

var gcpSecretsmanagerCmd = &cobra.Command{
	Use:     "gcpsecretsmanager",
	Aliases: []string{"gs", "gcpsecretsmanager", "gcp", "gcp-secret"},
	Short:   "Convert GCP Secrets-Manager secrets",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		secretVal := domain.SecretValues{
			Data:        GcpSecretsManager.GetSecret(args[0]),
			Name:        ParseStringFlag(cmd, "name"),
			Namespace:   ParseStringFlag(cmd, "namespace"),
			Labels:      ParseLabels(cmd),
			Annotations: ParseAnnotations(cmd),
		}
		KubeSeal.BuildSecretFile(secretVal)
	},
}

func init() {
	rootCmd.AddCommand(gcpSecretsmanagerCmd)
}
