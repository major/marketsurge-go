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

// fetchChartMarkupsFixture calls FetchChartMarkups against a test server that
// serves the response fixture and returns the parsed response.
func fetchChartMarkupsFixture(t *testing.T) FetchChartMarkupsResponse {
	t.Helper()

	respBytes, err := os.ReadFile("testdata/FetchChartMarkups/response.json")
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
	req := NewFetchChartMarkupsRequest("13-5320")
	resp, err := client.FetchChartMarkups(context.Background(), req)
	if err != nil {
		t.Fatalf("FetchChartMarkups() error = %v", err)
	}
	return *resp
}

func TestFetchChartMarkups(t *testing.T) {
	t.Parallel()

	resp := fetchChartMarkupsFixture(t)

	if resp.User == nil {
		t.Fatal("FetchChartMarkups().User is nil")
	}
	if resp.User.ChartMarkups == nil {
		t.Fatal("FetchChartMarkups().User.ChartMarkups is nil")
	}
	if got, want := len(resp.User.ChartMarkups.ChartMarkups), 2; got != want {
		t.Fatalf("len(ChartMarkups) = %d, want %d", got, want)
	}
}

func TestFetchChartMarkupsUser(t *testing.T) {
	t.Parallel()

	resp := fetchChartMarkupsFixture(t)

	if resp.User == nil || resp.User.ChartMarkups == nil {
		t.Fatal("User or ChartMarkups is nil")
	}
	if resp.User.ChartMarkups.CursorID == nil || *resp.User.ChartMarkups.CursorID != "cursor-abc123" {
		t.Errorf("CursorID = %v, want cursor-abc123", resp.User.ChartMarkups.CursorID)
	}
}

func TestFetchChartMarkupsItems(t *testing.T) {
	t.Parallel()

	resp := fetchChartMarkupsFixture(t)

	items := resp.User.ChartMarkups.ChartMarkups
	if len(items) < 2 {
		t.Fatalf("len(ChartMarkups) = %d, want >= 2", len(items))
	}

	first := items[0]
	if first.ID == nil || *first.ID != "markup-001" {
		t.Errorf("ChartMarkups[0].ID = %v, want markup-001", first.ID)
	}
	if first.Name == nil || *first.Name != "My Trend Lines" {
		t.Errorf("ChartMarkups[0].Name = %v, want My Trend Lines", first.Name)
	}
	if first.Frequency == nil || *first.Frequency != "DAILY" {
		t.Errorf("ChartMarkups[0].Frequency = %v, want DAILY", first.Frequency)
	}
	if first.Site == nil || *first.Site != "marketsurge" {
		t.Errorf("ChartMarkups[0].Site = %v, want marketsurge", first.Site)
	}
	if first.CreatedAt == nil || *first.CreatedAt != "2024-01-15T12:00:00.000Z" {
		t.Errorf("ChartMarkups[0].CreatedAt = %v, want 2024-01-15T12:00:00.000Z", first.CreatedAt)
	}
	if first.UpdatedAt == nil || *first.UpdatedAt != "2024-01-15T12:00:00.000Z" {
		t.Errorf("ChartMarkups[0].UpdatedAt = %v, want 2024-01-15T12:00:00.000Z", first.UpdatedAt)
	}
	if first.Data == nil || *first.Data != `{"version":1,"annotations":[]}` {
		t.Errorf("ChartMarkups[0].Data = %v, want {\"version\":1,\"annotations\":[]}", first.Data)
	}

	second := items[1]
	if second.ID == nil || *second.ID != "markup-002" {
		t.Errorf("ChartMarkups[1].ID = %v, want markup-002", second.ID)
	}
	if second.Frequency == nil || *second.Frequency != "WEEKLY" {
		t.Errorf("ChartMarkups[1].Frequency = %v, want WEEKLY", second.Frequency)
	}
}

