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

// chartMarketDataFixture calls ChartMarketData against a test server that
// serves the daily response fixture and returns the first ChartMarketDataItem.
func chartMarketDataFixture(t *testing.T) ChartMarketDataItem {
	t.Helper()

	respBytes, err := os.ReadFile("testdata/ChartMarketData/response.json")
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
	req := NewChartMarketDataRequest(
		[]string{"AAPL"},
		"2025-01-01T00:00:00.000Z",
		"2025-05-02T23:59:59.000Z",
		"ONE_DAY",
		"NYSE",
	)
	resp, err := client.ChartMarketData(context.Background(), req)
	if err != nil {
		t.Fatalf("ChartMarketData() error = %v", err)
	}
	if got, want := len(resp.MarketData), 1; got != want {
		t.Fatalf("len(MarketData) = %d, want %d", got, want)
	}
	return resp.MarketData[0]
}

// chartMarketDataWeeklyFixture calls ChartMarketDataWeekly against a test
// server that serves the weekly response fixture and returns the first
// ChartMarketDataItem.
func chartMarketDataWeeklyFixture(t *testing.T) ChartMarketDataItem {
	t.Helper()

	respBytes, err := os.ReadFile("testdata/ChartMarketDataWeekly/response.json")
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
	req := NewChartMarketDataWeeklyRequest(
		[]string{"AAPL"},
		"2024-05-01T00:00:00.000Z",
		"2025-05-02T23:59:59.000Z",
		"ONE_WEEK",
	)
	resp, err := client.ChartMarketDataWeekly(context.Background(), req)
	if err != nil {
		t.Fatalf("ChartMarketDataWeekly() error = %v", err)
	}
	if got, want := len(resp.MarketData), 1; got != want {
		t.Fatalf("len(MarketData) = %d, want %d", got, want)
	}
	return resp.MarketData[0]
}

func TestNewChartMarketDataRequest(t *testing.T) {
	t.Parallel()

	req := NewChartMarketDataRequest(
		[]string{"AAPL"},
		"2025-01-01T00:00:00.000Z",
		"2025-05-02T23:59:59.000Z",
		"ONE_DAY",
		"NYSE",
	)

	if got, want := req.SymbolDialectType, DefaultChartSymbolDialectType; got != want {
		t.Errorf("NewChartMarketDataRequest().SymbolDialectType = %q, want %q", got, want)
	}
	if len(req.Symbols) != 1 || req.Symbols[0] != "AAPL" {
		t.Errorf("NewChartMarketDataRequest().Symbols = %v, want [AAPL]", req.Symbols)
	}
	if got, want := req.Where.IncludeIntradayData, true; got != want {
		t.Errorf("NewChartMarketDataRequest().Where.IncludeIntradayData = %v, want %v", got, want)
	}
	if got, want := req.ExchangeName, "NYSE"; got != want {
		t.Errorf("NewChartMarketDataRequest().ExchangeName = %q, want %q", got, want)
	}
}

func TestNewChartMarketDataWeeklyRequest(t *testing.T) {
	t.Parallel()

	req := NewChartMarketDataWeeklyRequest(
		[]string{"AAPL"},
		"2024-05-01T00:00:00.000Z",
		"2025-05-02T23:59:59.000Z",
		"ONE_WEEK",
	)

	if got, want := req.SymbolDialectType, DefaultChartSymbolDialectType; got != want {
		t.Errorf("NewChartMarketDataWeeklyRequest().SymbolDialectType = %q, want %q", got, want)
	}
	if len(req.Symbols) != 1 || req.Symbols[0] != "AAPL" {
		t.Errorf("NewChartMarketDataWeeklyRequest().Symbols = %v, want [AAPL]", req.Symbols)
	}
	if got, want := req.Where.IncludeIntradayData, true; got != want {
		t.Errorf("NewChartMarketDataWeeklyRequest().Where.IncludeIntradayData = %v, want %v", got, want)
	}
}

