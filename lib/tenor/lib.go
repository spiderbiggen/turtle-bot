package tenor

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

type tenorEndpoint string

const (
	search   tenorEndpoint = "search"
	random   tenorEndpoint = "random"
	trending tenorEndpoint = "trending"
)

func getTenorUrl(endpoint tenorEndpoint, query string) (url *url.URL, err error) {
	sprintf := fmt.Sprintf("https://g.tenor.com/v1/%s", endpoint)
	if url, err = url.Parse(sprintf); err != nil {
		return
	}
	q := url.Query()
	q.Set("key", os.Getenv("TENOR_KEY"))
	q.Set("q", query)
	q.Set("locale", "en")
	q.Set("contentfilter", "off")
	q.Set("media_filter", "minimal")
	q.Set("limit", "50")
	url.RawQuery = q.Encode()
	return
}

func request(u *url.URL) (results []Results, err error) {
	var resp *http.Response
	if resp, err = http.Get(u.String()); err != nil {
		return
	}
	var response Response
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return
	}
	results = response.Results
	return
}

func Random(query string) (results []Results, err error) {
	var u *url.URL
	if u, err = getTenorUrl(random, query); err != nil {
		return
	}
	results, err = request(u)
	return
}

func Top(query string) (results []Results, err error) {
	var u *url.URL
	if u, err = getTenorUrl(search, query); err != nil {
		return
	}
	results, err = request(u)
	return
}

type Response struct {
	Results []Results `json:"results"`
	Next    string    `json:"next"`
}

type Results struct {
	ID    string `json:"id"`
	URL   string `json:"url"`
	Title string `json:"title"`
}
