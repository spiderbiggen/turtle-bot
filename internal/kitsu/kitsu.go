package kitsu

import (
	"context"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	jsonApiType       = "application/vnd.api+json"
	acceptHeader      = "Accept"
	contentTypeHeader = "Content-Type"
)

type Anime struct {
	ID             string     `json:"id,omitempty"`
	CreatedAt      *time.Time `json:"created_at,omitempty"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty"`
	Slug           string     `json:"slug,omitempty"`
	CanonicalTitle string     `json:"canonical_title,omitempty"`
	Synopsis       string     `json:"synopsis,omitempty"`
	Description    string     `json:"description,omitempty"`
	Cover          *ImageSet  `json:"cover,omitempty"`
	Poster         *ImageSet  `json:"poster,omitempty"`
}

type ImageSet struct {
	Tiny     *string `json:"tiny,omitempty"`
	Small    *string `json:"small,omitempty"`
	Medium   *string `json:"medium,omitempty"`
	Large    *string `json:"large,omitempty"`
	Original string  `json:"original,omitempty"`
}

type animeCollectionResult struct {
	Data []*struct {
		ID         string `json:"id,omitempty"`
		Attributes struct {
			CreatedAt      *time.Time `json:"created_at,omitempty"`
			UpdatedAt      *time.Time `json:"updated_at,omitempty"`
			Slug           string     `json:"slug,omitempty"`
			CanonicalTitle string     `json:"canonical_title,omitempty"`
			Synopsis       string     `json:"synopsis,omitempty"`
			Description    string     `json:"description,omitempty"`
			Cover          *ImageSet  `json:"coverImage,omitempty"`
			Poster         *ImageSet  `json:"posterImage,omitempty"`
		} `json:"attributes"`
	} `json:"data,omitempty"`
	Meta struct {
		Count int `json:"count,omitempty" json:"count,omitempty"`
	} `json:"meta"`
	Links struct {
		First string `json:"first,omitempty" json:"first,omitempty"`
		Next  string `json:"next,omitempty" json:"next,omitempty"`
		Last  string `json:"last,omitempty" json:"last,omitempty"`
	} `json:"links"`
}

type Kitsu struct {
	Client *http.Client
}

func New() Kitsu {
	return Kitsu{
		Client: http.DefaultClient,
	}
}

func (k *Kitsu) doRequest(req *http.Request) (*http.Response, error) {
	req.Header.Set(acceptHeader, jsonApiType)
	req.Header.Set(contentTypeHeader, jsonApiType)
	return k.Client.Do(req)
}

func (k *Kitsu) SearchAnime(ctx context.Context, title string) ([]*Anime, error) {
	u, err := url.Parse("https://kitsu.io/api/edge/anime")
	if err != nil {
		return nil, err
	}

	q := u.Query()
	if title != "" {
		q.Add("filter[text]", title)
	}
	q.Add("page[limit]", "5")
	q.Add("page[offset]", "0")
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := k.doRequest(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Warnf("Failed to close body. %v", err)
		}
	}(resp.Body)

	var result animeCollectionResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	anime := make([]*Anime, 0, len(result.Data))
	for _, d := range result.Data {
		anime = append(anime, &Anime{
			ID:             d.ID,
			CreatedAt:      d.Attributes.CreatedAt,
			UpdatedAt:      d.Attributes.UpdatedAt,
			Slug:           d.Attributes.Slug,
			CanonicalTitle: d.Attributes.CanonicalTitle,
			Synopsis:       d.Attributes.Synopsis,
			Description:    d.Attributes.Description,
			Cover:          d.Attributes.Cover,
			Poster:         d.Attributes.Poster,
		})
	}

	return anime, nil
}