func TestChartMarketDataTimeSeries(t *testing.T) {
	t.Parallel()
	item := chartMarketDataFixture(t)

	if item.Pricing == nil {
		t.Fatal("Pricing is nil")
	}
	if item.Pricing.TimeSeries == nil {
		t.Fatal("TimeSeries is nil")
	}
	if got, want := item.Pricing.TimeSeries.Period, "ONE_DAY"; got != want {
		t.Errorf("TimeSeries.Period = %q, want %q", got, want)
	}
	if got, want := len(item.Pricing.TimeSeries.DataPoints), 2; got != want {
		t.Fatalf("len(DataPoints) = %d, want %d", got, want)
	}

	dp := item.Pricing.TimeSeries.DataPoints[0]
	if got, want := dp.StartDateTime, "2025-05-01T13:30:00.000Z"; got != want {
		t.Errorf("DataPoints[0].StartDateTime = %q, want %q", got, want)
	}
	if dp.Last == nil || dp.Last.Value == nil {
		t.Fatal("DataPoints[0].Last is nil")
	}
	if got, want := *dp.Last.Value, 210.45; got != want {
		t.Errorf("DataPoints[0].Last.Value = %v, want %v", got, want)
	}
}

func TestChartMarketDataQuote(t *testing.T) {
	t.Parallel()
	item := chartMarketDataFixture(t)

	if item.Pricing == nil || item.Pricing.Quote == nil {
		t.Fatal("Quote is nil")
	}
	q := item.Pricing.Quote
	if q.TradeDateTime == nil || *q.TradeDateTime != "2025-05-02T20:00:00.000Z" {
		t.Errorf("Quote.TradeDateTime = %v, want 2025-05-02T20:00:00.000Z", q.TradeDateTime)
	}
	if q.QuoteType == nil || *q.QuoteType != "REGULAR" {
		t.Errorf("Quote.QuoteType = %v, want REGULAR", q.QuoteType)
	}
	if q.Last == nil || q.Last.Value == nil || *q.Last.Value != 212.30 {
		t.Errorf("Quote.Last.Value = %v, want 212.30", q.Last)
	}
	if q.Last == nil || q.Last.FormattedValue == nil || *q.Last.FormattedValue != "212.30" {
		t.Errorf("Quote.Last.FormattedValue = %v, want 212.30", q.Last)
	}
}

func TestChartMarketDataPremarketQuote(t *testing.T) {
	t.Parallel()
	item := chartMarketDataFixture(t)

	if item.Pricing == nil || item.Pricing.PremarketQuote == nil {
		t.Fatal("PremarketQuote is nil")
	}
	q := item.Pricing.PremarketQuote
	if q.QuoteType == nil || *q.QuoteType != "PRE_MARKET" {
		t.Errorf("PremarketQuote.QuoteType = %v, want PRE_MARKET", q.QuoteType)
	}
	if q.Last == nil || q.Last.Value == nil || *q.Last.Value != 212.50 {
		t.Errorf("PremarketQuote.Last.Value = %v, want 212.50", q.Last)
	}
	if q.Last == nil || q.Last.FormattedValue == nil || *q.Last.FormattedValue != "212.50" {
		t.Errorf("PremarketQuote.Last.FormattedValue = %v, want 212.50", q.Last)
	}
}

func TestChartMarketDataPostmarketQuote(t *testing.T) {
	t.Parallel()
	item := chartMarketDataFixture(t)

	if item.Pricing == nil || item.Pricing.PostmarketQuote == nil {
		t.Fatal("PostmarketQuote is nil")
	}
	q := item.Pricing.PostmarketQuote
	if q.QuoteType == nil || *q.QuoteType != "POST_MARKET" {
		t.Errorf("PostmarketQuote.QuoteType = %v, want POST_MARKET", q.QuoteType)
	}
	if q.Last == nil || q.Last.Value == nil || *q.Last.Value != 212.05 {
		t.Errorf("PostmarketQuote.Last.Value = %v, want 212.05", q.Last)
	}
	if q.Last == nil || q.Last.FormattedValue == nil || *q.Last.FormattedValue != "212.05" {
		t.Errorf("PostmarketQuote.Last.FormattedValue = %v, want 212.05", q.Last)
	}
}

