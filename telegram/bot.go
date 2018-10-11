package telegram

import (
	"github.com/pkg/errors"
	"gopkg.in/telegram-bot-api.v4"
)

// NewBot creates new telegram bot
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

// GetUpdatesChannel gets new updates updates
func (b *Bot) GetUpdatesChannel() (tgbotapi.UpdatesChannel, error) {
	updatesConfig := tgbotapi.NewUpdate(0)
	if b.Config.UpdatesTimeout > 0 {
		updatesConfig.Timeout = b.Config.UpdatesTimeout
	}
	return b.Bot.GetUpdatesChan(updatesConfig)
}

// SendMessage sends a telegram message
func (b *Bot) SendMessage(m Message) error {
	msg := tgbotapi.NewMessage(int64(m.ChatID), m.Text)
	msg.ReplyToMessageID = m.MessageID
	_, err := b.Bot.Send(msg)
	return errors.Wrapf(err, "could not send message: %v", m)
}
