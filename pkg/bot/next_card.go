package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"echoBot/pkg/models"
)

const (
	nextEmoji = "➡"
	likeEmoji = "👍🏻"
)

var (
	nextEmojiButton = tgbotapi.KeyboardButton{Text: nextEmoji}
	likeEmojiButton = tgbotapi.KeyboardButton{Text: likeEmoji}
	nextKeyBoard    = tgbotapi.NewReplyKeyboard([]tgbotapi.KeyboardButton{nextEmojiButton, likeEmojiButton})
)

func replyWithCard(u *models.User, to int64) (ret *tgbotapi.PhotoConfig) {
	ret = replyWithPhoto(u, to)
	ret.ReplyMarkup = nextKeyBoard
	return
}
