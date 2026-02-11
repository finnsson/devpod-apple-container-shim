# devpod-apple-container-shim

A CLI shim that implements the [DevPod custom driver](https://devpod.sh/docs/developing-providers/driver) interface on top of [Apple Container](https://github.com/apple/container) for macOS Apple Silicon.

## Prerequisites

- macOS on Apple Silicon (arm64)
- [Apple Container CLI](https://github.com/apple/container) installed at `/usr/local/bin/container`
- [DevPod](https://devpod.sh/) v0.6+
- Go 1.25+ (for building from source)

## Building

```bash
make build
```

The binary is output to `build/devpod-apple-container-shim`.

## Provider Configuration

There are two provider YAML files:

### provider-local.yaml — Local Development

Use this when developing and testing the shim locally. It references the binary from the local build directory via the `SHIM_PATH` option.

```bash
# Build the shim
make build

# Add the local provider
devpod provider add ./provider-local.yaml

# Point it to the local build
devpod provider set-options apple-container -o SHIM_PATH=$(pwd)/build/devpod-apple-container-shim

# Launch a workspace
devpod up --provider apple-container ubuntu:latest
```

### provider.yaml — Release / Distribution

Use this for end users. It downloads the shim binary from a GitHub release via `agent.binaries`. Users install with:

```bash
devpod provider add https://raw.githubusercontent.com/finnsson/devpod-apple-container-shim/main/provider.yaml
devpod up --provider apple-container ubuntu:latest
```

Both providers use `sh -c` wrappers for variable expansion (DevPod's `agent.custom` arrays use `exec.CommandContext` directly, so `${VAR}` won't expand without a shell).

## Tips

- **Git inside the container**: The base `ubuntu:latest` image doesn't include git. DevPod's agent will log errors about missing git. Either use an image with git (e.g. `mcr.microsoft.com/devcontainers/base:ubuntu`) or disable git credential injection:
  ```bash
  devpod context set-options default -o SSH_INJECT_GIT_CREDENTIALS=false
  ```

- **Cleaning up orphan containers**: If a `devpod up` is interrupted, orphaned Apple Container instances may remain. List and remove them with:
  ```bash
  container list --all
  container delete <container-id>
  ```