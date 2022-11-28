package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// ParseStringFlag gets a cobra context, and parses a given arg from cmd line arguments.
// To be used by each one of the secret providers implementations.
func ParseStringFlag(cmd *cobra.Command, flagName string) string {
	name, e := cmd.Flags().GetString(flagName)
	if e != nil {
		fmt.Printf("Unable to parse %v flag: %v", flagName, e)
	}

	return name
}

// ParseAnnotations gets a cobra context, and parses a list of annotations from cmd line arguments.
// To be used by each one of the secret providers implementations.
func ParseAnnotations(cmd *cobra.Command) map[string]string {
	annotations, e := cmd.Flags().GetStringToString("annotations")
	if e != nil {
		fmt.Printf("Unable to parse annoations: %v", e)
	}

	return annotations
}

// ParseLabels gets a cobra context, and parses a list of labels from cmd line arguments.
// To be used by each one of the secret providers implementations.
func ParseLabels(cmd *cobra.Command) map[string]string {
	labels, e := cmd.Flags().GetStringToString("labels")
	if e != nil {
		fmt.Printf("Unable to parse annoations: %v", e)
	}

	return labels
}
