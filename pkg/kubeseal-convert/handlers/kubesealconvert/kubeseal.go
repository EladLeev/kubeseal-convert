// nolint
package kubesealconvert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/eladleev/kubeseal-convert/pkg/kubeseal-convert/domain"
	"github.com/eladleev/kubeseal-convert/pkg/kubeseal-convert/interfaces"
)

type KubesealImpl struct {
}

func New() interfaces.KubeSeal {
	return &KubesealImpl{}
}

// Seal gets a raw k8s secret as a string, and run the kubeseal command
func (*KubesealImpl) Seal(secret string) {
	kubesealBinary, e := exec.LookPath("kubeseal")
	if e != nil {
		log.Fatalf("Unable to find kubeseal command: %v", e)
	}

	// assume that echo command is always available
	echoBin, _ := exec.LookPath("echo")
	echoCmd := &exec.Cmd{
		Path: echoBin,
		Args: []string{echoBin, secret},
	}

	cmdExec := &exec.Cmd{
		Path:   kubesealBinary,
		Args:   []string{kubesealBinary, "--format", "yaml"},
		Stdout: os.Stdout,
		Stderr: os.Stdout,
	}
	var cmdBuffer bytes.Buffer
	r, w := io.Pipe()
	echoCmd.Stdout = w
	cmdExec.Stdin = r
	cmdExec.Stdout = &cmdBuffer

	echoCmd.Start()
	cmdExec.Start()
	echoCmd.Wait()
	w.Close()
	cmdExec.Wait()

	_, err := io.Copy(os.Stdout, &cmdBuffer)
	if err != nil {
		log.Fatal(err)
	}
}

// BuildSecretFile generates a Sealed Secrets
func (impl *KubesealImpl) BuildSecretFile(secretValues domain.SecretValues) {
	rawSecret := buildSecret(secretValues)
	output, e := json.Marshal(&rawSecret)
	if e != nil {
		fmt.Printf("Unable to marshal secret: %v", e)
	}
	impl.Seal(string(output))
}
