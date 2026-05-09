package marketsurge

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// screenFixture calls Screen against a test server that serves the response
// fixture and returns the ScreenUser.
func screenFixture(t *testing.T) *ScreenUser {
	t.Helper()

	respBytes, err := os.ReadFile("testdata/Screen/response.json")
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
	req := NewScreenRequest("screen-Peter Lynch")
	coachScreen := true
	req.CoachScreen = &coachScreen

	resp, err := client.Screen(context.Background(), req)
	if err != nil {
		t.Fatalf("Screen() error = %v", err)
	}
	if resp.User == nil || resp.User.Screen == nil {
		t.Fatal("Screen response User or Screen is nil")
	}
	return resp.User
}

func TestScreenHappyPath(t *testing.T) {
	t.Parallel()

	user := screenFixture(t)

	if user.Screen == nil {
		t.Fatal("Screen is nil")
	}
	assertStringPtr(t, "Screen.ID", user.Screen.ID, "screen-Peter Lynch")
	assertStringPtr(t, "Screen.Name", user.Screen.Name, "Peter Lynch")
}

func TestScreenDetailFields(t *testing.T) {
	t.Parallel()

	screen := screenFixture(t).Screen

	assertStringPtr(t, "Site", screen.Site, "marketsurge")
	assertStringPtr(t, "Description", screen.Description, "Fundamental screen inspired by Peter Lynch")
	assertStringPtr(t, "Type", screen.Type, "STOCK_SCREEN")
	assertStringPtr(t, "CreatedAt", screen.CreatedAt, "2024-01-01T00:00:00Z")
	assertStringPtr(t, "UpdatedAt", screen.UpdatedAt, "2024-09-15T14:30:00Z")
}

func TestScreenFilterCriteria(t *testing.T) {
	t.Parallel()

	screen := screenFixture(t).Screen

	if screen.FilterCriteria == nil {
		t.Fatal("FilterCriteria is nil")
	}
	assertStringPtr(t, "FilterCriteria.Type", screen.FilterCriteria.Type, "AND")
	if got, want := len(screen.FilterCriteria.Terms), 2; got != want {
		t.Fatalf("len(FilterCriteria.Terms) = %d, want %d", got, want)
	}

	term := screen.FilterCriteria.Terms[0]
	if term.Left == nil {
		t.Fatal("Terms[0].Left is nil")
	}
	assertStringPtr(t, "Terms[0].Left.Name", term.Left.Name, "RSRating")
	assertStringPtr(t, "Terms[0].Operand", term.Operand, "GREATER_THAN_OR_EQUAL")
	if term.Right == nil {
		t.Fatal("Terms[0].Right is nil")
	}
	assertStringPtr(t, "Terms[0].Right.Value", term.Right.Value, "80")
}

func TestScreenResultConfig(t *testing.T) {
	t.Parallel()

	screen := screenFixture(t).Screen

	if screen.ResultConfig == nil {
		t.Fatal("ResultConfig is nil")
	}
	if screen.ResultConfig.Limit == nil {
		t.Fatal("ResultConfig.Limit is nil")
	}
	if got, want := *screen.ResultConfig.Limit, 500; got != want {
		t.Errorf("ResultConfig.Limit = %d, want %d", got, want)
	}
	if screen.ResultConfig.SortBy == nil {
		t.Fatal("ResultConfig.SortBy is nil")
	}
	assertStringPtr(t, "SortBy.Field", screen.ResultConfig.SortBy.Field, "RSRating")
	assertStringPtr(t, "SortBy.Direction", screen.ResultConfig.SortBy.Direction, "DESCENDING")
}

func TestScreenResultSummary(t *testing.T) {
	t.Parallel()

	screen := screenFixture(t).Screen

	if screen.Result == nil {
		t.Fatal("Result is nil")
	}
	if screen.Result.Count == nil {
		t.Fatal("Result.Count is nil")
	}
	if got, want := *screen.Result.Count, 42; got != want {
		t.Errorf("Result.Count = %d, want %d", got, want)
	}
	assertStringPtr(t, "Result.Description", screen.Result.Description, "Stocks matching Peter Lynch criteria")
	assertStringPtr(t, "Result.UpdatedAt", screen.Result.UpdatedAt, "2024-09-15T14:30:00Z")
}

func TestScreenSource(t *testing.T) {
	t.Parallel()

	screen := screenFixture(t).Screen

	if screen.Source == nil {
		t.Fatal("Source is nil")
	}
	if screen.Source.ExcludeMsrDatabase == nil {
		t.Fatal("Source.ExcludeMsrDatabase is nil")
	}
	if got, want := *screen.Source.ExcludeMsrDatabase, false; got != want {
		t.Errorf("Source.ExcludeMsrDatabase = %v, want %v", got, want)
	}
}

