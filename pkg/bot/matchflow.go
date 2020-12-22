package bot

import (
	"echoBot/pkg/models"
	"echoBot/pkg/store"
	"go.mongodb.org/mongo-driver/bson"
)

func (b *bot) dislike(user *models.User) (reply interface{}) {
	last, err := b.getLastUnseen(user)
	b.store.GetActions().DeleteEvents(store.Options{
		bson.E{"who", user.Id},
		bson.E{"whome", last},
	})
	b.store.GetActions().AddEvent(store.Entry{
		Who:   user.Id,
		Whome: last,
		Event: store.EventView,
	})
	// attempt to get new unseen user
	last, err = b.getLastUnseen(user)
	if err != nil {
		reply = user.ReplyWithText(allSeen)
		return reply
	}
	cand, _ := b.store.GetUser(last)
	return b.replyWithCard(cand, user.Id)
}

func (b *bot) like(user *models.User) (reply interface{}) {
	last, err := b.getLastUnseen(user)
	b.store.GetActions().DeleteEvents(store.Options{
		bson.E{"who", user.Id},
		bson.E{"whome", last},
	})
	b.store.GetActions().AddEvent(store.Entry{
		Who:   user.Id,
		Whome: last,
		Event: store.EventLike,
	})
	// attempt to get new unseen user
	last, err = b.getLastUnseen(user)
	if err != nil {
		reply = user.ReplyWithText(allSeen)
		return reply
	}
	cand, _ := b.store.GetUser(last)
	return b.replyWithCard(cand, user.Id)
}
