package stats

import (
	"encoding/json"
	"errors"
	"turtle-bot/internal/riot"
	"turtle-bot/internal/storage/models"
)

var (
	ErrInvalidParticipant = errors.New("invalid participant")
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

type ComparableResult struct {
	Current float64
	StatResult
}

type StatMap map[models.ChallengeType]StatResult
type ComparableMap map[models.ChallengeType]ComparableResult

func FindComparable(p *riot.Participant, s StatMap) (ComparableMap, error) {
	if p == nil {
		return nil, ErrInvalidParticipant
	}
	var pMap map[models.ChallengeType]interface{}
	b, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(b, &pMap); err != nil {
		return nil, err
	}

	challenges, ok := pMap["challenges"].(map[string]interface{})
	if !ok {
		return nil, ErrInvalidParticipant
	}
	results := make(ComparableMap)
	for k, hist := range s {
		cur, ok := challenges[string(k)].(float64)
		if !ok {
			cur, ok = pMap[k].(float64)
		}
		if ok {
			results[k] = ComparableResult{
				Current:    cur,
				StatResult: hist,
			}
		}
	}
	return results, nil
}

func (m ComparableMap) FilterMax() ComparableMap {
	r := make(ComparableMap)
	for k, cmp := range m {
		if cmp.Current > cmp.Max {
			r[k] = cmp
		}
	}
	return r
}

func (m ComparableMap) FilterMin() ComparableMap {
	r := make(ComparableMap)
	for k, cmp := range m {
		if cmp.Current < cmp.Min {
			r[k] = cmp
		}
	}
	return r
}

func (m ComparableMap) FilterAboveAverage() ComparableMap {
	r := make(ComparableMap)
	for k, cmp := range m {
		if cmp.Current > cmp.Average() {
			r[k] = cmp
		}
	}
	return r
}

func (m ComparableMap) FilterBelowAverage() ComparableMap {
	r := make(ComparableMap)
	for k, cmp := range m {
		if cmp.Current < cmp.Average() {
			r[k] = cmp
		}
	}
	return r
}
