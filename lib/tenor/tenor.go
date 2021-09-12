package tenor

import (
	"encoding/json"
	"fmt"
	"log"
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

func getTenorUrl(endpoint tenorEndpoint, query string) (*url.URL, error) {
	sprintf := fmt.Sprintf("https://g.tenor.com/v1/%s", endpoint)
	u, err := url.Parse(sprintf)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	q := u.Query()
	q.Set("key", os.Getenv("TENOR_KEY"))
	q.Set("q", query)
	q.Set("locale", "en")
	q.Set("contentfilter", "off")
	q.Set("media_filter", "minimal")
	q.Set("limit", "50")
	u.RawQuery = q.Encode()
	return u, nil
}

func Random(query string) []Results {
	u, err := getTenorUrl(random, query)
	resp, err := http.Get(u.String())
	if err != nil {
		log.Fatalln(err)
	}
	var response Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Fatal("ooopsss! an error occurred, please try again")
	}
	return response.Results
}

func Top(query string) []Results {
	u, err := getTenorUrl(search, query)

	resp, err := http.Get(u.String())
	if err != nil {
		log.Fatalln(err)
	}
	var response Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Fatal("ooopsss! an error occurred, please try again")
	}
	return response.Results
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
