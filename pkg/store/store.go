package store

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"

	"echoBot/pkg/models"
)

var (
	justOne = []bson.D{bson.D{{"$sample", bson.D{{"size", 1}}}}}
	many    = []bson.D{bson.D{{"$sample", bson.D{{"size", 10}}}}}
)

type Store interface {
	PutUser(model models.User) error
	GetUser(id int64) (*models.User, error)
	DeleteUser(id int64) error
	GetActions() Registry
	GetAllUsers() (ret []*models.User, err error)
	FindUser(filter bson.D) *models.User
	UpdUserField(id int64, field string, value interface{}) (err error)
}

type store struct {
	usersCollection *mongo.Collection
	registry        Registry
}

func (s *store) GetActions() Registry {
	return s.registry
}
func (s *store) PutUser(model models.User) error {
	_, err := s.usersCollection.InsertOne(context.TODO(), model)
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

func (s *store) FindUser(opt bson.D) *models.User {
	usr := s.usersCollection.FindOne(context.TODO(), opt)
	user := new(models.User)
	usr.Decode(user)
	return user
}
func (s *store) DeleteUser(id int64) (err error) {
	filter := bson.D{{"id", id}}
	_, err = s.usersCollection.DeleteOne(context.TODO(), filter)
	return err
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

func NewStore(users *mongo.Collection, registry *mongo.Collection) Store {
	return &store{
		usersCollection: users,
		registry:        NewRegistry(registry),
	}
}
