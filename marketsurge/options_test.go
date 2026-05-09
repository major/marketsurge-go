package marketsurge

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestOptionDefaults(t *testing.T) {
	t.Parallel()

	got := defaultClientConfig()

	want := clientConfig{
		httpClient:        nil,
		graphQLURL:        DefaultGraphQLURL,
		investorsBaseURL:  DefaultInvestorsBaseURL,
		userAgent:         DefaultUserAgent,
		responseBodyLimit: DefaultResponseBodyLimit,
		headers:           http.Header{},
	}

	if diff := cmp.Diff(want, got, cmp.AllowUnexported(clientConfig{})); diff != "" {
		t.Fatalf("defaultClientConfig() mismatch (-want +got):\n%s", diff)
	}
}

func TestOptionWithHTTPClient(t *testing.T) {
	t.Parallel()

	client := &http.Client{}
	var cfg clientConfig

	WithHTTPClient(client)(&cfg)

	if diff := cmp.Diff(client, cfg.httpClient); diff != "" {
		t.Fatalf("WithHTTPClient(client) mismatch (-want +got):\n%s", diff)
	}
}

func TestOptionWithGraphQLURL(t *testing.T) {
	t.Parallel()

	var cfg clientConfig

	WithGraphQLURL("https://example.com/graphql")(&cfg)

	if got, want := cfg.graphQLURL, "https://example.com/graphql"; got != want {
		t.Fatalf("WithGraphQLURL(url) graphQLURL = %q, want %q", got, want)
	}
}

func TestOptionWithInvestorsBaseURL(t *testing.T) {
	t.Parallel()

	var cfg clientConfig

	WithInvestorsBaseURL("https://example.com")(&cfg)

	if got, want := cfg.investorsBaseURL, "https://example.com"; got != want {
		t.Fatalf("WithInvestorsBaseURL(url) investorsBaseURL = %q, want %q", got, want)
	}
}

func TestOptionWithUserAgent(t *testing.T) {
	t.Parallel()

	var cfg clientConfig

	WithUserAgent("test-agent")(&cfg)

	if got, want := cfg.userAgent, "test-agent"; got != want {
		t.Fatalf("WithUserAgent(userAgent) userAgent = %q, want %q", got, want)
	}
}

func TestOptionWithResponseBodyLimit(t *testing.T) {
	t.Parallel()

	var cfg clientConfig

	WithResponseBodyLimit(123)(&cfg)

	if got, want := cfg.responseBodyLimit, int64(123); got != want {
		t.Fatalf("WithResponseBodyLimit(limit) responseBodyLimit = %d, want %d", got, want)
	}
}

func TestOptionWithJWT(t *testing.T) {
	t.Parallel()

	var cfg clientConfig

	WithJWT("token-123")(&cfg)

	provider, ok := cfg.jwtProvider.(StaticJWTProvider)
	if !ok {
		t.Fatalf("WithJWT(token) jwtProvider type = %T, want StaticJWTProvider", cfg.jwtProvider)
	}
	if got, want := provider.Token, "token-123"; got != want {
		t.Fatalf("WithJWT(token) provider.Token = %q, want %q", got, want)
	}

	gotToken, err := cfg.jwtProvider.JWT(context.Background())
	if err != nil {
		t.Fatalf("StaticJWTProvider.JWT() error = %v", err)
	}
	if got, want := gotToken, "token-123"; got != want {
		t.Fatalf("StaticJWTProvider.JWT() = %q, want %q", got, want)
	}
}

func TestOptionWithJWTProvider(t *testing.T) {
	t.Parallel()

	var cfg clientConfig
	provider := StaticJWTProvider{Token: "provider-token"}

	WithJWTProvider(provider)(&cfg)

	if diff := cmp.Diff(provider, cfg.jwtProvider); diff != "" {
		t.Fatalf("WithJWTProvider(provider) mismatch (-want +got):\n%s", diff)
	}
}

func TestOptionWithHeader(t *testing.T) {
	t.Parallel()

	var cfg clientConfig

	WithHeader("X-Test", "value")(&cfg)

	want := http.Header{}
	want.Set("X-Test", "value")
	if diff := cmp.Diff(want, cfg.headers); diff != "" {
		t.Fatalf("WithHeader(name, value) mismatch (-want +got):\n%s", diff)
	}
}

func TestOptionJWTProviderLastOneWins(t *testing.T) {
	t.Parallel()

	var cfg clientConfig

	WithJWT("first")(&cfg)
	WithJWTProvider(StaticJWTProvider{Token: "second"})(&cfg)

	provider, ok := cfg.jwtProvider.(StaticJWTProvider)
	if !ok {
		t.Fatalf("jwtProvider type = %T, want StaticJWTProvider", cfg.jwtProvider)
	}
	if got, want := provider.Token, "second"; got != want {
		t.Fatalf("jwtProvider.Token = %q, want %q", got, want)
	}
}

func TestStaticJWTProvider(t *testing.T) {
	t.Parallel()

	provider := StaticJWTProvider{Token: "static-token"}

	if _, ok := any(provider).(JWTProvider); !ok {
		t.Fatal("StaticJWTProvider does not satisfy JWTProvider")
	}

	token, err := provider.JWT(context.Background())
	if err != nil {
		t.Fatalf("StaticJWTProvider.JWT() error = %v", err)
	}
	if got, want := token, "static-token"; got != want {
		t.Fatalf("StaticJWTProvider.JWT() = %q, want %q", got, want)
	}
}
