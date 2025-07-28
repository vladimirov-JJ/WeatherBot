package handler

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

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
	case data == "main_menu":
		h.mainMenu(chatID)

	case data == "enter_city":
		h.bot.Send(tgbotapi.NewMessage(chatID, "Введи название города:"))
		userState[chatID] = "waiting_city"

	case data == "send_location":
		h.bot.Send(tgbotapi.NewMessage(chatID, "Отправь геолокацию:"))
		userState[chatID] = "waiting_location"

	case data == "drone_selection":
		h.droneMenu(chatID)

	case data == "drone_list":
		drone, err := h.db.ListDrone(update.CallbackQuery.From.ID)
		if err != nil {
			h.bot.Send(tgbotapi.NewMessage(chatID, "Ошибка получения списка БВС"))
			break
		}

		if len(drone) == 0 {
			h.bot.Send(tgbotapi.NewMessage(chatID, "У Вас еще нет ни одного БВС"))
			h.mainMenu(chatID)
			break
		}

		var rows [][]tgbotapi.InlineKeyboardButton
		for _, d := range drone {
			title := fmt.Sprintf("[ID:%d] %s (%dгр)", d.ID, d.Name, d.Weight)
			data := fmt.Sprintf("drone_edit_%d", d.ID)
			btn := tgbotapi.NewInlineKeyboardButtonData(title, data)
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(btn))
		}

		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("↩ Управление БВС", "drone_selection"),
		))

		markup := tgbotapi.NewInlineKeyboardMarkup(rows...)

		msg := tgbotapi.NewMessage(chatID, "Список ваших БВС (нажми, чтобы изменить [функционал отсутствует]):")
		msg.ReplyMarkup = markup
		h.bot.Send(msg)

	case data == "drone_add":
		h.bot.Send(tgbotapi.NewMessage(chatID, "Введи название БВС и вес в граммах через запятую, например \"Dji Mavic 4 PRO, 1063\""))
		userState[chatID] = "waiting_drone_add"

	case data == "drone_delete":
		h.bot.Send(tgbotapi.NewMessage(chatID, "Введи ID БВС, который нужно удалить:\n(\"-1\", если хочешь отменить действие)"))
		userState[chatID] = "waiting_drone_delete"

	case strings.HasPrefix(data, "copy_coords:"):
		coords := strings.TrimPrefix(data, "copy_coords:")
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("📌 Координаты:\n`%s`", coords))
		msg.ParseMode = "Markdown"
		h.bot.Send(msg)

	default:
		h.bot.Send(tgbotapi.NewMessage(chatID, "Неизвестное действие"))
		h.mainMenu(chatID)
	}

	h.bot.Request(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))
}

func (h *Handler) Callback(update tgbotapi.Update) {
	if update.CallbackQuery != nil {
		return
	}

	chatID := update.Message.Chat.ID
	userID := update.Message.From.ID

	switch userState[chatID] {

	case "waiting_city":
		city := update.Message.Text
		delete(userState, chatID)

		coordinates, err := h.owClient.Coordinates(city)
		if err != nil || city == "/start" {
			h.bot.Send(tgbotapi.NewMessage(chatID, "Ошибка получения города"))
			break
		}

		h.Message(chatID, userID, coordinates.Lat, coordinates.Lon)

	case "waiting_location":
		location := update.Message.Location
		delete(userState, chatID)

		if location == nil {
			h.bot.Send(tgbotapi.NewMessage(chatID, "Ошибка получения геолокации"))
			break
		}

		h.Message(chatID, userID, location.Latitude, location.Longitude)

	case "waiting_drone_add":
		delete(userState, chatID)

		input := strings.Split(update.Message.Text, ",")
		if len(input) != 2 {
			h.bot.Send(tgbotapi.NewMessage(chatID, "Ошибка получения данных"))
			break
		}

		nameDrone := input[0]
		weight, err := strconv.Atoi(strings.TrimSpace(input[1]))
		if err != nil {
			h.bot.Send(tgbotapi.NewMessage(chatID, "Ошибка получения веса"))
			break
		}

		result, err := h.db.SaveDrone(userID, nameDrone, weight)
		if err != nil {
			h.bot.Send(tgbotapi.NewMessage(chatID, "Ошибка добавления БВС, возможно, такой уже существует"))
			return
		}

		h.bot.Send(tgbotapi.NewMessage(chatID, result))
		h.droneMenu(chatID)
		return

	case "waiting_drone_delete":
		delete(userState, chatID)
		input := update.Message.Text

		droneID, err := strconv.Atoi(input)
		if err != nil {
			h.bot.Send(tgbotapi.NewMessage(chatID, "Ошибка получения ID"))
			return
		}

		if droneID == -1 {
			h.droneMenu(chatID)
			return
		}

		result, err := h.db.DeleteDrone(userID, droneID)
		if err != nil {
			h.bot.Send(tgbotapi.NewMessage(chatID, "Ошибка удаления БВС"))
			h.droneMenu(chatID)
			return
		}

		h.bot.Send(tgbotapi.NewMessage(chatID, result))
		h.droneMenu(chatID)
		return
	}

	h.mainMenu(chatID)
}

