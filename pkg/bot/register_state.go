package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	registerReply = "Начата регистрация"
)

type RegisterState struct {
}

func (r *RegisterState) Reply(message *tgbotapi.Message) (reply *tgbotapi.MessageConfig) {
	reply = &tgbotapi.MessageConfig{
		Text: registerReply,
	}
	return
}
