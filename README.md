# 🏛️ Processador de DARMs - Versão Go

Esta pasta contém a versão **Go** do Processador de DARMs, oferecendo melhor performance e um executável standalone.

## 🚀 Instalação Rápida

### Linux/macOS
```bash
cd goversion
chmod +x install.sh
./install.sh
```

### Windows
```cmd
cd goversion
install.bat
```

### Manual
```bash
cd goversion
go mod tidy
go build -o darm-processor
```

## 🎯 Como Usar

### Execução Básica
```bash
cd goversion
./darm-processor
```

### Com Makefile
```bash
cd goversion
make run
```

### Com Docker
```bash
cd goversion
docker-compose up
```

## 📁 Estrutura dos Arquivos

```
goversion/
├── 📄 main.go                    # Ponto de entrada
├── 📄 darm_processor.go          # Processador principal
├── 📄 config.go                  # Configurações
├── 📄 utils.go                   # Utilitários
├── 📄 test_darm_processor.go     # Testes
├── 📄 exemplo_uso.go             # Exemplos
├── 📄 go.mod                     # Dependências
├── 📄 Makefile                   # Automação
├── 📄 Dockerfile                 # Containerização
├── 📄 docker-compose.yml         # Orquestração
├── 📄 install.sh                 # Instalação Linux/macOS
├── 📄 install.bat                # Instalação Windows
├── 📄 config.json                # Configuração padrão
├── 📄 README_Go.md               # Documentação completa
└── 📄 RESUMO_VERSAO_GO.md        # Resumo técnico
```

## ⚡ Vantagens da Versão Go

- **5-10x mais rápido** que a versão Python
- **3-5x menos uso de memória**
- **Executável standalone** - não requer Go instalado
- **Tipagem estática** - menos bugs
- **Concorrência nativa** - suporte a goroutines

## 🧪 Testes

```bash
cd goversion
go test ./...
go test -bench=. ./...
```

## 📚 Documentação Completa

- **README_Go.md** - Documentação detalhada
- **RESUMO_VERSAO_GO.md** - Resumo técnico
- **exemplo_uso.go** - Exemplos práticos

## 🔧 Configuração

O arquivo `config.json` contém todas as configurações. Edite conforme necessário.

## 🐳 Docker

```bash
cd goversion
docker build -t darm-processor:1.0.0 .
docker run -v $(pwd)/darms:/app/darms -v $(pwd)/inserts:/app/inserts darm-processor:1.0.0
```

---

**🏛️ Processador de DARMs - Versão Go 1.0.0**

*Para documentação completa, consulte README_Go.md* 