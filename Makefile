# Nome do executável final
APP_NAME=goxterm
# Diretório de saída
DIST_DIR=dist

# Caminho do código-fonte principal
MAIN_FILE=main.go

# Flags de compilação (pode ajustar conforme necessário)
BUILD_FLAGS=

# Detecta sistema operacional para build cross-platform se quiser
OS=$(shell uname -s)

.PHONY: all build run clean

all: build

run:
	go run $(MAIN_FILE)

build:
	@echo "🔧 Buildando $(APP_NAME)..."
	@mkdir -p $(DIST_DIR)
	GOOS=$(shell go env GOOS) GOARCH=$(shell go env GOARCH) go build $(BUILD_FLAGS) -o $(DIST_DIR)/$(APP_NAME) $(MAIN_FILE)
	@echo "✅ Build completo: $(DIST_DIR)/$(APP_NAME)"

clean:
	@echo "🧹 Limpando arquivos de build..."
	rm -rf $(DIST_DIR)
	@echo "✅ Clean completo"
