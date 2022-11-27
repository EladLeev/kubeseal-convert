package kubesealconvert

import (
	"bytes"
	"io"
	"log"
	"os"
	"os/exec"
)

var cmdBuffer bytes.Buffer

// Kubeseal gets a raw k8s secret as a string, and run the kubeseal command
func Kubeseal(secret string) {
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

	r, w := io.Pipe()
	echoCmd.Stdout = w
	cmdExec.Stdin = r
	cmdExec.Stdout = &cmdBuffer

	echoCmd.Start()
	cmdExec.Start()
	echoCmd.Wait()
	w.Close()
	cmdExec.Wait()
	io.Copy(os.Stdout, &cmdBuffer)

}
