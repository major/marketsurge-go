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
// CoachTree fixtures and tests
// ---------------------------------------------------------------------------

// coachTreeFixture calls CoachTree against a test server that serves the
// response fixture and returns the CoachTreeUser.
func coachTreeFixture(t *testing.T) *CoachTreeUser {
	t.Helper()

	respBytes, err := os.ReadFile("testdata/CoachTree/response.json")
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
	req := NewCoachTreeRequest()
	resp, err := client.CoachTree(context.Background(), req)
	if err != nil {
		t.Fatalf("CoachTree() error = %v", err)
	}
	if resp.User == nil {
		t.Fatal("CoachTree().User is nil")
	}
	return resp.User
}

func TestCoachTreeHappyPath(t *testing.T) {
	t.Parallel()

	user := coachTreeFixture(t)

	if got, want := len(user.Watchlists), 2; got != want {
		t.Fatalf("len(Watchlists) = %d, want %d", got, want)
	}
	if got, want := len(user.Screens), 1; got != want {
		t.Fatalf("len(Screens) = %d, want %d", got, want)
	}
}

func TestCoachTreeNodeFields(t *testing.T) {
	t.Parallel()

	user := coachTreeFixture(t)

	// Folder node (watchlists[0]).
	folder := user.Watchlists[0]
	assertStringPtr(t, "Watchlists[0].ID", folder.ID, "folder-1")
	assertStringPtr(t, "Watchlists[0].Name", folder.Name, "My Watchlists")
	assertStringPtr(t, "Watchlists[0].Type", folder.Type, "FOLDER")
	assertStringPtr(t, "Watchlists[0].TreeType", folder.TreeType, "MSR_NAV")
	assertStringPtr(t, "Watchlists[0].ContentType", folder.ContentType, "WATCHLIST")
	if folder.ParentID != nil {
		t.Errorf("Watchlists[0].ParentID = %v, want nil", *folder.ParentID)
	}

	// Leaf node (watchlists[1]).
	leaf := user.Watchlists[1]
	assertStringPtr(t, "Watchlists[1].ID", leaf.ID, "leaf-2")
	assertStringPtr(t, "Watchlists[1].Name", leaf.Name, "Tech Leaders")
	assertStringPtr(t, "Watchlists[1].Type", leaf.Type, "LEAF")
	assertStringPtr(t, "Watchlists[1].ParentID", leaf.ParentID, "folder-1")
	assertStringPtr(t, "Watchlists[1].TreeType", leaf.TreeType, "MSR_NAV")
	assertStringPtr(t, "Watchlists[1].URL", leaf.URL, "/watchlist/tech-leaders")
	assertStringPtr(t, "Watchlists[1].ReferenceID", leaf.ReferenceID, "12345")
}

func TestCoachTreeChildFields(t *testing.T) {
	t.Parallel()

	user := coachTreeFixture(t)

	folder := user.Watchlists[0]
	if got, want := len(folder.Children), 1; got != want {
		t.Fatalf("len(Children) = %d, want %d", got, want)
	}
	child := folder.Children[0]
	assertStringPtr(t, "Children[0].ID", child.ID, "leaf-1")
	assertStringPtr(t, "Children[0].Name", child.Name, "Growth Stocks")
	assertStringPtr(t, "Children[0].Type", child.Type, "LEAF")
}

func TestCoachTreeRequestBody(t *testing.T) {
	t.Parallel()

	var gotReq struct {
		OperationName string `json:"operationName"`
		Variables     struct {
			Site     string `json:"site"`
			TreeType string `json:"treeType"`
		} `json:"variables"`
		Query string `json:"query"`
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&gotReq); err != nil {
			t.Errorf("decode request body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"user":{"watchlists":[],"screens":[]}}}`))
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL)
	req := NewCoachTreeRequest()
	_, err := client.CoachTree(context.Background(), req)
	if err != nil {
		t.Fatalf("CoachTree() error = %v", err)
	}

	if got, want := gotReq.OperationName, "CoachTree"; got != want {
		t.Errorf("operationName = %q, want %q", got, want)
	}
	if got, want := gotReq.Variables.Site, DefaultCoachTreeSite; got != want {
		t.Errorf("variables.site = %q, want %q", got, want)
	}
	if got, want := gotReq.Variables.TreeType, DefaultCoachTreeTreeType; got != want {
		t.Errorf("variables.treeType = %q, want %q", got, want)
	}
}

func TestCoachTreeStatusError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte("service unavailable"))
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL)
	req := NewCoachTreeRequest()
	_, err := client.CoachTree(context.Background(), req)

	var statusErr *StatusError
	if !errors.As(err, &statusErr) {
		t.Fatalf("CoachTree() error = %T, want *StatusError", err)
	}
	if got, want := statusErr.StatusCode, http.StatusServiceUnavailable; got != want {
		t.Errorf("StatusCode = %d, want %d", got, want)
	}
}

