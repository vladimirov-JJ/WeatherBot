package handler

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

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
