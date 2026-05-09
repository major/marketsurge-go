// Package browserauth discovers MarketSurge sessions from local browser cookie
// stores.
//
// This package is optional. Importing it adds browser-store and SQLite
// dependencies that the root marketsurge package intentionally avoids.
// Callers provide an explicit Firefox profile path, and browserauth reads only
// that profile's cookies.sqlite database.
package browserauth
