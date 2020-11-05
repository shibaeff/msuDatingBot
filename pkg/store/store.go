package store

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

const (
	sampleSize = 5
)

type Store interface {
	Put(model *UserModel) error
	// CheckExists() bool
	GetAny() *UserModel
	GetBunch() []*UserModel
}

type store struct {
	collection *mongo.Collection
}

func (s *store) Put(model *UserModel) error {
	_, err := s.collection.InsertOne(context.TODO(), *model)
	return err
}

func (s *store) GetAny() *UserModel {
	s.collection.CountDocuments()
}

func (s *store) GetBunch() []*UserModel {
	panic("implement me")
}

func NewStore(coll *mongo.Collection) Store {
	return &store{collection: coll}
}
