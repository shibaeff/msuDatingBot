package store

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"echoBot/pkg/bot"
)

type Store interface {
	Put(model *bot.User) error
	// CheckExists() bool
	PutLike(who int64, whome int64) error
	GetAny(for_id int64) (int64, error)
	GetBunch() (ret []int64, err error)
}

type store struct {
	usersCollection *mongo.Collection
	likesCollection *mongo.Collection
	seenCollection  *mongo.Collection
}

func (s *store) Put(model *bot.User) error {
	_, err := s.usersCollection.InsertOne(context.TODO(), *model)
	return err
}

func (s *store) PutLike(who, whome int64) (err error) {
	filter := bson.D{{"name", "Peter"}}
	change := bson.D{
		{
			"$push", bson.D{
				{"likes", 1},
			},
		},
	}
	_, err = s.likesCollection.UpdateOne(context.TODO(), filter, change)
	if err != nil {
		_, err = s.likesCollection.InsertOne(context.TODO(), Likes{
			who:   who,
			whome: []int64{whome},
		})
		return err
	}
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
		likesCollection: likes,
		seenCollection:  seen,
	}
}
