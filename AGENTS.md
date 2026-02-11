# AGENTS.md

## Project Overview

**devpod-apple-container-shim** is a Go CLI that implements the [DevPod custom driver](https://devpod.sh/docs/developing-providers/driver) interface, translating DevPod's container lifecycle commands into [Apple Container](https://github.com/apple/container) CLI calls.

- **Language**: Go 1.25+ (no external dependencies)
- **Target**: macOS Apple Silicon (darwin/arm64) only
- **Module**: `github.com/finnsson/devpod-apple-container-shim`

## Repository Structure

```
├── main.go                  # Entry point — dispatches subcommands to cmd/
├── cmd/                     # All command implementations
│   ├── types.go             # DevPod + Apple Container JSON type definitions
│   ├── find.go              # findDevContainer — list & match containers
│   ├── run.go               # runDevContainer — create & start containers
│   ├── command.go           # commandDevContainer — exec into containers
│   ├── start.go             # startDevContainer
│   ├── stop.go              # stopDevContainer
│   ├── delete.go            # deleteDevContainer
│   ├── arch.go              # targetArchitecture — outputs "arm64"
│   ├── logs.go              # getDevContainerLogs
│   └── container.go         # Helper functions for invoking Apple Container CLI
├── provider.yaml            # Release provider config (downloads binary from GitHub)
├── provider-local.yaml      # Local dev provider config (uses local build)
├── Makefile                 # Build, clean, install targets
├── vibe-context/            # Reference data for development
│   ├── container.markdown   # Apple Container CLI command reference
│   ├── container-list-all-format.json
│   └── container-inspect.json
└── .devcontainer.json       # DevContainer config for developing this project
```

## Build & Test

```bash
# Build the binary
make build

# Clean build artifacts
make clean

# Install to /usr/local/bin
make install

# Refresh vibe-context reference data
make vibe-context/container.markdown
make vibe-context/container-list-all-format.json
make vibe-context/container-inspect.json
```

## Local Development Workflow

```bash
# Build
make build

# Add the local provider to DevPod (first time only)
devpod provider add ./provider-local.yaml
devpod provider set-options apple-container -o SHIM_PATH=$(pwd)/build/devpod-apple-container-shim

# Test
devpod up --provider apple-container --debug ubuntu:latest

# After code changes, rebuild and retry
make build
devpod delete ubuntu-latest
devpod up --provider apple-container --debug ubuntu:latest
```

## Key Technical Details

### DevPod Custom Driver Protocol

The shim implements 8 commands that DevPod invokes via `agent.custom` in provider YAML:

| Command | Env Vars | Stdin/Stdout |
|---------|----------|--------------|
| `find <id>` | `DEVCONTAINER_ID` | stdout: ContainerDetails JSON (empty = not found) |
| `run <id>` | `DEVCONTAINER_ID`, `DEVCONTAINER_RUN_OPTIONS` (JSON) | stderr: logs |
| `command <id>` | `DEVCONTAINER_ID`, `DEVCONTAINER_USER`, `DEVCONTAINER_COMMAND` | stdin/stdout/stderr passthrough (binary-safe) |
| `start <id>` | `DEVCONTAINER_ID` | — |
| `stop <id>` | `DEVCONTAINER_ID` | — |
| `delete <id>` | `DEVCONTAINER_ID` | — |
| `arch` | — | stdout: architecture string (e.g. `arm64`) |
| `logs <id>` | `DEVCONTAINER_ID` | stdout: container logs |

### Apple Container JSON Schema

Apple Container's `list --format json` and `inspect` output uses a nested structure:

```json
{
  "configuration": {
    "id": "...",
    "labels": {},
    "image": { "reference": "..." },
    "initProcess": {
      "executable": "...",
      "environment": ["KEY=VALUE", ...],
      "workingDirectory": "...",
      "user": { "raw": { "userString": "..." } }
    }
  },
  "status": "running",
  "startedDate": 792529857.216959
}
```

- `startedDate` is **CFAbsoluteTime** (seconds since 2001-01-01). Convert to Unix by adding `978307200`.
- `environment` is `[]string` of `KEY=VALUE` pairs, not a map.

### Provider YAML Notes

- `agent.custom` command arrays use `exec.CommandContext` directly — **no shell expansion**. That's why all commands are wrapped in `["sh", "-c", "..."]`.
- `agent.binaries` injects a binary path as an env var (e.g. `$APPLE_CONTAINER_SHIM` or `$SHIM_PATH`).
- Docker-specific mount options (`consistency`, `bind-propagation`, etc.) are filtered out in `run.go` since Apple Container doesn't support them.

### Known Limitations

- Only macOS Apple Silicon is supported.
- `capAdd`, `securityOpt`, and `privileged` from DevPod's RunOptions are silently ignored.
- Minimal images (e.g. `ubuntu:latest`) lack git, causing DevPod agent to log git-related errors. Use `mcr.microsoft.com/devcontainers/base:ubuntu` or disable git injection with `devpod context set-options default -o SSH_INJECT_GIT_CREDENTIALS=false`.

## Relevant Links

- [DevPod Custom Driver Docs](https://devpod.sh/docs/developing-providers/driver)
- [DevPod Custom Driver Source (custom.go)](https://github.com/loft-sh/devpod/blob/main/pkg/driver/custom/custom.go)
- [Apple Container GitHub](https://github.com/apple/container)
- [Apple Container Command Reference](https://github.com/apple/container/blob/main/docs/command-reference.md)
