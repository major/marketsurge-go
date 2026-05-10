package marketsurge

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"time"
)

// queryOtherMarketData is the GraphQL query for the OtherMarketData endpoint.
// It contains {pattern_start_date} and {pattern_end_date} placeholders that
// are replaced with actual dates before sending.
const queryOtherMarketData = `query OtherMarketData(
  $symbols: [String!]!
  $symbolDialectType: MDSymbolDialectType!
  $upToHistoricalPeriodOffset: MDUpToQueryPeriodOffsetHistorical!
  $upToQueryPeriodOffset: MDUpToQueryPeriodOffsetFuture!
  $upToHistoricalPeriodForProfitMargin: MDUpToQueryPeriodOffsetHistorical!
) {
  marketData(symbols: $symbols, symbolDialectType: $symbolDialectType) {
    id
    originRequest {
      fromDialect
      symbol
    }
    ratings {
      compRating {
        value
        periodOffset
        period
      }
      rsRating(where: { periodOffset: { eq: CURRENT } }) {
        value
        periodOffset
        period
      }
      epsRating {
        value
        periodOffset
        period
      }
      smrRating {
        value
        periodOffset
        period
        letterValue
      }
      adRating(where: { periodOffset: { eq: CURRENT }, period: { eq: P12M } }) {
        letterValue
        period
        periodOffset
        value
      }
    }
    pricingStatistics {
      endOfDayStatistics {
        historicalPriceStatistics(where: { period: { eq: P1Q } }) {
          period
          periodOffset
          periodEndDate {
            value
            formattedValue
          }
          priceHighDate {
            value
            formattedValue
          }
          priceHigh {
            value
            formattedValue
          }
          priceLowDate {
            value
            formattedValue
          }
          priceLow {
            value
            formattedValue
          }
          priceClose {
            value
            formattedValue
          }
          pricePercentChange {
            value
            formattedValue
          }
        }
        pricingStartDate {
          value
        }
        pricingEndDate {
          value
        }
        volumeMovingAverages {
          value
          period
          periodOffset
        }
        avgDollarVolume50Day {
          value
          formattedValue
        }
        marketCapitalization {
          value
          formattedValue
        }
        averageTrueRangePercent {
          value
          formattedValue
          period
          periodOffset
        }
        antEvents {
          value
        }
        upDownVolumeRatio {
          value
          scalingFactor
          formattedValue
        }
        alpha {
          formattedValue
          scalingFactor
          value
        }
        beta {
          value
          scalingFactor
          formattedValue
        }
        shortInterest {
          daysToCover {
            formattedValue
            value
          }
          daysToCoverPercentChange {
            formattedValue
            value
          }
          percentOfFloat {
            value
            scalingFactor
            formattedValue
          }
          volume {
            value
            scalingFactor
            formattedValue
          }
        }
        blueDotDailyEvents {
          formattedValue
          value
        }
        blueDotWeeklyEvents {
          formattedValue
          value
        }
      }
      intradayStatistics {
        pricePercentChangeVs {
          formattedValue
          value
          subject
          period
        }
        volumePercentChangeVs {
          value
          formattedValue
          subject
          period
        }
        isDailyBlueDotEvent
        isWeeklyBlueDotEvent
        yield {
          formattedValue
          scalingFactor
          value
        }
        priceToCashFlowRatio {
          value
          scalingFactor
          formattedValue
        }
        forwardPriceToEarningsRatio {
          value
          scalingFactor
          formattedValue
        }
        priceToSalesRatio {
          value
          scalingFactor
          formattedValue
        }
        priceToEarningsRatio {
          formattedValue
          scalingFactor
          value
        }
        priceToEarningsVsSP500 {
          value
          scalingFactor
          formattedValue
        }
      }
    }
    corporateActions {
      dividendNextReportedExDate {
        formattedValue
        value
      }
      dividends {
        amount {
          formattedValue
        }
        changeIndicator
        exDate {
          value
        }
      }
      spinoffs {
        exDate {
          value
        }
      }
      splits {
        splitDate {
          value
        }
      }
    }
    symbology {
      company {
        companyName
        address
        address2
        phone
        businessDescription
        url
        city
        country
        stateProvince
      }
      instrument {
        subType
        ipoDate {
          value
        }
        ipoPrice {
          formattedValue
          value
          currencySymbolInfo {
            unitSymbol
          }
        }
      }
    }
    patternInfo {
      patterns(
        where: {
          baseStartDate: { value: { gte: "{pattern_start_date}" } }
          baseEndDate: { value: { lte: "{pattern_end_date}" } }
          periodicity: { eq: DAILY }
        }
      ) {
        ... on MDCupWithoutHandle {
          baseDepth {
            value
            scalingFactor
            formattedValue
          }
          avgVolumeRatePctOnPivot {
            value
            scalingFactor
            formattedValue
          }
          pricePctChangeOnPivot {
            value
            scalingFactor
            formattedValue
          }
          handleDepth {
            value
            scalingFactor
            formattedValue
          }
          handleLength
          cupLength
          cupEndDate {
            value
          }
          handleLowDate {
            value
          }
          handleStartDate {
            value
          }
          baseBottomDate {
            value
          }
          id
          baseLength
          baseNumber
          baseStatus
          patternType
          periodicity
          baseStage
          baseStartDate {
            value
          }
          baseEndDate {
            value
          }
          pivotPrice {
            currencySymbolInfo {
              mantissaPrecision
              unitSymbol
              isoCurrencyCode
              isSuffix
            }
            value
            scalingFactor
            formattedValue
          }
          pivotDate {
            value
          }
          pivotPriceDate {
            value
          }
          leftSideHighDate {
            value
          }
        }
        ... on MDCupWithHandle {
          baseDepth {
            value
            scalingFactor
            formattedValue
          }
          avgVolumeRatePctOnPivot {
            value
            scalingFactor
            formattedValue
          }
          pricePctChangeOnPivot {
            value
            scalingFactor
            formattedValue
          }
          handleDepth {
            value
            scalingFactor
            formattedValue
          }
          handleLength
          cupLength
          cupEndDate {
            value
          }
          handleLowDate {
            value
          }
          handleStartDate {
            value
          }
          baseBottomDate {
            value
          }
          id
          baseLength
          baseNumber
          baseStatus
          patternType
          periodicity
          baseStage
          baseStartDate {
            value
          }
          baseEndDate {
            value
          }
          pivotPrice {
            currencySymbolInfo {
              mantissaPrecision
              unitSymbol
              isoCurrencyCode
              isSuffix
            }
            value
            scalingFactor
            formattedValue
          }
          pivotDate {
            value
          }
          pivotPriceDate {
            value
          }
          leftSideHighDate {
            value
          }
        }
        ... on MDConsolidation {
          baseDepth {
            value
            scalingFactor
            formattedValue
          }
          avgVolumeRatePctOnPivot {
            value
            scalingFactor
            formattedValue
          }
          pricePctChangeOnPivot {
            value
            scalingFactor
            formattedValue
          }
          baseBottomDate {
            value
          }
          id
          baseLength
          baseNumber
          baseStatus
          patternType
          periodicity
          baseStage
          baseStartDate {
            value
          }
          baseEndDate {
            value
          }
          pivotPrice {
            value
            scalingFactor
            formattedValue
            currencySymbolInfo {
              mantissaPrecision
              unitSymbol
              isoCurrencyCode
              isSuffix
            }
          }
          pivotDate {
            value
          }
          pivotPriceDate {
            value
          }
          leftSideHighDate {
            value
          }
        }
        ... on MDFlatBase {
          baseDepth {
            value
            scalingFactor
            formattedValue
          }
          avgVolumeRatePctOnPivot {
            value
            scalingFactor
            formattedValue
          }
          pricePctChangeOnPivot {
            value
            scalingFactor
            formattedValue
          }
          baseBottomDate {
            value
          }
          id
          baseLength
          baseNumber
          baseStatus
          patternType
          periodicity
          baseStage
          baseStartDate {
            value
          }
          baseEndDate {
            value
          }
          pivotPrice {
            currencySymbolInfo {
              mantissaPrecision
              unitSymbol
              isoCurrencyCode
              isSuffix
            }
            value
            scalingFactor
            formattedValue
          }
          pivotDate {
            value
          }
          pivotPriceDate {
            value
          }
          leftSideHighDate {
            value
          }
        }
        ... on MDIpoBase {
          upBars
          blueBars
          stallBars
          upVolumeTotal {
            value
            scalingFactor
            formattedValue
          }
          downBars
          redBars
          supportBars
          downVolumeTotal {
            value
            scalingFactor
            formattedValue
          }
          volumePctChangeOnPivot {
            value
            scalingFactor
            formattedValue
          }
          baseDepth {
            value
            scalingFactor
            formattedValue
          }
          avgVolumeRatePctOnPivot {
            value
            scalingFactor
            formattedValue
          }
          pricePctChangeOnPivot {
            value
            scalingFactor
            formattedValue
          }
          baseBottomDate {
            value
          }
          id
          baseLength
          baseNumber
          baseStatus
          patternType
          periodicity
          baseStage
          baseStartDate {
            value
          }
          baseEndDate {
            value
          }
          pivotPrice {
            currencySymbolInfo {
              mantissaPrecision
              unitSymbol
              isoCurrencyCode
              isSuffix
            }
            value
            scalingFactor
            formattedValue
          }
          pivotDate {
            value
          }
          pivotPriceDate {
            value
          }
          leftSideHighDate {
            value
          }
        }
        ... on MDSaucerWithHandle {
          baseDepth {
            value
            scalingFactor
            formattedValue
          }
          avgVolumeRatePctOnPivot {
            value
            scalingFactor
            formattedValue
          }
          pricePctChangeOnPivot {
            value
            scalingFactor
            formattedValue
          }
          handleDepth {
            value
            scalingFactor
            formattedValue
          }
          handleLength
          cupLength
          cupEndDate {
            value
          }
          handleLowDate {
            value
          }
          handleStartDate {
            value
          }
          baseBottomDate {
            value
          }
          id
          baseLength
          baseNumber
          baseStatus
          patternType
          periodicity
          baseStage
          baseStartDate {
            value
          }
          baseEndDate {
            value
          }
          pivotPrice {
            currencySymbolInfo {
              mantissaPrecision
              unitSymbol
              isoCurrencyCode
              isSuffix
            }
            value
            scalingFactor
            formattedValue
          }
          pivotDate {
            value
          }
          pivotPriceDate {
            value
          }
          leftSideHighDate {
            value
          }
        }
        ... on MDSaucerWithoutHandle {
          baseDepth {
            value
            scalingFactor
            formattedValue
          }
          avgVolumeRatePctOnPivot {
            value
            scalingFactor
            formattedValue
          }
          pricePctChangeOnPivot {
            value
            scalingFactor
            formattedValue
          }
          handleDepth {
            value
            scalingFactor
            formattedValue
          }
          handleLength
          cupLength
          cupEndDate {
            value
          }
          handleLowDate {
            value
          }
          handleStartDate {
            value
          }
          baseBottomDate {
            value
          }
          id
          baseLength
          baseNumber
          baseStatus
          patternType
          periodicity
          baseStage
          baseStartDate {
            value
          }
          baseEndDate {
            value
          }
          pivotPrice {
            currencySymbolInfo {
              mantissaPrecision
              unitSymbol
              isoCurrencyCode
              isSuffix
            }
            value
            scalingFactor
            formattedValue
          }
          pivotDate {
            value
          }
          pivotPriceDate {
            value
          }
          leftSideHighDate {
            value
          }
        }
        ... on MDAscendingBase {
          avgVolumeRatePctOnPivot {
            value
            scalingFactor
            formattedValue
          }
          pricePctChangeOnPivot {
            value
            scalingFactor
            formattedValue
          }
          firstBottomDate {
            value
          }
          secondAscendingHighDate {
            value
          }
          secondBottomDate {
            value
          }
          thirdAscendingHighDate {
            value
          }
          thirdBottomDate {
            value
          }
          pullBack1Depth {
            value
            scalingFactor
            formattedValue
          }
          pullBack2Depth {
            value
            scalingFactor
            formattedValue
          }
          pullBack3Depth {
            value
            scalingFactor
            formattedValue
          }
          id
          baseLength
          baseNumber
          baseStatus
          patternType
          periodicity
          baseStage
          baseStartDate {
            value
          }
          baseEndDate {
            value
          }
          pivotPrice {
            currencySymbolInfo {
              mantissaPrecision
              unitSymbol
              isoCurrencyCode
              isSuffix
            }
            value
            scalingFactor
            formattedValue
          }
          pivotDate {
            value
          }
          pivotPriceDate {
            value
          }
          leftSideHighDate {
            value
          }
        }
        ... on MDDoubleBottom {
          baseDepth {
            value
            scalingFactor
            formattedValue
          }
          avgVolumeRatePctOnPivot {
            value
            scalingFactor
            formattedValue
          }
          pricePctChangeOnPivot {
            value
            scalingFactor
            formattedValue
          }
          firstBottomDate {
            value
          }
          secondBottomDate {
            value
          }
          midPeakDate {
            value
          }
          id
          baseLength
          baseNumber
          baseStatus
          patternType
          periodicity
          baseStage
          baseStartDate {
            value
          }
          baseEndDate {
            value
          }
          pivotPrice {
            currencySymbolInfo {
              mantissaPrecision
              unitSymbol
              isoCurrencyCode
              isSuffix
            }
            value
            scalingFactor
            formattedValue
          }
          pivotDate {
            value
          }
          pivotPriceDate {
            value
          }
          leftSideHighDate {
            value
          }
        }
      }
      tightAreas(
        where: {
          startDate: { value: { gte: "{pattern_start_date}" } }
          endDate: { value: { gte: "{pattern_start_date}" } }
        }
      ) {
        patternID
        startDate {
          value
        }
        endDate {
          value
        }
        length
      }
    }
    financials {
      epsDueDate {
        value
        formattedValue
      }
      epsDueDateStatus
      epsLastReportedDate {
        value
      }
      consensusFinancials {
        eps {
          reportedEarnings(upToHistoricalPeriodOffset: $upToHistoricalPeriodOffset) {
            value {
              value
            }
            percentChangeYOY {
              value
            }
            periodOffset
            period
            periodEndDate {
              value
            }
            effectiveDate {
              value
            }
            percentSurprise {
              value
            }
            surpriseAmount {
              value
            }
            quarterNumber
            fiscalYear
          }
          growthRate(where: { period: { eq: P1Y } }) {
            value
            scalingFactor
            period
            formattedValue
          }
          earningsStability
        }
        sales {
          reportedSales(upToHistoricalPeriodOffset: $upToHistoricalPeriodOffset) {
            value {
              value
            }
            percentChangeYOY {
              value
            }
            periodOffset
            period
            periodEndDate {
              value
            }
            effectiveDate {
              value
            }
            percentSurprise {
              value
            }
            surpriseAmount {
              value
            }
            quarterNumber
            fiscalYear
          }
          growthRate(where: { period: { eq: P3Y } }) {
            formattedValue
            period
            scalingFactor
            value
          }
        }
      }
      cashFlowPerShareLastYear {
        formattedValue
        value
      }
      profitMarginValues(upToHistoricalPeriodOffset: $upToHistoricalPeriodForProfitMargin) {
        period
        preTaxMargin {
          formattedValue
          scalingFactor
          value
        }
        afterTaxMargin {
          value
        }
        grossMargin {
          value
        }
        returnOnEquity {
          value
          formattedValue
        }
        periodEndDate {
          value
          formattedValue
        }
        periodOffset
      }
      estimates {
        epsEstimates(upToQueryPeriodOffset: $upToQueryPeriodOffset) {
          revisionDirection
          effectiveDate {
            value
          }
          period
          value {
            value
          }
          percentChangeYOY {
            value
          }
          periodEndDate {
            value
          }
          type
        }
        salesEstimates(upToQueryPeriodOffset: $upToQueryPeriodOffset) {
          revisionDirection
          effectiveDate {
            value
          }
          period
          value {
            value
          }
          percentChangeYOY {
            value
          }
          periodEndDate {
            value
          }
        }
      }
    }
    industry {
      name
      sector
      indCode
      groupRanks {
        value
        period
        periodOffset
      }
      groupRS {
        value
        periodOffset
        letterValue
        period
      }
      numberOfStocksInGroup
    }
    ownership {
      fundsFloatPercentHeld {
        value
        scalingFactor
        formattedValue
      }
    }
    fundamentals {
      researchAndDevelopmentPercentLastQtr {
        value
        scalingFactor
        formattedValue
      }
      newCEODate {
        value
      }
      debtPercent {
        formattedValue
      }
    }
  }
}`

