package cmd

import "fmt"

// Logs implements the getDevContainerLogs custom driver command.
// It fetches container logs by calling: container logs <container-id>
// stdout and stderr are piped directly to the current process.
func Logs(containerID string) error {
	if containerID == "" {
		return fmt.Errorf("container ID is required")
	}

	return RunContainerCmdPassthrough("logs", containerID)
}
