package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

const (
	helloMessage = "Добро пожаловать в бот знакомств"
)

type InitState struct {
}

func (i *InitState) Reply(message *tgbotapi.Message) (reply *tgbotapi.MessageConfig) {
	reply = &tgbotapi.MessageConfig{
		Text: helloMessage,
	}
	return
}

func NewInitState() *InitState {
	return &InitState{}
}
