package viacep

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrCEPNotFound = errors.New("can not find zipcode")
)

type response struct {
	Localidade string          `json:"localidade"`
	Erro       json.RawMessage `json:"erro"`
}

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient() *Client {
	return &Client{
		baseURL:    "https://viacep.com.br/ws",
		httpClient: &http.Client{},
	}
}

func NewClientWithBaseURL(baseURL string) *Client {
	return &Client{
		baseURL:    baseURL,
		httpClient: &http.Client{},
	}
}

func (c *Client) GetCity(cep string) (string, error) {
	url := fmt.Sprintf("%s/%s/json/", c.baseURL, cep)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusBadRequest {
		return "", ErrCEPNotFound
	}

	var result response
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	errVal := string(result.Erro)
	if errVal == "true" || errVal == `"true"` {
		return "", ErrCEPNotFound
	}

	return result.Localidade, nil
}
