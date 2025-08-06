package main

import (
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/p1relly/weatherbot/internal/handler"
	"github.com/p1relly/weatherbot/internal/logger"
	"github.com/p1relly/weatherbot/internal/openweather"
	"github.com/p1relly/weatherbot/internal/storage/sqlite"
)

func main() {
	log := logger.SetupLogger("./log/weatherbot.log")
	log.Info("Starting weatherbot")

	err := godotenv.Load(".env")
	if err != nil {
		log.Error("Error loading .env file")
		os.Exit(1)
	}
	log.Info("Loading .env file")

	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	log.Info("Connected TGBotAPI")

	bot.Debug = true
	log.Info("Authorized on account:", bot.Self.UserName)

	owClient := openweather.New(os.Getenv("OPENWEATHER_TOKEN"))
	log.Info("Connected OpenweatherAPI")

	sqliteStoragePath := "storage/storage.db"
	db, err := sqlite.New(sqliteStoragePath)
	if err != nil {
		log.Error("error storage:", err)
		os.Exit(1)
	}
	log.Info("Loading storage")
	// defer dbb.Close()

	botHandler := handler.New(bot, owClient, db)

	botHandler.Start(log)
}
