package cmd

import (
	"github.com/eladleev/kubeseal-convert/pkg/kubeseal-convert/domain"
	"github.com/spf13/cobra"
)

var azureKeyVaultCmd = &cobra.Command{
	Use:     "az",
	Aliases: []string{"azurekeyvault", "azure-key-vault", "azure"},
	Short:   "Convert Azure Key Vault secrets",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		secretVal := domain.SecretValues{
			Data:        AzureKeyVault.GetSecret(args[0]),
			Name:        ParseStringFlag(cmd, "name"),
			Namespace:   ParseStringFlag(cmd, "namespace"),
			Labels:      ParseLabels(cmd),
			Annotations: ParseAnnotations(cmd),
		}
		KubeSeal.BuildSecretFile(secretVal)
	},
}

func init() {
	rootCmd.AddCommand(azureKeyVaultCmd)
}
