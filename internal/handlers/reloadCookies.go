package handlers

import (
	"context"
	"os/exec"
	"strconv"
	"ytdlp-bot/internal/environment"
	"ytdlp-bot/internal/security"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func ReloadCookiesCommandHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if !security.CheckAuth(update.Message.From.ID) {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Пользователь с id " + strconv.FormatInt(update.Message.From.ID, 10) + " не авторизован",
		})
		return
	}
	workingDir := environment.Environment.WorkingDir
	cmd := exec.Command(workingDir + "/" + "firefoxCookiesUpdater/start.sh")
	cmd.Run()
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Куки обновлены",
	})
}
