package marketsurge

import "context"

// queryChartMarketData is the GraphQL query for the daily ChartMarketData endpoint.
const queryChartMarketData = `query ChartMarketData(
  $symbols: [String!]!
  $symbolDialectType: MDSymbolDialectType!
  $where: TimeSeriesFilterInput!
  $exchangeName: String!
) {
  marketData(symbols: $symbols, symbolDialectType: $symbolDialectType) {
    id
    originRequest {
      fromDialect
      symbol
    }
    pricing {
      timeSeries(where: $where) {
        period
        dataPoints {
          startDateTime
          endDateTime
          volume { value }
          last { value }
          low { value }
          high { value }
          open { value }
        }
      }
      quote {
        tradeDateTime
        timeliness
        quoteType
        volume { value formattedValue }
        percentChange { value formattedValue }
        netChange { value formattedValue }
        last { value formattedValue }
      }
      premarketQuote {
        last { value formattedValue }
        tradeDateTime
        timeliness
        volume { value formattedValue }
        percentChange { value formattedValue }
        quoteType
        netChange { value formattedValue }
      }
      postmarketQuote {
        volume { value formattedValue }
        tradeDateTime
        timeliness
        percentChange { formattedValue value }
        netChange { value formattedValue }
        quoteType
        last { value formattedValue }
      }
      currentMarketState
    }
  }
  exchangeData(exchangeName: $exchangeName) {
    city
    countryCode
    exchangeISO
    id
    holidays(
      where: {
        startDateTime: { gt: "2021-12-02T07:00:00.000Z" }
        endDateTime: { lt: "2026-03-27T23:55:25.000Z" }
      }
    ) {
      name
      holidayType
      description
      startDateTime
      endDateTime
    }
  }
}`

// queryChartMarketDataWeekly is the GraphQL query for the weekly ChartMarketData endpoint.
const queryChartMarketDataWeekly = `query ChartMarketData(
  $symbols: [String!]!
  $symbolDialectType: MDSymbolDialectType!
  $where: TimeSeriesFilterInput!
) {
  marketData(symbols: $symbols, symbolDialectType: $symbolDialectType) {
    id
    originRequest {
      fromDialect
      symbol
    }
    pricing {
      timeSeries(where: $where) {
        period
        dataPoints {
          startDateTime
          endDateTime
          volume { value }
          last { value }
          low { value }
          high { value }
          open { value }
        }
      }
      quote {
        tradeDateTime
        timeliness
        quoteType
        volume { value formattedValue }
        percentChange { value formattedValue }
        netChange { value formattedValue }
        last { value formattedValue }
      }
      premarketQuote {
        last { value formattedValue }
        tradeDateTime
        timeliness
        volume { value formattedValue }
        percentChange { value formattedValue }
        quoteType
        netChange { value formattedValue }
      }
      postmarketQuote {
        volume { value formattedValue }
        tradeDateTime
        timeliness
        percentChange { formattedValue value }
        netChange { value formattedValue }
        quoteType
        last { value formattedValue }
      }
      currentMarketState
    }
  }
}`

// DefaultChartSymbolDialectType is the default symbol dialect for chart queries.
const DefaultChartSymbolDialectType = "CHARTING"

// ---------------------------------------------------------------------------
// Request types
// ---------------------------------------------------------------------------

// ChartMarketDataRequest holds parameters for the daily ChartMarketData query.
type ChartMarketDataRequest struct {
	Symbols           []string
	SymbolDialectType string
	Where             TimeSeriesFilter
	ExchangeName      string
}

// TimeSeriesFilter holds the where clause for time series queries.
type TimeSeriesFilter struct {
	StartDateTime       string
	EndDateTime         string
	TimeSeriesType      string
	IncludeIntradayData bool
}

