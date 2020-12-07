package store

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

//const (
//	usersCollectionName = "users"
//	likesCollectionName = "likes"
//	seenCollectionName  = "seen"
//
//	putNewLike = "put new like success"
//	updLike    = "update like success"
//
//	noUsersGetAny   = "successful get any test when there are no users"
//	oneUserGetAny   = "successful test with only one user"
//	manyUsersGetAny = "successful test with two and more users"
//	getAllUsersTest = "successful get all the users"
//)
//
//var (
//	pasha = models.User{
//		Id:   whoID,
//		Name: "pasha",
//	}
//	ksyusha = models.User{
//		Id:   whomeID,
//		Name: "ksyusha",
//	}
//)
//
//func Test_store_Put(t *testing.T) {
//	collection, _ := prepareCollection(usersCollectionName)
//	ash := models.User{
//		Name: "Peter",
//	}
//	store := NewStore(collection, []*mongo.Collection{nil, nil, nil})
//	err := store.PutUser(&ash)
//	assert.NoError(t, err)
//
//	var result models.User
//	filter := bson.D{{"name", "Peter"}}
//	err = collection.FindOne(context.TODO(), filter).Decode(&result)
//	if err != nil {
//		log.Fatal(err)
//	}
//	assert.NoError(t, err)
//	assert.Equal(t, ash.Name, result.Name)
//}
//
//func Test_store_GetUser(t *testing.T) {
//	collection, _ := prepareCollection(usersCollectionName)
//	ash := models.User{
//		Name: "Peter",
//		Id:   whoID,
//	}
//	store := NewStore(collection, []*mongo.Collection{nil, nil, nil})
//	err := store.PutUser(&ash)
//	assert.NoError(t, err)
//	var result *models.User
//	result, err = store.GetUser(whoID)
//	assert.NoError(t, err)
//	assert.Equal(t, ash.Id, result.Id)
//	filter := bson.D{{"id", whoID}}
//	collection.DeleteOne(context.TODO(), filter)
//}
//
//func Test_store_UpdUserField(t *testing.T) {
//	collection, _ := prepareCollection(usersCollectionName)
//	ash := models.User{
//		Name: "Peter",
//		Id:   whoID,
//	}
//	store := NewStore(collection, []*mongo.Collection{nil, nil, nil})
//	err := store.PutUser(&ash)
//	assert.NoError(t, err)
//	var result *models.User
//	result, err = store.GetUser(whoID)
//	assert.NoError(t, err)
//	assert.Equal(t, ash.Id, result.Id)
//
//	ash.Name = "Pasha"
//	err = store.UpdUserField(ash.Id, "name", ash.Name)
//	assert.NoError(t, err)
//	err = store.DeleteUser(ash.Id)
//	assert.NoError(t, err)
//}
//
//func Test_store_PutLike(t *testing.T) {
//	collection, _ := prepareCollection(likesCollectionName)
//	st := NewStore(nil, []*mongo.Collection{collection, nil, nil})
//	t.Run(putNewLike, func(t *testing.T) {
//		err := st.PutLike(whoID, whomeID)
//		assert.NoError(t, err)
//		likes, err := st.GetLikes(whoID)
//		assert.NoError(t, err)
//		assert.Equal(t, int64(whoID), likes[0].Who)
//		filter := bson.D{
//			{
//				"who", whoID,
//			},
//		}
//		deleteRes, err := st.(*store).likesRegistry.(*registry).collection.DeleteOne(context.TODO(), filter)
//		assert.NoError(t, err)
//		assert.Equal(t, int64(1), deleteRes.DeletedCount)
//	})
//
//	t.Run(putNewLike, func(t *testing.T) {
//		err := st.PutLike(whoID, whomeID)
//		assert.NoError(t, err)
//		err = st.PutLike(whoID, anotherWhome)
//		assert.NoError(t, err)
//		likes, err := st.GetLikes(whoID)
//		var ids []int
//		for _, item := range likes {
//			ids = append(ids, int(item.Whome))
//		}
//		sort.Ints(ids)
//		i1 := sort.SearchInts(ids, whomeID)
//		assert.Equal(t, whomeID, ids[i1])
//		i2 := sort.SearchInts(ids, anotherWhome)
//		assert.Equal(t, anotherWhome, ids[i2])
//		filter := bson.D{
//			{
//				"who", whoID,
//			},
//		}
//		deleteRes, err := st.(*store).likesRegistry.(*registry).collection.DeleteOne(context.TODO(), filter)
//		assert.NoError(t, err)
//		assert.Equal(t, int64(1), deleteRes.DeletedCount)
//	})
//}
//
//func Test_store_GetAny(t *testing.T) {
//	usersCollection, _ := prepareCollection(usersCollectionName)
//	store := NewStore(usersCollection, []*mongo.Collection{nil, nil, nil})
//
//	t.Run(noUsersGetAny, func(t *testing.T) {
//		_, err := store.GetUser(whoID)
//		assert.Error(t, err)
//	})
//
//	t.Run(oneUserGetAny, func(t *testing.T) {
//		err := store.PutUser(&pasha)
//		assert.NoError(t, err)
//		user, err := store.GetUser(pasha.Id)
//		assert.NoError(t, err)
//		assert.Equal(t, pasha.Id, user.Id)
//		_, err = store.GetUser(whomeID)
//		assert.Error(t, err)
//		err = store.DeleteUser(pasha.Id)
//		assert.NoError(t, err)
//	})
//}
//
//func Test_store_GetAll(t *testing.T) {
//	usersCollection, _ := prepareCollection(usersCollectionName)
//	store := NewStore(usersCollection, []*mongo.Collection{nil, nil, nil})
//
//	t.Run(getAllUsersTest, func(t *testing.T) {
//		err := store.PutUser(&pasha)
//		assert.NoError(t, err)
//		err = store.PutUser(&ksyusha)
//		assert.NoError(t, err)
//		users, err := store.GetAllUsers()
//		assert.NoError(t, err)
//		assert.Equal(t, 2, len(users))
//		store.DeleteUser(pasha.Id)
//		store.DeleteUser(ksyusha.Id)
//	})
//}

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
