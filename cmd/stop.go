package cmd

import "fmt"

// Stop implements the stopDevContainer custom driver command.
// It stops a running container by calling: container stop <container-id>
func Stop(containerID string) error {
	if containerID == "" {
		return fmt.Errorf("container ID is required")
	}

	Logf("Stopping container: %s", containerID)
	return RunContainerCmdWithStdoutStderr("stop", containerID)
}
