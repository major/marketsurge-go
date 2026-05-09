package marketsurge

import (
	"context"
	"net/http"
)

const (
	DefaultGraphQLURL        = "https://shared-data.dowjones.io/gateway/graphql"
	DefaultInvestorsBaseURL  = "https://www.investors.com"
	DefaultUserAgent         = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:149.0) Gecko/20100101 Firefox/149.0"
	DefaultResponseBodyLimit = int64(10 * 1024 * 1024)
)

type JWTProvider interface {
	JWT(ctx context.Context) (string, error)
}

type StaticJWTProvider struct {
	Token string
}

func (p StaticJWTProvider) JWT(_ context.Context) (string, error) {
	return p.Token, nil
}

var _ JWTProvider = StaticJWTProvider{}

type Option func(*clientConfig)

type clientConfig struct {
	httpClient        *http.Client
	graphQLURL        string
	investorsBaseURL  string
	userAgent         string
	responseBodyLimit int64
	jwtProvider       JWTProvider
	headers           http.Header
}

func defaultClientConfig() clientConfig {
	return clientConfig{
		httpClient:        nil,
		graphQLURL:        DefaultGraphQLURL,
		investorsBaseURL:  DefaultInvestorsBaseURL,
		userAgent:         DefaultUserAgent,
		responseBodyLimit: DefaultResponseBodyLimit,
		headers:           http.Header{},
	}
}

func WithHTTPClient(client *http.Client) Option {
	return func(cfg *clientConfig) {
		cfg.httpClient = client
	}
}

func WithGraphQLURL(rawURL string) Option {
	return func(cfg *clientConfig) {
		cfg.graphQLURL = rawURL
	}
}

func WithInvestorsBaseURL(rawURL string) Option {
	return func(cfg *clientConfig) {
		cfg.investorsBaseURL = rawURL
	}
}

func WithUserAgent(userAgent string) Option {
	return func(cfg *clientConfig) {
		cfg.userAgent = userAgent
	}
}

func WithResponseBodyLimit(limit int64) Option {
	return func(cfg *clientConfig) {
		cfg.responseBodyLimit = limit
	}
}

func WithJWT(token string) Option {
	return func(cfg *clientConfig) {
		cfg.jwtProvider = StaticJWTProvider{Token: token}
	}
}

func WithJWTProvider(provider JWTProvider) Option {
	return func(cfg *clientConfig) {
		cfg.jwtProvider = provider
	}
}

func WithHeader(name, value string) Option {
	return func(cfg *clientConfig) {
		if cfg.headers == nil {
			cfg.headers = http.Header{}
		}
		cfg.headers.Set(name, value)
	}
}
