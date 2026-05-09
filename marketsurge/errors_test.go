package marketsurge

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestErrorTypes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		got  error
		want string
	}{
		{
			name: "status error",
			got: &StatusError{
				StatusCode: http.StatusUnauthorized,
				Status:     http.StatusText(http.StatusUnauthorized),
			},
			want: "marketsurge API error 401 Unauthorized",
		},
		{
			name: "decode error",
			got:  &DecodeError{Operation: "decode response", Err: errors.New("boom")},
			want: "decode response: boom",
		},
		{
			name: "body limit error",
			got:  &BodyLimitError{Limit: 8},
			want: "marketsurge response body exceeded 8 byte limit",
		},
		{
			name: "graphql field error",
			got:  GraphQLFieldError{Message: "bad field"},
			want: "bad field",
		},
		{
			name: "graphql error",
			got: &GraphQLError{Errors: []GraphQLFieldError{
				{Message: "first"},
				{Message: "second"},
			}},
			want: "first",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.got.Error(); got != tt.want {
				t.Fatalf("%T.Error() = %q, want %q", tt.got, got, tt.want)
			}
		})
	}
}

func TestErrorUnwrap(t *testing.T) {
	t.Parallel()

	decodeErr := errors.New("boom")
	statusErr := &StatusError{
		StatusCode: http.StatusTooManyRequests,
		Status:     http.StatusText(http.StatusTooManyRequests),
	}
	graphQLErr := &GraphQLError{
		Errors: []GraphQLFieldError{{Message: "first"}, {Message: "second"}},
	}

	tests := []struct {
		name string
		err  error
		want error
	}{
		{name: "status error", err: statusErr, want: nil},
		{
			name: "decode error",
			err:  &DecodeError{Operation: "decode response", Err: decodeErr},
			want: decodeErr,
		},
		{name: "body limit error", err: &BodyLimitError{Limit: 8}, want: nil},
		{name: "graphql field error", err: GraphQLFieldError{Message: "bad field"}, want: nil},
		{name: "graphql error", err: graphQLErr, want: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := errors.Unwrap(tt.err); !errors.Is(got, tt.want) {
				t.Fatalf("errors.Unwrap(%T) = %v, want %v", tt.err, got, tt.want)
			}
		})
	}

	if multi, ok := any(graphQLErr).(interface{ Unwrap() []error }); !ok {
		t.Fatal("GraphQLError does not implement multi-error unwrap")
	} else {
		got := multi.Unwrap()
		if len(got) != 2 {
			t.Fatalf("GraphQLError.Unwrap() len = %d, want 2", len(got))
		}
		if got[0].Error() != "first" || got[1].Error() != "second" {
			t.Fatalf("GraphQLError.Unwrap() = %v, want [first second]", got)
		}
	}
}

