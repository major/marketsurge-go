// Package marketsurge provides a Go client for the MarketSurge GraphQL API.
//
// The library keeps request and response handling strongly typed, so callers
// work with explicit structs instead of ad hoc maps or interface{} payloads.
// A Client is safe for concurrent use after construction.
//
// Example:
//
//	client, err := marketsurge.NewClient()
//	if err != nil {
//	    // handle error
//	}
//	resp, err := client.ChartMarketData(ctx, req)
//	_ = resp
//	_ = err
package marketsurge
