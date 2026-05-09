---
applyTo: ".coderabbit.yaml"
---

# CodeRabbit configuration instructions for marketsurge-go

## Review trigger

Keep `auto_review.enabled: false`. Reviews run on demand via PR comments (`@coderabbitai review`), not automatically on every push. This avoids noise on work-in-progress commits.

## Review profile

`profile: chill` is intentional. The goal is catching real bugs and contract violations, not style nitpicks that golangci-lint already handles.

## Path instructions focus areas

Path instructions in `.coderabbit.yaml` are scoped to the concerns that matter most for this library:

**MarketSurge GraphQL API contracts** - GraphQL operation names, query strings, variable structures, and response models must match captured browser traffic. Any drift from the captured traffic is a bug, not a style issue. Flag it.

**Static catalogs** - `marketsurge/columns.go` and `marketsurge/report_screens.go` contain captured API data. Wire names, report screen identifiers, spacing, and capitalization must remain exact. Catalog functions should return defensive copies, and `ColumnName` request helpers should not normalize, trim, or reject names that upstream may already accept.

**Explicit session and JWT handling** - `ExchangeJWT` uses cookies to obtain a JWT. Sessions created via `NewSession` carry that JWT. Neither cookies nor JWT values should appear in error messages, log output, or test assertions. Flag any exposure.

**Browser-auth isolation** - kooky, SQLite, and Firefox profile access belong in `marketsurge/browserauth/` only. The `marketsurge` package must not import them. If a change adds such an import to the `marketsurge` package, that's a hard violation.

**Typed error classification** - `marketsurge/errors.go` defines the error taxonomy for this library. `IsAuthError`, `IsRateLimited`, `RetryAfter`, `IsBodyLimit`, and `StatusCode` are part of the public API. Don't remove or weaken them. New error conditions should fit into this taxonomy or extend it with a new predicate.

**Context propagation** - Every method that performs an HTTP request must accept `context.Context` as its first parameter and pass it to the underlying HTTP request. Missing context propagation is a correctness bug.

**Response body limits** - `doGraphQL` enforces a configurable body limit before reading the response. Changes that bypass or remove this limit should be flagged.

## What to skip

Don't flag: gofmt formatting, import ordering, variable naming that passes golangci-lint, comment style on unexported symbols, or test helper organization that doesn't affect correctness.
