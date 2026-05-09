package marketsurge

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const defaultHTTPClientTimeout = 30 * time.Second

const (
	// DylanToken is the public API exchange token, not a secret.
	DylanToken = "x4ckyhshg90pdq6bwf6n1voijs7r3fdk" // #nosec G101 -- public API exchange token, not a secret.
	// ApolloClientName identifies this client to the GraphQL gateway.
	ApolloClientName = "marketsurge"
	// RefererURL is the browser referer expected by MarketSurge.
	RefererURL = "https://marketsurge.investors.com/"
	// OriginURL is the browser origin expected by MarketSurge.
	OriginURL = "https://marketsurge.investors.com"
)

const (
	headerApolloClientName    = "Apollographql-Client-Name"
	headerAuthorization       = "Authorization"
	headerContentType         = "Content-Type"
	headerDylanEntitlement    = "Dylan-Entitlement-Token"
	headerOrigin              = "Origin"
	headerReferer             = "Referer"
	headerUserAgent           = "User-Agent"
	mediaTypeApplicationJSON  = "application/json"
	authorizationBearerPrefix = "Bearer "
)

// Client sends authenticated requests to the MarketSurge GraphQL API.
type Client struct {
	httpClient        *http.Client
	graphQLURL        string
	investorsBaseURL  string
	userAgent         string
	responseBodyLimit int64
	jwtProvider       JWTProvider
	headers           http.Header
}

// NewClient creates a Client from options.
func NewClient(opts ...Option) (*Client, error) {
	cfg := defaultClientConfig()
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}

	if err := validateAbsoluteURL("GraphQL", cfg.graphQLURL); err != nil {
		return nil, err
	}
	if err := validateAbsoluteURL("investors base", cfg.investorsBaseURL); err != nil {
		return nil, err
	}
	if cfg.responseBodyLimit < 0 {
		return nil, errors.New("marketsurge: response body limit must be non-negative")
	}

	httpClient := cfg.httpClient
	if httpClient == nil {
		httpClient = &http.Client{Timeout: defaultHTTPClientTimeout}
	}

	return &Client{
		httpClient:        httpClient,
		graphQLURL:        cfg.graphQLURL,
		investorsBaseURL:  cfg.investorsBaseURL,
		userAgent:         cfg.userAgent,
		responseBodyLimit: cfg.responseBodyLimit,
		jwtProvider:       cfg.jwtProvider,
		headers:           cloneHeader(cfg.headers),
	}, nil
}

func (c *Client) doGraphQL(ctx context.Context, operationName string, variables any, query string, target any) error {
	if c.jwtProvider == nil {
		return errors.New("marketsurge: no JWT provider configured; use WithJWT or WithJWTProvider")
	}

	jwt, err := c.jwtProvider.JWT(ctx)
	if err != nil {
		return fmt.Errorf("marketsurge: get JWT: %w", err)
	}

	body, err := json.Marshal(GraphQLRequest[any]{
		OperationName: operationName,
		Variables:     variables,
		Query:         query,
	})
	if err != nil {
		return fmt.Errorf("marketsurge: marshal GraphQL request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.graphQLURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("marketsurge: create GraphQL request: %w", err)
	}
	c.setGraphQLHeaders(req, jwt)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("marketsurge: send GraphQL request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return c.statusError(resp)
	}

	responseBody, err := c.readBody(resp.Body)
	if err != nil {
		return err
	}

	return decodeGraphQLResponse(responseBody, target)
}

func validateAbsoluteURL(name string, rawURL string) error {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("marketsurge: invalid %s URL %q: %w", name, rawURL, err)
	}
	if parsed.Scheme == "" {
		return fmt.Errorf("marketsurge: invalid %s URL %q: missing scheme", name, rawURL)
	}
	return nil
}

func (c *Client) setGraphQLHeaders(req *http.Request, jwt string) {
	req.Header.Set(headerAuthorization, authorizationBearerPrefix+jwt)
	req.Header.Set(headerContentType, mediaTypeApplicationJSON)
	req.Header.Set(headerApolloClientName, ApolloClientName)
	req.Header.Set(headerDylanEntitlement, DylanToken)
	req.Header.Set(headerReferer, RefererURL)
	req.Header.Set(headerOrigin, OriginURL)
	req.Header.Set(headerUserAgent, c.userAgent)
	copyHeaderInto(req.Header, c.headers)
}

func (c *Client) statusError(resp *http.Response) error {
	body, err := readBodyUpToLimit(resp.Body, c.responseBodyLimit)
	if err != nil {
		return fmt.Errorf("marketsurge: read error response body: %w", err)
	}
	return &StatusError{
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
		Body:       body,
		Header:     cloneHeader(resp.Header),
	}
}

func (c *Client) readBody(body io.Reader) ([]byte, error) {
	data, err := readBodyUpToLimit(body, c.responseBodyLimit+1)
	if err != nil {
		return nil, fmt.Errorf("marketsurge: read GraphQL response body: %w", err)
	}
	if int64(len(data)) > c.responseBodyLimit {
		return nil, &BodyLimitError{Limit: c.responseBodyLimit}
	}
	return data, nil
}

func decodeGraphQLResponse(body []byte, target any) error {
	var response GraphQLResponse[json.RawMessage]
	if err := json.Unmarshal(body, &response); err != nil {
		return &DecodeError{Operation: "decode GraphQL response", Err: err}
	}
	if len(response.Errors) > 0 {
		return &GraphQLError{Errors: response.Errors}
	}
	if target == nil {
		return nil
	}
	if err := json.Unmarshal(response.Data, target); err != nil {
		return &DecodeError{Operation: "decode GraphQL data", Err: err}
	}
	return nil
}

func readBodyUpToLimit(body io.Reader, limit int64) ([]byte, error) {
	if limit <= 0 {
		return nil, nil
	}
	return io.ReadAll(io.LimitReader(body, limit))
}

func cloneHeader(header http.Header) http.Header {
	if header == nil {
		return http.Header{}
	}

	cloned := make(http.Header, len(header))
	copyHeaderInto(cloned, header)
	return cloned
}

func copyHeaderInto(dst http.Header, src http.Header) {
	for name, values := range src {
		dst.Del(name)
		for _, value := range values {
			dst.Add(name, value)
		}
	}
}
