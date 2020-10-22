package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	waitingMessage = "Ждем Вашего действия. Воспользуйтесь клавиатурой"
)

type WaitingState struct {
}

func (i *WaitingState) Reply(message *tgbotapi.Message) (reply *tgbotapi.MessageConfig) {
	reply = &tgbotapi.MessageConfig{
		Text: waitingMessage,
	}
	return
}

func NewWaitingState() *WaitingState {
	return &WaitingState{}
}
