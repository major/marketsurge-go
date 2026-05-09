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

// ---------------------------------------------------------------------------
// GetAllWatchlistNames
// ---------------------------------------------------------------------------

// watchlistNamesFixture calls GetAllWatchlistNames against a test server that
// serves the response fixture and returns the first WatchlistSummary.
func watchlistNamesFixture(t *testing.T) WatchlistSummary {
	t.Helper()

	respBytes, err := os.ReadFile("testdata/GetAllWatchlistNames/response.json")
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
	req := NewGetAllWatchlistNamesRequest()
	resp, err := client.GetAllWatchlistNames(context.Background(), req)
	if err != nil {
		t.Fatalf("GetAllWatchlistNames() error = %v", err)
	}
	if got, want := len(resp.Watchlists), 1; got != want {
		t.Fatalf("len(Watchlists) = %d, want %d", got, want)
	}
	return resp.Watchlists[0]
}

func TestNewGetAllWatchlistNamesRequest(t *testing.T) {
	t.Parallel()

	req := NewGetAllWatchlistNamesRequest()

	if got, want := req.Pub, DefaultWatchlistPub; got != want {
		t.Errorf("NewGetAllWatchlistNamesRequest().Pub = %q, want %q", got, want)
	}
}

func TestGetAllWatchlistNamesID(t *testing.T) {
	t.Parallel()
	item := watchlistNamesFixture(t)

	if got, want := item.ID, "12345"; got != want {
		t.Errorf("GetAllWatchlistNames().Watchlists[0].ID = %q, want %q", got, want)
	}
}

func TestGetAllWatchlistNamesMetadata(t *testing.T) {
	t.Parallel()
	item := watchlistNamesFixture(t)

	if got, want := item.Name, "My Watchlist"; got != want {
		t.Errorf("GetAllWatchlistNames().Watchlists[0].Name = %q, want %q", got, want)
	}
	if got, want := item.LastModifiedDateUtc, "2025-01-01T00:00:00Z"; got != want {
		t.Errorf("GetAllWatchlistNames().Watchlists[0].LastModifiedDateUtc = %q, want %q", got, want)
	}
	if got, want := item.Description, "Test watchlist"; got != want {
		t.Errorf("GetAllWatchlistNames().Watchlists[0].Description = %q, want %q", got, want)
	}
}

func TestGetAllWatchlistNamesRequestBody(t *testing.T) {
	t.Parallel()

	var gotReq struct {
		OperationName string `json:"operationName"`
		Variables     struct {
			Pub string `json:"pub"`
		} `json:"variables"`
		Query string `json:"query"`
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&gotReq); err != nil {
			t.Errorf("decode request body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"watchlists":[]}}`))
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL)
	req := NewGetAllWatchlistNamesRequest()
	_, err := client.GetAllWatchlistNames(context.Background(), req)
	if err != nil {
		t.Fatalf("GetAllWatchlistNames() error = %v", err)
	}

	if got, want := gotReq.OperationName, "GetAllWatchlistNames"; got != want {
		t.Errorf("operationName = %q, want %q", got, want)
	}
	if got, want := gotReq.Variables.Pub, DefaultWatchlistPub; got != want {
		t.Errorf("variables.pub = %q, want %q", got, want)
	}
}

func TestGetAllWatchlistNamesStatusError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte("service unavailable"))
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL)
	req := NewGetAllWatchlistNamesRequest()
	_, err := client.GetAllWatchlistNames(context.Background(), req)

	var statusErr *StatusError
	if !errors.As(err, &statusErr) {
		t.Fatalf("GetAllWatchlistNames() error = %T, want *StatusError", err)
	}
	if got, want := statusErr.StatusCode, http.StatusServiceUnavailable; got != want {
		t.Errorf("StatusCode = %d, want %d", got, want)
	}
}

func TestGetAllWatchlistNamesDecodeError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("{invalid json"))
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL)
	req := NewGetAllWatchlistNamesRequest()
	_, err := client.GetAllWatchlistNames(context.Background(), req)

	var decodeErr *DecodeError
	if !errors.As(err, &decodeErr) {
		t.Fatalf("GetAllWatchlistNames() error = %T, want *DecodeError", err)
	}
}