func TestJSONTags(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   any
		want string
	}{
		{
			name: "status error",
			in: &StatusError{
				StatusCode: 401,
				Status:     "Unauthorized",
				Body:       []byte("body"),
				Header:     http.Header{"Retry-After": []string{"30"}},
			},
			want: `{"statusCode":401,"status":"Unauthorized","body":"Ym9keQ==","header":{"Retry-After":["30"]}}`,
		},
		{
			name: "decode error",
			in:   &DecodeError{Operation: "decode response", Err: errors.New("boom")},
			want: `{"operation":"decode response","err":{}}`,
		},
		{
			name: "body limit error",
			in:   &BodyLimitError{Limit: 8},
			want: `{"limit":8}`,
		},
		{
			name: "graphql field error",
			in: &GraphQLFieldError{
				Message:    "bad field",
				Path:       []string{"query", "field"},
				Extensions: &GraphQLFieldErrorExtensions{Code: "BAD_FIELD"},
			},
			want: `{"message":"bad field","path":["query","field"],"extensions":{"code":"BAD_FIELD"}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gotBytes, err := json.Marshal(tt.in)
			if err != nil {
				t.Fatalf("json.Marshal(%T) returned error: %v", tt.in, err)
			}
			if got := string(gotBytes); got != tt.want {
				t.Fatalf("json.Marshal(%T) = %s, want %s", tt.in, got, tt.want)
			}
		})
	}
}

func TestStatusCodeHelper(t *testing.T) {
	t.Parallel()

	statusErr := &StatusError{
		StatusCode: http.StatusTooManyRequests,
		Status:     http.StatusText(http.StatusTooManyRequests),
	}

	if got := StatusCode(statusErr); got != http.StatusTooManyRequests {
		t.Fatalf("StatusCode(%v) = %d, want %d", statusErr, got, http.StatusTooManyRequests)
	}
	if got := StatusCode(errors.New("plain error")); got != 0 {
		t.Fatalf("StatusCode(plain error) = %d, want 0", got)
	}
	if !IsStatusCode(statusErr, http.StatusTooManyRequests) {
		t.Fatal("IsStatusCode(statusErr, 429) = false, want true")
	}
	if IsStatusCode(statusErr, http.StatusUnauthorized) {
		t.Fatal("IsStatusCode(statusErr, 401) = true, want false")
	}
	if IsStatusCode(errors.New("plain error"), http.StatusUnauthorized) {
		t.Fatal("IsStatusCode(plain error, 401) = true, want false")
	}
}

func TestIsAuthError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		err  error
		want bool
	}{
		{name: "nil", err: nil, want: false},
		{
			name: "wrapped non-status",
			err:  fmt.Errorf("wrap: %w", errors.New("invalid session")),
			want: false,
		},
		{
			name: "unauthorized status",
			err: &StatusError{
				StatusCode: http.StatusUnauthorized,
				Status:     http.StatusText(http.StatusUnauthorized),
			},
			want: true,
		},
		{
			name: "forbidden status",
			err: &StatusError{
				StatusCode: http.StatusForbidden,
				Status:     http.StatusText(http.StatusForbidden),
			},
			want: true,
		},
		{
			name: "rate limited status",
			err: &StatusError{
				StatusCode: http.StatusTooManyRequests,
				Status:     http.StatusText(http.StatusTooManyRequests),
			},
			want: false,
		},
		{name: "plain error", err: errors.New("plain error"), want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := IsAuthError(tt.err); got != tt.want {
				t.Fatalf("IsAuthError(%v) = %t, want %t", tt.err, got, tt.want)
			}
		})
	}
}

func TestIsRateLimited(t *testing.T) {
	t.Parallel()

	if !IsRateLimited(
		&StatusError{
			StatusCode: http.StatusTooManyRequests,
			Status:     http.StatusText(http.StatusTooManyRequests),
		},
	) {
		t.Fatal("IsRateLimited(429) = false, want true")
	}
	if IsRateLimited(
		&StatusError{
			StatusCode: http.StatusUnauthorized,
			Status:     http.StatusText(http.StatusUnauthorized),
		},
	) {
		t.Fatal("IsRateLimited(401) = true, want false")
	}
	if IsRateLimited(errors.New("plain error")) {
		t.Fatal("IsRateLimited(plain error) = true, want false")
	}
}

func TestIsBodyLimit(t *testing.T) {
	t.Parallel()

	if !IsBodyLimit(&BodyLimitError{Limit: 8}) {
		t.Fatal("IsBodyLimit(body limit error) = false, want true")
	}
	if !IsBodyLimit(fmt.Errorf("wrap: %w", &BodyLimitError{Limit: 8})) {
		t.Fatal("IsBodyLimit(wrapped body limit error) = false, want true")
	}
	if IsBodyLimit(errors.New("plain error")) {
		t.Fatal("IsBodyLimit(plain error) = true, want false")
	}
}

func TestRetryAfter(t *testing.T) {
	t.Parallel()

	future := time.Now().Add(2 * time.Hour).UTC().Truncate(time.Second)
	tests := []struct {
		name      string
		err       error
		wantAfter time.Duration
	}{
		{
			name:      "delta seconds",
			err:       &StatusError{Header: http.Header{"Retry-After": []string{"30"}}},
			wantAfter: 30 * time.Second,
		},
		{
			name: "http date",
			err: &StatusError{
				Header: http.Header{"Retry-After": []string{future.Format(http.TimeFormat)}},
			},
			wantAfter: 2 * time.Hour,
		},
		{
			name: "invalid header",
			err:  &StatusError{Header: http.Header{"Retry-After": []string{"later"}}},
		},
		{name: "missing header", err: &StatusError{}},
		{name: "plain error", err: errors.New("plain error")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := RetryAfter(tt.err)
			if tt.wantAfter == 0 {
				if got != 0 {
					t.Fatalf("RetryAfter(%v) = %v, want 0", tt.err, got)
				}
				return
			}
			if got <= 0 || got > tt.wantAfter {
				t.Fatalf(
					"RetryAfter(%v) = %v, want positive duration up to %v",
					tt.err,
					got,
					tt.wantAfter,
				)
			}
		})
	}
}
