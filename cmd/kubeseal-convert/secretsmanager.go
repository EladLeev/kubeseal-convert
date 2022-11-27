package cmd

import (
	kubesealconvert "github.com/eladleev/kubeseal-convert/pkg/kubeseal-convert"
	"github.com/spf13/cobra"
)

var secretsmanagerCmd = &cobra.Command{
	Use:     "secretsmanager",
	Aliases: []string{"sm", "secretmanager", "aws", "aws-secret"},
	Short:   "Convert AWS Secrets-Manager secrets",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		secretData := kubesealconvert.Secretsmanager(args[0], cmd)

		// Build Sealed Secrets using AWS Secrets Manager data
		kubesealconvert.BuildSecretFile(secretData)

	},
}

func init() {
	rootCmd.AddCommand(secretsmanagerCmd)
}
