# marketsurge-go review instructions

marketsurge-go is a Go client library for authenticated MarketSurge GraphQL endpoints. The public API exposes typed request parameters, typed response structs, context propagation on every request method, explicit session and JWT handling, and typed errors with predicate helpers.

Focus reviews on bugs, security issues, data loss, broken API contracts, and project conventions. Don't nitpick formatting or style that gofmt and golangci-lint already enforce.

## Project invariants

- Module path: `github.com/major/marketsurge-go`
- Public packages: root package and `browserauth/`
- The root package must stay free of browser-store, SQLite, and kooky dependencies. Those belong exclusively in `browserauth/`.
- Library code must not call `os.Exit`, write user-facing output, read hidden config files, or inspect environment variables.
- Every public method that performs an HTTP request takes `context.Context` as its first parameter.
- GraphQL operation names, query strings, variable structures, and response models must match captured MarketSurge browser traffic. Don't invent or simplify them.
- Preserve the typed errors and predicate helpers from `errors.go`: `IsAuthError`, `IsRateLimited`, `RetryAfter`, `IsBodyLimit`, `StatusCode`.
- All exported identifiers must have Go doc comments that are useful, not just restatements of the name.

## Security and session safety

- Flag any code that exposes cookie values, JWT tokens, or browser profile paths in log statements, error messages, test output, or documentation examples.
- Prefer explicit sessions via `NewSession` for service use. `browserauth.SessionFromFirefox` is for desktop automation only, not server-side code.
- Browser-backed auth must stay isolated to the `browserauth/` package. The root package must not import kooky or SQLite.
- Avoid silent fallback behavior around HTTP status handling, GraphQL error decoding, and body limits. Return clear typed errors instead.

## API boundary expectations

- Each endpoint method wraps `doGraphQL` with typed request and response structs specific to that endpoint.
- GraphQL queries embedded as string constants must match captured traffic verbatim. Don't paraphrase or reformat them.
- Do not silently mutate caller-provided cookies, headers, slices, or maps. Copy defensively at boundaries.

## Testing expectations

- Use stdlib `testing` only. No testify, no require, no assert packages.
- Every top-level test and subtest must call `t.Parallel()`.
- Use `t.Helper()` on all test helper functions.
- Use `httptest.NewServer()` for HTTP mocks with inline request validation handlers.
- Prefer table-driven subtests with `t.Run()`.
- Verify request body JSON structure: check `operationName`, `variables`, and `query` fields.
- Use fixture-based response testing for GraphQL responses.
- Use `assertStringPtr` and `assertIntPtr` helpers instead of bare if-checks where those helpers exist.
- Coverage should stay at or above 90%.

## Build and lint expectations

- CI runs `go test` with `-race` and `-coverprofile`, `golangci-lint run`, and `go build ./...`.
- GoReleaser is configured for source-only release (library, no `main` package, no binaries).
- `//nolint` directives must name the specific linter and include a brief explanation.
- US English spelling is enforced in comments and identifiers.
