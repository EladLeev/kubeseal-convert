// nolint
package kubesealconvert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"

	log "github.com/sirupsen/logrus"

	"github.com/eladleev/kubeseal-convert/pkg/kubeseal-convert/domain"
	"github.com/eladleev/kubeseal-convert/pkg/kubeseal-convert/interfaces"
)

var noKubesealHelper = "Follow: https://github.com/bitnami-labs/sealed-secrets?tab=readme-ov-file#kubeseal to install the tool."

type KubesealImpl struct {
}

// checkKubesealBinary verifies if kubeseal is installed and returns its path
func checkKubesealBinary() string {
	kubesealBinary, err := exec.LookPath("kubeseal")
	if err != nil {
		log.Fatalf("kubeseal command not found: %v\n%v", err, noKubesealHelper)
	}
	log.Debugf("kubesealBinary: %v", kubesealBinary)
	return kubesealBinary
}

// runCommandWithInput sets up and runs a command with the provided input, and returns the output.
func runCommandWithInput(cmdPath string, cmdArgs []string, input string) (string, error) {
	cmd := exec.Command(cmdPath, cmdArgs...)
	log.Debugf("cmd: %v, cmdArgs: %v", cmd, cmdArgs)

	var outputBuffer, errorBuffer bytes.Buffer
	cmd.Stdout = &outputBuffer
	cmd.Stderr = &errorBuffer

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return "", fmt.Errorf("failed to get stdin pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return "", fmt.Errorf("failed to start command: %w", err)
	}

	_, err = io.WriteString(stdin, input)
	if err != nil {
		return "", fmt.Errorf("failed to write to stdin: %w", err)
	}
	stdin.Close()

	if err := cmd.Wait(); err != nil {
		return "", fmt.Errorf("command failed: %w\nStderr: %s", err, errorBuffer.String())
	}

	return outputBuffer.String(), nil
}

func New() interfaces.KubeSeal {
	return &KubesealImpl{}
}

// Seal gets a raw k8s secret as a string, and run the kubeseal command
func (*KubesealImpl) Seal(secret string) {
	kubesealBinary := checkKubesealBinary()
	args := []string{"--format", "yaml"}

	output, err := runCommandWithInput(kubesealBinary, args, secret)
	if err != nil {
		log.Fatalf("Sealing failed: %v", err)
	}

	fmt.Print(output)
}

// RawSeal gets a raw k8s secret as a string, and run the kubeseal command, using the RAW mode
func (*KubesealImpl) RawSeal(secretValues domain.SecretValues) {
	log.Debug("using raw mode")

	kubesealBinary := checkKubesealBinary()

	secretData, err := json.Marshal(secretValues.Data)
	if err != nil {
		log.Fatalf("Failed to marshal secret data: %v", err)
	}

	args := []string{"--raw"}
	if secretValues.Namespace != "" {
		args = append(args, "--namespace", secretValues.Namespace)
	}
	if secretValues.Name != "" {
		args = append(args, "--name", secretValues.Name)
	}

	log.Debugf("kubesealBinary: %v, args: %v", kubesealBinary, args)
	output, err := runCommandWithInput(kubesealBinary, args, string(secretData))
	if err != nil {
		log.Fatalf("Raw sealing failed: %v", err)
	}

	fmt.Print(output)
}

// BuildSecretFile generates a Sealed Secrets
func (impl *KubesealImpl) BuildSecretFile(secretValues domain.SecretValues, useRaw bool) {
	if useRaw {
		impl.RawSeal(secretValues)
		return
	}

	rawSecret := buildKubeSecret(secretValues)
	output, err := json.Marshal(&rawSecret)
	if err != nil {
		log.Fatalf("Unable to marshal secret: %v", err)
	}

	log.Debugf("rawSecret: %v, outputSecret: %v", rawSecret, output)
	impl.Seal(string(output))
}
