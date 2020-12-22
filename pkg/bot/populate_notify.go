package bot

import (
	"echoBot/pkg/models"
	"echoBot/pkg/store"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	notifyMsg = &tgbotapi.MessageConfig{
		Text: "Появился 1 новый пользователь!",
	}
)

func (b *bot) populateNotify(newuser *models.User) {
	users, _ := b.store.GetAllUsers()
	for _, user := range users {
		if !b.ensureGender(newuser, user) {
			continue
		}
		b.store.GetActions().AddEvent(store.Entry{
			Who:   newuser.Id,
			Whome: user.Id,
			Event: store.EventUseen,
		})
		b.store.GetActions().AddEvent(store.Entry{
			Who:   user.Id,
			Whome: newuser.Id,
			Event: store.EventUseen,
		})
		// notify another user
		notifyMsg.ChatID = user.Id
		b.api.Send(notifyMsg)
	}
}