// Default variable values for the OtherMarketData query.
const (
	DefaultSymbolDialectType                   = "CHARTING"
	DefaultUpToHistoricalPeriodForProfitMargin = "P12Q_AGO"
	DefaultUpToHistoricalPeriodOffset          = "P24Q_AGO"
	DefaultUpToQueryPeriodOffset               = "P4Q_FUTURE"
)

// OtherMarketDataRequest holds parameters for the OtherMarketData query.
type OtherMarketDataRequest struct {
	Symbols                             []string
	SymbolDialectType                   string
	UpToHistoricalPeriodForProfitMargin string
	UpToHistoricalPeriodOffset          string
	UpToQueryPeriodOffset               string
	PatternStartDate                    string
	PatternEndDate                      string
}

// NewOtherMarketDataRequest creates an OtherMarketDataRequest with sensible
// defaults for the given symbols. PatternStartDate defaults to 4 years ago
// and PatternEndDate defaults to today.
func NewOtherMarketDataRequest(symbols ...string) OtherMarketDataRequest {
	now := time.Now().UTC()
	return OtherMarketDataRequest{
		Symbols:                             symbols,
		SymbolDialectType:                   DefaultSymbolDialectType,
		UpToHistoricalPeriodForProfitMargin: DefaultUpToHistoricalPeriodForProfitMargin,
		UpToHistoricalPeriodOffset:          DefaultUpToHistoricalPeriodOffset,
		UpToQueryPeriodOffset:               DefaultUpToQueryPeriodOffset,
		PatternStartDate:                    now.AddDate(-4, 0, 0).Format("2006-01-02"),
		PatternEndDate:                      now.Format("2006-01-02"),
	}
}

