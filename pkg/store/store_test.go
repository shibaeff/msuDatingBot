package store

import (
	"context"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"echoBot/pkg/bot"
)

const (
	usersCollectionName = "another"
)

type Trainer struct {
	Name string
	Age  int
	City string
}

func Test_store_Put(t *testing.T) {
	collection, _ := prepareCollection(usersCollectionName)
	ash := bot.User{
		Name: "Peter",
	}
	store := NewStore(collection, nil, nil)
	err := store.Put(&ash)
	assert.NoError(t, err)

	var result Trainer
	filter := bson.D{{"name", "Peter"}}
	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}
	assert.NoError(t, err)
	assert.Equal(t, ash.Name, result.Name)
}

func prepareCollection(name string) (col *mongo.Collection, conn *mongo.Client) {
	client, _ := mongo.NewClient(options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	err := client.Connect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	col = client.Database("another").Collection(name)
	return
}