func TestChartMarketDataExchangeData(t *testing.T) {
	t.Parallel()

	respBytes, err := os.ReadFile("testdata/ChartMarketData/response.json")
	if err != nil {
		t.Fatalf("os.ReadFile(response.json) error = %v", err)
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(respBytes)
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL)
	req := NewChartMarketDataRequest(
		[]string{"AAPL"},
		"2025-01-01T00:00:00.000Z",
		"2025-05-02T23:59:59.000Z",
		"ONE_DAY",
		"NYSE",
	)
	resp, err := client.ChartMarketData(context.Background(), req)
	if err != nil {
		t.Fatalf("ChartMarketData() error = %v", err)
	}

	if resp.ExchangeData == nil {
		t.Fatal("ExchangeData is nil")
	}
	ed := resp.ExchangeData
	if ed.City == nil || *ed.City != "New York" {
		t.Errorf("ExchangeData.City = %v, want %q", ed.City, "New York")
	}
	if ed.CountryCode == nil || *ed.CountryCode != "US" {
		t.Errorf("ExchangeData.CountryCode = %v, want %q", ed.CountryCode, "US")
	}
	if ed.ExchangeISO == nil || *ed.ExchangeISO != "XNYS" {
		t.Errorf("ExchangeData.ExchangeISO = %v, want %q", ed.ExchangeISO, "XNYS")
	}
	if ed.ID == nil || *ed.ID != "NYSE" {
		t.Errorf("ExchangeData.ID = %v, want %q", ed.ID, "NYSE")
	}
	if got, want := len(ed.Holidays), 1; got != want {
		t.Fatalf("len(Holidays) = %d, want %d", got, want)
	}
	if got, want := ed.Holidays[0].Name, "Independence Day"; got != want {
		t.Errorf("Holidays[0].Name = %q, want %q", got, want)
	}
}

func TestChartMarketDataRequestBody(t *testing.T) {
	t.Parallel()

	var gotReq struct {
		OperationName string `json:"operationName"`
		Variables     struct {
			Symbols           []string `json:"symbols"`
			SymbolDialectType string   `json:"symbolDialectType"`
			Where             struct {
				StartDateTime struct {
					Eq string `json:"eq"`
				} `json:"startDateTime"`
				EndDateTime struct {
					Eq string `json:"eq"`
				} `json:"endDateTime"`
				TimeSeriesType struct {
					Eq string `json:"eq"`
				} `json:"timeSeriesType"`
				IncludeIntradayData bool `json:"includeIntradayData"`
			} `json:"where"`
			ExchangeName string `json:"exchangeName"`
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
	req := NewChartMarketDataRequest(
		[]string{"AAPL"},
		"2025-01-01T00:00:00.000Z",
		"2025-05-02T23:59:59.000Z",
		"ONE_DAY",
		"NYSE",
	)
	_, err := client.ChartMarketData(context.Background(), req)
	if err != nil {
		t.Fatalf("ChartMarketData() error = %v", err)
	}

	if got, want := gotReq.OperationName, "ChartMarketData"; got != want {
		t.Errorf("operationName = %q, want %q", got, want)
	}
	if len(gotReq.Variables.Symbols) != 1 || gotReq.Variables.Symbols[0] != "AAPL" {
		t.Errorf("variables.symbols = %v, want [AAPL]", gotReq.Variables.Symbols)
	}
	if got, want := gotReq.Variables.SymbolDialectType, DefaultChartSymbolDialectType; got != want {
		t.Errorf("variables.symbolDialectType = %q, want %q", got, want)
	}
	if got, want := gotReq.Variables.ExchangeName, "NYSE"; got != want {
		t.Errorf("variables.exchangeName = %q, want %q", got, want)
	}
	if got, want := gotReq.Variables.Where.StartDateTime.Eq, "2025-01-01T00:00:00.000Z"; got != want {
		t.Errorf("variables.where.startDateTime.eq = %q, want %q", got, want)
	}
	if got, want := gotReq.Variables.Where.TimeSeriesType.Eq, "ONE_DAY"; got != want {
		t.Errorf("variables.where.timeSeriesType.eq = %q, want %q", got, want)
	}
	if got, want := gotReq.Variables.Where.IncludeIntradayData, true; got != want {
		t.Errorf("variables.where.includeIntradayData = %v, want %v", got, want)
	}
}

func TestChartMarketDataStatusError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte("service unavailable"))
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL)
	req := NewChartMarketDataRequest([]string{"AAPL"}, "", "", "", "")
	_, err := client.ChartMarketData(context.Background(), req)

	var statusErr *StatusError
	if !errors.As(err, &statusErr) {
		t.Fatalf("ChartMarketData() error = %T, want *StatusError", err)
	}
	if got, want := statusErr.StatusCode, http.StatusServiceUnavailable; got != want {
		t.Errorf("StatusCode = %d, want %d", got, want)
	}
}

