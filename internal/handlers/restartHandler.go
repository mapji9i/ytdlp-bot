package handlers

import (
	"context"
	"os"
	"os/exec"
	"strconv"
	"ytdlp-bot/internal/data"
	"ytdlp-bot/internal/environment"
	"ytdlp-bot/internal/security"

	// "log"
	// "os"
	// "os/exec"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func RestartCommandHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if !security.CheckAuth(update.Message.From.ID) {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Пользователь с id " + strconv.FormatInt(update.Message.From.ID, 10) + " не авторизован",
		})
		return
	}
	sendYesNoDialogInlineKeyBoard(ctx, b, update.Message.Chat.ID, "restart", update.Message.ID, "Вы действительно хотите перезагрузить бота")
}

func RestartCallbackHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	dialogHandler(ctx, b, update, selfRestart)

}

func selfRestart(ctx context.Context, b *bot.Bot, update *models.Update) {
	KillOtherProcesses()
	data.RemoveNextProcessNotEqualArgsFromDB(-1)
	switch environment.Environment.Stage {
	case "DEBUG":
		workingDir := environment.Environment.WorkingDir
		cmd := exec.Command(workingDir+"/"+environment.Environment.RunScriptName, "-r")
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		b.Close(ctx)
		cmd.Start()
		os.Exit(0)
	case "DEPLOY":
		exec.Command("systemctl", "restart", "bot-app.service").Run()
	}
}

func KillOtherProcesses() {
	pid := os.Getpid()
	for true {
		proc, err := data.RemoveNextProcessNotEqualArgsFromDB(pid)
		if err != nil {
			break
		}
		exec.Command("kill", strconv.Itoa(proc.PID)).Start()
	}
}