func TestScreenRequestBody(t *testing.T) {
	t.Parallel()

	var gotReq struct {
		OperationName string `json:"operationName"`
		Variables     struct {
			Site        string `json:"site"`
			ScreenID    string `json:"screenId"`
			CoachScreen *bool  `json:"coachScreen"`
		} `json:"variables"`
		Query string `json:"query"`
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&gotReq); err != nil {
			t.Errorf("decode request body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"user":{"screen":null}}}`))
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL)
	req := NewScreenRequest("screen-Peter Lynch")
	coachScreen := true
	req.CoachScreen = &coachScreen

	_, err := client.Screen(context.Background(), req)
	if err != nil {
		t.Fatalf("Screen() error = %v", err)
	}

	if got, want := gotReq.OperationName, "Screen"; got != want {
		t.Errorf("operationName = %q, want %q", got, want)
	}
	if got, want := gotReq.Variables.Site, DefaultScreenSite; got != want {
		t.Errorf("variables.site = %q, want %q", got, want)
	}
	if got, want := gotReq.Variables.ScreenID, "screen-Peter Lynch"; got != want {
		t.Errorf("variables.screenId = %q, want %q", got, want)
	}
	if gotReq.Variables.CoachScreen == nil || !*gotReq.Variables.CoachScreen {
		t.Errorf("variables.coachScreen = %v, want true", gotReq.Variables.CoachScreen)
	}
}

func TestScreenRequestFixture(t *testing.T) {
	t.Parallel()

	wantBytes, err := os.ReadFile("testdata/Screen/request.json")
	if err != nil {
		t.Fatalf("os.ReadFile(request.json) error = %v", err)
	}

	var gotReqBytes []byte
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotReqBytes, err = io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("read request body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"user":{"screen":null}}}`))
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL)
	req := NewScreenRequest("screen-Peter Lynch")
	coachScreen := true
	req.CoachScreen = &coachScreen

	_, err = client.Screen(context.Background(), req)
	if err != nil {
		t.Fatalf("Screen() error = %v", err)
	}

	var want, got any
	if unmarshalErr := json.Unmarshal(wantBytes, &want); unmarshalErr != nil {
		t.Fatalf("unmarshal request fixture: %v", unmarshalErr)
	}
	if unmarshalErr := json.Unmarshal(gotReqBytes, &got); unmarshalErr != nil {
		t.Fatalf("unmarshal sent request: %v", unmarshalErr)
	}

	wantJSON, marshalErr := json.Marshal(want)
	if marshalErr != nil {
		t.Fatalf("marshal want: %v", marshalErr)
	}
	gotJSON, marshalErr := json.Marshal(got)
	if marshalErr != nil {
		t.Fatalf("marshal got: %v", marshalErr)
	}
	if string(gotJSON) != string(wantJSON) {
		t.Errorf("request body mismatch:\ngot:  %s\nwant: %s", gotJSON, wantJSON)
	}
}

func TestScreenStatusError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte("service unavailable"))
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL)
	req := NewScreenRequest("screen-123")
	_, err := client.Screen(context.Background(), req)

	var statusErr *StatusError
	if !errors.As(err, &statusErr) {
		t.Fatalf("Screen() error = %T, want *StatusError", err)
	}
	if got, want := statusErr.StatusCode, http.StatusServiceUnavailable; got != want {
		t.Errorf("StatusCode = %d, want %d", got, want)
	}
}

func TestScreenDecodeError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("{invalid json"))
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL)
	req := NewScreenRequest("screen-123")
	_, err := client.Screen(context.Background(), req)

	var decodeErr *DecodeError
	if !errors.As(err, &decodeErr) {
		t.Fatalf("Screen() error = %T, want *DecodeError", err)
	}
}

func TestScreenGraphQLError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"errors":[{"message":"not authorized","extensions":{"code":"UNAUTHORIZED"}}]}`))
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL)
	req := NewScreenRequest("screen-123")
	_, err := client.Screen(context.Background(), req)

	var gqlErr *GraphQLError
	if !errors.As(err, &gqlErr) {
		t.Fatalf("Screen() error = %T, want *GraphQLError", err)
	}
	if got, want := gqlErr.Error(), "not authorized"; got != want {
		t.Errorf("GraphQLError.Error() = %q, want %q", got, want)
	}
}

func TestScreenBodyLimit(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"user":{"screen":null}}}`))
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL, WithResponseBodyLimit(5))
	req := NewScreenRequest("screen-123")
	_, err := client.Screen(context.Background(), req)

	if !IsBodyLimit(err) {
		t.Fatalf("Screen() error = %v, want BodyLimitError", err)
	}
}

func TestScreenNoAuth(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		t.Error("Screen() unexpectedly sent request without auth")
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(WithGraphQLURL(srv.URL))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	req := NewScreenRequest("screen-123")
	_, err = client.Screen(context.Background(), req)
	if err == nil {
		t.Fatal("Screen() error = nil, want auth error")
	}
	if got, want := err.Error(), "no JWT provider configured"; !strings.Contains(got, want) {
		t.Errorf("error = %q, want substring %q", got, want)
	}
}

func TestNewScreenRequestDefaults(t *testing.T) {
	t.Parallel()

	req := NewScreenRequest("screen-abc-123")
	if got, want := req.Site, DefaultScreenSite; got != want {
		t.Errorf("Site = %q, want %q", got, want)
	}
	if got, want := req.ScreenID, "screen-abc-123"; got != want {
		t.Errorf("ScreenID = %q, want %q", got, want)
	}
	if req.CoachScreen != nil {
		t.Errorf("CoachScreen = %v, want nil", req.CoachScreen)
	}
}
