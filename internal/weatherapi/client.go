package weatherapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type response struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey:     apiKey,
		baseURL:    "http://api.weatherapi.com/v1",
		httpClient: &http.Client{},
	}
}

func NewClientWithBaseURL(baseURL, apiKey string) *Client {
	return &Client{
		apiKey:     apiKey,
		baseURL:    baseURL,
		httpClient: &http.Client{},
	}
}

func (c *Client) GetCurrentTemp(city string) (float64, error) {
	reqURL := fmt.Sprintf("%s/current.json?key=%s&q=%s&aqi=no",
		c.baseURL, c.apiKey, url.QueryEscape(city))

	resp, err := c.httpClient.Get(reqURL)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("weather API returned status %d", resp.StatusCode)
	}

	var result response
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	return result.Current.TempC, nil
}