func TestCoachTreeDecodeError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("{invalid json"))
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL)
	req := NewCoachTreeRequest()
	_, err := client.CoachTree(context.Background(), req)

	var decodeErr *DecodeError
	if !errors.As(err, &decodeErr) {
		t.Fatalf("CoachTree() error = %T, want *DecodeError", err)
	}
}

func TestCoachTreeGraphQLError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"errors":[{"message":"not authorized","extensions":{"code":"UNAUTHORIZED"}}]}`))
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL)
	req := NewCoachTreeRequest()
	_, err := client.CoachTree(context.Background(), req)

	var gqlErr *GraphQLError
	if !errors.As(err, &gqlErr) {
		t.Fatalf("CoachTree() error = %T, want *GraphQLError", err)
	}
	if got, want := gqlErr.Error(), "not authorized"; got != want {
		t.Errorf("GraphQLError.Error() = %q, want %q", got, want)
	}
}

func TestCoachTreeBodyLimit(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"user":{"watchlists":[],"screens":[]}}}`))
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL, WithResponseBodyLimit(5))
	req := NewCoachTreeRequest()
	_, err := client.CoachTree(context.Background(), req)

	if !IsBodyLimit(err) {
		t.Fatalf("CoachTree() error = %v, want BodyLimitError", err)
	}
}

func TestCoachTreeNoAuth(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		t.Error("CoachTree() unexpectedly sent request without auth")
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(WithGraphQLURL(srv.URL))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	req := NewCoachTreeRequest()
	_, err = client.CoachTree(context.Background(), req)
	if err == nil {
		t.Fatal("CoachTree() error = nil, want auth error")
	}
	if got, want := err.Error(), "no JWT provider configured"; !strings.Contains(got, want) {
		t.Errorf("error = %q, want substring %q", got, want)
	}
}

// ---------------------------------------------------------------------------
// IndustryGroupRS fixtures and tests
// ---------------------------------------------------------------------------

// industryGroupRSFixture calls IndustryGroupRS against a test server that
// serves the response fixture and returns the first IndustryGroupRSItem.
func industryGroupRSFixture(t *testing.T) IndustryGroupRSItem {
	t.Helper()

	respBytes, err := os.ReadFile("testdata/IndustryGroupRS/response.json")
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
	req := NewIndustryGroupRSRequest("AAPL")
	resp, err := client.IndustryGroupRS(context.Background(), req)
	if err != nil {
		t.Fatalf("IndustryGroupRS() error = %v", err)
	}
	if got, want := len(resp.MarketData), 1; got != want {
		t.Fatalf("len(MarketData) = %d, want %d", got, want)
	}
	return resp.MarketData[0]
}

func TestIndustryGroupRSHappyPath(t *testing.T) {
	t.Parallel()

	item := industryGroupRSFixture(t)

	if item.OriginRequest == nil {
		t.Fatal("OriginRequest is nil")
	}
	assertStringPtr(t, "OriginRequest.Symbol", item.OriginRequest.Symbol, "AAPL")

	if item.Industry == nil {
		t.Fatal("Industry is nil")
	}
	if got, want := len(item.Industry.GroupRS), 1; got != want {
		t.Fatalf("len(GroupRS) = %d, want %d", got, want)
	}
	if item.Industry.GroupRS[0].Value == nil {
		t.Fatal("GroupRS[0].Value is nil")
	}
	if got, want := *item.Industry.GroupRS[0].Value, 85; got != want {
		t.Errorf("GroupRS[0].Value = %d, want %d", got, want)
	}
}

func TestIndustryGroupRSNodeFields(t *testing.T) {
	t.Parallel()

	item := industryGroupRSFixture(t)

	if item.OriginRequest == nil {
		t.Fatal("OriginRequest is nil")
	}
	if item.OriginRequest.Symbol == nil {
		t.Fatal("OriginRequest.Symbol is nil")
	}
	if got, want := *item.OriginRequest.Symbol, "AAPL"; got != want {
		t.Errorf("OriginRequest.Symbol = %q, want %q", got, want)
	}

	if item.Industry == nil {
		t.Fatal("Industry is nil")
	}
	if got := len(item.Industry.GroupRS); got < 1 {
		t.Fatal("GroupRS is empty")
	}
	if item.Industry.GroupRS[0].Value == nil {
		t.Fatal("GroupRS[0].Value is nil")
	}
	if got, want := *item.Industry.GroupRS[0].Value, 85; got != want {
		t.Errorf("GroupRS[0].Value = %d, want %d", got, want)
	}
}

