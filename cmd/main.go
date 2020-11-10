package main

import (
	"context"
	"log"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"echoBot/pkg/bot"
	"echoBot/pkg/store"
)

const (
	dbName              = "main"
	usersCollectionName = "users"
	likes               = "likes"
	seen                = "seen"
)

var (
	client, _ = mongo.NewClient(options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
)

func PrepareCollection(client *mongo.Client, name string) (conn *mongo.Collection) {
	conn = client.Database(dbName).Collection(name)
	return
}

func main() {
	err := client.Connect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	users := PrepareCollection(client, usersCollectionName)
	likes := PrepareCollection(client, likes)
	seen := PrepareCollection(client, seen)

	store := store.NewStore(users, likes, seen)

	api, err := tgbotapi.NewBotAPI("1327834524:AAFSH9KVrRiowoqo8uCGdm5EfBIk9Hdxurs")
	if err != nil {
		log.Panic(err)
	}

	api.Debug = true

	log.Printf("Authorized on account %s", api.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := api.GetUpdatesChan(u)
	Bot := bot.NewBot(store)
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
