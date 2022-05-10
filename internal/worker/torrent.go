package worker

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
	"weeb_bot/internal/kitsu"
	"weeb_bot/internal/nyaa"
)

var lastCheck = time.Now()

func NyaaCheck(ctx context.Context, discord *discordgo.Session) {
	checkTime := time.Now()
	episodes, err := nyaa.Episodes(ctx)
	if err != nil {
		log.Fatal(err)
	}
	c := map[string]string{
		"825808364649971712": "825808364649971715",
	}
	a := map[string][]string{
		"Tate no Yuusha no Nariagari S2":                    {"825808364649971712"},
		"Spy x Family":                                      {"825808364649971712"},
		"Gaikotsu Kishi-sama, Tadaima Isekai e Odekakechuu": {"825808364649971712"},
	}
	wg := sync.WaitGroup{}
	for _, group := range episodes {
		if group.FirstPublishedDate.Before(lastCheck) {
			continue
		}
		wg.Add(1)
		go func(ctx context.Context, group nyaa.Group, wg *sync.WaitGroup) {
			defer wg.Done()
			var embed *discordgo.MessageEmbed
			if guilds, ok := a[group.AnimeTitle]; ok {
				for _, guild := range guilds {
					if channel, ok := c[guild]; ok {
						if embed == nil {
							embed = makeEmbed(ctx, group)
						}
						_, err := discord.ChannelMessageSendEmbed(channel, embed)
						if err != nil {
							log.Errorln(err)
						}
					}
				}
			}
		}(ctx, group, &wg)
	}
	wg.Wait()
	lastCheck = checkTime
}

func coverImage(i ...*kitsu.ImageSet) string {
	for _, imageSet := range i {
		if imageSet == nil {
			continue
		}
		if imageSet.Medium != nil {
			return *imageSet.Medium
		} else {
			return imageSet.Original
		}
	}
	return ""
}

func makeEmbed(ctx context.Context, g nyaa.Group) *discordgo.MessageEmbed {
	title := g.AnimeTitle
	if g.Episode.Number != 0 {
		title = fmt.Sprintf("%s Ep %d", g.AnimeTitle, g.Episode.Number)
	}

	k := kitsu.New()
	var image *discordgo.MessageEmbedImage
	if r, _ := k.SearchAnime(ctx, g.AnimeTitle); len(r) > 0 {
		if imageUrl := coverImage(r[0].Cover, r[0].Poster); imageUrl != "" {
			image = &discordgo.MessageEmbedImage{URL: imageUrl}
		}
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
