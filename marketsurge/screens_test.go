package marketsurge

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// ---------------------------------------------------------------------------
// MarketDataAdhocScreen fixtures and tests
// ---------------------------------------------------------------------------

// adhocScreenFixture calls MarketDataAdhocScreen against a test server that
// serves the response fixture and returns the AdhocScreenResult.
func adhocScreenFixture(t *testing.T) *AdhocScreenResult {
	t.Helper()

	respBytes, err := os.ReadFile("testdata/MarketDataAdhocScreen/response.json")
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
	columns := []AdhocScreenResponseColumn{
		{Name: "Symbol"},
		{Name: "CompanyName"},
	}
	req := NewMarketDataAdhocScreenRequest(columns)
	req.IncludeSource = AdhocScreenIncludeSource{
		ScreenID: &AdhocScreenID{ID: 46, Dialect: "MS_LIST_ID"},
	}
	resp, err := client.MarketDataAdhocScreen(context.Background(), req)
	if err != nil {
		t.Fatalf("MarketDataAdhocScreen() error = %v", err)
	}
	if resp.MarketDataAdhocScreen == nil {
		t.Fatal("MarketDataAdhocScreen response is nil")
	}
	return resp.MarketDataAdhocScreen
}

func TestNewMarketDataAdhocScreenRequest(t *testing.T) {
	t.Parallel()

	columns := AdhocScreenColumns(ColumnSymbol)
	req := NewMarketDataAdhocScreenRequest(columns)

	if got, want := req.CorrelationTag, DefaultAdhocScreenCorrelationTag; got != want {
		t.Errorf("CorrelationTag = %q, want %q", got, want)
	}
	if got, want := req.PageSize, DefaultAdhocScreenPageSize; got != want {
		t.Errorf("PageSize = %d, want %d", got, want)
	}
	if got, want := req.ResultLimit, DefaultAdhocScreenResultLimit; got != want {
		t.Errorf("ResultLimit = %d, want %d", got, want)
	}
	if got, want := req.PageSkip, 0; got != want {
		t.Errorf("PageSkip = %d, want %d", got, want)
	}
	if got, want := req.ResultType, DefaultAdhocScreenResultType; got != want {
		t.Errorf("ResultType = %q, want %q", got, want)
	}
	if got, want := len(req.ResponseColumns), 1; got != want {
		t.Errorf("len(ResponseColumns) = %d, want %d", got, want)
	}
	if req.AdhocQuery != nil {
		t.Errorf("AdhocQuery = %v, want nil", req.AdhocQuery)
	}
}

func TestNewMarketDataAdhocScreenRequestCopiesColumns(t *testing.T) {
	t.Parallel()

	columns := AdhocScreenColumns(ColumnSymbol)
	req := NewMarketDataAdhocScreenRequest(columns)
	columns[0].Name = ColumnPrice

	if got, want := req.ResponseColumns[0].Name, ColumnName(ColumnSymbol); got != want {
		t.Errorf("ResponseColumns[0].Name after caller mutation = %q, want %q", got, want)
	}
}

func TestMarketDataAdhocScreenCorrelationTag(t *testing.T) {
	t.Parallel()

	result := adhocScreenFixture(t)

	if result.CorrelationTag == nil {
		t.Fatal("CorrelationTag is nil")
	}
	if got, want := *result.CorrelationTag, "marketsurge"; got != want {
		t.Errorf("CorrelationTag = %q, want %q", got, want)
	}
}

func TestMarketDataAdhocScreenElapsedTime(t *testing.T) {
	t.Parallel()

	result := adhocScreenFixture(t)

	if result.ElapsedTime == nil {
		t.Fatal("ElapsedTime is nil")
	}
	if got, want := *result.ElapsedTime, "42ms"; got != want {
		t.Errorf("ElapsedTime = %q, want %q", got, want)
	}
}

