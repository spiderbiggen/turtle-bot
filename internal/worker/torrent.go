package worker

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
	animeApi "turtle-bot/internal/anime"
	kitsuApi "turtle-bot/internal/kitsu"
	"turtle-bot/internal/storage/postgres"
)

type TorrentWorker struct {
	db        *postgres.Client
	kitsu     *kitsuApi.Client
	anime     *animeApi.Client
	lastCheck time.Time
	entryID   cron.EntryID
}

func NewTorrent(db *postgres.Client, kitsu *kitsuApi.Client, anime *animeApi.Client) TorrentWorker {
	return TorrentWorker{db: db, kitsu: kitsu, anime: anime, lastCheck: time.Now()}
}

func (w *TorrentWorker) Schedule(cron *cron.Cron, session *discordgo.Session) (err error) {
	if w.entryID == 0 {
		now := time.Now()
		startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		w.lastCheck = startOfDay
		w.entryID, err = cron.AddFunc("1-59/5 * * * *", func() {
			timeout, cancelFunc := context.WithTimeout(context.Background(), 1*time.Minute)
			defer cancelFunc()
			if err := w.Run(timeout, session); err != nil {
				log.Error(err)
			}
		})
		entry := cron.Entry(w.entryID)
		log.Debugf("Scheduled Torrent Worker with id: %d, first run at: %s", w.entryID, entry.Next)
	}
	return
}

func (w *TorrentWorker) Run(ctx context.Context, session *discordgo.Session) error {
	results, err := w.anime.SearchAnime(ctx, "")
	if err != nil {
		return fmt.Errorf("failed to get episodes from nyaa: %v", err)
	}
	results = w.filterAnime(results)
	log.Debugf("found %d new episodes", len(results))
	if len(results) == 0 {
		return nil
	}

	wg := sync.WaitGroup{}
	var checkTime time.Time
	for _, group := range results {
		for _, d := range group.Downloads {
			if d.PublishedDate.After(checkTime) {
				checkTime = d.PublishedDate
			}
		}
		wg.Add(1)
		go w.sendToGuilds(ctx, session, group, &wg)
	}
	wg.Wait()
	w.lastCheck = checkTime
	return nil
}

func (w *TorrentWorker) filterAnime(results []animeApi.DownloadsResult) []animeApi.DownloadsResult {
	var filtered []animeApi.DownloadsResult
	for _, result := range results {
		download, found := findHdDownload(result)
		if !found || !download.PublishedDate.After(w.lastCheck) {
			continue
		}
		filtered = append(filtered, result)
	}
	return filtered
}

func findHdDownload(anime animeApi.DownloadsResult) (animeApi.Download, bool) {
	for _, download := range anime.Downloads {
		if download.Resolution == "1080p" || download.Resolution == "2160p" {
			return download, true
		}
	}
	return animeApi.Download{}, false
}

func (w *TorrentWorker) sendToGuilds(ctx context.Context, s *discordgo.Session, group animeApi.DownloadsResult, wg *sync.WaitGroup) {
	defer wg.Done()
	aSubs, err := w.db.GetSubscriptions(ctx, group.Title)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Debugf("found no subscriptions for %v", group.Title)
			return
		}
		log.Errorf("Failed to get subscriptions: %v", err)
		return
	}

	log.Debugf("found %d subscriptions for %s", len(aSubs.Subscriptions), group.Title)
	embed := w.makeEmbed(group, aSubs.Anime)
	for _, sub := range aSubs.Subscriptions {
		if _, err := s.ChannelMessageSendEmbed(sub.ChannelID, &embed); err != nil {
			log.Errorf("Failed to send download embed: %v", err)
		}
	}
}

func (w *TorrentWorker) makeEmbed(g animeApi.DownloadsResult, anime postgres.Anime) discordgo.MessageEmbed {
	title := g.Title
	if g.Episode != 0 {
		title = fmt.Sprintf("%s Ep %d", title, g.Episode)
	}

	var image *discordgo.MessageEmbedImage
	if anime.ImageURL != "" {
		image = &discordgo.MessageEmbedImage{URL: anime.ImageURL}
	}

	fields := make([]*discordgo.MessageEmbedField, 0, len(g.Downloads))
	for _, d := range g.Downloads {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   d.Resolution,
			Value:  fmt.Sprintf("[torrent](%s)\n[comments](%s)", d.Torrent, d.Comments),
			Inline: true,
		})
	}

	return discordgo.MessageEmbed{
		Type:   discordgo.EmbedTypeRich,
		Title:  title,
		Fields: fields,
		Image:  image,
	}
}
