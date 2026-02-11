package cmd

import (
	"encoding/json"
	"os"
	"strings"
	"time"
)

// CFAbsoluteTimeEpoch is the offset between Unix epoch (1970-01-01) and
// Core Foundation absolute time epoch (2001-01-01) in seconds.
const CFAbsoluteTimeEpoch = 978307200

// Find implements the findDevContainer custom driver command.
// It lists all containers, finds one matching by configuration.id or
// the dev.containers.id label, and outputs DevPod-compatible ContainerDetails JSON.
// If no container is found, it outputs nothing (empty stdout) which DevPod interprets as "not found".
func Find(containerID string) error {
	if containerID == "" {
		return nil
	}

	// List all containers in JSON format
	out, err := RunContainerCmd("list", "--all", "--format", "json")
	if err != nil {
		Logf("Warning: failed to list containers: %v", err)
		return nil
	}

	// Parse the JSON output — Apple Container returns an array
	var containers []AppleContainerEntry
	if err := json.Unmarshal(out, &containers); err != nil {
		Logf("Warning: failed to parse container list: %v (raw: %s)", err, string(out[:min(len(out), 200)]))
		return nil
	}

	// Find container matching by configuration.id or by dev.containers.id label
	var matched *AppleContainerEntry
	for i := range containers {
		c := &containers[i]
		// Match by configuration.id (exact match or prefix)
		if c.Configuration.ID == containerID || strings.HasPrefix(c.Configuration.ID, containerID) {
			matched = c
			break
		}
		// Match by dev.containers.id label
		if c.Configuration.Labels != nil {
			if id, ok := c.Configuration.Labels["dev.containers.id"]; ok && id == containerID {
				matched = c
				break
			}
		}
	}

	if matched == nil {
		// No matching container — output nothing
		return nil
	}

	// Convert to DevPod ContainerDetails and output
	details := convertAppleToDevPod(matched)
	encoder := json.NewEncoder(os.Stdout)
	return encoder.Encode(details)
}

// convertAppleToDevPod converts an Apple Container entry to the DevPod ContainerDetails format.
func convertAppleToDevPod(ac *AppleContainerEntry) *ContainerDetails {
	// Ensure labels is never nil — DevPod dereferences it
	labels := ac.Configuration.Labels
	if labels == nil {
		labels = make(map[string]string)
	}

	// Convert CFAbsoluteTime to RFC3339
	startedAt := ""
	if ac.StartedDate > 0 {
		unixTime := ac.StartedDate + CFAbsoluteTimeEpoch
		sec := int64(unixTime)
		nsec := int64((unixTime - float64(sec)) * 1e9)
		startedAt = time.Unix(sec, nsec).UTC().Format(time.RFC3339)
	}

	// Extract user string
	user := ""
	if ac.Configuration.InitProcess.User != nil &&
		ac.Configuration.InitProcess.User.Raw != nil {
		user = ac.Configuration.InitProcess.User.Raw.UserString
	}

	return &ContainerDetails{
		ID:      ac.Configuration.ID,
		Created: startedAt, // Apple Container doesn't have a separate created time
		State: ContainerDetailsState{
			Status:    normalizeState(ac.Status),
			StartedAt: startedAt,
		},
		Config: ContainerDetailsConfig{
			Labels:      labels,
			WorkingDir:  ac.Configuration.InitProcess.WorkingDirectory,
			LegacyUser:  user,
			LegacyImage: ac.Configuration.Image.Reference,
		},
	}
}

// normalizeState converts Apple Container state strings to what DevPod expects.
// DevPod checks for "running" specifically (case-insensitive).
func normalizeState(state string) string {
	s := strings.ToLower(strings.TrimSpace(state))
	switch s {
	case "running":
		return "running"
	case "stopped", "exited", "created":
		return "exited"
	default:
		if s != "" {
			return s
		}
		return ""
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
