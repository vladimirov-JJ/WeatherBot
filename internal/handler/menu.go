package handler

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func mainMenu() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Ввести город", "enter_city"),
			tgbotapi.NewInlineKeyboardButtonData("Отправить геолокацию", "send_location"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("Открыть сайт", "https://t.me/ivanvladimirov"),
		),
	)
}

func (h *Handler) MainMenuCallBack(update tgbotapi.Update) {
	callback := update.CallbackQuery
	chatID := callback.Message.Chat.ID
	data := callback.Data

	switch {
	case data == "enter_city":
		userState[chatID] = "waiting_city"
		h.bot.Send(tgbotapi.NewMessage(chatID, "Введи название города:"))

	case data == "send_location":
		userState[chatID] = "waiting_location"
		h.bot.Send(tgbotapi.NewMessage(chatID, "Отправь координаты:"))

	case strings.HasPrefix(data, "copy_coords:"):
		coords := strings.TrimPrefix(data, "copy_coords:")
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("📌 Координаты:\n`%s`", coords))
		msg.ParseMode = "Markdown"
		h.bot.Send(msg)

	default:
		h.bot.Send(tgbotapi.NewMessage(chatID, "Неизвестное действие 🤔"))
	}

	h.bot.Request(tgbotapi.NewCallback(callback.ID, ""))
}
