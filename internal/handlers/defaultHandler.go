package handlers

import (
	"context"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"ytdlp-bot/internal/api"
	"ytdlp-bot/internal/data"
	"ytdlp-bot/internal/security"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func DefaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {

	if !security.CheckAuth(update.Message.From.ID) {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Пользователь с id " + strconv.FormatInt(update.Message.From.ID, 10) + " не авторизован",
		})
		return
	}

	if update.Message.From.ID == b.ID() {
		return
	}

	if update.Message.Text == "/start" {
		b.DeleteMessage(ctx, &bot.DeleteMessageParams{
			ChatID:    update.Message.Chat.ID,
			MessageID: update.Message.ID,
		})
		return
	}

	entities := getEntitiesFromLinkFacade(update.Message.Text)

	if len(entities) == 0 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Неизвестная ссылка",
		})
	} else {
		data.PutEntitiesToDB(entities)
		b.DeleteMessage(ctx, &bot.DeleteMessageParams{
			ChatID:    update.Message.Chat.ID,
			MessageID: update.Message.ID,
		})
	}

}

//yt-dlp -t mp4 -o "%(title)s.%(ext)s" -N 4 --cookies-from-browser chrome -r 40M --no-part  If3YIcLZ7sc

func getEntitiesFromLinkFacade(link string) []data.DatabaseEntity {
	if strings.Contains(link, "/watch") {
		return getVideoEntityFromLinkToVideo(link)
	} else if strings.Contains(link, "/playlist") {
		return getVideoEntitiesFromLinkToPlaylist(link)
	} else if strings.Contains(link, "//youtu.be") {
		return getVideoEntityShortLinkToVideo(link)
	} else if strings.Contains(link, "/live") {
		return getVideoEntityShortLinkToVideo(link)
	} else if strings.Contains(link, "/shorts") {
		return getShortEntityFromLinkToVideo(link)
	}
	return nil
}

func getShortEntityFromLinkToVideo(link string) []data.DatabaseEntity {
	u, err := url.Parse(link)
	if err != nil {
		return nil
	}
	pathArr := strings.Split(u.Path, "/")
	videoId := pathArr[len(pathArr)-1]
	channelTitle, videoTitle := api.GetChannelTitleAndVideoTitleFromVideoID(videoId)

	short := data.DatabaseEntity{
		VideoID:       videoId,
		VideoTitle:    videoTitleFilter(videoTitle),
		ChannelTitle:  channelTitle,
		PlaylistTitle: "Shorts",
		InWork:        0,
	}
	return []data.DatabaseEntity{short}
}

func getVideoEntityShortLinkToVideo(link string) []data.DatabaseEntity {
	u, err := url.Parse(link)
	if err != nil {
		return nil
	}
	pathArr := strings.Split(u.Path, "/")
	videoId := pathArr[len(pathArr)-1]

	channelTitle, videoTitle := api.GetChannelTitleAndVideoTitleFromVideoID(videoId)

	entity := data.DatabaseEntity{
		VideoID:       videoId,
		VideoTitle:    videoTitleFilter(videoTitle),
		ChannelTitle:  channelTitle,
		PlaylistTitle: "",
		InWork:        0,
	}
	return []data.DatabaseEntity{entity}
}

func getVideoEntityFromLinkToVideo(link string) []data.DatabaseEntity {
	u, err := url.Parse(link)
	if err != nil {
		return nil
	}
	querry, _ := url.ParseQuery(u.RawQuery)
	videoId := querry["v"]

	channelTitle, videoTitle := api.GetChannelTitleAndVideoTitleFromVideoID(videoId[0])
	playlistTitle := ""

	playlistId, ok := querry["list"]
	if ok {
		playlistTitle = api.GetPlaylistTitleFromListID(playlistId[0])
	}

	entity := data.DatabaseEntity{
		VideoID:       videoId[0],
		VideoTitle:    videoTitleFilter(videoTitle),
		ChannelTitle:  channelTitle,
		PlaylistTitle: playlistTitle,
		InWork:        0,
	}
	return []data.DatabaseEntity{entity}
}

func getVideoEntitiesFromLinkToPlaylist(link string) []data.DatabaseEntity {
	u, err := url.Parse(link)
	if err != nil {
		return nil
	}
	querry, _ := url.ParseQuery(u.RawQuery)

	result := []data.DatabaseEntity{}

	playlistIdArr, ok := querry["list"]
	if ok {
		playlistId := playlistIdArr[0]
		ids, titles, channelTitle := api.GetVideoIdsFromPlaylistItems(playlistId)
		playlistTitle := api.GetPlaylistTitleFromListID(playlistId)
		if ids != nil {
			for i, id := range ids {
				result = append(result, data.DatabaseEntity{
					VideoID:       id,
					VideoTitle:    videoTitleFilter(titles[i]),
					ChannelTitle:  channelTitle,
					PlaylistTitle: playlistTitle,
					InWork:        0,
				})
			}
		}
		return result
	}
	return nil
}

func videoTitleFilter(title string) string {
	title = strings.ReplaceAll(title, " | ", "-")
	re := regexp.MustCompile("#.+(,|$|\t)?")

	title = re.ReplaceAllString(title, "")
	re = regexp.MustCompile(`[^0-9a-zA-Zа-яА-Я,!?:;%-\\ ]`)
	title = re.ReplaceAllString(title, "")
	title = strings.TrimSpace(title)
	if title == "" {
		title = bot.RandomString(12)
	}
	return title
}
