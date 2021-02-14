package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.mongodb.org/mongo-driver/bson"
	"log"
)

func (b *bot) feedback(text string) {
	for _, admin := range b.adminsList {
		adm := b.store.FindUser(
			bson.D{
				{"username", admin},
			},
		)
		if adm == nil {
			continue
		}
		_, err := b.api.Send(&tgbotapi.MessageConfig{
			Text: text,
			BaseChat: tgbotapi.BaseChat{
				ChatID: adm.Id,
			},
		})
		if err != nil {
			log.Println(err)
			continue
		}
	}
}
