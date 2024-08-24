package cmd

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"gotest.tools/assert"

	"github.com/eladleev/kubeseal-convert/mocks"
)

func TestAzureKeyVaultCmd(t *testing.T) {
	rootCmd.AddCommand(azureKeyVaultCmd)

	// mock kubeseal
	mockKubeSeal := mocks.NewKubeSeal(t)
	mockKubeSeal.On("BuildSecretFile", mock.AnythingOfType("domain.SecretValues"), mock.AnythingOfType("bool")).Return()
	KubeSeal = mockKubeSeal

	// mock azurekeyvault
	mockAzureKeyVault := mocks.NewAzureKeyVault(t)
	mockAzureKeyVault.On("GetSecrets", mock.Anything).Return(map[string]interface{}{
		"key": "value",
	}, nil)
	AzureKeyVault = mockAzureKeyVault

	// test az command
	output, err := ExecuteCommand(rootCmd, "az", "dev/secret", "--name", "blabla")
	assert.NilError(t, err)
	assert.Equal(t, "", output)

	// Verify that the mocks were called as expected
	mockAzureKeyVault.AssertExpectations(t)
	mockKubeSeal.AssertExpectations(t)
}
