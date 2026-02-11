package cmd

import "fmt"

// Delete implements the deleteDevContainer custom driver command.
// It deletes a container by calling: container delete --force <container-id>
// The --force flag ensures we can delete even running containers.
func Delete(containerID string) error {
	if containerID == "" {
		return fmt.Errorf("container ID is required")
	}

	Logf("Deleting container: %s", containerID)
	return RunContainerCmdWithStdoutStderr("delete", "--force", containerID)
}
