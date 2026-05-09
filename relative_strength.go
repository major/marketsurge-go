package marketsurge

import "context"

// queryRSRatingRIPanel is the GraphQL query for the RSRatingRIPanel endpoint.
const queryRSRatingRIPanel = `query RSRatingRIPanel(
  $symbols: [String!]!
  $symbolDialectType: MDSymbolDialectType!
) {
  marketData(symbols: $symbols, symbolDialectType: $symbolDialectType) {
    id
    originRequest {
      fromDialect
      symbol
    }
    ratings {
      rsRating {
        letterValue
        period
        periodOffset
        value
      }
    }
    pricingStatistics {
      intradayStatistics {
        rsLineNewHigh
      }
    }
  }
}`

// DefaultRSRatingRIPanelSymbolDialectType is the default symbol dialect for
// RSRatingRIPanel queries.
const DefaultRSRatingRIPanelSymbolDialectType = "CHARTING"

// RSRatingRIPanelRequest holds parameters for the RSRatingRIPanel query.
type RSRatingRIPanelRequest struct {
	Symbols           []string
	SymbolDialectType string
}

// NewRSRatingRIPanelRequest creates an RSRatingRIPanelRequest with sensible
// defaults for the given symbols.
func NewRSRatingRIPanelRequest(symbols ...string) RSRatingRIPanelRequest {
	return RSRatingRIPanelRequest{
		Symbols:           symbols,
		SymbolDialectType: DefaultRSRatingRIPanelSymbolDialectType,
	}
}

// rsRatingRIPanelVariables holds the GraphQL variables sent with the query.
type rsRatingRIPanelVariables struct {
	Symbols           []string `json:"symbols"`
	SymbolDialectType string   `json:"symbolDialectType"`
}

// ---------------------------------------------------------------------------
// Response types
// ---------------------------------------------------------------------------

// RSRatingRIPanelResponse is the top-level response from the RSRatingRIPanel query.
type RSRatingRIPanelResponse struct {
	MarketData []RSRatingRIPanelItem `json:"marketData"`
}

// RSRatingRIPanelItem represents a single symbol's RS rating data.
type RSRatingRIPanelItem struct {
	ID                *string                    `json:"id"`
	OriginRequest     *RSRatingOriginRequest     `json:"originRequest"`
	Ratings           *RSRatingRatings           `json:"ratings"`
	PricingStatistics *RSRatingPricingStatistics `json:"pricingStatistics"`
}

// RSRatingOriginRequest holds the original request dialect and symbol.
type RSRatingOriginRequest struct {
	FromDialect *string `json:"fromDialect"`
	Symbol      *string `json:"symbol"`
}

// RSRatingRatings holds the RS rating data.
type RSRatingRatings struct {
	RSRating []RSRatingSnapshot `json:"rsRating"`
}

// RSRatingSnapshot holds a single RS rating value at a specific period.
type RSRatingSnapshot struct {
	LetterValue  *string `json:"letterValue"`
	Period       *string `json:"period"`
	PeriodOffset *string `json:"periodOffset"`
	Value        *int    `json:"value"`
}

// RSRatingPricingStatistics holds pricing statistics data.
type RSRatingPricingStatistics struct {
	IntradayStatistics *RSRatingIntradayStatistics `json:"intradayStatistics"`
}

// RSRatingIntradayStatistics holds intraday statistics data.
type RSRatingIntradayStatistics struct {
	RSLineNewHigh *bool `json:"rsLineNewHigh"`
}

// ---------------------------------------------------------------------------
// Client method
// ---------------------------------------------------------------------------

// RSRatingRIPanel fetches RS rating data for the requested symbols.
func (c *Client) RSRatingRIPanel(ctx context.Context, req RSRatingRIPanelRequest) (*RSRatingRIPanelResponse, error) {
	vars := rsRatingRIPanelVariables(req)

	var resp RSRatingRIPanelResponse
	if err := c.doGraphQL(ctx, "RSRatingRIPanel", vars, queryRSRatingRIPanel, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
