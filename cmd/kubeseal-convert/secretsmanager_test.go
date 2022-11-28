package cmd

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/mock"

	"github.com/eladleev/kubeseal-convert/mocks"
	"gotest.tools/assert"
)

func executeCommand(root *cobra.Command, args ...string) (output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)

	err = root.Execute()
	if err != nil {
		fmt.Println(err)
	}

	return buf.String(), err
}

func TestGetResult(t *testing.T) {
	rootCmd.AddCommand(secretsmanagerCmd)

	// mock kubeseal
	mockKubeSeal := mocks.NewKubeSeal(t)
	mockKubeSeal.On("BuildSecretFile", mock.Anything).Return()
	KubeSeal = mockKubeSeal

	// mock secretsmanager
	mockSecretsManager := mocks.NewSecretsManager(t)
	mockSecretsManager.On("GetSecret", mock.Anything).Return(nil)
	SecretsManager = mockSecretsManager
	// test command
	output, _ := executeCommand(rootCmd, "sm", "dev/secret", "--name", "blabla")
	assert.Equal(t, "", output)
}