// ChartMarketDataWeeklyRequest holds parameters for the weekly ChartMarketData query.
// It is identical to ChartMarketDataRequest but without ExchangeName.
type ChartMarketDataWeeklyRequest struct {
	Symbols           []string
	SymbolDialectType string
	Where             TimeSeriesFilter
}

// NewChartMarketDataRequest creates a ChartMarketDataRequest with sensible defaults.
func NewChartMarketDataRequest(
	symbols []string,
	startDate, endDate, timeSeriesType, exchangeName string,
) ChartMarketDataRequest {
	return ChartMarketDataRequest{
		Symbols:           symbols,
		SymbolDialectType: DefaultChartSymbolDialectType,
		Where: TimeSeriesFilter{
			StartDateTime:       startDate,
			EndDateTime:         endDate,
			TimeSeriesType:      timeSeriesType,
			IncludeIntradayData: true,
		},
		ExchangeName: exchangeName,
	}
}

// NewChartMarketDataWeeklyRequest creates a ChartMarketDataWeeklyRequest
// with sensible defaults.
func NewChartMarketDataWeeklyRequest(
	symbols []string,
	startDate, endDate, timeSeriesType string,
) ChartMarketDataWeeklyRequest {
	return ChartMarketDataWeeklyRequest{
		Symbols:           symbols,
		SymbolDialectType: DefaultChartSymbolDialectType,
		Where: TimeSeriesFilter{
			StartDateTime:       startDate,
			EndDateTime:         endDate,
			TimeSeriesType:      timeSeriesType,
			IncludeIntradayData: true,
		},
	}
}

// ---------------------------------------------------------------------------
// GraphQL variables (unexported)
// ---------------------------------------------------------------------------

type chartMarketDataVariables struct {
	Symbols           []string              `json:"symbols"`
	SymbolDialectType string                `json:"symbolDialectType"`
	Where             timeSeriesFilterInput `json:"where"`
	ExchangeName      string                `json:"exchangeName"`
}

type chartMarketDataWeeklyVariables struct {
	Symbols           []string              `json:"symbols"`
	SymbolDialectType string                `json:"symbolDialectType"`
	Where             timeSeriesFilterInput `json:"where"`
}

type timeSeriesFilterInput struct {
	StartDateTime       eqFilter `json:"startDateTime"`
	EndDateTime         eqFilter `json:"endDateTime"`
	TimeSeriesType      eqFilter `json:"timeSeriesType"`
	IncludeIntradayData bool     `json:"includeIntradayData"`
}

type eqFilter struct {
	Eq string `json:"eq"`
}

// ---------------------------------------------------------------------------
// Response types
// ---------------------------------------------------------------------------

// ChartMarketDataResponse is the top-level response from the ChartMarketData query.
type ChartMarketDataResponse struct {
	MarketData   []ChartMarketDataItem `json:"marketData"`
	ExchangeData *ChartExchangeData    `json:"exchangeData"` // nil for weekly
}

// ChartMarketDataItem represents a single symbol's chart market data.
type ChartMarketDataItem struct {
	ID            string              `json:"id"`
	OriginRequest *ChartOriginRequest `json:"originRequest"`
	Pricing       *ChartPricing       `json:"pricing"`
}

// ChartOriginRequest holds the original request dialect and symbol.
type ChartOriginRequest struct {
	FromDialect string `json:"fromDialect"`
	Symbol      string `json:"symbol"`
}

// ChartPricing holds pricing data including time series and quotes.
type ChartPricing struct {
	TimeSeries         *ChartTimeSeries `json:"timeSeries"`
	Quote              *ChartQuote      `json:"quote"`
	PremarketQuote     *ChartQuote      `json:"premarketQuote"`
	PostmarketQuote    *ChartQuote      `json:"postmarketQuote"`
	CurrentMarketState *string          `json:"currentMarketState"`
}

// ChartTimeSeries holds time series data with a period and data points.
type ChartTimeSeries struct {
	Period     string           `json:"period"`
	DataPoints []ChartDataPoint `json:"dataPoints"`
}

