# AGENTS.md - marketsurge-go

Go client library for the MarketSurge GraphQL API. Pure library (no CLI, no binaries). Module path: `github.com/major/marketsurge-go`. Requires Go 1.26+.

Unofficial project, not affiliated with Investor's Business Daily or MarketSurge.

## Layout

```text
marketsurge-go/
  marketsurge/           # Public library package (all source + tests)
    browserauth/         # Subpackage isolating kooky/SQLite browser deps
      testdata/          # Firefox cookie DB fixtures
    testdata/            # Fixture dirs per endpoint (2 files each)
  .github/
    copilot-instructions.md  # Copilot review config
    instructions/            # Additional Copilot instruction files
    workflows/               # CI, release, CodeQL, Scorecard
  .golangci.yml          # Linter config (maratori golden base)
  .goreleaser.yml        # Source-only release config
  .coderabbit.yaml       # CodeRabbit review config
  Makefile               # build, test, lint, vuln, clean, release
  renovate.json          # Dependency management
```

Not a standard Go app layout. No `cmd/`, `internal/`, `pkg/` directories. No `main.go`, no `func main()`, no `init()` functions.

## Build and test

```bash
make build           # go build ./...
make test            # gotestsum with -race -coverprofile, then coverage-check
make coverage-check  # Fails if total coverage < 90%
make lint            # golangci-lint run ./...
make vuln            # govulncheck ./...
make clean           # Remove coverage.out, junit.xml, dist/
make release VERSION=v0.3.0  # Must be on main, clean tree, test+lint pass, signed tag
```

Coverage floor: 90%. Tests produce `junit.xml` via gotestsum. Race detector is always on.

## Linting

Config: `.golangci.yml` (v2 format, based on [maratori golden config](https://github.com/maratori/golangci-lint-config)).

Key settings:
- `golines` max line length: 120
- `goimports` local prefix: `github.com/major/marketsurge-go`
- `nolintlint`: requires specific linter name + explanation (except funlen, gocognit, golines)
- `funlen`: 100 lines / 50 statements max
- `gocognit`: min complexity 20
- `nakedret`: forbidden (max-func-lines: 0)
- `depguard`: blocks deprecated packages (golang/protobuf, satori/uuid, math/rand, log)
- `testpackage`: enforced (use `_test` package)
- `paralleltest`: enforced
- Many linters relaxed in `_test.go` files (bodyclose, dupl, errcheck, funlen, goconst, gosec, noctx)

## CI

GitHub Actions workflows:
- `ci.yml`: build, test (with race + coverage), lint, vuln check
- `release.yml`: GoReleaser source-only release
- `codeql.yml`: CodeQL security scanning
- `scorecard.yml`: OpenSSF Scorecard

## Release

GoReleaser v2, source-only (no binary targets). SHA256 checksums, cosign keyless signing. Release tags are signed GPG tags created via `make release VERSION=vX.Y.Z` on `main` only.

## Project invariants

- The `marketsurge` package must not import kooky, SQLite, or browser-store dependencies. Those belong in `marketsurge/browserauth/` only. An import smoke test enforces this at CI time.
- Library code must not call `os.Exit`, read environment variables, write user-facing output, or read config files.
- Every public method that performs an HTTP request takes `context.Context` as its first parameter.
- GraphQL operation names, query strings, variable structures, and response models must match captured MarketSurge browser traffic verbatim. Do not invent or simplify them.
- Typed errors (`StatusError`, `GraphQLError`, `DecodeError`, `BodyLimitError`) and predicate helpers (`IsAuthError`, `IsRateLimited`, `RetryAfter`, `IsBodyLimit`, `StatusCode`) are API contract.
- All exported identifiers must have useful Go doc comments, not restatements of the name.
- `//nolint` directives must name the specific linter and include a brief explanation.
- US English spelling enforced in comments and identifiers.

## Testing conventions

- stdlib `testing` only. No testify, no require, no assert.
- `t.Parallel()` on every top-level test and subtest.
- `t.Helper()` on all test helpers.
- `httptest.NewServer` with inline request-validation handlers.
- Table-driven subtests with `t.Run()`.
- Request body JSON verified for `operationName`, `variables`, and `query` fields.
- Fixture-based response testing from `testdata/` directories.
- Use `assertStringPtr` and `assertIntPtr` helpers where they exist.
- Coverage must stay at or above 90%.

## Security rules

- Never expose cookie values, JWT tokens, or browser profile paths in log statements, error messages, test output, or documentation examples.
- `NewSession` is for explicit session construction. `browserauth.SessionFromFirefox` is for desktop automation only, not server-side code.
- Browser-backed auth stays isolated in `marketsurge/browserauth/`.
- No silent fallback behavior around HTTP status handling, GraphQL error decoding, or body limits. Return clear typed errors.

## Dependencies

Direct:
- `github.com/browserutils/kooky` - Firefox cookie reading (browserauth only)
- `github.com/google/go-cmp` - Test comparisons

## Maintenance

Keep these files in sync when project conventions, structure, or rules change:
- `AGENTS.md` (this file and subdirectory copies)
- `README.md` (endpoint table, examples, contributing rules)
- `.github/copilot-instructions.md` (Copilot review instructions)
- `.coderabbit.yaml` (CodeRabbit review config and path instructions)
- `.github/instructions/*.instructions.md` (supplemental Copilot instructions)
