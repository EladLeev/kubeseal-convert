package cmd

import (
	"testing"

	"github.com/stretchr/testify/mock"

	"gotest.tools/assert"

	"github.com/eladleev/kubeseal-convert/mocks"
)

func TestSecretManagerCmd(t *testing.T) {
	rootCmd.AddCommand(secretsmanagerCmd)

	// mock kubeseal
	mockKubeSeal := mocks.NewKubeSeal(t)
	mockKubeSeal.On("BuildSecretFile", mock.Anything, mock.AnythingOfType("bool")).Return()
	KubeSeal = mockKubeSeal

	// mock secretsmanager
	mockSecretsManager := mocks.NewSecretsManager(t)
	mockSecretsManager.On("GetSecret", mock.Anything).Return(nil)
	SecretsManager = mockSecretsManager
	// test sm command
	output, _ := ExecuteCommand(rootCmd, "sm", "dev/secret", "--name", "blabla")
	assert.Equal(t, "", output)
}