// otherMarketDataVariables holds the GraphQL variables sent with the query.
type otherMarketDataVariables struct {
	Symbols                             []string `json:"symbols"`
	SymbolDialectType                   string   `json:"symbolDialectType"`
	UpToHistoricalPeriodForProfitMargin string   `json:"upToHistoricalPeriodForProfitMargin"`
	UpToHistoricalPeriodOffset          string   `json:"upToHistoricalPeriodOffset"`
	UpToQueryPeriodOffset               string   `json:"upToQueryPeriodOffset"`
}

// ---------------------------------------------------------------------------
// Value types
// ---------------------------------------------------------------------------

// MDFormattedFloat holds a numeric value with an optional formatted string.
type MDFormattedFloat struct {
	Value          *float64 `json:"value"`
	FormattedValue *string  `json:"formattedValue"`
}

// MDScaledFloat holds a numeric value with scaling factor and formatted string.
type MDScaledFloat struct {
	Value          *float64 `json:"value"`
	FormattedValue *string  `json:"formattedValue"`
	ScalingFactor  *float64 `json:"scalingFactor"`
}

// MDDateValue holds a single date string value.
type MDDateValue struct {
	Value string `json:"value"`
}

// MDFormattedString holds a string value with an optional formatted representation.
type MDFormattedString struct {
	Value          *string `json:"value"`
	FormattedValue *string `json:"formattedValue"`
}

