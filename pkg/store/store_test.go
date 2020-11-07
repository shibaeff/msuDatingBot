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
	usersCollectionName = "users"
	likesCollectionName = "likes"
	seenCollectionName  = "seen"

	putNewLike = "put new like success"
	updLike    = "update like success"
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

func Test_store_PutLike(t *testing.T) {
	collection, _ := prepareCollection(likesCollectionName)
	st := NewStore(nil, collection, nil)
	t.Run(putNewLike, func(t *testing.T) {
		err := st.PutLike(whoID, whomeID)
		assert.NoError(t, err)
		likes, err := st.GetLikes(whoID)
		assert.NoError(t, err)
		assert.Equal(t, int64(whoID), likes.Who)
		filter := bson.D{
			{
				"who", whoID,
			},
		}
		deleteRes, err := st.(*store).likesRegistry.(*registry).collection.DeleteOne(context.TODO(), filter)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), deleteRes.DeletedCount)
	})

	t.Run(putNewLike, func(t *testing.T) {
		err := st.PutLike(whoID, whomeID)
		assert.NoError(t, err)
		err = st.PutLike(whoID, anotherWhome)
		assert.NoError(t, err)
		likes, err := st.GetLikes(whoID)
		assert.NoError(t, err)
		assert.Equal(t, int64(whoID), likes.Who)
		assert.Equal(t, int64(whomeID), likes.Whome[0])
		assert.Equal(t, int64(anotherWhome), likes.Whome[1])
		filter := bson.D{
			{
				"who", whoID,
			},
		}
		deleteRes, err := st.(*store).likesRegistry.(*registry).collection.DeleteOne(context.TODO(), filter)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), deleteRes.DeletedCount)
	})
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
