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

// rsRatingRIPanelFixture calls RSRatingRIPanel against a test server that
// serves the response fixture and returns the first RSRatingRIPanelItem.
func rsRatingRIPanelFixture(t *testing.T) RSRatingRIPanelItem {
	t.Helper()

	respBytes, err := os.ReadFile("testdata/RSRatingRIPanel/response.json")
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
	req := NewRSRatingRIPanelRequest("AAPL")
	resp, err := client.RSRatingRIPanel(context.Background(), req)
	if err != nil {
		t.Fatalf("RSRatingRIPanel() error = %v", err)
	}
	if got, want := len(resp.MarketData), 1; got != want {
		t.Fatalf("len(MarketData) = %d, want %d", got, want)
	}
	return resp.MarketData[0]
}

func TestNewRSRatingRIPanelRequest(t *testing.T) {
	t.Parallel()

	req := NewRSRatingRIPanelRequest("AAPL")

	if got, want := req.SymbolDialectType, DefaultRSRatingRIPanelSymbolDialectType; got != want {
		t.Errorf("NewRSRatingRIPanelRequest().SymbolDialectType = %q, want %q", got, want)
	}
	if len(req.Symbols) != 1 || req.Symbols[0] != "AAPL" {
		t.Errorf("NewRSRatingRIPanelRequest().Symbols = %v, want [AAPL]", req.Symbols)
	}
}

func TestRSRatingRIPanelID(t *testing.T) {
	t.Parallel()
	item := rsRatingRIPanelFixture(t)

	if item.ID == nil || *item.ID != "AAPL" {
		t.Errorf("RSRatingRIPanel().MarketData[0].ID = %v, want AAPL", item.ID)
	}
}

func TestRSRatingRIPanelOriginRequest(t *testing.T) {
	t.Parallel()
	item := rsRatingRIPanelFixture(t)

	if item.OriginRequest == nil {
		t.Fatal("OriginRequest is nil")
	}
	if item.OriginRequest.FromDialect == nil || *item.OriginRequest.FromDialect != "CHARTING" {
		t.Errorf("OriginRequest.FromDialect = %v, want CHARTING", item.OriginRequest.FromDialect)
	}
	if item.OriginRequest.Symbol == nil || *item.OriginRequest.Symbol != "AAPL" {
		t.Errorf("OriginRequest.Symbol = %v, want AAPL", item.OriginRequest.Symbol)
	}
}

func TestRSRatingRIPanelRatings(t *testing.T) {
	t.Parallel()
	item := rsRatingRIPanelFixture(t)

	if item.Ratings == nil || len(item.Ratings.RSRating) < 1 {
		t.Fatal("Ratings.RSRating missing")
	}
	rating := item.Ratings.RSRating[0]
	if rating.LetterValue == nil || *rating.LetterValue != "A" {
		t.Errorf("RSRating[0].LetterValue = %v, want A", rating.LetterValue)
	}
	if rating.Period == nil || *rating.Period != "DAILY" {
		t.Errorf("RSRating[0].Period = %v, want DAILY", rating.Period)
	}
	if rating.PeriodOffset == nil || *rating.PeriodOffset != "CURRENT" {
		t.Errorf("RSRating[0].PeriodOffset = %v, want CURRENT", rating.PeriodOffset)
	}
	if rating.Value == nil || *rating.Value != 92 {
		t.Errorf("RSRating[0].Value = %v, want 92", rating.Value)
	}
}

func TestRSRatingRIPanelIntradayStatistics(t *testing.T) {
	t.Parallel()
	item := rsRatingRIPanelFixture(t)

	if item.PricingStatistics == nil || item.PricingStatistics.IntradayStatistics == nil {
		t.Fatal("PricingStatistics.IntradayStatistics is nil")
	}
	if item.PricingStatistics.IntradayStatistics.RSLineNewHigh == nil {
		t.Fatal("RSLineNewHigh is nil")
	}
	if got, want := *item.PricingStatistics.IntradayStatistics.RSLineNewHigh, true; got != want {
		t.Errorf("RSLineNewHigh = %v, want %v", got, want)
	}
}