// MDValueWrapper wraps a single numeric value, used for nested {value} objects.
type MDValueWrapper struct {
	Value *float64 `json:"value"`
}

// MDCurrencyValue holds a currency amount with symbol information.
type MDCurrencyValue struct {
	Value              *float64              `json:"value"`
	FormattedValue     *string               `json:"formattedValue"`
	ScalingFactor      *float64              `json:"scalingFactor"`
	CurrencySymbolInfo *MDCurrencySymbolInfo `json:"currencySymbolInfo"`
}

// MDCurrencySymbolInfo describes how to format a currency value.
type MDCurrencySymbolInfo struct {
	MantissaPrecision *int    `json:"mantissaPrecision"`
	UnitSymbol        *string `json:"unitSymbol"`
	IsoCurrencyCode   *string `json:"isoCurrencyCode"`
	IsSuffix          *bool   `json:"isSuffix"`
}

// ---------------------------------------------------------------------------
// Response types
// ---------------------------------------------------------------------------

// OtherMarketDataResponse is the top-level response from the OtherMarketData query.
type OtherMarketDataResponse struct {
	MarketData []MarketDataItem `json:"marketData"`
}

// MarketDataItem represents a single symbol's market data.
type MarketDataItem struct {
	ID                *string              `json:"id"`
	OriginRequest     *MDOriginRequest     `json:"originRequest"`
	Ratings           *MDRatings           `json:"ratings"`
	PricingStatistics *MDPricingStatistics `json:"pricingStatistics"`
	CorporateActions  *MDCorporateActions  `json:"corporateActions"`
	Symbology         *MDSymbology         `json:"symbology"`
	PatternInfo       *MDPatternInfo       `json:"patternInfo"`
	Financials        *MDFinancials        `json:"financials"`
	Industry          *MDIndustry          `json:"industry"`
	Ownership         *MDOwnership         `json:"ownership"`
	Fundamentals      *MDFundamentals      `json:"fundamentals"`
}

