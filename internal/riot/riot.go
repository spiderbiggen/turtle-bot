package riot

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
	"turtle-bot/internal/limiter"
)

var (
	ErrInvalidStatus = errors.New("invalid status code")
)

type Client struct {
	ApiKey  string
	Http    *http.Client
	Limiter limiter.Limiter
}

func New(apiKey string) *Client {
	intervalWindow := limiter.NewIntervalWindow(
		limiter.Limit{Count: 20, Interval: 1 * time.Second},
		limiter.Limit{Count: 100, Interval: 2 * time.Minute},
	)
	intervalWindow.StartCleanup()
	return &Client{
		ApiKey:  apiKey,
		Http:    http.DefaultClient,
		Limiter: intervalWindow,
	}
}

func (c *Client) request(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	req.Header.Set("X-Riot-Token", c.ApiKey)
	if batch, ok := ctx.Value("batch").(bool); !ok || !batch {
		err := c.Limiter.Wait(ctx)
		if err != nil {
			return nil, err
		}
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