func TestMarketDataAdhocScreenCounts(t *testing.T) {
	t.Parallel()

	result := adhocScreenFixture(t)

	if result.NumberOfInstrumentsInSource == nil {
		t.Fatal("NumberOfInstrumentsInSource is nil")
	}
	if got, want := *result.NumberOfInstrumentsInSource, 8500; got != want {
		t.Errorf("NumberOfInstrumentsInSource = %d, want %d", got, want)
	}
	if result.NumberOfMatchingInstruments == nil {
		t.Fatal("NumberOfMatchingInstruments is nil")
	}
	if got, want := *result.NumberOfMatchingInstruments, 2; got != want {
		t.Errorf("NumberOfMatchingInstruments = %d, want %d", got, want)
	}
}

func TestMarketDataAdhocScreenAdhocQueryString(t *testing.T) {
	t.Parallel()

	result := adhocScreenFixture(t)

	if result.AdhocQueryString == nil {
		t.Fatal("AdhocQueryString is nil")
	}
	if got, want := *result.AdhocQueryString, "CompositeRating >= 90"; got != want {
		t.Errorf("AdhocQueryString = %q, want %q", got, want)
	}
}

func TestMarketDataAdhocScreenAdhocQuery(t *testing.T) {
	t.Parallel()

	result := adhocScreenFixture(t)

	if result.AdhocQuery == nil {
		t.Fatal("AdhocQuery is nil")
	}
	if got, want := len(result.AdhocQuery.Terms), 1; got != want {
		t.Fatalf("len(AdhocQuery.Terms) = %d, want %d", got, want)
	}

	term := result.AdhocQuery.Terms[0]
	if term.NumberOfMatchingInstruments == nil {
		t.Fatal("term.NumberOfMatchingInstruments is nil")
	}
	if got, want := *term.NumberOfMatchingInstruments, 2; got != want {
		t.Errorf("term.NumberOfMatchingInstruments = %d, want %d", got, want)
	}
	if term.Ordinal == nil {
		t.Fatal("term.Ordinal is nil")
	}
	if got, want := *term.Ordinal, 1; got != want {
		t.Errorf("term.Ordinal = %d, want %d", got, want)
	}
	assertStringPtr(t, "term.Operand", term.Operand, ">=")

	if term.Left == nil {
		t.Fatal("term.Left is nil")
	}
	assertStringPtr(t, "term.Left.Name", term.Left.Name, "CompositeRating")
	if term.Left.MDItemID == nil {
		t.Fatal("term.Left.MDItemID is nil")
	}
	if got, want := term.Left.MDItemID.String(), "2001"; got != want {
		t.Errorf("term.Left.MDItemID = %q, want %q", got, want)
	}

	if term.Right == nil {
		t.Fatal("term.Right is nil")
	}
	assertStringPtr(t, "term.Right.Value", term.Right.Value, "90")
	if term.Right.MaximumValue != nil {
		t.Errorf("term.Right.MaximumValue = %v, want nil", *term.Right.MaximumValue)
	}
	if term.Right.MinimumValue != nil {
		t.Errorf("term.Right.MinimumValue = %v, want nil", *term.Right.MinimumValue)
	}
}

func TestMarketDataAdhocScreenResponseValues(t *testing.T) {
	t.Parallel()

	result := adhocScreenFixture(t)

	if got, want := len(result.ResponseValues), 2; got != want {
		t.Fatalf("len(ResponseValues) = %d, want %d", got, want)
	}

	row := result.ResponseValues[0]
	if got, want := len(row), 2; got != want {
		t.Fatalf("len(row[0]) = %d, want %d", got, want)
	}

	cell := row[0]
	if cell.Value == nil {
		t.Fatal("cell.Value is nil")
	}
	if got, want := *cell.Value, "AAPL"; got != want {
		t.Errorf("cell.Value = %q, want %q", got, want)
	}
	if cell.MDItem == nil {
		t.Fatal("cell.MDItem is nil")
	}
	if cell.MDItem.Name == nil {
		t.Fatal("cell.MDItem.Name is nil")
	}
	if got, want := *cell.MDItem.Name, "Symbol"; got != want {
		t.Errorf("cell.MDItem.Name = %q, want %q", got, want)
	}
	if cell.MDItem.MDItemID == nil {
		t.Fatal("cell.MDItem.MDItemID is nil")
	}
	if got, want := cell.MDItem.MDItemID.String(), "1001"; got != want {
		t.Errorf("cell.MDItem.MDItemID = %q, want %q", got, want)
	}
}

