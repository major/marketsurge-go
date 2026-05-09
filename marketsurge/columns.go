package marketsurge

// Updating the column catalog from a HAR file
//
// The column definitions below are extracted from the MarketSurge web
// application's webpack bundle. MarketSurge does not expose a public API
// for listing available columns, so the only way to get this data is by
// capturing browser traffic and parsing the JavaScript bundle.
//
// Step 1: Capture a HAR file
//
//   Open Chrome/Firefox DevTools, go to the Network tab, and navigate to
//   marketsurge.investors.com. Use the site normally (open screens, run
//   reports) so the full app bundle loads. Export the network log as a HAR
//   file ("Save All as HAR").
//
// Step 2: Find the webpack bundle entry
//
//   Search the HAR file for entries whose URL matches "_app-*.js" or look
//   for the largest JavaScript response (typically 500KB+). The column
//   definitions live in the main application chunk, not in vendor bundles.
//
// Step 3: Locate the column definitions object
//
//   In the response body, search for a pattern like:
//
//     mdItemId:483,id:483,categoryId:9
//
//   or search for "EPSRating" or "SmartSelect". The columns are defined in
//   a single large JavaScript object literal (roughly 120KB of minified JS).
//   The object key for each column is a snake_case identifier; the value is
//   an object with fields: name, displayName, mdItemId, id, categoryId,
//   categoryName, instrumentCode, instrumentName, itemDescription, toolTip.
//
// Step 4: Extract and convert to JSON
//
//   The raw JS object is not valid JSON. You need to:
//
//   - Quote all object keys (name: -> "name":)
//   - Replace !0 with true and !1 with false
//   - Replace hex escapes like \xae with their Unicode equivalents
//   - Replace function references (p.toInteger, p.toDecimal, arrow
//     functions used as formatters) with null or a string placeholder
//   - Handle null values (already valid JSON)
//   - Wrap the entire block in { } if it is not already
//
//   This is the hardest part. The minified JS uses shorthand that breaks
//   naive regex-based conversion. Consider pasting the object into a
//   browser console and running JSON.stringify() on it, or use a JS
//   parser (e.g. Node.js) to evaluate the object in a sandbox.
//
// Step 5: Update this file
//
//   Compare the extracted JSON against the existing constants and catalog
//   entries below. Add new columns, update changed descriptions, and
//   remove any that no longer appear. Columns whose wire name contains
//   spaces cannot be Go constants and appear only in the Columns() catalog.
//   Duplicate wire names (same name, different internal IDs) get one
//   constant and one catalog entry for the primary (stock) variant.
//
// The webpack bundle also contains formatter functions for each column
// (p.toInteger, p.toDecimal, arrow functions for date formatting, etc.).
// These are not captured here because they are presentation logic, not
// data contract. The wire format always returns string values regardless
// of the column's display type.

// ColumnInfo describes a data column available in MarketSurge screen and
// report queries. Use the Name field as values in [RunScreenResponseColumn]
// and [AdhocScreenResponseColumn] fields.
type ColumnInfo struct {
	// Name is the wire name sent in GraphQL query responseColumns.
	Name string

	// DisplayName is the human-readable label shown in the MarketSurge UI.
	DisplayName string

	// Description explains what the column measures. Empty for columns
	// where MarketSurge provides no description.
	Description string

	// Category groups related columns (e.g. "Earnings", "Price & Volume").
	// Empty for uncategorized columns.
	Category string
}

// Column categories group related columns in the MarketSurge UI.
const (
	catSmartSelectRating    = "SmartSelect Rating"
	catEarnings             = "Earnings"
	catSales                = "Sales"
	catIndustrySector       = "Industry & Sector"
	catSharesHoldings       = "Shares & Holdings"
	catPriceVolume          = "Price & Volume"
	catMarginsRatios        = "Margins & Ratios"
	catGeneralStockCriteria = "General Stock Criteria"
	catIBDContent           = "IBD Content"
	catMiscellaneous        = "Miscellaneous"
	catGeneralFundCriteria  = "General Fund Criteria"
	catPerformance          = "Performance"
	catRisk                 = "Risk"
	catTaxEfficiency        = "Tax Efficiency"
	catExpenses             = "Expenses"
	catAlerts               = "Alerts"
)

// Repeated display names shared across multiple catalog entries.
const (
	displayEventDate    = "Event Date"
	displayGap          = "Gap %"
	displayGapFilled    = "Gap Filled"
	displayPctChgYTD    = "% Chg YTD"
	displayIdxPctChgQTD = "Index % Chg QTD"
)

// Repeated descriptions shared across multiple catalog entries.
const (
	descAfterTaxMarginRatio = "A financial performance ratio, calculated by dividing net income after taxes by net sales."
)

// SmartSelect Rating columns.
const (
	// ColumnEPSRating is the Earnings Per Share (EPS) Rating column.
	ColumnEPSRating = "EPSRating"

	// ColumnRSRating is the Relative Price Strength (RS) Rating column.
	ColumnRSRating = "RSRating"

	// ColumnIndustryGroupRSRatingLetter is the Industry Group Relative Strength (Grp RS) Rating column.
	ColumnIndustryGroupRSRatingLetter = "IndustryGroupRSRatingLetter"

	// ColumnSMRRating is the Sales + Profit Margins + ROE (SMR) Rating column.
	ColumnSMRRating = "SMRRating"

	// ColumnAccDisRating is the Accumulation/Distribution (Acc/Dis) Rating column.
	ColumnAccDisRating = "AccDisRating"

	// ColumnCompositeRating is the SmartSelect Composite Rating column.
	ColumnCompositeRating = "CompositeRating"
)

// Earnings columns.
const (
	// ColumnEarningsEffectiveDate is the Earnings Reported Last Date column.
	ColumnEarningsEffectiveDate = "EarningsEffectiveDate"

	// ColumnEPSDueDate is the EPS Due Date column.
	ColumnEPSDueDate = "EPSDueDate"

	// ColumnEpsPctChgYoyLastReportedQ is the % Change in Latest Quarter's EPS vs. Same Quarter Prior Year column.
	ColumnEpsPctChgYoyLastReportedQ = "EpsPctChgYoyLastReportedQ"

	// ColumnEpsPctChgYoy1QAgo is the % Change in 1 Quarter Ago EPS vs. Same Quarter Prior Year column.
	ColumnEpsPctChgYoy1QAgo = "EpsPctChgYoy1QAgo"

	// ColumnEpsPctChgYoy2QAgo is the % Change in 2 Quarters Ago EPS vs. Same Quarter Prior Year column.
	ColumnEpsPctChgYoy2QAgo = "EpsPctChgYoy2QAgo"

	// ColumnEpsPctChgYoy3QAgo is the % Change in 3 Quarters Ago EPS vs. Same Quarter Prior Year column.
	ColumnEpsPctChgYoy3QAgo = "EpsPctChgYoy3QAgo"

	// ColumnEPSTrailing4Q is the EPS, Trailing 4 Qtrs column.
	ColumnEPSTrailing4Q = "EPSTrailing4Q"

	// ColumnEPSValueLastReportedY is the EPS, Fiscal Year - Lst Reptd Yr column.
	ColumnEPSValueLastReportedY = "EPSValueLastReportedY"

	// ColumnEPSValue1YAgo is the EPS, Fiscal Year - 1 Yr Ago column.
	ColumnEPSValue1YAgo = "EPSValue1YAgo"

	// ColumnEPSValue2YAgo is the EPS, Fiscal Year - 2 Yrs Ago column.
	ColumnEPSValue2YAgo = "EPSValue2YAgo"

	// ColumnEPSValue3YAgo is the EPS, Fiscal Year - 3 Yrs Ago column.
	ColumnEPSValue3YAgo = "EPSValue3YAgo"

	// ColumnEPSValue4YAgo is the EPS, Fiscal Year - 4 Yrs Ago column.
	ColumnEPSValue4YAgo = "EPSValue4YAgo"

	// ColumnEPSValue5YAgo is the EPS, Fiscal Year - 5 Yrs Ago column.
	ColumnEPSValue5YAgo = "EPSValue5YAgo"

	// ColumnEPSValue6YAgo is the EPS, Fiscal Year - 6 Yrs Ago column.
	ColumnEPSValue6YAgo = "EPSValue6YAgo"

	// ColumnEpsAnnualPctChgLastReportedY is the % Change in Fiscal Year EPS vs. Prior Year, Lst Reptd Yr column.
	ColumnEpsAnnualPctChgLastReportedY = "EpsAnnualPctChgLastReportedY"

	// ColumnEpsPctChg1YAgo is the % Change in Fiscal Year EPS vs. Prior Year, 1 Yr Ago column.
	ColumnEpsPctChg1YAgo = "EpsPctChg1YAgo"

	// ColumnEpsEstPctChgQ1 is the % Increase in Next Quarter's EPS Estimate From Same Quarter Prior Year Actual EPS
	// column.
	ColumnEpsEstPctChgQ1 = "EpsEstPctChgQ1"

	// ColumnEpsEstPctChgY1 is the % Increase in Current Annual EPS Estimate vs. Last Reported Annual EPS column.
	ColumnEpsEstPctChgY1 = "EpsEstPctChgY1"

	// ColumnEpsEstPctChgY2 is the % Increase in Next Year's Annual EPS Estimate vs. Current Annual EPS Estimate column.
	ColumnEpsEstPctChgY2 = "EpsEstPctChgY2"

	// ColumnSustainableGrowthPct is the Sustainable Growth Model Projected EPS Growth % column.
	ColumnSustainableGrowthPct = "SustainableGrowthPct"

	// ColumnEpsAccelerationLast3Q is the EPS Have Been Accelerating Over the Latest 3 Quarters (%) column.
	ColumnEpsAccelerationLast3Q = "EpsAccelerationLast3Q"

	// ColumnEpsPctGrowthRate1Y is the Annual % EPS Growth Rate of Last 1 Year column.
	ColumnEpsPctGrowthRate1Y = "EpsPctGrowthRate1Y"

	// ColumnEpsPctGrowthRate3Y is the Annual % EPS Growth Rate of Last 3 Years column.
	ColumnEpsPctGrowthRate3Y = "EpsPctGrowthRate3Y"

	// ColumnEpsPctGrowthRate5Y is the Annual % EPS Growth Rate of Last 5 Years column.
	ColumnEpsPctGrowthRate5Y = "EpsPctGrowthRate5Y"

	// ColumnEpsPctGrowthRate5YPctRank is the Annual % EPS Growth Rate of Last 5 Years Percent Rank column.
	ColumnEpsPctGrowthRate5YPctRank = "EpsPctGrowthRate5YPctRank"

	// ColumnEarningsStability is the Earnings Stability column.
	ColumnEarningsStability = "EarningsStability"

	// ColumnEpsAvgPctChgLast2Q is the Average EPS % Change for Last 2 Quarters column.
	ColumnEpsAvgPctChgLast2Q = "EpsAvgPctChgLast2Q"

	// ColumnEpsAvgPctChgLast3Q is the Average EPS % Change for Last 3 Quarters column.
	ColumnEpsAvgPctChgLast3Q = "EpsAvgPctChgLast3Q"

	// ColumnEpsAvgPctChgLast4Q is the Average EPS % Change for Last 4 Quarters column.
	ColumnEpsAvgPctChgLast4Q = "EpsAvgPctChgLast4Q"

	// ColumnEpsAvgPctChg5Q is the Average EPS % Change for Last 5 Quarters column.
	ColumnEpsAvgPctChg5Q = "EpsAvgPctChg5Q"

	// ColumnEpsAvgPctChg6Q is the Average EPS % Change for Last 6 Quarters column.
	ColumnEpsAvgPctChg6Q = "EpsAvgPctChg6Q"

	// ColumnEPSSurprisePct is the Last Quarter % Earnings Surprise column.
	ColumnEPSSurprisePct = "EPSSurprisePct"

	// ColumnEPSTrailing4QGtrEPS4YAgo is the EPS, Trailing 4 Qtrs Greater Than EPS, 4 Yrs Ago column.
	ColumnEPSTrailing4QGtrEPS4YAgo = "EPSTrailing4QGtrEPS4YAgo"

	// ColumnEPSLastYGtrEPS4YAgo is the EPS, Lst Fiscal Yr Greater Than EPS, 4 Yrs Ago column.
	ColumnEPSLastYGtrEPS4YAgo = "EPSLastYGtrEPS4YAgo"

	// ColumnEPSTrailing4QGeqEPSLastY is the EPS, Trailing 4 Qtrs Greater Than or Equal to EPS, Lst Fiscal Yr column.
	ColumnEPSTrailing4QGeqEPSLastY = "EPSTrailing4QGeqEPSLastY"

	// ColumnEPSPctChgLastQGtr3YGrowth is the EPS % Change Current Quarter Greater Than 3-Yr EPS Growth Rate column.
	ColumnEPSPctChgLastQGtr3YGrowth = "EPSPctChgLastQGtr3YGrowth"

	// ColumnIs3YGrowthGte5YGrowth is the EPS 3-Yr Growth Rate Greater Than or Equal to 5-Yr Growth Rate column.
	ColumnIs3YGrowthGte5YGrowth = "Is3YGrowthGte5YGrowth"
)