func TestGetAllWatchlistNamesGraphQLError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"errors":[{"message":"not authorized","extensions":{"code":"UNAUTHORIZED"}}]}`))
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL)
	req := NewGetAllWatchlistNamesRequest()
	_, err := client.GetAllWatchlistNames(context.Background(), req)

	var gqlErr *GraphQLError
	if !errors.As(err, &gqlErr) {
		t.Fatalf("GetAllWatchlistNames() error = %T, want *GraphQLError", err)
	}
	if got, want := gqlErr.Error(), "not authorized"; got != want {
		t.Errorf("GraphQLError.Error() = %q, want %q", got, want)
	}
}

func TestGetAllWatchlistNamesBodyLimit(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"watchlists":[]}}`))
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL, WithResponseBodyLimit(5))
	req := NewGetAllWatchlistNamesRequest()
	_, err := client.GetAllWatchlistNames(context.Background(), req)

	if !IsBodyLimit(err) {
		t.Fatalf("GetAllWatchlistNames() error = %v, want BodyLimitError", err)
	}
}

func TestGetAllWatchlistNamesNoAuth(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		t.Error("GetAllWatchlistNames() unexpectedly sent request without auth")
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(WithGraphQLURL(srv.URL))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	req := NewGetAllWatchlistNamesRequest()
	_, err = client.GetAllWatchlistNames(context.Background(), req)
	if err == nil {
		t.Fatal("GetAllWatchlistNames() error = nil, want auth error")
	}
	if got, want := err.Error(), "no JWT provider configured"; !strings.Contains(got, want) {
		t.Errorf("error = %q, want substring %q", got, want)
	}
}

func TestGetAllWatchlistNamesFixtureSanitized(t *testing.T) {
	t.Parallel()

	for _, name := range []string{
		"testdata/GetAllWatchlistNames/request.json",
		"testdata/GetAllWatchlistNames/response.json",
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

// ---------------------------------------------------------------------------
// FlaggedSymbols
// ---------------------------------------------------------------------------

// flaggedSymbolsFixture calls FlaggedSymbols against a test server that serves
// the response fixture and returns the WatchlistDetail.
func flaggedSymbolsFixture(t *testing.T) WatchlistDetail {
	t.Helper()

	respBytes, err := os.ReadFile("testdata/FlaggedSymbols/response.json")
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
	req := NewFlaggedSymbolsRequest("12345")
	resp, err := client.FlaggedSymbols(context.Background(), req)
	if err != nil {
		t.Fatalf("FlaggedSymbols() error = %v", err)
	}
	return resp.Watchlist
}

func TestNewFlaggedSymbolsRequest(t *testing.T) {
	t.Parallel()

	req := NewFlaggedSymbolsRequest("12345")

	if got, want := req.Pub, DefaultWatchlistPub; got != want {
		t.Errorf("NewFlaggedSymbolsRequest().Pub = %q, want %q", got, want)
	}
	if got, want := req.WatchlistID, "12345"; got != want {
		t.Errorf("NewFlaggedSymbolsRequest().WatchlistID = %q, want %q", got, want)
	}
}

func TestFlaggedSymbolsID(t *testing.T) {
	t.Parallel()
	wl := flaggedSymbolsFixture(t)

	if got, want := wl.ID, "12345"; got != want {
		t.Errorf("FlaggedSymbols().Watchlist.ID = %q, want %q", got, want)
	}
}

func TestFlaggedSymbolsMetadata(t *testing.T) {
	t.Parallel()
	wl := flaggedSymbolsFixture(t)

	if got, want := wl.Name, "My Watchlist"; got != want {
		t.Errorf("FlaggedSymbols().Watchlist.Name = %q, want %q", got, want)
	}
	if got, want := wl.LastModifiedDateUtc, "2025-01-01T00:00:00Z"; got != want {
		t.Errorf("FlaggedSymbols().Watchlist.LastModifiedDateUtc = %q, want %q", got, want)
	}
	if got, want := wl.Description, "Test watchlist"; got != want {
		t.Errorf("FlaggedSymbols().Watchlist.Description = %q, want %q", got, want)
	}
}

func TestFlaggedSymbolsItems(t *testing.T) {
	t.Parallel()
	wl := flaggedSymbolsFixture(t)

	if got, want := len(wl.Items), 2; got != want {
		t.Fatalf("len(FlaggedSymbols().Watchlist.Items) = %d, want %d", got, want)
	}
	if got, want := wl.Items[0].Key, "AAPL"; got != want {
		t.Errorf("Items[0].Key = %q, want %q", got, want)
	}
	if got, want := wl.Items[0].DowJonesKey, "US:AAPL"; got != want {
		t.Errorf("Items[0].DowJonesKey = %q, want %q", got, want)
	}
	if got, want := wl.Items[1].Key, "MSFT"; got != want {
		t.Errorf("Items[1].Key = %q, want %q", got, want)
	}
	if got, want := wl.Items[1].DowJonesKey, "US:MSFT"; got != want {
		t.Errorf("Items[1].DowJonesKey = %q, want %q", got, want)
	}
}

func TestFlaggedSymbolsRequestBody(t *testing.T) {
	t.Parallel()

	var gotReq struct {
		OperationName string `json:"operationName"`
		Variables     struct {
			Pub         string `json:"pub"`
			WatchlistID string `json:"watchlistId"`
		} `json:"variables"`
		Query string `json:"query"`
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&gotReq); err != nil {
			t.Errorf("decode request body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(
			[]byte(`{"data":{"watchlist":{"id":"","name":"","lastModifiedDateUtc":"","description":"","items":[]}}}`),
		)
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL)
	req := NewFlaggedSymbolsRequest("12345")
	_, err := client.FlaggedSymbols(context.Background(), req)
	if err != nil {
		t.Fatalf("FlaggedSymbols() error = %v", err)
	}

	if got, want := gotReq.OperationName, "FlaggedSymbols"; got != want {
		t.Errorf("operationName = %q, want %q", got, want)
	}
	if got, want := gotReq.Variables.Pub, DefaultWatchlistPub; got != want {
		t.Errorf("variables.pub = %q, want %q", got, want)
	}
	if got, want := gotReq.Variables.WatchlistID, "12345"; got != want {
		t.Errorf("variables.watchlistId = %q, want %q", got, want)
	}
}

func TestFlaggedSymbolsStatusError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte("service unavailable"))
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL)
	req := NewFlaggedSymbolsRequest("12345")
	_, err := client.FlaggedSymbols(context.Background(), req)

	var statusErr *StatusError
	if !errors.As(err, &statusErr) {
		t.Fatalf("FlaggedSymbols() error = %T, want *StatusError", err)
	}
	if got, want := statusErr.StatusCode, http.StatusServiceUnavailable; got != want {
		t.Errorf("StatusCode = %d, want %d", got, want)
	}
}

