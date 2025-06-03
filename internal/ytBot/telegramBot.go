package YtBot

import (
	"bytes"
	"context"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"time"
	"ytdlp-bot/internal/data"
	"ytdlp-bot/internal/environment"
	"ytdlp-bot/internal/handlers"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

var ytBot *bot.Bot
var err error
var ctx context.Context
var cancel context.CancelFunc
var ChatId int64
var UserId int64

func RunBot() {
	exec.Command("source setenv").Start()

	env := environment.GetEnvironment()
	ChatId = env.ChatId
	UserId = env.RootUserId
	ctx, cancel = signal.NotifyContext(context.Background(), os.Interrupt)

	log.Printf("Starting bot (pid=%d)\n", os.Getpid())
	opts := []bot.Option{
		bot.WithDefaultHandler(handlers.DefaultHandler),
		bot.WithMessageTextHandler("/restart", bot.MatchTypeExact, handlers.RestartCommandHandler),
		bot.WithMessageTextHandler("/stop", bot.MatchTypeExact, handlers.StopCommandHandler),
		bot.WithMessageTextHandler("/reload_cookies", bot.MatchTypeExact, handlers.ReloadCookiesCommandHandler),
		bot.WithMessageTextHandler("/update_ytdlp", bot.MatchTypeExact, handlers.UpdateYtDlpCommandHandler),

		bot.WithCallbackQueryDataHandler("restart#dialog_button", bot.MatchTypePrefix, handlers.RestartCallbackHandler),
		bot.WithCallbackQueryDataHandler("stop#dialog_button", bot.MatchTypePrefix, handlers.StopCallbackHandler),
	}

	ytBot, err = bot.New(env.TelegramToken, opts...)
	initCommands()

}
func createAndPutToDbMessage(message *models.Message) {
	messageIdForDB := data.Message{
		ID:   message.ID,
		Date: time.Now().UTC().Format("2006-01-02 15:04:05"),
	}
	data.PutMessageToDB(messageIdForDB)
}
func SendMessage(text string) int {
	message, err := ytBot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    ChatId,
		Text:      text,
		ParseMode: models.ParseModeMarkdown,
	})
	if err != nil {
		log.Printf("Ошибка при отправке сообщения \"%s\" в чат  %d", text, ChatId)
		log.Printf(err.Error())
	}
	createAndPutToDbMessage(message)
	return message.ID
}

func EditMessageCreateIfNotExist(messageID int, text string) int {
	messageId := messageID
	editParams := bot.EditMessageTextParams{
		ChatID:    ChatId,
		MessageID: messageID,
		Text:      text,
	}
	_, err := ytBot.EditMessageText(ctx, &editParams)
	if err != nil {
		log.Printf("Редактируемое сообщение с id=%d не существует", messageID)
		messageId = SendMessage(text)
	}
	return messageId
}
func SendDocument(caption string, filename string, filePath string) {
	fileData, errReadFile := os.ReadFile(filePath)
	if errReadFile != nil {
		log.Printf("error read file, %v\n", filePath)
		return
	}

	params := &bot.SendDocumentParams{
		ChatID:   ChatId,
		Document: &models.InputFileUpload{Filename: filename, Data: bytes.NewReader(fileData)},
		Caption:  caption,
	}
	message, err := ytBot.SendDocument(ctx, params)
	if err != nil {
		log.Printf("Ошибка при отправке документа \"%s\" в чат  %d", filePath, ChatId)
	}
	createAndPutToDbMessage(message)
}

func SendVideo(caption string, filename string, filePath string) {
	fileData, errReadFile := os.ReadFile(filePath)
	if errReadFile != nil {
		log.Printf("error read file, %v\n", filePath)
		return
	}

	params := &bot.SendVideoParams{
		ChatID:  ChatId,
		Video:   &models.InputFileUpload{Filename: filename, Data: bytes.NewReader(fileData)},
		Caption: caption,
	}
	message, err := ytBot.SendVideo(ctx, params)
	if err != nil {
		log.Printf("Ошибка при отправке видео \"%s\" в чат  %d", filePath, ChatId)
	}
	createAndPutToDbMessage(message)
}

func DeleteMessages(messageIDs []int) {
	ytBot.DeleteMessages(ctx, &bot.DeleteMessagesParams{
		ChatID:     ChatId,
		MessageIDs: messageIDs,
	})
}
func DeleteMessage(messageID int) {
	ytBot.DeleteMessages(ctx, &bot.DeleteMessagesParams{
		ChatID:     ChatId,
		MessageIDs: []int{messageID},
	})
	data.DeleteMessageWithId(messageID)
}
func DeleteAllMessages() {
	log.Print("Очистка истории сообщений")
	DeleteMessages(data.GetAllMessagesId(true))
}

func GetBot() *bot.Bot {
	return ytBot
}

func GetCtx() context.Context {
	return ctx
}

func GetCancel() context.CancelFunc {
	return cancel
}
func PinMessage(messageID int) {
	ytBot.PinChatMessage(ctx, &bot.PinChatMessageParams{
		ChatID:    ChatId,
		MessageID: messageID,
	})
}

func initCommands() {

	parameters := bot.SetMyCommandsParams{
		Commands: []models.BotCommand{
			{
				Command:     "restart",
				Description: "Перезагрузить бота",
			},
			{
				Command:     "stop",
				Description: "Остановить бота",
			},
			{
				Command:     "reload_cookies",
				Description: "Перезагрузить куки",
			},
			{
				Command:     "update_ytdlp",
				Description: "Обновить утилиту yt-dlp",
			},
		},
	}
	ytBot.SetMyCommands(ctx, &parameters)
}
