package marketsurge

import "context"

// queryMarketDataAdhocScreen is the GraphQL query for the MarketDataAdhocScreen endpoint.
const queryMarketDataAdhocScreen = `query MarketDataAdhocScreen(
  $correlationTag: String!
  $adhocQuery: MDAdhocQueryInput
  $responseColumns: [MDAdhocScreenerDataItemInput!]!
  $resultLimit: Int!
  $pageSize: Int!
  $pageSkip: Int
  $includeSource: MDScreenerDataSourceInput!
  $resultType: MDScreenerResultType
) {
  marketDataAdhocScreen(
    correlationTag: $correlationTag
    adhocQuery: $adhocQuery
    resultLimit: $resultLimit
    pageSize: $pageSize
    pageSkip: $pageSkip
    includeSource: $includeSource
    responseDataPoints: $responseColumns
    resultType: $resultType
  ) {
    correlationTag
    elapsedTime
    errorValues
    numberOfInstrumentsInSource
    numberOfMatchingInstruments
    adhocQueryString
    adhocQuery {
      terms {
        numberOfMatchingInstruments
        ordinal
        left {
          name
          mdItemID
        }
        operand
        right {
          value
          maximumValue
          minimumValue
        }
      }
    }
    responseValues {
      value
      mdItem {
        mdItemID
        name
      }
    }
  }
}`

// queryRunScreen is the GraphQL query for the RunScreen endpoint.
const queryRunScreen = `query RunScreen($input: ScreenResultInput!) {
  user {
    runScreen(input: $input) {
      numberOfMatchingInstruments
      responseValues {
        value
        mdItem {
          name
          mdItemID
        }
      }
    }
  }
}`

// queryScreens is the GraphQL query for the Screens endpoint.
const queryScreens = `query Screens($site: Site!, $type: ScreenType, $sortDir: SortDirInput) {
  user {
    screens(site: $site, type: $type, sortDir: $sortDir) {
      site
      id
      name
      type
      source {
        id
        type
        pub
      }
      updatedAt
      filterCriteria
      description
      createdAt
    }
  }
}`

// ---------------------------------------------------------------------------
// Defaults
// ---------------------------------------------------------------------------

const (
	// DefaultAdhocScreenCorrelationTag is the default correlation tag for adhoc screen queries.
	DefaultAdhocScreenCorrelationTag = "marketsurge"
	// DefaultAdhocScreenPageSize is the default page size for adhoc screen queries.
	DefaultAdhocScreenPageSize = 1000
	// DefaultAdhocScreenResultLimit is the default result limit for adhoc screen queries.
	DefaultAdhocScreenResultLimit = 1000000
	// DefaultAdhocScreenResultType is the default result type for adhoc screen queries.
	DefaultAdhocScreenResultType = "RESULT_WITH_EXPRESSION_COUNTS"
	// DefaultRunScreenCorrelationTag is the default correlation tag for run screen queries.
	DefaultRunScreenCorrelationTag = "marketsurge"
	// DefaultRunScreenPageSize is the default page size for run screen queries.
	DefaultRunScreenPageSize = 1000
	// DefaultRunScreenResultLimit is the default result limit for run screen queries.
	DefaultRunScreenResultLimit = 1000000
	// DefaultRunScreenSite is the default site for run screen queries.
	DefaultRunScreenSite = "marketsurge"
	// DefaultScreensSite is the default site for screens queries.
	DefaultScreensSite = "marketsurge"
)

// ---------------------------------------------------------------------------
// MarketDataAdhocScreen request types
// ---------------------------------------------------------------------------

// AdhocScreenSortInformation specifies the sort direction and priority for a
// response column in an adhoc screen query.
type AdhocScreenSortInformation struct {
	Direction string `json:"direction"`
	Order     string `json:"order"`
}

