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

// navTreeFixture calls NavTree against a test server that serves the
// response fixture and returns the NavTreeUser.
func navTreeFixture(t *testing.T) *NavTreeUser {
	t.Helper()

	respBytes, err := os.ReadFile("testdata/NavTree/response.json")
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
	req := NewNavTreeRequest()
	resp, err := client.NavTree(context.Background(), req)
	if err != nil {
		t.Fatalf("NavTree() error = %v", err)
	}
	if resp.User == nil {
		t.Fatal("NavTree().User is nil")
	}
	return resp.User
}

func TestNavTreeHappyPath(t *testing.T) {
	t.Parallel()

	user := navTreeFixture(t)

	if got, want := len(user.NavTree), 3; got != want {
		t.Fatalf("len(NavTree) = %d, want %d", got, want)
	}
}

func TestNavTreeFolderFields(t *testing.T) {
	t.Parallel()

	user := navTreeFixture(t)

	folder := user.NavTree[0]
	assertStringPtr(t, "NavTree[0].ID", folder.ID, "folder-reports")
	assertStringPtr(t, "NavTree[0].Name", folder.Name, "My Reports")
	assertStringPtr(t, "NavTree[0].Type", folder.Type, "SYSTEM_FOLDER")
	assertStringPtr(t, "NavTree[0].TreeType", folder.TreeType, "MSR_NAV")
	assertStringPtr(t, "NavTree[0].ContentType", folder.ContentType, "REPORTS")
	if folder.ParentID != nil {
		t.Errorf("NavTree[0].ParentID = %v, want nil", *folder.ParentID)
	}
	if folder.URL != nil {
		t.Errorf("NavTree[0].URL = %v, want nil", *folder.URL)
	}
	if folder.ReferenceID != nil {
		t.Errorf("NavTree[0].ReferenceID = %v, want nil", *folder.ReferenceID)
	}
}

func TestNavTreeLeafFields(t *testing.T) {
	t.Parallel()

	user := navTreeFixture(t)

	leaf := user.NavTree[1]
	assertStringPtr(t, "NavTree[1].ID", leaf.ID, "report-120")
	assertStringPtr(t, "NavTree[1].Name", leaf.Name, "Minervini Trend - 5 Months")
	assertStringPtr(t, "NavTree[1].ParentID", leaf.ParentID, "folder-reports")
	assertStringPtr(t, "NavTree[1].Type", leaf.Type, "REPORTS_SCREEN")
	assertStringPtr(t, "NavTree[1].TreeType", leaf.TreeType, "MSR_NAV")
	assertStringPtr(t, "NavTree[1].URL", leaf.URL, "/report/minervini-5m")
	assertStringPtr(t, "NavTree[1].ReferenceID", leaf.ReferenceID, `{"originalId":120,"isCoachAccount":false}`)
}

func TestNavTreeChildFields(t *testing.T) {
	t.Parallel()

	user := navTreeFixture(t)

	folder := user.NavTree[0]
	if got, want := len(folder.Children), 2; got != want {
		t.Fatalf("len(Children) = %d, want %d", got, want)
	}
	child := folder.Children[0]
	assertStringPtr(t, "Children[0].ID", child.ID, "report-120")
	assertStringPtr(t, "Children[0].Name", child.Name, "Minervini Trend - 5 Months")
	assertStringPtr(t, "Children[0].Type", child.Type, "REPORTS_SCREEN")
}

func TestNavTreeRequestBody(t *testing.T) {
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
		_, _ = w.Write([]byte(`{"data":{"user":{"navTree":[]}}}`))
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL)
	req := NewNavTreeRequest()
	_, err := client.NavTree(context.Background(), req)
	if err != nil {
		t.Fatalf("NavTree() error = %v", err)
	}

	if got, want := gotReq.OperationName, "NavTree"; got != want {
		t.Errorf("operationName = %q, want %q", got, want)
	}
	if got, want := gotReq.Variables.Site, DefaultNavTreeSite; got != want {
		t.Errorf("variables.site = %q, want %q", got, want)
	}
	if got, want := gotReq.Variables.TreeType, DefaultNavTreeTreeType; got != want {
		t.Errorf("variables.treeType = %q, want %q", got, want)
	}
}

func TestNavTreeRequestFixture(t *testing.T) {
	t.Parallel()

	wantBytes, err := os.ReadFile("testdata/NavTree/request.json")
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
		_, _ = w.Write([]byte(`{"data":{"user":{"navTree":[]}}}`))
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL)
	req := NewNavTreeRequest()
	_, err = client.NavTree(context.Background(), req)
	if err != nil {
		t.Fatalf("NavTree() error = %v", err)
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

func TestNavTreeStatusError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte("service unavailable"))
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL)
	req := NewNavTreeRequest()
	_, err := client.NavTree(context.Background(), req)

	var statusErr *StatusError
	if !errors.As(err, &statusErr) {
		t.Fatalf("NavTree() error = %T, want *StatusError", err)
	}
	if got, want := statusErr.StatusCode, http.StatusServiceUnavailable; got != want {
		t.Errorf("StatusCode = %d, want %d", got, want)
	}
}

func TestNavTreeDecodeError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("{invalid json"))
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL)
	req := NewNavTreeRequest()
	_, err := client.NavTree(context.Background(), req)

	var decodeErr *DecodeError
	if !errors.As(err, &decodeErr) {
		t.Fatalf("NavTree() error = %T, want *DecodeError", err)
	}
}

func TestNavTreeGraphQLError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"errors":[{"message":"not authorized","extensions":{"code":"UNAUTHORIZED"}}]}`))
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL)
	req := NewNavTreeRequest()
	_, err := client.NavTree(context.Background(), req)

	var gqlErr *GraphQLError
	if !errors.As(err, &gqlErr) {
		t.Fatalf("NavTree() error = %T, want *GraphQLError", err)
	}
	if got, want := gqlErr.Error(), "not authorized"; got != want {
		t.Errorf("GraphQLError.Error() = %q, want %q", got, want)
	}
}

func TestNavTreeBodyLimit(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"user":{"navTree":[]}}}`))
	}))
	t.Cleanup(srv.Close)

	client := newGraphQLTestClient(t, srv.URL, WithResponseBodyLimit(5))
	req := NewNavTreeRequest()
	_, err := client.NavTree(context.Background(), req)

	if !IsBodyLimit(err) {
		t.Fatalf("NavTree() error = %v, want BodyLimitError", err)
	}
}

func TestNavTreeNoAuth(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		t.Error("NavTree() unexpectedly sent request without auth")
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(WithGraphQLURL(srv.URL))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	req := NewNavTreeRequest()
	_, err = client.NavTree(context.Background(), req)
	if err == nil {
		t.Fatal("NavTree() error = nil, want auth error")
	}
	if got, want := err.Error(), "no JWT provider configured"; !strings.Contains(got, want) {
		t.Errorf("error = %q, want substring %q", got, want)
	}
}

func TestNewNavTreeRequestDefaults(t *testing.T) {
	t.Parallel()

	req := NewNavTreeRequest()
	if got, want := req.Site, DefaultNavTreeSite; got != want {
		t.Errorf("Site = %q, want %q", got, want)
	}
	if got, want := req.TreeType, DefaultNavTreeTreeType; got != want {
		t.Errorf("TreeType = %q, want %q", got, want)
	}
}
