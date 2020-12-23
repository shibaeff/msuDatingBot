package main

import (
	"context"
	"echoBot/pkg/bot"
	"echoBot/pkg/store"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const (
	dbName = "main"
)

var (
	client, _ = mongo.NewClient(options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
)

func PrepareCollection(client *mongo.Client, name string) (conn *mongo.Collection) {
	conn = client.Database(dbName).Collection(name)
	return
}

func deleteAllRecords(collection *mongo.Collection) {
	collection.DeleteMany(context.TODO(), bson.E{})
}

func main() {
	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	client.Connect(ctx)
	users := PrepareCollection(client, "users")
	actions := PrepareCollection(client, "actions")
	str := store.NewStore(users, actions)
	usrs, _ := str.GetAllUsers()
	for _, usr1 := range usrs {
		for _, usr2 := range usrs {
			if bot.EnsureGender(usr1, usr2) {
				str.GetActions().AddEvent(store.Entry{
					Who:   usr1.Id,
					Whome: usr2.Id,
					Event: store.EventUseen,
				})
				str.GetActions().AddEvent(store.Entry{
					Who:   usr2.Id,
					Whome: usr1.Id,
					Event: store.EventUseen,
				})
			}
		}
	}
}