func (h *Handler) DroneRecommendations(chatID, userID int64, text *string, weather *openweather.WeatherResponse) {
	drone, err := h.db.ListDrone(userID)
	if err != nil {
		h.bot.Send(tgbotapi.NewMessage(chatID, "Ошибка получения списка БВС"))
		return
	}

	var recommendations string
	for _, d := range drone {
		R := flightRisk(d.Weight, weather, time.Now())
		if R >= 0.4 {
			recommendations += fmt.Sprintf("[⚠️] БВС %s не рекомендован к полёту [%.2f]\n", d.Name, R)
		}
	}
	if recommendations != "" {
		*text += "\nРекмендации:\n" + recommendations
	}
}

func (h *Handler) Message(chatID, userID int64, Lat, Lon float64) {
	weather, err := h.owClient.Weather(Lat, Lon)
	if err != nil {
		h.bot.Send(tgbotapi.NewMessage(chatID, "Ошибка получения погоды"))
		return
	}

	text := formatter.MessageWeather(weather)

	h.DroneRecommendations(chatID, userID, &text, &weather)

	msg := tgbotapi.NewMessage(chatID, text)
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

	h.bot.Send(msg)
}

func flightRisk(mGrams int, weather *openweather.WeatherResponse, now time.Time) float64 {
	mRef := 1000.0
	k, γ := 9.5, 0.7
	Wthr := k * math.Pow(float64(mGrams)/mRef, γ)
	Gthr := 1.2 * Wthr

	var w float64
	if weather.Wind.Speed <= Wthr {
		w = 0
	} else {
		w = (weather.Wind.Speed - Wthr) / Wthr
		if w > 1 {
			w = 1
		}
	}

	var g float64
	if weather.Wind.Gust <= Gthr {
		g = 0
	} else {
		g = (weather.Wind.Gust - Gthr) / Gthr
		if g > 1 {
			g = 1
		}
	}

	p := math.Min(weather.Rain.OneH/5.0, 1)
	v := 1 - math.Min(float64(weather.Visibility)/2000.0, 1)

	Topt, Trange := 17.5, 12.5
	var tn float64
	if weather.Main.Temp >= Topt-Trange && weather.Main.Temp <= Topt+Trange {
		tn = 0
	} else {
		tn = math.Min(math.Abs(weather.Main.Temp-Topt)/Trange, 1)
	}

	var d float64
	hour := now.Hour()
	switch {
	case hour >= 6 && hour < 18:
		d = 0
	case hour >= 18 && hour < 20:
		d = 0.5
	default:
		d = 1
	}

	R := 0.35*w + 0.15*g + 0.20*p + 0.10*v + 0.05*tn + 0.05*d
	return R
}
