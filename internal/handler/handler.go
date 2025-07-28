package handler

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/p1relly/weatherbot/internal/openweather"
	"github.com/p1relly/weatherbot/internal/storage/sqlite"
)

type Handler struct {
	bot      *tgbotapi.BotAPI
	owClient *openweather.OpenWeatherClient
	db       *sqlite.Storage
}

func New(bot *tgbotapi.BotAPI, owClient *openweather.OpenWeatherClient, db *sqlite.Storage) *Handler {
	return &Handler{
		bot:      bot,
		owClient: owClient,
		db:       db,
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
			h.CallbackQuery(update)
		}
	}
}
