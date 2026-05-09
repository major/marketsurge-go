package marketsurge

import "context"

// queryOwnership is the GraphQL query for the Ownership endpoint.
const queryOwnership = `query Ownership($symbols: [String!]!, $symbolDialectType: MDSymbolDialectType!) {
  marketData(symbols: $symbols, symbolDialectType: $symbolDialectType) {
    ownership {
      fundsFloatPercentHeld {
        formattedValue
      }
      fundOwnershipSummary {
        date {
          value
        }
        numberOfFundsHeld {
          formattedValue
        }
      }
    }
  }
}`

// DefaultOwnershipSymbolDialectType is the default symbol dialect for Ownership queries.
const DefaultOwnershipSymbolDialectType = "CHARTING"

// OwnershipRequest holds parameters for the Ownership query.
type OwnershipRequest struct {
	Symbols           []string
	SymbolDialectType string
}

// NewOwnershipRequest creates an OwnershipRequest with sensible defaults
// for the given symbols.
func NewOwnershipRequest(symbols ...string) OwnershipRequest {
	return OwnershipRequest{
		Symbols:           symbols,
		SymbolDialectType: DefaultOwnershipSymbolDialectType,
	}
}

// ownershipVariables holds the GraphQL variables sent with the query.
type ownershipVariables struct {
	Symbols           []string `json:"symbols"`
	SymbolDialectType string   `json:"symbolDialectType"`
}

// ---------------------------------------------------------------------------
// Response types
// ---------------------------------------------------------------------------

// OwnershipResponse is the top-level response from the Ownership query.
type OwnershipResponse struct {
	MarketData []OwnershipItem `json:"marketData"`
}

// OwnershipItem represents a single symbol's ownership data.
type OwnershipItem struct {
	Ownership *OwnershipData `json:"ownership"`
}

// OwnershipData holds fund ownership statistics.
type OwnershipData struct {
	FundsFloatPercentHeld *OwnershipFormattedValue    `json:"fundsFloatPercentHeld"`
	FundOwnershipSummary  []OwnershipQuarterlySummary `json:"fundOwnershipSummary"`
}

// OwnershipFormattedValue holds a formatted string value.
type OwnershipFormattedValue struct {
	FormattedValue *string `json:"formattedValue"`
}

// OwnershipQuarterlySummary holds fund ownership data for a single quarter.
type OwnershipQuarterlySummary struct {
	Date              *OwnershipDateValue      `json:"date"`
	NumberOfFundsHeld *OwnershipFormattedValue `json:"numberOfFundsHeld"`
}

// OwnershipDateValue holds a single date string value.
type OwnershipDateValue struct {
	Value string `json:"value"`
}

// ---------------------------------------------------------------------------
// Client method
// ---------------------------------------------------------------------------

// Ownership fetches fund ownership data for the requested symbols.
func (c *Client) Ownership(ctx context.Context, req OwnershipRequest) (*OwnershipResponse, error) {
	vars := ownershipVariables(req)

	var resp OwnershipResponse
	if err := c.doGraphQL(ctx, "Ownership", vars, queryOwnership, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