// MDOriginRequest identifies the original symbol request.
type MDOriginRequest struct {
	FromDialect *string `json:"fromDialect"`
	Symbol      *string `json:"symbol"`
}

// ---------------------------------------------------------------------------
// Ratings
// ---------------------------------------------------------------------------

// MDRatings holds all rating categories for a symbol.
type MDRatings struct {
	CompRating []MDRating `json:"compRating"`
	RSRating   []MDRating `json:"rsRating"`
	EPSRating  []MDRating `json:"epsRating"`
	SMRRating  []MDRating `json:"smrRating"`
	ADRating   []MDRating `json:"adRating"`
}

// MDRating represents a single rating value with period metadata.
type MDRating struct {
	Value        *int    `json:"value"`
	PeriodOffset *string `json:"periodOffset"`
	Period       *string `json:"period"`
	LetterValue  *string `json:"letterValue"`
}

// ---------------------------------------------------------------------------
// Pricing statistics
// ---------------------------------------------------------------------------

// MDPricingStatistics groups end-of-day and intraday pricing data.
type MDPricingStatistics struct {
	EndOfDayStatistics *MDEndOfDayStatistics `json:"endOfDayStatistics"`
	IntradayStatistics *MDIntradayStatistics `json:"intradayStatistics"`
}

// MDEndOfDayStatistics holds end-of-day pricing metrics.
type MDEndOfDayStatistics struct {
	HistoricalPriceStatistics []MDHistoricalPriceStatistic `json:"historicalPriceStatistics"`
	PricingStartDate          *MDDateValue                 `json:"pricingStartDate"`
	PricingEndDate            *MDDateValue                 `json:"pricingEndDate"`
	VolumeMovingAverages      []MDVolumeMovingAverage      `json:"volumeMovingAverages"`
	AvgDollarVolume50Day      *MDFormattedFloat            `json:"avgDollarVolume50Day"`
	MarketCapitalization      *MDFormattedFloat            `json:"marketCapitalization"`
	AverageTrueRangePercent   []MDAverageTrueRangePercent  `json:"averageTrueRangePercent"`
	AntEvents                 []MDDateValue                `json:"antEvents"`
	UpDownVolumeRatio         *MDScaledFloat               `json:"upDownVolumeRatio"`
	Alpha                     *MDScaledFloat               `json:"alpha"`
	Beta                      *MDScaledFloat               `json:"beta"`
	ShortInterest             *MDShortInterest             `json:"shortInterest"`
	BlueDotDailyEvents        []MDFormattedString          `json:"blueDotDailyEvents"`
	BlueDotWeeklyEvents       []MDFormattedString          `json:"blueDotWeeklyEvents"`
}

// MDHistoricalPriceStatistic holds price statistics for a single period.
type MDHistoricalPriceStatistic struct {
	Period             *string            `json:"period"`
	PeriodOffset       *string            `json:"periodOffset"`
	PeriodEndDate      *MDFormattedString `json:"periodEndDate"`
	PriceHighDate      *MDFormattedString `json:"priceHighDate"`
	PriceHigh          *MDFormattedFloat  `json:"priceHigh"`
	PriceLowDate       *MDFormattedString `json:"priceLowDate"`
	PriceLow           *MDFormattedFloat  `json:"priceLow"`
	PriceClose         *MDFormattedFloat  `json:"priceClose"`
	PricePercentChange *MDFormattedFloat  `json:"pricePercentChange"`
}

// MDVolumeMovingAverage holds a volume moving average with period metadata.
type MDVolumeMovingAverage struct {
	Value        *float64 `json:"value"`
	Period       *string  `json:"period"`
	PeriodOffset *string  `json:"periodOffset"`
}

// MDAverageTrueRangePercent holds ATR percentage with period metadata.
type MDAverageTrueRangePercent struct {
	Value          *float64 `json:"value"`
	FormattedValue *string  `json:"formattedValue"`
	Period         *string  `json:"period"`
	PeriodOffset   *string  `json:"periodOffset"`
}

// MDShortInterest holds short interest metrics.
type MDShortInterest struct {
	DaysToCover              *MDFormattedFloat `json:"daysToCover"`
	DaysToCoverPercentChange *MDFormattedFloat `json:"daysToCoverPercentChange"`
	PercentOfFloat           *MDScaledFloat    `json:"percentOfFloat"`
	Volume                   *MDScaledFloat    `json:"volume"`
}

// MDIntradayStatistics holds intraday pricing metrics.
type MDIntradayStatistics struct {
	PricePercentChangeVs        []MDPercentChangeVs `json:"pricePercentChangeVs"`
	VolumePercentChangeVs       []MDPercentChangeVs `json:"volumePercentChangeVs"`
	IsDailyBlueDotEvent         *bool               `json:"isDailyBlueDotEvent"`
	IsWeeklyBlueDotEvent        *bool               `json:"isWeeklyBlueDotEvent"`
	Yield                       *MDScaledFloat      `json:"yield"`
	PriceToCashFlowRatio        *MDScaledFloat      `json:"priceToCashFlowRatio"`
	ForwardPriceToEarningsRatio *MDScaledFloat      `json:"forwardPriceToEarningsRatio"`
	PriceToSalesRatio           *MDScaledFloat      `json:"priceToSalesRatio"`
	PriceToEarningsRatio        *MDScaledFloat      `json:"priceToEarningsRatio"`
	PriceToEarningsVsSP500      *MDScaledFloat      `json:"priceToEarningsVsSP500"`
}

