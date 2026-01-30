## Summary

Adds `init` command that initializes new projects with interactive scaffolding, allowing choice between frameworks (Next.js/Elysia), linters (ESLint/Biome), and optional BetterAuth configuration.

## Related Issue

Closes (N/A)

## Why

The CLI needed a simple way to create new projects integrated with AbacatePay. Before, developers had to manually clone templates and configure everything by hand. The `init` command automates this entire process with an interactive onboarding flow.

## What changed

- Created `cmd/init.go` with interactive onboarding flow using prompts from the `style` package
- Created `internal/scaffold/config.go` with `Config` structure and validation methods
- Created `internal/scaffold/git.go` for repository clone operations
- Created `internal/scaffold/package.go` for `package.json` manipulation and merging
- Created `internal/scaffold/fs.go` for file and directory copy operations
- Created `internal/scaffold/scaffold.go` with scaffolding orchestration via `ProjectBuilder`
- Implemented layer-based composition system that combines base templates + linters + features (BetterAuth)

## Breaking changes

- [ ] Yes
- [X] No

## Checklist

- [ ] Docs updated
- [X] CI passing
- [X] I followed the CONTRIBUTING guidelines
- [X] I added or updated tests (if applicable)

## Additional context

The command uses the `github.com/albuquerquesz/abacatepay-templates` repository as the source for templates. The layer-based composition architecture allows adding new frameworks, linters, and features without creating a combinatorial explosion of templates.
