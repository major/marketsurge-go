package marketsurge

import (
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNewSessionDefensivelyCopiesCookies(t *testing.T) {
	original := []*http.Cookie{{Name: "session", Value: "before", Unparsed: []string{"a"}}}

	session := NewSession(original)
	original[0].Value = "after"
	original[0].Unparsed[0] = "b"
	original[0] = &http.Cookie{Name: "extra", Value: "cookie"}

	want := []*http.Cookie{{Name: "session", Value: "before", Unparsed: []string{"a"}}}
	if diff := cmp.Diff(want, session.Cookies); diff != "" {
		t.Fatalf("NewSession(cookies) mismatch (-want +got):\n%s", diff)
	}
}

func TestNewSessionNilCookiesReturnsEmptySlice(t *testing.T) {
	session := NewSession(nil)

	if session.Cookies == nil {
		t.Fatal("NewSession(nil) returned nil cookie slice, want empty non-nil slice")
	}
	if got, want := len(session.Cookies), 0; got != want {
		t.Fatalf("NewSession(nil) len(Cookies) = %d, want %d", got, want)
	}
}
