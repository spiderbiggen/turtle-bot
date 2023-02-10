package anime

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

var (
	ErrHttpError = errors.New("http error")
)

type Download struct {
	Comments      string    `json:"comments"`
	Resolution    string    `json:"resolution"`
	Torrent       string    `json:"torrent"`
	FileName      string    `json:"file_name"`
	PublishedDate time.Time `json:"published_date"`
}

type DownloadsResult struct {
	Title         string     `json:"title"`
	Episode       int        `json:"episode"`
	PublishedDate time.Time  `json:"published_date"`
	Downloads     []Download `json:"downloads"`
}

type Client struct {
	http *http.Client
}

func New() *Client {
	return &Client{
		http: http.DefaultClient,
	}
}

func (k *Client) SearchAnime(ctx context.Context, title string) ([]DownloadsResult, error) {
	u, err := url.Parse("https://api.spiderbiggen.com/anime/downloads")
	if err != nil {
		return nil, err
	}

	q := u.Query()
	if title != "" {
		q.Add("title", title)
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := k.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("%w, invalid status code: %d", ErrHttpError, resp.StatusCode)
	}

	var result []DownloadsResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}
