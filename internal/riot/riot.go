package riot

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/time/rate"
	"net/http"
)

var (
	ErrInvalidStatus     = errors.New("invalid status code")
	ErrRateLimitExceeded = errors.New("rate limit exceeded")
)

type Client struct {
	ApiKey string
	Http   *http.Client
	Rate   *rate.Limiter
}

func New(apiKey string) *Client {
	return &Client{
		ApiKey: apiKey,
		Http:   http.DefaultClient,
		Rate:   rate.NewLimiter(rate.Limit(20), 20),
	}
}

func (c *Client) request(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	req.Header.Set("X-Riot-Token", c.ApiKey)
	err = c.Rate.Wait(ctx)
	if err != nil {
		return nil, ErrRateLimitExceeded
	}
	resp, err := c.Http.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return resp, fmt.Errorf("%w: %d", ErrInvalidStatus, resp.StatusCode)
	}
	return resp, nil
}
