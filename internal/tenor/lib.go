package tenor

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/time/rate"
	"net/http"
	"net/url"
	"time"
)

var (
	ErrNoQuery           = errors.New("no search query")
	ErrInvalidStatus     = errors.New("invalid status code")
	ErrRateLimitExceeded = errors.New("rate limit exceeded")
)

type Client struct {
	key  string
	Http *http.Client
	Rate *rate.Limiter
}

func New(key string) *Client {
	return &Client{
		key:  key,
		Http: http.DefaultClient,
		Rate: rate.NewLimiter(rate.Every(1*time.Second), 1),
	}
}

type Response struct {
	Results ResultList `json:"results"`
	Next    string     `json:"next"`
}

type Result struct {
	ID    string `json:"id"`
	URL   string `json:"url"`
	Title string `json:"title"`
}

type ResultList []*Result

func (c *Client) Search(ctx context.Context, query string, opts ...SearchOpt) (ResultList, error) {
	if query == "" {
		return nil, ErrNoQuery
	}
	p := tenorSearchParameters{Query: query}
	for _, opt := range opts {
		opt(&p)
	}
	return c.request(ctx, p)
}

type SearchOpt func(query *tenorSearchParameters)

func WithLocale(locale string) func(query *tenorSearchParameters) {
	return func(query *tenorSearchParameters) {
		query.Locale = locale
	}
}

func WithNoContentFilter() func(query *tenorSearchParameters) {
	return func(query *tenorSearchParameters) {
		query.ContentFilter = off
	}
}

func WithLowContentFilter() func(query *tenorSearchParameters) {
	return func(query *tenorSearchParameters) {
		query.ContentFilter = low
	}
}

func WithMediumContentFilter() func(query *tenorSearchParameters) {
	return func(query *tenorSearchParameters) {
		query.ContentFilter = medium
	}
}

func WithHighContentFilter() func(query *tenorSearchParameters) {
	return func(query *tenorSearchParameters) {
		query.ContentFilter = high
	}
}

func WithMinimalMediaFilter() func(query *tenorSearchParameters) {
	return func(query *tenorSearchParameters) {
		query.MediaFilter = minimal
	}
}

func WithBasicMediaFilter() func(query *tenorSearchParameters) {
	return func(query *tenorSearchParameters) {
		query.MediaFilter = basic
	}
}

func WithLimit(limit uint8) func(query *tenorSearchParameters) {
	return func(query *tenorSearchParameters) {
		query.Limit = limit
	}
}

func WithPosition(pos uint8) func(query *tenorSearchParameters) {
	return func(query *tenorSearchParameters) {
		query.Position = pos
	}
}

func WithRandom(random bool) func(query *tenorSearchParameters) {
	return func(query *tenorSearchParameters) {
		query.Random = random
	}
}

type contentFilter string
type mediaFilter string

const (
	off    contentFilter = "off"
	low    contentFilter = "low"
	medium contentFilter = "medium"
	high   contentFilter = "high"
)

// deprecated
const (
	minimal mediaFilter = "minimal"
	basic   mediaFilter = "basic"
)

type tenorSearchParameters struct {
	Query         string
	Locale        string
	Random        bool
	ContentFilter contentFilter
	MediaFilter   mediaFilter
	Limit         uint8
	Position      uint8
}

func (c *Client) request(ctx context.Context, t tenorSearchParameters) (ResultList, error) {
	u, err := c.url(t)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	err = c.Rate.Wait(ctx)
	if err != nil {
		return nil, err
	}
	res, err := c.Http.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode > 299 {
		return nil, fmt.Errorf("%w: got %v", ErrInvalidStatus, res.Status)
	}
	var response Response
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response.Results, nil
}

func (c *Client) url(t tenorSearchParameters) (*url.URL, error) {
	url, err := url.Parse("https://tenor.googleapis.com/v2/search")
	if err != nil {
		return nil, err
	}
	q := url.Query()
	q.Set("key", c.key)
	q.Set("q", t.Query)
	q.Set("random", fmt.Sprintf("%t", t.Random))
	q.Set("pos", fmt.Sprintf("%d", t.Position))

	if t.Locale != "" {
		q.Set("locale", t.Locale)
	}

	if t.ContentFilter != "" {
		q.Set("contentfilter", string(t.ContentFilter))
	}

	q.Set("media_filter", "gif,gif_transparent")
	// TODO
	//if t.MediaFilter != "" {
	//	q.Set("media_filter", string(t.MediaFilter))
	//}

	if t.Limit != 0 {
		q.Set("limit", fmt.Sprintf("%d", t.Limit))
	}

	url.RawQuery = q.Encode()
	return url, nil
}
