package main

import (
	"os"
	"time"
	"ytdlp-bot/internal/data"
	"ytdlp-bot/internal/executors"
	YtBot "ytdlp-bot/internal/ytBot"
)

func main() {

	defer executors.GetService().Close()
	defer YtBot.GetCancel()
	YtBot.RunBot()
	YtBot.DeleteAllMessages()
	botStartInfo()
	go executors.StartReader()
	proc := data.RegisteredProcess{
		PID:  os.Getpid(),
		Date: time.Now().Format("2006-01-02 15:04:05"),
	}

	data.PutProcessToDB(proc)
	YtBot.GetBot().Start(YtBot.GetCtx())
}
func botStartInfo() {
	text := "Бот запущен"

	for _, arg := range os.Args {
		if arg == "-r" {
			text = "Бот перезапущен"
			break
		}
	}
	YtBot.SendMessage(text)
}
