# AGENTS.md - marketsurge package

Public library package providing a typed Go client for the MarketSurge GraphQL API.

## Architecture

### Client lifecycle

`NewClient` accepts functional options (`Option` funcs) and returns a `*Client`. The client is safe for concurrent use after construction. All HTTP requests flow through `doGraphQL`, which handles JSON marshaling, header injection, response body limiting, and GraphQL error decoding.

Key types:
- `Client` - holds HTTP client, URLs, JWT provider, headers, body limit
- `clientConfig` - internal config struct populated by `Option` functions
- `Session` - carries defensively copied browser cookies for JWT exchange
- `JWTProvider` - interface for token refresh; `StaticJWTProvider` is the simple case

### Functional options

Options defined in `options.go`: `WithHTTPClient`, `WithGraphQLURL`, `WithInvestorsBaseURL`, `WithUserAgent`, `WithResponseBodyLimit`, `WithJWT`, `WithJWTProvider`, `WithHeader`. All operate on `*clientConfig`. Invalid values fail at `NewClient` time, not at request time.

### Error types

Defined in `errors.go`, these are API contract:
- `StatusError` - non-2xx HTTP response, carries status code, body, and headers
- `GraphQLError` - wraps `[]GraphQLFieldError`, implements multi-unwrap
- `DecodeError` - JSON unmarshal failure, wraps underlying error
- `BodyLimitError` - response exceeded configured limit

Predicate helpers: `IsAuthError`, `IsRateLimited`, `RetryAfter`, `IsBodyLimit`, `StatusCode`, `IsStatusCode`.

### GraphQL transport

`doGraphQL` is the single transport path for all endpoint methods. Flow: get JWT from provider, marshal `GraphQLRequest`, create HTTP request with context, set headers, send, check status, read body with limit, decode `GraphQLResponse`, extract data into target struct.

`GraphQLRequest[V]` and `GraphQLResponse[T]` are generic envelope types in `graphql.go`.

### Auth flow

`ExchangeJWT` in `auth.go` exchanges browser cookies for a JWT via the investors.com `/client` endpoint. Returns `ClientInfoResponse` with JWT, login status, and user name fields. Does not require a JWT provider on the client.

## File organization

Each endpoint gets its own file pair:

| Source file | Test file | Endpoint methods |
|---|---|---|
| `market_data.go` | `market_data_test.go` | `OtherMarketData` |
| `fundamentals.go` | `fundamentals_test.go` | `Fundamentals` |
| `ownership.go` | `ownership_test.go` | `Ownership` |
| `relative_strength.go` | `relative_strength_test.go` | `RSRatingRIPanel` |
| `charts.go` | `charts_test.go` | `ChartMarketData`, `ChartMarketDataWeekly` |
| `chart_markups.go` | `chart_markups_test.go` | `FetchChartMarkups` |
| `watchlists.go` | `watchlists_test.go` | `GetAllWatchlistNames`, `FlaggedSymbols` |
| `screens.go` | `screens_test.go` | `MarketDataAdhocScreen`, `RunScreen`, `Screens` |
| `coach_tree.go` | `coach_tree_test.go` | `CoachTree`, `IndustryGroupRS` |

Infrastructure files: `client.go`, `options.go`, `errors.go`, `graphql.go`, `auth.go`, `session.go`, `doc.go`.

Special test: `import_smoke_test.go` verifies the root package does not pull in kooky/SQLite/browser-store transitive deps.

### Adding a new endpoint

1. Create `<endpoint>.go` with request/response structs, a `New<Endpoint>Request` constructor, and a method on `*Client` that calls `doGraphQL`.
2. GraphQL operation name, query string, and variable struct must match captured browser traffic exactly.
3. Create `<endpoint>_test.go` with table-driven subtests, `httptest.NewServer`, and fixture files under `testdata/<endpoint>/`.
4. Each `testdata/<endpoint>/` directory holds two files: the request fixture and the response fixture.
5. Update the endpoint table in `README.md`.

## Dependency boundary

This package must not import:
- `github.com/browserutils/kooky` or any kooky subpackage
- Any SQLite driver
- Any keyring or dbus package

The import smoke test (`import_smoke_test.go`) enforces this via `go list -deps`. Browser-store dependencies belong exclusively in the `browserauth` subpackage.

## Testing patterns

- All tests use the `marketsurge_test` package (external test package enforced by linter).
- Every test and subtest calls `t.Parallel()`.
- HTTP mocks use `httptest.NewServer` with inline handlers that validate request method, path, headers, and JSON body structure.
- Request body validation checks `operationName`, `variables`, and `query` fields match expected values.
- Response fixtures live in `testdata/<endpoint>/` directories.
- Use `assertStringPtr` and `assertIntPtr` helpers for pointer field assertions.
- `cmp.Diff` from `github.com/google/go-cmp` for struct comparisons.
