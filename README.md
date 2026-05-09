# marketsurge-go

Go client library for the [MarketSurge](https://marketsurge.investors.com/) GraphQL API. Covers market data, charts, watchlists, screens, and more with typed responses, functional options, and structured error handling.

> [!NOTE]
> This project is unofficial and is not affiliated with, sponsored by, or endorsed by Investor's Business Daily or MarketSurge in any way.

## Installation

```bash
go get github.com/major/marketsurge-go
```

Requires Go 1.26 or later.

## Quick start

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/major/marketsurge-go/marketsurge"
)

func main() {
	client, err := marketsurge.NewClient(
		marketsurge.WithJWT("your-jwt-token"),
	)
	if err != nil {
		log.Fatal(err)
	}

	req := marketsurge.NewOtherMarketDataRequest("AAPL", "MSFT")
	resp, err := client.OtherMarketData(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range resp.MarketData {
		if item.OriginRequest != nil && item.OriginRequest.Symbol != nil {
			fmt.Println(*item.OriginRequest.Symbol)
		}
	}
}
```

## Auth

MarketSurge uses a JWT that you exchange from browser session cookies. The flow is:

1. Log into MarketSurge in your browser.
2. Extract the session cookies.
3. Call `ExchangeJWT` to get a JWT.
4. Pass the JWT to `NewClient` with `WithJWT`.

```go
package main

import (
	"context"
	"log"
	"net/http"

	"github.com/major/marketsurge-go/marketsurge"
)

func main() {
	// Build a session from cookies you extracted from your browser.
	cookies := []*http.Cookie{
		{Name: "IBD_USER", Value: "..."},
		{Name: "IBD_SESSION", Value: "..."},
	}
	session := marketsurge.NewSession(cookies)

	// Create a client without a JWT first -- ExchangeJWT doesn't need one.
	client, err := marketsurge.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	info, err := client.ExchangeJWT(context.Background(), session)
	if err != nil {
		log.Fatal(err)
	}

	// Now create a client with the JWT for API calls.
	client, err = marketsurge.NewClient(
		marketsurge.WithJWT(info.JWT),
	)
	if err != nil {
		log.Fatal(err)
	}

	_ = client
}
```

`ExchangeJWT` returns a `*ClientInfoResponse` with `JWT`, `GivenName`, `FamilyName`, and `IsLoggedIn`. It returns an error if the cookies are invalid or the account isn't logged in.

For long-running processes, implement `JWTProvider` and pass it with `WithJWTProvider` to refresh the token on demand.

## browserauth

The `browserauth` subpackage reads cookies directly from a local Firefox profile, so you don't have to copy them manually.

```go
package main

import (
	"context"
	"log"

	"github.com/major/marketsurge-go/marketsurge"
	"github.com/major/marketsurge-go/marketsurge/browserauth"
)

func main() {
	// Pass the path to your Firefox profile directory.
	// On Linux this is typically ~/.mozilla/firefox/<profile>.default-release/
	session, err := browserauth.SessionFromFirefox("/home/user/.mozilla/firefox/abc123.default-release")
	if err != nil {
		log.Fatal(err)
	}

	client, err := marketsurge.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	info, err := client.ExchangeJWT(context.Background(), *session)
	if err != nil {
		log.Fatal(err)
	}

	client, err = marketsurge.NewClient(marketsurge.WithJWT(info.JWT))
	if err != nil {
		log.Fatal(err)
	}

	_ = client
}
```

`SessionFromFirefox` reads `cookies.sqlite` from the given profile directory. Close Firefox before calling it, or you'll get `ErrCookieDBLocked`.

## Endpoint coverage

| Method | Wire operation | File | Notes |
|--------|---------------|------|-------|
| `OtherMarketData` | `OtherMarketData` | `marketsurge/market_data.go` | |
| `Fundamentals` | `FundermentalDataBox` | `marketsurge/fundamentals.go` | upstream typo in wire name |
| `Ownership` | `Ownership` | `marketsurge/ownership.go` | |
| `RSRatingRIPanel` | `RSRatingRIPanel` | `marketsurge/relative_strength.go` | |
| `ChartMarketData` | `ChartMarketData` | `marketsurge/charts.go` | daily periodicity |
| `ChartMarketDataWeekly` | `ChartMarketData` | `marketsurge/charts.go` | weekly periodicity |
| `FetchChartMarkups` | `FetchChartMarkups` | `marketsurge/chart_markups.go` | |
| `GetAllWatchlistNames` | `GetAllWatchlistNames` | `marketsurge/watchlists.go` | |
| `FlaggedSymbols` | `FlaggedSymbols` | `marketsurge/watchlists.go` | |
| `MarketDataAdhocScreen` | `MarketDataAdhocScreen` | `marketsurge/screens.go` | |
| `RunScreen` | `RunScreen` | `marketsurge/screens.go` | |
| `Screens` | `Screens` | `marketsurge/screens.go` | |
| `CoachTree` | `CoachTree` | `marketsurge/coach_tree.go` | |
| `IndustryGroupRS` | `IndustryGroupRS` | `marketsurge/coach_tree.go` | |
| Phase 5 operations | various | - | not implemented - undiscoverable from agent sources |

Each endpoint has a `NewXxxRequest` constructor with sensible defaults. All methods take `context.Context` as the first argument.

## Error handling

HTTP errors come back as `*StatusError`. GraphQL-level errors come back as `*GraphQLError`. Use the helper functions to classify errors without type-asserting yourself:

```go
import (
	"fmt"
	"time"

	"github.com/major/marketsurge-go/marketsurge"
)

resp, err := client.OtherMarketData(ctx, req)
if err != nil {
	switch {
	case marketsurge.IsAuthError(err):
		// 401 or 403 -- re-authenticate and retry.
		fmt.Println("auth failed, re-exchange JWT")
	case marketsurge.IsRateLimited(err):
		// 429 -- back off before retrying.
		wait := marketsurge.RetryAfter(err)
		if wait == 0 {
			wait = 5 * time.Second
		}
		fmt.Printf("rate limited, retry after %s\n", wait)
	case marketsurge.IsBodyLimit(err):
		// Response exceeded the configured body limit.
		fmt.Println("response too large")
	default:
		// Check the raw status code for anything else.
		code := marketsurge.StatusCode(err)
		fmt.Printf("HTTP %d: %v\n", code, err)
	}
	return
}
```

`*GraphQLError` wraps one or more `GraphQLFieldError` values. Use `errors.As` or `errors.Is` to inspect individual field errors.

Response bodies are capped at 10 MiB by default. Override with `WithResponseBodyLimit`.

## Contributing

1. Write tests first. All new endpoints need table-driven tests with `httptest` fixtures.
2. Run `golangci-lint run` before opening a PR. The config is in `.golangci.yml`.
3. Run `go test ./...` and make sure the import smoke test passes.
4. Use [Conventional Commits](https://www.conventionalcommits.org/) for commit messages.
5. Keep the endpoint coverage table in this README up to date.

## License

[Apache License 2.0](LICENSE)
