package cmd

import (
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/eladleev/kubeseal-convert/mocks"
	"gotest.tools/assert"
)

func TestSecretManagerCmd(t *testing.T) {
	rootCmd.AddCommand(secretsmanagerCmd)

	// mock kubeseal
	mockKubeSeal := mocks.NewKubeSeal(t)
	mockKubeSeal.On("BuildSecretFile", mock.Anything).Return()
	KubeSeal = mockKubeSeal

	// mock secretsmanager
	mockSecretsManager := mocks.NewSecretsManager(t)
	mockSecretsManager.On("GetSecret", mock.Anything).Return(nil)
	SecretsManager = mockSecretsManager
	//test sm command
	output, _ := ExecuteCommand(rootCmd, "sm", "dev/secret", "--name", "blabla")
	assert.Equal(t, "", output)
}
