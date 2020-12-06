package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

type Match struct {
	Upd  interface{}
	Msg1 *tgbotapi.MessageConfig
	Msg2 *tgbotapi.MessageConfig
}
