package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

// ContainerBinaryPath returns the path to the Apple Container CLI binary.
// It checks the CONTAINER_PATH env var first, then defaults to /usr/local/bin/container.
func ContainerBinaryPath() string {
	if p := os.Getenv("CONTAINER_PATH"); p != "" {
		return p
	}
	return "/usr/local/bin/container"
}

// RunContainerCmd executes the container CLI with the given arguments and returns stdout.
func RunContainerCmd(args ...string) ([]byte, error) {
	cmd := exec.Command(ContainerBinaryPath(), args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("container %v failed: %w\nstderr: %s", args, err, stderr.String())
	}
	return stdout.Bytes(), nil
}

// RunContainerCmdPassthrough executes the container CLI with stdin/stdout/stderr
// connected directly to the current process. This is critical for commands like
// `exec` where DevPod sends binary data through these streams.
func RunContainerCmdPassthrough(args ...string) error {
	cmd := exec.Command(ContainerBinaryPath(), args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// RunContainerCmdWithStdoutStderr executes the container CLI and pipes stdout/stderr
// to the current process stdout/stderr. Used for commands where DevPod captures output.
func RunContainerCmdWithStdoutStderr(args ...string) error {
	cmd := exec.Command(ContainerBinaryPath(), args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Logf prints a debug message to stderr (DevPod captures stderr as log lines).
func Logf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
}
