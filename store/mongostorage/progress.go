package mongostorage

import (
	"context"

	"github.com/CyrilKuzmin/itpath69/models"
	"github.com/CyrilKuzmin/itpath69/store"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (m *MongoStorage) OpenModules(ctx context.Context, username string, amount int) error {
	_, err := m.users.UpdateOne(ctx, bson.M{
		"username": username,
	}, bson.D{
		{"$inc", bson.D{{"modules_opened", amount}}},
	}, options.Update().SetUpsert(true))
	if err != nil {
		return store.ErrInternal(err)
	}
	return nil
}

func (m *MongoStorage) IncrementCompletedModules(ctx context.Context, username string) error {
	_, err := m.users.UpdateOne(ctx, bson.M{
		"username": username,
	}, bson.D{
		{"$inc", bson.D{{"modules_completed", 1}}},
	}, options.Update().SetUpsert(true))
	if err != nil {
		return store.ErrInternal(err)
	}
	return nil
}

func (m *MongoStorage) GetModulesMeta(ctx context.Context, amount int) ([]models.ModuleMeta, error) {
	res := make([]models.ModuleMeta, 0)
	opts := options.Find().SetProjection(bson.D{{"meta", 1}, {"_id", 0}})
	// {"_id", bson.D{{"$lte", amount}}
	cur, err := m.modules.Find(ctx, bson.D{{"_id", bson.D{{"$lte", amount}}}}, opts)
	if err != nil {
		return nil, store.ErrInternal(err)
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var m models.Module
		err := cur.Decode(&m)
		if err != nil {
			return nil, store.ErrInternal(err)
		}
		res = append(res, m.Meta)
	}
	if err := cur.Err(); err != nil {
		return nil, store.ErrInternal(err)
	}
	return res, nil
}
