package nyaa

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"time"
)

type Episode struct {
	AnimeTitle string
	Number     uint16
	Decimal    uint16
	Version    uint16
}

type Download struct {
	Comments      string
	Resolution    string
	Torrent       string
	FileName      string
	PublishedDate *time.Time
}

type Group struct {
	Episode
	FirstPublishedDate *time.Time
	Downloads          []Download
}

type Client struct {
	http *http.Client
}

func New() *Client {
	return &Client{http: http.DefaultClient}
}

func (c *Client) Episodes(ctx context.Context) ([]Group, error) {
	type episodesResult struct {
		Resolution string
		Results    resultSet
		Err        error
	}

	ch := make(chan episodesResult)
	m := make(map[string]resultSet)
	ctc, cancel := context.WithCancel(ctx)
	defer cancel()
	for _, res := range resolutions {
		go func(wg chan episodesResult, res string) {
			u, err := getUrl(res)
			if err != nil {
				ch <- episodesResult{Err: err}
				return
			}
			a, err := c.getAnime(ctc, u)
			if err != nil {

				ch <- episodesResult{Err: err}
				return
			}
			ch <- episodesResult{Resolution: res, Results: a}
		}(ch, res)
	}

	for range resolutions {
		select {
		case r := <-ch:
			if r.Err != nil {
				log.Warnf("Failed to get anime for %s", r.Err)
				continue
			}
			m[r.Resolution] = r.Results
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	return groupResults(m), nil
}

func groupResults(s map[string]resultSet) []Group {
	m := make(map[Episode]*Group, len(s))
	for _, set := range s {
		for _, result := range set {
			if result == nil {
				log.Warnf("Missing episode in %v", set)
				continue
			}
			if a, ok := m[result.Episode]; ok {
				a.Downloads = append(a.Downloads, result.Download)
			} else {
				m[result.Episode] = &Group{
					Episode:            result.Episode,
					FirstPublishedDate: result.PublishedDate,
					Downloads:          []Download{result.Download},
				}
			}
		}
	}
	r := make([]Group, 0, len(m))
	for _, group := range m {
		r = append(r, *group)
	}

	return r
}

func getUrl(resolution string) (*url.URL, error) {
	u, err := url.Parse("https://nyaa.si/")
	if err != nil {
		return nil, err
	}
	v := u.Query()
	v.Set("page", "rss")
	v.Set("c", "1_2")
	v.Set("q", fmt.Sprintf("[SubsPlease] (%s)", resolution))
	u.RawQuery = v.Encode()
	return u, nil
}