func TestMarketDataAdhocScreenRequestBody(t *testing.T) {
	t.Parallel()

	wantBytes, err := os.ReadFile("testdata/MarketDataAdhocScreen/request.json")
	if err != nil {
		t.Fatalf("os.ReadFile(request.json) error = %v", err)
	}

	var captured []byte
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captured, _ = readBodyUpToLimit(r.Body, 1<<20)
		respBytes, _ := os.ReadFile("testdata/MarketDataAdhocScreen/response.json")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(respBytes)
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL)
	columns := []AdhocScreenResponseColumn{
		{Name: "Symbol"},
		{Name: "CompanyName"},
	}
	req := NewMarketDataAdhocScreenRequest(columns)
	req.IncludeSource = AdhocScreenIncludeSource{
		ScreenID: &AdhocScreenID{ID: 46, Dialect: "MS_LIST_ID"},
	}
	_, err = client.MarketDataAdhocScreen(context.Background(), req)
	if err != nil {
		t.Fatalf("MarketDataAdhocScreen() error = %v", err)
	}

	var want, got any
	if err = json.Unmarshal(wantBytes, &want); err != nil {
		t.Fatalf("unmarshal want: %v", err)
	}
	if err = json.Unmarshal(captured, &got); err != nil {
		t.Fatalf("unmarshal got: %v", err)
	}
	wantJSON, _ := json.Marshal(want)
	gotJSON, _ := json.Marshal(got)
	if string(wantJSON) != string(gotJSON) {
		t.Errorf("request body mismatch:\ngot:  %s\nwant: %s", gotJSON, wantJSON)
	}
}

func TestMarketDataAdhocScreenStatusError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL)
	req := NewMarketDataAdhocScreenRequest(nil)
	_, err := client.MarketDataAdhocScreen(context.Background(), req)
	assertErrorType(t, err, &StatusError{})
}

func TestMarketDataAdhocScreenDecodeError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`not json`))
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL)
	req := NewMarketDataAdhocScreenRequest(nil)
	_, err := client.MarketDataAdhocScreen(context.Background(), req)
	assertErrorType(t, err, &DecodeError{})
}

func TestMarketDataAdhocScreenGraphQLError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"errors":[{"message":"bad request"}]}`))
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL)
	req := NewMarketDataAdhocScreenRequest(nil)
	_, err := client.MarketDataAdhocScreen(context.Background(), req)
	assertErrorType(t, err, &GraphQLError{})
}

func TestMarketDataAdhocScreenBodyLimitError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"marketDataAdhocScreen":` + strings.Repeat("x", 1024) + `}}`))
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL, WithResponseBodyLimit(64))
	req := NewMarketDataAdhocScreenRequest(nil)
	_, err := client.MarketDataAdhocScreen(context.Background(), req)
	assertErrorType(t, err, &BodyLimitError{})
}

func TestMarketDataAdhocScreenNoAuth(t *testing.T) {
	t.Parallel()

	client, err := NewClient()
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	req := NewMarketDataAdhocScreenRequest(nil)
	_, err = client.MarketDataAdhocScreen(context.Background(), req)
	if err == nil {
		t.Fatal("expected error for missing auth, got nil")
	}
}

