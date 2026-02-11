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
