package main

import (
	"log"

	"echoBot/pkg/bot"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	api, err := tgbotapi.NewBotAPI("1327834524:AAFSH9KVrRiowoqo8uCGdm5EfBIk9Hdxurs")
	if err != nil {
		log.Panic(err)
	}

	api.Debug = true

	log.Printf("Authorized on account %s", api.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := api.GetUpdatesChan(u)
	Bot := bot.NewBot()
	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		reply := Bot.Reply(update.Message)
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		//buttons := []tgbotapi.KeyboardButton{tgbotapi.KeyboardButton{Text: "Hello",},}
		//msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		//msg.ReplyToMessageID = update.Message.MessageID
		//msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(buttons)
		reply.ChatID = update.Message.Chat.ID
		api.Send(reply)
	}
}
