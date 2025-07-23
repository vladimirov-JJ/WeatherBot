package handler

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/p1relly/weatherbot/internal/openweather"
)

type Handler struct {
	bot      *tgbotapi.BotAPI
	owClient *openweather.OpenWeatherClient
}

func New(bot *tgbotapi.BotAPI, owClient *openweather.OpenWeatherClient) *Handler {
	return &Handler{
		bot:      bot,
		owClient: owClient,
	}
}

func (h *Handler) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := h.bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message != nil {
			h.Callback(update)
		} else if update.CallbackQuery != nil {
			h.MainMenuCallBack(update)
		}
	}
}
