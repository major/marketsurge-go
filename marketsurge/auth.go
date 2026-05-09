package marketsurge

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const (
	jwtExchangePath                 = "/client"
	headerEncryptedDocumentKey      = "X-Encrypted-Document-Key"
	headerOriginalHost              = "X-Original-Host"
	headerOriginalReferrer          = "X-Original-Referrer"
	headerOriginalURL               = "X-Original-Url"
	marketsurgeOriginalHost         = "marketsurge.investors.com"
	marketsurgeOriginalURL          = "/mstool"
	operationDecodeJWTExchangeReply = "decode JWT exchange response"
)

// ClientInfoResponse holds the response from the JWT exchange endpoint.
type ClientInfoResponse struct {
	IsLoggedIn bool   `json:"isLoggedIn"`
	JWT        string `json:"jwt"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
}

// ExchangeJWT exchanges browser session cookies for a JWT token.
func (c *Client) ExchangeJWT(ctx context.Context, session Session) (*ClientInfoResponse, error) {
	if len(session.Cookies) == 0 {
		return nil, errors.New("marketsurge: ExchangeJWT: session has no cookies")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.investorsBaseURL+jwtExchangePath, nil)
	if err != nil {
		return nil, fmt.Errorf("marketsurge: create JWT exchange request: %w", err)
	}
	c.setJWTExchangeHeaders(req)
	for _, cookie := range session.Cookies {
		if cookie == nil {
			continue
		}
		req.AddCookie(cookie)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("marketsurge: send JWT exchange request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, c.statusError(resp)
	}

	responseBody, err := c.readJWTExchangeBody(resp.Body)
	if err != nil {
		return nil, err
	}

	var info ClientInfoResponse
	if decodeErr := json.Unmarshal(responseBody, &info); decodeErr != nil {
		return nil, &DecodeError{Operation: operationDecodeJWTExchangeReply, Err: decodeErr}
	}
	if !info.IsLoggedIn {
		return nil, errors.New(
			"marketsurge: ExchangeJWT: not logged in: check that you are signed into MarketSurge in the browser",
		)
	}
	if info.JWT == "" {
		return nil, errors.New("marketsurge: ExchangeJWT: JWT not found in exchange response")
	}

	return &info, nil
}

func (c *Client) setJWTExchangeHeaders(req *http.Request) {
	req.Header.Set(headerUserAgent, c.userAgent)
	req.Header.Set(headerEncryptedDocumentKey, "")
	req.Header.Set(headerOriginalHost, marketsurgeOriginalHost)
	req.Header.Set(headerOriginalReferrer, "")
	req.Header.Set(headerOriginalURL, marketsurgeOriginalURL)
	req.Header.Set(headerReferer, RefererURL)
	req.Header.Set(headerOrigin, OriginURL)
}

func (c *Client) readJWTExchangeBody(body io.Reader) ([]byte, error) {
	data, err := readBodyUpToLimit(body, c.responseBodyLimit+1)
	if err != nil {
		return nil, fmt.Errorf("marketsurge: read JWT exchange response body: %w", err)
	}
	if int64(len(data)) > c.responseBodyLimit {
		return nil, &BodyLimitError{Limit: c.responseBodyLimit}
	}
	return data, nil
}
