# Dockerfile para Processador de DARMs em Go

# Estágio de build
FROM golang:1.21-alpine AS builder

# Instalar dependências do sistema
RUN apk add --no-cache git ca-certificates tzdata

# Definir diretório de trabalho
WORKDIR /app

# Copiar arquivos de dependências
COPY go.mod go.sum ./

# Baixar dependências
RUN go mod download

# Copiar código fonte
COPY . .

# Compilar o aplicativo
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o darm-processor .

# Estágio final
FROM alpine:latest

# Instalar ca-certificates para HTTPS
RUN apk --no-cache add ca-certificates tzdata

# Criar usuário não-root
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Definir diretório de trabalho
WORKDIR /app

# Copiar o binário do estágio de build
COPY --from=builder /app/darm-processor .

# Copiar arquivos de configuração (se existirem)
COPY --from=builder /app/config.json* ./

# Criar diretórios necessários
RUN mkdir -p darms inserts && \
    chown -R appuser:appgroup /app

# Mudar para usuário não-root
USER appuser

# Expor porta (se necessário)
# EXPOSE 8080

# Definir variáveis de ambiente
ENV GO_ENV=production
ENV TZ=America/Sao_Paulo

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ["./darm-processor", "--health-check"] || exit 1

# Comando padrão
CMD ["./darm-processor"]

# Labels
LABEL maintainer="rodrigosardinha"
LABEL version="1.0.0"
LABEL description="Processador de DARMs em Go"
LABEL org.opencontainers.image.source="https://github.com/rodrigosardinha/gerador-query-darm" 