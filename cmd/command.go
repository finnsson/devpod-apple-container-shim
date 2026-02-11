package cmd

import (
	"fmt"
	"os"
	"os/exec"
)

// Command implements the commandDevContainer custom driver command.
// It executes a command inside a running container by calling:
//
//	container exec -i -u $DEVCONTAINER_USER <container-id> sh -c "$DEVCONTAINER_COMMAND"
//
// This is the MOST CRITICAL command — DevPod pipes binary data (compressed tar,
// agent binaries, SSH tunnel data) through stdin/stdout. The streams MUST be
// faithfully proxied without any buffering or modification.
func Command(containerID string) error {
	if containerID == "" {
		return fmt.Errorf("container ID is required")
	}

	user := os.Getenv("DEVCONTAINER_USER")
	command := os.Getenv("DEVCONTAINER_COMMAND")

	if command == "" {
		return fmt.Errorf("DEVCONTAINER_COMMAND environment variable is required")
	}

	args := []string{"exec"}

	// Always use -i to keep stdin open (DevPod sends data through it)
	args = append(args, "-i")

	// Set user if specified
	if user != "" {
		args = append(args, "-u", user)
	}

	// Container ID
	args = append(args, containerID)

	// The command to execute: sh -c "<command>"
	args = append(args, "sh", "-c", command)

	// Execute with full stdin/stdout/stderr passthrough
	cmd := exec.Command(ContainerBinaryPath(), args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		return fmt.Errorf("exec failed: %w", err)
	}

	return nil
}
