package marketsurge

import "context"

// queryFetchChartMarkups is the GraphQL query for the FetchChartMarkups endpoint.
const queryFetchChartMarkups = `query FetchChartMarkups($site: Site!, $dowJonesKey: String, $frequency: ChartMarkupFrequencyInput, $dateStart: String, $dateEnd: String, $cursorId: String, $limit: Int, $sortDir: SortDirInput) {
  user {
    chartMarkups(
      site: $site
      dowJonesKey: $dowJonesKey
      frequency: $frequency
      dateStart: $dateStart
      dateEnd: $dateEnd
      cursorId: $cursorId
      limit: $limit
      sortDir: $sortDir
    ) {
      cursorId
      chartMarkups {
        createdAt
        data
        frequency
        id
        name
        site
        updatedAt
      }
    }
  }
}`

// DefaultFetchChartMarkupsSite is the default site for FetchChartMarkups queries.
const DefaultFetchChartMarkupsSite = "marketsurge"

// FetchChartMarkupsRequest holds parameters for the FetchChartMarkups query.
type FetchChartMarkupsRequest struct {
	Site        string
	DowJonesKey string
	Frequency   string
	SortDir     string
	DateStart   *string
	DateEnd     *string
	CursorID    *string
	Limit       *int
}

// NewFetchChartMarkupsRequest creates a FetchChartMarkupsRequest with sensible
// defaults for the given Dow Jones key.
func NewFetchChartMarkupsRequest(dowJonesKey string) FetchChartMarkupsRequest {
	return FetchChartMarkupsRequest{
		Site:        DefaultFetchChartMarkupsSite,
		DowJonesKey: dowJonesKey,
	}
}

// fetchChartMarkupsVariables holds the GraphQL variables sent with the query.
type fetchChartMarkupsVariables struct {
	Site        string  `json:"site"`
	DowJonesKey string  `json:"dowJonesKey"`
	Frequency   string  `json:"frequency,omitempty"`
	SortDir     string  `json:"sortDir,omitempty"`
	DateStart   *string `json:"dateStart,omitempty"`
	DateEnd     *string `json:"dateEnd,omitempty"`
	CursorID    *string `json:"cursorId,omitempty"`
	Limit       *int    `json:"limit,omitempty"`
}

// ---------------------------------------------------------------------------
// Response types
// ---------------------------------------------------------------------------

// FetchChartMarkupsResponse is the top-level response from the FetchChartMarkups query.
type FetchChartMarkupsResponse struct {
	User *FetchChartMarkupsUser `json:"user"`
}

// FetchChartMarkupsUser holds the user-scoped data.
type FetchChartMarkupsUser struct {
	ChartMarkups *FetchChartMarkupsList `json:"chartMarkups"`
}

// FetchChartMarkupsList holds a paginated list of chart markups.
type FetchChartMarkupsList struct {
	CursorID     *string            `json:"cursorId"`
	ChartMarkups []FetchChartMarkup `json:"chartMarkups"`
}

// FetchChartMarkup represents a single user-saved chart markup.
type FetchChartMarkup struct {
	CreatedAt *string `json:"createdAt"`
	Data      *string `json:"data"`
	Frequency *string `json:"frequency"`
	ID        *string `json:"id"`
	Name      *string `json:"name"`
	Site      *string `json:"site"`
	UpdatedAt *string `json:"updatedAt"`
}

// ---------------------------------------------------------------------------
// Client method
// ---------------------------------------------------------------------------

// FetchChartMarkups fetches user-saved chart markups for the requested parameters.
func (c *Client) FetchChartMarkups(
	ctx context.Context,
	req FetchChartMarkupsRequest,
) (*FetchChartMarkupsResponse, error) {
	vars := fetchChartMarkupsVariables(req)

	var resp FetchChartMarkupsResponse
	if err := c.doGraphQL(ctx, "FetchChartMarkups", vars, queryFetchChartMarkups, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
