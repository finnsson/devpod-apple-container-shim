package cmd

import "fmt"

// Start implements the startDevContainer custom driver command.
// It starts a stopped container by calling: container start <container-id>
func Start(containerID string) error {
	if containerID == "" {
		return fmt.Errorf("container ID is required")
	}

	Logf("Starting container: %s", containerID)
	return RunContainerCmdWithStdoutStderr("start", containerID)
}
