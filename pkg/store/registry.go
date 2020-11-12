package store

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Registry interface {
	AddToList(who, whome int64) error
	GetList(whose int64) (*Entry, error)
	DeleteItem(int642 int64) error
}

type registry struct {
	collection *mongo.Collection
}

type Entry struct {
	Who   int64
	Whome []int64
}

func (r *registry) AddToList(who, whome int64) (err error) {
	filter := bson.D{{"who", who}}
	change := bson.D{
		{
			"$push", bson.D{
				{"whome", whome},
			},
		},
	}
	upd, err := r.collection.UpdateOne(context.TODO(), filter, change)
	log.Printf("Update result %d", upd.ModifiedCount)
	if err != nil || upd.MatchedCount == 0 {
		_, err = r.collection.InsertOne(context.TODO(), Entry{
			Who:   who,
			Whome: []int64{whome},
		})
		return err
	}
	return
}

func (r *registry) GetList(whose int64) (item *Entry, err error) {
	filter := bson.D{{"who", whose}}
	res := r.collection.FindOne(context.TODO(), filter)
	item = new(Entry)
	err = res.Decode(item)
	if err != nil {
		return nil, err
	}
	return
}

func (r *registry) DeleteItem(whose int64) (err error) {
	filter := bson.D{{"who", whose}}
	_, err = r.collection.DeleteOne(context.TODO(), filter)
	return err
}

func NewRegistry(collection *mongo.Collection) Registry {
	return &registry{collection: collection}
}
