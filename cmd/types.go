package cmd

// DevPod ContainerDetails - the JSON structure DevPod expects from findDevContainer
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

// DevPod RunOptions - JSON received via DEVCONTAINER_RUN_OPTIONS env var
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

// Apple Container list --format json output structures
type AppleContainerListEntry struct {
	ID     string            `json:"id"`
	Name   string            `json:"name"`
	Image  string            `json:"image"`
	State  string            `json:"state"`
	Status string            `json:"status"`
	Labels map[string]string `json:"labels"`
}

// Apple Container inspect output structure
type AppleContainerInspect struct {
	ID            string                       `json:"id"`
	Name          string                       `json:"name"`
	Image         string                       `json:"image"`
	State         string                       `json:"state"`
	Status        string                       `json:"status"`
	CreatedAt     string                       `json:"createdAt"`
	StartedAt     string                       `json:"startedAt"`
	Labels        map[string]string            `json:"labels"`
	Config        *AppleContainerInspectConfig `json:"config"`
	ProcessConfig *AppleContainerProcessConfig `json:"processConfig"`
}

type AppleContainerInspectConfig struct {
	Image      string            `json:"image"`
	Env        map[string]string `json:"env"`
	Labels     map[string]string `json:"labels"`
	WorkingDir string            `json:"workingDir"`
	User       string            `json:"user"`
	Entrypoint []string          `json:"entrypoint"`
	Cmd        []string          `json:"cmd"`
}

type AppleContainerProcessConfig struct {
	User string            `json:"user"`
	Env  map[string]string `json:"env"`
}
