package marketsurge

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// ownershipFixture calls Ownership against a test server that serves
// the response fixture and returns the first OwnershipItem.
func ownershipFixture(t *testing.T) OwnershipItem {
	t.Helper()

	respBytes, err := os.ReadFile("testdata/Ownership/response.json")
	if err != nil {
		t.Fatalf("os.ReadFile(response.json) error = %v", err)
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.Method, http.MethodPost; got != want {
			t.Errorf("request method = %s, want %s", got, want)
		}
		assertHeader(t, r.Header, "Authorization", "Bearer "+testJWT)
		assertHeader(t, r.Header, "Content-Type", "application/json")
		assertHeader(t, r.Header, "Apollographql-Client-Name", ApolloClientName)
		assertHeader(t, r.Header, "Dylan-Entitlement-Token", DylanToken)

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(respBytes)
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL)
	req := NewOwnershipRequest("AAPL")
	resp, err := client.Ownership(context.Background(), req)
	if err != nil {
		t.Fatalf("Ownership() error = %v", err)
	}
	if got, want := len(resp.MarketData), 1; got != want {
		t.Fatalf("len(MarketData) = %d, want %d", got, want)
	}
	return resp.MarketData[0]
}

func TestNewOwnershipRequest(t *testing.T) {
	t.Parallel()

	req := NewOwnershipRequest("AAPL")

	if got, want := req.SymbolDialectType, DefaultOwnershipSymbolDialectType; got != want {
		t.Errorf("NewOwnershipRequest().SymbolDialectType = %q, want %q", got, want)
	}
	if len(req.Symbols) != 1 || req.Symbols[0] != "AAPL" {
		t.Errorf("NewOwnershipRequest().Symbols = %v, want [AAPL]", req.Symbols)
	}
}

func TestOwnershipFundsFloatPercentHeld(t *testing.T) {
	t.Parallel()
	item := ownershipFixture(t)

	if item.Ownership == nil {
		t.Fatal("Ownership is nil")
	}
	if item.Ownership.FundsFloatPercentHeld == nil {
		t.Fatal("FundsFloatPercentHeld is nil")
	}
	if item.Ownership.FundsFloatPercentHeld.FormattedValue == nil {
		t.Fatal("FundsFloatPercentHeld.FormattedValue is nil")
	}
	if got, want := *item.Ownership.FundsFloatPercentHeld.FormattedValue, "62.3%"; got != want {
		t.Errorf("FundsFloatPercentHeld.FormattedValue = %q, want %q", got, want)
	}
}

func TestOwnershipFundOwnershipSummary(t *testing.T) {
	t.Parallel()
	item := ownershipFixture(t)

	if item.Ownership == nil {
		t.Fatal("Ownership is nil")
	}
	if got := len(item.Ownership.FundOwnershipSummary); got < 1 {
		t.Fatal("FundOwnershipSummary is empty")
	}
	summary := item.Ownership.FundOwnershipSummary[0]
	if summary.Date == nil || summary.Date.Value != "2026-03-31" {
		t.Errorf("FundOwnershipSummary[0].Date = %v, want 2026-03-31", summary.Date)
	}
	if summary.NumberOfFundsHeld == nil || summary.NumberOfFundsHeld.FormattedValue == nil {
		t.Fatal("FundOwnershipSummary[0].NumberOfFundsHeld is nil")
	}
	if got, want := *summary.NumberOfFundsHeld.FormattedValue, "5,432"; got != want {
		t.Errorf("NumberOfFundsHeld.FormattedValue = %q, want %q", got, want)
	}
}

func TestOwnershipRequestBody(t *testing.T) {
	t.Parallel()

	var gotReq struct {
		OperationName string `json:"operationName"`
		Variables     struct {
			Symbols           []string `json:"symbols"`
			SymbolDialectType string   `json:"symbolDialectType"`
		} `json:"variables"`
		Query string `json:"query"`
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&gotReq); err != nil {
			t.Errorf("decode request body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"marketData":[]}}`))
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL)
	req := NewOwnershipRequest("AAPL")
	_, err := client.Ownership(context.Background(), req)
	if err != nil {
		t.Fatalf("Ownership() error = %v", err)
	}

	if got, want := gotReq.OperationName, "Ownership"; got != want {
		t.Errorf("operationName = %q, want %q", got, want)
	}
	if len(gotReq.Variables.Symbols) != 1 || gotReq.Variables.Symbols[0] != "AAPL" {
		t.Errorf("variables.symbols = %v, want [AAPL]", gotReq.Variables.Symbols)
	}
	if got, want := gotReq.Variables.SymbolDialectType, DefaultOwnershipSymbolDialectType; got != want {
		t.Errorf("variables.symbolDialectType = %q, want %q", got, want)
	}
}

func TestOwnershipStatusError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte("service unavailable"))
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL)
	req := NewOwnershipRequest("AAPL")
	_, err := client.Ownership(context.Background(), req)

	var statusErr *StatusError
	if !errors.As(err, &statusErr) {
		t.Fatalf("Ownership() error = %T, want *StatusError", err)
	}
	if got, want := statusErr.StatusCode, http.StatusServiceUnavailable; got != want {
		t.Errorf("StatusCode = %d, want %d", got, want)
	}
}

func TestOwnershipDecodeError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("{invalid json"))
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL)
	req := NewOwnershipRequest("AAPL")
	_, err := client.Ownership(context.Background(), req)

	var decodeErr *DecodeError
	if !errors.As(err, &decodeErr) {
		t.Fatalf("Ownership() error = %T, want *DecodeError", err)
	}
}

func TestOwnershipGraphQLError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"errors":[{"message":"not authorized","extensions":{"code":"UNAUTHORIZED"}}]}`))
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL)
	req := NewOwnershipRequest("AAPL")
	_, err := client.Ownership(context.Background(), req)

	var gqlErr *GraphQLError
	if !errors.As(err, &gqlErr) {
		t.Fatalf("Ownership() error = %T, want *GraphQLError", err)
	}
	if got, want := gqlErr.Error(), "not authorized"; got != want {
		t.Errorf("GraphQLError.Error() = %q, want %q", got, want)
	}
}

func TestOwnershipBodyLimitError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"marketData":[]}}`))
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL, WithResponseBodyLimit(5))
	req := NewOwnershipRequest("AAPL")
	_, err := client.Ownership(context.Background(), req)

	if !IsBodyLimit(err) {
		t.Fatalf("Ownership() error = %v, want BodyLimitError", err)
	}
}

func TestOwnershipNoAuth(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		t.Error("Ownership() unexpectedly sent request without auth")
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(WithGraphQLURL(srv.URL))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	req := NewOwnershipRequest("AAPL")
	_, err = client.Ownership(context.Background(), req)
	if err == nil {
		t.Fatal("Ownership() error = nil, want auth error")
	}
	if got, want := err.Error(), "no JWT provider configured"; !strings.Contains(got, want) {
		t.Errorf("error = %q, want substring %q", got, want)
	}
}

func TestOwnershipFixtureSanitized(t *testing.T) {
	t.Parallel()

	for _, name := range []string{
		"testdata/Ownership/request.json",
		"testdata/Ownership/response.json",
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			data, err := os.ReadFile(name)
			if err != nil {
				t.Fatalf("os.ReadFile(%s) error = %v", name, err)
			}
			content := string(data)
			for _, secret := range []string{"Bearer ", "password", "cookie", "session_id"} {
				if strings.Contains(strings.ToLower(content), strings.ToLower(secret)) {
					t.Errorf("fixture %s contains potential secret %q", name, secret)
				}
			}
		})
	}
}
