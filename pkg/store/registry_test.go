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

	simpleTest = "simple test success"
)

func Test_registry_AddToList(t *testing.T) {
	collection, _ := prepareCollection(collectionName)
	st := NewRegistry(collection)
	t.Run(simpleTest, func(t *testing.T) {
		err := st.AddToList(whoID, whomeID)
		assert.NoError(t, err)
		likes, err := st.GetList(whoID)
		assert.NoError(t, err)
		assert.Equal(t, int64(whoID), likes.Who)
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
		assert.Equal(t, int64(whoID), likes.Who)
		assert.Equal(t, int64(whomeID), likes.Whome[0])
		assert.Equal(t, int64(anotherWhome), likes.Whome[1])
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
