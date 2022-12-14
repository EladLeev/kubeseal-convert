package cmd

import (
	"fmt"
	"os"

	"github.com/eladleev/kubeseal-convert/pkg/kubeseal-convert/handlers/kubesealconvert"
	"github.com/eladleev/kubeseal-convert/pkg/kubeseal-convert/handlers/secretsmanager"
	"github.com/eladleev/kubeseal-convert/pkg/kubeseal-convert/handlers/vault"
	"github.com/spf13/cobra"
)

var (
	version         = "undefined"
	secretName      string
	secretNamespace string
	rootCmd         = &cobra.Command{
		Use:     "kubeseal-convert",
		Short:   "kubeseal-convert - a simple CLI to transform external secrets into Sealed Secrets",
		Long:    "kubeseal-convert is used to convert external secrets into Sealed Secrets objects, and help you adopt Sealed Secrets more easily.",
		Version: version,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	KubeSeal       = kubesealconvert.New()
	SecretsManager = secretsmanager.New()
	Vault          = vault.New()
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing the command.'%s'", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&secretName, "name", "n", "", "The Sealed Secret name (required)")
	rootCmd.PersistentFlags().StringVar(&secretNamespace, "namespace", "", "The Sealed Secret namespace. If not specified, taken from k8s context.")

	rootCmd.PersistentFlags().StringToStringP("annotations", "a", map[string]string{}, "Set k8s annotations")
	rootCmd.PersistentFlags().StringToStringP("labels", "l", map[string]string{}, "Set k8s lables")
}
