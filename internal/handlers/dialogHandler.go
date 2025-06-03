package handlers

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func sendYesNoDialogInlineKeyBoard(ctx context.Context, b *bot.Bot, chatId int64, command string, reasonMessageId int, dialogMessageText string) {
	kb := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "Да", CallbackData: fmt.Sprintf("%s#dialog_button_yes#%d", command, reasonMessageId)},
				{Text: "Нет", CallbackData: fmt.Sprintf("%s#dialog_button_no#%d", command, reasonMessageId)},
			},
		},
	}
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatId,
		Text:        dialogMessageText,
		ReplyMarkup: kb,
	})
}
func dialogHandler(ctx context.Context, b *bot.Bot, update *models.Update, fun func(ctx context.Context, b *bot.Bot, update *models.Update)) {
	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})
	calbackDataArr := strings.Split(update.CallbackQuery.Data, "#")
	initMSGId, _ := strconv.Atoi(calbackDataArr[2])
	calbackData := calbackDataArr[1]

	switch calbackData {
	case "dialog_button_yes":
		deleteInitMessages(ctx, b, update.CallbackQuery.Message.Message.Chat.ID, []int{initMSGId, update.CallbackQuery.Message.Message.ID})
		fun(ctx, b, update)
	case "dialog_button_no":
		deleteInitMessages(ctx, b, update.CallbackQuery.Message.Message.Chat.ID, []int{initMSGId, update.CallbackQuery.Message.Message.ID})
	}
}

func deleteInitMessages(ctx context.Context, b *bot.Bot, chatId int64, messageIDs []int) {
	b.DeleteMessages(ctx, &bot.DeleteMessagesParams{
		ChatID:     chatId,
		MessageIDs: messageIDs,
	})
}
