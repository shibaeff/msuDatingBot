package bot

import (
	"echoBot/pkg/models"
	"echoBot/pkg/store"
	"go.mongodb.org/mongo-driver/bson"
)

func (b *bot) reset(user *models.User) {
	// remove all user's records from the registry
	b.store.GetActions().DeleteEvents(store.Options{
		bson.E{
			"who", user.Id,
		},
	})
	// repopulate user's preferences
	users, _ := b.store.GetAllUsers()
	for _, another := range users {
		b.store.GetActions().AddEvent(store.Entry{
			Who:   user.Id,
			Whome: another.Id,
			Event: store.EventUseen,
		})
	}
}
