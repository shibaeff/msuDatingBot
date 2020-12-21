package bot

import (
	"echoBot/pkg/store"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.mongodb.org/mongo-driver/bson"
)

func (b *bot) feedback(text string) {
	for _, admin := range b.adminsList {
		adm := b.store.FindUser(store.Options{
			bson.E{
				"username", admin,
			},
		})
		b.api.Send(&tgbotapi.MessageConfig{
			Text: text,
			BaseChat: tgbotapi.BaseChat{
				ChatID: adm.Id,
			},
		})
	}
}
