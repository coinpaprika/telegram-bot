package telegram

import "gopkg.in/telegram-bot-api.v4"

type BotConfig struct {
	Token          string
	Debug          bool
	UpdatesTimeout int
}

type Bot struct {
	Bot    *tgbotapi.BotAPI
	Config BotConfig
}

type Message struct {
	ChatID    int
	MessageID int
	Text      string
}