// Sales columns.
const (
	// ColumnSalesValueCurrentY is the Current Total Annual Sales (mil) column.
	ColumnSalesValueCurrentY = "SalesValueCurrentY"

	// ColumnSalesPctChgYoy1QAgo is the % Change Latest Quarter's Sales vs. Same Quarter Prior Year column.
	ColumnSalesPctChgYoy1QAgo = "SalesPctChgYoy1QAgo"

	// ColumnSalesPctChgLastY is the % Change in Latest Fiscal Year Sales vs. Prior Year column.
	ColumnSalesPctChgLastY = "SalesPctChgLastY"

	// ColumnSalesAccelerationLast2Q is the Sales Have Been Accelerating Over the Latest 2 Quarters (%) column.
	ColumnSalesAccelerationLast2Q = "SalesAccelerationLast2Q"

	// ColumnSalesAccelerationLast3Q is the Sales Have Been Accelerating Over the Latest 3 Quarters (%) column.
	ColumnSalesAccelerationLast3Q = "SalesAccelerationLast3Q"

	// ColumnSalesGrowthLast3Y is the Annual % Sales Growth Rate of Last 3 Years column.
	ColumnSalesGrowthLast3Y = "SalesGrowthLast3Y"

	// ColumnSalesGrowthLast5Y is the Annual % Sales Growth Rate of Last 5 Years column.
	ColumnSalesGrowthLast5Y = "SalesGrowthLast5Y"

	// ColumnSalesAvgPctChgLast2Q is the Average Sales % Change for Last 2 Quarters column.
	ColumnSalesAvgPctChgLast2Q = "SalesAvgPctChgLast2Q"

	// ColumnSalesAvgPctChgLast3Q is the Average Sales % Change for Last 3 Quarters column.
	ColumnSalesAvgPctChgLast3Q = "SalesAvgPctChgLast3Q"

	// ColumnSalesAvgPctChgLast4Q is the Average Sales % Change for Last 4 Quarters column.
	ColumnSalesAvgPctChgLast4Q = "SalesAvgPctChgLast4Q"

	// ColumnSalesAvgPctChgLast5Q is the Average Sales % Change for Last 5 Quarters column.
	ColumnSalesAvgPctChgLast5Q = "SalesAvgPctChgLast5Q"

	// ColumnSalesAvgPctChgLast6Q is the Average Sales % Change for Last 6 Quarters column.
	ColumnSalesAvgPctChgLast6Q = "SalesAvgPctChgLast6Q"
)

// Industry & Sector columns.
const (
	// ColumnIndustryGroupRank is the Company's Industry Group Rank column.
	ColumnIndustryGroupRank = "IndustryGroupRank"

	// ColumnIndustryName is the Industry Group Name column.
	ColumnIndustryName = "IndustryName"

	// ColumnSectorName is the Broad Sectors column.
	ColumnSectorName = "SectorName"
)

// Shares & Holdings columns.
const (
	// ColumnSharesOutstanding is the Shares Outstanding (1000s) column.
	ColumnSharesOutstanding = "SharesOutstanding"

	// ColumnSharesInFloat1000s is the Shares in Float (1000s) column.
	ColumnSharesInFloat1000s = "SharesInFloat1000s"

	// ColumnMarketCapIntraday is the Market Capitalization (mil) column.
	ColumnMarketCapIntraday = "MarketCapIntraday"

	// ColumnEnterpriseValue is the Enterprise Value (mil) column.
	ColumnEnterpriseValue = "EnterpriseValue"

	// ColumnPctHeldByFunds is the % of Stock Owned by Mutual Funds column.
	ColumnPctHeldByFunds = "PctHeldByFunds"

	// ColumnPctChgInFundOwnership is the % of the Number of Mutual Funds Owning for Current Quarter vs. Previous Quarter
	// column.
	ColumnPctChgInFundOwnership = "PctChgInFundOwnership"

	// ColumnFundsNumberHoldingOwnership is the Sponsorship, Mutual Funds Holding Position, Lst Reptd Qtr column.
	ColumnFundsNumberHoldingOwnership = "FundsNumberHoldingOwnership"
)

// Price & Volume columns.
const (
	// ColumnPrice is the Current Price column.
	ColumnPrice = "Price"

	// ColumnPricePctOff52WHigh is the Current Price in Relation to 52-Week High column.
	ColumnPricePctOff52WHigh = "PricePctOff52WHigh"

	// ColumnPricePctChg is the Price % Change column.
	ColumnPricePctChg = "PricePctChg"

	// ColumnPriceNetChg is the Price $ Chg column.
	ColumnPriceNetChg = "PriceNetChg"

	// ColumnRSLineNewHigh is the RS Line Making New High column.
	ColumnRSLineNewHigh = "RSLineNewHigh"

	// ColumnRSLineNewLow is the RS Line Making New Low column.
	ColumnRSLineNewLow = "RSLineNewLow"

	// ColumnRSRating3M is the Relative Strength, 3 Month column.
	ColumnRSRating3M = "RSRating3M"

	// ColumnRSRating6M is the Relative Strength, 6 Month column.
	ColumnRSRating6M = "RSRating6M"

	// ColumnTrailing26WPctPerfomanceVsSP500 is the Price, % Change vs S&P 500, 26 Weeks column.
	ColumnTrailing26WPctPerfomanceVsSP500 = "Trailing26WPctPerfomanceVsSP500"

	// ColumnVolume is the Volume, Intraday (1000s) column.
	ColumnVolume = "Volume"

	// ColumnWeeklyClosingRange is the Weekly Closing Range column.
	ColumnWeeklyClosingRange = "WeeklyClosingRange"

	// ColumnDailyClosingRange is the Daily Closing Range column.
	ColumnDailyClosingRange = "DailyClosingRange"

	// ColumnVolumeAvg50Day is the Current 50-Day Average Volume (1000s) column.
	ColumnVolumeAvg50Day = "VolumeAvg50Day"

	// ColumnVolumePctChgVs50DAvgVolume is the % Increase in Current Day's Volume vs. 50-Day Average Volume column.
	ColumnVolumePctChgVs50DAvgVolume = "VolumePctChgVs50DAvgVolume"

	// ColumnUpDownVolumeRatio is the Up / Down Volume Ratio column.
	ColumnUpDownVolumeRatio = "UpDownVolumeRatio"

	// ColumnVolumeDollarAvg50D is the Current 50- Day Average Daily $ Volume (1000s) column.
	ColumnVolumeDollarAvg50D = "VolumeDollarAvg50D"

	// ColumnPricePctChgVs50DaySMA is the Current Price % Above or Below 50-Day Moving Average Price column.
	ColumnPricePctChgVs50DaySMA = "PricePctChgVs50DaySMA"

	// ColumnPricePctChgVs200DaySMA is the Current Price % Above or Below 200-Day Moving Average Price column.
	ColumnPricePctChgVs200DaySMA = "PricePctChgVs200DaySMA"

	// ColumnPricePctChgCurrentWeek is the % Price Is Up in the Current Week column.
	ColumnPricePctChgCurrentWeek = "PricePctChgCurrentWeek"

	// ColumnPricePctChgLast1M is the % Price Is Up in the Latest 1 Month column.
	ColumnPricePctChgLast1M = "PricePctChgLast1M"

	// ColumnPricePctChgLast3M is the % Price Is Up in the Latest 3 Months column.
	ColumnPricePctChgLast3M = "PricePctChgLast3M"

	// ColumnPricePctChgLast6M is the % Price Is Up in the Latest 6 Months column.
	ColumnPricePctChgLast6M = "PricePctChgLast6M"

	// ColumnPricePctChgLast12M is the % Price Is Up in the Latest 12 Months column.
	ColumnPricePctChgLast12M = "PricePctChgLast12M"

	// ColumnPricePctChgYTD is the % Price Is Up for Year to Date column.
	ColumnPricePctChgYTD = "PricePctChgYTD"

	// ColumnAlpha is the Current Alpha column.
	ColumnAlpha = "Alpha"

	// ColumnBeta is the Current Beta column.
	ColumnBeta = "Beta"

	// ColumnATR30D is the Average True Range (30 Days) column.
	ColumnATR30D = "ATR30D"

	// ColumnPricePercentChangeVs10DaySMA is the Current Price % Above or Below 10-Day Moving Average Price column.
	ColumnPricePercentChangeVs10DaySMA = "PricePercentChangeVs10DaySMA"

	// ColumnPricePercentChangeVs21DaySMA is the Current Price % Above or Below 21-Day Moving Average Price column.
	ColumnPricePercentChangeVs21DaySMA = "PricePercentChangeVs21DaySMA"

	// ColumnPricePercentChangeVs150DaySMA is the Current Price % Above or Below 150-Day Moving Average Price column.
	ColumnPricePercentChangeVs150DaySMA = "PricePercentChangeVs150DaySMA"

	// ColumnVolumePctChg10W is the % Increase in Weekly Volume vs. 10-Week Average Volume column.
	ColumnVolumePctChg10W = "VolumePctChg10W"

	// ColumnIsSMA10DAboveSMA21DAndSMA21DAboveSMA50 is the 10 Day > 21 Day > 50 Day column.
	ColumnIsSMA10DAboveSMA21DAndSMA21DAboveSMA50 = "IsSMA10DAboveSMA21DAndSMA21DAboveSMA50"

	// ColumnIsSMA50DayAboveSMA150DayAndSMA150DayAboveSMA200Day is the 50-Day > 150-Day > 200-Day column.
	ColumnIsSMA50DayAboveSMA150DayAndSMA150DayAboveSMA200Day = "IsSMA50DayAboveSMA150DayAndSMA150DayAboveSMA200Day"

	// ColumnVolumeUp5D is the Current day's Volume greater than previous 5 days' Volume column.
	ColumnVolumeUp5D = "VolumeUp5D"

	// ColumnVolumeUp10D is the Current day's Volume greater than previous 10 days' Volume column.
	ColumnVolumeUp10D = "VolumeUp10D"

	// ColumnVolumeUp20D is the Current day's Volume greater than previous 20 days' Volume column.
	ColumnVolumeUp20D = "VolumeUp20D"

	// ColumnRSLineWithin5PctOfHigh is the RS Line Within 5% of New (52 week) High column.
	ColumnRSLineWithin5PctOfHigh = "RSLineWithin5PctOfHigh"

	// ColumnATRPct21D is the The average daily true percentage range over the past 21 days column.
	ColumnATRPct21D = "ATRPct21D"

	// ColumnATRPct30D is the The average daily true percentage range over the past 30 days column.
	ColumnATRPct30D = "ATRPct30D"

	// ColumnATRPct50D is the The average daily true percentage range over the past 50 days column.
	ColumnATRPct50D = "ATRPct50D"
)

