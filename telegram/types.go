package telegram

import "gopkg.in/telegram-bot-api.v4"

// BotConfig configuration of the bot
type BotConfig struct {
	Token          string
	Debug          bool
	UpdatesTimeout int
}

// Bot telegram interaction client
type Bot struct {
	Bot    *tgbotapi.BotAPI
	Config BotConfig
}

// Message a telegram message struct
type Message struct {
	ChatID    int
	MessageID int
	Text      string
}
