package marketsurge

import "context"

// queryGetAllWatchlistNames is the GraphQL query for listing all watchlists.
const queryGetAllWatchlistNames = `query GetAllWatchlistNames($pub: String!) {
  watchlists(pub: $pub) {
    id
    name
    lastModifiedDateUtc
    description
  }
}`

// queryFlaggedSymbols is the GraphQL query for fetching a single watchlist with items.
const queryFlaggedSymbols = `query FlaggedSymbols($pub: String!, $watchlistId: ID!) {
  watchlist(pub: $pub, watchlistId: $watchlistId) {
    id
    name
    lastModifiedDateUtc
    description
    items {
      key
      dowJonesKey
    }
  }
}`

// DefaultWatchlistPub is the default publication identifier for watchlist queries.
const DefaultWatchlistPub = "msr"

// ---------------------------------------------------------------------------
// GetAllWatchlistNames
// ---------------------------------------------------------------------------

// GetAllWatchlistNamesRequest holds parameters for the GetAllWatchlistNames query.
type GetAllWatchlistNamesRequest struct {
	Pub string
}

// NewGetAllWatchlistNamesRequest creates a GetAllWatchlistNamesRequest with
// sensible defaults.
func NewGetAllWatchlistNamesRequest() GetAllWatchlistNamesRequest {
	return GetAllWatchlistNamesRequest{Pub: DefaultWatchlistPub}
}

// getAllWatchlistNamesVariables holds the GraphQL variables sent with the query.
type getAllWatchlistNamesVariables struct {
	Pub string `json:"pub"`
}

// GetAllWatchlistNamesResponse is the top-level response from the
// GetAllWatchlistNames query.
type GetAllWatchlistNamesResponse struct {
	Watchlists []WatchlistSummary `json:"watchlists"`
}

// WatchlistSummary represents a single watchlist entry returned by
// GetAllWatchlistNames.
type WatchlistSummary struct {
	ID                  string `json:"id"`
	Name                string `json:"name"`
	LastModifiedDateUtc string `json:"lastModifiedDateUtc"`
	Description         string `json:"description"`
}

// GetAllWatchlistNames fetches all watchlist names for the given publication.
func (c *Client) GetAllWatchlistNames(
	ctx context.Context,
	req GetAllWatchlistNamesRequest,
) (*GetAllWatchlistNamesResponse, error) {
	vars := getAllWatchlistNamesVariables(req)

	var resp GetAllWatchlistNamesResponse
	if err := c.doGraphQL(ctx, "GetAllWatchlistNames", vars, queryGetAllWatchlistNames, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ---------------------------------------------------------------------------
// FlaggedSymbols
// ---------------------------------------------------------------------------

// FlaggedSymbolsRequest holds parameters for the FlaggedSymbols query.
type FlaggedSymbolsRequest struct {
	Pub         string
	WatchlistID string
}

// NewFlaggedSymbolsRequest creates a FlaggedSymbolsRequest with sensible
// defaults for the given watchlist ID.
func NewFlaggedSymbolsRequest(watchlistID string) FlaggedSymbolsRequest {
	return FlaggedSymbolsRequest{Pub: DefaultWatchlistPub, WatchlistID: watchlistID}
}

// flaggedSymbolsVariables holds the GraphQL variables sent with the query.
type flaggedSymbolsVariables struct {
	Pub         string `json:"pub"`
	WatchlistID string `json:"watchlistId"`
}

// FlaggedSymbolsResponse is the top-level response from the FlaggedSymbols query.
type FlaggedSymbolsResponse struct {
	Watchlist WatchlistDetail `json:"watchlist"`
}

// WatchlistDetail represents a watchlist with its items.
type WatchlistDetail struct {
	ID                  string          `json:"id"`
	Name                string          `json:"name"`
	LastModifiedDateUtc string          `json:"lastModifiedDateUtc"`
	Description         string          `json:"description"`
	Items               []WatchlistItem `json:"items"`
}

// WatchlistItem represents a single symbol in a watchlist.
type WatchlistItem struct {
	Key         string `json:"key"`
	DowJonesKey string `json:"dowJonesKey"`
}

// FlaggedSymbols fetches the symbols in the specified watchlist.
func (c *Client) FlaggedSymbols(
	ctx context.Context,
	req FlaggedSymbolsRequest,
) (*FlaggedSymbolsResponse, error) {
	vars := flaggedSymbolsVariables(req)

	var resp FlaggedSymbolsResponse
	if err := c.doGraphQL(ctx, "FlaggedSymbols", vars, queryFlaggedSymbols, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
