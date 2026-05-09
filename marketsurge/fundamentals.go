package marketsurge

import "context"

// queryFundermentalDataBox is the GraphQL query for the FundermentalDataBox
// endpoint. The operation name preserves the upstream API's spelling.
const queryFundermentalDataBox = `query FundermentalDataBox(
  $symbols: [String!]!
  $symbolDialectType: MDSymbolDialectType!
  $upToHistoricalPeriodOffset: MDUpToQueryPeriodOffsetHistorical!
  $upToQueryPeriodOffset: MDUpToQueryPeriodOffsetFuture!
  $reportedSalesUpToHistoricalPeriod2: MDUpToQueryPeriodOffsetHistorical!
  $salesEstimatesUpToQueryPeriod2: MDUpToQueryPeriodOffsetFuture!
) {
  marketData(symbols: $symbols, symbolDialectType: $symbolDialectType) {
    financials {
      consensusFinancials {
        eps {
          reportedEarnings(upToHistoricalPeriodOffset: $upToHistoricalPeriodOffset) {
            value {
              formattedValue
              value
            }
            percentChangeYOY {
              formattedValue
              value
            }
            periodOffset
            periodEndDate {
              value
            }
          }
        }
        sales {
          reportedSales(upToHistoricalPeriodOffset: $reportedSalesUpToHistoricalPeriod2) {
            value {
              formattedValue
              value
            }
            percentChangeYOY {
              formattedValue
              value
            }
            periodEndDate {
              value
            }
            periodOffset
          }
        }
      }
      estimates {
        epsEstimates(upToQueryPeriodOffset: $upToQueryPeriodOffset) {
          value {
            value
            formattedValue
          }
          percentChangeYOY {
            value
            formattedValue
          }
          periodOffset
          period
          revisionDirection
        }
        salesEstimates(upToQueryPeriodOffset: $salesEstimatesUpToQueryPeriod2) {
          value {
            value
            formattedValue
          }
          percentChangeYOY {
            value
            formattedValue
          }
          periodOffset
          period
        }
      }
    }
    id
    symbology {
      company {
        companyName
      }
      instrument {
        symbols {
          value
          type
        }
      }
    }
  }
}`

// Default variable values for the FundermentalDataBox query.
const (
	DefaultFundamentalsSymbolDialectType                  = "CHARTING"
	DefaultFundamentalsUpToHistoricalPeriodOffset         = "P7Y_AGO"
	DefaultFundamentalsUpToQueryPeriodOffset              = "P2Y_FUTURE"
	DefaultFundamentalsReportedSalesUpToHistoricalPeriod2 = "P7Y_AGO"
	DefaultFundamentalsSalesEstimatesUpToQueryPeriod2     = "P2Y_FUTURE"
)

// FundamentalsRequest holds parameters for the FundermentalDataBox query.
type FundamentalsRequest struct {
	Symbols                            []string
	SymbolDialectType                  string
	UpToHistoricalPeriodOffset         string
	UpToQueryPeriodOffset              string
	ReportedSalesUpToHistoricalPeriod2 string
	SalesEstimatesUpToQueryPeriod2     string
}

// NewFundamentalsRequest creates a FundamentalsRequest with sensible defaults
// for the given symbols.
func NewFundamentalsRequest(symbols ...string) FundamentalsRequest {
	return FundamentalsRequest{
		Symbols:                            symbols,
		SymbolDialectType:                  DefaultFundamentalsSymbolDialectType,
		UpToHistoricalPeriodOffset:         DefaultFundamentalsUpToHistoricalPeriodOffset,
		UpToQueryPeriodOffset:              DefaultFundamentalsUpToQueryPeriodOffset,
		ReportedSalesUpToHistoricalPeriod2: DefaultFundamentalsReportedSalesUpToHistoricalPeriod2,
		SalesEstimatesUpToQueryPeriod2:     DefaultFundamentalsSalesEstimatesUpToQueryPeriod2,
	}
}

// fundamentalsVariables holds the GraphQL variables sent with the query.
type fundamentalsVariables struct {
	Symbols                            []string `json:"symbols"`
	SymbolDialectType                  string   `json:"symbolDialectType"`
	UpToHistoricalPeriodOffset         string   `json:"upToHistoricalPeriodOffset"`
	UpToQueryPeriodOffset              string   `json:"upToQueryPeriodOffset"`
	ReportedSalesUpToHistoricalPeriod2 string   `json:"reportedSalesUpToHistoricalPeriod2"`
	SalesEstimatesUpToQueryPeriod2     string   `json:"salesEstimatesUpToQueryPeriod2"`
}

