package formatter

import (
	"fmt"
	"time"

	"github.com/p1relly/weatherbot/internal/openweather"
)

func unixToTime(unixTime int, timezoneOffset int) string {
	loc := time.FixedZone("offset", timezoneOffset)
	return time.Unix(int64(unixTime), 0).In(loc).Format("15:04")
}

func translateWeather(description string) string {
	descriptionMap := map[string]string{
		"clear sky":            "ясное небо",
		"few clouds":           "малооблачно",
		"scattered clouds":     "облачно",
		"broken clouds":        "облачно с прояснениями",
		"overcast clouds":      "пасмурно",
		"light rain":           "небольшой дождь",
		"moderate rain":        "умеренный дождь",
		"heavy intensity rain": "сильный дождь",
		"drizzle":              "морось",
		"thunderstorm":         "гроза",
		"light snow":           "небольшой снег",
		"heavy snow":           "сильный снег",
		"mist":                 "туман",
	}

	return descriptionMap[description]
}

func windDirection(deg int) string {
	directions := []string{
		"северный", "северо-восточный", "восточный", "юго-восточный",
		"южный", "юго-западный", "западный", "северо-западный",
	}
	idx := int((float64(deg)+22.5)/45.0) % 8
	return directions[idx]
}

func MessageWeather(weather openweather.WeatherResponse) string {
	return fmt.Sprintf(`
📍 Погода в г. %s, %s (Часовой пояс: UTC %+d)

🌡️ Температура: %.0f℃, %s 
	🔻%.0f℃ 🔺%.0f℃
	• Давление: %d гПа
	• Влажность: %d%%

💨 Ветер: %s (%d°), %.1f м/с (порывы до %.1f м/с)

🌫️ Видимость: %d км
🌧️ Осадки за 1ч: %.2f мм

🌅 Восход: %s
🌇 Закат: %s

📌 Координаты:
	%.6f°N (Lat), %.6f°E (Lon)
`,
		weather.Name,
		weather.Sys.Country,
		weather.Timezone/3600,

		weather.Main.Temp,
		translateWeather(weather.Weather[0].Description),
		weather.Main.TempMin,
		weather.Main.TempMax,
		weather.Main.Pressure,
		weather.Main.Humidity,

		windDirection(weather.Wind.Deg),
		weather.Wind.Deg,
		weather.Wind.Speed,
		weather.Wind.Gust,

		weather.Visibility/1000, // m -> km
		weather.Rain.OneH,

		unixToTime(weather.Sys.Sunrise, weather.Timezone),
		unixToTime(weather.Sys.Sunset, weather.Timezone),

		weather.Coord.Lat, weather.Coord.Lon,
	)
}
