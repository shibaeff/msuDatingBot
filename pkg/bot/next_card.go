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
	nextData        = nextEmoji
	likeData        = likeEmoji
	nextEmojiButton = tgbotapi.InlineKeyboardButton{Text: nextEmoji, CallbackData: &nextData}
	likeEmojiButton = tgbotapi.InlineKeyboardButton{Text: likeEmoji, CallbackData: &likeData}
	nextKeyBoard    = tgbotapi.NewInlineKeyboardMarkup([]tgbotapi.InlineKeyboardButton{nextEmojiButton, likeEmojiButton})
)

func (b *bot) replyWithCard(candidate *models.User, whome int64) (ret *tgbotapi.PhotoConfig) {
	ret = candidate.ReplyWithPhoto()
	ret.ChatID = whome
	ret.ReplyMarkup = nextKeyBoard
	return
}
