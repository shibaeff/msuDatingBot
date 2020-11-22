package store

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	collectionName = "collection"
	whoID          = 1
	whomeID        = 2
	anotherWhome   = 3

	simpleTest       = "simple test success"
	simpleDeleteTest = "simple delete test"
)

func Test_registry_AddToList(t *testing.T) {
	collection, _ := prepareCollection(collectionName)
	st := NewRegistry(collection)
	t.Run(simpleTest, func(t *testing.T) {
		err := st.AddToList(whoID, whomeID)
		assert.NoError(t, err)
		likes, err := st.GetList(whoID)
		assert.NoError(t, err)
		assert.Equal(t, int64(whoID), likes[0].Who)
		filter := bson.D{
			{
				"who", whoID,
			},
		}
		deleteRes, err := st.(*registry).collection.DeleteOne(context.TODO(), filter)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), deleteRes.DeletedCount)
	})

	t.Run(simpleTest, func(t *testing.T) {
		err := st.AddToList(whoID, whomeID)
		assert.NoError(t, err)
		err = st.AddToList(whoID, anotherWhome)
		assert.NoError(t, err)
		likes, err := st.GetList(whoID)
		assert.NoError(t, err)
		assert.Equal(t, int64(whoID), likes[0].Who)
		assert.Equal(t, int64(whomeID), likes[0].Whome)
		assert.Equal(t, int64(anotherWhome), likes[1].Whome)
		filter := bson.D{
			{
				"who", whoID,
			},
		}
		deleteRes, err := st.(*registry).collection.DeleteOne(context.TODO(), filter)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), deleteRes.DeletedCount)
	})
}

func TestRegistry_DeleteItem(t *testing.T) {
	collection, _ := prepareCollection(collectionName)
	reg := NewRegistry(collection)
	t.Run(simpleDeleteTest, func(t *testing.T) {
		err := reg.AddToList(1, 2)
		assert.NoError(t, err)
		entry, err := reg.GetList(1)
		assert.NoError(t, err)
		//assert.Equal(t, 1, int(entry[0].Who))
		//assert.Equal(t, 2, int(entry[0].Whome))
		err = reg.DeleteItems(1)
		assert.NoError(t, err)
		entry, err = reg.GetList(1)
		// assert.Error(t, err)
		assert.Nil(t, entry)
	})
}
