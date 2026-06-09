## Summary

- Fixes Go installation by using the canonical module path and adding `cmd/abacatepay` so Go users can install the expected `abacatepay` binary.
- Narrows the CLI to the requested supported surface: auth, webhook listening, and dev-mode payment simulation.
- Updates payment simulation to API v2: `POST /v2/transparents/simulate-payment?id=<charge-id>`.
- Removes unsupported commands/features from the current codebase until there is documented API v2 parity.
- Updates README installation and usage docs.

Closes #32.

## Validation

- `go mod tidy`
- `go test ./...`
- `go build ./...`
- `GOBIN=/tmp/abacatepay-cli-bin go install .` and root CLI help smoke test
- `GOBIN=/tmp/abacatepay-bin go install ./cmd/abacatepay` and `abacatepay` help smoke test

## Notes

The `--local` flag is kept as a backwards-compatible no-op because API v2 uses the same public endpoint for production and dev mode; the API key determines the environment.

The old root Go install path still works and produces `abacatepay-cli`; the documented Go path now installs `abacatepay`.
