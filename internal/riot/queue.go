package riot

import (
	"context"
	"encoding/json"
	"net/http"
)

type Queue struct {
	ID          uint32 `json:"queueId"`
	Map         string `json:"map"`
	Description string `json:"description"`
	Notes       string `json:"notes"`
}

func Queues(ctx context.Context, httpClient *http.Client) ([]Queue, error) {
	url := "https://static.developer.riotgames.com/docs/lol/queues.json"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var q []Queue
	if err := json.NewDecoder(resp.Body).Decode(&q); err != nil {
		return nil, err
	}
	return q, nil
}

func (c *Client) Queues(ctx context.Context) ([]Queue, error) {
	return Queues(ctx, c.Http)
}
