package worker

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
	kitsuApi "weeb_bot/internal/kitsu"
	"weeb_bot/internal/nyaa"
	"weeb_bot/internal/storage/postgres"
)

type nyaaWorker struct {
	db        *postgres.Client
	kitsu     *kitsuApi.Client
	lastCheck time.Time
}

func NyaaCheck(db *postgres.Client, kitsu *kitsuApi.Client) Worker {
	w := nyaaWorker{db: db, kitsu: kitsu, lastCheck: time.Now()}
	return func(ctx context.Context, s *discordgo.Session) {
		checkTime := time.Now()
		episodes, err := nyaa.Episodes(ctx)
		if err != nil {
			log.Errorf("Failed to get episodes from nyaa: %v", err)
			return
		}

		wg := sync.WaitGroup{}
		for _, group := range episodes {
			if group.FirstPublishedDate.Before(w.lastCheck) {
				continue
			}
			wg.Add(1)
			go w.sendToGuilds(ctx, s, group, &wg)
		}
		wg.Wait()
		w.lastCheck = checkTime
	}
}

func (w *nyaaWorker) sendToGuilds(ctx context.Context, s *discordgo.Session, group nyaa.Group, wg *sync.WaitGroup) {
	defer wg.Done()
	var embed *discordgo.MessageEmbed
	aSubs, err := w.db.GetSubscriptions(ctx, group.AnimeTitle)
	if err != nil {
		log.Warnf("Failed to get subscriptions: %v", err)
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

func (w *nyaaWorker) makeEmbed(g nyaa.Group, anime *postgres.Anime) *discordgo.MessageEmbed {
	title := anime.CanonicalTitle
	if g.Episode.Number != 0 {
		title = fmt.Sprintf("%s Ep %d", g.AnimeTitle, g.Episode.Number)
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