func TestMarketDataAdhocScreenFixtureSanitized(t *testing.T) {
	t.Parallel()

	for _, name := range []string{
		"testdata/MarketDataAdhocScreen/request.json",
		"testdata/MarketDataAdhocScreen/response.json",
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
// RunScreen fixtures and tests
// ---------------------------------------------------------------------------

// runScreenFixture calls RunScreen against a test server that serves
// the response fixture and returns the RunScreenResult.
func runScreenFixture(t *testing.T) *RunScreenResult {
	t.Helper()

	respBytes, err := os.ReadFile("testdata/RunScreen/response.json")
	if err != nil {
		t.Fatalf("os.ReadFile(response.json) error = %v", err)
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.Method, http.MethodPost; got != want {
			t.Errorf("request method = %s, want %s", got, want)
		}
		assertHeader(t, r.Header, "Authorization", "Bearer "+testJWT)
		assertHeader(t, r.Header, "Content-Type", "application/json")

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(respBytes)
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL)
	columns := []RunScreenResponseColumn{
		{Name: "Symbol"},
		{Name: "CompanyName"},
	}
	req := NewRunScreenRequest("screen-abc-123", columns)
	resp, err := client.RunScreen(context.Background(), req)
	if err != nil {
		t.Fatalf("RunScreen() error = %v", err)
	}
	if resp.User == nil || resp.User.RunScreen == nil {
		t.Fatal("RunScreen response User or RunScreen is nil")
	}
	return resp.User.RunScreen
}

func TestNewRunScreenRequest(t *testing.T) {
	t.Parallel()

	req := NewRunScreenRequest("screen-123", RunScreenColumns(ColumnSymbol))

	if got, want := req.Input.CorrelationTag, DefaultRunScreenCorrelationTag; got != want {
		t.Errorf("CorrelationTag = %q, want %q", got, want)
	}
	if got, want := req.Input.CoachAccount, true; got != want {
		t.Errorf("CoachAccount = %v, want %v", got, want)
	}
	if got, want := req.Input.PageSize, DefaultRunScreenPageSize; got != want {
		t.Errorf("PageSize = %d, want %d", got, want)
	}
	if got, want := req.Input.ResultLimit, DefaultRunScreenResultLimit; got != want {
		t.Errorf("ResultLimit = %d, want %d", got, want)
	}
	if got, want := req.Input.ScreenID, "screen-123"; got != want {
		t.Errorf("ScreenID = %q, want %q", got, want)
	}
	if got, want := req.Input.Site, DefaultRunScreenSite; got != want {
		t.Errorf("Site = %q, want %q", got, want)
	}
	if got, want := req.Input.Skip, 0; got != want {
		t.Errorf("Skip = %d, want %d", got, want)
	}
	if req.Input.IncludeSource == nil {
		t.Fatal("IncludeSource is nil")
	}
	if req.Input.IncludeSource.Source != nil {
		t.Errorf("IncludeSource.Source = %v, want nil", req.Input.IncludeSource.Source)
	}
	if got, want := len(req.Input.ResponseColumns), 1; got != want {
		t.Errorf("len(ResponseColumns) = %d, want %d", got, want)
	}
	if got, want := req.Input.ResponseColumns[0].Name, ColumnName(ColumnSymbol); got != want {
		t.Errorf("ResponseColumns[0].Name = %q, want %q", got, want)
	}
}

func TestNewRunScreenRequestCopiesColumns(t *testing.T) {
	t.Parallel()

	columns := RunScreenColumns(ColumnSymbol)
	req := NewRunScreenRequest("screen-123", columns)
	columns[0].Name = ColumnPrice

	if got, want := req.Input.ResponseColumns[0].Name, ColumnName(ColumnSymbol); got != want {
		t.Errorf("ResponseColumns[0].Name after caller mutation = %q, want %q", got, want)
	}
}

func TestRunScreenMatchingInstruments(t *testing.T) {
	t.Parallel()

	result := runScreenFixture(t)

	if result.NumberOfMatchingInstruments == nil {
		t.Fatal("NumberOfMatchingInstruments is nil")
	}
	if got, want := *result.NumberOfMatchingInstruments, 3; got != want {
		t.Errorf("NumberOfMatchingInstruments = %d, want %d", got, want)
	}
}

func TestRunScreenResponseValues(t *testing.T) {
	t.Parallel()

	result := runScreenFixture(t)

	if got, want := len(result.ResponseValues), 2; got != want {
		t.Fatalf("len(ResponseValues) = %d, want %d", got, want)
	}

	row := result.ResponseValues[0]
	if got, want := len(row), 2; got != want {
		t.Fatalf("len(row[0]) = %d, want %d", got, want)
	}

	cell := row[0]
	if cell.Value == nil {
		t.Fatal("cell.Value is nil")
	}
	if got, want := *cell.Value, "NVDA"; got != want {
		t.Errorf("cell.Value = %q, want %q", got, want)
	}
	if cell.MDItem == nil {
		t.Fatal("cell.MDItem is nil")
	}
	if cell.MDItem.Name == nil {
		t.Fatal("cell.MDItem.Name is nil")
	}
	if got, want := *cell.MDItem.Name, "Symbol"; got != want {
		t.Errorf("cell.MDItem.Name = %q, want %q", got, want)
	}
	if cell.MDItem.MDItemID == nil {
		t.Fatal("cell.MDItem.MDItemID is nil")
	}
	if got, want := cell.MDItem.MDItemID.String(), "2001"; got != want {
		t.Errorf("cell.MDItem.MDItemID = %q, want %q", got, want)
	}
}

func TestRunScreenRequestBody(t *testing.T) {
	t.Parallel()

	wantBytes, err := os.ReadFile("testdata/RunScreen/request.json")
	if err != nil {
		t.Fatalf("os.ReadFile(request.json) error = %v", err)
	}

	var captured []byte
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captured, _ = readBodyUpToLimit(r.Body, 1<<20)
		respBytes, _ := os.ReadFile("testdata/RunScreen/response.json")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(respBytes)
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL)
	columns := []RunScreenResponseColumn{
		{Name: "Symbol"},
		{Name: "CompanyName"},
	}
	req := NewRunScreenRequest("screen-abc-123", columns)
	_, err = client.RunScreen(context.Background(), req)
	if err != nil {
		t.Fatalf("RunScreen() error = %v", err)
	}

	var want, got any
	if err = json.Unmarshal(wantBytes, &want); err != nil {
		t.Fatalf("unmarshal want: %v", err)
	}
	if err = json.Unmarshal(captured, &got); err != nil {
		t.Fatalf("unmarshal got: %v", err)
	}
	wantJSON, _ := json.Marshal(want)
	gotJSON, _ := json.Marshal(got)
	if string(wantJSON) != string(gotJSON) {
		t.Errorf("request body mismatch:\ngot:  %s\nwant: %s", gotJSON, wantJSON)
	}
}

func TestRunScreenStatusError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL)
	req := NewRunScreenRequest("screen-123", nil)
	_, err := client.RunScreen(context.Background(), req)
	assertErrorType(t, err, &StatusError{})
}

