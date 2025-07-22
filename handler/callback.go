package handler

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/p1relly/weatherbot/clients/openweather"
	"github.com/p1relly/weatherbot/handler/format"
)

var userState = make(map[int64]string)

func (h *Handler) Callback(update tgbotapi.Update) {
	if update.CallbackQuery != nil {
		return
	}

	chatID := update.Message.Chat.ID

	switch userState[chatID] {
	case "waiting_city":
		city := update.Message.Text
		delete(userState, chatID) // сброс состояния

		coordinates, err := h.owClient.Coordinates(city)
		if err != nil {
			h.bot.Send(tgbotapi.NewMessage(chatID, "Город не найден"))
			return
		}

		weather, err := h.owClient.Weather(coordinates.Lat, coordinates.Lon)
		if err != nil {
			h.bot.Send(tgbotapi.NewMessage(chatID, "Ошибка получения погоды"))
			return
		}

		msg := message(chatID, weather)
		h.bot.Send(msg)

	case "waiting_location":
		location := update.Message.Location
		delete(userState, chatID) // сброс состояния

		weather, err := h.owClient.Weather(location.Latitude, location.Longitude)
		if err != nil {
			h.bot.Send(tgbotapi.NewMessage(chatID, "Ошибка получения погоды"))
			return
		}

		msg := message(chatID, weather)
		h.bot.Send(msg)
	}

	menuMsg := tgbotapi.NewMessage(chatID, "Выберите действие:")
	menuMsg.ReplyMarkup = mainMenu()
	h.bot.Send(menuMsg)
}

func message(chatID int64, weather openweather.WeatherResponse) tgbotapi.MessageConfig {
	msgWeather := format.MessageWeather(weather)
	msg := tgbotapi.NewMessage(chatID, msgWeather)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Скопировать координаты",
				fmt.Sprintf("copy_coords:%.6f %.6f", weather.Coord.Lat, weather.Coord.Lon)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("Открыть в Google Maps", fmt.Sprintf("https://maps.google.com/?q=%.6f %.6f", weather.Coord.Lat, weather.Coord.Lon)),
		),
	)

	return msg
}