// MDPercentChangeVs holds a percent change relative to a subject/period.
type MDPercentChangeVs struct {
	Value          *float64 `json:"value"`
	FormattedValue *string  `json:"formattedValue"`
	Subject        *string  `json:"subject"`
	Period         *string  `json:"period"`
}

// ---------------------------------------------------------------------------
// Corporate actions
// ---------------------------------------------------------------------------

// MDCorporateActions holds dividend, split, and spinoff history.
type MDCorporateActions struct {
	DividendNextReportedExDate *MDFormattedString `json:"dividendNextReportedExDate"`
	Dividends                  []MDDividend       `json:"dividends"`
	Spinoffs                   []MDSpinoff        `json:"spinoffs"`
	Splits                     []MDSplit          `json:"splits"`
}

// MDDividend represents a single dividend event.
type MDDividend struct {
	Amount          *MDFormattedFloat `json:"amount"`
	ChangeIndicator *string           `json:"changeIndicator"`
	ExDate          *MDDateValue      `json:"exDate"`
}

// MDSpinoff represents a single spinoff event.
type MDSpinoff struct {
	ExDate *MDDateValue `json:"exDate"`
}

// MDSplit represents a single stock split event.
type MDSplit struct {
	SplitDate *MDDateValue `json:"splitDate"`
}

// ---------------------------------------------------------------------------
// Symbology
// ---------------------------------------------------------------------------

// MDSymbology holds company and instrument information.
type MDSymbology struct {
	Company    *MDCompany    `json:"company"`
	Instrument *MDInstrument `json:"instrument"`
}

// UnmarshalJSON accepts both historical object-shaped symbology data and the
// array-shaped data now returned by some live MarketSurge responses.
func (m *MDSymbology) UnmarshalJSON(data []byte) error {
	type mdSymbology MDSymbology
	var raw struct {
		*mdSymbology

		Company    json.RawMessage `json:"company"`
		Instrument json.RawMessage `json:"instrument"`
	}
	raw.mdSymbology = (*mdSymbology)(m)

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if err := m.decodeCompany(raw.Company); err != nil {
		return err
	}
	if err := m.decodeInstrument(raw.Instrument); err != nil {
		return err
	}

	return nil
}

func (m *MDSymbology) decodeCompany(data json.RawMessage) error {
	m.Company = nil

	data = bytes.TrimSpace(data)
	if len(data) == 0 || bytes.Equal(data, []byte("null")) {
		return nil
	}

	if data[0] == '[' {
		var companies []MDCompany
		if err := json.Unmarshal(data, &companies); err != nil {
			return err
		}
		if len(companies) == 0 {
			return nil
		}
		m.Company = &companies[0]
		return nil
	}

	var company MDCompany
	if err := json.Unmarshal(data, &company); err != nil {
		return err
	}
	m.Company = &company
	return nil
}

func (m *MDSymbology) decodeInstrument(data json.RawMessage) error {
	m.Instrument = nil

	data = bytes.TrimSpace(data)
	if len(data) == 0 || bytes.Equal(data, []byte("null")) {
		return nil
	}

	if data[0] == '[' {
		var instruments []MDInstrument
		if err := json.Unmarshal(data, &instruments); err != nil {
			return err
		}
		if len(instruments) == 0 {
			return nil
		}
		m.Instrument = &instruments[0]
		return nil
	}

	var instrument MDInstrument
	if err := json.Unmarshal(data, &instrument); err != nil {
		return err
	}
	m.Instrument = &instrument
	return nil
}

// MDCompany holds company profile information.
type MDCompany struct {
	CompanyName         *string `json:"companyName"`
	Address             *string `json:"address"`
	Address2            *string `json:"address2"`
	Phone               *string `json:"phone"`
	BusinessDescription *string `json:"businessDescription"`
	URL                 *string `json:"url"`
	City                *string `json:"city"`
	Country             *string `json:"country"`
	StateProvince       *string `json:"stateProvince"`
}

// MDInstrument holds instrument metadata.
type MDInstrument struct {
	SubType  *string          `json:"subType"`
	IPODate  *MDDateValue     `json:"ipoDate"`
	IPOPrice *MDCurrencyValue `json:"ipoPrice"`
}

// ---------------------------------------------------------------------------
// Pattern info
// ---------------------------------------------------------------------------

// MDPatternInfo holds chart patterns and tight areas.
type MDPatternInfo struct {
	Patterns   []MDPattern   `json:"patterns"`
	TightAreas []MDTightArea `json:"tightAreas"`
}