func TestRunScreenDecodeError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`not json`))
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL)
	req := NewRunScreenRequest("screen-123", nil)
	_, err := client.RunScreen(context.Background(), req)
	assertErrorType(t, err, &DecodeError{})
}

func TestRunScreenGraphQLError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"errors":[{"message":"bad request"}]}`))
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL)
	req := NewRunScreenRequest("screen-123", nil)
	_, err := client.RunScreen(context.Background(), req)
	assertErrorType(t, err, &GraphQLError{})
}

func TestRunScreenBodyLimitError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"user":{"runScreen":` + strings.Repeat("x", 1024) + `}}}`))
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL, WithResponseBodyLimit(64))
	req := NewRunScreenRequest("screen-123", nil)
	_, err := client.RunScreen(context.Background(), req)
	assertErrorType(t, err, &BodyLimitError{})
}

func TestRunScreenNoAuth(t *testing.T) {
	t.Parallel()

	client, err := NewClient()
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	req := NewRunScreenRequest("screen-123", nil)
	_, err = client.RunScreen(context.Background(), req)
	if err == nil {
		t.Fatal("expected error for missing auth, got nil")
	}
}

func TestRunScreenFixtureSanitized(t *testing.T) {
	t.Parallel()

	for _, name := range []string{
		"testdata/RunScreen/request.json",
		"testdata/RunScreen/response.json",
	} {
		data, err := os.ReadFile(name)
		if err != nil {
			t.Fatalf("os.ReadFile(%s) error = %v", name, err)
		}
		content := strings.ToLower(string(data))
		for _, secret := range []string{"cookie", "session_id", "password"} {
			if strings.Contains(content, secret) {
				t.Errorf("fixture %s contains potential secret %q", name, secret)
			}
		}
	}
}

