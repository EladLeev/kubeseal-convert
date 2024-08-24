package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/eladleev/kubeseal-convert/pkg/kubeseal-convert/domain"
)

var azureKeyVaultCmd = &cobra.Command{
	Use:     "akv",
	Aliases: []string{"azurekeyvault", "azure-key-vault", "az", "azure"},
	Short:   "Convert Azure Key Vault secrets",
	Args:    cobra.ExactArgs(1),
	PreRun:  toggleDebug,
	Run: func(cmd *cobra.Command, args []string) {
		secretVal := domain.SecretValues{
			Data:        AzureKeyVault.GetSecrets(args[0]),
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
	rootCmd.AddCommand(azureKeyVaultCmd)
}
