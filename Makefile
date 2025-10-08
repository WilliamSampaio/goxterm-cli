APP_NAME=goxterm
DIST_DIR=build
MAIN_FILE=main.go
BUILD_FLAGS=

BIN_INSTALL_DIR = /usr/local/bin
FILES_INSTALL_DIR = /usr/local/share

ZIP ?= zip

OS=$(shell uname -s)

.PHONY: all build install uninstall clean

all: build

build: clean
	@mkdir -p $(DIST_DIR)
ifeq ($(DOCKER),1)
	@echo "🔧 Building $(APP_NAME) with Docker..."
	@docker build --build-arg APP_NAME="$(APP_NAME)" --build-arg MAIN_FILE="$(MAIN_FILE)" -t $(APP_NAME)-bin-build .
	@docker create --name $(APP_NAME)-bin-temp-container $(APP_NAME)-bin-build
	@docker cp $(APP_NAME)-bin-temp-container:/app/bin/$(APP_NAME) ./build
	@docker rm $(APP_NAME)-bin-temp-container
else
	@echo "🔧 Building $(APP_NAME)..."
	GOOS=$(shell go env GOOS) GOARCH=$(shell go env GOARCH) go build $(BUILD_FLAGS) -o $(DIST_DIR)/$(APP_NAME) $(MAIN_FILE)
endif
	@command -v $(ZIP) > /dev/null 2>&1 || { \
		echo "❌ O comando '${ZIP}' não está instalado. Por favor, instale-o e tente novamente."; \
		exit 1; \
	}
	@echo "📦 Packaging Chrome extension..."
	@cd extension && zip -r ../$(DIST_DIR)/$(APP_NAME)-extension-chrome.zip .
# 	@echo "📦 Packaging Firefox extension..."
# 	@cd extension && zip -r ../$(DIST_DIR)/$(APP_NAME)-extension-firefox.xpi .
	@echo "✅ Build complete: ./$(DIST_DIR)"

install: build
	@echo "📦 Installing $(APP_NAME) in $(BIN_INSTALL_DIR)..."
	@sudo install -Dm 0755 $(DIST_DIR)/$(APP_NAME) $(BIN_INSTALL_DIR)/$(APP_NAME)
	@echo "📦 Installing $(APP_NAME) in $(FILES_INSTALL_DIR)..."
	@sudo mkdir -p $(FILES_INSTALL_DIR)/$(APP_NAME)
	@sudo cp -r assets $(FILES_INSTALL_DIR)/$(APP_NAME)/
	@sudo cp -r pages $(FILES_INSTALL_DIR)/$(APP_NAME)/
	@echo "✅ Installation complete. Now you can use the command '$(APP_NAME)' directly."

uninstall:
	@echo "🗑️ Uninstalling $(APP_NAME) from $(BIN_INSTALL_DIR)..."
	@sudo rm -f $(BIN_INSTALL_DIR)/$(APP_NAME)
	@echo "✅ Uninstallation complete."

clean:
	@echo "🧹 Cleaning up build files..."
	@rm -rf $(DIST_DIR)
	@echo "✅ Complete clean"
