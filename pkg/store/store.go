package store

import (
	"context"
	"errors"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"echoBot/pkg/models"
)

var (
	justOne = []bson.D{bson.D{{"$sample", bson.D{{"size", 1}}}}}
	many    = []bson.D{bson.D{{"$sample", bson.D{{"size", 10}}}}}
)

type Store interface {
	PutUser(model *models.User) error
	GetUser(id int64) (*models.User, error)
	DeleteUser(id int64) error
	// CheckExists() bool
	PutLike(who int64, whome int64) error
	GetLikes(whose int64) ([]Entry, error)
	PutUnseen(who int64, whome int64) error
	GetUnseen(whose int64) ([]Entry, error)
	GetSeen(whose int64) ([]Entry, error)
	GetAny(for_id int64) (*models.User, error)
	GetBunch(n int) (ret []*models.User, err error)
	GetMatchesRegistry() Registry
	UpdUserField(id int64, field string, value interface{}) (err error)
	DeleteFromRegistires(id int64) error
	GetAllUsers() (ret []*models.User, err error)
	GetUnseenRegistry() Registry
	GetSeenRegistry() Registry
	GetAdmin(username string) (user *models.User, err error)
}

type store struct {
	usersCollection *mongo.Collection
	likesRegistry   Registry
	unseenRegistry  Registry
	matchesRegistry Registry
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

func (s *store) GetAdmin(username string) (user *models.User, err error) {
	filter := bson.D{{"username", username}}
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

func (s *store) GetLikes(whose int64) (likes []Entry, err error) {
	likes, err = s.likesRegistry.GetList(whose)
	return
}

func (s *store) PutUnseen(who, whome int64) (err error) {
	err = s.unseenRegistry.AddToList(who, whome)
	return
}

func (s *store) GetUnseen(whose int64) (seen []Entry, err error) {
	seen, err = s.unseenRegistry.GetList(whose)
	return
}

func (s *store) GetSeen(whose int64) (seen []Entry, err error) {
	seen, err = s.seenRegistry.GetList(whose)
	return
}

func (s *store) GetAny(for_id int64) (user *models.User, err error) {
	users, err := s.GetBunch(5)
	for i, user := range users {
		if user.Id == for_id {
			users = remove(users, i)
		}
	}
	if err != nil {
		return
	}
	if len(users) > 0 {
		user = users[0]
		return
	}
	return nil, errors.New("no users")
}

func (s *store) GetAllUsers() (ret []*models.User, err error) {
	empty := bson.D{}
	cur, err := s.usersCollection.Find(context.TODO(), empty)
	if err != nil {
		return
	}
	for cur.Next(context.TODO()) {
		user := new(models.User)
		if err = cur.Decode(user); err != nil {
			return nil, err
		}
		ret = append(ret, user)
	}
	return
}

func (s *store) GetBunch(n int) (ret []*models.User, err error) {
	many = []bson.D{bson.D{{"$sample", bson.D{{"size", n}}}}}
	cur, err := s.usersCollection.Aggregate(context.TODO(), many)
	if err != nil {
		return
	}
	for cur.Next(context.TODO()) {
		user := new(models.User)
		if err = cur.Decode(user); err != nil {
			return nil, err
		}
		ret = append(ret, user)
	}
	return
}

func (s *store) GetMatchesRegistry() Registry {
	return s.matchesRegistry
}

func (s *store) DeleteFromRegistires(id int64) (err error) {
	s.matchesRegistry.DeleteItems(id)
	s.likesRegistry.DeleteItems(id)
	s.unseenRegistry.DeleteItems(id)
	return nil
}

func (s *store) GetUnseenRegistry() Registry {
	return s.unseenRegistry
}

func (s *store) GetSeenRegistry() Registry {
	return s.seenRegistry
}

func remove(s []*models.User, i int) []*models.User {
	s[i] = s[len(s)-1]
	// We do not need to put s[i] at the end, as it will be discarded anyway
	return s[:len(s)-1]
}

func NewStore(users *mongo.Collection, registries []*mongo.Collection) Store {
	return &store{
		usersCollection: users,
		likesRegistry:   NewRegistry(registries[0]),
		unseenRegistry:  NewRegistry(registries[1]),
		matchesRegistry: NewRegistry(registries[2]),
		seenRegistry:    NewRegistry(registries[3]),
	}
}
