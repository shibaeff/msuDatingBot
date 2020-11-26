package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

type MessageToAdmin struct {
	*tgbotapi.MessageConfig
	AdminsList []string
}
