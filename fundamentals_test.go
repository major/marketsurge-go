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

// fundamentalsFixture calls Fundamentals against a test server that serves
// the response fixture and returns the first FundamentalsItem.
func fundamentalsFixture(t *testing.T) FundamentalsItem {
	t.Helper()

	respBytes, err := os.ReadFile("testdata/FundermentalDataBox/response.json")
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
	req := NewFundamentalsRequest("AAPL")
	resp, err := client.Fundamentals(context.Background(), req)
	if err != nil {
		t.Fatalf("Fundamentals() error = %v", err)
	}
	if got, want := len(resp.MarketData), 1; got != want {
		t.Fatalf("len(MarketData) = %d, want %d", got, want)
	}
	return resp.MarketData[0]
}

func TestNewFundamentalsRequest(t *testing.T) {
	t.Parallel()

	req := NewFundamentalsRequest("AAPL")

	if got, want := req.SymbolDialectType, DefaultFundamentalsSymbolDialectType; got != want {
		t.Errorf("NewFundamentalsRequest().SymbolDialectType = %q, want %q", got, want)
	}
	if got, want := req.UpToHistoricalPeriodOffset, DefaultFundamentalsUpToHistoricalPeriodOffset; got != want {
		t.Errorf("NewFundamentalsRequest().UpToHistoricalPeriodOffset = %q, want %q", got, want)
	}
	if got, want := req.UpToQueryPeriodOffset, DefaultFundamentalsUpToQueryPeriodOffset; got != want {
		t.Errorf("NewFundamentalsRequest().UpToQueryPeriodOffset = %q, want %q", got, want)
	}
	if got, want := req.ReportedSalesUpToHistoricalPeriod2, DefaultFundamentalsReportedSalesUpToHistoricalPeriod2; got != want {
		t.Errorf("NewFundamentalsRequest().ReportedSalesUpToHistoricalPeriod2 = %q, want %q", got, want)
	}
	if got, want := req.SalesEstimatesUpToQueryPeriod2, DefaultFundamentalsSalesEstimatesUpToQueryPeriod2; got != want {
		t.Errorf("NewFundamentalsRequest().SalesEstimatesUpToQueryPeriod2 = %q, want %q", got, want)
	}
	if len(req.Symbols) != 1 || req.Symbols[0] != "AAPL" {
		t.Errorf("NewFundamentalsRequest().Symbols = %v, want [AAPL]", req.Symbols)
	}
}

func TestFundamentalsID(t *testing.T) {
	t.Parallel()
	item := fundamentalsFixture(t)

	if item.ID == nil || *item.ID != "AAPL" {
		t.Errorf("Fundamentals().MarketData[0].ID = %v, want AAPL", item.ID)
	}
}

func TestFundamentalsSymbology(t *testing.T) {
	t.Parallel()
	item := fundamentalsFixture(t)

	if item.Symbology == nil || item.Symbology.Company == nil {
		t.Fatal("Symbology.Company is nil")
	}
	if item.Symbology.Company.CompanyName == nil || *item.Symbology.Company.CompanyName != "Apple Inc." {
		t.Errorf("CompanyName = %v, want Apple Inc.", item.Symbology.Company.CompanyName)
	}
	if item.Symbology.Instrument == nil || len(item.Symbology.Instrument.Symbols) != 1 {
		t.Fatal("Symbology.Instrument.Symbols missing")
	}
	sym := item.Symbology.Instrument.Symbols[0]
	if sym.Value == nil || *sym.Value != "AAPL" {
		t.Errorf("Symbol.Value = %v, want AAPL", sym.Value)
	}
	if sym.Type == nil || *sym.Type != "CHARTING" {
		t.Errorf("Symbol.Type = %v, want CHARTING", sym.Type)
	}
}

func TestFundamentalsReportedEarnings(t *testing.T) {
	t.Parallel()
	item := fundamentalsFixture(t)

	if item.Financials == nil || item.Financials.ConsensusFinancials == nil {
		t.Fatal("Financials.ConsensusFinancials is nil")
	}
	eps := item.Financials.ConsensusFinancials.EPS
	if eps == nil || len(eps.ReportedEarnings) < 1 {
		t.Fatal("EPS.ReportedEarnings missing")
	}
	re := eps.ReportedEarnings[0]
	if re.Value == nil || re.Value.Value == nil {
		t.Fatal("ReportedEarnings[0].Value is nil")
	}
	if got, want := *re.Value.Value, 1.65; got != want {
		t.Errorf("ReportedEarnings[0].Value.Value = %v, want %v", got, want)
	}
	if re.PercentChangeYOY == nil || re.PercentChangeYOY.Value == nil {
		t.Fatal("ReportedEarnings[0].PercentChangeYOY is nil")
	}
	if got, want := *re.PercentChangeYOY.Value, 12.5; got != want {
		t.Errorf("ReportedEarnings[0].PercentChangeYOY.Value = %v, want %v", got, want)
	}
	if re.PeriodOffset == nil || *re.PeriodOffset != "CURRENT" {
		t.Errorf("ReportedEarnings[0].PeriodOffset = %v, want CURRENT", re.PeriodOffset)
	}
}