// ChartDataPoint represents a single data point in a time series.
type ChartDataPoint struct {
	StartDateTime string      `json:"startDateTime"`
	EndDateTime   string      `json:"endDateTime"`
	Volume        *ChartValue `json:"volume"`
	Last          *ChartValue `json:"last"`
	Low           *ChartValue `json:"low"`
	High          *ChartValue `json:"high"`
	Open          *ChartValue `json:"open"`
}

// ChartValue holds a single numeric value.
type ChartValue struct {
	Value *float64 `json:"value"`
}

// ChartQuote holds real-time or extended-hours quote data.
type ChartQuote struct {
	TradeDateTime *string              `json:"tradeDateTime"`
	Timeliness    *string              `json:"timeliness"`
	QuoteType     *string              `json:"quoteType"`
	Volume        *ChartFormattedValue `json:"volume"`
	PercentChange *ChartFormattedValue `json:"percentChange"`
	NetChange     *ChartFormattedValue `json:"netChange"`
	Last          *ChartFormattedValue `json:"last"`
}

// ChartFormattedValue holds a numeric value with its formatted string.
type ChartFormattedValue struct {
	Value          *float64 `json:"value"`
	FormattedValue *string  `json:"formattedValue"`
}

// ChartExchangeData holds exchange information and holidays.
type ChartExchangeData struct {
	City        *string                `json:"city"`
	CountryCode *string                `json:"countryCode"`
	ExchangeISO *string                `json:"exchangeISO"`
	ID          *string                `json:"id"`
	Holidays    []ChartExchangeHoliday `json:"holidays"`
}

// ChartExchangeHoliday represents a single exchange holiday.
type ChartExchangeHoliday struct {
	Name          string  `json:"name"`
	HolidayType   *string `json:"holidayType"`
	Description   *string `json:"description"`
	StartDateTime string  `json:"startDateTime"`
	EndDateTime   string  `json:"endDateTime"`
}

// ---------------------------------------------------------------------------
// Client methods
// ---------------------------------------------------------------------------

// ChartMarketData fetches daily chart market data for the requested symbols.
func (c *Client) ChartMarketData(
	ctx context.Context,
	req ChartMarketDataRequest,
) (*ChartMarketDataResponse, error) {
	vars := chartMarketDataVariables{
		Symbols:           req.Symbols,
		SymbolDialectType: req.SymbolDialectType,
		Where: timeSeriesFilterInput{
			StartDateTime:       eqFilter{Eq: req.Where.StartDateTime},
			EndDateTime:         eqFilter{Eq: req.Where.EndDateTime},
			TimeSeriesType:      eqFilter{Eq: req.Where.TimeSeriesType},
			IncludeIntradayData: req.Where.IncludeIntradayData,
		},
		ExchangeName: req.ExchangeName,
	}

	var resp ChartMarketDataResponse
	if err := c.doGraphQL(ctx, "ChartMarketData", vars, queryChartMarketData, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ChartMarketDataWeekly fetches weekly chart market data for the requested symbols.
func (c *Client) ChartMarketDataWeekly(
	ctx context.Context,
	req ChartMarketDataWeeklyRequest,
) (*ChartMarketDataResponse, error) {
	vars := chartMarketDataWeeklyVariables{
		Symbols:           req.Symbols,
		SymbolDialectType: req.SymbolDialectType,
		Where: timeSeriesFilterInput{
			StartDateTime:       eqFilter{Eq: req.Where.StartDateTime},
			EndDateTime:         eqFilter{Eq: req.Where.EndDateTime},
			TimeSeriesType:      eqFilter{Eq: req.Where.TimeSeriesType},
			IncludeIntradayData: req.Where.IncludeIntradayData,
		},
	}

	var resp ChartMarketDataResponse
	if err := c.doGraphQL(ctx, "ChartMarketData", vars, queryChartMarketDataWeekly, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
