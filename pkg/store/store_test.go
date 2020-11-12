package store

import (
	"context"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"echoBot/pkg/models"
)

const (
	usersCollectionName = "users"
	likesCollectionName = "likes"
	seenCollectionName  = "seen"

	putNewLike = "put new like success"
	updLike    = "update like success"

	noUsersGetAny   = "successful get any test when there are no users"
	oneUserGetAny   = "successful test with only one user"
	manyUsersGetAny = "successful test with two and more users"
)

var (
	pasha = models.User{
		Id:   whoID,
		Name: "pasha",
	}
	ksyusha = models.User{
		Id:   whomeID,
		Name: "ksyusha",
	}
)

func Test_store_Put(t *testing.T) {
	collection, _ := prepareCollection(usersCollectionName)
	ash := models.User{
		Name: "Peter",
	}
	store := NewStore(collection, []*mongo.Collection{nil, nil, nil})
	err := store.PutUser(&ash)
	assert.NoError(t, err)

	var result models.User
	filter := bson.D{{"name", "Peter"}}
	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}
	assert.NoError(t, err)
	assert.Equal(t, ash.Name, result.Name)
}

func Test_store_GetUser(t *testing.T) {
	collection, _ := prepareCollection(usersCollectionName)
	ash := models.User{
		Name: "Peter",
		Id:   whoID,
	}
	store := NewStore(collection, []*mongo.Collection{nil, nil, nil})
	err := store.PutUser(&ash)
	assert.NoError(t, err)
	var result *models.User
	result, err = store.GetUser(whoID)
	assert.NoError(t, err)
	assert.Equal(t, ash.Id, result.Id)
	filter := bson.D{{"id", whoID}}
	collection.DeleteOne(context.TODO(), filter)
}

func Test_store_UpdUserField(t *testing.T) {
	collection, _ := prepareCollection(usersCollectionName)
	ash := models.User{
		Name: "Peter",
		Id:   whoID,
	}
	store := NewStore(collection, []*mongo.Collection{nil, nil, nil})
	err := store.PutUser(&ash)
	assert.NoError(t, err)
	var result *models.User
	result, err = store.GetUser(whoID)
	assert.NoError(t, err)
	assert.Equal(t, ash.Id, result.Id)

	ash.Name = "Pasha"
	err = store.UpdUserField(ash.Id, "name", ash.Name)
	assert.NoError(t, err)
	err = store.DeleteUser(ash.Id)
	assert.NoError(t, err)
}

func Test_store_PutLike(t *testing.T) {
	collection, _ := prepareCollection(likesCollectionName)
	st := NewStore(nil, []*mongo.Collection{collection, nil, nil})
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

func Test_store_GetAny(t *testing.T) {
	usersCollection, _ := prepareCollection(usersCollectionName)
	store := NewStore(usersCollection, []*mongo.Collection{nil, nil})

	t.Run(noUsersGetAny, func(t *testing.T) {
		_, err := store.GetUser(whoID)
		assert.Error(t, err)
	})

	t.Run(oneUserGetAny, func(t *testing.T) {
		err := store.PutUser(&pasha)
		assert.NoError(t, err)
		user, err := store.GetUser(pasha.Id)
		assert.NoError(t, err)
		assert.Equal(t, pasha.Id, user.Id)
		_, err = store.GetUser(whomeID)
		assert.Error(t, err)
		err = store.DeleteUser(pasha.Id)
		assert.NoError(t, err)
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