// AdhocScreenResponseColumn identifies a data column to include in the adhoc screen response.
// Set SortInformation to control server-side sorting of the results.
type AdhocScreenResponseColumn struct {
	Name            string                      `json:"name"`
	SortInformation *AdhocScreenSortInformation `json:"sortInformation,omitempty"`
}

// AdhocScreenIncludeSource specifies the data source for an adhoc screen query.
// Use ScreenID to run a predefined report, or Instruments to screen specific symbols.
type AdhocScreenIncludeSource struct {
	ScreenID    *AdhocScreenID          `json:"screenId,omitempty"`
	Instruments *AdhocScreenInstruments `json:"instruments,omitempty"`
}

// AdhocScreenID identifies a predefined screen by ID and dialect.
type AdhocScreenID struct {
	ID      int    `json:"id"`
	Dialect string `json:"dialect"`
}

// AdhocScreenInstruments specifies symbols and dialect for an instrument-based screen.
type AdhocScreenInstruments struct {
	Symbols []string `json:"symbols"`
	Dialect string   `json:"dialect"`
}

// MarketDataAdhocScreenRequest holds parameters for the MarketDataAdhocScreen query.
type MarketDataAdhocScreenRequest struct {
	CorrelationTag  string
	ResponseColumns []AdhocScreenResponseColumn
	AdhocQuery      *string
	IncludeSource   AdhocScreenIncludeSource
	PageSize        int
	ResultLimit     int
	PageSkip        int
	ResultType      string
}

// NewMarketDataAdhocScreenRequest creates a MarketDataAdhocScreenRequest with sensible defaults
// for the given response columns.
func NewMarketDataAdhocScreenRequest(responseColumns []AdhocScreenResponseColumn) MarketDataAdhocScreenRequest {
	return MarketDataAdhocScreenRequest{
		CorrelationTag:  DefaultAdhocScreenCorrelationTag,
		ResponseColumns: responseColumns,
		IncludeSource:   AdhocScreenIncludeSource{},
		PageSize:        DefaultAdhocScreenPageSize,
		ResultLimit:     DefaultAdhocScreenResultLimit,
		PageSkip:        0,
		ResultType:      DefaultAdhocScreenResultType,
	}
}

// adhocScreenVariables holds the GraphQL variables sent with the query.
type adhocScreenVariables struct {
	CorrelationTag  string                      `json:"correlationTag"`
	ResponseColumns []AdhocScreenResponseColumn `json:"responseColumns"`
	AdhocQuery      *string                     `json:"adhocQuery"`
	IncludeSource   AdhocScreenIncludeSource    `json:"includeSource"`
	PageSize        int                         `json:"pageSize"`
	ResultLimit     int                         `json:"resultLimit"`
	PageSkip        int                         `json:"pageSkip"`
	ResultType      string                      `json:"resultType"`
}

// ---------------------------------------------------------------------------
// RunScreen request types
// ---------------------------------------------------------------------------

// RunScreenSortInformation specifies the sort direction and priority for a
// response column in a RunScreen query.
type RunScreenSortInformation struct {
	Direction string `json:"direction"`
	Order     string `json:"order"`
}

// RunScreenResponseColumn identifies a data column to include in the RunScreen response.
// Set SortInformation to control server-side sorting of the results.
type RunScreenResponseColumn struct {
	Name            string                    `json:"name"`
	SortInformation *RunScreenSortInformation `json:"sortInformation,omitempty"`
}

// RunScreenIncludeSource specifies the data source universe for a RunScreen query.
// Set Source to restrict results to a specific source (e.g. "IBD_STOCKS").
// A nil RunScreenIncludeSource pointer on RunScreenInput sends null on the wire.
type RunScreenIncludeSource struct {
	Source *string `json:"source,omitempty"`
}

