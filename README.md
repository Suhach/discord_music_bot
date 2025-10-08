# Discord Music Bot

Музыкальный бот для Discord с поддержкой YouTube, Spotify и SoundCloud.

## Возможности

- 🎵 Поиск и воспроизведение музыки из YouTube, Spotify и SoundCloud
- 📋 Система очередей
- ⏸️ Управление воспроизведением (пауза, скип, стоп)
- 🎮 Удобный интерфейс с кнопками

## Требования

### Системные зависимости
- **Go 1.21+**
- **FFmpeg** - для обработки аудио
- **yt-dlp** - для скачивания из YouTube/SoundCloud

### Установка зависимостей

#### Ubuntu/Debian
```bash
sudo apt update
sudo apt install ffmpeg
sudo pip install yt-dlp
```

#### macOS
```bash
brew install ffmpeg
brew install yt-dlp
```

#### Windows
1. Скачайте FFmpeg с https://ffmpeg.org/download.html
2. Установите yt-dlp: `pip install yt-dlp`

## Установка бота

### 1. Клонируйте репозиторий
```bash
git clone https://github.com/Artem/DC_bot.git
cd DC_bot
```

### 2. Получите API ключи

#### Discord Bot Token
1. Зайдите на https://discord.com/developers/applications
2. Создайте новое приложение
3. Перейдите в раздел "Bot"
4. Нажмите "Add Bot"
5. Скопируйте токен
6. В разделе "Privileged Gateway Intents" включите:
   - Server Members Intent
   - Message Content Intent

#### Spotify API
1. Зайдите на https://developer.spotify.com/dashboard
2. Создайте новое приложение
3. Скопируйте Client ID и Client Secret

#### YouTube API
1. Зайдите на https://console.cloud.google.com/
2. Создайте новый проект
3. Включите "YouTube Data API v3"
4. Создайте API ключ в разделе "Credentials"

#### SoundCloud (опционально)
1. Зайдите на https://soundcloud.com/you/apps
2. Зарегистрируйте новое приложение
3. Скопируйте Client ID

### 3. Настройте .env файл
```bash
cp .env.example .env
nano .env
```

Заполните все значения вашими ключами:
```env
DISCORD_TOKEN=your_discord_bot_token
SPOTIFY_CLIENT_ID=your_spotify_client_id
SPOTIFY_SECRET=your_spotify_secret
YOUTUBE_API_KEY=your_youtube_api_key
SOUNDCLOUD_CLIENT_ID=your_soundcloud_client_id
```

### 4. Установите Go зависимости
```bash
go mod download
```

### 5. Запустите бота
```bash
go run cmd/bot/main.go
```

Или скомпилируйте:
```bash
go build -o bot cmd/bot/main.go
./bot
```

## Добавление бота на сервер

1. Перейдите на https://discord.com/developers/applications
2. Выберите ваше приложение
3. Перейдите в "OAuth2" → "URL Generator"
4. Выберите scopes:
   - `bot`
   - `applications.commands`
5. Выберите Bot Permissions:
   - Send Messages
   - Connect
   - Speak
   - Use Slash Commands
6. Скопируйте сгенерированную ссылку и откройте в браузере
7. Добавьте бота на ваш сервер

## Команды

- `/play_s [название]` - Поиск в Spotify
- `/play_y [название]` - Поиск в YouTube
- `/play_sc [название]` - Поиск в SoundCloud
- `/skip` - Пропустить текущий трек
- `/pause` - Пауза/возобновить
- `/stop` - Остановить и отключить бота
- `/queue` - Показать очередь

## Использование

1. Зайдите в голосовой канал
2. Используйте команду `/play_y название песни`
3. Выберите трек из списка, нажав на кнопку
4. Бот подключится к вашему каналу и начнет воспроизведение

## Возможные проблемы

### FFmpeg не найден
```bash
which ffmpeg  # Проверьте установку
```

### yt-dlp не работает
```bash
yt-dlp --update  # Обновите до последней версии
```

### Бот не подключается к голосовому каналу
- Проверьте права бота на сервере
- Убедитесь, что вы находитесь в голосовом канале

### Spotify не работает
- Проверьте правильность Client ID и Secret
- Убедитесь, что приложение активно в Spotify Dashboard

## Структура проекта

```
├── cmd/bot/main.go          # Точка входа
├── internal/
│   ├── bot/bot.go           # Основная логика бота
│   ├── player/player.go     # Аудио плеер и очередь
│   ├── spotify/spotify.go   # Интеграция со Spotify
│   ├── youtube/youtube.go   # Интеграция с YouTube
│   └── soundcloud/soundcloud.go  # Интеграция с SoundCloud
├── pkg/config/config.go     # Конфигурация
└── .env                     # Переменные окружения
```

## Разработка

### Добавление новых команд
Добавьте команду в массив `commands` в `internal/bot/bot.go` и создайте обработчик.

### Тестирование
```bash
go test ./...
```

## Лицензия

MIT

## Поддержка

При возникновении проблем создайте issue на GitHub.
