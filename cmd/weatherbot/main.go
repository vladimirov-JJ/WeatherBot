package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/p1relly/weatherbot/internal/handler"
	"github.com/p1relly/weatherbot/internal/openweather"
	"github.com/p1relly/weatherbot/internal/storage/sqlite"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	owClient := openweather.New(os.Getenv("OPENWEATHER_TOKEN"))

	sqliteStoragePath := "storage/storage.db"
	db, err := sqlite.New(sqliteStoragePath)
	if err != nil {
		log.Fatalf("error storage: %s", err)
	}
	// defer dbb.Close()

	botHandler := handler.New(bot, owClient, db)

	botHandler.Start()
}
