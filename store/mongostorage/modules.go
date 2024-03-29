package mongostorage

import (
	"context"

	"github.com/CyrilKuzmin/itpath69/internal/service/module"
	"github.com/CyrilKuzmin/itpath69/store"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (m *MongoStorage) saveStuff(ctx context.Context, col *mongo.Collection, items []interface{}) error {
	col.Drop(ctx)
	_, err := col.InsertMany(ctx, items, nil)
	if err != nil {
		return store.ErrInternal(err)
	}
	return nil
}

func (m *MongoStorage) SaveModules(ctx context.Context, modules []module.Module) error {
	items := make([]interface{}, 0)
	for _, s := range modules {
		items = append(items, s)
	}
	return m.saveStuff(ctx, m.modules, items)
}

func (m *MongoStorage) GetModulesMeta(ctx context.Context, courseId string, amount int) ([]module.ModuleMeta, error) {
	res := make([]module.ModuleMeta, 0)
	opts := options.Find().SetProjection(bson.D{{"meta", 1}, {"_id", 0}})
	// {"_id", bson.D{{"$lte", amount}}
	cur, err := m.modules.Find(ctx, bson.D{
		{"courseid", courseId},
		{"_id", bson.D{{"$lte", amount}}},
	}, opts)
	if err != nil {
		return nil, store.ErrInternal(err)
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var m module.Module
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

func (m *MongoStorage) GetModule(ctx context.Context, courseId string, id int) (module.Module, error) {
	var res module.Module
	err := m.modules.FindOne(ctx, bson.D{
		{"_id", id},
		{"courseid", courseId},
	}).Decode(&res)
	if err == mongo.ErrNoDocuments {
		return res, store.ErrModuleNotFound(id)
	}
	if err != nil {
		return res, store.ErrInternal(err)
	}
	return res, nil
}