func TestChartMarketDataDecodeError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("{invalid json"))
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL)
	req := NewChartMarketDataRequest([]string{"AAPL"}, "", "", "", "")
	_, err := client.ChartMarketData(context.Background(), req)

	var decodeErr *DecodeError
	if !errors.As(err, &decodeErr) {
		t.Fatalf("ChartMarketData() error = %T, want *DecodeError", err)
	}
}

func TestChartMarketDataGraphQLError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"errors":[{"message":"not authorized","extensions":{"code":"UNAUTHORIZED"}}]}`))
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL)
	req := NewChartMarketDataRequest([]string{"AAPL"}, "", "", "", "")
	_, err := client.ChartMarketData(context.Background(), req)

	var gqlErr *GraphQLError
	if !errors.As(err, &gqlErr) {
		t.Fatalf("ChartMarketData() error = %T, want *GraphQLError", err)
	}
	if got, want := gqlErr.Error(), "not authorized"; got != want {
		t.Errorf("GraphQLError.Error() = %q, want %q", got, want)
	}
}

func TestChartMarketDataBodyLimitError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"marketData":[]}}`))
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL, WithResponseBodyLimit(5))
	req := NewChartMarketDataRequest([]string{"AAPL"}, "", "", "", "")
	_, err := client.ChartMarketData(context.Background(), req)

	if !IsBodyLimit(err) {
		t.Fatalf("ChartMarketData() error = %v, want BodyLimitError", err)
	}
}

func TestChartMarketDataNoAuth(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		t.Error("ChartMarketData() unexpectedly sent request without auth")
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(WithGraphQLURL(srv.URL))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	req := NewChartMarketDataRequest([]string{"AAPL"}, "", "", "", "")
	_, err = client.ChartMarketData(context.Background(), req)
	if err == nil {
		t.Fatal("ChartMarketData() error = nil, want auth error")
	}
	if got, want := err.Error(), "no JWT provider configured"; !strings.Contains(got, want) {
		t.Errorf("error = %q, want substring %q", got, want)
	}
}

