package store

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"

	"echoBot/pkg/bot"
)

type Store interface {
	Put(model *bot.User) error
	// CheckExists() bool
	PutLike(who int64, whome int64) error
	GetLikes(whose int64) (*Entry, error)
	PutSeen(who int64, whome int64) error
	GetSeen(whose int64) (*Entry, error)
	GetAny(for_id int64) (int64, error)
	GetBunch() (ret []int64, err error)
}

type store struct {
	usersCollection *mongo.Collection
	likesRegistry   Registry
	seenRegistry    Registry
}

func (s *store) Put(model *bot.User) error {
	_, err := s.usersCollection.InsertOne(context.TODO(), *model)
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

func (s *store) GetAny(for_id int64) (int64, error) {
	panic("nothing")
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
