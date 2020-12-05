package bot

import (
	"echoBot/pkg/models"
	"echoBot/pkg/store"
	"encoding/json"
	"io/ioutil"
	"log"
)

type DumpItem struct {
	Likes   []store.Entry
	Seen    []store.Entry
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
		likes, _ := b.store.GetLikes(user.Id)
		seen, _ := b.store.GetSeen(user.Id)
		matches, _ := b.store.GetMatchesRegistry().GetList(user.Id)
		for _, like := range likes {
			item.Likes = append(item.Likes, like)
		}
		for _, view := range seen {
			item.Seen = append(item.Seen, view)
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
