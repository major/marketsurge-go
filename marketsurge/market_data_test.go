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

// otherMarketDataFixture calls OtherMarketData against a test server that
// serves the response fixture and returns the first MarketDataItem.
func otherMarketDataFixture(t *testing.T) MarketDataItem {
	t.Helper()

	respBytes, err := os.ReadFile("testdata/OtherMarketData/response.json")
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
	req := NewOtherMarketDataRequest("AAPL")
	resp, err := client.OtherMarketData(context.Background(), req)
	if err != nil {
		t.Fatalf("OtherMarketData() error = %v", err)
	}
	if got, want := len(resp.MarketData), 1; got != want {
		t.Fatalf("len(MarketData) = %d, want %d", got, want)
	}
	return resp.MarketData[0]
}

func TestNewOtherMarketDataRequest(t *testing.T) {
	t.Parallel()

	req := NewOtherMarketDataRequest("AAPL")

	if got, want := req.SymbolDialectType, DefaultSymbolDialectType; got != want {
		t.Errorf("SymbolDialectType = %q, want %q", got, want)
	}
	if got, want := req.UpToHistoricalPeriodForProfitMargin, DefaultUpToHistoricalPeriodForProfitMargin; got != want {
		t.Errorf("UpToHistoricalPeriodForProfitMargin = %q, want %q", got, want)
	}
	if got, want := req.UpToHistoricalPeriodOffset, DefaultUpToHistoricalPeriodOffset; got != want {
		t.Errorf("UpToHistoricalPeriodOffset = %q, want %q", got, want)
	}
	if got, want := req.UpToQueryPeriodOffset, DefaultUpToQueryPeriodOffset; got != want {
		t.Errorf("UpToQueryPeriodOffset = %q, want %q", got, want)
	}
	if req.PatternStartDate == "" {
		t.Error("PatternStartDate is empty, want non-empty date")
	}
	if req.PatternEndDate == "" {
		t.Error("PatternEndDate is empty, want non-empty date")
	}
	if len(req.Symbols) != 1 || req.Symbols[0] != "AAPL" {
		t.Errorf("Symbols = %v, want [AAPL]", req.Symbols)
	}
}

func TestOtherMarketDataID(t *testing.T) {
	t.Parallel()
	item := otherMarketDataFixture(t)

	if item.ID == nil || *item.ID != "AAPL" {
		t.Errorf("MarketData[0].ID = %v, want AAPL", item.ID)
	}
}

func TestOtherMarketDataRatings(t *testing.T) {
	t.Parallel()
	item := otherMarketDataFixture(t)

	if item.Ratings == nil {
		t.Fatal("Ratings is nil")
	}
	if got := len(item.Ratings.CompRating); got != 1 {
		t.Fatalf("len(CompRating) = %d, want 1", got)
	}
	if got, want := *item.Ratings.CompRating[0].Value, 95; got != want {
		t.Errorf("CompRating[0].Value = %d, want %d", got, want)
	}
	if got := len(item.Ratings.SMRRating); got != 1 {
		t.Fatalf("len(SMRRating) = %d, want 1", got)
	}
	if item.Ratings.SMRRating[0].LetterValue == nil || *item.Ratings.SMRRating[0].LetterValue != "A" {
		t.Errorf("SMRRating[0].LetterValue = %v, want A", item.Ratings.SMRRating[0].LetterValue)
	}
}

func TestOtherMarketDataPricing(t *testing.T) {
	t.Parallel()
	item := otherMarketDataFixture(t)

	if item.PricingStatistics == nil || item.PricingStatistics.EndOfDayStatistics == nil {
		t.Fatal("PricingStatistics.EndOfDayStatistics is nil")
	}
	eod := item.PricingStatistics.EndOfDayStatistics
	if eod.MarketCapitalization == nil || eod.MarketCapitalization.Value == nil {
		t.Fatal("MarketCapitalization is nil")
	}
	if got, want := *eod.MarketCapitalization.Value, 3200000000000.0; got != want {
		t.Errorf("MarketCapitalization.Value = %v, want %v", got, want)
	}
}

