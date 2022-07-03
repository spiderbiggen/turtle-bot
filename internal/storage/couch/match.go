package couch

import (
	"context"
	log "github.com/sirupsen/logrus"
	"weeb_bot/internal/riot"
)

func (c *Client) AddMatch(ctx context.Context, match *riot.Match) error {
	db, err := c.DB()
	if err != nil {
		return err
	}
	_, err = db.Put(ctx, match.Metadata.MatchID, match)
	return err
}

func (c *Client) FilterMatchIds(ctx context.Context, ids []string) ([]string, error) {
	db, err := c.DB()
	if err != nil {
		return nil, err
	}
	resultSet := db.Find(ctx, map[string]interface{}{
		"selector": map[string]interface{}{
			"_id": map[string]interface{}{"$in": ids},
		},
		"limit":  100,
		"fields": []string{"_id"},
	})
	defer resultSet.Close()
	if err != nil {
		log.Errorf("%#v", err)
		return nil, err
	}

	storedIds := make(map[string]interface{})
	for resultSet.Next() {
		var result struct {
			ID string `json:"_id"`
		}
		err := resultSet.ScanDoc(&result)
		if err != nil {
			log.Errorf("%v", err)
			return nil, err
		}
		storedIds[result.ID] = nil
	}
	if resultSet.Err() != nil {
		log.Errorf("%#v", resultSet.Err())
		return nil, resultSet.Err()
	}
	r := make([]string, 0, len(ids)-len(storedIds))
	for _, id := range ids {
		if _, ok := storedIds[id]; !ok {
			r = append(r, id)
		}
	}
	log.Printf("%v", storedIds)
	log.Printf("%v", r)

	return r, nil
}
