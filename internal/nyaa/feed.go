package nyaa

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"
)

var (
	ErrHttpStatus = errors.New("invalid status code ")
	pattern       = regexp.MustCompile(`^\[.*?] (.*) - (\d+)(?:\.(\d+))?(?:[vV](\d+?))? \((\d+?p)\) \[.*?].mkv`)
	resolutions   = []string{"1080p", "720p", "540p", "480p"}
)

type rssFeed struct {
	XMLName xml.Name `xml:"rss"`
	Channel rssChannel
}

type rssChannel struct {
	XMLName     xml.Name  `xml:"channel"`
	Title       string    `xml:"title"`
	Description string    `xml:"description"`
	Items       []rssItem `xml:"item"`
}

type rssItem struct {
	XMLName xml.Name `xml:"item"`
	Title   string   `xml:"title"`
	Link    string   `xml:"link"`
	GUID    string   `xml:"guid"`
	PubDate *rssTime `xml:"pubDate"`
}

type rssTime struct {
	time.Time
}

func (c *rssTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	err := d.DecodeElement(&v, &start)
	if err != nil {
		return nil
	}
	parse, err := time.Parse(time.RFC1123Z, v)
	if err != nil {
		return err
	}
	*c = rssTime{parse.UTC()}
	return nil
}

type rssResult struct {
	Episode
	Download
}

type resultSet []*rssResult

func (s resultSet) Find(episode Episode) *rssResult {
	for _, result := range s {
		if result.Episode == episode {
			return result
		}
	}
	return nil
}

func (c *Client) getAnime(ctx context.Context, url *url.URL) (resultSet, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("%w: %d", ErrHttpStatus, resp.StatusCode)
	}
	return parseFeed(resp.Body)
}

func parseFeed(reader io.Reader) (resultSet, error) {
	var feed rssFeed
	if err := xml.NewDecoder(reader).Decode(&feed); err != nil {
		return nil, err
	}

	r := make([]*rssResult, 0, len(feed.Channel.Items))
	for _, item := range feed.Channel.Items {
		matches := pattern.FindStringSubmatch(item.Title)
		if matches == nil {
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
				Comments:   item.GUID,
				Resolution: matches[5],
				Torrent:    item.Link,
				FileName:   matches[0],
			},
		}
		if item.PubDate != nil {
			result.Download.PublishedDate = &item.PubDate.Time
		}
		r = append(r, &result)
	}
	return r, nil
}
