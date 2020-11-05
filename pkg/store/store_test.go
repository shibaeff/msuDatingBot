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

func Test_store_Put(t *testing.T) {
	collection, _ := prepareCollection()
	store := NewStore(collection)
	//defer func() {
	//	client.Disconnect(context.TODO())
	//}()
	user := bot.NewUser("pavel", "vmk", "m",
		"f", "nohing", 123, "link")
	userModel := NewUserModel(user)
	err := store.Put(userModel)
	assert.NoError(t, err)
	filter := bson.D{{"name", "pavel"}}
	var another UserModel
	err = collection.FindOne(context.TODO(), filter).Decode(&another)
	assert.NoError(t, err)
	assert.Equal(t, userModel.User.GetId(), another.User.GetId())
}

func prepareCollection() (col *mongo.Collection, conn *mongo.Client) {
	client, _ := mongo.NewClient(options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	err := client.Connect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	col = client.Database("another").Collection(usersCollectionName)
	return
}
