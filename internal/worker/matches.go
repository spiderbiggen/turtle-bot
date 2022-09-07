package worker

import (
	"context"
	"github.com/go-kivik/kivik/v4"
	log "github.com/sirupsen/logrus"
	"net/http"
	"sync/atomic"
	"time"
	"turtle-bot/internal/queue"
	"turtle-bot/internal/riot"
	"turtle-bot/internal/storage/couch"
)

type MatchWorker struct {
	api      *riot.Client
	db       *couch.Client
	started  atomic.Bool
	interval time.Duration
	ids      queue.Queue[string]
	cancel   func()
}

type MatchQueue interface {
	AddMatchIds(ids ...string)
	Start()
}

func NewMatchWorker(api *riot.Client, db *couch.Client) *MatchWorker {
	return &MatchWorker{
		api:      api,
		db:       db,
		interval: 15 * time.Second,
		ids:      queue.New[string](20),
	}
}

func (m *MatchWorker) Start() {
	if m.started.CompareAndSwap(false, true) {
		c := make(chan interface{})
		m.cancel = func() { c <- nil }
		go func() {
			ticker := time.NewTicker(m.interval)
			for {
				select {
				case <-ticker.C:
					m.getMatches()
				case <-c:
					m.cancel = nil
					m.started.Store(false)
					return
				}
			}
		}()
	}
}

func (m *MatchWorker) Cancel() {
	if m.cancel != nil {
		m.cancel()
	}
}

func (m *MatchWorker) AddMatchIds(ids ...string) {
	mIds := make([]*string, len(ids))
	for i, id := range ids {
		id2 := id
		mIds[i] = &id2
	}
	m.ids.AppendMany(mIds)
}

func (m *MatchWorker) getMatches() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	ids := m.ids.PopN(10)

	c := make(chan *riot.Match)
	for i, id := range ids {
		go func(i int, id string) {
			match, err := m.api.Match(ctx, riot.EUW1, id)
			if err != nil {
				log.Errorf("Failed to get match: %v", err)
				m.ids.Append(&id)
				c <- nil
				return
			}
			c <- match
		}(i, id)
	}
	for range ids {
		select {
		case match := <-c:
			if match != nil {
				ctx2, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
				if err := m.db.AddMatch(ctx2, match); err != nil {
					if kivik.StatusCode(err) == http.StatusConflict {
						log.Warnf("Tried to insert already existing match %s", match.Metadata.MatchID)
					} else {
						m.ids.Append(&match.Metadata.MatchID)
						log.Errorf("Failed to store match: %v", err)
					}
				}
				cancel()
			}
		}
	}

}