// ---------------------------------------------------------------------------
// Screens fixtures and tests
// ---------------------------------------------------------------------------

// screensFixture calls Screens against a test server that serves
// the response fixture and returns the screen entries.
func screensFixture(t *testing.T) []ScreenEntry {
	t.Helper()

	respBytes, err := os.ReadFile("testdata/Screens/response.json")
	if err != nil {
		t.Fatalf("os.ReadFile(response.json) error = %v", err)
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.Method, http.MethodPost; got != want {
			t.Errorf("request method = %s, want %s", got, want)
		}
		assertHeader(t, r.Header, "Authorization", "Bearer "+testJWT)
		assertHeader(t, r.Header, "Content-Type", "application/json")

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(respBytes)
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL)
	req := NewScreensRequest()
	resp, err := client.Screens(context.Background(), req)
	if err != nil {
		t.Fatalf("Screens() error = %v", err)
	}
	if resp.User == nil {
		t.Fatal("Screens response User is nil")
	}
	return resp.User.Screens
}

func TestNewScreensRequest(t *testing.T) {
	t.Parallel()

	req := NewScreensRequest()

	if got, want := req.Site, DefaultScreensSite; got != want {
		t.Errorf("Site = %q, want %q", got, want)
	}
	if req.Type != nil {
		t.Errorf("Type = %v, want nil", req.Type)
	}
	if req.SortDir != nil {
		t.Errorf("SortDir = %v, want nil", req.SortDir)
	}
}

func TestScreensEntries(t *testing.T) {
	t.Parallel()

	screens := screensFixture(t)

	if got, want := len(screens), 2; got != want {
		t.Fatalf("len(screens) = %d, want %d", got, want)
	}

	s := screens[0]
	if s.Name == nil {
		t.Fatal("Name is nil")
	}
	if got, want := *s.Name, "Growth Leaders"; got != want {
		t.Errorf("Name = %q, want %q", got, want)
	}
	if s.ID == nil {
		t.Fatal("ID is nil")
	}
	if got, want := *s.ID, "scr-001"; got != want {
		t.Errorf("ID = %q, want %q", got, want)
	}
	if s.Type == nil {
		t.Fatal("Type is nil")
	}
	if got, want := *s.Type, "CUSTOM"; got != want {
		t.Errorf("Type = %q, want %q", got, want)
	}
}

func TestScreensFilterCriteria(t *testing.T) {
	t.Parallel()

	screens := screensFixture(t)
	s := screens[0]

	if s.FilterCriteria == nil {
		t.Fatal("FilterCriteria is nil")
	}
	assertStringPtr(t, "FilterCriteria.Type", s.FilterCriteria.Type, "AND")
	if got, want := len(s.FilterCriteria.Terms), 1; got != want {
		t.Fatalf("len(FilterCriteria.Terms) = %d, want %d", got, want)
	}

	term := s.FilterCriteria.Terms[0]
	if term.Left == nil {
		t.Fatal("Terms[0].Left is nil")
	}
	assertStringPtr(t, "Terms[0].Left.Name", term.Left.Name, "CompositeRating")
	assertStringPtr(t, "Terms[0].Operand", term.Operand, "GREATER_THAN_OR_EQUAL")
	if term.Right == nil {
		t.Fatal("Terms[0].Right is nil")
	}
	assertStringPtr(t, "Terms[0].Right.Value", term.Right.Value, "90")
}

