// Sample Go code for user authorization
package api

import (
	"log"
	"ytdlp-bot/internal/environment"

	"golang.org/x/net/context"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func GetChannelTitleAndVideoTitleFromVideoID(videoId string) (string, string) {
	ctx := context.Background()

	service, _ := youtube.NewService(ctx, option.WithAPIKey(environment.Environment.GoogleToken))
	call := service.Videos.List([]string{"snippet"})
	call = call.Id(videoId)
	response, _ := call.Do()
	return response.Items[0].Snippet.ChannelTitle, response.Items[0].Snippet.Title
}

func GetPlaylistTitleFromListID(listId string) string {
	ctx := context.Background()

	service, _ := youtube.NewService(ctx, option.WithAPIKey(environment.Environment.GoogleToken))

	call := service.Playlists.List([]string{"snippet"})

	call = call.Id(listId)
	response, _ := call.Do()
	return response.Items[0].Snippet.Title
}

func GetVideoIdsFromPlaylistItems(listId string) ([]string, []string, string) {
	ctx := context.Background()

	service, _ := youtube.NewService(ctx, option.WithAPIKey(environment.Environment.GoogleToken))

	call := service.PlaylistItems.List([]string{"snippet"})

	call = call.PlaylistId(listId).MaxResults(300)

	response, err := call.Do()
	if err != nil {
		log.Panic(err)
	}
	var videoIds = []string{}
	var videoTitles = []string{}
	for _, item := range response.Items {
		videoIds = append(videoIds, item.Snippet.ResourceId.VideoId)
		videoTitles = append(videoTitles, item.Snippet.Title)
	}
	channelTitle := response.Items[0].Snippet.ChannelTitle
	return videoIds, videoTitles, channelTitle
}
