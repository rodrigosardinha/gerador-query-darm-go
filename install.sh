#!/bin/bash

# Script de instalação para Processador de DARMs em Go
# Compatível com Linux e macOS

set -e

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Variáveis
PROJECT_NAME="gerador-query-darm-go"
BINARY_NAME="darm-processor"
VERSION="1.0.0"

# Função para imprimir mensagens coloridas
print_message() {
    echo -e "${2}${1}${NC}"
}

# Função para verificar se um comando existe
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Função para detectar o sistema operacional
detect_os() {
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        echo "linux"
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        echo "macos"
    elif [[ "$OSTYPE" == "cygwin" ]] || [[ "$OSTYPE" == "msys" ]] || [[ "$OSTYPE" == "win32" ]]; then
        echo "windows"
    else
        echo "unknown"
    fi
}

# Função para verificar requisitos
check_requirements() {
    print_message "🔍 Verificando requisitos do sistema..." "$BLUE"
    
    # Verificar Go
    if ! command_exists go; then
        print_message "❌ Go não está instalado!" "$RED"
        print_message "💡 Instale o Go em: https://golang.org/dl/" "$YELLOW"
        exit 1
    fi
    
    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    print_message "✅ Go encontrado: $GO_VERSION" "$GREEN"
    
    # Verificar Git
    if ! command_exists git; then
        print_message "❌ Git não está instalado!" "$RED"
        print_message "💡 Instale o Git para clonar o repositório" "$YELLOW"
        exit 1
    fi
    
    print_message "✅ Git encontrado: $(git --version)" "$GREEN"
    
    # Verificar Make (opcional)
    if command_exists make; then
        print_message "✅ Make encontrado" "$GREEN"
    else
        print_message "⚠️ Make não encontrado (opcional)" "$YELLOW"
    fi
    
    # Verificar Docker (opcional)
    if command_exists docker; then
        print_message "✅ Docker encontrado" "$GREEN"
    else
        print_message "⚠️ Docker não encontrado (opcional)" "$YELLOW"
    fi
}

# Função para instalar dependências
install_dependencies() {
    print_message "📦 Instalando dependências..." "$BLUE"
    
    # Baixar dependências
    go mod download
    
    # Verificar dependências
    go mod tidy
    
    print_message "✅ Dependências instaladas!" "$GREEN"
}

# Função para criar diretórios
create_directories() {
    print_message "📁 Criando diretórios necessários..." "$BLUE"
    
    mkdir -p darms inserts build
    
    print_message "✅ Diretórios criados!" "$GREEN"
}

# Função para criar arquivo de configuração
create_config() {
    print_message "⚙️ Criando arquivo de configuração padrão..." "$BLUE"
    
    cat > config.json << EOF
{
  "database": {
    "host": "localhost",
    "port": 3306,
    "database": "silfae",
    "username": "root",
    "password": "",
    "charset": "latin1"
  },
  "paths": {
    "base_dir": ".",
    "darms_dir": "darms",
    "output_dir": "inserts",
    "temp_dir": "temp"
  },
  "sql": {
    "encoding": "latin1",
    "batch_size": 100,
    "use_transaction": true,
    "use_ignore": true
  },
  "logging": {
    "level": "info",
    "format": "text",
    "output_file": ""
  }
}
EOF
    
    print_message "✅ Arquivo config.json criado!" "$GREEN"
}

# Função para compilar o projeto
build_project() {
    print_message "🔨 Compilando o projeto..." "$BLUE"
    
    # Compilar
    go build -ldflags="-X main.version=$VERSION" -o build/$BINARY_NAME .
    
    print_message "✅ Executável criado: build/$BINARY_NAME" "$GREEN"
}

# Função para executar testes
run_tests() {
    print_message "🧪 Executando testes..." "$BLUE"
    
    # Executar testes
    go test -v ./...
    
    print_message "✅ Testes concluídos!" "$GREEN"
}

# Função para executar benchmarks
run_benchmarks() {
    print_message "⚡ Executando benchmarks..." "$BLUE"
    
    # Executar benchmarks
    go test -bench=. ./...
    
    print_message "✅ Benchmarks concluídos!" "$GREEN"
}

# Função para criar arquivo .gitkeep
create_gitkeep() {
    print_message "📝 Criando arquivos .gitkeep..." "$BLUE"
    
    touch darms/.gitkeep
    touch inserts/.gitkeep
    
    print_message "✅ Arquivos .gitkeep criados!" "$GREEN"
}

# Função para mostrar informações do sistema
show_system_info() {
    print_message "💻 Informações do sistema:" "$CYAN"
    echo "  Sistema Operacional: $(detect_os)"
    echo "  Go Version: $(go version)"
    echo "  Git Version: $(git --version)"
    echo "  Arquitetura: $(uname -m)"
    echo "  Kernel: $(uname -r)"
}

# Função para mostrar próximos passos
show_next_steps() {
    print_message "🎉 Instalação concluída com sucesso!" "$GREEN"
    print_message "" "$NC"
    print_message "📋 Próximos passos:" "$YELLOW"
    print_message "  1. Coloque PDFs dos DARMs na pasta darms/" "$NC"
    print_message "  2. Execute: ./build/$BINARY_NAME" "$NC"
    print_message "  3. Verifique os arquivos gerados na pasta inserts/" "$NC"
    print_message "" "$NC"
    print_message "🛠️ Comandos úteis:" "$YELLOW"
    print_message "  make help        - Mostra ajuda do Makefile" "$NC"
    print_message "  make test        - Executa testes" "$NC"
    print_message "  make run         - Executa o programa" "$NC"
    print_message "  make build       - Recompila o executável" "$NC"
    print_message "  make clean       - Remove arquivos de build" "$NC"
    print_message "" "$NC"
    print_message "📚 Documentação:" "$YELLOW"
    print_message "  README_Go.md     - Documentação completa" "$NC"
    print_message "  exemplo_uso.go   - Exemplos de uso" "$NC"
}

# Função principal
main() {
    print_message "🏛️ Processador de DARMs - Versão Go $VERSION" "$CYAN"
    print_message "🚀 Iniciando instalação..." "$BLUE"
    
    # Verificar requisitos
    check_requirements
    
    # Mostrar informações do sistema
    show_system_info
    
    # Instalar dependências
    install_dependencies
    
    # Criar diretórios
    create_directories
    
    # Criar arquivo de configuração
    create_config
    
    # Criar arquivos .gitkeep
    create_gitkeep
    
    # Compilar projeto
    build_project
    
    # Executar testes
    run_tests
    
    # Executar benchmarks
    run_benchmarks
    
    # Mostrar próximos passos
    show_next_steps
}

# Verificar se o script está sendo executado no diretório correto
if [[ ! -f "go.mod" ]]; then
    print_message "❌ Erro: Execute este script no diretório raiz do projeto!" "$RED"
    print_message "💡 Certifique-se de que o arquivo go.mod existe no diretório atual." "$YELLOW"
    exit 1
fi

# Executar função principal
main "$@" 