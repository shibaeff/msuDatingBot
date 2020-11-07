package store

//
//import (
//	"context"
//	"log"
//
//	"go.mongodb.org/mongo-driver/bson"
//	"go.mongodb.org/mongo-driver/mongo"
//
//	"echoBot/pkg/bot"
//)
//
//type Store interface {
//	Put(model *bot.User) error
//	// CheckExists() bool
//	PutLike(who int64, whome int64) error
//	GetLikes(whose int64) (*Likes, error)
//	PutSeen(who int64, whome int64) error
//	GetSeen(whose int64) (*Seen, error)
//	GetAny(for_id int64) (int64, error)
//	GetBunch() (ret []int64, err error)
//}
//
//type store struct {
//	usersCollection *mongo.Collection
//	likesCollection *mongo.Collection
//	seenCollection  *mongo.Collection
//}
//
//func (s *store) Put(model *bot.User) error {
//	_, err := s.usersCollection.InsertOne(context.TODO(), *model)
//	return err
//}
//
//func (s *store) PutLike(who, whome int64) (err error) {
//	filter := bson.D{{"who", who}}
//	change := bson.D{
//		{
//			"$push", bson.D{
//				{"whome", whome},
//			},
//		},
//	}
//	upd, err := s.likesCollection.UpdateOne(context.TODO(), filter, change)
//	log.Printf("Update result %d", upd.ModifiedCount)
//	if err != nil || upd.MatchedCount == 0 {
//		_, err = s.likesCollection.InsertOne(context.TODO(), Likes{
//			Who:   who,
//			Whome: []int64{whome},
//		})
//		return err
//	}
//	return
//}
//
//func (s *store) GetLikes(whose int64) (likes *Likes, err error) {
//	filter := bson.D{{"who", whose}}
//	res := s.likesCollection.FindOne(context.TODO(), filter)
//	likes = new(Likes)
//	err = res.Decode(likes)
//	if err != nil {
//		return nil, err
//	}
//	return
//}
//
//func (s *store) PutSeen(who, whome int64) (err error) {
//	filter := bson.D{{"who", who}}
//	change := bson.D{
//		{
//			"$push", bson.D{
//			{"whome", whome},
//		},
//		},
//	}
//	upd, err := s.seenCollection.UpdateOne(context.TODO(), filter, change)
//	log.Printf("Update result %d", upd.ModifiedCount)
//	if err != nil || upd.MatchedCount == 0 {
//		_, err = s.seenCollection.InsertOne(context.TODO(), Seen{
//			Who:   who,
//			Whome: []int64{whome},
//		})
//		return err
//	}
//	return
//}
//
//func (s *store) GetSeen(whose int64) (likes *Seen, err error) {
//	filter := bson.D{{"who", whose}}
//	res := s.seenCollection.FindOne(context.TODO(), filter)
//	likes = new(Seen)
//	err = res.Decode(likes)
//	if err != nil {
//		return nil, err
//	}
//	return
//}
//
//func (s *store) GetAny(for_id int64) (int64, error) {
//	panic("nothing")
//}
//
//func (s *store) GetBunch() (ret []int64, err error) {
//	panic("nothing")
//}
//
//func NewStore(users *mongo.Collection, likes *mongo.Collection, seen *mongo.Collection) Store {
//	return &store{
//		usersCollection: users,
//		likesCollection: likes,
//		seenCollection:  seen,
//	}
//}