// Margins & Ratios columns.
const (
	// ColumnProfitMarginAfterTaxLastReportedQ is the After Tax Profit Margin (Most Recent Reported Quarter) column.
	ColumnProfitMarginAfterTaxLastReportedQ = "ProfitMarginAfterTaxLastReportedQ"

	// ColumnProfitMarginAvgAfterTax2Q is the Average After Tax Profit Margin (Last 2 Quarters) column.
	ColumnProfitMarginAvgAfterTax2Q = "ProfitMarginAvgAfterTax2Q"

	// ColumnProfitMarginAvgAfterTax3Q is the Average After Tax Profit Margin (Last 3 Quarters) column.
	ColumnProfitMarginAvgAfterTax3Q = "ProfitMarginAvgAfterTax3Q"

	// ColumnProfitMarginAvgAfterTax4Q is the Average After Tax Profit Margin (Last 4 Quarters) column.
	ColumnProfitMarginAvgAfterTax4Q = "ProfitMarginAvgAfterTax4Q"

	// ColumnProfitMarginAvgAfterTax5Q is the Average After Tax Profit Margin (Last 5 Quarters) column.
	ColumnProfitMarginAvgAfterTax5Q = "ProfitMarginAvgAfterTax5Q"

	// ColumnProfitMarginAvgAfterTax6Q is the Average After Tax Profit Margin (Last 6 Quarters) column.
	ColumnProfitMarginAvgAfterTax6Q = "ProfitMarginAvgAfterTax6Q"

	// ColumnAfterTaxMarginAccelerationLast3Q is the Companies With Accelerating After Tax Profit Margins During Latest 3
	// Quarters (%) column.
	ColumnAfterTaxMarginAccelerationLast3Q = "AfterTaxMarginAccelerationLast3Q"

	// ColumnProfitMarginPreTaxLastReportedY is the Pre-tax Annual Margins (Latest Fiscal Year Reported) column.
	ColumnProfitMarginPreTaxLastReportedY = "ProfitMarginPreTaxLastReportedY"

	// ColumnOpMargGeqIndMedian is the Operating Margin Greater Than or Equal to Industry Median column.
	ColumnOpMargGeqIndMedian = "OpMargGeqIndMedian"

	// ColumnProfitMarginGeqIndustryMedian is the Profit Margin Greater Than or Equal to Industry Median column.
	ColumnProfitMarginGeqIndustryMedian = "ProfitMarginGeqIndustryMedian"

	// ColumnCashFlowVsEPSPctLastQ is the Cash Flow vs. EPS % Difference - Last Reported Quarter column.
	ColumnCashFlowVsEPSPctLastQ = "CashFlowVsEPSPctLastQ"

	// ColumnCashFlowVsEPSPctLastY is the Cash Flow vs. EPS % Difference - Last Reported Year column.
	ColumnCashFlowVsEPSPctLastY = "CashFlowVsEPSPctLastY"

	// ColumnROE is the ROE (Latest Fiscal Year Reported) column.
	ColumnROE = "ROE"

	// ColumnROEAvg5Y is the ROE, 5-Year Average column.
	ColumnROEAvg5Y = "ROEAvg5Y"

	// ColumnDebtToEquityRatioLastY is the Debt % (Latest Fiscal Year Reported) column.
	ColumnDebtToEquityRatioLastY = "DebtToEquityRatioLastY"

	// ColumnPriceEarningsRatioForward is the Estimated P/E Going Out 1 Year column.
	ColumnPriceEarningsRatioForward = "PriceEarningsRatioForward"

	// ColumnPriceEarningsRatio is the Current Trailing 12 Month P/E column.
	ColumnPriceEarningsRatio = "PriceEarningsRatio"

	// ColumnPEVsSP500PE is the P/E vs Current S&P 500 P/E column.
	ColumnPEVsSP500PE = "PEVsSP500PE"

	// ColumnPEPercentileRank is the P/E Percent Rank column.
	ColumnPEPercentileRank = "PEPercentileRank"

	// ColumnPEPercentileRankInGroup is the P/E Number Rank in Group column.
	ColumnPEPercentileRankInGroup = "PEPercentileRankInGroup"

	// ColumnPELessThan5YAvg is the P/E Less Than 5-Year Average column.
	ColumnPELessThan5YAvg = "PELessThan5YAvg"

	// ColumnPEGRatio is the P/E to Earnings Growth Rate (PEG) column.
	ColumnPEGRatio = "PEGRatio"

	// ColumnPEGDividendAdjusted is the P/E to Earnings Growth Rate Plus Dividend Yield (Dividend-Adjusted PEG) column.
	ColumnPEGDividendAdjusted = "PEGDividendAdjusted"

	// ColumnPriceToSales is the Current Price to Sales Ratio column.
	ColumnPriceToSales = "PriceToSales"

	// ColumnPriceToBookRatio is the Current Price To Book Value column.
	ColumnPriceToBookRatio = "PriceToBookRatio"

	// ColumnPriceToCashFlow is the Current Price to Cash Flow column.
	ColumnPriceToCashFlow = "PriceToCashFlow"

	// ColumnEnterpriseValueToFreeCashFlow is the Current Enterprise Value to Free Cash Flow column.
	ColumnEnterpriseValueToFreeCashFlow = "EnterpriseValueToFreeCashFlow"

	// ColumnCurrentRatio is the Current Ratio column.
	ColumnCurrentRatio = "CurrentRatio"

	// ColumnYieldPct is the Current Yield % column.
	ColumnYieldPct = "YieldPct"

	// ColumnLongTermDebtToWorkingCapitalRatio is the Long-term Debt to Working Capital Ratio column.
	ColumnLongTermDebtToWorkingCapitalRatio = "LongTermDebtToWorkingCapitalRatio"

	// ColumnLiabilitesToAssetsLessThanIndustryMedian is the Total Liabilities/Total Assets Ratio Less Than Industry
	// Median column.
	ColumnLiabilitesToAssetsLessThanIndustryMedian = "LiabilitesToAssetsLessThanIndustryMedian"
)

// General Stock Criteria columns.
const (
	// ColumnExchange is the Exchange column.
	ColumnExchange = "Exchange"

	// ColumnIsAdr is the ADR column.
	ColumnIsAdr = "IsAdr"

	// ColumnIsETFOrClosedEndFund is the ETF/Closed - End Fund column.
	ColumnIsETFOrClosedEndFund = "IsETFOrClosedEndFund"

	// ColumnIsETF is the ETF column.
	ColumnIsETF = "IsETF"

	// ColumnIPODate is the IPO Date column.
	ColumnIPODate = "IPODate"

	// ColumnIncorpDate is the Date of Incorporation column.
	ColumnIncorpDate = "IncorpDate"

	// ColumnCity is the Company Headquarters - City (US and Canada only) or Country column.
	ColumnCity = "City"

	// ColumnCompanyDescription is the Company Description column.
	ColumnCompanyDescription = "CompanyDescription"

	// ColumnState is the Headquarters - State or Province column.
	ColumnState = "State"

	// ColumnNewCEOLast12M is the New CEO in Latest 12 Months column.
	ColumnNewCEOLast12M = "NewCEOLast12M"
)

// IBD Content columns.
const (
	// ColumnIBD50Flag is the IBD 50 Top Rated Stocks column.
	ColumnIBD50Flag = "IBD50Flag"

	// ColumnIBDNewAmericaFlag is the IBD New America Index column.
	ColumnIBDNewAmericaFlag = "IBDNewAmericaFlag"

	// ColumnIBD8585Flag is the IBD 85-85 Index column.
	ColumnIBD8585Flag = "IBD8585Flag"

	// ColumnIBDBigCap20Flag is the IBD Big Cap 20 column.
	ColumnIBDBigCap20Flag = "IBDBigCap20Flag"
)

// Miscellaneous columns.
const (
	// ColumnAccDisRatingPreviousW is the Accumulation/Distribution Letter Rating - Prior Week column.
	ColumnAccDisRatingPreviousW = "AccDisRatingPreviousW"

	// ColumnDividendAmountNextReported is the DIVIDEND, NEXT REPORTED AMOUNT column.
	ColumnDividendAmountNextReported = "DividendAmountNextReported"

	// ColumnDividendDateNextReported is the DIVIDEND, NEXT REPORTED EX-DIV DATE column.
	ColumnDividendDateNextReported = "DividendDateNextReported"

	// ColumnShortVolume is the Short Vol column.
	ColumnShortVolume = "ShortVolume"

	// ColumnShortInterestPctOfFloat is the SHORT INTEREST, % OF FLOAT column.
	ColumnShortInterestPctOfFloat = "ShortInterestPctOfFloat"

	// ColumnIndustryMarketValue is the Industry Group Market Value column.
	ColumnIndustryMarketValue = "IndustryMarketValue"

	// ColumnIndustryGroupNumberOfNewHighs is the Industry Group - Number of New High Stocks column.
	ColumnIndustryGroupNumberOfNewHighs = "IndustryGroupNumberOfNewHighs"

	// ColumnIndustryGroupPctOfNewHighs is the Industry Group - % New Highs in Group column.
	ColumnIndustryGroupPctOfNewHighs = "IndustryGroupPctOfNewHighs"

	// ColumnIndustryGroupNumberOfNewLows is the Industry Group - Number of New Low Stocks column.
	ColumnIndustryGroupNumberOfNewLows = "IndustryGroupNumberOfNewLows"

	// ColumnIndustryGroupPctOfNewLows is the Industry Group - % New Lows in Group column.
	ColumnIndustryGroupPctOfNewLows = "IndustryGroupPctOfNewLows"

	// ColumnIndustryGroupRankLastW is the Industry Group Rank, Last Week column.
	ColumnIndustryGroupRankLastW = "IndustryGroupRankLastW"

	// ColumnIndustryGroupRank12WAgo is the Industry Group Rank, 3 Months Ago column.
	ColumnIndustryGroupRank12WAgo = "IndustryGroupRank12WAgo"

	// ColumnIndustryGroupRank26WAgo is the Industry Group Rank, 6 Months Ago column.
	ColumnIndustryGroupRank26WAgo = "IndustryGroupRank26WAgo"

	// ColumnIndexPctChg5D is the Index % Chg vs 5 Days - Intraday column.
	ColumnIndexPctChg5D = "IndexPctChg5D"

	// ColumnAlertLastMarkUpDate is the Last Mark-up Date column.
	ColumnAlertLastMarkUpDate = "AlertLastMarkUpDate"

	// ColumnCompaniesInIndustryGroup is the Industry Group Number Of Companies In Industry Group column.
	ColumnCompaniesInIndustryGroup = "CompaniesInIndustryGroup"

	// ColumnPriceHighPreviousW is the Price, High Prior Week column.
	ColumnPriceHighPreviousW = "PriceHighPreviousW"

	// ColumnVolumeShortDaysCurrent is the current column.
	ColumnVolumeShortDaysCurrent = "VolumeShortDaysCurrent"

	// ColumnVolumeShortDays1PeriodAgo is the one month column.
	ColumnVolumeShortDays1PeriodAgo = "VolumeShortDays1PeriodAgo"

	// ColumnVolumeShortDays2PeriodsAgo is the two months column.
	ColumnVolumeShortDays2PeriodsAgo = "VolumeShortDays2PeriodsAgo"

	// ColumnShortInterestPctChg is the SHORT INTEREST, % CHG MO TO MO column.
	ColumnShortInterestPctChg = "ShortInterestPctChg"

	// ColumnPriceHigh52W is the 52-Week High Price column.
	ColumnPriceHigh52W = "PriceHigh52W"

	// ColumnPriceLow52W is the 52 Week Price Low column.
	ColumnPriceLow52W = "PriceLow52W"
)

