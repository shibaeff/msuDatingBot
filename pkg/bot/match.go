package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

type Match struct {
	Msg1 *tgbotapi.MessageConfig
	Msg2 *tgbotapi.MessageConfig
}
