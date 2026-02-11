package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// Run implements the runDevContainer custom driver command.
// It reads DEVCONTAINER_RUN_OPTIONS (JSON) from the environment, parses it into
// RunOptions, and translates it to an Apple Container `container run` command.
//
// Key translations:
//   - --name <containerID>: use the workspace ID as the container name
//   - -d: detach (always, DevPod expects the container to run in background)
//   - -e KEY=VALUE: environment variables
//   - -l KEY=VALUE: labels
//   - -u USER: user
//   - --entrypoint CMD: entrypoint override
//   - --mount type=<>,source=<>,target=<>: mounts
//   - Unsupported Docker flags (capAdd, securityOpt, privileged) are silently ignored
func Run(containerID string) error {
	if containerID == "" {
		return fmt.Errorf("container ID is required")
	}

	runOptionsJSON := os.Getenv("DEVCONTAINER_RUN_OPTIONS")
	if runOptionsJSON == "" {
		return fmt.Errorf("DEVCONTAINER_RUN_OPTIONS environment variable is required")
	}

	var opts RunOptions
	if err := json.Unmarshal([]byte(runOptionsJSON), &opts); err != nil {
		return fmt.Errorf("failed to parse DEVCONTAINER_RUN_OPTIONS: %w", err)
	}

	if opts.Image == "" {
		return fmt.Errorf("image is required in run options")
	}

	// Build the container run command
	args := []string{"run"}

	// Always run detached
	args = append(args, "-d")

	// Name the container with the workspace ID
	args = append(args, "--name", containerID)

	// Set user
	if opts.User != "" {
		args = append(args, "-u", opts.User)
	}

	// Set entrypoint
	if opts.Entrypoint != "" {
		args = append(args, "--entrypoint", opts.Entrypoint)
	}

	// Set environment variables
	for key, value := range opts.Env {
		args = append(args, "-e", key+"="+value)
	}

	// Set labels
	for _, label := range opts.Labels {
		args = append(args, "-l", label)
	}

	// Set workspace mount
	if opts.WorkspaceMount != nil {
		mountArg := buildMountArg(opts.WorkspaceMount)
		if mountArg != "" {
			args = append(args, "--mount", mountArg)
		}
	}

	// Set additional mounts
	for _, mount := range opts.Mounts {
		mountArg := buildMountArg(mount)
		if mountArg != "" {
			args = append(args, "--mount", mountArg)
		}
	}

	// Log ignored options
	if len(opts.CapAdd) > 0 {
		Logf("Warning: capAdd is not supported by Apple Container, ignoring: %v", opts.CapAdd)
	}
	if len(opts.SecurityOpt) > 0 {
		Logf("Warning: securityOpt is not supported by Apple Container, ignoring: %v", opts.SecurityOpt)
	}
	if opts.Privileged != nil && *opts.Privileged {
		Logf("Warning: privileged mode is not supported by Apple Container, ignoring")
	}

	// Image (must come after all flags)
	args = append(args, opts.Image)

	// Command arguments (after the image)
	if len(opts.Cmd) > 0 {
		args = append(args, opts.Cmd...)
	}

	Logf("Running: %s %s", ContainerBinaryPath(), strings.Join(args, " "))

	// Execute and pipe stdout/stderr to our stdout/stderr (DevPod logs these)
	return RunContainerCmdWithStdoutStderr(args...)
}

// buildMountArg converts a DevPod Mount struct to an Apple Container --mount flag value.
// Format: type=<type>,source=<source>,target=<target>[,readonly]
func buildMountArg(m *Mount) string {
	if m == nil || m.Target == "" {
		return ""
	}

	parts := []string{}

	mountType := m.Type
	if mountType == "" {
		mountType = "bind"
	}
	parts = append(parts, "type="+mountType)

	if m.Source != "" {
		parts = append(parts, "source="+m.Source)
	}

	parts = append(parts, "target="+m.Target)

	// Append any extra options from Other
	for _, o := range m.Other {
		parts = append(parts, o)
	}

	return strings.Join(parts, ",")
}
