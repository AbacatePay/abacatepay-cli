<h1 align="center">AbacatePay CLI</h1>

<p align="center">
  CLI oficial do AbacatePay para desenvolvimento local, webhooks e testes rápidos via terminal.
</p>

<p align="center">
  <sub>
    <samp>
      <a href="#instalação">Instalação</a> •
      <a href="#uso-rápido">Uso</a> •
      <a href="#autenticação">Autenticação</a> •
      <a href="#ambientes">Ambientes</a>
    </samp>
  </sub>
</p>

<h2 align="center">Instalação</h2>

<h3 align="center">Go (recomendado)</h3>

```bash
go install github.com/AbacatePay/abacatepay-cli@latest
```

<p align="center">O binário automáticamente será instalado como <em>abacatepay</em> com um aliás <em>abkt</em>.</p>

<h3 align="center">Homebrew (macOS / Linux)</h3>

```bash
brew install --build-from-source github.com/AbacatePay/abacatepay-cli
```

<h2 align="center">Uso Rápido</h2>

```bash
abacatepay login
```

<p align="center">Após a autenticação <em>(OAuth Device Flow)</em>, você deve informar a URL do seu servidor local, então a CLI encaminhará todos os webhooks para você.</p>

<p align="center">Todos os comandos podem ser usados com a seguinte sintaxe:</p>

```bash
abacatepay <command> [...flags] [...args]
```

<p align="center">Use a flag <code>-h</code> para obter informações detalhadas sobre cada comando.</p>

<h2 align="center">Ambientes</h2>

<p align="center">Atualmente a CLI suporta dois ambientes, o de <em>produção</em> (Padrão) e <em>teste</em>.</p>

<h3 align="center">Produção (Padrão)</h3>

```bash
# API: https://api.abacatepay.com
# WebSocket: wss://ws.abacatepay.com/ws

abacatepay login
```

<h3 align="center">Servidor de Teste</h3>

<p align="center">Para usar o modo de desenvolvimento, use a flag <code>-l</code></p>


```bash
# API: http://191.252.202.128:8080
# WebSocket: ws://191.252.202.128:8080/ws

abacatepay login -l
```

<h2 align="center">Autenticação</h2>

<p align="center">A CLI usa <em>OAuth2 Device Flow</em>, sem necessidade de copiar tokens manualmente, apenas seguindo o fluxo abaixo</p>

1. Use `abacatepay login`
2. Abra a URL exibida no navegador
3. Autorize o acesso na sua conta AbacatePay
4. A CLI detecta a autorização automaticamente
5. Informe a URL para encaminhar webhooks (Ou pressione a tecla <em>Enter</em> para usar o padrão)

<h3 align="center">Armazenamento</h3>

<p align="center">O token da sua conta é armazenado com segurança no <em>keyring nativo</em> do seu sistema operacional:</p>

- **macOS**: Keychain
- **Linux**: gnome-keyring ou kwallet
- **Windows**: Credential Manager

<p align="center">Caso o token não consiga ser salvo, você deverá instalar o keyring no seu sistema operacional (Linux)</p>

```bash
# Debian/Ubuntu
sudo apt install gnome-keyring

# Fedora
sudo dnf install gnome-keyring
```

<h2 align="center">Logs</h2>

<p align="center">Todos os logs são salvos ni caminho <code>~/.abacatepay/logs/</code> com uma rotação automática de 10mb por arquivo, 5 backups e 30 dias de retenção</p>

- **abacatepay.log** - Log geral (JSON)
- **transactions.log** - Webhooks recebidos e encaminhados

<h3 align="center">Use jq para analisar os logs</h3>

```bash
# Erros
cat ~/.abacatepay/logs/abacatepay.log | jq 'select(.level=="ERROR")'

# Webhooks recebidos
cat ~/.abacatepay/logs/transactions.log | jq 'select(.msg=="webhook_received")'

# Tempo médio de encaminhamento
cat ~/.abacatepay/logs/transactions.log | jq 'select(.msg=="webhook_forwarded") | .duration_ms' | jq -s 'add/length'
```

<h2 align="center">Documentação</h2>

<p align="center">Para uma documentação completa, veja a <a href="https://docs.abacatepay.com/pages/cli">documentação oficial</a>.</p>