func TestChartMarketDataFixtureSanitized(t *testing.T) {
	t.Parallel()

	for _, name := range []string{
		"testdata/ChartMarketData/request.json",
		"testdata/ChartMarketData/response.json",
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

func TestChartMarketDataWeeklyTimeSeries(t *testing.T) {
	t.Parallel()
	item := chartMarketDataWeeklyFixture(t)

	if item.Pricing == nil {
		t.Fatal("Pricing is nil")
	}
	if item.Pricing.TimeSeries == nil {
		t.Fatal("TimeSeries is nil")
	}
	if got, want := item.Pricing.TimeSeries.Period, "ONE_WEEK"; got != want {
		t.Errorf("TimeSeries.Period = %q, want %q", got, want)
	}
	if got, want := len(item.Pricing.TimeSeries.DataPoints), 1; got != want {
		t.Fatalf("len(DataPoints) = %d, want %d", got, want)
	}

	dp := item.Pricing.TimeSeries.DataPoints[0]
	if got, want := dp.StartDateTime, "2025-04-28T13:30:00.000Z"; got != want {
		t.Errorf("DataPoints[0].StartDateTime = %q, want %q", got, want)
	}
	if dp.Last == nil || dp.Last.Value == nil {
		t.Fatal("DataPoints[0].Last is nil")
	}
	if got, want := *dp.Last.Value, 212.30; got != want {
		t.Errorf("DataPoints[0].Last.Value = %v, want %v", got, want)
	}
}

func TestChartMarketDataWeeklyRequestBody(t *testing.T) {
	t.Parallel()

	var gotReq struct {
		OperationName string `json:"operationName"`
		Variables     struct {
			Symbols           []string `json:"symbols"`
			SymbolDialectType string   `json:"symbolDialectType"`
			ExchangeName      *string  `json:"exchangeName"`
			Where             struct {
				StartDateTime struct {
					Eq string `json:"eq"`
				} `json:"startDateTime"`
				EndDateTime struct {
					Eq string `json:"eq"`
				} `json:"endDateTime"`
				TimeSeriesType struct {
					Eq string `json:"eq"`
				} `json:"timeSeriesType"`
				IncludeIntradayData bool `json:"includeIntradayData"`
			} `json:"where"`
		} `json:"variables"`
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
	req := NewChartMarketDataWeeklyRequest(
		[]string{"AAPL"},
		"2024-05-01T00:00:00.000Z",
		"2025-05-02T23:59:59.000Z",
		"ONE_WEEK",
	)
	_, err := client.ChartMarketDataWeekly(context.Background(), req)
	if err != nil {
		t.Fatalf("ChartMarketDataWeekly() error = %v", err)
	}

	if got, want := gotReq.OperationName, "ChartMarketData"; got != want {
		t.Errorf("operationName = %q, want %q", got, want)
	}
	if gotReq.Variables.ExchangeName != nil {
		t.Errorf("variables.exchangeName = %v, want absent", *gotReq.Variables.ExchangeName)
	}
	if got, want := gotReq.Variables.Where.TimeSeriesType.Eq, "ONE_WEEK"; got != want {
		t.Errorf("variables.where.timeSeriesType.eq = %q, want %q", got, want)
	}
}

func TestChartMarketDataWeeklyNoExchangeData(t *testing.T) {
	t.Parallel()

	respBytes, err := os.ReadFile("testdata/ChartMarketDataWeekly/response.json")
	if err != nil {
		t.Fatalf("os.ReadFile(response.json) error = %v", err)
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(respBytes)
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL)
	req := NewChartMarketDataWeeklyRequest(
		[]string{"AAPL"},
		"2024-05-01T00:00:00.000Z",
		"2025-05-02T23:59:59.000Z",
		"ONE_WEEK",
	)
	resp, err := client.ChartMarketDataWeekly(context.Background(), req)
	if err != nil {
		t.Fatalf("ChartMarketDataWeekly() error = %v", err)
	}
	if resp.ExchangeData != nil {
		t.Errorf("ExchangeData = %v, want nil", resp.ExchangeData)
	}
}

func TestChartMarketDataWeeklyStatusError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte("service unavailable"))
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL)
	req := NewChartMarketDataWeeklyRequest([]string{"AAPL"}, "", "", "")
	_, err := client.ChartMarketDataWeekly(context.Background(), req)

	var statusErr *StatusError
	if !errors.As(err, &statusErr) {
		t.Fatalf("ChartMarketDataWeekly() error = %T, want *StatusError", err)
	}
	if got, want := statusErr.StatusCode, http.StatusServiceUnavailable; got != want {
		t.Errorf("StatusCode = %d, want %d", got, want)
	}
}

func TestChartMarketDataWeeklyNoAuth(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		t.Error("ChartMarketDataWeekly() unexpectedly sent request without auth")
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(WithGraphQLURL(srv.URL))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	req := NewChartMarketDataWeeklyRequest([]string{"AAPL"}, "", "", "")
	_, err = client.ChartMarketDataWeekly(context.Background(), req)
	if err == nil {
		t.Fatal("ChartMarketDataWeekly() error = nil, want auth error")
	}
	if got, want := err.Error(), "no JWT provider configured"; !strings.Contains(got, want) {
		t.Errorf("error = %q, want substring %q", got, want)
	}
}

func TestChartMarketDataWeeklyFixtureSanitized(t *testing.T) {
	t.Parallel()

	for _, name := range []string{
		"testdata/ChartMarketDataWeekly/request.json",
		"testdata/ChartMarketDataWeekly/response.json",
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
