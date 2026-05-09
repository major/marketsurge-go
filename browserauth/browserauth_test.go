package browserauth

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	"github.com/browserutils/kooky"

	marketsurge "github.com/major/marketsurge-go"
)

const (
	fixtureSessionCookie = "IBD_SESSION"
	fixtureSessionValue  = "fixture-session"
	fixtureUserCookie    = "IBD_USER"
	fixtureUserValue     = "fixture-user"
)

// TestSessionFromFirefoxReadsFixtureDatabase verifies that Firefox profile and
// cookie database paths both produce sessions from investors.com cookies.
func TestSessionFromFirefoxReadsFixtureDatabase(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		profilePath string
	}{
		{
			name:        "profile directory",
			profilePath: "testdata",
		},
		{
			name:        "sqlite database path",
			profilePath: filepath.Join("testdata", "cookies.sqlite"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			session, err := SessionFromFirefox(tt.profilePath)
			if err != nil {
				t.Fatalf("SessionFromFirefox(%q) error = %v, want nil", tt.profilePath, err)
			}
			if session == nil {
				t.Fatalf("SessionFromFirefox(%q) session = nil, want session", tt.profilePath)
			}

			assertSessionCookieValues(t, session, map[string]string{
				fixtureSessionCookie: fixtureSessionValue,
				fixtureUserCookie:    fixtureUserValue,
			})
		})
	}
}

// TestSessionFromFirefoxReportsMissingPath verifies missing profile and cookie
// database paths return a matchable not-found error.
func TestSessionFromFirefoxReportsMissingPath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		profilePath string
	}{
		{
			name:        "missing directory",
			profilePath: filepath.Join("testdata", "missing-profile"),
		},
		{
			name:        "missing sqlite database",
			profilePath: filepath.Join("testdata", "missing.sqlite"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			session, err := SessionFromFirefox(tt.profilePath)
			if !errors.Is(err, ErrProfileNotFound) {
				t.Fatalf(
					"SessionFromFirefox(%q) error = %v, want ErrProfileNotFound",
					tt.profilePath,
					err,
				)
			}
			if session != nil {
				t.Fatalf("SessionFromFirefox(%q) session = %#v, want nil", tt.profilePath, session)
			}
		})
	}
}

// TestSessionFromFirefoxReportsNoInvestorsCookies verifies readable databases
// without investors.com cookies return a matchable no-cookie error.
func TestSessionFromFirefoxReportsNoInvestorsCookies(t *testing.T) {
	t.Parallel()

	profilePath := filepath.Join("testdata", "no_investors.sqlite")
	session, err := SessionFromFirefox(profilePath)
	if !errors.Is(err, ErrNoCookies) {
		t.Fatalf("SessionFromFirefox(%q) error = %v, want ErrNoCookies", profilePath, err)
	}
	if session != nil {
		t.Fatalf("SessionFromFirefox(%q) session = %#v, want nil", profilePath, session)
	}
}

// TestSessionFromFirefoxReportsLockedDatabase verifies locked cookie database
// errors are surfaced as ErrCookieDBLocked.
func TestSessionFromFirefoxReportsLockedDatabase(t *testing.T) {
	original := readCookies
	readCookies = func(context.Context, string, ...kooky.Filter) ([]*kooky.Cookie, error) {
		return nil, errors.New("database is locked")
	}
	t.Cleanup(func() { readCookies = original })

	profilePath := filepath.Join("testdata", "cookies.sqlite")
	session, err := SessionFromFirefox(profilePath)
	if !errors.Is(err, ErrCookieDBLocked) {
		t.Fatalf("SessionFromFirefox(%q) error = %v, want ErrCookieDBLocked", profilePath, err)
	}
	if session != nil {
		t.Fatalf("SessionFromFirefox(%q) session = %#v, want nil", profilePath, session)
	}
}

// TestSessionFromFirefoxWrapsReadErrors verifies non-lock read errors are
// wrapped for callers.
func TestSessionFromFirefoxWrapsReadErrors(t *testing.T) {
	original := readCookies
	wantErr := errors.New("boom")
	readCookies = func(context.Context, string, ...kooky.Filter) ([]*kooky.Cookie, error) {
		return nil, wantErr
	}
	t.Cleanup(func() { readCookies = original })

	profilePath := filepath.Join("testdata", "cookies.sqlite")
	session, err := SessionFromFirefox(profilePath)
	if !errors.Is(err, wantErr) {
		t.Fatalf("SessionFromFirefox(%q) error = %v, want wrapped %v", profilePath, err, wantErr)
	}
	if session != nil {
		t.Fatalf("SessionFromFirefox(%q) session = %#v, want nil", profilePath, session)
	}
}

// TestIsCookieDBLockedMatchesSQLiteBusyErrors verifies locked SQLite errors
// are classified before they are wrapped for callers.
func TestIsCookieDBLockedMatchesSQLiteBusyErrors(t *testing.T) {
	t.Parallel()

	if !isCookieDBLocked(errors.New("sqlite: database is locked")) {
		t.Fatal("isCookieDBLocked(database locked) = false, want true")
	}
}

// TestHttpCookiesFromKookySkipsNilEntries verifies nil cookie wrappers are
// ignored while valid cookies are preserved.
func TestHttpCookiesFromKookySkipsNilEntries(t *testing.T) {
	t.Parallel()

	cookies := []*kooky.Cookie{
		nil,
		func() *kooky.Cookie {
			cookie := &kooky.Cookie{}
			cookie.Name = "test"
			cookie.Value = "val"
			return cookie
		}(),
		nil,
	}

	got := httpCookiesFromKooky(cookies)
	if len(got) != 1 {
		t.Fatalf("httpCookiesFromKooky(%#v) len = %d, want 1", cookies, len(got))
	}
	if got[0] == nil || got[0].Name != "test" || got[0].Value != "val" {
		t.Fatalf("httpCookiesFromKooky(%#v)[0] = %#v, want cookie test=val", cookies, got[0])
	}
}

// assertSessionCookieValues checks the returned session has exactly the wanted
// cookie names and values.
func assertSessionCookieValues(t *testing.T, session *marketsurge.Session, want map[string]string) {
	t.Helper()

	got := sessionCookieValues(session)
	if len(got) != len(want) {
		t.Fatalf("sessionCookieValues() count = %d, want %d; values = %#v", len(got), len(want), got)
	}
	for name, wantValue := range want {
		if gotValue := got[name]; gotValue != wantValue {
			t.Errorf("sessionCookieValues()[%q] = %q, want %q", name, gotValue, wantValue)
		}
	}
}

// sessionCookieValues converts a session's cookies into a map for focused test
// assertions.
func sessionCookieValues(session *marketsurge.Session) map[string]string {
	if session == nil {
		return nil
	}

	values := make(map[string]string, len(session.Cookies))
	for _, cookie := range session.Cookies {
		if cookie == nil {
			continue
		}
		values[cookie.Name] = cookie.Value
	}
	return values
}
