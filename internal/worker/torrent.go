package worker

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
	"turtle-bot/internal/anime"
	kitsuApi "turtle-bot/internal/kitsu"
	"turtle-bot/internal/storage/postgres"
)

type nyaaWorker struct {
	db        *postgres.Client
	kitsu     *kitsuApi.Client
	lastCheck time.Time
}

func NyaaCheck(db *postgres.Client, kitsu *kitsuApi.Client, anime *anime.Client, startTimes ...time.Time) Worker {
	w := nyaaWorker{db: db, kitsu: kitsu, lastCheck: time.Now()}
	if len(startTimes) > 0 {
		w.lastCheck = startTimes[0]
	}

	return func(ctx context.Context, s *discordgo.Session) {
		var checkTime time.Time
		episodes, err := anime.SearchAnime(ctx, "")
		if err != nil {
			log.Errorf("Failed to get episodes from nyaa: %v", err)
			return
		}

		wg := sync.WaitGroup{}
		for _, group := range episodes {
			for _, d := range group.Downloads {
				if d.PublishedDate.After(checkTime) {
					checkTime = d.PublishedDate
				}
			}
			d := group.Downloads[0]
			if d.Resolution == "1080p" && d.PublishedDate.Before(w.lastCheck) {
				continue
			}
			wg.Add(1)
			go w.sendToGuilds(ctx, s, group, &wg)
		}
		wg.Wait()
		w.lastCheck = checkTime
	}
}

func (w *nyaaWorker) sendToGuilds(ctx context.Context, s *discordgo.Session, group anime.DownloadsResult, wg *sync.WaitGroup) {
	defer wg.Done()
	var embed *discordgo.MessageEmbed
	aSubs, err := w.db.GetSubscriptions(ctx, group.Title)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return
		}
		log.Errorf("Failed to get subscriptions: %v", err)
		return
	}
	for _, sub := range aSubs.Subs {
		if embed == nil {
			embed = w.makeEmbed(group, aSubs.Anime)
		}
		_, err := s.ChannelMessageSendEmbed(sub.ChannelID, embed)
		if err != nil {
			log.Errorf("Failed to send download embed: %v", err)
		}
	}
}

func (w *nyaaWorker) makeEmbed(g anime.DownloadsResult, anime *postgres.Anime) *discordgo.MessageEmbed {
	title := anime.CanonicalTitle
	if g.Episode != 0 {
		title = fmt.Sprintf("%s Ep %d", g.Title, g.Episode)
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

	return &discordgo.MessageEmbed{
		Type:   discordgo.EmbedTypeRich,
		Title:  title,
		Fields: fields,
		Image:  image,
	}
}
