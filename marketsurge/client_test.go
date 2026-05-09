package marketsurge

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

const (
	testJWT       = "jwt-token"
	testUserAgent = "test-agent"
	testQuery     = "query TestOperation { ok }"
)

func TestNewClient(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		opts    []Option
		wantErr bool
	}{
		{name: "zero options succeeds"},
		{name: "invalid GraphQL URL", opts: []Option{WithGraphQLURL("://bad")}, wantErr: true},
		{name: "invalid investors base URL", opts: []Option{WithInvestorsBaseURL("://bad")}, wantErr: true},
		{name: "negative body limit", opts: []Option{WithResponseBodyLimit(-1)}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client, err := NewClient(tt.opts...)
			if gotErr := err != nil; gotErr != tt.wantErr {
				t.Fatalf("NewClient(%s) error = %v, want error presence = %t", tt.name, err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if client == nil {
				t.Fatal("NewClient(zero options) = nil, want Client")
			}
			if client.httpClient == nil {
				t.Fatal("NewClient(zero options).httpClient = nil, want default HTTP client")
			}
			if got, want := client.httpClient.Timeout, 30*time.Second; got != want {
				t.Fatalf("NewClient(zero options).httpClient.Timeout = %v, want %v", got, want)
			}
		})
	}
}

func TestDoGraphQL(t *testing.T) {
	t.Parallel()

	requestSeen := make(chan struct{}, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestSeen <- struct{}{}

		if got, want := r.Method, http.MethodPost; got != want {
			t.Errorf("doGraphQL() method = %s, want %s", got, want)
		}
		if got, want := r.URL.Path, "/graphql"; got != want {
			t.Errorf("doGraphQL() path = %s, want %s", got, want)
		}
		assertHeader(t, r.Header, "Authorization", "Bearer "+testJWT)
		assertHeader(t, r.Header, "Content-Type", "application/json")
		assertHeader(t, r.Header, "Apollographql-Client-Name", ApolloClientName)
		assertHeader(t, r.Header, "Dylan-Entitlement-Token", DylanToken)
		assertHeader(t, r.Header, "Referer", RefererURL)
		assertHeader(t, r.Header, "Origin", OriginURL)
		assertHeader(t, r.Header, "User-Agent", testUserAgent)
		assertHeader(t, r.Header, "X-Test", "extra")

		var got GraphQLRequest[map[string]string]
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Errorf("json.Decode(GraphQLRequest) error = %v", err)
		}
		if got.OperationName != "TestOperation" {
			t.Errorf("GraphQLRequest.OperationName = %q, want %q", got.OperationName, "TestOperation")
		}
		if got.Query != testQuery {
			t.Errorf("GraphQLRequest.Query = %q, want %q", got.Query, testQuery)
		}
		if got.Variables["symbol"] != "AAPL" {
			t.Errorf("GraphQLRequest.Variables[symbol] = %q, want %q", got.Variables["symbol"], "AAPL")
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"ok":true}}`))
	}))
	t.Cleanup(server.Close)

	client := newGraphQLTestClient(t, server.URL+"/graphql", WithHeader("X-Test", "extra"))
	var target struct {
		OK bool `json:"ok"`
	}

	err := client.doGraphQL(
		context.Background(),
		"TestOperation",
		map[string]string{"symbol": "AAPL"},
		testQuery,
		&target,
	)
	if err != nil {
		t.Fatalf("Client.doGraphQL(ctx, operation, variables, query, target) error = %v", err)
	}
	if !target.OK {
		t.Fatal("Client.doGraphQL(ctx, operation, variables, query, target) target.OK = false, want true")
	}
	select {
	case <-requestSeen:
	default:
		t.Fatal("Client.doGraphQL(ctx, operation, variables, query, target) did not send request")
	}
}

func TestDoGraphQLErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		statusCode int
		body       string
		limit      int64
		wantStatus int
		wantErr    any
	}{
		{
			name:       "status error",
			statusCode: http.StatusServiceUnavailable,
			body:       "service unavailable",
			wantStatus: http.StatusServiceUnavailable,
			wantErr:    &StatusError{},
		},
		{name: "decode error", statusCode: http.StatusOK, body: `{`, wantErr: &DecodeError{}},
		{
			name:       "graphql error",
			statusCode: http.StatusOK,
			body:       `{"errors":[{"message":"bad field","path":["quote"],"extensions":{"code":"BAD_FIELD"}}]}`,
			wantErr:    &GraphQLError{},
		},
		{
			name:       "body limit error",
			statusCode: http.StatusOK,
			body:       `{"data":{"ok":true}}`,
			limit:      5,
			wantErr:    &BodyLimitError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte(tt.body))
			}))
			t.Cleanup(server.Close)

			opts := []Option{}
			if tt.limit > 0 {
				opts = append(opts, WithResponseBodyLimit(tt.limit))
			}
			client := newGraphQLTestClient(t, server.URL, opts...)
			var target struct {
				OK bool `json:"ok"`
			}

			err := client.doGraphQL(context.Background(), "TestOperation", struct{}{}, testQuery, &target)
			assertErrorType(t, err, tt.wantErr)

			if tt.wantStatus == 0 {
				return
			}
			var statusErr *StatusError
			if !errors.As(err, &statusErr) {
				t.Fatalf("Client.doGraphQL(status %d) error = %T, want *StatusError", tt.statusCode, err)
			}
			if got := statusErr.StatusCode; got != tt.wantStatus {
				t.Fatalf("StatusError.StatusCode = %d, want %d", got, tt.wantStatus)
			}
			if got, want := string(statusErr.Body), tt.body; got != want {
				t.Fatalf("StatusError.Body = %q, want %q", got, want)
			}
		})
	}
}

func TestDoGraphQLNoAuth(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		t.Error("Client.doGraphQL(ctx, operation, variables, query, target) unexpectedly sent request")
	}))
	t.Cleanup(server.Close)

	client, err := NewClient(WithGraphQLURL(server.URL))
	if err != nil {
		t.Fatalf("NewClient(WithGraphQLURL(server.URL)) error = %v", err)
	}

	var target struct{}
	err = client.doGraphQL(context.Background(), "TestOperation", struct{}{}, testQuery, &target)
	if err == nil {
		t.Fatal("Client.doGraphQL(ctx, operation, variables, query, target) error = nil, want auth error")
	}
	if got, want := err.Error(), "no JWT provider configured"; !strings.Contains(got, want) {
		t.Fatalf("Client.doGraphQL(ctx, operation, variables, query, target) error = %q, want substring %q", got, want)
	}
}

func assertHeader(t *testing.T, header http.Header, name string, want string) {
	t.Helper()

	if got := header.Get(name); got != want {
		t.Errorf("request header %s = %q, want %q", name, got, want)
	}
}

func assertErrorType(t *testing.T, err error, want any) {
	t.Helper()

	if err == nil {
		t.Fatalf("Client.doGraphQL(ctx, operation, variables, query, target) error = nil, want %T", want)
	}
	switch want.(type) {
	case *StatusError:
		var target *StatusError
		if !errors.As(err, &target) {
			t.Fatalf("Client.doGraphQL(ctx, operation, variables, query, target) error = %T, want *StatusError", err)
		}
	case *DecodeError:
		var target *DecodeError
		if !errors.As(err, &target) {
			t.Fatalf("Client.doGraphQL(ctx, operation, variables, query, target) error = %T, want *DecodeError", err)
		}
	case *GraphQLError:
		var target *GraphQLError
		if !errors.As(err, &target) {
			t.Fatalf("Client.doGraphQL(ctx, operation, variables, query, target) error = %T, want *GraphQLError", err)
		}
	case *BodyLimitError:
		var target *BodyLimitError
		if !errors.As(err, &target) {
			t.Fatalf("Client.doGraphQL(ctx, operation, variables, query, target) error = %T, want *BodyLimitError", err)
		}
	default:
		t.Fatalf("assertErrorType(err, want) unsupported want type %T", want)
	}
}

func newGraphQLTestClient(t *testing.T, graphQLURL string, opts ...Option) *Client {
	t.Helper()

	clientOpts := []Option{WithGraphQLURL(graphQLURL), WithJWT(testJWT), WithUserAgent(testUserAgent)}
	clientOpts = append(clientOpts, opts...)
	client, err := NewClient(clientOpts...)
	if err != nil {
		t.Fatalf("NewClient(test options) error = %v", err)
	}
	return client
}
