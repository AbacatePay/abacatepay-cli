<h1 align="center">AbacatePay CLI</h1>

<p align="center">
  CLI oficial da AbacatePay para autenticação, recebimento de webhooks e simulação de pagamentos em modo de desenvolvimento.
</p>

<p align="center">
  <sub>
    <samp>
      <a href="#instalação">Instalação</a> •
      <a href="#uso-rápido">Uso rápido</a> •
      <a href="#autenticação">Autenticação</a> •
      <a href="#webhooks">Webhooks</a> •
      <a href="#pagamentos-em-dev-mode">Pagamentos em dev mode</a>
    </samp>
  </sub>
</p>

## Instalação

### Script de instalação (Linux e macOS)

```bash
curl -fsSL https://abacatepay.com/install.sh | sh
```

Baixa o binário `abacatepay` do último release publicado, valida o checksum e instala em `/usr/local/bin` (use `ABACATEPAY_INSTALL_DIR=<dir>` para instalar em outro lugar).

### Go

```bash
go install github.com/AbacatePay/abacatepay-cli/cmd/abacatepay@latest
```

Isso instala o binário `abacatepay`.

O comando antigo também funciona, mas instala o binário com o nome do repositório:

```bash
go install github.com/AbacatePay/abacatepay-cli@latest
# binário: abacatepay-cli
```

### Binários prontos

Os releases publicam binários para Linux, macOS e Windows nas arquiteturas `amd64` e `arm64` em [Releases](https://github.com/AbacatePay/abacatepay-cli/releases).

## Uso rápido

```bash
abacatepay login
abacatepay listen --forward-to http://localhost:3000/webhooks/abacatepay
abacatepay payments simulate <charge-id>
```

Use `abacatepay <command> -h` para ver as flags de cada comando.

## Autenticação

```bash
abacatepay login
```

Abre um link no navegador para você aprovar o acesso da CLI na sua conta. O token da sessão fica salvo no keyring nativo do sistema operacional. Se nenhum nome for informado com `--name`, a CLI usa o perfil `default`.

Para sair: `abacatepay logout`.

## Webhooks

Receba eventos da AbacatePay e encaminhe para sua aplicação local:

```bash
abacatepay listen --forward-to http://localhost:3000/webhooks/abacatepay
```

`--forward-to` define a URL local que receberá os eventos (padrão: `http://localhost:3000/webhooks/abacatepay`).

A CLI só recebe eventos de contas/cobranças em `devMode` — produção nunca chega no relay. A CLI assina os eventos encaminhados com o header `X-Webhook-Signature`, usando a mesma chave pública HMAC documentada em [Segurança de Webhooks](https://docs.abacatepay.com/pages/webhooks/security#2-assinatura-hmac). Essa chave é fixa — a verificação que você já implementou seguindo a documentação funciona sem alteração.

## Pagamentos em dev mode

Simule o pagamento de uma cobrança transparente criada com uma chave de dev mode:

```bash
abacatepay payments simulate <charge-id>
```

Esse comando chama a API v2:

```http
POST /v2/transparents/simulate-payment?id=<charge-id>
```

Ele só funciona para cobranças em `devMode`. Em produção, o pagamento precisa acontecer pelo fluxo real.

## Output

A flag global `--output` permite trocar o formato de saída:

```bash
abacatepay payments simulate <charge-id> --output json
```

Formatos disponíveis: `text`, `json` e `table`.

## Escopo atual

Para manter a CLI simples e alinhada com a API v2, o escopo atual é:

- login/logout;
- receber webhooks;
- simular pagamento de cobranças transparentes em dev mode.

Recursos fora desse escopo foram removidos da codebase até existir paridade clara com a API v2 e documentação atualizada.
