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
	nextEmojiButton = tgbotapi.InlineKeyboardButton{Text: nextEmoji}
	likeEmojiButton = tgbotapi.InlineKeyboardButton{Text: likeEmoji}
	nextKeyBoard    = tgbotapi.NewInlineKeyboardMarkup([]tgbotapi.KeyboardButton{nextEmojiButton, likeEmojiButton})
)

func replyWithCard(u *models.User, to int64) (ret *tgbotapi.PhotoConfig) {
	ret = replyWithPhoto(u, to)
	ret.ReplyMarkup = nextKeyBoard
	return
}
