package mongostorage

import (
	"context"

	"github.com/CyrilKuzmin/itpath69/models"
	"github.com/CyrilKuzmin/itpath69/store"
	"go.mongodb.org/mongo-driver/mongo"
)

func (m *MongoStorage) saveStuff(ctx context.Context, col *mongo.Collection, items []interface{}) error {
	col.Drop(ctx)
	_, err := col.InsertMany(ctx, items, nil)
	if err != nil {
		return store.ErrInternal(err)
	}
	return nil
}

func (m *MongoStorage) SaveModules(ctx context.Context, modules []models.Module) error {
	items := make([]interface{}, 0)
	for _, s := range modules {
		items = append(items, s)
	}
	return m.saveStuff(ctx, m.modules, items)
}
