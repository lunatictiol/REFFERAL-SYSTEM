package users

import "go.mongodb.org/mongo-driver/v2/mongo"

type UserStore struct {
	store *mongo.Client
}

func NewStore(store *mongo.Client) *UserStore {
	return &UserStore{store: store}
}

func (u *UserStore) Ping() {
	println("yayyyyy")
}
