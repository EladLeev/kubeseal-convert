package cmd

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"gotest.tools/assert"

	"github.com/eladleev/kubeseal-convert/mocks"
)

func TestVaultCmd(t *testing.T) {
	rootCmd.AddCommand(vaultCmd)

	// mock kubeseal
	mockKubeSeal := mocks.NewKubeSeal(t)
	mockKubeSeal.On("BuildSecretFile", mock.AnythingOfType("domain.SecretValues"), mock.AnythingOfType("bool")).Return()
	KubeSeal = mockKubeSeal

	// mock vault
	mockVault := mocks.NewVault(t)
	mockVault.On("GetSecret", mock.Anything).Return(map[string]interface{}{
		"key": "value",
	}, nil)
	Vault = mockVault

	// test vault command
	output, err := ExecuteCommand(rootCmd, "vault", "dev/secret", "--name", "blabla")
	assert.NilError(t, err)
	assert.Equal(t, "", output)

	// Verify that the mocks were called as expected
	mockVault.AssertExpectations(t)
	mockKubeSeal.AssertExpectations(t)
}
