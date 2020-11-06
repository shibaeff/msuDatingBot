package store

import (
	"context"
	"errors"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"echoBot/pkg/bot"
)

type Store interface {
	Put(model *bot.User) error
	// CheckExists() bool
	GetAny() (int64, error)
	GetBunch() (ret []int64, err error)
}

type store struct {
	collection *mongo.Collection
}

func (s *store) Put(model *bot.User) error {
	_, err := s.collection.InsertOne(context.TODO(), *model)
	return err
}

func (s *store) GetAny() (int64, error) {
	panic("nothing")
}

func (s *store) GetBunch() (ret []int64, err error) {
	panic("nothing")
}

func NewStore(coll *mongo.Collection) Store {
	return &store{collection: coll}
}
