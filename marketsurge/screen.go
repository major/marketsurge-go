package marketsurge

import "context"

// queryScreen is the GraphQL query for the Screen endpoint, which fetches
// a single screen definition by ID.
const queryScreen = `query Screen($site: Site!, $screenId: ID!, $coachScreen: Boolean) {
  user {
    screen(site: $site, screenId: $screenId, coachScreen: $coachScreen) {
      id
      name
      site
      description
      filterCriteria
      resultConfig {
        limit
        sortBy {
          field
          direction
        }
      }
      result {
        count
        description
        updatedAt
      }
      type
      source {
        excludeMsrDatabase
      }
      createdAt
      updatedAt
    }
  }
}`

// ---------------------------------------------------------------------------
// Defaults
// ---------------------------------------------------------------------------

const (
	// DefaultScreenSite is the default site for Screen queries.
	DefaultScreenSite = "marketsurge"
)

// ---------------------------------------------------------------------------
// Screen request types
// ---------------------------------------------------------------------------

// ScreenRequest holds parameters for the Screen query.
type ScreenRequest struct {
	Site        string
	ScreenID    string
	CoachScreen *bool
}

// NewScreenRequest creates a ScreenRequest with sensible defaults for the
// given screen ID.
func NewScreenRequest(screenID string) ScreenRequest {
	return ScreenRequest{
		Site:     DefaultScreenSite,
		ScreenID: screenID,
	}
}

// screenVariables holds the GraphQL variables sent with the Screen query.
type screenVariables struct {
	Site        string `json:"site"`
	ScreenID    string `json:"screenId"`
	CoachScreen *bool  `json:"coachScreen,omitempty"`
}

// ---------------------------------------------------------------------------
// Screen response types
// ---------------------------------------------------------------------------

// ScreenResponse is the top-level response from the Screen query.
type ScreenResponse struct {
	User *ScreenUser `json:"user"`
}

// ScreenUser wraps the user-scoped Screen result.
type ScreenUser struct {
	Screen *ScreenDetail `json:"screen"`
}

// ScreenDetail holds the full definition of a single screen, including its
// filter criteria, result configuration, and latest result summary.
type ScreenDetail struct {
	ID             *string               `json:"id"`
	Name           *string               `json:"name"`
	Site           *string               `json:"site"`
	Description    *string               `json:"description"`
	FilterCriteria *ScreenFilterCriteria `json:"filterCriteria"`
	ResultConfig   *ScreenResultConfig   `json:"resultConfig"`
	Result         *ScreenResultSummary  `json:"result"`
	Type           *string               `json:"type"`
	Source         *ScreenDetailSource   `json:"source"`
	CreatedAt      *string               `json:"createdAt"`
	UpdatedAt      *string               `json:"updatedAt"`
}

// ScreenFilterCriteria describes the filter rules applied to a screen.
type ScreenFilterCriteria struct {
	Terms []ScreenFilterTerm `json:"terms"`
	Type  *string            `json:"type"`
}

// ScreenFilterTerm represents a single filter condition within the criteria.
type ScreenFilterTerm struct {
	Left    *ScreenFilterTermLeft  `json:"left"`
	Operand *string                `json:"operand"`
	Right   *ScreenFilterTermRight `json:"right"`
}

// ScreenFilterTermLeft identifies the data field on the left side of a
// filter condition.
type ScreenFilterTermLeft struct {
	Name *string `json:"name"`
}

// ScreenFilterTermRight holds the comparison value(s) on the right side of
// a filter condition.
type ScreenFilterTermRight struct {
	Value *string `json:"value"`
}

// ScreenResultConfig holds the result configuration for a screen, including
// the result limit and sort order.
type ScreenResultConfig struct {
	Limit  *int          `json:"limit"`
	SortBy *ScreenSortBy `json:"sortBy"`
}

// ScreenSortBy specifies the sort field and direction within a screen's
// result configuration.
type ScreenSortBy struct {
	Field     *string `json:"field"`
	Direction *string `json:"direction"`
}

// ScreenResultSummary holds summary information about the latest screen run,
// including the number of matching instruments.
type ScreenResultSummary struct {
	Count       *int    `json:"count"`
	Description *string `json:"description"`
	UpdatedAt   *string `json:"updatedAt"`
}

// ScreenDetailSource holds source configuration for a single screen.
type ScreenDetailSource struct {
	ExcludeMsrDatabase *bool `json:"excludeMsrDatabase"`
}

// ---------------------------------------------------------------------------
// Client methods
// ---------------------------------------------------------------------------

// Screen fetches a single screen definition by its ID.
func (c *Client) Screen(ctx context.Context, req ScreenRequest) (*ScreenResponse, error) {
	vars := screenVariables(req)

	var resp ScreenResponse
	if err := c.doGraphQL(ctx, "Screen", vars, queryScreen, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