func TestFundamentalsReportedSales(t *testing.T) {
	t.Parallel()
	item := fundamentalsFixture(t)

	if item.Financials == nil || item.Financials.ConsensusFinancials == nil {
		t.Fatal("Financials.ConsensusFinancials is nil")
	}
	sales := item.Financials.ConsensusFinancials.Sales
	if sales == nil || len(sales.ReportedSales) != 1 {
		t.Fatal("Sales.ReportedSales missing")
	}
	rs := sales.ReportedSales[0]
	if rs.Value == nil || rs.Value.Value == nil {
		t.Fatal("ReportedSales[0].Value is nil")
	}
	if got, want := *rs.Value.Value, 95200000000.0; got != want {
		t.Errorf("ReportedSales[0].Value.Value = %v, want %v", got, want)
	}
	if rs.PeriodOffset == nil || *rs.PeriodOffset != "CURRENT" {
		t.Errorf("ReportedSales[0].PeriodOffset = %v, want CURRENT", rs.PeriodOffset)
	}
}

func TestFundamentalsEPSEstimates(t *testing.T) {
	t.Parallel()
	item := fundamentalsFixture(t)

	if item.Financials == nil || item.Financials.Estimates == nil {
		t.Fatal("Financials.Estimates is nil")
	}
	if got := len(item.Financials.Estimates.EPSEstimates); got != 1 {
		t.Fatalf("len(EPSEstimates) = %d, want 1", got)
	}
	est := item.Financials.Estimates.EPSEstimates[0]
	if est.Value == nil || est.Value.Value == nil {
		t.Fatal("EPSEstimates[0].Value is nil")
	}
	if got, want := *est.Value.Value, 1.72; got != want {
		t.Errorf("EPSEstimates[0].Value.Value = %v, want %v", got, want)
	}
	if est.RevisionDirection == nil || *est.RevisionDirection != "UP" {
		t.Errorf("EPSEstimates[0].RevisionDirection = %v, want UP", est.RevisionDirection)
	}
	if est.Period == nil || *est.Period != "P1Q" {
		t.Errorf("EPSEstimates[0].Period = %v, want P1Q", est.Period)
	}
}

func TestFundamentalsSalesEstimates(t *testing.T) {
	t.Parallel()
	item := fundamentalsFixture(t)

	if item.Financials == nil || item.Financials.Estimates == nil {
		t.Fatal("Financials.Estimates is nil")
	}
	if got := len(item.Financials.Estimates.SalesEstimates); got != 1 {
		t.Fatalf("len(SalesEstimates) = %d, want 1", got)
	}
	est := item.Financials.Estimates.SalesEstimates[0]
	if est.Value == nil || est.Value.Value == nil {
		t.Fatal("SalesEstimates[0].Value is nil")
	}
	if got, want := *est.Value.Value, 98500000000.0; got != want {
		t.Errorf("SalesEstimates[0].Value.Value = %v, want %v", got, want)
	}
	if est.Period == nil || *est.Period != "P1Q" {
		t.Errorf("SalesEstimates[0].Period = %v, want P1Q", est.Period)
	}
}

