package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strings"
)

func (b *bot) users(n int) (reply *tgbotapi.MessageConfig) {
	users, _ := b.store.GetAllUsers()
	if len(users) > n {
		users = users[len(users)-n:]
	}
	var raw []string
	for _, user := range users {
		raw = append(raw, "@"+user.UserName)
	}
	return &tgbotapi.MessageConfig{Text: strings.Join(raw, "\n")}
}
