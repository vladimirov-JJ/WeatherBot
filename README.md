# WeatherBot

![GitHub package.json version](https://img.shields.io/badge/version-1.0.3-blue)
![Go](https://img.shields.io/badge/go-1.24.2-blue)
![Telegram Bot](https://img.shields.io/badge/telegram_bot-blueviolet)
![Pet Project](https://img.shields.io/badge/pet--project-personal-brightgreen)

## Описание
Telegram-бот на **Go** для получения данных о погоде и автоматической оценки пригодности полётов для различных моделей дронов.

Коротко: принимает название города или геопозицию из Telegram, получает метеоданные от OpenWeather, сверяет их с ограничениями моделей дронов, сохраняет данные в локальной SQLite и поддерживает развёртывание через GitHub Actions.

## Основные возможности

- Получение текущей погоды по названию города или по геопозиции (OpenWeather API).
- Оценка пригодности полёта для сохранённых моделей дронов.
- Хранение характеристик дронов и историй запросов в SQLite.
- Простая конфигурация через `.env` (токены и пути к БД).
- Шаблон для CI/CD и деплоя (GitHub Actions + деплой по SSH).

## Как пользоваться

- Отправьте боту название города — он ответит текущей погодой и рекомендацией по нежелательным полетам.
- Отправьте геопозицию — бот обрабатывает координаты и возвращает те же данные.
- В ответе указываются: температура, влажность, скорость и направление ветра, предупреждения для моделей дронов.

## Технологии
- Язык: Go 1.24.2+
- БД: SQLite
- API: [Telegram Bot](https://github.com/go-telegram-bot-api/telegram-bot-api), [OpenWeather](https://openweathermap.org/)
- CI/CD: GitHub Actions

## Установка и запуск

1. Клонируйте репозиторий:

```bash
git clone https://github.com/vladimirov-JJ/WeatherBot.git
cd WeatherBot
```

2. Создайте файл `.env` рядом с бинарником или в корне проекта (пример ниже).

3. Создайте ключи в github: (settings -> secrets and variables)
   ```bash
   BOT_TOKEN
   OPENWEATHER_TOKEN
   DEPLOY_SSH_KEY

### Пример `.env`

```env
TELEGRAM_TOKEN=your_telegram_bot_token
OPENWEATHER_KEY=your_openweather_api_key
DB_PATH=./weatherbot.db