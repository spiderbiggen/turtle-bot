package couch

import (
	"context"
	log "github.com/sirupsen/logrus"
	"math"
	"turtle-bot/internal/stats"
	"turtle-bot/internal/storage/models"
)

func (c *Client) GetAverages(ctx context.Context, puuid string) (stats.StatMap, error) {
	return c.GetQueueAverages(ctx, -1, puuid)
}

func (c *Client) GetQueueAverages(ctx context.Context, queueID int, puuID string) (stats.StatMap, error) {
	db, err := c.DB()
	if err != nil {
		return nil, err
	}

	var qMin, qMax interface{} = &queueID, &queueID
	if queueID < 0 {
		qMin = nil
		qMax = map[string]interface{}{}
	}
	rs := db.Query(
		ctx,
		"matches", "player_average",
		map[string]interface{}{
			"start_key": []interface{}{puuID, qMin, nil},
			"end_key":   []interface{}{puuID, qMax, map[string]interface{}{}},
			"reduce":    true,
			"group":     true,
		},
	)
	defer rs.Close()
	result := make(stats.StatMap)
	for rs.Next() {
		var key []interface{}
		var value stats.StatResult
		_ = rs.ScanKey(&key)
		err := rs.ScanValue(&value)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		if k, ok := key[len(key)-1].(string); ok {
			mKey := models.ChallengeType(k)
			if v, ok := result[mKey]; ok {
				result[mKey] = stats.StatResult{
					Sum:   v.Sum + value.Sum,
					Count: v.Count + value.Count,
					Min:   math.Min(v.Min, value.Min),
					Max:   math.Max(v.Max, value.Max),
				}
			} else {
				result[mKey] = value
			}
		}
	}
	if rs.Err() != nil {
		return nil, rs.Err()
	}
	return result, nil
}
