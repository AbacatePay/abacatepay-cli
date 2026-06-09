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

Os releases publicam binários para Linux, macOS e Windows nas arquiteturas `amd64` e `arm64`.

## Uso rápido

```bash
abacatepay login --key "$ABACATEPAY_API_KEY"
abacatepay listen --forward-to http://localhost:3000/webhooks/abacatepay
abacatepay payments simulate <charge-id>
```

Use `abacatepay <command> -h` para ver as flags de cada comando.

## Autenticação

A CLI aceita login por chave de API ou pelo fluxo de device login:

```bash
# Recomendado para automações e ambientes sem navegador
abacatepay login --key "$ABACATEPAY_API_KEY"

# Fluxo interativo pelo navegador
abacatepay login
```

A chave fica salva no keyring nativo do sistema operacional. Se nenhum nome for informado com `--name`, a CLI usa o perfil `default`.

A API v2 usa o mesmo endpoint público para produção e desenvolvimento. O ambiente é determinado pela chave utilizada: chaves de dev mode simulam transações; chaves de produção operam em produção.

## Webhooks

Receba eventos da AbacatePay e encaminhe para sua aplicação local:

```bash
abacatepay listen --forward-to http://localhost:3000/webhooks/abacatepay
```

Flags úteis:

| Flag | Descrição | Padrão |
| ---- | --------- | ------ |
| `--forward-to` | URL local que receberá os eventos | `http://localhost:3000/webhooks/abacatepay` |
| `--mock` | Gera webhooks locais sem conectar na API | `false` |

A CLI assina os eventos encaminhados com o header `X-Abacate-Signature`, facilitando testes locais de validação de assinatura.

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
