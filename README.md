# WeatherBot

![GitHub package.json version](https://img.shields.io/badge/version-1.0.3-blue)
![Go](https://img.shields.io/badge/go-1.24.2-blue)
![Telegram Bot](https://img.shields.io/badge/telegram_bot-blueviolet)

## Описание
**WeatherBot** — это простой Telegram‑бот, который позволяет:
- Узнавать текущую погоду по названию города или отправленной геопозиции;
- Проверять совместимость полётов для разных моделей дронов на основе базы данных.  

Проект написан на Golang, использует Telegram Bot API и внешнее Weather API, а также встроенную SQLite‑базу данных для хранения информации о дронах.

## Возможности
- Получение прогноза погоды (температура, влажность, скорость ветра и т. д.).
- Определение, на каких моделях дронов полёты не рекомендованы (в зависимости от параметров погоды).
- Автоматическая обработка текстовых запросов и геопозиции пользователя.

## Технологии
- Go 1.24.2+
- [telegram-bot](https://github.com/go-telegram-bot-api/telegram-bot-api)
- [dotenv](https://github.com/joho/godotenv) (для работы с переменными окружения)
- `SQLite` (встроенная БД)

## Установка и запуск

1. Клонируйте репозиторий:
   ```bash
   git clone https://github.com/p1Relly/weatherbot.git
   cd weatherbot

2. Создайте ключи в github: (settings -> secrets and variables)
   ```bash
   BOT_TOKEN
   OPENWEATHER_TOKEN
   DEPLOY_SSH_KEY

3. Соберите и запустите бота:
   ```bash
   go build -o weatherbot .
   ./weatherbot