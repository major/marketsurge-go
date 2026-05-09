# AGENTS.md - browserauth package

Optional subpackage that reads MarketSurge sessions from local Firefox cookie stores. Importing this package pulls in kooky and SQLite dependencies that the root `marketsurge` package intentionally avoids.

## Isolation boundary

This package exists to keep browser-store dependencies out of the root `marketsurge` package. The root package's import smoke test (`import_smoke_test.go`) enforces that `marketsurge` never transitively imports kooky, SQLite, keyring, dbus, or browserutils. All browser cookie access must go through this package.

## Public API

- `SessionFromFirefox(profilePath string) (*marketsurge.Session, error)` - reads `cookies.sqlite` from a Firefox profile directory or database file path, returns a session with defensively copied cookies.

## Sentinel errors

- `ErrProfileNotFound` - profile directory does not exist
- `ErrNoCookies` - no investors.com cookies found in the profile
- `ErrCookieDBLocked` - Firefox is holding a lock on cookies.sqlite (close Firefox first)

## Test seam

The package-level `readCookies` variable (declared as `firefox.ReadCookies`) is the test seam for substituting cookie reads without touching real browser stores. Tests override this variable to inject fixture data. The `//nolint:gochecknoglobals` directive on this variable is the only nolint in the codebase.

## Test fixtures

`testdata/` contains a minimal Firefox cookies.sqlite database for testing. Tests use this fixture instead of real browser profiles.

## Security

- Never expose Firefox profile paths, raw cookie values, or JWT tokens in error messages, log output, or test assertions.
- This package is for desktop automation only, not server-side code.
- `SessionFromFirefox` does not discover or modify Firefox profiles. Callers provide an explicit path.
