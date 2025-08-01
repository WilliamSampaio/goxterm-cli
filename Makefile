APP_NAME=goxterm
DIST_DIR=build
MAIN_FILE=main.go
BUILD_FLAGS=

OS=$(shell uname -s)

.PHONY: all build clean install

all: build

build:
	@echo "🔧 Building $(APP_NAME)..."
	@mkdir -p $(DIST_DIR)
	GOOS=$(shell go env GOOS) GOARCH=$(shell go env GOARCH) go build $(BUILD_FLAGS) -o $(DIST_DIR)/$(APP_NAME) $(MAIN_FILE)
	@echo "✅ Build complete: $(DIST_DIR)/$(APP_NAME)"

install: build
	@echo "📦 Installing $(APP_NAME) in /usr/local/bin..."
	@sudo cp $(DIST_DIR)/$(APP_NAME) /usr/local/bin/$(APP_NAME)
	@echo "✅ Installation complete. Now you can use the command '$(APP_NAME)' directly."

clean:
	@echo "🧹 Cleaning up build files..."
	@rm -rf $(DIST_DIR)
	@echo "✅ Complete clean"
