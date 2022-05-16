package riot

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrInvalidStatus = errors.New("invalid status code")
)

type Client struct {
	ApiKey string
	Http   *http.Client
}

func New(apiKey string) *Client {
	return &Client{
		ApiKey: apiKey,
		Http:   http.DefaultClient,
	}
}

func (c *Client) request(context context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(context, http.MethodGet, url, nil)

	req.Header.Set("X-Riot-Token", c.ApiKey)
	resp, err := c.Http.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return resp, fmt.Errorf("%w: %d", ErrInvalidStatus, resp.StatusCode)
	}
	return resp, nil
}