// RunScreenInput holds the input parameters for a RunScreen query.
type RunScreenInput struct {
	CorrelationTag  string                    `json:"correlationTag"`
	CoachAccount    bool                      `json:"coachAccount"`
	IncludeSource   *RunScreenIncludeSource   `json:"includeSource"`
	PageSize        int                       `json:"pageSize"`
	ResultLimit     int                       `json:"resultLimit"`
	ScreenID        string                    `json:"screenId"`
	Site            string                    `json:"site"`
	Skip            int                       `json:"skip"`
	ResponseColumns []RunScreenResponseColumn `json:"responseColumns"`
}

// RunScreenRequest holds parameters for the RunScreen query.
type RunScreenRequest struct {
	Input RunScreenInput
}

// NewRunScreenRequest creates a RunScreenRequest with sensible defaults
// for the given screen ID and response columns.
func NewRunScreenRequest(screenID string, responseColumns []RunScreenResponseColumn) RunScreenRequest {
	return RunScreenRequest{
		Input: RunScreenInput{
			CorrelationTag:  DefaultRunScreenCorrelationTag,
			CoachAccount:    true,
			PageSize:        DefaultRunScreenPageSize,
			ResultLimit:     DefaultRunScreenResultLimit,
			ScreenID:        screenID,
			Site:            DefaultRunScreenSite,
			Skip:            0,
			ResponseColumns: responseColumns,
		},
	}
}

// runScreenVariables holds the GraphQL variables sent with the query.
type runScreenVariables struct {
	Input RunScreenInput `json:"input"`
}

// ---------------------------------------------------------------------------
// Screens request types
// ---------------------------------------------------------------------------

// ScreensRequest holds parameters for the Screens query.
type ScreensRequest struct {
	Site    string
	Type    *string
	SortDir *string
}

// NewScreensRequest creates a ScreensRequest with sensible defaults.
func NewScreensRequest() ScreensRequest {
	return ScreensRequest{
		Site: DefaultScreensSite,
	}
}

// screensVariables holds the GraphQL variables sent with the query.
type screensVariables struct {
	Site    string  `json:"site"`
	Type    *string `json:"type,omitempty"`
	SortDir *string `json:"sortDir,omitempty"`
}

// ---------------------------------------------------------------------------
// Response types - MarketDataAdhocScreen
// ---------------------------------------------------------------------------

// AdhocScreenResponse is the top-level response from the MarketDataAdhocScreen query.
type AdhocScreenResponse struct {
	MarketDataAdhocScreen *AdhocScreenResult `json:"marketDataAdhocScreen"`
}

// AdhocScreenResult holds the result of an adhoc screen query.
type AdhocScreenResult struct {
	CorrelationTag              *string             `json:"correlationTag"`
	ElapsedTime                 *string             `json:"elapsedTime"`
	ErrorValues                 []string            `json:"errorValues"`
	NumberOfInstrumentsInSource *int                `json:"numberOfInstrumentsInSource"`
	NumberOfMatchingInstruments *int                `json:"numberOfMatchingInstruments"`
	AdhocQueryString            *string             `json:"adhocQueryString"`
	AdhocQuery                  *AdhocQueryResult   `json:"adhocQuery"`
	ResponseValues              [][]AdhocScreenCell `json:"responseValues"`
}

// AdhocQueryResult holds the parsed filter criteria returned by an adhoc screen query.
type AdhocQueryResult struct {
	Terms []AdhocQueryTerm `json:"terms"`
}

// AdhocQueryTerm represents a single filter term in an adhoc screen query.
type AdhocQueryTerm struct {
	NumberOfMatchingInstruments *int                 `json:"numberOfMatchingInstruments"`
	Ordinal                     *int                 `json:"ordinal"`
	Left                        *AdhocQueryTermLeft  `json:"left"`
	Operand                     *string              `json:"operand"`
	Right                       *AdhocQueryTermRight `json:"right"`
}

// AdhocQueryTermLeft identifies the market data item being filtered.
type AdhocQueryTermLeft struct {
	Name     *string `json:"name"`
	MDItemID *string `json:"mdItemID"`
}