func TestFlaggedSymbolsDecodeError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("{invalid json"))
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL)
	req := NewFlaggedSymbolsRequest("12345")
	_, err := client.FlaggedSymbols(context.Background(), req)

	var decodeErr *DecodeError
	if !errors.As(err, &decodeErr) {
		t.Fatalf("FlaggedSymbols() error = %T, want *DecodeError", err)
	}
}

func TestFlaggedSymbolsGraphQLError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"errors":[{"message":"not authorized","extensions":{"code":"UNAUTHORIZED"}}]}`))
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL)
	req := NewFlaggedSymbolsRequest("12345")
	_, err := client.FlaggedSymbols(context.Background(), req)

	var gqlErr *GraphQLError
	if !errors.As(err, &gqlErr) {
		t.Fatalf("FlaggedSymbols() error = %T, want *GraphQLError", err)
	}
	if got, want := gqlErr.Error(), "not authorized"; got != want {
		t.Errorf("GraphQLError.Error() = %q, want %q", got, want)
	}
}

func TestFlaggedSymbolsBodyLimit(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(
			[]byte(`{"data":{"watchlist":{"id":"","name":"","lastModifiedDateUtc":"","description":"","items":[]}}}`),
		)
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL, WithResponseBodyLimit(5))
	req := NewFlaggedSymbolsRequest("12345")
	_, err := client.FlaggedSymbols(context.Background(), req)

	if !IsBodyLimit(err) {
		t.Fatalf("FlaggedSymbols() error = %v, want BodyLimitError", err)
	}
}

func TestFlaggedSymbolsNoAuth(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		t.Error("FlaggedSymbols() unexpectedly sent request without auth")
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(WithGraphQLURL(srv.URL))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	req := NewFlaggedSymbolsRequest("12345")
	_, err = client.FlaggedSymbols(context.Background(), req)
	if err == nil {
		t.Fatal("FlaggedSymbols() error = nil, want auth error")
	}
	if got, want := err.Error(), "no JWT provider configured"; !strings.Contains(got, want) {
		t.Errorf("error = %q, want substring %q", got, want)
	}
}

func TestFlaggedSymbolsFixtureSanitized(t *testing.T) {
	t.Parallel()

	for _, name := range []string{
		"testdata/FlaggedSymbols/request.json",
		"testdata/FlaggedSymbols/response.json",
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