func TestFundamentalsRequestBody(t *testing.T) {
	t.Parallel()

	var gotReq struct {
		OperationName string `json:"operationName"`
		Variables     struct {
			Symbols                            []string `json:"symbols"`
			SymbolDialectType                  string   `json:"symbolDialectType"`
			UpToHistoricalPeriodOffset         string   `json:"upToHistoricalPeriodOffset"`
			UpToQueryPeriodOffset              string   `json:"upToQueryPeriodOffset"`
			ReportedSalesUpToHistoricalPeriod2 string   `json:"reportedSalesUpToHistoricalPeriod2"`
			SalesEstimatesUpToQueryPeriod2     string   `json:"salesEstimatesUpToQueryPeriod2"`
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
	req := NewFundamentalsRequest("AAPL")
	_, err := client.Fundamentals(context.Background(), req)
	if err != nil {
		t.Fatalf("Fundamentals() error = %v", err)
	}

	if got, want := gotReq.OperationName, "FundermentalDataBox"; got != want {
		t.Errorf("operationName = %q, want %q", got, want)
	}
	if len(gotReq.Variables.Symbols) != 1 || gotReq.Variables.Symbols[0] != "AAPL" {
		t.Errorf("variables.symbols = %v, want [AAPL]", gotReq.Variables.Symbols)
	}
	if got, want := gotReq.Variables.SymbolDialectType, DefaultFundamentalsSymbolDialectType; got != want {
		t.Errorf("variables.symbolDialectType = %q, want %q", got, want)
	}
	if got, want := gotReq.Variables.UpToHistoricalPeriodOffset, DefaultFundamentalsUpToHistoricalPeriodOffset; got != want {
		t.Errorf("variables.upToHistoricalPeriodOffset = %q, want %q", got, want)
	}
	if got, want := gotReq.Variables.UpToQueryPeriodOffset, DefaultFundamentalsUpToQueryPeriodOffset; got != want {
		t.Errorf("variables.upToQueryPeriodOffset = %q, want %q", got, want)
	}
	if got, want := gotReq.Variables.ReportedSalesUpToHistoricalPeriod2, DefaultFundamentalsReportedSalesUpToHistoricalPeriod2; got != want {
		t.Errorf("variables.reportedSalesUpToHistoricalPeriod2 = %q, want %q", got, want)
	}
	if got, want := gotReq.Variables.SalesEstimatesUpToQueryPeriod2, DefaultFundamentalsSalesEstimatesUpToQueryPeriod2; got != want {
		t.Errorf("variables.salesEstimatesUpToQueryPeriod2 = %q, want %q", got, want)
	}
}

func TestFundamentalsStatusError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte("service unavailable"))
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL)
	req := NewFundamentalsRequest("AAPL")
	_, err := client.Fundamentals(context.Background(), req)

	var statusErr *StatusError
	if !errors.As(err, &statusErr) {
		t.Fatalf("Fundamentals() error = %T, want *StatusError", err)
	}
	if got, want := statusErr.StatusCode, http.StatusServiceUnavailable; got != want {
		t.Errorf("StatusCode = %d, want %d", got, want)
	}
}

func TestFundamentalsDecodeError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("{invalid json"))
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL)
	req := NewFundamentalsRequest("AAPL")
	_, err := client.Fundamentals(context.Background(), req)

	var decodeErr *DecodeError
	if !errors.As(err, &decodeErr) {
		t.Fatalf("Fundamentals() error = %T, want *DecodeError", err)
	}
}

func TestFundamentalsGraphQLError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"errors":[{"message":"not authorized","extensions":{"code":"UNAUTHORIZED"}}]}`))
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL)
	req := NewFundamentalsRequest("AAPL")
	_, err := client.Fundamentals(context.Background(), req)

	var gqlErr *GraphQLError
	if !errors.As(err, &gqlErr) {
		t.Fatalf("Fundamentals() error = %T, want *GraphQLError", err)
	}
	if got, want := gqlErr.Error(), "not authorized"; got != want {
		t.Errorf("GraphQLError.Error() = %q, want %q", got, want)
	}
}

func TestFundamentalsBodyLimitError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"marketData":[]}}`))
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL, WithResponseBodyLimit(5))
	req := NewFundamentalsRequest("AAPL")
	_, err := client.Fundamentals(context.Background(), req)

	if !IsBodyLimit(err) {
		t.Fatalf("Fundamentals() error = %v, want BodyLimitError", err)
	}
}

func TestFundamentalsNoAuth(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		t.Error("Fundamentals() unexpectedly sent request without auth")
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(WithGraphQLURL(srv.URL))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	req := NewFundamentalsRequest("AAPL")
	_, err = client.Fundamentals(context.Background(), req)
	if err == nil {
		t.Fatal("Fundamentals() error = nil, want auth error")
	}
	if got, want := err.Error(), "no JWT provider configured"; !strings.Contains(got, want) {
		t.Errorf("error = %q, want substring %q", got, want)
	}
}

func TestFundamentalsFixtureSanitized(t *testing.T) {
	t.Parallel()

	for _, name := range []string{
		"testdata/FundermentalDataBox/request.json",
		"testdata/FundermentalDataBox/response.json",
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
