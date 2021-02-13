package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"echoBot/pkg/bot"
	"echoBot/pkg/store"
)

const (
	dbName                = "main"
	usersCollectionName   = "users"
	actionsCollectionName = "actions"

	checkUrl = "https://umsu.me/answer?tg_id=%d"
)

var (
	client, _ = mongo.NewClient(options.Client().ApplyURI("mongodb+srv://user:123@userscluster.whxir.mongodb.net/users?retryWrites=true&w=majority"))
)

func PrepareCollection(client *mongo.Client, name string) (conn *mongo.Collection) {
	conn = client.Database(dbName).Collection(name)
	return
}

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func switchReply(api *tgbotapi.BotAPI, reply interface{}) {
	switch v := reply.(type) {
	case *tgbotapi.MessageConfig:
		if _, err := api.Send(v); err != nil {
			log.Println(err)
		}
	case *tgbotapi.PhotoConfig:
		api.Send(v)
	case *tgbotapi.DocumentConfig:
		api.Send(v)
	}
}

var banned, allowed map[int64]bool

func ban(name int64) {
	if banned == nil {
		banned = make(map[int64]bool)
	}
	banned[name] = true
}
func checkUserStatus(id int64) bool {
	url := fmt.Sprintf(checkUrl, id)
	resp, _ := http.Get(url)
	bytes, _ := ioutil.ReadAll(resp.Body)
	status := string(bytes)
	if status == "3" || status == "0" {
		ban(int64(id))
		return false
	} else {
		if allowed == nil {
			allowed = make(map[int64]bool)
		}
		allowed[int64(id)] = true
	}
	return true
}

var (
	api *tgbotapi.BotAPI
)

func bannedReply(update tgbotapi.Update) {
	api.Send(tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID: update.Message.Chat.ID,
		},
		Text: "К сожалению, верификация не пройдена ",
	})
}

func main() {
	log.Println("hello")
	err := client.Connect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	users := PrepareCollection(client, usersCollectionName)
	actions := PrepareCollection(client, actionsCollectionName)
	store := store.NewStore(users, actions)

	token, exists := os.LookupEnv("TELEGRAM_TOKEN")
	if !exists {
		log.Panic("Not telegram key specified!")
	}
	api, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	logFile, err := os.OpenFile("log.txt", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	readFile, err := os.OpenFile("log.txt", os.O_RDONLY, 0666)
	defer readFile.Close()
	if err != nil {
		panic(err)
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)

	api.Debug = true

	log.Printf("Authorized on account %s", api.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// updates, err := api.GetUpdatesChan(u)
	admins, ok := os.LookupEnv("ADMINS")
	if !ok {
		log.Fatal("cannot load admins")
	}
	admins_list := strings.Split(admins, " ")
	Bot := bot.NewBot(store, api, readFile, admins_list)
	defer logFile.Close()
	updates, _ := api.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil && update.CallbackQuery == nil { // ignore any non-Message Updates
			continue
		}
		if strings.HasPrefix(update.Message.Text, "/strike") && Bot.EnsureAdmin(update.Message.Chat.UserName) {
			split := strings.Split(update.Message.Text, " ")
			if len(split) == 2 {
				user := Bot.GetStore().FindUser(bson.D{
					{"username", split[1]},
				})
				if banned == nil {
					banned = make(map[int64]bool)
				}
				banned[user.Id] = true
				api.Send(tgbotapi.MessageConfig{
					BaseChat: tgbotapi.BaseChat{
						ChatID: update.Message.Chat.ID,
					},
					Text: "Пользователь забанен",
				})
				continue
			}
		}
		var reply interface{}
		if update.Message != nil {
			ctx, _ := context.WithTimeout(context.Background(), time.Second)
			go func() {
				// register or ban
				if update.Message.Text == "/start" {
					reply, err = Bot.ReplyMessage(ctx, update.Message)
					if err != nil {
						log.Fatal(err)
					}
					switchReply(api, reply)
					if checkUserStatus(update.Message.Chat.ID) {
						api.Send(tgbotapi.MessageConfig{
							BaseChat: tgbotapi.BaseChat{
								ChatID: update.Message.Chat.ID,
							},
							Text: "Верификация успешно пройдена.\n\nДля продолжения нажмите /register",
						})
					} else {
						bannedReply(update)
					}
				} else {
					_, allow := allowed[update.Message.Chat.ID]
					_, block := banned[update.Message.Chat.ID]
					if block {
						bannedReply(update)
						return
					}
					if allow || !block {
						reply, err = Bot.ReplyMessage(ctx, update.Message)
						if err != nil {
							log.Fatal(err)
						}
						switchReply(api, reply)
					} else {
						api.Send(tgbotapi.MessageConfig{
							BaseChat: tgbotapi.BaseChat{
								ChatID: update.Message.Chat.ID,
							},
							Text: "Пройдите верификацию по ссылке в /start",
						})
					}
				}
			}()
		} else {
			ctx, _ := context.WithTimeout(context.Background(), time.Second)
			go func() {
				reply, _ = Bot.HandleCallbackQuery(ctx, update.CallbackQuery)
				switchReply(api, reply)
				if err != nil {
					log.Fatal(err)
				}
			}()
		}
	}
}
