package browserauth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/browserutils/kooky"
	"github.com/browserutils/kooky/browser/firefox"

	marketsurge "github.com/major/marketsurge-go"
)

const investorsDomain = "investors.com"

var (
	// ErrProfileNotFound reports that the requested Firefox profile directory does not exist.
	ErrProfileNotFound = errors.New("firefox profile not found")

	// ErrNoCookies reports that no investors.com cookies were found in the Firefox profile.
	ErrNoCookies = errors.New("investors.com cookies not found")

	// ErrCookieDBLocked reports that Firefox is holding a lock on cookies.sqlite.
	ErrCookieDBLocked = errors.New("firefox cookie database locked")
)

//nolint:gochecknoglobals // Test seam for substituting Firefox cookie reads without touching browser stores.
var readCookies = firefox.ReadCookies

// SessionFromFirefox reads investors.com cookies from a Firefox profile and
// returns a MarketSurge session built from those cookies.
//
// The profilePath argument must point to a Firefox profile directory or a
// cookies.sqlite database. The function does not discover or modify profiles.
func SessionFromFirefox(profilePath string) (*marketsurge.Session, error) {
	cookieDBPath, err := resolveDBPath(profilePath)
	if err != nil {
		return nil, err
	}

	kookyCookies, err := readCookies(
		context.Background(),
		cookieDBPath,
		kooky.Valid,
		kooky.DomainHasSuffix(investorsDomain),
	)
	if err != nil {
		if isCookieDBLocked(err) {
			return nil, fmt.Errorf("%w: read %s: %w", ErrCookieDBLocked, cookieDBPath, err)
		}
		return nil, fmt.Errorf("read firefox cookies from %s: %w", cookieDBPath, err)
	}

	httpCookies := httpCookiesFromKooky(kookyCookies)
	if len(httpCookies) == 0 {
		return nil, fmt.Errorf("%w in %s", ErrNoCookies, cookieDBPath)
	}

	session := marketsurge.NewSession(httpCookies)
	return &session, nil
}

func resolveDBPath(profilePath string) (string, error) {
	info, err := os.Stat(profilePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("%w: %s", ErrProfileNotFound, profilePath)
		}
		return "", fmt.Errorf("stat firefox profile path %s: %w", profilePath, err)
	}
	if info.IsDir() {
		return filepath.Join(profilePath, "cookies.sqlite"), nil
	}
	return profilePath, nil
}

// httpCookiesFromKooky converts kooky's cookie wrappers into standard library
// cookies for marketsurge.NewSession.
func httpCookiesFromKooky(cookies []*kooky.Cookie) []*http.Cookie {
	httpCookies := make([]*http.Cookie, 0, len(cookies))
	for _, cookie := range cookies {
		if cookie == nil {
			continue
		}

		// #nosec G124 - cookies are read from the user's Firefox profile and sent back to MarketSurge only.
		httpCookies = append(httpCookies, &cookie.Cookie)
	}
	return httpCookies
}

// isCookieDBLocked reports whether err looks like a SQLite lock failure from
// Firefox keeping cookies.sqlite open.
func isCookieDBLocked(err error) bool {
	message := strings.ToLower(err.Error())
	return strings.Contains(message, "database is locked") ||
		strings.Contains(message, "database table is locked") ||
		strings.Contains(message, "sqlite_busy")
}
