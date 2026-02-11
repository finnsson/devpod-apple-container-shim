.PHONY: build clean install

BINARY_NAME := devpod-apple-container-shim
BUILD_DIR := ./build

build:
	@echo "Building $(BINARY_NAME)..."
	GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME) .

clean:
	rm -rf $(BUILD_DIR)

install: build
	@echo "Installing $(BINARY_NAME) to /usr/local/bin/..."
	cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)

# Build for local testing without installing
test-build: build
	@echo "Binary available at $(BUILD_DIR)/$(BINARY_NAME)"
	@echo "Test with: $(BUILD_DIR)/$(BINARY_NAME) arch"

lint:
	go vet -lostcancel=false ./...
	go tool staticcheck ./...
	deadcode
	govulncheck
	nilaway

deadcode: # identify dead code
	go tool deadcode ./cmd/api/...

govulncheck: # identify vulnerabilities in golang dependencies
	go tool govulncheck ./...

nilaway: # identify nil pointer dereferences or unnecessary nil checks
	go tool nilaway -include-pkgs="github.com/finnsson/devpod-apple-container-shim" -test=false  ./...

vibe-context/container.markdown:
	curl https://raw.githubusercontent.com/apple/container/refs/heads/main/docs/command-reference.md \
		-o vibe-context/container.markdown

vibe-context/container-list-all-format.json:
	container list --all --format json > vibe-context/container-list-all-format.json

vibe-context/container-inspect.json: vibe-context/container-list-all-format.json
	container inspect $$(jq -r '.[].configuration.id' vibe-context/container-list-all-format.json) > vibe-context/container-inspect.json