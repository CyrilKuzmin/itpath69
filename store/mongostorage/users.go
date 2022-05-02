package mongostorage

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"time"

	"github.com/CyrilKuzmin/itpath69/internal/domain/users"
	"github.com/CyrilKuzmin/itpath69/store"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (m *MongoStorage) CreateUser(ctx context.Context, username, password string) (*users.User, error) {
	var user *users.User
	err := m.users.FindOne(ctx, bson.D{{"username", username}}).Decode(&user)
	if err == nil {
		return nil, store.ErrUserAlreadyExists(username)
	}
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return nil, store.ErrInternal(err)
		}
	}
	pwMd5 := md5.Sum([]byte(password))
	user = &users.User{
		Id:           uuid.New().String(),
		Username:     username,
		PasswordHash: hex.EncodeToString(pwMd5[:]),
		CreatedAt:    time.Now(),
		Modules:      map[int]users.ModuleProgress{},
	}
	_, err = m.users.InsertOne(ctx, user)
	if err != nil {
		return nil, store.ErrInternal(err)
	}
	return user, nil
}

func (m *MongoStorage) CheckUserPassword(ctx context.Context, username, password string) (*users.User, error) {
	var user *users.User
	pwMd5 := md5.Sum([]byte(password))
	err := m.users.FindOne(ctx, bson.D{
		{"username", username},
		{"passwordhash", hex.EncodeToString(pwMd5[:])},
	}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, store.ErrUserNotFound(username)
	}
	if err != nil {
		return nil, store.ErrInternal(err)
	}
	return user, nil
}

func (m *MongoStorage) GetUser(ctx context.Context, username string) (*users.User, error) {
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
