package bot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	registerCommand = "/register"
)

type Bot interface {
	Reply(message *tgbotapi.Message) *tgbotapi.MessageConfig
}

type bot struct {
	state State
}

func (b *bot) Reply(message *tgbotapi.Message) (reply *tgbotapi.MessageConfig) {
	reply = b.state.Reply(message)
	if b.state.(*InitState).name == initName {
		log.Printf("State switched to %s")
		b.state = &WaitingState{}
	}
	if message.Text == registerCommand {
		b.state = &RegisterState{}
	}
	return
}

func NewBot() (b Bot) {
	b = &bot{
		state: &InitState{},
	}
	return b
}
