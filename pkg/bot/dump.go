package bot

import (
	"echoBot/pkg/models"
	"echoBot/pkg/store"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"io/ioutil"
	"log"
)

type DumpItem struct {
	Likes   []store.Entry
	View    []store.Entry
	Matches []store.Entry
}

type UserData struct {
	Model  models.User
	Arrays DumpItem
}

func (b *bot) dumpEntire() {
	users, _ := b.store.GetAllUsers()
	var items []UserData
	for _, user := range users {
		var item DumpItem
		likes, _ := b.store.GetActions().GetEvents(store.Options{
			bson.E{
				"who", user.Id,
			},
			bson.E{
				"event", store.EventLike,
			},
		})
		seen, _ := b.store.GetActions().GetEvents(store.Options{
			bson.E{
				"who", user.Id,
			},
			bson.E{
				"event", store.EventView,
			},
		})
		matches, _ := b.store.GetActions().GetEvents(store.Options{
			bson.E{
				"who", user.Id,
			},
			bson.E{
				"event", store.EventMatch,
			},
		})
		for _, like := range likes {
			item.Likes = append(item.Likes, like)
		}
		for _, view := range seen {
			item.View = append(item.View, view)
		}
		for _, match := range matches {
			item.Matches = append(item.Matches, match)
		}
		items = append(items, UserData{
			Model:  *user,
			Arrays: item,
		})
	}
	c, err := json.Marshal(items)
	if err != nil {
		panic(err)
	}
	_ = ioutil.WriteFile("dump.json", c, 0644)
	log.Println(items)
}