func TestRSRatingRIPanelRequestBody(t *testing.T) {
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
	req := NewRSRatingRIPanelRequest("AAPL")
	_, err := client.RSRatingRIPanel(context.Background(), req)
	if err != nil {
		t.Fatalf("RSRatingRIPanel() error = %v", err)
	}

	if got, want := gotReq.OperationName, "RSRatingRIPanel"; got != want {
		t.Errorf("operationName = %q, want %q", got, want)
	}
	if len(gotReq.Variables.Symbols) != 1 || gotReq.Variables.Symbols[0] != "AAPL" {
		t.Errorf("variables.symbols = %v, want [AAPL]", gotReq.Variables.Symbols)
	}
	if got, want := gotReq.Variables.SymbolDialectType, DefaultRSRatingRIPanelSymbolDialectType; got != want {
		t.Errorf("variables.symbolDialectType = %q, want %q", got, want)
	}
}

func TestRSRatingRIPanelStatusError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte("service unavailable"))
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL)
	req := NewRSRatingRIPanelRequest("AAPL")
	_, err := client.RSRatingRIPanel(context.Background(), req)

	var statusErr *StatusError
	if !errors.As(err, &statusErr) {
		t.Fatalf("RSRatingRIPanel() error = %T, want *StatusError", err)
	}
	if got, want := statusErr.StatusCode, http.StatusServiceUnavailable; got != want {
		t.Errorf("StatusCode = %d, want %d", got, want)
	}
}

func TestRSRatingRIPanelDecodeError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("{invalid json"))
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL)
	req := NewRSRatingRIPanelRequest("AAPL")
	_, err := client.RSRatingRIPanel(context.Background(), req)

	var decodeErr *DecodeError
	if !errors.As(err, &decodeErr) {
		t.Fatalf("RSRatingRIPanel() error = %T, want *DecodeError", err)
	}
}

func TestRSRatingRIPanelGraphQLError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"errors":[{"message":"not authorized","extensions":{"code":"UNAUTHORIZED"}}]}`))
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL)
	req := NewRSRatingRIPanelRequest("AAPL")
	_, err := client.RSRatingRIPanel(context.Background(), req)

	var gqlErr *GraphQLError
	if !errors.As(err, &gqlErr) {
		t.Fatalf("RSRatingRIPanel() error = %T, want *GraphQLError", err)
	}
	if got, want := gqlErr.Error(), "not authorized"; got != want {
		t.Errorf("GraphQLError.Error() = %q, want %q", got, want)
	}
}

func TestRSRatingRIPanelBodyLimitError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"marketData":[]}}`))
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL, WithResponseBodyLimit(5))
	req := NewRSRatingRIPanelRequest("AAPL")
	_, err := client.RSRatingRIPanel(context.Background(), req)

	if !IsBodyLimit(err) {
		t.Fatalf("RSRatingRIPanel() error = %v, want BodyLimitError", err)
	}
}

func TestRSRatingRIPanelNoAuth(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		t.Error("RSRatingRIPanel() unexpectedly sent request without auth")
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(WithGraphQLURL(srv.URL))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	req := NewRSRatingRIPanelRequest("AAPL")
	_, err = client.RSRatingRIPanel(context.Background(), req)
	if err == nil {
		t.Fatal("RSRatingRIPanel() error = nil, want auth error")
	}
	if got, want := err.Error(), "no JWT provider configured"; !strings.Contains(got, want) {
		t.Errorf("error = %q, want substring %q", got, want)
	}
}

func TestRSRatingRIPanelFixtureSanitized(t *testing.T) {
	t.Parallel()

	for _, name := range []string{
		"testdata/RSRatingRIPanel/request.json",
		"testdata/RSRatingRIPanel/response.json",
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
