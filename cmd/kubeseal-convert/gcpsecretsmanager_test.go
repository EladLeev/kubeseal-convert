package cmd

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"gotest.tools/assert"

	"github.com/eladleev/kubeseal-convert/mocks"
)

func TestGcpSecretsManagerCmd(t *testing.T) {
	rootCmd.AddCommand(gcpSecretsmanagerCmd)

	// mock kubeseal
	mockKubeSeal := mocks.NewKubeSeal(t)
	mockKubeSeal.On("BuildSecretFile", mock.AnythingOfType("domain.SecretValues"), mock.AnythingOfType("bool")).
		Return()
	KubeSeal = mockKubeSeal

	// mock gcpsecretsmanager
	mockGCP := mocks.NewGcpSecretsManager(t)
	mockGCP.On("GetSecret", mock.Anything, mock.AnythingOfType("int")).
		Return(map[string]interface{}{"key": "value"})
	GcpSecretsManager = mockGCP

	// test gcp command
	output, err := ExecuteCommand(
		rootCmd,
		"gcp",
		"projects/my-project/secrets/my-secret/versions/1",
		"--name",
		"blabla",
	)
	assert.NilError(t, err)
	assert.Equal(t, "", output)

	// Verify that the mocks were called as expected
	mockGCP.AssertExpectations(t)
	mockKubeSeal.AssertExpectations(t)
}
