# WeatherBot

![GitHub package.json version](https://img.shields.io/badge/version-1.0.0-blue)
![Python](https://img.shields.io/badge/python-3.8%2B-green)
![Telegram Bot](https://img.shields.io/badge/telegram-bot-blueviolet)

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
- [telegram-bot](github.com/go-telegram-bot-api/telegram-bot-api/v5)
- [Requests](https://pypi.org/project/requests/) (для HTTP‑запросов к Weather API)
- SQLite (встроенная БД)
- `dotenv` (для работы с переменными окружения)

## Установка и запуск

1. Клонируйте репозиторий:
   ```bash
   git clone https://github.com/p1Relly/weatherbot.git
   cd weatherbot