// MDPattern is a flat union of all chart pattern types. Fields specific to a
// particular pattern type (cup, saucer, ascending base, etc.) are nil when
// the pattern is a different type.
type MDPattern struct {
	// Common fields present in most pattern types.
	ID                      *string          `json:"id"`
	PatternType             *string          `json:"patternType"`
	Periodicity             *string          `json:"periodicity"`
	BaseStage               *string          `json:"baseStage"`
	BaseNumber              *int             `json:"baseNumber"`
	BaseStatus              *string          `json:"baseStatus"`
	BaseLength              *int             `json:"baseLength"`
	BaseDepth               *MDScaledFloat   `json:"baseDepth"`
	BaseStartDate           *MDDateValue     `json:"baseStartDate"`
	BaseEndDate             *MDDateValue     `json:"baseEndDate"`
	BaseBottomDate          *MDDateValue     `json:"baseBottomDate"`
	LeftSideHighDate        *MDDateValue     `json:"leftSideHighDate"`
	PivotPrice              *MDCurrencyValue `json:"pivotPrice"`
	PivotDate               *MDDateValue     `json:"pivotDate"`
	PivotPriceDate          *MDDateValue     `json:"pivotPriceDate"`
	AvgVolumeRatePctOnPivot *MDScaledFloat   `json:"avgVolumeRatePctOnPivot"`
	PricePctChangeOnPivot   *MDScaledFloat   `json:"pricePctChangeOnPivot"`

	// Cup/Saucer-specific fields.
	HandleDepth     *MDScaledFloat `json:"handleDepth"`
	HandleLength    *int           `json:"handleLength"`
	CupLength       *int           `json:"cupLength"`
	CupEndDate      *MDDateValue   `json:"cupEndDate"`
	HandleLowDate   *MDDateValue   `json:"handleLowDate"`
	HandleStartDate *MDDateValue   `json:"handleStartDate"`

	// IPO base-specific fields.
	UpBars                 *int           `json:"upBars"`
	BlueBars               *int           `json:"blueBars"`
	StallBars              *int           `json:"stallBars"`
	DownBars               *int           `json:"downBars"`
	RedBars                *int           `json:"redBars"`
	SupportBars            *int           `json:"supportBars"`
	UpVolumeTotal          *MDScaledFloat `json:"upVolumeTotal"`
	DownVolumeTotal        *MDScaledFloat `json:"downVolumeTotal"`
	VolumePctChangeOnPivot *MDScaledFloat `json:"volumePctChangeOnPivot"`

	// Ascending base-specific fields.
	FirstBottomDate         *MDDateValue   `json:"firstBottomDate"`
	SecondAscendingHighDate *MDDateValue   `json:"secondAscendingHighDate"`
	SecondBottomDate        *MDDateValue   `json:"secondBottomDate"`
	ThirdAscendingHighDate  *MDDateValue   `json:"thirdAscendingHighDate"`
	ThirdBottomDate         *MDDateValue   `json:"thirdBottomDate"`
	PullBack1Depth          *MDScaledFloat `json:"pullBack1Depth"`
	PullBack2Depth          *MDScaledFloat `json:"pullBack2Depth"`
	PullBack3Depth          *MDScaledFloat `json:"pullBack3Depth"`

	// Double bottom-specific fields.
	MidPeakDate *MDDateValue `json:"midPeakDate"`
}

// MDTightArea represents a tight price consolidation area on a chart.
type MDTightArea struct {
	PatternID *int         `json:"patternID"`
	StartDate *MDDateValue `json:"startDate"`
	EndDate   *MDDateValue `json:"endDate"`
	Length    *int         `json:"length"`
}

// ---------------------------------------------------------------------------
// Financials
// ---------------------------------------------------------------------------

// MDFinancials holds earnings, sales, margins, and estimate data.
type MDFinancials struct {
	EPSDueDate               *MDFormattedString     `json:"epsDueDate"`
	EPSDueDateStatus         *string                `json:"epsDueDateStatus"`
	EPSLastReportedDate      *MDDateValue           `json:"epsLastReportedDate"`
	ConsensusFinancials      *MDConsensusFinancials `json:"consensusFinancials"`
	CashFlowPerShareLastYear *MDFormattedFloat      `json:"cashFlowPerShareLastYear"`
	ProfitMarginValues       []MDProfitMarginValue  `json:"profitMarginValues"`
	Estimates                *MDEstimates           `json:"estimates"`
}

// MDConsensusFinancials groups EPS and sales consensus data.
type MDConsensusFinancials struct {
	EPS   *MDConsensusEPS   `json:"eps"`
	Sales *MDConsensusSales `json:"sales"`
}

// MDConsensusEPS holds consensus EPS data.
type MDConsensusEPS struct {
	ReportedEarnings  []MDReportedPeriod `json:"reportedEarnings"`
	GrowthRate        []MDGrowthRate     `json:"growthRate"`
	EarningsStability *int               `json:"earningsStability"`
}

// MDConsensusSales holds consensus sales data.
type MDConsensusSales struct {
	ReportedSales []MDReportedPeriod `json:"reportedSales"`
	GrowthRate    []MDGrowthRate     `json:"growthRate"`
}

// MDReportedPeriod holds a single reported earnings or sales period.
type MDReportedPeriod struct {
	Value            *MDValueWrapper `json:"value"`
	PercentChangeYOY *MDValueWrapper `json:"percentChangeYOY"`
	PeriodOffset     *string         `json:"periodOffset"`
	Period           *string         `json:"period"`
	PeriodEndDate    *MDDateValue    `json:"periodEndDate"`
	EffectiveDate    *MDDateValue    `json:"effectiveDate"`
	PercentSurprise  *MDValueWrapper `json:"percentSurprise"`
	SurpriseAmount   *MDValueWrapper `json:"surpriseAmount"`
	QuarterNumber    *int            `json:"quarterNumber"`
	FiscalYear       *int            `json:"fiscalYear"`
}

