package store

import (
	"context"
	"errors"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Store interface {
	Put(model *UserModel) error
	// CheckExists() bool
	GetAny() (*UserModel, error)
	GetBunch() (ret []*UserModel, err error)
}

type store struct {
	collection *mongo.Collection
}

func (s *store) Put(model *UserModel) error {
	_, err := s.collection.InsertOne(context.TODO(), *model)
	return err
}

func (s *store) GetAny() (*UserModel, error) {
	pipeline := []bson.D{bson.D{{"$sample", bson.D{{"size", 1}}}}}
	cur, err := s.collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}
	defer func() {
		e := cur.Close(context.TODO())
		if e != nil {
			log.Fatal(e)
		}
	}()
	for cur.Next(context.TODO()) {
		var user UserModel
		err = cur.Decode(&user)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		return &user, nil
	}
	return nil, errors.New("no object found")
}

func (s *store) GetBunch() (ret []*UserModel, err error) {
	pipeline := []bson.D{bson.D{{"$sample", bson.D{{"size", 5}}}}}
	cur, err := s.collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}
	defer func() {
		e := cur.Close(context.TODO())
		if e != nil {
			log.Fatal(e)
		}
	}()
	for cur.Next(context.TODO()) {
		var user UserModel
		err = cur.Decode(&user)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		ret = append(ret, &user)
	}
	return ret, nil
}

func NewStore(coll *mongo.Collection) Store {
	return &store{collection: coll}
}
