package tenor

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

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

func Random(query string, opts ...Opt) (ResultList, error) {
	return newQuery(random, query, opts...).request()
}

func Search(query string, opts ...Opt) (ResultList, error) {
	return newQuery(search, query, opts...).request()
}

func Trending(opts ...Opt) (ResultList, error) {
	return newQuery(trending, "", opts...).request()
}

type Opt func(query *tenorQuery)

func WithLocale(locale string) func(query *tenorQuery) {
	return func(query *tenorQuery) {
		query.Locale = locale
	}
}

func WithNoContentFilter() func(query *tenorQuery) {
	return func(query *tenorQuery) {
		query.ContentFilter = off
	}
}

func WithLowContentFilter() func(query *tenorQuery) {
	return func(query *tenorQuery) {
		query.ContentFilter = low
	}
}

func WithMediumContentFilter() func(query *tenorQuery) {
	return func(query *tenorQuery) {
		query.ContentFilter = medium
	}
}

func WithHighContentFilter() func(query *tenorQuery) {
	return func(query *tenorQuery) {
		query.ContentFilter = high
	}
}

func WithMinimalMediaFilter() func(query *tenorQuery) {
	return func(query *tenorQuery) {
		query.MediaFilter = minimal
	}
}

func WithBasicMediaFilter() func(query *tenorQuery) {
	return func(query *tenorQuery) {
		query.MediaFilter = basic
	}
}

func WithLimit(limit uint8) func(query *tenorQuery) {
	return func(query *tenorQuery) {
		query.Limit = &limit
	}
}

func WithPosition(pos uint8) func(query *tenorQuery) {
	return func(query *tenorQuery) {
		query.Position = &pos
	}
}

type endpoint string
type contentFilter string
type mediaFilter string

const (
	search   endpoint = "search"
	random   endpoint = "random"
	trending endpoint = "trending"
)

const (
	off    contentFilter = "off"
	low    contentFilter = "low"
	medium contentFilter = "medium"
	high   contentFilter = "high"
)

const (
	minimal mediaFilter = "minimal"
	basic   mediaFilter = "basic"
)

type tenorQuery struct {
	Query         string
	Locale        string
	Endpoint      endpoint
	ContentFilter contentFilter
	MediaFilter   mediaFilter
	Limit         *uint8
	Position      *uint8
}

func (t tenorQuery) Url() (url *url.URL, err error) {
	sprintf := fmt.Sprintf("https://g.tenor.com/v1/%s", t.Endpoint)
	if url, err = url.Parse(sprintf); err != nil {
		return
	}
	q := url.Query()
	q.Set("key", os.Getenv("TENOR_KEY"))

	if t.Query != "" {
		q.Set("q", t.Query)
	}

	if t.Locale != "" {
		q.Set("locale", t.Locale)
	}

	if t.ContentFilter != "" {
		q.Set("contentfilter", string(t.ContentFilter))
	}

	if t.MediaFilter != "" {
		q.Set("media_filter", string(t.MediaFilter))
	} else {
		q.Set("media_filter", string(minimal))
	}

	if t.Limit != nil {
		q.Set("limit", strconv.FormatUint(uint64(*t.Limit), 10))
	}

	if t.Position != nil {
		q.Set("pos", strconv.FormatUint(uint64(*t.Position), 10))
	}
	url.RawQuery = q.Encode()
	return
}

func newQuery(e endpoint, q string, opts ...Opt) tenorQuery {
	t := tenorQuery{Endpoint: e, Query: q}
	for _, opt := range opts {
		opt(&t)
	}
	return t
}

func (t tenorQuery) request() (ResultList, error) {
	u, err := t.Url()
	if err != nil {
		return nil, err
	}
	res, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode > 299 {
		return nil, fmt.Errorf("invalid status: got %v", res.Status)
	}
	var response Response
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response.Results, nil
}
