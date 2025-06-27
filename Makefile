# Makefile para Processador de DARMs em Go

# Variáveis
BINARY_NAME=darm-processor
MAIN_FILE=main.go
BUILD_DIR=build
VERSION=1.0.0
GOOS?=$(shell go env GOOS)
GOARCH?=$(shell go env GOARCH)

# Cores para output
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[0;33m
BLUE=\033[0;34m
PURPLE=\033[0;35m
CYAN=\033[0;36m
WHITE=\033[0;37m
NC=\033[0m # No Color

# Comandos principais
.PHONY: help build clean test run install deps lint format docker-build docker-run

# Ajuda
help:
	@echo "$(CYAN)🏛️ Processador de DARMs - Versão Go$(NC)"
	@echo "$(YELLOW)Comandos disponíveis:$(NC)"
	@echo "  $(GREEN)make help$(NC)        - Mostra esta ajuda"
	@echo "  $(GREEN)make build$(NC)       - Compila o executável"
	@echo "  $(GREEN)make clean$(NC)       - Remove arquivos de build"
	@echo "  $(GREEN)make test$(NC)        - Executa testes"
	@echo "  $(GREEN)make run$(NC)         - Executa o programa"
	@echo "  $(GREEN)make install$(NC)     - Instala dependências"
	@echo "  $(GREEN)make lint$(NC)        - Executa linter"
	@echo "  $(GREEN)make format$(NC)      - Formata código"
	@echo "  $(GREEN)make docker-build$(NC) - Build Docker"
	@echo "  $(GREEN)make docker-run$(NC)   - Executa Docker"
	@echo "  $(GREEN)make release$(NC)     - Build para múltiplas plataformas"

# Instalar dependências
install:
	@echo "$(BLUE)📦 Instalando dependências...$(NC)"
	go mod download
	go mod tidy
	@echo "$(GREEN)✅ Dependências instaladas!$(NC)"

# Compilar
build:
	@echo "$(BLUE)🔨 Compilando $(BINARY_NAME)...$(NC)"
	@mkdir -p $(BUILD_DIR)
	go build -ldflags="-X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_FILE)
	@echo "$(GREEN)✅ Executável criado: $(BUILD_DIR)/$(BINARY_NAME)$(NC)"

# Compilar para múltiplas plataformas
release:
	@echo "$(BLUE)🚀 Criando builds para múltiplas plataformas...$(NC)"
	@mkdir -p $(BUILD_DIR)
	
	@echo "$(YELLOW)📱 Windows (amd64)...$(NC)"
	GOOS=windows GOARCH=amd64 go build -ldflags="-X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_FILE)
	
	@echo "$(YELLOW)🐧 Linux (amd64)...$(NC)"
	GOOS=linux GOARCH=amd64 go build -ldflags="-X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_FILE)
	
	@echo "$(YELLOW)🍎 macOS (amd64)...$(NC)"
	GOOS=darwin GOARCH=amd64 go build -ldflags="-X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_FILE)
	
	@echo "$(YELLOW)🍎 macOS (arm64)...$(NC)"
	GOOS=darwin GOARCH=arm64 go build -ldflags="-X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_FILE)
	
	@echo "$(GREEN)✅ Builds criados em $(BUILD_DIR)/$(NC)"

# Limpar
clean:
	@echo "$(YELLOW)🧹 Limpando arquivos de build...$(NC)"
	rm -rf $(BUILD_DIR)
	go clean
	@echo "$(GREEN)✅ Limpeza concluída!$(NC)"

# Executar testes
test:
	@echo "$(BLUE)🧪 Executando testes...$(NC)"
	go test -v ./...
	@echo "$(GREEN)✅ Testes concluídos!$(NC)"

# Executar testes com cobertura
test-coverage:
	@echo "$(BLUE)🧪 Executando testes com cobertura...$(NC)"
	go test -cover ./...
	@echo "$(GREEN)✅ Testes com cobertura concluídos!$(NC)"

# Executar benchmarks
bench:
	@echo "$(BLUE)⚡ Executando benchmarks...$(NC)"
	go test -bench=. ./...
	@echo "$(GREEN)✅ Benchmarks concluídos!$(NC)"

# Executar
run:
	@echo "$(BLUE)🚀 Executando $(BINARY_NAME)...$(NC)"
	go run $(MAIN_FILE)

# Executar com configuração
run-config:
	@echo "$(BLUE)🚀 Executando $(BINARY_NAME) com configuração...$(NC)"
	go run $(MAIN_FILE) -config=config.json

# Linter
lint:
	@echo "$(BLUE)🔍 Executando linter...$(NC)"
	golangci-lint run
	@echo "$(GREEN)✅ Linter concluído!$(NC)"

# Formatar código
format:
	@echo "$(BLUE)🎨 Formatando código...$(NC)"
	go fmt ./...
	go vet ./...
	@echo "$(GREEN)✅ Código formatado!$(NC)"

# Verificar dependências
deps:
	@echo "$(BLUE)📋 Verificando dependências...$(NC)"
	go mod verify
	go list -m all
	@echo "$(GREEN)✅ Dependências verificadas!$(NC)"

# Criar diretórios necessários
setup:
	@echo "$(BLUE)📁 Criando diretórios necessários...$(NC)"
	mkdir -p darms inserts
	@echo "$(GREEN)✅ Diretórios criados!$(NC)"

# Docker build
docker-build:
	@echo "$(BLUE)🐳 Criando imagem Docker...$(NC)"
	docker build -t $(BINARY_NAME):$(VERSION) .
	@echo "$(GREEN)✅ Imagem Docker criada!$(NC)"

# Docker run
docker-run:
	@echo "$(BLUE)🐳 Executando container Docker...$(NC)"
	docker run -v $(PWD)/darms:/app/darms -v $(PWD)/inserts:/app/inserts $(BINARY_NAME):$(VERSION)
	@echo "$(GREEN)✅ Container Docker executado!$(NC)"

# Criar arquivo de configuração padrão
config:
	@echo "$(BLUE)⚙️ Criando arquivo de configuração padrão...$(NC)"
	@echo '{"database":{"host":"localhost","port":3306,"database":"silfae","username":"root","password":"","charset":"latin1"},"paths":{"base_dir":".","darms_dir":"darms","output_dir":"inserts","temp_dir":"temp"},"sql":{"encoding":"latin1","batch_size":100,"use_transaction":true,"use_ignore":true},"logging":{"level":"info","format":"text","output_file":""}}' > config.json
	@echo "$(GREEN)✅ Arquivo config.json criado!$(NC)"

# Verificar versão
version:
	@echo "$(CYAN)Versão: $(VERSION)$(NC)"
	@echo "$(CYAN)Go: $(shell go version)$(NC)"
	@echo "$(CYAN)OS: $(GOOS)/$(GOARCH)$(NC)"

# Desenvolvimento
dev: install setup config
	@echo "$(GREEN)🎉 Ambiente de desenvolvimento configurado!$(NC)"
	@echo "$(YELLOW)Próximos passos:$(NC)"
	@echo "  1. Coloque PDFs dos DARMs na pasta darms/"
	@echo "  2. Execute: make run"
	@echo "  3. Verifique os arquivos gerados na pasta inserts/"

# Build completo
all: clean install format lint test build
	@echo "$(GREEN)🎉 Build completo concluído!$(NC)"

# Instalação completa
full-install: install setup config
	@echo "$(GREEN)🎉 Instalação completa concluída!$(NC)" 