func TestIndustryGroupRSChildFields(t *testing.T) {
	t.Parallel()

	item := industryGroupRSFixture(t)

	if item.Industry == nil {
		t.Fatal("Industry is nil")
	}
	if got, want := len(item.Industry.GroupRS), 1; got != want {
		t.Fatalf("len(GroupRS) = %d, want %d", got, want)
	}
	rs := item.Industry.GroupRS[0]
	if rs.Value == nil {
		t.Fatal("rs.Value is nil")
	}
	if got, want := *rs.Value, 85; got != want {
		t.Errorf("GroupRS[0].Value = %d, want %d", got, want)
	}
}

func TestIndustryGroupRSRequestBody(t *testing.T) {
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
	req := NewIndustryGroupRSRequest("AAPL")
	_, err := client.IndustryGroupRS(context.Background(), req)
	if err != nil {
		t.Fatalf("IndustryGroupRS() error = %v", err)
	}

	if got, want := gotReq.OperationName, "IndustryGroupRS"; got != want {
		t.Errorf("operationName = %q, want %q", got, want)
	}
	if len(gotReq.Variables.Symbols) != 1 || gotReq.Variables.Symbols[0] != "AAPL" {
		t.Errorf("variables.symbols = %v, want [AAPL]", gotReq.Variables.Symbols)
	}
	if got, want := gotReq.Variables.SymbolDialectType, DefaultIndustryGroupRSSymbolDialectType; got != want {
		t.Errorf("variables.symbolDialectType = %q, want %q", got, want)
	}
}

func TestIndustryGroupRSStatusError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte("service unavailable"))
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL)
	req := NewIndustryGroupRSRequest("AAPL")
	_, err := client.IndustryGroupRS(context.Background(), req)

	var statusErr *StatusError
	if !errors.As(err, &statusErr) {
		t.Fatalf("IndustryGroupRS() error = %T, want *StatusError", err)
	}
	if got, want := statusErr.StatusCode, http.StatusServiceUnavailable; got != want {
		t.Errorf("StatusCode = %d, want %d", got, want)
	}
}

func TestIndustryGroupRSDecodeError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("{invalid json"))
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL)
	req := NewIndustryGroupRSRequest("AAPL")
	_, err := client.IndustryGroupRS(context.Background(), req)

	var decodeErr *DecodeError
	if !errors.As(err, &decodeErr) {
		t.Fatalf("IndustryGroupRS() error = %T, want *DecodeError", err)
	}
}

func TestIndustryGroupRSGraphQLError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"errors":[{"message":"not authorized","extensions":{"code":"UNAUTHORIZED"}}]}`))
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL)
	req := NewIndustryGroupRSRequest("AAPL")
	_, err := client.IndustryGroupRS(context.Background(), req)

	var gqlErr *GraphQLError
	if !errors.As(err, &gqlErr) {
		t.Fatalf("IndustryGroupRS() error = %T, want *GraphQLError", err)
	}
	if got, want := gqlErr.Error(), "not authorized"; got != want {
		t.Errorf("GraphQLError.Error() = %q, want %q", got, want)
	}
}

func TestIndustryGroupRSBodyLimit(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"marketData":[]}}`))
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL, WithResponseBodyLimit(5))
	req := NewIndustryGroupRSRequest("AAPL")
	_, err := client.IndustryGroupRS(context.Background(), req)

	if !IsBodyLimit(err) {
		t.Fatalf("IndustryGroupRS() error = %v, want BodyLimitError", err)
	}
}

func TestIndustryGroupRSNoAuth(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		t.Error("IndustryGroupRS() unexpectedly sent request without auth")
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(WithGraphQLURL(srv.URL))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	req := NewIndustryGroupRSRequest("AAPL")
	_, err = client.IndustryGroupRS(context.Background(), req)
	if err == nil {
		t.Fatal("IndustryGroupRS() error = nil, want auth error")
	}
	if got, want := err.Error(), "no JWT provider configured"; !strings.Contains(got, want) {
		t.Errorf("error = %q, want substring %q", got, want)
	}
}

// ---------------------------------------------------------------------------
// Test helpers
// ---------------------------------------------------------------------------

// assertStringPtr checks that a *string pointer is non-nil and equals want.
func assertStringPtr(t *testing.T, field string, got *string, want string) {
	t.Helper()
	if got == nil {
		t.Fatalf("%s is nil, want %q", field, want)
	}
	if *got != want {
		t.Errorf("%s = %q, want %q", field, *got, want)
	}
}
