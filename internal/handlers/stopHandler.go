package handlers

import (
	"context"
	"os/exec"
	"strconv"
	"ytdlp-bot/internal/data"
	"ytdlp-bot/internal/environment"
	"ytdlp-bot/internal/security"

	"os"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func StopCommandHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if !security.CheckAuth(update.Message.From.ID) {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Пользователь с id " + strconv.FormatInt(update.Message.From.ID, 10) + " не авторизован",
		})
		return
	}
	sendYesNoDialogInlineKeyBoard(ctx, b, update.Message.Chat.ID, "stop", update.Message.ID, "Вы действительно хотите остановить бота")
}
func StopCallbackHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	dialogHandler(ctx, b, update, selfStop)
}

func selfStop(ctx context.Context, b *bot.Bot, update *models.Update) {
	KillOtherProcesses()
	data.RemoveNextProcessNotEqualArgsFromDB(-1)
	switch environment.Environment.Stage {
	case "DEBUG":
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.CallbackQuery.Message.Message.Chat.ID,
			Text:   "Бот остановлен",
		})
		b.Close(ctx)
		os.Exit(0)
	case "DEPLOY":
		exec.Command("systemctl", "stop", "bot-app.service").Run()
	}

}
