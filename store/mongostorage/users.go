package mongostorage

import (
	"context"

	"github.com/CyrilKuzmin/itpath69/internal/service/users"
	"github.com/CyrilKuzmin/itpath69/store"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (m *MongoStorage) CreateUser(ctx context.Context, user *users.User) error {
	err := m.users.FindOne(ctx, bson.D{{"username", user.Username}}).Decode(&user)
	if err == nil {
		return store.ErrUserAlreadyExists(user.Username)
	}
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return store.ErrInternal(err)
		}
	}
	_, err = m.users.InsertOne(ctx, user)
	if err != nil {
		return store.ErrInternal(err)
	}
	return nil
}

func (m *MongoStorage) CheckUserPassword(ctx context.Context, username, passwordHash string) error {
	var user *users.User
	err := m.users.FindOne(ctx, bson.D{
		{"username", username},
		{"passwordhash", passwordHash},
	}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return store.ErrUserNotFound(username)
	}
	if err != nil {
		return store.ErrInternal(err)
	}
	return nil
}

func (m *MongoStorage) GetUserByName(ctx context.Context, username string) (*users.User, error) {
	var user *users.User
	err := m.users.FindOne(ctx, bson.D{
		{"username", username},
	}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, store.ErrUserNotFound(username)
	}
	if err != nil {
		return nil, store.ErrInternal(err)
	}
	return user, nil
}
