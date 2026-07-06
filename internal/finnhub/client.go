package finnhub

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Quote struct {
	CurrentPrice  float64 `json:"c"`
	Change        float64 `json:"d"`
	PercentChange float64 `json:"dp"`
	High          float64 `json:"h"`
	Low           float64 `json:"l"`
	Open          float64 `json:"o"`
	PreviousClose float64 `json:"pc"`
}

type Client struct {
	apiKey     string
	httpClient *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey:     apiKey,
		httpClient: &http.Client{},
	}
}

// takes a ticker and gets the stock info from finnhub
func (c *Client) GetQuote(ticker string) (*Quote, error) {
	endpoint := "https://finnhub.io/api/v1/quote"

	params := url.Values{}
	params.Add("symbol", ticker)
	params.Add("token", c.apiKey)

	resp, err := c.httpClient.Get(endpoint + "?" + params.Encode())
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var quote Quote
	if err := json.Unmarshal(body, &quote); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &quote, nil
}