func TestNewFetchChartMarkupsRequest(t *testing.T) {
	t.Parallel()

	req := NewFetchChartMarkupsRequest("13-5320")

	if got, want := req.Site, DefaultFetchChartMarkupsSite; got != want {
		t.Errorf("NewFetchChartMarkupsRequest().Site = %q, want %q", got, want)
	}
	if got, want := req.DowJonesKey, "13-5320"; got != want {
		t.Errorf("NewFetchChartMarkupsRequest().DowJonesKey = %q, want %q", got, want)
	}
	if got, want := req.Frequency, ""; got != want {
		t.Errorf("NewFetchChartMarkupsRequest().Frequency = %q, want %q", got, want)
	}
	if got, want := req.SortDir, ""; got != want {
		t.Errorf("NewFetchChartMarkupsRequest().SortDir = %q, want %q", got, want)
	}
	if req.DateStart != nil {
		t.Errorf("NewFetchChartMarkupsRequest().DateStart = %v, want nil", req.DateStart)
	}
	if req.DateEnd != nil {
		t.Errorf("NewFetchChartMarkupsRequest().DateEnd = %v, want nil", req.DateEnd)
	}
	if req.CursorID != nil {
		t.Errorf("NewFetchChartMarkupsRequest().CursorID = %v, want nil", req.CursorID)
	}
	if req.Limit != nil {
		t.Errorf("NewFetchChartMarkupsRequest().Limit = %v, want nil", req.Limit)
	}
}

func TestFetchChartMarkupsNullFields(t *testing.T) {
	t.Parallel()

	list := fetchChartMarkupsFixture(t)
	items := list.User.ChartMarkups.ChartMarkups
	if len(items) < 2 {
		t.Fatalf("len(ChartMarkups) = %d, want >= 2", len(items))
	}

	second := items[1]
	if second.Name != nil {
		t.Errorf("ChartMarkups[1].Name = %q, want nil", *second.Name)
	}
	if second.UpdatedAt != nil {
		t.Errorf("ChartMarkups[1].UpdatedAt = %q, want nil", *second.UpdatedAt)
	}
}

func TestFetchChartMarkupsRequestBodyOptionals(t *testing.T) {
	t.Parallel()

	var gotReq struct {
		Variables struct {
			Frequency string  `json:"frequency"`
			SortDir   string  `json:"sortDir"`
			DateStart *string `json:"dateStart"`
			DateEnd   *string `json:"dateEnd"`
			CursorID  *string `json:"cursorId"`
			Limit     *int    `json:"limit"`
		} `json:"variables"`
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&gotReq); err != nil {
			t.Errorf("decode request body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"user":{"chartMarkups":{"cursorId":null,"chartMarkups":[]}}}}`))
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL)
	start := "2024-01-01"
	end := "2024-12-31"
	cursor := "cur-xyz"
	limit := 50
	req := FetchChartMarkupsRequest{
		Site:        DefaultFetchChartMarkupsSite,
		DowJonesKey: "13-5320",
		Frequency:   "DAILY",
		SortDir:     "DESC",
		DateStart:   &start,
		DateEnd:     &end,
		CursorID:    &cursor,
		Limit:       &limit,
	}
	_, err := client.FetchChartMarkups(context.Background(), req)
	if err != nil {
		t.Fatalf("FetchChartMarkups() error = %v", err)
	}

	if got, want := gotReq.Variables.Frequency, "DAILY"; got != want {
		t.Errorf("variables.frequency = %q, want %q", got, want)
	}
	if got, want := gotReq.Variables.SortDir, "DESC"; got != want {
		t.Errorf("variables.sortDir = %q, want %q", got, want)
	}
	if gotReq.Variables.DateStart == nil || *gotReq.Variables.DateStart != "2024-01-01" {
		t.Errorf("variables.dateStart = %v, want %q", gotReq.Variables.DateStart, "2024-01-01")
	}
	if gotReq.Variables.DateEnd == nil || *gotReq.Variables.DateEnd != "2024-12-31" {
		t.Errorf("variables.dateEnd = %v, want %q", gotReq.Variables.DateEnd, "2024-12-31")
	}
	if gotReq.Variables.CursorID == nil || *gotReq.Variables.CursorID != "cur-xyz" {
		t.Errorf("variables.cursorId = %v, want %q", gotReq.Variables.CursorID, "cur-xyz")
	}
	if gotReq.Variables.Limit == nil || *gotReq.Variables.Limit != 50 {
		t.Errorf("variables.limit = %v, want %d", gotReq.Variables.Limit, 50)
	}
}

func TestFetchChartMarkupsStatusError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte("service unavailable"))
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL)
	req := NewFetchChartMarkupsRequest("13-5320")
	_, err := client.FetchChartMarkups(context.Background(), req)

	var statusErr *StatusError
	if !errors.As(err, &statusErr) {
		t.Fatalf("FetchChartMarkups() error = %T, want *StatusError", err)
	}
	if got, want := statusErr.StatusCode, http.StatusServiceUnavailable; got != want {
		t.Errorf("StatusCode = %d, want %d", got, want)
	}
}

