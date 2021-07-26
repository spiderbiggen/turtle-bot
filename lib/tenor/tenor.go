package tenor

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
)

func Random(query string) []Results {
	u, err := url.Parse("https://g.tenor.com/v1/random")
	if err != nil {
		log.Fatalln(err)
	}
	q := u.Query()
	// TODO
	q.Set("key", os.Getenv("TENOR_KEY"))
	q.Set("q", "sleep well")
	q.Set("locale", "en")
	q.Set("contentfilter", "off")
	q.Set("media_filter", "minimal")
	q.Set("limt", "50")
	u.RawQuery = q.Encode()

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
