package cmd

import (
	"github.com/eladleev/kubeseal-convert/pkg/kubeseal-convert/domain"
	"github.com/spf13/cobra"
)

var vaultCmd = &cobra.Command{
	Use:     "vault",
	Aliases: []string{"vlt", "vault"},
	Short:   "Convert Hashicorp Vault secrets",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		secretVal := domain.SecretValues{
			Data:        Vault.GetSecret(args[0]),
			Name:        ParseStringFlag(cmd, "name"),
			Namespace:   ParseStringFlag(cmd, "namespace"),
			Labels:      ParseLabels(cmd),
			Annotations: ParseAnnotations(cmd),
		}
		KubeSeal.BuildSecretFile(secretVal)
	},
}

func init() {
	rootCmd.AddCommand(vaultCmd)
}