func TestFetchChartMarkupsDecodeError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("{invalid json"))
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL)
	req := NewFetchChartMarkupsRequest("13-5320")
	_, err := client.FetchChartMarkups(context.Background(), req)

	var decodeErr *DecodeError
	if !errors.As(err, &decodeErr) {
		t.Fatalf("FetchChartMarkups() error = %T, want *DecodeError", err)
	}
}

func TestFetchChartMarkupsGraphQLError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"errors":[{"message":"not authorized","extensions":{"code":"UNAUTHORIZED"}}]}`))
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL)
	req := NewFetchChartMarkupsRequest("13-5320")
	_, err := client.FetchChartMarkups(context.Background(), req)

	var gqlErr *GraphQLError
	if !errors.As(err, &gqlErr) {
		t.Fatalf("FetchChartMarkups() error = %T, want *GraphQLError", err)
	}
	if got, want := gqlErr.Error(), "not authorized"; got != want {
		t.Errorf("GraphQLError.Error() = %q, want %q", got, want)
	}
}

func TestFetchChartMarkupsBodyLimit(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"user":{"chartMarkups":{"cursorId":null,"chartMarkups":[]}}}}`))
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL, WithResponseBodyLimit(5))
	req := NewFetchChartMarkupsRequest("13-5320")
	_, err := client.FetchChartMarkups(context.Background(), req)

	if !IsBodyLimit(err) {
		t.Fatalf("FetchChartMarkups() error = %v, want BodyLimitError", err)
	}
}

func TestFetchChartMarkupsNoAuth(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		t.Error("FetchChartMarkups() unexpectedly sent request without auth")
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(WithGraphQLURL(srv.URL))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	req := NewFetchChartMarkupsRequest("13-5320")
	_, err = client.FetchChartMarkups(context.Background(), req)
	if err == nil {
		t.Fatal("FetchChartMarkups() error = nil, want auth error")
	}
	if got, want := err.Error(), "no JWT provider configured"; !strings.Contains(got, want) {
		t.Errorf("error = %q, want substring %q", got, want)
	}
}

func TestFetchChartMarkupsRequestBody(t *testing.T) {
	t.Parallel()

	var gotReq struct {
		OperationName string `json:"operationName"`
		Variables     struct {
			Site        string `json:"site"`
			DowJonesKey string `json:"dowJonesKey"`
		} `json:"variables"`
		Query string `json:"query"`
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&gotReq); err != nil {
			t.Errorf("decode request body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"user":{"chartMarkups":{"cursorId":null,"chartMarkups":[]}}}}`))
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL)
	req := NewFetchChartMarkupsRequest("13-5320")
	_, err := client.FetchChartMarkups(context.Background(), req)
	if err != nil {
		t.Fatalf("FetchChartMarkups() error = %v", err)
	}

	if got, want := gotReq.OperationName, "FetchChartMarkups"; got != want {
		t.Errorf("operationName = %q, want %q", got, want)
	}
	if got, want := gotReq.Variables.Site, DefaultFetchChartMarkupsSite; got != want {
		t.Errorf("variables.site = %q, want %q", got, want)
	}
	if got, want := gotReq.Variables.DowJonesKey, "13-5320"; got != want {
		t.Errorf("variables.dowJonesKey = %q, want %q", got, want)
	}
}

func TestFetchChartMarkupsFixtureSanitized(t *testing.T) {
	t.Parallel()

	for _, name := range []string{
		"testdata/FetchChartMarkups/request.json",
		"testdata/FetchChartMarkups/response.json",
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
