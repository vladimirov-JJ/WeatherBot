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
	err := godotenv.Load("/root/weatherbot/.env")
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

	botHandler := handler.New(bot, owClient)

	sqliteStoragePath := "/root/weatherbot/"
	storage, err := sqlite.New(sqliteStoragePath)
	if err != nil {
		log.Fatalf("error storaggggge: %s", err)
	}

	_, err = storage.SaveURL("https://google.com", "google")
	if err != nil {
		log.Fatalf("error storageee: %s", err)
		return
	}

	// _, err = storage.SaveURL("https://google.com", "google")
	// if err != nil {
	// 	log.Fatalf("error storage: %s", err)
	// 	return
	// }

	botHandler.Start()
}
