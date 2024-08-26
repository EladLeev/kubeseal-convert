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
	mockKubeSeal.On("BuildSecretFile", mock.AnythingOfType("domain.SecretValues"), mock.AnythingOfType("bool")).Return()
	KubeSeal = mockKubeSeal

	// mock secretsmanager
	mockSecretsManager := mocks.NewSecretsManager(t)

	mockSecretsManager.On("GetSecret", mock.Anything, mock.AnythingOfType("int")).Return(map[string]interface{}{"key": "value"}, nil)
	SecretsManager = mockSecretsManager

	// test sm command
	output, err := ExecuteCommand(rootCmd, "sm", "dev/secret", "--name", "blabla")
	assert.NilError(t, err)
	assert.Equal(t, "", output)

	// Verify that the mocks were called as expected
	mockSecretsManager.AssertExpectations(t)
	mockKubeSeal.AssertExpectations(t)
}