// Mutual Fund - General columns.
const (
	// ColumnFundRankVsAllFunds is the IBD Rank vs All Funds column.
	ColumnFundRankVsAllFunds = "FundRankVsAllFunds"

	// ColumnFundNetAssets is the Net Assets column.
	ColumnFundNetAssets = "FundNetAssets"

	// ColumnFundObjective is the Objective column.
	ColumnFundObjective = "FundObjective"

	// ColumnFundManagerStartYear is the Mngr Start Year column.
	ColumnFundManagerStartYear = "FundManagerStartYear"

	// ColumnFundClosedToNewInvestors is the Closed to New Investors column.
	ColumnFundClosedToNewInvestors = "FundClosedToNewInvestors"

	// ColumnNAV is the NAV column.
	ColumnNAV = "NAV"

	// ColumnNAVChg is the NAV Change column.
	ColumnNAVChg = "NAVChg"

	// ColumnFundNavPctChg is the NAV % Change column.
	ColumnFundNavPctChg = "FundNavPctChg"

	ColumnHoldingsPctFundAssetsHeld = "HoldingsPctFundAssetsHeld"

	ColumnSharesHeldPct = "SharesHeldPct"
)

// Mutual Fund - Performance columns.
const (
	// ColumnFundYTDReturn is the Year-to-Date Return column.
	ColumnFundYTDReturn = "FundYTDReturn"

	// ColumnFundTotalReturn1M is the 1-Month Total Return column.
	ColumnFundTotalReturn1M = "FundTotalReturn1M"

	// ColumnFundTotalReturn3M is the 3-Month Total Return column.
	ColumnFundTotalReturn3M = "FundTotalReturn3M"

	// ColumnFundTotalReturn1Y is the 1-Year Total Return column.
	ColumnFundTotalReturn1Y = "FundTotalReturn1Y"

	// ColumnFundTotalReturn3Y is the 3-Year Total Return column.
	ColumnFundTotalReturn3Y = "FundTotalReturn3Y"

	// ColumnFundTotalReturn5Y is the 5-Year Total Return column.
	ColumnFundTotalReturn5Y = "FundTotalReturn5Y"

	// ColumnFundTotalReturn10Y is the 10-Year Total Return column.
	ColumnFundTotalReturn10Y = "FundTotalReturn10Y"
)

// Mutual Fund - Risk columns.
const (
	// ColumnFundStandardDeviation3Y is the Standard Deviation-3 Year column.
	ColumnFundStandardDeviation3Y = "FundStandardDeviation3Y"
)

// Mutual Fund - Tax Efficiency columns.
const (
	// ColumnFundTurnoverPct is the Turnover column.
	ColumnFundTurnoverPct = "FundTurnoverPct"
)

// Mutual Fund - Expenses columns.
const (
	// ColumnFundFrontEndLoad is the Front End Load column.
	ColumnFundFrontEndLoad = "FundFrontEndLoad"

	// ColumnFundExpenseRatio is the Expense Ratio column.
	ColumnFundExpenseRatio = "FundExpenseRatio"
)

// Alerts columns.
const (
	ColumnAlertPrice = "AlertPrice"

	ColumnAlertType = "AlertType"

	ColumnAlertCreatedDate = "AlertCreatedDate"

	ColumnAlertUpSidePerc = "AlertUpSidePerc"

	ColumnAlertDownSidePerc = "AlertDownSidePerc"

	ColumnAlertTriggered = "AlertTriggered"

	ColumnChgPercVsAlertPrice = "ChgPercVsAlertPrice"
)

// Other columns.
const (
	ColumnSymbol = "Symbol"

	ColumnCompanyName = "CompanyName"

	ColumnListRank = "ListRank"

	ColumnEPSPlusRSRating = "EPSPlusRSRating"

	ColumnBlueDotMostRecentDate = "BlueDotMostRecentDate"

	ColumnBlueDotCount45Day = "BlueDotCount45Day"

	ColumnRSLineMostRecentNewHighDate = "RSLineMostRecentNewHighDate"

	ColumnRSLineNewHighCount45Day = "RSLineNewHighCount45Day"

	Column50DayBreakOnVolumeMostRecentDate = "50DayBreakOnVolumeMostRecentDate"

	ColumnPullbackTo10WeekLineMostRecentDate = "PullbackTo10WeekLineMostRecentDate"

	ColumnEarningsGapDownPct = "EarningsGapDownPct"

	ColumnAntEventMostRecentDate = "AntEventMostRecentDate"

	ColumnTightAreasMostRecentDate = "TightAreasMostRecentDate"

	ColumnBasesFormingMostRecentDate = "BasesFormingMostRecentDate"

	ColumnPricePctChgVsPivotWeekly = "PricePctChgVsPivotWeekly"

	ColumnAntsEventCount6M = "AntsEventCount6M"

	ColumnBreakawayGapPct = "BreakawayGapPct"

	ColumnBaseTypeWeekly = "BaseTypeWeekly"

	ColumnBaseStageWeekly = "BaseStageWeekly"

	ColumnWeeklyPivotWeek = "WeeklyPivotWeek"

	ColumnPricePctChgToPivotWeekly = "PricePctChgToPivotWeekly"

	ColumnBreakawayGapMostRecentDate = "BreakawayGapMostRecentDate"

	ColumnIsBreakawayGapFilled = "IsBreakawayGapFilled"

	ColumnPricePctChgVsPivotDaily = "PricePctChgVsPivotDaily"

	ColumnBaseTypeDaily = "BaseTypeDaily"

	ColumnBaseStageDaily = "BaseStageDaily"

	ColumnGapUpEarningsMostRecentDate = "GapUpEarningsMostRecentDate"

	ColumnIsGapUpEarningsGapFilled = "IsGapUpEarningsGapFilled"

	ColumnGapDownEarningsMostRecentDate = "GapDownEarningsMostRecentDate"

	ColumnIsGapDownEarningsGapFilled = "IsGapDownEarningsGapFilled"

	ColumnEarningsGapUpPct = "EarningsGapUpPct"

	ColumnDatePriceHigh52W = "DatePriceHigh52W"

	ColumnPricePctChgQTD = "PricePctChgQTD"
)

// Columns returns a copy of the full column catalog. Each entry includes
// the wire name, display name, description, and category.
//
// Some wire names contain spaces and cannot be represented as Go constants.
// Those columns are only available through this catalog. Callers can filter
// or search the returned slice by any field.
func Columns() []ColumnInfo {
	result := make([]ColumnInfo, len(columns))
	copy(result, columns)
	return result
}

