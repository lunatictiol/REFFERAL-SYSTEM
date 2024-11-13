package users

import (
	"context"
	"errors"

	"github.com/lunatictiol/referal-system/config"
	"github.com/lunatictiol/referal-system/models"
	"go.mongodb.org/mongo-driver/v2/bson"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UserStore struct {
	store *mongo.Client
}

func NewStore(store *mongo.Client) *UserStore {
	return &UserStore{store: store}
}

func (u *UserStore) Ping() {
	println("yayyyyy")
}
func (s *UserStore) CheckIfEmailExisits(email string) (models.User, error) {

	psa := s.store.Database("REFFERAL-SYSTEM").Collection("users")
	var user models.User
	err := psa.FindOne(context.Background(), bson.D{{Key: "email", Value: email}}).Decode(&user)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return user, nil
	} else if err != nil {
		println('1')
		config.LogFatal("Error finding user")
	}
	return user, nil

}

// register user
func (s *UserStore) RegisterUser(u models.RegisterUserPayload) (interface{}, error) {

	psa := s.store.Database("REFFERAL-SYSTEM").Collection("users")
	sp, err := psa.InsertOne(context.Background(), u)
	if err != nil {
		println('2')
		return bson.NilObjectID, err
	}
	return sp.InsertedID, nil

}
