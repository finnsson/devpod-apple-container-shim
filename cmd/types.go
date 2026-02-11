// Package cmd implements the DevPod custom driver commands for Apple Container.
package cmd

// ContainerDetails is the JSON structure DevPod expects from findDevContainer.
type ContainerDetails struct {
	ID      string                 `json:"ID,omitempty"`
	Created string                 `json:"Created,omitempty"`
	State   ContainerDetailsState  `json:"State,omitempty"`
	Config  ContainerDetailsConfig `json:"Config,omitempty"`
}

type ContainerDetailsState struct {
	Status    string `json:"Status,omitempty"`
	StartedAt string `json:"StartedAt,omitempty"`
}

type ContainerDetailsConfig struct {
	Labels      map[string]string `json:"Labels,omitempty"`
	WorkingDir  string            `json:"WorkingDir,omitempty"`
	LegacyUser  string            `json:"User,omitempty"`
	LegacyImage string            `json:"Image,omitempty"`
}

// RunOptions is the JSON received via the DEVCONTAINER_RUN_OPTIONS env var.
type RunOptions struct {
	UID            string            `json:"uid,omitempty"`
	Image          string            `json:"image,omitempty"`
	User           string            `json:"user,omitempty"`
	Entrypoint     string            `json:"entrypoint,omitempty"`
	Cmd            []string          `json:"cmd,omitempty"`
	Env            map[string]string `json:"env,omitempty"`
	CapAdd         []string          `json:"capAdd,omitempty"`
	SecurityOpt    []string          `json:"securityOpt,omitempty"`
	Labels         []string          `json:"labels,omitempty"`
	Privileged     *bool             `json:"privileged,omitempty"`
	WorkspaceMount *Mount            `json:"workspaceMount,omitempty"`
	Mounts         []*Mount          `json:"mounts,omitempty"`
}

type Mount struct {
	Type     string   `json:"type,omitempty"`
	Source   string   `json:"source,omitempty"`
	Target   string   `json:"target,omitempty"`
	External bool     `json:"external,omitempty"`
	Other    []string `json:"other,omitempty"`
}

// AppleContainerEntry represents a single entry from Apple Container's
// `container list --format json` or `container inspect` output.
type AppleContainerEntry struct {
	Configuration AppleContainerConfig `json:"configuration"`
	Status        string               `json:"status"`
	StartedDate   float64              `json:"startedDate"` // CFAbsoluteTime: seconds since 2001-01-01
	Networks      []interface{}        `json:"networks,omitempty"`
}

type AppleContainerConfig struct {
	ID          string                    `json:"id"`
	Labels      map[string]string         `json:"labels,omitempty"`
	Image       AppleContainerImage       `json:"image"`
	InitProcess AppleContainerInitProcess `json:"initProcess"`
	Platform    *AppleContainerPlatform   `json:"platform,omitempty"`
}

type AppleContainerImage struct {
	Reference  string                    `json:"reference"`
	Descriptor *AppleContainerDescriptor `json:"descriptor,omitempty"`
}

type AppleContainerDescriptor struct {
	MediaType string `json:"mediaType,omitempty"`
	Digest    string `json:"digest,omitempty"`
	Size      int64  `json:"size,omitempty"`
}

type AppleContainerInitProcess struct {
	Executable       string              `json:"executable,omitempty"`
	Arguments        []string            `json:"arguments,omitempty"`
	Environment      []string            `json:"environment,omitempty"`
	WorkingDirectory string              `json:"workingDirectory,omitempty"`
	User             *AppleContainerUser `json:"user,omitempty"`
}

type AppleContainerUser struct {
	Raw *AppleContainerRawUser `json:"raw,omitempty"`
}

type AppleContainerRawUser struct {
	UserString string `json:"userString,omitempty"`
}

type AppleContainerPlatform struct {
	Architecture string `json:"architecture,omitempty"`
	OS           string `json:"os,omitempty"`
}