func TestOtherMarketDataPatterns(t *testing.T) {
	t.Parallel()
	item := otherMarketDataFixture(t)

	if item.PatternInfo == nil {
		t.Fatal("PatternInfo is nil")
	}
	if got := len(item.PatternInfo.Patterns); got != 1 {
		t.Fatalf("len(Patterns) = %d, want 1", got)
	}
	pattern := item.PatternInfo.Patterns[0]
	if pattern.PatternType == nil || *pattern.PatternType != "CUP_WITH_HANDLE" {
		t.Errorf("PatternType = %v, want CUP_WITH_HANDLE", pattern.PatternType)
	}
	if pattern.HandleLength == nil || *pattern.HandleLength != 12 {
		t.Errorf("HandleLength = %v, want 12", pattern.HandleLength)
	}
	if pattern.PivotPrice == nil || pattern.PivotPrice.CurrencySymbolInfo == nil {
		t.Fatal("PivotPrice.CurrencySymbolInfo is nil")
	}
	if got, want := *pattern.PivotPrice.CurrencySymbolInfo.IsoCurrencyCode, "USD"; got != want {
		t.Errorf("IsoCurrencyCode = %q, want %q", got, want)
	}
	if got := len(item.PatternInfo.TightAreas); got != 1 {
		t.Fatalf("len(TightAreas) = %d, want 1", got)
	}
	if item.PatternInfo.TightAreas[0].PatternID == nil || *item.PatternInfo.TightAreas[0].PatternID != 1 {
		t.Errorf("TightAreas[0].PatternID = %v, want 1", item.PatternInfo.TightAreas[0].PatternID)
	}
}

func TestOtherMarketDataFinancials(t *testing.T) {
	t.Parallel()
	item := otherMarketDataFixture(t)

	if item.Financials == nil || item.Financials.ConsensusFinancials == nil {
		t.Fatal("Financials.ConsensusFinancials is nil")
	}
	eps := item.Financials.ConsensusFinancials.EPS
	if eps == nil || len(eps.ReportedEarnings) != 1 {
		t.Fatal("EPS.ReportedEarnings missing")
	}
	if eps.ReportedEarnings[0].Value == nil || eps.ReportedEarnings[0].Value.Value == nil {
		t.Fatal("ReportedEarnings[0].Value is nil")
	}
	if got, want := *eps.ReportedEarnings[0].Value.Value, 1.65; got != want {
		t.Errorf("ReportedEarnings[0].Value.Value = %v, want %v", got, want)
	}
	if eps.EarningsStability == nil || *eps.EarningsStability != 3 {
		t.Errorf("EarningsStability = %v, want 3", eps.EarningsStability)
	}
}

func TestOtherMarketDataIndustry(t *testing.T) {
	t.Parallel()
	item := otherMarketDataFixture(t)

	if item.Industry == nil || item.Industry.Name == nil {
		t.Fatal("Industry.Name is nil")
	}
	if got, want := *item.Industry.Name, "Comp-Hardware/Peripherals"; got != want {
		t.Errorf("Industry.Name = %q, want %q", got, want)
	}
	if item.Industry.NumberOfStocksInGroup == nil || *item.Industry.NumberOfStocksInGroup != 25 {
		t.Errorf("NumberOfStocksInGroup = %v, want 25", item.Industry.NumberOfStocksInGroup)
	}
}

func TestOtherMarketDataIndustryNumericIndCode(t *testing.T) {
	t.Parallel()

	const response = `{
		"data": {
			"marketData": [
				{
					"id": "AMD",
					"industry": {
						"name": "Electronics-Semiconductor Fabless",
						"indCode": 123
					}
				}
			]
		}
	}`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(response))
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL)
	req := NewOtherMarketDataRequest("AMD")
	resp, err := client.OtherMarketData(context.Background(), req)
	if err != nil {
		t.Fatalf("OtherMarketData() error = %v", err)
	}
	if got, want := len(resp.MarketData), 1; got != want {
		t.Fatalf("len(MarketData) = %d, want %d", got, want)
	}

	industry := resp.MarketData[0].Industry
	if industry == nil || industry.IndCode == nil {
		t.Fatal("MarketData[0].Industry.IndCode is nil")
	}
	if got, want := *industry.IndCode, "123"; got != want {
		t.Errorf("MarketData[0].Industry.IndCode = %q, want %q", got, want)
	}
}

func TestOtherMarketDataSymbology(t *testing.T) {
	t.Parallel()
	item := otherMarketDataFixture(t)

	if item.Symbology == nil || item.Symbology.Company == nil {
		t.Fatal("Symbology.Company is nil")
	}
	if item.Symbology.Company.CompanyName == nil || *item.Symbology.Company.CompanyName != "Apple Inc." {
		t.Errorf("CompanyName = %v, want Apple Inc.", item.Symbology.Company.CompanyName)
	}
}

