APP_NAME=goxterm
DIST_DIR=build
MAIN_FILE=main.go
BUILD_FLAGS=

BIN_INSTALL_DIR = /usr/local/bin
WEBAPP_INSTALL_DIR = /usr/local/share

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
	@echo "🔧 Building webapp with Docker..."
	@docker build -t $(APP_NAME)-webapp-build ./webapp
	@docker create --name $(APP_NAME)-webapp-temp-container $(APP_NAME)-webapp-build
	@docker cp $(APP_NAME)-webapp-temp-container:/app/dist ./build/dist
	@docker rm $(APP_NAME)-webapp-temp-container
else
	@echo "🔧 Building $(APP_NAME)..."
	GOOS=$(shell go env GOOS) GOARCH=$(shell go env GOARCH) go build $(BUILD_FLAGS) -o $(DIST_DIR)/$(APP_NAME) $(MAIN_FILE)
	@echo "🔧 Building webapp..."
	@cd webapp && npm install && npx vite build && cp -R dist ../$(DIST_DIR)
endif
	@echo "✅ Build complete: ./$(DIST_DIR)"

install: build
	@echo "📦 Installing $(APP_NAME) in $(BIN_INSTALL_DIR)..."
	@sudo install -Dm 0755 $(DIST_DIR)/$(APP_NAME) $(BIN_INSTALL_DIR)/$(APP_NAME)
	@echo "📦 Installing $(APP_NAME) in $(WEBAPP_INSTALL_DIR)..."
	@sudo mkdir -p $(WEBAPP_INSTALL_DIR)/$(APP_NAME)
	@sudo cp -r $(DIST_DIR)/dist/* $(WEBAPP_INSTALL_DIR)/$(APP_NAME)/
	@echo "✅ Installation complete. Now you can use the command '$(APP_NAME)' directly."

uninstall:
	@echo "🗑️ Uninstalling $(APP_NAME) from $(BIN_INSTALL_DIR)..."
	@sudo rm -f $(BIN_INSTALL_DIR)/$(APP_NAME)
	@echo "🗑️ Uninstalling $(APP_NAME) from $(WEBAPP_INSTALL_DIR)..."
	@sudo rm -rf $(WEBAPP_INSTALL_DIR)/$(APP_NAME)
	@echo "✅ Uninstallation complete."

clean:
	@echo "🧹 Cleaning up build files..."
	@rm -rf $(DIST_DIR)
	@echo "✅ Complete clean"
