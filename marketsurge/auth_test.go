package marketsurge

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestExchangeJWT(t *testing.T) {
	t.Parallel()

	requestSeen := make(chan struct{}, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestSeen <- struct{}{}

		if got, want := r.Method, http.MethodGet; got != want {
			t.Errorf("Client.ExchangeJWT(ctx, session) method = %s, want %s", got, want)
		}
		if got, want := r.URL.Path, "/client"; got != want {
			t.Errorf("Client.ExchangeJWT(ctx, session) path = %s, want %s", got, want)
		}

		assertHeader(t, r.Header, "User-Agent", testUserAgent)
		assertEmptyHeader(t, r.Header, "X-Encrypted-Document-Key")
		assertHeader(t, r.Header, "X-Original-Host", "marketsurge.investors.com")
		assertEmptyHeader(t, r.Header, "X-Original-Referrer")
		assertHeader(t, r.Header, "X-Original-Url", "/mstool")
		assertHeader(t, r.Header, "Referer", RefererURL)
		assertHeader(t, r.Header, "Origin", OriginURL)

		assertRequestCookie(t, r, "ibd-session", "session-value")
		assertRequestCookie(t, r, "auth", "auth-value")

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"isLoggedIn":true,"jwt":"jwt-token","given_name":"Major","family_name":"Hayden"}`))
	}))
	t.Cleanup(server.Close)

	client := newAuthTestClient(t, server.URL)
	session := NewSession([]*http.Cookie{
		{Name: "ibd-session", Value: "session-value"},
		{Name: "auth", Value: "auth-value"},
	})

	got, err := client.ExchangeJWT(context.Background(), session)
	if err != nil {
		t.Fatalf("Client.ExchangeJWT(ctx, session) error = %v", err)
	}

	want := &ClientInfoResponse{
		IsLoggedIn: true,
		JWT:        "jwt-token",
		GivenName:  "Major",
		FamilyName: "Hayden",
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("Client.ExchangeJWT(ctx, session) mismatch (-want +got):\n%s", diff)
	}
	select {
	case <-requestSeen:
	default:
		t.Fatal("Client.ExchangeJWT(ctx, session) did not send request")
	}
}

func TestExchangeJWTErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		session    Session
		statusCode int
		body       string
		wantErr    any
		wantNoSend bool
	}{
		{
			name:       "no cookies",
			session:    NewSession(nil),
			statusCode: http.StatusOK,
			body:       `{"isLoggedIn":true,"jwt":"jwt-token"}`,
			wantErr:    errStringContains("session has no cookies"),
			wantNoSend: true,
		},
		{
			name:       "non-2xx status",
			session:    NewSession([]*http.Cookie{{Name: "ibd-session", Value: "session-value"}}),
			statusCode: http.StatusUnauthorized,
			body:       "unauthorized",
			wantErr:    &StatusError{},
		},
		{
			name:       "not logged in",
			session:    NewSession([]*http.Cookie{{Name: "ibd-session", Value: "session-value"}}),
			statusCode: http.StatusOK,
			body:       `{"isLoggedIn":false,"jwt":"jwt-token"}`,
			wantErr:    errStringContains("not logged in"),
		},
		{
			name:       "empty jwt",
			session:    NewSession([]*http.Cookie{{Name: "ibd-session", Value: "session-value"}}),
			statusCode: http.StatusOK,
			body:       `{"isLoggedIn":true,"jwt":""}`,
			wantErr:    errStringContains("JWT not found"),
		},
		{
			name:       "malformed json",
			session:    NewSession([]*http.Cookie{{Name: "ibd-session", Value: "session-value"}}),
			statusCode: http.StatusOK,
			body:       `{`,
			wantErr:    &DecodeError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var requestCount atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				requestCount.Add(1)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte(tt.body))
			}))
			t.Cleanup(server.Close)

			client := newAuthTestClient(t, server.URL)
			_, err := client.ExchangeJWT(context.Background(), tt.session)
			assertExchangeJWTError(t, err, tt.wantErr)

			if tt.wantNoSend && requestCount.Load() != 0 {
				t.Fatalf("Client.ExchangeJWT(ctx, %s) sent %d requests, want 0", tt.name, requestCount.Load())
			}
		})
	}
}

type errStringContains string

func assertEmptyHeader(t *testing.T, header http.Header, name string) {
	t.Helper()

	values, ok := header[http.CanonicalHeaderKey(name)]
	if !ok {
		t.Errorf("request header %s missing, want empty value", name)
		return
	}
	if len(values) != 1 || values[0] != "" {
		t.Errorf("request header %s = %q, want empty value", name, values)
	}
}

func assertRequestCookie(t *testing.T, r *http.Request, name string, want string) {
	t.Helper()

	cookie, err := r.Cookie(name)
	if err != nil {
		t.Errorf("request cookie %s error = %v, want value %q", name, err, want)
		return
	}
	if got := cookie.Value; got != want {
		t.Errorf("request cookie %s = %q, want %q", name, got, want)
	}
}

func assertExchangeJWTError(t *testing.T, err error, want any) {
	t.Helper()

	if err == nil {
		t.Fatalf("Client.ExchangeJWT(ctx, session) error = nil, want %T", want)
	}
	switch want := want.(type) {
	case *StatusError:
		var target *StatusError
		if !errors.As(err, &target) {
			t.Fatalf("Client.ExchangeJWT(ctx, session) error = %T, want *StatusError", err)
		}
	case *DecodeError:
		var target *DecodeError
		if !errors.As(err, &target) {
			t.Fatalf("Client.ExchangeJWT(ctx, session) error = %T, want *DecodeError", err)
		}
	case errStringContains:
		if !strings.Contains(err.Error(), string(want)) {
			t.Fatalf("Client.ExchangeJWT(ctx, session) error = %q, want substring %q", err, want)
		}
	default:
		t.Fatalf("assertExchangeJWTError(err, want) unsupported want type %T", want)
	}
}

func newAuthTestClient(t *testing.T, investorsBaseURL string) *Client {
	t.Helper()

	client, err := NewClient(WithInvestorsBaseURL(investorsBaseURL), WithUserAgent(testUserAgent))
	if err != nil {
		t.Fatalf("NewClient(auth test options) error = %v", err)
	}
	return client
}
