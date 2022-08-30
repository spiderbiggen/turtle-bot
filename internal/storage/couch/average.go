package couch

import (
	"context"
	log "github.com/sirupsen/logrus"
	"math"
)

type StatResult struct {
	Sum   float64 `json:"sum"`
	Count int     `json:"count"`
	Min   float64 `json:"min"`
	Max   float64 `json:"max"`
}

func (r *StatResult) Average() float64 {
	return r.Sum / float64(r.Count)
}

func (c *Client) GetAverages(ctx context.Context, puuid string) (map[string]StatResult, error) {
	return c.GetQueueAverages(ctx, -1, puuid)
}

func (c *Client) GetQueueAverages(ctx context.Context, queueID int, puuID string) (map[string]StatResult, error) {
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
	result := make(map[string]StatResult)
	for rs.Next() {
		var key []interface{}
		var value StatResult
		_ = rs.ScanKey(&key)
		err := rs.ScanValue(&value)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		if k, ok := key[len(key)-1].(string); ok {
			if v, ok := result[k]; ok {
				result[k] = StatResult{
					Sum:   v.Sum + value.Sum,
					Count: v.Count + value.Count,
					Min:   math.Min(v.Min, value.Min),
					Max:   math.Max(v.Max, value.Max),
				}
			} else {
				result[k] = value
			}
		}
	}
	if rs.Err() != nil {
		return nil, rs.Err()
	}
	return result, nil
}
