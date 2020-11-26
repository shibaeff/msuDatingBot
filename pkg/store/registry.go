package store

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Registry interface {
	AddToList(who, whome int64) error
	GetList(whose int64) ([]Entry, error)
	DeleteItems(int64) error
	DeleteItem(int64, int64) error
	IsPresent(who int64, whome int64) bool
}

type registry struct {
	collection *mongo.Collection
}

type Entry struct {
	Who   int64
	Whome int64
}

func (r *registry) AddToList(who, whome int64) (err error) {
	entry := Entry{who, whome}
	_, err = r.collection.InsertOne(context.TODO(), entry)
	return
}

func (r *registry) GetList(whose int64) (items []Entry, err error) {
	filter := bson.D{{"who", whose}}
	cur, err := r.collection.Find(context.TODO(), filter)
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

func (r *registry) IsPresent(who int64, whome int64) bool {
	filter := bson.D{{"who", who}, {"whome", whome}}
	cur := r.collection.FindOne(context.TODO(), filter)
	if cur == nil {
		return false
	}
	return true
}

func (r *registry) DeleteItems(whose int64) (err error) {
	filter := bson.D{{"who", whose}}
	_, err = r.collection.DeleteMany(context.TODO(), filter)
	return err
}

func (r *registry) DeleteItem(who, whome int64) (err error) {
	filter := bson.D{{"who", who}, {"whome", whome}}
	_, err = r.collection.DeleteMany(context.TODO(), filter)
	return err
}

func NewRegistry(collection *mongo.Collection) Registry {
	return &registry{collection: collection}
}
