package cmd

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"gotest.tools/assert"

	"github.com/eladleev/kubeseal-convert/mocks"
)

func TestVaultCmd(t *testing.T) {
	rootCmd.AddCommand(secretsmanagerCmd)

	// mock kubeseal
	mockKubeSeal := mocks.NewKubeSeal(t)
	mockKubeSeal.On("BuildSecretFile", mock.Anything, mock.AnythingOfType("bool")).Return()
	KubeSeal = mockKubeSeal

	// mock vault
	mockVault := mocks.NewVault(t)
	mockVault.On("GetSecret", mock.Anything).Return(map[string]interface{}{
		"key": "value",
	})
	Vault = mockVault

	// test vault command
	output, _ := ExecuteCommand(rootCmd, "vault", "dev/secret", "--name", "blabla")
	assert.Equal(t, "", output)
}