// ---------------------------------------------------------------------------
// Value types
// ---------------------------------------------------------------------------

// FundamentalsFormattedValue holds a numeric value with an optional formatted string.
type FundamentalsFormattedValue struct {
	Value          *float64 `json:"value"`
	FormattedValue *string  `json:"formattedValue"`
}

// FundamentalsDateValue holds a single date string value.
type FundamentalsDateValue struct {
	Value string `json:"value"`
}

// ---------------------------------------------------------------------------
// Response types
// ---------------------------------------------------------------------------

// FundamentalsResponse is the top-level response from the FundermentalDataBox query.
type FundamentalsResponse struct {
	MarketData []FundamentalsItem `json:"marketData"`
}

// FundamentalsItem represents a single symbol's fundamental data.
type FundamentalsItem struct {
	ID         *string                 `json:"id"`
	Financials *FundamentalsFinancials `json:"financials"`
	Symbology  *FundamentalsSymbology  `json:"symbology"`
}

// FundamentalsFinancials holds consensus financials and estimates.
type FundamentalsFinancials struct {
	ConsensusFinancials *FundamentalsConsensus `json:"consensusFinancials"`
	Estimates           *FundamentalsEstimates `json:"estimates"`
}

// FundamentalsConsensus groups EPS and sales consensus data.
type FundamentalsConsensus struct {
	EPS   *FundamentalsConsensusEPS   `json:"eps"`
	Sales *FundamentalsConsensusSales `json:"sales"`
}

// FundamentalsConsensusEPS holds consensus EPS data with reported earnings.
type FundamentalsConsensusEPS struct {
	ReportedEarnings []FundamentalsReportedPeriod `json:"reportedEarnings"`
}

// FundamentalsConsensusSales holds consensus sales data with reported sales.
type FundamentalsConsensusSales struct {
	ReportedSales []FundamentalsReportedPeriod `json:"reportedSales"`
}

// FundamentalsReportedPeriod holds a single reported earnings or sales period.
type FundamentalsReportedPeriod struct {
	Value            *FundamentalsFormattedValue `json:"value"`
	PercentChangeYOY *FundamentalsFormattedValue `json:"percentChangeYOY"`
	PeriodOffset     *string                     `json:"periodOffset"`
	PeriodEndDate    *FundamentalsDateValue      `json:"periodEndDate"`
}

// FundamentalsEstimates groups EPS and sales estimate data.
type FundamentalsEstimates struct {
	EPSEstimates   []FundamentalsEstimate `json:"epsEstimates"`
	SalesEstimates []FundamentalsEstimate `json:"salesEstimates"`
}

// FundamentalsEstimate holds a single earnings or sales estimate.
type FundamentalsEstimate struct {
	Value             *FundamentalsFormattedValue `json:"value"`
	PercentChangeYOY  *FundamentalsFormattedValue `json:"percentChangeYOY"`
	PeriodOffset      *string                     `json:"periodOffset"`
	Period            *string                     `json:"period"`
	RevisionDirection *string                     `json:"revisionDirection"`
}

// ---------------------------------------------------------------------------
// Symbology
// ---------------------------------------------------------------------------

// FundamentalsSymbology holds company and instrument information.
type FundamentalsSymbology struct {
	Company    *FundamentalsCompany    `json:"company"`
	Instrument *FundamentalsInstrument `json:"instrument"`
}

// FundamentalsCompany holds company profile information.
type FundamentalsCompany struct {
	CompanyName *string `json:"companyName"`
}

// FundamentalsInstrument holds instrument metadata.
type FundamentalsInstrument struct {
	Symbols []FundamentalsSymbol `json:"symbols"`
}

// FundamentalsSymbol holds a symbol value and its dialect type.
type FundamentalsSymbol struct {
	Value *string `json:"value"`
	Type  *string `json:"type"`
}

// ---------------------------------------------------------------------------
// Client method
// ---------------------------------------------------------------------------

// Fundamentals fetches fundamental financial data for the requested symbols.
func (c *Client) Fundamentals(ctx context.Context, req FundamentalsRequest) (*FundamentalsResponse, error) {
	vars := fundamentalsVariables(req)

	var resp FundamentalsResponse
	if err := c.doGraphQL(ctx, "FundermentalDataBox", vars, queryFundermentalDataBox, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
