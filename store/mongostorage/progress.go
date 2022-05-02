package mongostorage

import (
	"context"
	"strconv"

	"github.com/CyrilKuzmin/itpath69/internal/domain/users"
	"github.com/CyrilKuzmin/itpath69/store"
	"go.mongodb.org/mongo-driver/bson"
)

func (m *MongoStorage) UpdateProgress(ctx context.Context,
	username string, progress map[int]users.ModuleProgress) error {
	newProgress := make(bson.M)
	for i, pr := range progress {
		newProgress[strconv.Itoa(i)] = pr
	}
	_, err := m.users.UpdateOne(ctx,
		bson.M{"username": username},
		bson.D{
			{"$set", bson.M{"modules": newProgress}},
		})
	if err != nil {
		return store.ErrInternal(err)
	}
	return nil
}
