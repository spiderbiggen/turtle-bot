package couch

import (
	"context"
	log "github.com/sirupsen/logrus"
	"strings"
)

func (c *Client) GetAverages(ctx context.Context, puuid string) (map[string]float64, error) {
	db, err := c.DB()
	if err != nil {
		return nil, err
	}
	rs := db.Query(
		ctx,
		"matches", "player_average",
		map[string]interface{}{
			"start_key": []interface{}{puuid, 0},
			"end_key":   []interface{}{puuid, map[string]interface{}{}},
			"reduce":    true,
			"group":     true,
		},
	)
	defer rs.Close()
	result := make(map[string]float64)
	for rs.Next() {
		var key []string
		var value float64
		_ = rs.ScanKey(&key)
		err := rs.ScanValue(&value)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		k := strings.Join(key[1:], ":")
		result[k] = value
	}
	if rs.Err() != nil {
		return nil, rs.Err()
	}
	return result, nil
}