// AdhocQueryTermRight holds the comparison values for a filter term.
type AdhocQueryTermRight struct {
	Value        *string `json:"value"`
	MaximumValue *string `json:"maximumValue"`
	MinimumValue *string `json:"minimumValue"`
}

// AdhocScreenCell represents a single cell in an adhoc screen response row.
type AdhocScreenCell struct {
	Value  *string          `json:"value"`
	MDItem *AdhocScreenItem `json:"mdItem"`
}

// AdhocScreenItem identifies a market data item in a screen cell.
type AdhocScreenItem struct {
	MDItemID *string `json:"mdItemID"`
	Name     *string `json:"name"`
}

// ---------------------------------------------------------------------------
// Response types - RunScreen
// ---------------------------------------------------------------------------

// RunScreenResponse is the top-level response from the RunScreen query.
type RunScreenResponse struct {
	User *RunScreenUser `json:"user"`
}

// RunScreenUser wraps the user-scoped RunScreen result.
type RunScreenUser struct {
	RunScreen *RunScreenResult `json:"runScreen"`
}

// RunScreenResult holds the result of running a named screen.
type RunScreenResult struct {
	NumberOfMatchingInstruments *int              `json:"numberOfMatchingInstruments"`
	ResponseValues              [][]RunScreenCell `json:"responseValues"`
}

// RunScreenCell represents a single cell in a run screen response row.
type RunScreenCell struct {
	Value  *string        `json:"value"`
	MDItem *RunScreenItem `json:"mdItem"`
}

// RunScreenItem identifies a market data item in a run screen cell.
type RunScreenItem struct {
	Name     *string `json:"name"`
	MDItemID *string `json:"mdItemID"`
}

// ---------------------------------------------------------------------------
// Response types - Screens
// ---------------------------------------------------------------------------

// ScreensResponse is the top-level response from the Screens query.
type ScreensResponse struct {
	User *ScreensUser `json:"user"`
}

// ScreensUser wraps the user-scoped Screens result.
type ScreensUser struct {
	Screens []ScreenEntry `json:"screens"`
}

// ScreenEntry represents a saved screen definition.
type ScreenEntry struct {
	Site           *string       `json:"site"`
	ID             *string       `json:"id"`
	Name           *string       `json:"name"`
	Type           *string       `json:"type"`
	Source         *ScreenSource `json:"source"`
	UpdatedAt      *string       `json:"updatedAt"`
	FilterCriteria *string       `json:"filterCriteria"`
	Description    *string       `json:"description"`
	CreatedAt      *string       `json:"createdAt"`
}

// ScreenSource represents a data source linked to a saved screen.
type ScreenSource struct {
	ID   *string `json:"id"`
	Type *string `json:"type"`
	Pub  *string `json:"pub"`
}

// ---------------------------------------------------------------------------
// Client methods
// ---------------------------------------------------------------------------

// MarketDataAdhocScreen runs an ad-hoc screen query against the MarketSurge screener.
func (c *Client) MarketDataAdhocScreen(
	ctx context.Context,
	req MarketDataAdhocScreenRequest,
) (*AdhocScreenResponse, error) {
	vars := adhocScreenVariables(req)

	var resp AdhocScreenResponse
	if err := c.doGraphQL(ctx, "MarketDataAdhocScreen", vars, queryMarketDataAdhocScreen, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// RunScreen runs a saved screen by its ID.
func (c *Client) RunScreen(ctx context.Context, req RunScreenRequest) (*RunScreenResponse, error) {
	vars := runScreenVariables(req)

	var resp RunScreenResponse
	if err := c.doGraphQL(ctx, "RunScreen", vars, queryRunScreen, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Screens lists saved screen definitions for the authenticated user.
func (c *Client) Screens(ctx context.Context, req ScreensRequest) (*ScreensResponse, error) {
	vars := screensVariables(req)

	var resp ScreensResponse
	if err := c.doGraphQL(ctx, "Screens", vars, queryScreens, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
