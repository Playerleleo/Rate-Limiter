# Rate Limiter em Go

Este é um rate limiter implementado em Go que permite limitar requisições baseado em IP ou token de acesso.

## Funcionalidades

- Limitação de requisições por IP
- Limitação de requisições por token de acesso
- Configuração via variáveis de ambiente
- Armazenamento em Redis
- Middleware para fácil integração com servidores HTTP

## Requisitos

- Go 1.16+
- Docker e Docker Compose
- Redis

## Configuração

As configurações são feitas através do arquivo `.env` na raiz do projeto:

```env
# Configurações do Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# Configurações de Limite por IP
IP_REQUESTS_PER_SECOND=5
IP_BLOCK_DURATION_SECONDS=300

# Configurações de Limite por Token
TOKEN_REQUESTS_PER_SECOND=10
TOKEN_BLOCK_DURATION_SECONDS=300

# Configurações do Servidor
SERVER_PORT=8080
```

## Executando o Projeto

1. Inicie o Redis usando Docker Compose:
```bash
docker-compose up -d
```

2. Execute o servidor:
```bash
go run main.go
```

## Uso

O rate limiter pode ser usado de duas formas:

1. Limitação por IP: Aplica-se automaticamente a todas as requisições
2. Limitação por Token: Adicione o header `API_KEY: <seu-token>` na requisição

### Exemplo de Requisição com Token

```bash
curl -H "API_KEY: seu-token-aqui" http://localhost:8080/
```

## Respostas

- 200: Requisição aceita
- 429: Limite de requisições excedido
- 500: Erro interno do servidor