func TestScreensSource(t *testing.T) {
	t.Parallel()

	screens := screensFixture(t)
	s := screens[0]

	if s.Source == nil {
		t.Fatal("Source is nil")
	}
	if s.Source.ID == nil {
		t.Fatal("Source.ID is nil")
	}
	if got, want := *s.Source.ID, "src-001"; got != want {
		t.Errorf("Source.ID = %q, want %q", got, want)
	}
	if s.Source.Type == nil {
		t.Fatal("Source.Type is nil")
	}
	if got, want := *s.Source.Type, "USER"; got != want {
		t.Errorf("Source.Type = %q, want %q", got, want)
	}
	if s.Source.Pub == nil {
		t.Fatal("Source.Pub is nil")
	}
	if got, want := *s.Source.Pub, "msr"; got != want {
		t.Errorf("Source.Pub = %q, want %q", got, want)
	}
}

func TestScreensNullSource(t *testing.T) {
	t.Parallel()

	screens := screensFixture(t)
	s := screens[1]

	if s.Source != nil {
		t.Errorf("Source = %v, want nil", s.Source)
	}
	if s.Description != nil {
		t.Errorf("Description = %v, want nil", s.Description)
	}
}

func TestScreensRequestBody(t *testing.T) {
	t.Parallel()

	wantBytes, err := os.ReadFile("testdata/Screens/request.json")
	if err != nil {
		t.Fatalf("os.ReadFile(request.json) error = %v", err)
	}

	var captured []byte
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captured, _ = readBodyUpToLimit(r.Body, 1<<20)
		respBytes, _ := os.ReadFile("testdata/Screens/response.json")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(respBytes)
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL)
	req := NewScreensRequest()
	_, err = client.Screens(context.Background(), req)
	if err != nil {
		t.Fatalf("Screens() error = %v", err)
	}

	var want, got any
	if err = json.Unmarshal(wantBytes, &want); err != nil {
		t.Fatalf("unmarshal want: %v", err)
	}
	if err = json.Unmarshal(captured, &got); err != nil {
		t.Fatalf("unmarshal got: %v", err)
	}
	wantJSON, _ := json.Marshal(want)
	gotJSON, _ := json.Marshal(got)
	if string(wantJSON) != string(gotJSON) {
		t.Errorf("request body mismatch:\ngot:  %s\nwant: %s", gotJSON, wantJSON)
	}
}

func TestScreensStatusError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL)
	req := NewScreensRequest()
	_, err := client.Screens(context.Background(), req)
	assertErrorType(t, err, &StatusError{})
}

func TestScreensDecodeError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`not json`))
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL)
	req := NewScreensRequest()
	_, err := client.Screens(context.Background(), req)
	assertErrorType(t, err, &DecodeError{})
}

func TestScreensGraphQLError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"errors":[{"message":"bad request"}]}`))
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL)
	req := NewScreensRequest()
	_, err := client.Screens(context.Background(), req)
	assertErrorType(t, err, &GraphQLError{})
}

func TestScreensBodyLimitError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"user":{"screens":` + strings.Repeat("x", 1024) + `}}}`))
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL, WithResponseBodyLimit(64))
	req := NewScreensRequest()
	_, err := client.Screens(context.Background(), req)
	assertErrorType(t, err, &BodyLimitError{})
}

func TestScreensNoAuth(t *testing.T) {
	t.Parallel()

	client, err := NewClient()
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	req := NewScreensRequest()
	_, err = client.Screens(context.Background(), req)
	if err == nil {
		t.Fatal("expected error for missing auth, got nil")
	}
}

func TestScreensFixtureSanitized(t *testing.T) {
	t.Parallel()

	for _, name := range []string{
		"testdata/Screens/request.json",
		"testdata/Screens/response.json",
	} {
		data, err := os.ReadFile(name)
		if err != nil {
			t.Fatalf("os.ReadFile(%s) error = %v", name, err)
		}
		content := strings.ToLower(string(data))
		for _, secret := range []string{"cookie", "session_id", "password"} {
			if strings.Contains(content, secret) {
				t.Errorf("fixture %s contains potential secret %q", name, secret)
			}
		}
	}
}
