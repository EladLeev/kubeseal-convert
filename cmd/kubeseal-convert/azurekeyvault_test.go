package cmd

import (
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/eladleev/kubeseal-convert/mocks"
	"gotest.tools/assert"
)

func TestAzureKeyVaultCmd(t *testing.T) {
	rootCmd.AddCommand(azureKeyVaultCmd)

	// mock kubeseal
	mockKubeSeal := mocks.NewKubeSeal(t)
	mockKubeSeal.On("BuildSecretFile", mock.Anything).Return()
	KubeSeal = mockKubeSeal

	// mock azurekeyvault
	mockAzureKeyVault := mocks.NewAzureKeyVault(t)
	mockAzureKeyVault.On("GetSecrets", mock.Anything).Return(map[string]interface{}{
		"key": "value",
	})
	AzureKeyVault = mockAzureKeyVault

	//test az command
	output, _ := ExecuteCommand(rootCmd, "az", "dev/secret", "--name", "blabla")
	assert.Equal(t, "", output)
}
