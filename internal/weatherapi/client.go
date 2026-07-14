package weatherapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
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

func normalizeLocation(location string) string {
	replacer := strings.NewReplacer(
		"á", "a", "à", "a", "â", "a", "ä", "a", "ã", "a", "å", "a",
		"ç", "c",
		"é", "e", "è", "e", "ê", "e", "ë", "e",
		"í", "i", "ì", "i", "î", "i", "ï", "i",
		"ó", "o", "ò", "o", "ô", "o", "ö", "o", "õ", "o", "ø", "o",
		"ú", "u", "ù", "u", "û", "u", "ü", "u",
		"ý", "y", "ÿ", "y",
		"ñ", "n",
		"ß", "ss",
		"æ", "ae",
		"œ", "oe",
	)

	return strings.Join(strings.Fields(replacer.Replace(strings.TrimSpace(location))), " ")
}

func (c *Client) GetCurrentTemp(city string) (float64, error) {
	normalizedCity := normalizeLocation(city)
	query := normalizedCity
	if strings.Contains(normalizedCity, ",") {
		query = normalizedCity
	} else {
		parts := strings.Fields(normalizedCity)
		if len(parts) > 1 {
			query = strings.Join(parts[:len(parts)-1], " ") + ", " + parts[len(parts)-1]
		}
	}
	reqURL := fmt.Sprintf("%s/current.json?key=%s&q=%s&aqi=no",
		c.baseURL, c.apiKey, url.QueryEscape(query))

	resp, err := c.httpClient.Get(reqURL)
	if err != nil {
		fmt.Printf("weatherapi request failed: url=%s err=%v\n", reqURL, err)
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("weatherapi returned non-200: url=%s status=%d\n", reqURL, resp.StatusCode)
		return 0, fmt.Errorf("weather API returned status %d", resp.StatusCode)
	}

	var result response
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Printf("weatherapi decode failed: url=%s err=%v\n", reqURL, err)
		return 0, err
	}

	return result.Current.TempC, nil
}
