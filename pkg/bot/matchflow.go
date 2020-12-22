package bot

import (
	"echoBot/pkg/models"
	"echoBot/pkg/store"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	match = "У вас мэтч!\n"
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
	cand, _ := b.store.GetUser(last)
	b.match(user, cand)
	// attempt to get new unseen user
	last, err = b.getLastUnseen(user)
	if err != nil {
		reply = user.ReplyWithText(allSeen)
		return reply
	}
	cand, _ = b.store.GetUser(last)
	return b.replyWithCard(cand, user.Id)
}

func (b *bot) match(user1, user2 *models.User) {
	e1, _ := b.store.GetActions().GetEvents(store.Options{
		bson.E{"who", user1.Id},
		bson.E{"whome", user2.Id},
		bson.E{"event", store.EventLike},
	})
	e2, _ := b.store.GetActions().GetEvents(store.Options{
		bson.E{"who", user2.Id},
		bson.E{"whome", user1.Id},
		bson.E{"event", store.EventLike},
	})
	if len(e1) == 1 && len(e2) == 1 {
		b.store.GetActions().AddEvent(store.Entry{
			Who:   user1.Id,
			Whome: user2.Id,
			Event: store.EventMatch,
		})
		b.store.GetActions().AddEvent(store.Entry{
			Who:   user2.Id,
			Whome: user1.Id,
			Event: store.EventMatch,
		})
		b.api.Send(&tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID: user1.Id,
			},
			Text: match + user2.GetLink(),
		})
		b.api.Send(&tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID: user2.Id,
			},
			Text: match + user1.GetLink(),
		})
	}
}
