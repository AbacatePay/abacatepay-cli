# abacatepay-cli

CLI para receber webhooks do AbacatePay em ambiente de desenvolvimento local.

## Instalação

```bash
go build -o abacatepay-cli ./cmd
```

## Uso Rápido

```bash
abacatepay-cli login
```

Após autenticar, a CLI pergunta a URL para encaminhar webhooks e já começa a escutar.

## Ambientes

A CLI suporta dois ambientes: **produção** (padrão) e **teste**.

### Produção (Padrão)

```bash
# API: https://api.abacatepay.com
# WebSocket: wss://ws.abacatepay.com/ws

abacatepay-cli login
```

### Servidor de Teste

Para desenvolvimento e testes, use a flag `-l`:

```bash
# API: http://191.252.202.128:8080
# WebSocket: ws://191.252.202.128:8080/ws

abacatepay-cli -l login
```

## Comandos

| Comando  | Descrição                                     |
| -------- | --------------------------------------------- |
| `login`  | Autenticar e iniciar listener                 |
| `logout` | Remover credenciais                           |
| `status` | Verificar se está autenticado                 |
| `listen` | Escutar webhooks (requer autenticação prévia) |

### Flags do Login

| Flag          | Descrição                             |
| ------------- | ------------------------------------- |
| `-f <url>`    | URL para encaminhar (pula o prompt)   |
| `--no-listen` | Apenas autentica, não inicia listener |

### Flags Globais

| Flag            | Descrição              |
| --------------- | ---------------------- |
| `-v, --verbose` | Logs detalhados        |
| `-l, --local`   | Usar servidor de teste |
| `--version`     | Versão da CLI          |

## Autenticação

A CLI usa OAuth2 Device Flow:

1. Execute `abacatepay-cli login`
2. Abra a URL exibida no navegador
3. Autorize o acesso na sua conta AbacatePay
4. A CLI detecta a autorização automaticamente
5. Informe a URL para encaminhar webhooks (ou pressione Enter para usar o padrão)

O token é armazenado no keyring nativo do sistema operacional:

- **macOS**: Keychain
- **Linux**: gnome-keyring ou kwallet
- **Windows**: Credential Manager

## Eventos Disponíveis

| Evento            | Descrição            |
| ----------------- | -------------------- |
| `billing.paid`    | Pagamento confirmado |
| `withdraw.done`   | Saque concluído      |
| `withdraw.failed` | Saque falhou         |

Consulte a [documentação de webhooks](https://docs.abacatepay.com/pages/webhooks) para detalhes dos payloads.

## Exemplo Completo

```bash
# Terminal 1: Servidor local
node server.js

# Terminal 2: CLI
abacatepay-cli login
# Abre URL no navegador, autoriza, informa URL do webhook
# Webhooks são encaminhados automaticamente

# Para parar: Ctrl+C
```

## Build Multi-plataforma

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o abacatepay-cli-linux-amd64 ./cmd

# macOS Intel
GOOS=darwin GOARCH=amd64 go build -o abacatepay-cli-darwin-amd64 ./cmd

# macOS Apple Silicon
GOOS=darwin GOARCH=arm64 go build -o abacatepay-cli-darwin-arm64 ./cmd

# Windows
GOOS=windows GOARCH=amd64 go build -o abacatepay-cli.exe ./cmd
```

## Troubleshooting

### "não autenticado"

Execute `abacatepay-cli login`.

### Token não salvo (Linux)

Instale o gnome-keyring:

```bash
# Debian/Ubuntu
sudo apt install gnome-keyring

# Fedora
sudo dnf install gnome-keyring
```

### WebSocket não conecta

1. Verifique autenticação: `abacatepay-cli status`
2. Use modo verbose: `abacatepay-cli -v login`

## Estrutura do Projeto

```
abacatepay-cli/
├── cmd/
│   └── main.go
├── internal/
│   ├── config/
│   ├── client/
│   ├── auth/
│   ├── logger/
│   └── webhook/
└── go.mod
```

## Logs

Os logs são salvos em `~/.abacatepay/logs/`:

- **abacatepay.log** - Log geral (JSON)
- **transactions.log** - Webhooks recebidos e encaminhados

Rotação automática: 10 MB por arquivo, 5 backups, 30 dias de retenção.

### Análise com jq

```bash
# Erros
cat ~/.abacatepay/logs/abacatepay.log | jq 'select(.level=="ERROR")'

# Webhooks recebidos
cat ~/.abacatepay/logs/transactions.log | jq 'select(.msg=="webhook_received")'

# Tempo médio de encaminhamento
cat ~/.abacatepay/logs/transactions.log | jq 'select(.msg=="webhook_forwarded") | .duration_ms' | jq -s 'add/length'
```

## Dependências

- [cobra](https://github.com/spf13/cobra) - Framework CLI
- [resty](https://github.com/go-resty/resty) - HTTP client
- [keyring](https://github.com/99designs/keyring) - Armazenamento seguro
- [gorilla/websocket](https://github.com/gorilla/websocket) - WebSocket client
- [lumberjack](https://github.com/natefinch/lumberjack) - Rotação de logs

## Licença

MIT

## Suporte

- Documentação: https://docs.abacatepay.com/
- Issues: https://github.com/AbacatePay/abacatepay-cli/issues
