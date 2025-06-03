package executors

import (
	"context"
	"log"
	"time"
	"ytdlp-bot/internal/data"
	YtBot "ytdlp-bot/internal/ytBot"
)

func scheduleCleanMessage(ctx context.Context) {
	for {
		nowTime := time.Now().UTC()
		log.Print(nowTime)
		messages := data.GetAllMessagesIdExcludeInput(messageId.GetValue())
		log.Print(messages)
		for _, message := range messages {
			messageCreationTime, _ := time.Parse("2006-01-02 15:04:05 ", message.Date)
			log.Print(messageCreationTime)
			deltaTime := nowTime.Sub(messageCreationTime)
			log.Print(deltaTime)
			if deltaTime >= time.Hour {
				YtBot.DeleteMessage(message.ID)
			}
		}
		time.Sleep(10 * time.Minute)
	}

}
