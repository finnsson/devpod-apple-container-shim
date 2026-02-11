package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// Find implements the findDevContainer custom driver command.
// It lists all containers, filters by the dev.containers.id label matching the workspaceID,
// then inspects the matching container and outputs DevPod-compatible ContainerDetails JSON.
// If no container is found, it outputs nothing (empty stdout) which DevPod interprets as "not found".
func Find(containerID string) error {
	if containerID == "" {
		// No container to find — output nothing
		return nil
	}

	// List all containers in JSON format
	out, err := RunContainerCmd("list", "--all", "--format", "json")
	if err != nil {
		// If listing fails, it may mean no containers exist. Output nothing.
		Logf("Warning: failed to list containers: %v", err)
		return nil
	}

	// Parse the JSON output — Apple Container returns an array of container entries
	var containers []AppleContainerListEntry
	if err := json.Unmarshal(out, &containers); err != nil {
		// Try parsing as newline-delimited JSON
		containers, err = parseContainerListNDJSON(out)
		if err != nil {
			Logf("Warning: failed to parse container list: %v", err)
			return nil
		}
	}

	// Filter by label dev.containers.id=<containerID>
	var matchedID string
	for _, c := range containers {
		if c.Labels != nil {
			if id, ok := c.Labels["dev.containers.id"]; ok && id == containerID {
				matchedID = c.ID
				if matchedID == "" {
					matchedID = c.Name
				}
				break
			}
		}
	}

	if matchedID == "" {
		// No matching container found — output nothing
		return nil
	}

	// Inspect the matched container
	inspectOut, err := RunContainerCmd("inspect", matchedID)
	if err != nil {
		return fmt.Errorf("failed to inspect container %s: %w", matchedID, err)
	}

	// Parse the inspect output
	details, err := parseInspectOutput(inspectOut, matchedID)
	if err != nil {
		return fmt.Errorf("failed to parse inspect output: %w", err)
	}

	// Output the DevPod ContainerDetails JSON
	encoder := json.NewEncoder(os.Stdout)
	return encoder.Encode(details)
}

// parseContainerListNDJSON tries to parse newline-delimited JSON entries
func parseContainerListNDJSON(data []byte) ([]AppleContainerListEntry, error) {
	var containers []AppleContainerListEntry
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || line == "[" || line == "]" {
			continue
		}
		// Strip trailing comma
		line = strings.TrimRight(line, ",")
		var c AppleContainerListEntry
		if err := json.Unmarshal([]byte(line), &c); err != nil {
			continue
		}
		containers = append(containers, c)
	}
	if len(containers) == 0 {
		return nil, fmt.Errorf("no containers parsed from output")
	}
	return containers, nil
}

// parseInspectOutput parses Apple Container inspect JSON and converts to DevPod ContainerDetails.
// Apple Container inspect may return a single object or an array.
func parseInspectOutput(data []byte, containerID string) (*ContainerDetails, error) {
	trimmed := strings.TrimSpace(string(data))

	// Try as a single object first
	var single AppleContainerInspect
	if err := json.Unmarshal([]byte(trimmed), &single); err == nil && single.ID != "" {
		return convertToContainerDetails(&single), nil
	}

	// Try as an array
	var arr []AppleContainerInspect
	if err := json.Unmarshal([]byte(trimmed), &arr); err == nil && len(arr) > 0 {
		return convertToContainerDetails(&arr[0]), nil
	}

	// Try as a generic map and extract what we can
	var generic map[string]interface{}
	if err := json.Unmarshal([]byte(trimmed), &generic); err == nil {
		return convertGenericToContainerDetails(generic, containerID), nil
	}

	return nil, fmt.Errorf("could not parse inspect output: %s", trimmed[:min(len(trimmed), 200)])
}

func convertToContainerDetails(ac *AppleContainerInspect) *ContainerDetails {
	state := normalizeState(ac.State)
	if state == "" {
		state = normalizeState(ac.Status)
	}

	labels := ac.Labels
	if labels == nil && ac.Config != nil {
		labels = ac.Config.Labels
	}

	user := ""
	if ac.Config != nil {
		user = ac.Config.User
	}
	if user == "" && ac.ProcessConfig != nil {
		user = ac.ProcessConfig.User
	}

	image := ac.Image
	if image == "" && ac.Config != nil {
		image = ac.Config.Image
	}

	workingDir := ""
	if ac.Config != nil {
		workingDir = ac.Config.WorkingDir
	}

	id := ac.ID
	if id == "" {
		id = ac.Name
	}

	return &ContainerDetails{
		ID:      id,
		Created: ac.CreatedAt,
		State: ContainerDetailsState{
			Status:    state,
			StartedAt: ac.StartedAt,
		},
		Config: ContainerDetailsConfig{
			Labels:      labels,
			WorkingDir:  workingDir,
			LegacyUser:  user,
			LegacyImage: image,
		},
	}
}

func convertGenericToContainerDetails(data map[string]interface{}, containerID string) *ContainerDetails {
	details := &ContainerDetails{
		ID: containerID,
	}

	if id, ok := data["id"].(string); ok && id != "" {
		details.ID = id
	}
	if name, ok := data["name"].(string); ok && name != "" && details.ID == "" {
		details.ID = name
	}
	if created, ok := data["createdAt"].(string); ok {
		details.Created = created
	}
	if state, ok := data["state"].(string); ok {
		details.State.Status = normalizeState(state)
	}
	if status, ok := data["status"].(string); ok && details.State.Status == "" {
		details.State.Status = normalizeState(status)
	}
	if startedAt, ok := data["startedAt"].(string); ok {
		details.State.StartedAt = startedAt
	}

	// Extract labels
	if labels, ok := data["labels"].(map[string]interface{}); ok {
		details.Config.Labels = make(map[string]string)
		for k, v := range labels {
			if s, ok := v.(string); ok {
				details.Config.Labels[k] = s
			}
		}
	}

	return details
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
