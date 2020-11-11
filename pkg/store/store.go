package store

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"echoBot/pkg/models"
)

var (
	justOne = []bson.D{bson.D{{"$sample", bson.D{{"size", 1}}}}}
)

type Store interface {
	PutUser(model *models.User) error
	GetUser(id int64) (*models.User, error)
	DeleteUser(id int64) error
	// CheckExists() bool
	PutLike(who int64, whome int64) error
	GetLikes(whose int64) (*Entry, error)
	PutSeen(who int64, whome int64) error
	GetSeen(whose int64) (*Entry, error)
	GetAny(for_id int64) (*models.User, error)
	GetBunch() (ret []int64, err error)
	UpdUserField(id int64, field string, value interface{}) (err error)
}

type store struct {
	usersCollection *mongo.Collection
	likesRegistry   Registry
	seenRegistry    Registry
}

func (s *store) PutUser(model *models.User) error {
	_, err := s.usersCollection.InsertOne(context.TODO(), *model)
	return err
}

func (s *store) GetUser(id int64) (user *models.User, err error) {
	filter := bson.D{{"id", id}}
	user = new(models.User)
	err = s.usersCollection.FindOne(context.TODO(), filter).Decode(user)
	return
}

func (s *store) UpdUserField(id int64, field string, value interface{}) (err error) {
	filter := bson.D{{"id", id}}
	pipeline := bson.D{
		{"$set", bson.D{{field, value}}},
	}
	res, err := s.usersCollection.UpdateOne(context.TODO(), filter, pipeline)
	if err == nil {
		log.Printf("modified %d documents\n", res.ModifiedCount)
	}
	return
}

func (s *store) DeleteUser(id int64) (err error) {
	filter := bson.D{{"id", id}}
	_, err = s.usersCollection.DeleteOne(context.TODO(), filter)
	return err
}

func (s *store) PutLike(who, whome int64) (err error) {
	err = s.likesRegistry.AddToList(who, whome)
	return
}

func (s *store) GetLikes(whose int64) (likes *Entry, err error) {
	likes, err = s.likesRegistry.GetList(whose)
	return
}

func (s *store) PutSeen(who, whome int64) (err error) {
	err = s.seenRegistry.AddToList(who, whome)
	return
}

func (s *store) GetSeen(whose int64) (seen *Entry, err error) {
	seen, err = s.seenRegistry.GetList(whose)
	return
}

func (s *store) GetAny(for_id int64) (user *models.User, err error) {
	user = new(models.User)
	err = s.usersCollection.FindOne(context.TODO(), justOne).Decode(user)
	if err != nil {
		return nil, err
	}
	return
}

func (s *store) GetBunch() (ret []int64, err error) {
	panic("nothing")
}

func NewStore(users *mongo.Collection, likes *mongo.Collection, seen *mongo.Collection) Store {
	return &store{
		usersCollection: users,
		likesRegistry:   NewRegistry(likes),
		seenRegistry:    NewRegistry(seen),
	}
}
