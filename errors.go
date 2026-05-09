package marketsurge

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// StatusError reports a non-2xx HTTP response.
type StatusError struct {
	StatusCode int         `json:"statusCode"`
	Status     string      `json:"status"`
	Body       []byte      `json:"body"`
	Header     http.Header `json:"header"`
}

// Error returns a human-readable HTTP status failure.
func (e *StatusError) Error() string {
	if e == nil {
		return "marketsurge API error <nil>"
	}
	status := e.Status
	if status == "" {
		status = http.StatusText(e.StatusCode)
	}
	if status == "" {
		return fmt.Sprintf("marketsurge API error %d", e.StatusCode)
	}
	return fmt.Sprintf("marketsurge API error %d %s", e.StatusCode, status)
}

// Unwrap returns nil, StatusError does not wrap another error.
func (e *StatusError) Unwrap() error { return nil }

// DecodeError reports a failed decode or transform operation.
type DecodeError struct {
	Operation string `json:"operation"`
	Err       error  `json:"err"`
}

// Error returns a human-readable decode failure.
func (e *DecodeError) Error() string {
	if e == nil {
		return "decode error: <nil>"
	}
	return fmt.Sprintf("%s: %v", e.Operation, e.Err)
}

// Unwrap returns the underlying decode error.
func (e *DecodeError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

// BodyLimitError reports that a body exceeded the configured limit.
type BodyLimitError struct {
	Limit int64 `json:"limit"`
}

// Error returns a human-readable body-size failure.
func (e *BodyLimitError) Error() string {
	if e == nil {
		return "marketsurge response body exceeded <nil> byte limit"
	}
	return fmt.Sprintf("marketsurge response body exceeded %d byte limit", e.Limit)
}

// GraphQLFieldErrorExtensions holds typed GraphQL error extensions.
type GraphQLFieldErrorExtensions struct {
	Code string `json:"code"`
}

// GraphQLFieldError represents a single GraphQL error entry.
type GraphQLFieldError struct {
	Message    string                       `json:"message"`
	Path       []string                     `json:"path"`
	Extensions *GraphQLFieldErrorExtensions `json:"extensions"`
}

// Error returns the GraphQL error message.
func (e GraphQLFieldError) Error() string { return e.Message }

// GraphQLError wraps a set of GraphQL field errors.
type GraphQLError struct {
	Errors []GraphQLFieldError `json:"errors"`
}

// Error returns the first GraphQL error message.
func (e *GraphQLError) Error() string {
	if e == nil || len(e.Errors) == 0 {
		return ""
	}
	return e.Errors[0].Message
}

// Unwrap returns all field errors as error values.
func (e *GraphQLError) Unwrap() []error {
	if e == nil || len(e.Errors) == 0 {
		return nil
	}
	out := make([]error, 0, len(e.Errors))
	for _, fieldErr := range e.Errors {
		out = append(out, fieldErr)
	}
	return out
}

// StatusCode returns the HTTP status code from err when it contains a StatusError.
func StatusCode(err error) int {
	var statusErr *StatusError
	if !errors.As(err, &statusErr) || statusErr == nil {
		return 0
	}
	return statusErr.StatusCode
}

// IsStatusCode reports whether err contains a StatusError with code.
func IsStatusCode(err error, code int) bool {
	return StatusCode(err) == code
}

// IsAuthError reports whether err indicates unauthorized or forbidden access.
func IsAuthError(err error) bool {
	return IsStatusCode(err, http.StatusUnauthorized) || IsStatusCode(err, http.StatusForbidden)
}

// IsRateLimited reports whether err indicates HTTP 429.
func IsRateLimited(err error) bool {
	return IsStatusCode(err, http.StatusTooManyRequests)
}

// IsBodyLimit reports whether err contains a BodyLimitError.
func IsBodyLimit(err error) bool {
	var limitErr *BodyLimitError
	return errors.As(err, &limitErr) && limitErr != nil
}

// RetryAfter returns the parsed Retry-After header from a StatusError.
func RetryAfter(err error) time.Duration {
	var statusErr *StatusError
	if !errors.As(err, &statusErr) || statusErr == nil || statusErr.Header == nil {
		return 0
	}

	raw := strings.TrimSpace(statusErr.Header.Get("Retry-After"))
	if raw == "" {
		return 0
	}

	if secs, parseErr := strconv.ParseInt(raw, 10, 64); parseErr == nil && secs >= 0 {
		return time.Duration(secs) * time.Second
	}

	retryAt, parseErr := http.ParseTime(raw)
	if parseErr != nil {
		return 0
	}
	if d := time.Until(retryAt); d > 0 {
		return d
	}
	return 0
}
