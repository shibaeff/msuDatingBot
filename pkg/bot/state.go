package bot

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

type State interface {
	Reply(message *tgbotapi.Message) (reply *tgbotapi.MessageConfig)
}

var ()