func TestOtherMarketDataSymbologyCompanyArray(t *testing.T) {
	t.Parallel()

	const response = `{
		"data": {
			"marketData": [
				{
					"id": "AMD",
					"symbology": {
						"company": [
							{
								"companyName": "Advanced Micro Devices, Inc.",
								"address": "2485 Augustine Drive",
								"address2": null,
								"phone": "408-749-4000",
								"businessDescription": "Designs semiconductor products.",
								"url": "https://www.amd.com",
								"city": "Santa Clara",
								"country": "US",
								"stateProvince": "CA"
							}
						]
					}
				}
			]
		}
	}`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(response))
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL)
	req := NewOtherMarketDataRequest("AMD")
	resp, err := client.OtherMarketData(context.Background(), req)
	if err != nil {
		t.Fatalf("OtherMarketData() error = %v", err)
	}
	if got, want := len(resp.MarketData), 1; got != want {
		t.Fatalf("len(MarketData) = %d, want %d", got, want)
	}

	if resp.MarketData[0].Symbology == nil {
		t.Fatal("MarketData[0].Symbology is nil")
	}
	company := resp.MarketData[0].Symbology.Company
	if company == nil || company.CompanyName == nil {
		t.Fatal("MarketData[0].Symbology.Company.CompanyName is nil")
	}
	if got, want := *company.CompanyName, "Advanced Micro Devices, Inc."; got != want {
		t.Errorf("MarketData[0].Symbology.Company.CompanyName = %q, want %q", got, want)
	}
}

func TestOtherMarketDataSymbologyInstrumentArray(t *testing.T) {
	t.Parallel()

	const response = `{
		"data": {
			"marketData": [
				{
					"id": "AMD",
					"symbology": {
						"instrument": [
							{
								"subType": "COMMON_STOCK",
								"ipoDate": {"value": "1972-09-27"}
							}
						]
					}
				}
			]
		}
	}`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(response))
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL)
	req := NewOtherMarketDataRequest("AMD")
	resp, err := client.OtherMarketData(context.Background(), req)
	if err != nil {
		t.Fatalf("OtherMarketData() error = %v", err)
	}
	if got, want := len(resp.MarketData), 1; got != want {
		t.Fatalf("len(MarketData) = %d, want %d", got, want)
	}

	if resp.MarketData[0].Symbology == nil {
		t.Fatal("MarketData[0].Symbology is nil")
	}
	instrument := resp.MarketData[0].Symbology.Instrument
	if instrument == nil || instrument.SubType == nil {
		t.Fatal("MarketData[0].Symbology.Instrument.SubType is nil")
	}
	if got, want := *instrument.SubType, "COMMON_STOCK"; got != want {
		t.Errorf("MarketData[0].Symbology.Instrument.SubType = %q, want %q", got, want)
	}
	if instrument.IPODate == nil {
		t.Fatal("MarketData[0].Symbology.Instrument.IPODate is nil")
	}
	if got, want := instrument.IPODate.Value, "1972-09-27"; got != want {
		t.Errorf("MarketData[0].Symbology.Instrument.IPODate.Value = %q, want %q", got, want)
	}
}

func TestOtherMarketDataOwnership(t *testing.T) {
	t.Parallel()
	item := otherMarketDataFixture(t)

	if item.Ownership == nil || item.Ownership.FundsFloatPercentHeld == nil {
		t.Fatal("Ownership.FundsFloatPercentHeld is nil")
	}
	if got, want := *item.Ownership.FundsFloatPercentHeld.Value, 62.5; got != want {
		t.Errorf("FundsFloatPercentHeld.Value = %v, want %v", got, want)
	}
}

