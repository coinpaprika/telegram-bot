package telegram

import (
	"github.com/pkg/errors"
	"gopkg.in/telegram-bot-api.v4"
)

type BotConfig struct {
	Token          string
	Debug          bool
	UpdatesTimeout int
}

type Bot struct {
	Bot    *tgbotapi.BotAPI
	Config BotConfig
}

func NewBot(c BotConfig) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(c.Token)
	if err != nil {
		return nil, errors.Wrap(err, "could not create telegram bot")
	}

	bot.Debug = c.Debug

	return &Bot{
		Bot:    bot,
		Config: c,
	}, nil
}
func (b *Bot) GetUpdatesChannel() (tgbotapi.UpdatesChannel, error) {
	updatesConfig := tgbotapi.NewUpdate(0)
	if b.Config.UpdatesTimeout > 0 {
		updatesConfig.Timeout = b.Config.UpdatesTimeout
	}
	return b.Bot.GetUpdatesChan(updatesConfig)
}