// columns is the full catalog of MarketSurge data columns, extracted from
// the web application's webpack bundle. See the file header comment for
// instructions on updating this list.
//
//nolint:gochecknoglobals // catalog data table, not mutable global state
var columns = []ColumnInfo{
	// SmartSelect Rating
	{
		Name:        ColumnEPSRating,
		DisplayName: "EPS Rating",
		Description: "A proprietary rating that evaluates a company’s earnings performance against all of the other companies in the market.",
		Category:    catSmartSelectRating,
	},
	{
		Name:        ColumnRSRating,
		DisplayName: "RS Rating",
		Description: "Measures a stock’s price movement over the last 12 months and compares it to all other stocks on the NYSE, AMEX & NASDAQ exchanges to measure and rate its price performance.",
		Category:    catSmartSelectRating,
	},
	{
		Name: ColumnIndustryGroupRSRatingLetter, DisplayName: "Ind Group RS",
		Description: "Measures a stock’s industry group performance over the past 6 months.",
		Category:    catSmartSelectRating,
	},
	{
		Name:        ColumnSMRRating,
		DisplayName: "SMR Rating",
		Description: "Using an A to E rating (A being highest) this measurement combines sales growth, pre-tax margins, after-tax margins and ROE.",
		Category:    catSmartSelectRating,
	},
	{
		Name:        ColumnAccDisRating,
		DisplayName: "A/D Rating",
		Description: "Allows you to interpret buying and selling activity in a stock. Using an A to E scale; A = Heavy Buying, E = Heavy Selling.",
		Category:    catSmartSelectRating,
	},
	{
		Name: ColumnCompositeRating, DisplayName: "Comp Rating",
		Description: "A rating that combines the other 5 SmartSelect Ratings into one, easy-to-use rating.",
		Category:    catSmartSelectRating,
	},

	// Earnings
	{
		Name: ColumnEarningsEffectiveDate, DisplayName: "EPS Lst Rptd",
		Description: "Earnings Reported Last Date",
		Category:    catEarnings,
	},
	{Name: ColumnEPSDueDate, DisplayName: "EPS Due Date", Description: "EPS Due Date", Category: catEarnings},
	{
		Name:        ColumnEpsPctChgYoyLastReportedQ,
		DisplayName: "EPS % Chg Last Qtr",
		Description: "Percentage change in earnings per share compared to the same quarter of the previous year. Amount is based on continuing operations.",
		Category:    catEarnings,
	},
	{
		Name:        ColumnEpsPctChgYoy1QAgo,
		DisplayName: "EPS % Chg 1 Q Ago",
		Description: "Percentage change in earnings per share as of one quarter ago compared to the same quarter of the previous year. Amount is based on continuing operations.",
		Category:    catEarnings,
	},
	{
		Name:        ColumnEpsPctChgYoy2QAgo,
		DisplayName: "EPS % Chg 2 Q Ago",
		Description: "Percentage change in earnings per share as of two quarters ago compared to the same quarter of the previous year. Amount is based on continuing operations.",
		Category:    catEarnings,
	},
	{
		Name:        ColumnEpsPctChgYoy3QAgo,
		DisplayName: "EPS % Chg 3 Q Ago",
		Description: "Percentage change in earnings per share as of three quarters ago compared to the same quarter of the previous year. Amount is based on continuing operations.",
		Category:    catEarnings,
	},
	{
		Name: ColumnEPSTrailing4Q, DisplayName: "EPS Trailing 4 Qtrs",
		Description: "Latest reported trailing four-quarters' earnings per share.",
		Category:    catEarnings,
	},
	{
		Name: ColumnEPSValueLastReportedY, DisplayName: "Fiscal EPS Lst Yr",
		Description: "Latest reported fiscal year earnings per share.",
		Category:    catEarnings,
	},
	{
		Name: ColumnEPSValue1YAgo, DisplayName: "Fiscal EPS 1 Yr Ago",
		Description: "Fiscal year earnings per share as of one year ago.",
		Category:    catEarnings,
	},
	{
		Name: ColumnEPSValue2YAgo, DisplayName: "Fiscal EPS 2 Yrs Ago",
		Description: "Fiscal year earnings per share as of two years ago.",
		Category:    catEarnings,
	},
	{
		Name: ColumnEPSValue3YAgo, DisplayName: "Fiscal EPS 3 Yrs Ago",
		Description: "Fiscal year earnings per share as of three years ago.",
		Category:    catEarnings,
	},
	{
		Name: ColumnEPSValue4YAgo, DisplayName: "Fiscal EPS 4 Yrs Ago",
		Description: "Fiscal year earnings per share as of four years ago.",
		Category:    catEarnings,
	},
	{
		Name: ColumnEPSValue5YAgo, DisplayName: "Fiscal EPS 5 Yrs Ago",
		Description: "Fiscal year earnings per share as of five years ago.",
		Category:    catEarnings,
	},
	{
		Name: ColumnEPSValue6YAgo, DisplayName: "Fiscal EPS 6 Yrs Ago",
		Description: "Fiscal year earnings per share as of six years ago.",
		Category:    catEarnings,
	},
	{
		Name:        ColumnEpsAnnualPctChgLastReportedY,
		DisplayName: "EPS % Chg Lst Yr",
		Description: "The percentage change of the latest fiscal year's earnings per share versus that of the prior year.",
		Category:    catEarnings,
	},
	{
		Name:        ColumnEpsPctChg1YAgo,
		DisplayName: "EPS % Chg 1 Yr Ago",
		Description: "The percentage change of the fiscal year's earnings per share from one year ago versus that of the prior year.",
		Category:    catEarnings,
	},
	{
		Name: ColumnEpsEstPctChgQ1, DisplayName: "EPS Est Cur Qtr %",
		Description: "Use this estimate to project whether your stocks exhibit positive indications of growth.",
		Category:    catEarnings,
	},
	{
		Name: ColumnEpsEstPctChgY1, DisplayName: "EPS Est Cur Yr %",
		Description: "Use this estimate to project whether your stocks exhibit positive indications of growth.",
		Category:    catEarnings,
	},
	{
		Name: ColumnEpsEstPctChgY2, DisplayName: "EPS Est Next Yr %",
		Description: "This percentage allows you to project long-term increase in growth.",
		Category:    catEarnings,
	},
	{
		Name:        ColumnSustainableGrowthPct,
		DisplayName: "Sustainable Growth %",
		Description: "This percentage allows you see expected growth based on a company reinvesting its earnings into the company.",
		Category:    catEarnings,
	},
	{
		Name:        ColumnEpsAccelerationLast3Q,
		DisplayName: "EPS Accel 3 Qtrs",
		Description: "Checking these criteria limits your stock list only to companies whose earnings have been accelerating over the last 3 quarters.",
		Category:    catEarnings,
	},
	{
		Name:        ColumnEpsPctGrowthRate1Y,
		DisplayName: "EPS % Growth 1 Yr",
		Description: "EPS growth rate is calculated by using a least squares regression fit over a 1 year period of earnings per share based on a trailing four-quarter count.",
		Category:    catEarnings,
	},
	{
		Name:        ColumnEpsPctGrowthRate3Y,
		DisplayName: "EPS % Growth 3 Yr",
		Description: "EPS growth rate is calculated by using a least squares regression fit over a 3 year period of earnings per share based on a trailing four-quarter count.",
		Category:    catEarnings,
	},
	{
		Name:        ColumnEpsPctGrowthRate5Y,
		DisplayName: "EPS % Growth 5 Yr",
		Description: "EPS growth rate is calculated by using a least squares regression fit over a 5 year period of earnings per share based on a trailing four-quarter count.",
		Category:    catEarnings,
	},
	{
		Name:        ColumnEpsPctGrowthRate5YPctRank,
		DisplayName: "EPS % Growth 5 Yr Pct Rnk",
		Description: "Percentile rank from 99 (highest) to 1 (lowest) of the 5-year annual % EPS growth rate, calculated by using a least squares regression fit over a 5 year period of earnings per share based on a trailing four-quarter count.",
		Category:    catEarnings,
	},
	{
		Name:        ColumnEarningsStability,
		DisplayName: "Earnings Stability",
		Description: "Indicates in percentage form one standard deviation of the variability of the 3 to 5 years earnings history with a scale ranging from 1 to 99.",
		Category:    catEarnings,
	},
	{
		Name: ColumnEpsAvgPctChgLast2Q, DisplayName: "Avg EPS % Chg 2Q",
		Description: "The average earnings per share (EPS) percentage change of the last 2 quarters.",
		Category:    catEarnings,
	},
	{
		Name: ColumnEpsAvgPctChgLast3Q, DisplayName: "Avg EPS % Chg 3Q",
		Description: "The average earnings per share (EPS) percentage change of the last 3 quarters.",
		Category:    catEarnings,
	},
	{
		Name: ColumnEpsAvgPctChgLast4Q, DisplayName: "Avg EPS % Chg 4Q",
		Description: "The average earnings per share (EPS) percentage change of the last 4 quarters.",
		Category:    catEarnings,
	},
	{
		Name: ColumnEpsAvgPctChg5Q, DisplayName: "Avg EPS % Chg 5Q",
		Description: "The average earnings per share (EPS) percentage change of the last 5 quarters.",
		Category:    catEarnings,
	},
	{
		Name: ColumnEpsAvgPctChg6Q, DisplayName: "Avg EPS % Chg 6Q",
		Description: "The average earnings per share (EPS) percentage change of the last 6 quarters.",
		Category:    catEarnings,
	},
	{
		Name: ColumnEPSSurprisePct, DisplayName: "EPS Surprise",
		Description: "The % difference between latest reported quarterly EPS and the composite estimate.",
		Category:    catEarnings,
	},
	{
		Name:        ColumnEPSTrailing4QGtrEPS4YAgo,
		DisplayName: "EPS Trl 4Q Gtr EPS 4 Yrs Ago",
		Description: "Indicates the latest reported trailing four-quarters' earnings per share (EPS) are greater than the EPS from four years ago.",
		Category:    catEarnings,
	},
	{
		Name:        ColumnEPSLastYGtrEPS4YAgo,
		DisplayName: "EPS Lst Yr Gtr EPS 4 Yrs Ago",
		Description: "Indicates the latest reported fiscal year's earnings per share (EPS) is greater than the EPS from four years ago.",
		Category:    catEarnings,
	},
	{
		Name:        ColumnEPSTrailing4QGeqEPSLastY,
		DisplayName: "EPS Trl 4Q Geq EPS Lst Fiscal Yr",
		Description: "Indicates the latest reported trailing four-quarters' earnings per share (EPS) are greater than or equal to the EPS from the last fiscal year.",
		Category:    catEarnings,
	},
	{
		Name:        ColumnEPSPctChgLastQGtr3YGrowth,
		DisplayName: "EPS % Chg Lst Q Gtr 3-Yr Growth",
		Description: "Indicates the current quarter's earnings per share (EPS) % change is greater than the 3-year EPS growth rate.",
		Category:    catEarnings,
	},
	{
		Name:        ColumnIs3YGrowthGte5YGrowth,
		DisplayName: "3-Yr EPS Growth Geq 5-Yr",
		Description: "Indicates the 3-year earnings per share (EPS) growth rate is greater than or equal to the 5-year EPS growth rate.",
		Category:    catEarnings,
	},

	// Sales
	{
		Name: ColumnSalesValueCurrentY, DisplayName: "Annual Sales (mil)",
		Description: "Total sales revenue reported by the company for the latest fiscal year.",
		Category:    catSales,
	},
	{
		Name: ColumnSalesPctChgYoy1QAgo, DisplayName: "Sales % Chg Lst Qtr",
		Description: "Most recent quarterly sales amount versus sales amount of same quarter in the previous year.",
		Category:    catSales,
	},
	{
		Name: ColumnSalesPctChgLastY, DisplayName: "Sales % Chg Lst Yr",
		Description: "Most recent fiscal year sales amount versus sales amount of the previous fiscal year.",
		Category:    catSales,
	},
	{
		Name: ColumnSalesAccelerationLast2Q, DisplayName: "Sales Accel 2 Qtrs",
		Description: "An increase in the sales growth rate quarter over quarter for the past two quarters.",
		Category:    catSales,
	},
	{
		Name: ColumnSalesAccelerationLast3Q, DisplayName: "Sales Accel 3 Qtrs",
		Description: "An increase in the sales growth rate quarter over quarter.",
		Category:    catSales,
	},
	{
		Name:        ColumnSalesGrowthLast3Y,
		DisplayName: "Sales Growth 3 Yr",
		Description: "A company's rate of increase in revenues (sales). Growth rate will be calculated only if a minimum of 11 quarters of sales exist.",
		Category:    catSales,
	},
	{
		Name: ColumnSalesGrowthLast5Y, DisplayName: "Sales Growth 5 Yr",
		Description: "Determinant of a company’s earnings growth over the last 5 years.",
		Category:    catSales,
	},
	{
		Name: ColumnSalesAvgPctChgLast2Q, DisplayName: "Avg Sales % Chg 2Q",
		Description: "The average sales percentage change of the last 2 quarters.",
		Category:    catSales,
	},
	{
		Name: ColumnSalesAvgPctChgLast3Q, DisplayName: "Avg Sales % Chg 3Q",
		Description: "The average sales percentage change of the last 3 quarters.",
		Category:    catSales,
	},
	{
		Name: ColumnSalesAvgPctChgLast4Q, DisplayName: "Avg Sales % Chg 4Q",
		Description: "The average sales percentage change of the last 4 quarters.",
		Category:    catSales,
	},
	{
		Name: ColumnSalesAvgPctChgLast5Q, DisplayName: "Avg Sales % Chg 5Q",
		Description: "The average sales percentage change of the last 5 quarters.",
		Category:    catSales,
	},
	{
		Name: ColumnSalesAvgPctChgLast6Q, DisplayName: "Avg Sales % Chg 6Q",
		Description: "The average sales percentage change of the last 6 quarters.",
		Category:    catSales,
	},

	// Industry & Sector
	{
		Name: ColumnIndustryGroupRank, DisplayName: "Ind Group Rank",
		Description: "A proprietary ranking of the Industry Group to which the company has been classified.",
		Category:    catIndustrySector,
	},
	{
		Name: ColumnIndustryName, DisplayName: "Industry Name",
		Description: "Companies are classified into one of 145 Industry Groups.",
		Category:    catIndustrySector,
	},
	{
		Name: ColumnSectorName, DisplayName: "Sector",
		Description: "Companies in the database are classified as one of the 33 broad market sectors.",
		Category:    catIndustrySector,
	},

	// Shares & Holdings
	{
		Name: ColumnSharesOutstanding, DisplayName: "Shares (1000s)",
		Description: "The number of common shares outstanding for the primary class of stock.",
		Category:    catSharesHoldings,
	},
	{
		Name: ColumnSharesInFloat1000s, DisplayName: "Shares in Float (1000s)",
		Description: "Shares in Float (1000s)",
		Category:    catSharesHoldings,
	},
	{
		Name: ColumnMarketCapIntraday, DisplayName: "Market Cap (mil)",
		Description: "The total market value of a company.",
		Category:    catSharesHoldings,
	},
	{
		Name: ColumnEnterpriseValue, DisplayName: "Enterprise Val (mil)",
		Description: "Market value plus debt, less cash and securities.",
		Category:    catSharesHoldings,
	},
	{
		Name: ColumnPctHeldByFunds, DisplayName: "Funds %",
		Description: "Represents the % of stock (in float) that is owned by mutual funds.",
		Category:    catSharesHoldings,
	},
	{
		Name:        ColumnPctChgInFundOwnership,
		DisplayName: "Funds % Increase",
		Description: "Compares how many Mutual Funds owned shares in the stock for the recent quarter vs. the previous quarter.",
		Category:    catSharesHoldings,
	},
	{
		Name: ColumnFundsNumberHoldingOwnership, DisplayName: "Number of Funds",
		Description: "The total number of mutual funds holding a position in the stock as of the last reporting date",
		Category:    catSharesHoldings,
	},

	// Price & Volume
	{
		Name: ColumnPrice, DisplayName: "Current Price",
		Description: "The price a security most recently traded for.",
		Category:    catPriceVolume,
	},
	{
		Name:        ColumnPricePctOff52WHigh,
		DisplayName: "% Off High",
		Description: "Use the left and right sliders to specify a percentage range for stocks trading at a new 52-week high.",
		Category:    catPriceVolume,
	},
	{
		Name: ColumnPricePctChg, DisplayName: "Price % Chg",
		Description: "Percentage change, intraday price versus previous day's closing price.",
		Category:    catPriceVolume,
	},
	{
		Name: ColumnPriceNetChg, DisplayName: "Price $ Chg",
		Description: "The intraday price minus the last closing price of the stock.",
		Category:    catPriceVolume,
	},
	{
		Name:        ColumnRSLineNewHigh,
		DisplayName: "RS Line New High",
		Description: "An indicator that the relative strength line has reached a new 52-week high on an intraday basis.",
		Category:    catPriceVolume,
	},
	{
		Name: ColumnRSLineNewLow, DisplayName: "RS Line New Low",
		Description: "An indicator that the relative strength line has reached a new 52-week low on an intraday basis.",
		Category:    catPriceVolume,
	},
	{
		Name:        ColumnRSRating3M,
		DisplayName: "RS 3-Month Rating",
		Description: "Measures a stock’s price movement over the last 3 months vs. other stocks meeting minimum liquidity and market cap restrictions.",
		Category:    catPriceVolume,
	},
	{
		Name:        ColumnRSRating6M,
		DisplayName: "RS 6-Month Rating",
		Description: "Measures a stock’s price movement over the last 6 months vs. other stocks meeting minimum liquidity and market cap restrictions.",
		Category:    catPriceVolume,
	},
	{
		Name:        "Timeliness Rating",
		DisplayName: "Timeliness Rating",
		Description: "A proprietary rating based upon recent earnings changes and price peformance indicating possible or potential relative price performance over the next 12 months.",
		Category:    catPriceVolume,
	},
	{
		Name:        ColumnTrailing26WPctPerfomanceVsSP500,
		DisplayName: "Trl 26 Wk % Perf vs S&P 500",
		Description: "The difference between the stock’s price percent change and the trailing 26-week percent change of the S&P500 index.",
		Category:    catPriceVolume,
	},
	{
		Name: ColumnVolume, DisplayName: "Volume (1000s)",
		Description: "The volume traded on a 20-minute delayed basis for the current trading day.",
		Category:    catPriceVolume,
	},
	{
		Name:        ColumnWeeklyClosingRange,
		DisplayName: "Weekly Closing Range",
		Description: "The percentage that a stock is able to recover FROM the low of the week. This figure is provided on the Weekly timeframe.",
		Category:    catPriceVolume,
	},
	{
		Name:        ColumnDailyClosingRange,
		DisplayName: "Daily Closing Range",
		Description: "The percentage that a stock is able to recover FROM the low of the day. This figure is provided on the Daily timeframe.",
		Category:    catPriceVolume,
	},
	{
		Name: ColumnVolumeAvg50Day, DisplayName: "50-Day Avg Vol (1000s)",
		Description: "The average number of shares traded each day over the last 50 trading days.",
		Category:    catPriceVolume,
	},
	{
		Name: ColumnVolumePctChgVs50DAvgVolume, DisplayName: "Vol % Chg vs 50-Day",
		Description: "Comparison measure of the volume for that day against the 50-Day average volume.",
		Category:    catPriceVolume,
	},
	{
		Name:        ColumnUpDownVolumeRatio,
		DisplayName: "Up/Down Vol",
		Description: "A 50-day ratio that is derived by dividing total volume on up days by the total volume on down days.",
		Category:    catPriceVolume,
	},
	{
		Name: ColumnVolumeDollarAvg50D, DisplayName: "50-Day Avg $ Vol (1000s)",
		Description: "The 50-day average volume multiplied by the latest closing price.",
		Category:    catPriceVolume,
	},
	{
		Name: ColumnPricePctChgVs50DaySMA, DisplayName: "Price vs 50-Day",
		Description: "The amount, in percentage, a stock's price is above or below its 50-Day Moving Average Price.",
		Category:    catPriceVolume,
	},
	{
		Name: ColumnPricePctChgVs200DaySMA, DisplayName: "Price vs 200-Day",
		Description: "The amount, in percentage, a stock's price is above or below its 200-Day Moving Average Price",
		Category:    catPriceVolume,
	},
	{
		Name: ColumnPricePctChgCurrentWeek, DisplayName: "% Chg Cur Week",
		Description: "The amount a stock's price increased from last Friday's close to the intraday price.",
		Category:    catPriceVolume,
	},
	{
		Name: ColumnPricePctChgLast1M, DisplayName: "% Chg 1 Month",
		Description: "The amount a stock's price increased over the last 1 month.",
		Category:    catPriceVolume,
	},
	{
		Name: ColumnPricePctChgLast3M, DisplayName: "% Chg 3 Months",
		Description: "The amount a stock's price increased over the last 3 months.",
		Category:    catPriceVolume,
	},
	{
		Name: ColumnPricePctChgLast6M, DisplayName: "% Chg 6 Months",
		Description: "The amount a stock's price increase over the last 6 months.",
		Category:    catPriceVolume,
	},
	{
		Name: ColumnPricePctChgLast12M, DisplayName: "% Chg 12 Months",
		Description: "The amount a stock's price has increased over the last 12 months.",
		Category:    catPriceVolume,
	},
	{
		Name: ColumnPricePctChgYTD, DisplayName: displayPctChgYTD,
		Description: "The amount, in percentage, a stock's price has increased for the year to date.",
		Category:    catPriceVolume,
	},
	{
		Name:        ColumnAlpha,
		DisplayName: "Alpha",
		Description: "A measurement which uses a fixed S&P 500 value to determine how much a stock would have appreciated or depreciated over the past 12 months.",
		Category:    catPriceVolume,
	},
	{
		Name:        ColumnBeta,
		DisplayName: "Beta",
		Description: "Measures a stock's price volatility relative to price performance of the S&P 500 Index, over a 12-month period.",
		Category:    catPriceVolume,
	},
	{
		Name: ColumnATR30D, DisplayName: "Avg True Range",
		Description: "The average true range over the last 30 trading days.",
		Category:    catPriceVolume,
	},
	{
		Name: ColumnPricePercentChangeVs10DaySMA, DisplayName: "Price vs 10-Day",
		Description: "The amount, in percentage, a stock's price is above or below its 10-Day Moving Average Price.",
		Category:    catPriceVolume,
	},
	{
		Name: ColumnPricePercentChangeVs21DaySMA, DisplayName: "Price vs 21-Day",
		Description: "The amount, in percentage, a stock's price is above or below its 21-Day Moving Average Price.",
		Category:    catPriceVolume,
	},
	{
		Name: ColumnPricePercentChangeVs150DaySMA, DisplayName: "Price vs 150-Day",
		Description: "The amount, in percentage, a stock's price is above or below its 150-Day Moving Average Price.",
		Category:    catPriceVolume,
	},
	{
		Name: ColumnVolumePctChg10W, DisplayName: "Vol % Chg vs 10-Week",
		Description: "Comparison measure of the volume for that weekly against the 10-Week average volume.",
		Category:    catPriceVolume,
	},
	{
		Name: ColumnIsSMA10DAboveSMA21DAndSMA21DAboveSMA50, DisplayName: "10 Day > 21 Day > 50 Day",
		Description: "10 Day > 21 Day > 50 Day",
		Category:    catPriceVolume,
	},
	{
		Name: ColumnIsSMA50DayAboveSMA150DayAndSMA150DayAboveSMA200Day, DisplayName: "50-Day > 150-Day > 200-Day",
		Description: "50-Day > 150-Day > 200-Day",
		Category:    catPriceVolume,
	},
	{
		Name: ColumnVolumeUp5D, DisplayName: "Current day's Volume greater than previous 5 days' Volume",
		Description: "Current day's Volume greater than previous 5 days' Volume",
		Category:    catPriceVolume,
	},
	{
		Name: ColumnVolumeUp10D, DisplayName: "Current day's Volume greater than previous 10 days' Volume",
		Description: "Current day's Volume greater than previous 10 days' Volume",
		Category:    catPriceVolume,
	},
	{
		Name: ColumnVolumeUp20D, DisplayName: "Current day's Volume greater than previous 20 days' Volume",
		Description: "Current day's Volume greater than previous 20 days' Volume",
		Category:    catPriceVolume,
	},
	{
		Name:        ColumnRSLineWithin5PctOfHigh,
		DisplayName: "RS Line Within 5% of New High",
		Description: "An indicator that the relative strength line has neared the latest 52-week high(within 5%) on an intraday basis.",
		Category:    catPriceVolume,
	},
	{
		Name: ColumnATRPct21D, DisplayName: "21 Day ATR %",
		Description: "The average daily true percentage range over the past 21 days.",
		Category:    catPriceVolume,
	},
	{
		Name: ColumnATRPct30D, DisplayName: "30 Day ATR %",
		Description: "The average daily true percentage range over the past 30 days.",
		Category:    catPriceVolume,
	},
	{
		Name: ColumnATRPct50D, DisplayName: "50 Day ATR %",
		Description: "The average daily true percentage range over the past 50 days.",
		Category:    catPriceVolume,
	},

	// Margins & Ratios
	{
		Name:        ColumnProfitMarginAfterTaxLastReportedQ,
		DisplayName: "AT Margin",
		Description: "Determined by dividing quarterly income by quarterly sales, this figure should be among the highest in the industry.",
		Category:    catMarginsRatios,
	},
	{
		Name: ColumnProfitMarginAvgAfterTax2Q, DisplayName: "Avg AT Margin 2Q",
		Description: descAfterTaxMarginRatio,
		Category:    catMarginsRatios,
	},
	{
		Name: ColumnProfitMarginAvgAfterTax3Q, DisplayName: "Avg AT Margin 3Q",
		Description: "A financial performance ratio - calculated by dividing net income after taxes by net sales.",
		Category:    catMarginsRatios,
	},
	{
		Name: ColumnProfitMarginAvgAfterTax4Q, DisplayName: "Avg AT Margin 4Q",
		Description: descAfterTaxMarginRatio,
		Category:    catMarginsRatios,
	},
	{
		Name: ColumnProfitMarginAvgAfterTax5Q, DisplayName: "Avg AT Margin 5Q",
		Description: descAfterTaxMarginRatio,
		Category:    catMarginsRatios,
	},
	{
		Name: ColumnProfitMarginAvgAfterTax6Q, DisplayName: "Avg AT Margin 6Q",
		Description: descAfterTaxMarginRatio,
		Category:    catMarginsRatios,
	},
	{
		Name:        ColumnAfterTaxMarginAccelerationLast3Q,
		DisplayName: "AT Margin Accel",
		Description: "Checking the box limits the screen to return only those companies whose After Tax Profit Margin has been accelerating over the latest 3 quarters.",
		Category:    catMarginsRatios,
	},
	{
		Name:        ColumnProfitMarginPreTaxLastReportedY,
		DisplayName: "Pre-tax Margins",
		Description: "Expressed as a percentage, this figure divides operational income before taxes by fiscal year sales.",
		Category:    catMarginsRatios,
	},
	{
		Name:        ColumnOpMargGeqIndMedian,
		DisplayName: "Op Marg Geq Ind Median",
		Description: "Indicates that the company's operating margin is greater than or equal to the median operating margin for the company's industry.",
		Category:    catMarginsRatios,
	},
	{
		Name:        ColumnProfitMarginGeqIndustryMedian,
		DisplayName: "Prof Marg Geq Ind Median",
		Description: "Indicates that the company's profit margin is greater than or equal to the median profit margin for the company's industry.",
		Category:    catMarginsRatios,
	},
	{
		Name:        ColumnCashFlowVsEPSPctLastQ,
		DisplayName: "CF vs EPS % Last Qtr",
		Description: "Percentage difference between operating cash flow vs. earnings per share (EPS) from the last completed quarter.",
		Category:    catMarginsRatios,
	},
	{
		Name:        ColumnCashFlowVsEPSPctLastY,
		DisplayName: "CF vs EPS % Last Yr",
		Description: "Percentage difference between operating cash flow per share vs. earnings per share (EPS) from the last completed fiscal year.",
		Category:    catMarginsRatios,
	},
	{
		Name:        ColumnROE,
		DisplayName: "ROE",
		Description: "Annual income divided by average of the latest fiscal year and the prior year’s stockholder’s equity.",
		Category:    catMarginsRatios,
	},
	{
		Name: ColumnROEAvg5Y, DisplayName: "ROE 5-Yr Avg",
		Description: "The 5-year average return on equity.",
		Category:    catMarginsRatios,
	},
	{
		Name: ColumnDebtToEquityRatioLastY, DisplayName: "Debt %",
		Description: "Long term debt divided by shareholders’ equity.",
		Category:    catMarginsRatios,
	},
	{
		Name:        ColumnPriceEarningsRatioForward,
		DisplayName: "Forward P/E",
		Description: "A measure of the price-to-earnings ratio (P/E) using forecasted earnings for the P/E calculation.",
		Category:    catMarginsRatios,
	},
	{
		Name:        ColumnPriceEarningsRatio,
		DisplayName: "P/E",
		Description: "The sum of a company's price-to-earnings, calculated by taking the current stock price and dividing it by the trailing earnings per share for the past 12 months.",
		Category:    catMarginsRatios,
	},
	{
		Name: ColumnPEVsSP500PE, DisplayName: "P/E vs S&P 500 P/E (%)",
		Description: "Current P/E ratio of the stock divided by the average P/E of the stocks comprising the S&P 500.",
		Category:    catMarginsRatios,
	},
	{
		Name:        ColumnPEPercentileRank,
		DisplayName: "P/E Percent Rank",
		Description: "Stocks are arranged in order of highest P/E ratio and assigned to percentiles 99 (highest PEs) to 1 (lowest PEs).",
		Category:    catMarginsRatios,
	},
	{
		Name:        ColumnPEPercentileRankInGroup,
		DisplayName: "P/E Ratio Rank in Grp",
		Description: "The stocks in an industry group are arranged in order of trailing four-quarters Price to Earnings ratios (lowest to highest ) and assigned a 99 (best) to 1 (worst) rank.",
		Category:    catMarginsRatios,
	},
	{
		Name: ColumnPELessThan5YAvg, DisplayName: "P/E Lss 5-Yr Avg",
		Description: "Indicates that the company's current P/E ratio is lower than the five-year average P/E ratio.",
		Category:    catMarginsRatios,
	},
	{
		Name: ColumnPEGRatio, DisplayName: "PEG",
		Description: "Current P/E divided by five-year earnings growth rate.",
		Category:    catMarginsRatios,
	},
	{
		Name: ColumnPEGDividendAdjusted, DisplayName: "Dividend-Adjusted PEG",
		Description: "Current P/E divided by five-year earnings growth rate plus current dividend yield.",
		Category:    catMarginsRatios,
	},
	{
		Name: ColumnPriceToSales, DisplayName: "Price to Sales",
		Description: "Market Value divided by four quarter sales",
		Category:    catMarginsRatios,
	},
	{
		Name:        ColumnPriceToBookRatio,
		DisplayName: "Price to Book",
		Description: "Compares a stock's value in the market (determined by the current stock price) to the value of total company assets less total company liabilities (book value).",
		Category:    catMarginsRatios,
	},
	{
		Name: ColumnPriceToCashFlow, DisplayName: "Price to CF",
		Description: "Market Value divided by latest reported cash flow.",
		Category:    catMarginsRatios,
	},
	{
		Name: ColumnEnterpriseValueToFreeCashFlow, DisplayName: "EV to FCF",
		Description: "Company's current enterprise value divided by the company's current free cash flow.",
		Category:    catMarginsRatios,
	},
	{
		Name: ColumnCurrentRatio, DisplayName: "Current Ratio",
		Description: "Current assets divided by current liabilities.",
		Category:    catMarginsRatios,
	},
	{
		Name: ColumnYieldPct, DisplayName: "Yield %",
		Description: "The annualized dividend divided by Friday’s closing price.",
		Category:    catMarginsRatios,
	},
	{
		Name: ColumnLongTermDebtToWorkingCapitalRatio, DisplayName: "LT Debt to Working Cap",
		Description: "Total amount of long-term debt divided by the working capital ratio.",
		Category:    catMarginsRatios,
	},
	{
		Name:        ColumnLiabilitesToAssetsLessThanIndustryMedian,
		DisplayName: "Liab to Assets Lss Ind Median",
		Description: "Indicates that the company's total liabilities to total assets ratio is less than the industry median. This ratio is calculated by dividing the total liabilities by the total assets.",
		Category:    catMarginsRatios,
	},

	// General Stock Criteria
	{
		Name: ColumnExchange, DisplayName: "Exchange",
		Description: "The primary stock exchange on which the stock is traded.",
		Category:    catGeneralStockCriteria,
	},
	{
		Name: ColumnIsAdr, DisplayName: "ADR",
		Description: "Indicates that a stock is an American Depositary Receipt.",
		Category:    catGeneralStockCriteria,
	},
	{
		Name: ColumnIsETFOrClosedEndFund, DisplayName: "ETF/Closed - End Fund",
		Description: "Indicates that the security is an Exchange-Traded or Closed-End Fund.",
		Category:    catGeneralStockCriteria,
	},
	{
		Name: ColumnIsETF, DisplayName: "ETF",
		Description: "Indicates that the security is an Exchange-Traded Fund.",
		Category:    catGeneralStockCriteria,
	},
	{
		Name: "Options Exchange", DisplayName: "Options Exchange",
		Description: "The exchange on which options for a particular stock are traded.",
		Category:    catGeneralStockCriteria,
	},
	{
		Name: ColumnIPODate, DisplayName: "IPO Date",
		Description: "Search by year or over several years.",
		Category:    catGeneralStockCriteria,
	},
	{
		Name: ColumnIncorpDate, DisplayName: "Incorp Date",
		Description: "Year in which the company was incorporated.",
		Category:    catGeneralStockCriteria,
	},
	{
		Name: ColumnCity, DisplayName: "City",
		Description: "Screen stocks by geographic location, either domestic (U.S. and Canada) or international.",
		Category:    catGeneralStockCriteria,
	},
	{
		Name: ColumnCompanyDescription, DisplayName: "Company Description",
		Description: "Screens stocks by keywords contained in the company description.",
		Category:    catGeneralStockCriteria,
	},
	{
		Name: ColumnState, DisplayName: "State",
		Description: "Screen stocks by state or province headquarters.",
		Category:    catGeneralStockCriteria,
	},
	{
		Name: ColumnNewCEOLast12M, DisplayName: "New CEO 12 Months",
		Description: "New management can lead to increased growth.",
		Category:    catGeneralStockCriteria,
	},

	// IBD Content
	{
		Name: ColumnIBD50Flag, DisplayName: "IBD 50",
		Description: "Indicates that a stock is a component of the IBD 50.",
		Category:    catIBDContent,
	},
	{
		Name: ColumnIBDNewAmericaFlag, DisplayName: "IBD New America",
		Description: "Indicates that a stock is a component of the IBD New America Index.",
		Category:    catIBDContent,
	},
	{
		Name: ColumnIBD8585Flag, DisplayName: "IBD 85-85 Index",
		Description: "Indicates that a stock is a component of the IBD 85-85 Index.",
		Category:    catIBDContent,
	},
	{
		Name: ColumnIBDBigCap20Flag, DisplayName: "IBD Big Cap 20",
		Description: "Indicates that a stock is a component of the IBD Big Cap 20 Index.",
		Category:    catIBDContent,
	},

	// Miscellaneous
	{
		Name: ColumnAccDisRatingPreviousW, DisplayName: "A/D Rating - Pr Wk",
		Description: "Accumulation/Distribution Letter Rating - Prior Week",
		Category:    catMiscellaneous,
	},
	{
		Name: ColumnDividendAmountNextReported, DisplayName: "Expected X Dividend Amount",
		Description: "DIVIDEND, NEXT REPORTED AMOUNT",
		Category:    catMiscellaneous,
	},
	{
		Name: ColumnDividendDateNextReported, DisplayName: "Ex-Dividend Date",
		Description: "DIVIDEND, NEXT REPORTED EX-DIV DATE",
		Category:    catMiscellaneous,
	},
	{Name: ColumnShortVolume, DisplayName: "Short Volume", Description: "Short Vol", Category: catMiscellaneous},
	{
		Name: ColumnShortInterestPctOfFloat, DisplayName: "Shrt Int % of Float",
		Description: "SHORT INTEREST, % OF FLOAT",
		Category:    catMiscellaneous,
	},
	{
		Name: ColumnIndustryMarketValue, DisplayName: "Ind Mkt Val (bil)",
		Description: "Industry Group Market Value",
		Category:    catMiscellaneous,
	},
	{
		Name: ColumnIndustryGroupNumberOfNewHighs, DisplayName: "# New Highs in Group",
		Description: "Industry Group - Number of New High Stocks",
		Category:    catMiscellaneous,
	},
	{
		Name: ColumnIndustryGroupPctOfNewHighs, DisplayName: "% New Highs in Group",
		Description: "Industry Group - % New Highs in Group",
		Category:    catMiscellaneous,
	},
	{
		Name: ColumnIndustryGroupNumberOfNewLows, DisplayName: "# New Lows in Group",
		Description: "Industry Group - Number of New Low Stocks",
		Category:    catMiscellaneous,
	},
	{
		Name: ColumnIndustryGroupPctOfNewLows, DisplayName: "% New Lows in Group",
		Description: "Industry Group - % New Lows in Group",
		Category:    catMiscellaneous,
	},
	{
		Name: "Ind Grp % Chg 3 Mo", DisplayName: "Ind Grp % Chg 3 Mo",
		Description: "Industry Group % Chg 3 Months - Intraday",
		Category:    catMiscellaneous,
	},
	{
		Name: ColumnIndustryGroupRankLastW, DisplayName: "Ind Grp Rnk Last Week",
		Description: "Industry Group Rank, Last Week",
		Category:    catMiscellaneous,
	},
	{
		Name: ColumnIndustryGroupRank12WAgo, DisplayName: "Ind Grp Rnk 3 Mo Ago",
		Description: "Industry Group Rank, 3 Months Ago",
		Category:    catMiscellaneous,
	},
	{
		Name: ColumnIndustryGroupRank26WAgo, DisplayName: "Ind Grp Rnk 6 Mo Ago",
		Description: "Industry Group Rank, 6 Months Ago",
		Category:    catMiscellaneous,
	},
	{
		Name: "Index % Chg", DisplayName: "Index % Chg",
		Description: "Index Intraday % Chg Vs Last Close",
		Category:    catMiscellaneous,
	},
	{
		Name: "Index % Chg 13 Wks", DisplayName: "Index % Chg 13 Wks",
		Description: "Index % Chg vs 13 Wks - Intraday",
		Category:    catMiscellaneous,
	},
	{
		Name: "Index % Chg 52 Wks", DisplayName: "Index % Chg 52 Wks",
		Description: "Index % Chg vs 52 Wks - Intraday",
		Category:    catMiscellaneous,
	},
	{
		Name: "Index % Chg 26 Wks", DisplayName: "Index % Chg 26 Wks",
		Description: "Index % Chg vs 26 Wks - Intraday",
		Category:    catMiscellaneous,
	},
	{
		Name: "Index % Chg 4 Wks", DisplayName: "Index % Chg 4 Wks",
		Description: "Index % Chg vs 4 Wks - Intraday",
		Category:    catMiscellaneous,
	},
	{
		Name: ColumnIndexPctChg5D, DisplayName: "Index % Chg 5 Days",
		Description: "Index % Chg vs 5 Days - Intraday",
		Category:    catMiscellaneous,
	},
	{
		Name: "Index % Chg 8 Wks", DisplayName: "Index % Chg 8 Wks",
		Description: "Index % Chg vs 8 Wks - Intraday",
		Category:    catMiscellaneous,
	},
	{
		Name: displayIdxPctChgQTD, DisplayName: displayIdxPctChgQTD,
		Description: "Index % Chg vs QTD - Intraday",
		Category:    catMiscellaneous,
	},
	{
		Name: "Index % Chg YTD", DisplayName: "Index % Chg YTD",
		Description: "Index % Chg vs YTD - Intraday",
		Category:    catMiscellaneous,
	},
	{
		Name: ColumnAlertLastMarkUpDate, DisplayName: "Last Mark-up Date",
		Description: "Last Mark-up Date",
		Category:    catMiscellaneous,
	},
	{
		Name: ColumnCompaniesInIndustryGroup, DisplayName: "Number of Stocks",
		Description: "Industry Group Number Of Companies In Industry Group",
		Category:    catMiscellaneous,
	},
	{
		Name: displayPctChgYTD, DisplayName: displayPctChgYTD,
		Description: "Price, % Chg Year To Date - Intraday",
		Category:    catMiscellaneous,
	},
	{
		Name: ColumnPriceHighPreviousW, DisplayName: "Pr Wk High($)",
		Description: "Price, High Prior Week",
		Category:    catMiscellaneous,
	},
	{
		Name: ColumnVolumeShortDaysCurrent, DisplayName: "Days Vol Short Current",
		Description: "current",
		Category:    catMiscellaneous,
	},
	{
		Name: ColumnVolumeShortDays1PeriodAgo, DisplayName: "Days Vol Short 1 Period Ago",
		Description: "one month",
		Category:    catMiscellaneous,
	},
	{
		Name: ColumnVolumeShortDays2PeriodsAgo, DisplayName: "Days Vol Short 2 Periods Ago",
		Description: "two months",
		Category:    catMiscellaneous,
	},
	{
		Name: ColumnShortInterestPctChg, DisplayName: "Shrt Int % Chg",
		Description: "SHORT INTEREST, % CHG MO TO MO",
		Category:    catMiscellaneous,
	},
	{
		Name: "Timeliness Rating - Pr Wk", DisplayName: "Timeliness Rating - Pr Wk",
		Description: "Timeliness Rating - Prior Week",
		Category:    catMiscellaneous,
	},
	{
		Name:        ColumnPriceHigh52W,
		DisplayName: "52-Wk High",
		Description: "52-Week High Price",
		Category:    catMiscellaneous,
	},
	{Name: ColumnPriceLow52W, DisplayName: "52-Wk Low", Description: "52 Week Price Low", Category: catMiscellaneous},

	// Mutual Fund - General
	{
		Name:        ColumnFundRankVsAllFunds,
		DisplayName: "Rank vs All Funds",
		Description: "Proprietary ranking that measures the fund's cumulative total return compared to all other funds.",
		Category:    catGeneralFundCriteria,
	},
	{
		Name: ColumnFundNetAssets, DisplayName: "Net Assets (mil)",
		Description: "The total value of the fund's portfolio less liabilities.",
		Category:    catGeneralFundCriteria,
	},
	{
		Name:        ColumnFundObjective,
		DisplayName: "Objective",
		Description: "The investment objective of the fund is assigned based on the primary investment objective detailed in the prospectus.",
		Category:    catGeneralFundCriteria,
	},
	{
		Name: ColumnFundManagerStartYear, DisplayName: "Mngr Start Year",
		Description: "The year a new manager took over a fund.",
		Category:    catGeneralFundCriteria,
	},
	{
		Name: ColumnFundClosedToNewInvestors, DisplayName: "Closed New Investors",
		Description: "Select the checkbox to screen only for stocks that are closed to new investors.",
		Category:    catGeneralFundCriteria,
	},
	{
		Name: ColumnNAV, DisplayName: "NAV",
		Description: "The net asset value per share represents the fund's market price.",
		Category:    catGeneralFundCriteria,
	},
	{
		Name: ColumnNAVChg, DisplayName: "NAV Chg",
		Description: "The change in the NAV on one day from the NAV on the previous day.",
		Category:    catGeneralFundCriteria,
	},
	{
		Name: ColumnFundNavPctChg, DisplayName: "NAV % Chg",
		Description: "The percentage change in the NAV on one day from the NAV on the previous day.",
		Category:    catGeneralFundCriteria,
	},
	{Name: ColumnHoldingsPctFundAssetsHeld, DisplayName: "% of Fund", Category: catGeneralFundCriteria},
	{Name: ColumnSharesHeldPct, DisplayName: "% of Shares", Category: catGeneralFundCriteria},

	// Mutual Fund - Performance
	{
		Name: ColumnFundYTDReturn, DisplayName: "YTD Return",
		Description: "Includes reinvestment of any dividend, capital gain, or change in the fund’s NAV.",
		Category:    catPerformance,
	},
	{
		Name: ColumnFundTotalReturn1M, DisplayName: "1-Month Return",
		Description: "Total return for the last full month.",
		Category:    catPerformance,
	},
	{
		Name: ColumnFundTotalReturn3M, DisplayName: "3-Month Return",
		Description: "Total return for the last three months",
		Category:    catPerformance,
	},
	{
		Name: ColumnFundTotalReturn1Y, DisplayName: "1-Yr Return",
		Description: "Total return for the last year.",
		Category:    catPerformance,
	},
	{
		Name: ColumnFundTotalReturn3Y, DisplayName: "3-Yr Return",
		Description: "Cumulative total return for the last three years.",
		Category:    catPerformance,
	},
	{
		Name: ColumnFundTotalReturn5Y, DisplayName: "5-Yr Return",
		Description: "Cumulative total return for the last five years.",
		Category:    catPerformance,
	},
	{
		Name: ColumnFundTotalReturn10Y, DisplayName: "10-Yr Return",
		Description: "Cumulative total return for the last ten years.",
		Category:    catPerformance,
	},

	// Mutual Fund - Risk
	{
		Name:        ColumnFundStandardDeviation3Y,
		DisplayName: "Standard Dev-3 Yr",
		Description: "Standard deviation (the square root of the variance) is one of the best-known measures of risk. It is a statistical measure of the variability of returns around the mean total return. A larger standard deviation implies greater volatility or risk, all other factors being equal.",
		Category:    catRisk,
	},

	// Mutual Fund - Tax Efficiency
	{
		Name:        ColumnFundTurnoverPct,
		DisplayName: "Turnover",
		Description: "Measure of the fund’s trading activity, which is computed by taking the lesser of purchases or sales (excluding all securities with maturities of less than one year) and dividing by average monthly net assets.",
		Category:    catTaxEfficiency,
	},

	// Mutual Fund - Expenses
	{
		Name: ColumnFundFrontEndLoad, DisplayName: "Front End Load",
		Description: "The initial sales charge or front-end load is a deduction made from each investment in the fund.",
		Category:    catExpenses,
	},
	{
		Name:        ColumnFundExpenseRatio,
		DisplayName: "Expense Ratio",
		Description: "The cumulative representation of all of a fund's annual fund operating expenses, expressed as a percentage of the fund's average net assets.",
		Category:    catExpenses,
	},

	// Alerts
	{Name: ColumnAlertPrice, DisplayName: "Alert Price", Category: catAlerts},
	{Name: ColumnAlertType, DisplayName: "Alert Type", Category: catAlerts},
	{Name: ColumnAlertCreatedDate, DisplayName: "Created Date", Category: catAlerts},
	{Name: ColumnAlertUpSidePerc, DisplayName: "% to Upside", Category: catAlerts},
	{Name: ColumnAlertDownSidePerc, DisplayName: "% to Downside", Category: catAlerts},
	{Name: ColumnAlertTriggered, DisplayName: "Triggered Date", Category: catAlerts},
	{Name: ColumnChgPercVsAlertPrice, DisplayName: "% Change vs Alert", Category: catAlerts},

	// Other
	{Name: ColumnSymbol, DisplayName: ColumnSymbol, Category: ""},
	{Name: ColumnCompanyName, DisplayName: "Name", Category: ""},
	{Name: ColumnListRank, DisplayName: "#", Category: ""},
	{Name: ColumnEPSPlusRSRating, DisplayName: "#", Category: ""},
	{Name: "Add To List Date", DisplayName: "Add To List Date", Category: ""},
	{Name: ColumnBlueDotMostRecentDate, DisplayName: displayEventDate, Category: ""},
	{Name: ColumnBlueDotCount45Day, DisplayName: "Count (last 45 days)", Category: ""},
	{Name: ColumnRSLineMostRecentNewHighDate, DisplayName: displayEventDate, Category: ""},
	{Name: ColumnRSLineNewHighCount45Day, DisplayName: "Count (last 45 days)", Category: ""},
	{Name: Column50DayBreakOnVolumeMostRecentDate, DisplayName: displayEventDate, Category: ""},
	{Name: ColumnPullbackTo10WeekLineMostRecentDate, DisplayName: displayEventDate, Category: ""},
	{Name: ColumnEarningsGapDownPct, DisplayName: displayGap, Category: ""},
	{Name: ColumnAntEventMostRecentDate, DisplayName: displayEventDate, Category: ""},
	{Name: ColumnTightAreasMostRecentDate, DisplayName: displayEventDate, Category: ""},
	{Name: ColumnBasesFormingMostRecentDate, DisplayName: displayEventDate, Category: ""},
	{Name: ColumnPricePctChgVsPivotWeekly, DisplayName: "Weekly % Chg vs Pivot", Category: ""},
	{Name: ColumnAntsEventCount6M, DisplayName: "Count (last 6 months)", Category: ""},
	{Name: ColumnBreakawayGapPct, DisplayName: displayGap, Category: ""},
	{Name: ColumnBaseTypeWeekly, DisplayName: "Weekly Base Type", Category: ""},
	{Name: ColumnBaseStageWeekly, DisplayName: "Weekly Base Stage", Category: ""},
	{Name: ColumnWeeklyPivotWeek, DisplayName: "Weekly Pivot Week", Category: ""},
	{Name: ColumnPricePctChgToPivotWeekly, DisplayName: "Weekly % Chg to Pivot", Category: ""},
	{Name: ColumnBreakawayGapMostRecentDate, DisplayName: displayEventDate, Category: ""},
	{Name: ColumnIsBreakawayGapFilled, DisplayName: displayGapFilled, Category: ""},
	{Name: ColumnPricePctChgVsPivotDaily, DisplayName: "Daily % Chg vs Pivot", Category: ""},
	{Name: ColumnBaseTypeDaily, DisplayName: "Daily Base Type", Category: ""},
	{Name: ColumnBaseStageDaily, DisplayName: "Daily Base Stage", Category: ""},
	{Name: ColumnGapUpEarningsMostRecentDate, DisplayName: displayEventDate, Category: ""},
	{Name: ColumnIsGapUpEarningsGapFilled, DisplayName: displayGapFilled, Category: ""},
	{Name: ColumnGapDownEarningsMostRecentDate, DisplayName: displayEventDate, Category: ""},
	{Name: ColumnIsGapDownEarningsGapFilled, DisplayName: displayGapFilled, Category: ""},
	{Name: ColumnEarningsGapUpPct, DisplayName: displayGap, Category: ""},
	{Name: ColumnDatePriceHigh52W, DisplayName: "Date of New 52-Wk High (Pr Wk)", Category: ""},
	{Name: ColumnPricePctChgQTD, DisplayName: displayIdxPctChgQTD, Category: ""},
}
