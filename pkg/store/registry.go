package store

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	EventLike    = "like"
	EventMatch   = "match"
	EventView    = "view"
	EventDislike = "dislike"
	EventUseen   = "unseen"
)

type Options []bson.E
type Registry interface {
	AddEvent(Entry) error
	GetEvents([]bson.E) ([]Entry, error)
	DeleteEvents([]bson.E) error
}

type registry struct {
	collection *mongo.Collection
}

type Entry struct {
	Who   int64
	Whome int64
	Event string
}

func (r *registry) AddEvent(e Entry) (err error) {
	_, err = r.collection.InsertOne(context.TODO(), e)
	return
}

func (r *registry) GetEvents(options []bson.E) (items []Entry, err error) {
	cur, err := r.collection.Find(context.TODO(), options)
	if err != nil {
		return
	}
	for cur.Next(context.TODO()) {
		item := new(Entry)
		err = cur.Decode(item)
		if err != nil {
			return nil, err
		}
		items = append(items, *item)
	}
	return
}

func (r *registry) IsPresent(options []bson.E) bool {
	cur := r.collection.FindOne(context.TODO(), options)
	if cur == nil {
		return false
	}
	var item Entry
	err := cur.Decode(&item)
	return err != nil
}

func (r *registry) DeleteEvents(options []bson.E) (err error) {
	_, err = r.collection.DeleteMany(context.TODO(), options)
	return err
}

func NewRegistry(collection *mongo.Collection) Registry {
	return &registry{collection: collection}
}
