package nyaa

import (
	"context"
	"fmt"
	"github.com/mmcdole/gofeed"
	log "github.com/sirupsen/logrus"
	"net/url"
	"regexp"
	"strconv"
	"sync"
	"time"
)

var (
	pattern     = regexp.MustCompile(`^\[.*?] (.*) - (\d+)(?:\.(\d+))?(?:[vV](\d+?))? \((\d+?p)\) \[.*?].mkv`)
	resolutions = []string{"1080p", "720p", "540p", "480p"}
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

type rssResult struct {
	Episode
	Download
}

type resultSet []*rssResult

type resultSets []resultSet

func (s resultSet) Find(episode Episode) *rssResult {
	for _, result := range s {
		if result.Episode == episode {
			return result
		}
	}
	return nil
}

func GetAnime(ctx context.Context) ([]Group, error) {
	results := make(resultSets, len(resolutions), len(resolutions))
	wg := sync.WaitGroup{}
	for i, res := range resolutions {
		u, err := getUrl(res)
		if err != nil {
			return nil, err
		}
		wg.Add(1)
		go func(wg *sync.WaitGroup, index int, u *url.URL) {
			a, err := getAnime(ctx, u)
			if err != nil {
				log.Warn(err)
			}
			results[index] = a
			wg.Done()
		}(&wg, i, u)
	}
	wg.Wait()
	return results.Group(), nil
}

func (s resultSets) Group() []Group {
	m := make(map[Episode]*Group, len(s))
	for _, set := range s {
		for _, result := range set {
			if result == nil {
				fmt.Printf("Missing episode in %v\n", set)
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

func getAnime(ctx context.Context, url *url.URL) (resultSet, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURLWithContext(url.String(), ctx)
	if err != nil {
		return nil, err
	}
	r := make([]*rssResult, 0, len(feed.Items))
	for _, item := range feed.Items {
		matches := pattern.FindStringSubmatch(item.Title)
		if matches == nil {
			//log.Warnf("Failed to match %s", item.AnimeTitle)
			continue
		}

		var episode, decimal, version uint64
		episode, _ = strconv.ParseUint(matches[2], 10, 16)
		decimal, _ = strconv.ParseUint(matches[3], 10, 16)
		version, _ = strconv.ParseUint(matches[4], 10, 16)

		result := rssResult{
			Episode{
				AnimeTitle: matches[1],
				Number:     uint16(episode),
				Decimal:    uint16(decimal),
				Version:    uint16(version),
			},
			Download{
				Comments:      item.GUID,
				Resolution:    matches[5],
				Torrent:       item.Link,
				FileName:      matches[0],
				PublishedDate: item.PublishedParsed,
			},
		}
		r = append(r, &result)
	}
	return r, nil
}

func getUrl(resolution string) (*url.URL, error) {
	u, err := url.Parse("https://nyaa.si/")
	if err != nil {
		return nil, err
	}
	v := u.Query()
	v.Set("page", "rss")
	v.Set("c", "1_2")
	v.Set("q", fmt.Sprintf("[SubsPlease] %s", resolution))
	u.RawQuery = v.Encode()
	return u, nil
}
