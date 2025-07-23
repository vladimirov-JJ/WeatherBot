package handler

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/p1relly/weatherbot/internal/formatter"
	"github.com/p1relly/weatherbot/internal/openweather"
)

var userState = make(map[int64]string)

func (h *Handler) CallbackQuery(update tgbotapi.Update) {
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

		msg := messageWithCoordinates(chatID, weather)
		h.bot.Send(msg)

	case "waiting_location":
		location := update.Message.Location
		delete(userState, chatID) // сброс состояния

		weather, err := h.owClient.Weather(location.Latitude, location.Longitude)
		if err != nil {
			h.bot.Send(tgbotapi.NewMessage(chatID, "Ошибка получения погоды"))
			return
		}

		msg := messageWithCoordinates(chatID, weather)
		h.bot.Send(msg)
	}

	h.mainMenu(chatID)
}

func messageWithCoordinates(chatID int64, weather openweather.WeatherResponse) tgbotapi.MessageConfig {
	msgWeather := formatter.MessageWeather(weather)
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
