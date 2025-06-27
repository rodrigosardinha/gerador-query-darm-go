@echo off
setlocal enabledelayedexpansion

REM Script de instalação para Processador de DARMs em Go
REM Compatível com Windows

REM Variáveis
set PROJECT_NAME=gerador-query-darm-go
set BINARY_NAME=darm-processor
set VERSION=1.0.0

REM Cores para output (Windows)
set RED=[91m
set GREEN=[92m
set YELLOW=[93m
set BLUE=[94m
set PURPLE=[95m
set CYAN=[96m
set NC=[0m

REM Função para imprimir mensagens coloridas
:print_message
echo %~2%~1%NC%
goto :eof

REM Função para verificar se um comando existe
:command_exists
where %~1 >nul 2>&1
if %errorlevel% equ 0 (
    exit /b 0
) else (
    exit /b 1
)

REM Função para verificar requisitos
:check_requirements
call :print_message "🔍 Verificando requisitos do sistema..." "%BLUE%"
    
REM Verificar Go
call :command_exists go
if %errorlevel% neq 0 (
    call :print_message "❌ Go não está instalado!" "%RED%"
    call :print_message "💡 Instale o Go em: https://golang.org/dl/" "%YELLOW%"
    exit /b 1
)
    
for /f "tokens=3" %%i in ('go version') do set GO_VERSION=%%i
set GO_VERSION=!GO_VERSION:go=!
call :print_message "✅ Go encontrado: !GO_VERSION!" "%GREEN%"
    
REM Verificar Git
call :command_exists git
if %errorlevel% neq 0 (
    call :print_message "❌ Git não está instalado!" "%RED%"
    call :print_message "💡 Instale o Git para clonar o repositório" "%YELLOW%"
    exit /b 1
)
    
for /f "tokens=*" %%i in ('git --version') do set GIT_VERSION=%%i
call :print_message "✅ Git encontrado: !GIT_VERSION!" "%GREEN%"
    
REM Verificar Make (opcional)
call :command_exists make
if %errorlevel% equ 0 (
    call :print_message "✅ Make encontrado" "%GREEN%"
) else (
    call :print_message "⚠️ Make não encontrado (opcional)" "%YELLOW%"
)
    
REM Verificar Docker (opcional)
call :command_exists docker
if %errorlevel% equ 0 (
    call :print_message "✅ Docker encontrado" "%GREEN%"
) else (
    call :print_message "⚠️ Docker não encontrado (opcional)" "%YELLOW%"
)
goto :eof

REM Função para instalar dependências
:install_dependencies
call :print_message "📦 Instalando dependências..." "%BLUE%"
    
REM Baixar dependências
go mod download
    
REM Verificar dependências
go mod tidy
    
call :print_message "✅ Dependências instaladas!" "%GREEN%"
goto :eof

REM Função para criar diretórios
:create_directories
call :print_message "📁 Criando diretórios necessários..." "%BLUE%"
    
if not exist "darms" mkdir darms
if not exist "inserts" mkdir inserts
if not exist "build" mkdir build
    
call :print_message "✅ Diretórios criados!" "%GREEN%"
goto :eof

REM Função para criar arquivo de configuração
:create_config
call :print_message "⚙️ Criando arquivo de configuração padrão..." "%BLUE%"
    
echo {> config.json
echo   "database": {>> config.json
echo     "host": "localhost",>> config.json
echo     "port": 3306,>> config.json
echo     "database": "silfae",>> config.json
echo     "username": "root",>> config.json
echo     "password": "",>> config.json
echo     "charset": "latin1">> config.json
echo   },>> config.json
echo   "paths": {>> config.json
echo     "base_dir": ".",>> config.json
echo     "darms_dir": "darms",>> config.json
echo     "output_dir": "inserts",>> config.json
echo     "temp_dir": "temp">> config.json
echo   },>> config.json
echo   "sql": {>> config.json
echo     "encoding": "latin1",>> config.json
echo     "batch_size": 100,>> config.json
echo     "use_transaction": true,>> config.json
echo     "use_ignore": true>> config.json
echo   },>> config.json
echo   "logging": {>> config.json
echo     "level": "info",>> config.json
echo     "format": "text",>> config.json
echo     "output_file": "">> config.json
echo   }>> config.json
echo }>> config.json
    
call :print_message "✅ Arquivo config.json criado!" "%GREEN%"
goto :eof

REM Função para compilar o projeto
:build_project
call :print_message "🔨 Compilando o projeto..." "%BLUE%"
    
REM Compilar
go build -ldflags="-X main.version=%VERSION%" -o build\%BINARY_NAME%.exe .
    
call :print_message "✅ Executável criado: build\%BINARY_NAME%.exe" "%GREEN%"
goto :eof

REM Função para executar testes
:run_tests
call :print_message "🧪 Executando testes..." "%BLUE%"
    
REM Executar testes
go test -v ./...
    
call :print_message "✅ Testes concluídos!" "%GREEN%"
goto :eof

REM Função para executar benchmarks
:run_benchmarks
call :print_message "⚡ Executando benchmarks..." "%BLUE%"
    
REM Executar benchmarks
go test -bench=. ./...
    
call :print_message "✅ Benchmarks concluídos!" "%GREEN%"
goto :eof

REM Função para criar arquivo .gitkeep
:create_gitkeep
call :print_message "📝 Criando arquivos .gitkeep..." "%BLUE%"
    
if not exist "darms\.gitkeep" type nul > darms\.gitkeep
if not exist "inserts\.gitkeep" type nul > inserts\.gitkeep
    
call :print_message "✅ Arquivos .gitkeep criados!" "%GREEN%"
goto :eof

REM Função para mostrar informações do sistema
:show_system_info
call :print_message "💻 Informações do sistema:" "%CYAN%"
echo   Sistema Operacional: Windows
for /f "tokens=*" %%i in ('go version') do echo   Go Version: %%i
for /f "tokens=*" %%i in ('git --version') do echo   Git Version: %%i
echo   Arquitetura: %PROCESSOR_ARCHITECTURE%
echo   Versão do Windows: %OS%
goto :eof

REM Função para mostrar próximos passos
:show_next_steps
call :print_message "🎉 Instalação concluída com sucesso!" "%GREEN%"
echo.
call :print_message "📋 Próximos passos:" "%YELLOW%"
echo   1. Coloque PDFs dos DARMs na pasta darms\
echo   2. Execute: build\%BINARY_NAME%.exe
echo   3. Verifique os arquivos gerados na pasta inserts\
echo.
call :print_message "🛠️ Comandos úteis:" "%YELLOW%"
echo   make help        - Mostra ajuda do Makefile
echo   make test        - Executa testes
echo   make run         - Executa o programa
echo   make build       - Recompila o executável
echo   make clean       - Remove arquivos de build
echo.
call :print_message "📚 Documentação:" "%YELLOW%"
echo   README_Go.md     - Documentação completa
echo   exemplo_uso.go   - Exemplos de uso
goto :eof

REM Função principal
:main
call :print_message "🏛️ Processador de DARMs - Versão Go %VERSION%" "%CYAN%"
call :print_message "🚀 Iniciando instalação..." "%BLUE%"
    
REM Verificar requisitos
call :check_requirements
if %errorlevel% neq 0 exit /b 1
    
REM Mostrar informações do sistema
call :show_system_info
    
REM Instalar dependências
call :install_dependencies
    
REM Criar diretórios
call :create_directories
    
REM Criar arquivo de configuração
call :create_config
    
REM Criar arquivos .gitkeep
call :create_gitkeep
    
REM Compilar projeto
call :build_project
    
REM Executar testes
call :run_tests
    
REM Executar benchmarks
call :run_benchmarks
    
REM Mostrar próximos passos
call :show_next_steps
goto :eof

REM Verificar se o script está sendo executado no diretório correto
if not exist "go.mod" (
    call :print_message "❌ Erro: Execute este script no diretório raiz do projeto!" "%RED%"
    call :print_message "💡 Certifique-se de que o arquivo go.mod existe no diretório atual." "%YELLOW%"
    exit /b 1
)

REM Executar função principal
call :main
pause 