// MDGrowthRate holds a growth rate with scaling and period metadata.
type MDGrowthRate struct {
	Value          *float64 `json:"value"`
	ScalingFactor  *float64 `json:"scalingFactor"`
	Period         *string  `json:"period"`
	FormattedValue *string  `json:"formattedValue"`
}

// MDProfitMarginValue holds profit margin metrics for a single period.
type MDProfitMarginValue struct {
	Period         *string            `json:"period"`
	PreTaxMargin   *MDScaledFloat     `json:"preTaxMargin"`
	AfterTaxMargin *MDValueWrapper    `json:"afterTaxMargin"`
	GrossMargin    *MDValueWrapper    `json:"grossMargin"`
	ReturnOnEquity *MDFormattedFloat  `json:"returnOnEquity"`
	PeriodEndDate  *MDFormattedString `json:"periodEndDate"`
	PeriodOffset   *string            `json:"periodOffset"`
}

// MDEstimates groups EPS and sales estimate data.
type MDEstimates struct {
	EPSEstimates   []MDEstimate `json:"epsEstimates"`
	SalesEstimates []MDEstimate `json:"salesEstimates"`
}

// MDEstimate holds a single earnings or sales estimate.
type MDEstimate struct {
	RevisionDirection *string         `json:"revisionDirection"`
	EffectiveDate     *MDDateValue    `json:"effectiveDate"`
	Period            *string         `json:"period"`
	Value             *MDValueWrapper `json:"value"`
	PercentChangeYOY  *MDValueWrapper `json:"percentChangeYOY"`
	PeriodEndDate     *MDDateValue    `json:"periodEndDate"`
	Type              *string         `json:"type"`
}

// ---------------------------------------------------------------------------
// Industry, ownership, fundamentals
// ---------------------------------------------------------------------------

// MDIndustry holds industry group information.
type MDIndustry struct {
	Name                  *string       `json:"name"`
	Sector                *string       `json:"sector"`
	IndCode               *string       `json:"indCode"`
	GroupRanks            []MDGroupRank `json:"groupRanks"`
	GroupRS               []MDGroupRS   `json:"groupRS"`
	NumberOfStocksInGroup *int          `json:"numberOfStocksInGroup"`
}

// UnmarshalJSON accepts industry codes as either strings or numbers. Live
// MarketSurge responses have used both shapes for the same GraphQL field.
func (m *MDIndustry) UnmarshalJSON(data []byte) error {
	type mdIndustry MDIndustry
	var raw struct {
		*mdIndustry

		IndCode json.RawMessage `json:"indCode"`
	}
	raw.mdIndustry = (*mdIndustry)(m)

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	return m.decodeIndCode(raw.IndCode)
}

func (m *MDIndustry) decodeIndCode(data json.RawMessage) error {
	m.IndCode = nil
	data = bytes.TrimSpace(data)
	if len(data) == 0 || bytes.Equal(data, []byte("null")) {
		return nil
	}

	if data[0] == '"' {
		var value string
		if err := json.Unmarshal(data, &value); err != nil {
			return err
		}
		m.IndCode = &value
		return nil
	}

	var value json.Number
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	code := value.String()
	m.IndCode = &code
	return nil
}

// MDGroupRank holds an industry group rank for a period.
type MDGroupRank struct {
	Value        *int    `json:"value"`
	Period       *string `json:"period"`
	PeriodOffset *string `json:"periodOffset"`
}

// MDGroupRS holds an industry group relative strength for a period.
type MDGroupRS struct {
	Value        *int    `json:"value"`
	PeriodOffset *string `json:"periodOffset"`
	LetterValue  *string `json:"letterValue"`
	Period       *string `json:"period"`
}

// MDOwnership holds fund ownership metrics.
type MDOwnership struct {
	FundsFloatPercentHeld *MDScaledFloat `json:"fundsFloatPercentHeld"`
}

// MDFundamentals holds fundamental financial data.
type MDFundamentals struct {
	ResearchAndDevelopmentPercentLastQtr *MDScaledFloat    `json:"researchAndDevelopmentPercentLastQtr"`
	NewCEODate                           *MDDateValue      `json:"newCEODate"`
	DebtPercent                          *MDFormattedFloat `json:"debtPercent"`
}

// ---------------------------------------------------------------------------
// Client method
// ---------------------------------------------------------------------------

// OtherMarketData fetches market data for the requested symbols.
func (c *Client) OtherMarketData(ctx context.Context, req OtherMarketDataRequest) (*OtherMarketDataResponse, error) {
	query := strings.ReplaceAll(queryOtherMarketData, "{pattern_start_date}", req.PatternStartDate)
	query = strings.ReplaceAll(query, "{pattern_end_date}", req.PatternEndDate)

	vars := otherMarketDataVariables{
		Symbols:                             req.Symbols,
		SymbolDialectType:                   req.SymbolDialectType,
		UpToHistoricalPeriodForProfitMargin: req.UpToHistoricalPeriodForProfitMargin,
		UpToHistoricalPeriodOffset:          req.UpToHistoricalPeriodOffset,
		UpToQueryPeriodOffset:               req.UpToQueryPeriodOffset,
	}

	var resp OtherMarketDataResponse
	if err := c.doGraphQL(ctx, "OtherMarketData", vars, query, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