func TestOtherMarketDataRequestBody(t *testing.T) {
	t.Parallel()

	var gotReq struct {
		OperationName string `json:"operationName"`
		Variables     struct {
			Symbols                             []string `json:"symbols"`
			SymbolDialectType                   string   `json:"symbolDialectType"`
			UpToHistoricalPeriodForProfitMargin string   `json:"upToHistoricalPeriodForProfitMargin"`
			UpToHistoricalPeriodOffset          string   `json:"upToHistoricalPeriodOffset"`
			UpToQueryPeriodOffset               string   `json:"upToQueryPeriodOffset"`
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
	req := NewOtherMarketDataRequest("AAPL")
	_, err := client.OtherMarketData(context.Background(), req)
	if err != nil {
		t.Fatalf("OtherMarketData() error = %v", err)
	}

	if got, want := gotReq.OperationName, "OtherMarketData"; got != want {
		t.Errorf("operationName = %q, want %q", got, want)
	}
	if len(gotReq.Variables.Symbols) != 1 || gotReq.Variables.Symbols[0] != "AAPL" {
		t.Errorf("variables.symbols = %v, want [AAPL]", gotReq.Variables.Symbols)
	}
	if got, want := gotReq.Variables.SymbolDialectType, DefaultSymbolDialectType; got != want {
		t.Errorf("variables.symbolDialectType = %q, want %q", got, want)
	}
	if got, want := gotReq.Variables.UpToHistoricalPeriodForProfitMargin, DefaultUpToHistoricalPeriodForProfitMargin; got != want {
		t.Errorf("variables.upToHistoricalPeriodForProfitMargin = %q, want %q", got, want)
	}
	if got, want := gotReq.Variables.UpToHistoricalPeriodOffset, DefaultUpToHistoricalPeriodOffset; got != want {
		t.Errorf("variables.upToHistoricalPeriodOffset = %q, want %q", got, want)
	}
	if got, want := gotReq.Variables.UpToQueryPeriodOffset, DefaultUpToQueryPeriodOffset; got != want {
		t.Errorf("variables.upToQueryPeriodOffset = %q, want %q", got, want)
	}
	// Verify pattern date placeholders were replaced in the query.
	if strings.Contains(gotReq.Query, "{pattern_start_date}") {
		t.Error("query still contains {pattern_start_date} placeholder")
	}
	if strings.Contains(gotReq.Query, "{pattern_end_date}") {
		t.Error("query still contains {pattern_end_date} placeholder")
	}
}

func TestOtherMarketDataStatusError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte("service unavailable"))
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL)
	req := NewOtherMarketDataRequest("AAPL")
	_, err := client.OtherMarketData(context.Background(), req)

	var statusErr *StatusError
	if !errors.As(err, &statusErr) {
		t.Fatalf("OtherMarketData() error = %T, want *StatusError", err)
	}
	if got, want := statusErr.StatusCode, http.StatusServiceUnavailable; got != want {
		t.Errorf("StatusCode = %d, want %d", got, want)
	}
}

func TestOtherMarketDataDecodeError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("{invalid json"))
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL)
	req := NewOtherMarketDataRequest("AAPL")
	_, err := client.OtherMarketData(context.Background(), req)

	var decodeErr *DecodeError
	if !errors.As(err, &decodeErr) {
		t.Fatalf("OtherMarketData() error = %T, want *DecodeError", err)
	}
}

func TestOtherMarketDataGraphQLError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"errors":[{"message":"not authorized","extensions":{"code":"UNAUTHORIZED"}}]}`))
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL)
	req := NewOtherMarketDataRequest("AAPL")
	_, err := client.OtherMarketData(context.Background(), req)

	var gqlErr *GraphQLError
	if !errors.As(err, &gqlErr) {
		t.Fatalf("OtherMarketData() error = %T, want *GraphQLError", err)
	}
	if got, want := gqlErr.Error(), "not authorized"; got != want {
		t.Errorf("GraphQLError.Error() = %q, want %q", got, want)
	}
}

func TestOtherMarketDataBodyLimitError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"marketData":[]}}`))
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL, WithResponseBodyLimit(5))
	req := NewOtherMarketDataRequest("AAPL")
	_, err := client.OtherMarketData(context.Background(), req)

	if !IsBodyLimit(err) {
		t.Fatalf("OtherMarketData() error = %v, want BodyLimitError", err)
	}
}

func TestOtherMarketDataNoAuth(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		t.Error("OtherMarketData() unexpectedly sent request without auth")
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(WithGraphQLURL(srv.URL))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	req := NewOtherMarketDataRequest("AAPL")
	_, err = client.OtherMarketData(context.Background(), req)
	if err == nil {
		t.Fatal("OtherMarketData() error = nil, want auth error")
	}
	if got, want := err.Error(), "no JWT provider configured"; !strings.Contains(got, want) {
		t.Errorf("error = %q, want substring %q", got, want)
	}
}

func TestOtherMarketDataFixtureSanitized(t *testing.T) {
	t.Parallel()

	for _, name := range []string{
		"testdata/OtherMarketData/request.json",
		"testdata/OtherMarketData/response.json",
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
