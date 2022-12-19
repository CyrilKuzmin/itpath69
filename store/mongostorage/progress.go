package mongostorage

import (
	"context"
	"time"

	"github.com/CyrilKuzmin/itpath69/internal/service/progress"
	"github.com/CyrilKuzmin/itpath69/store"
	"go.mongodb.org/mongo-driver/bson"
)

func (m *MongoStorage) CreateCourseProgress(ctx context.Context, data []progress.ModuleProgress) error {
	items := make([]interface{}, len(data))
	for i, s := range data {
		items[i] = s
	}
	_, err := m.progress.InsertMany(ctx, items, nil)
	if err != nil {
		return store.ErrInternal(err)
	}
	return nil
}

func (m *MongoStorage) GetUserProgress(ctx context.Context, userId, courseId string) ([]progress.ModuleProgress, error) {
	res := make([]progress.ModuleProgress, 0)
	// {"_id", bson.D{{"$lte", amount}}
	cur, err := m.progress.Find(ctx, bson.M{
		"userid":   userId,
		"courseid": courseId,
	}, nil)
	if err != nil {
		return nil, store.ErrInternal(err)
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var m progress.ModuleProgress
		err := cur.Decode(&m)
		if err != nil {
			return nil, store.ErrInternal(err)
		}
		res = append(res, m)
	}
	if err := cur.Err(); err != nil {
		return nil, store.ErrInternal(err)
	}
	return res, nil
}

func (m *MongoStorage) GetModuleProgress(ctx context.Context, userId, courseId string, moduleId int) (progress.ModuleProgress, error) {
	var res progress.ModuleProgress
	err := m.progress.FindOne(ctx, bson.M{
		"userid":   userId,
		"courseid": courseId,
		"moduleid": moduleId,
	}, nil).Decode(&res)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (m *MongoStorage) MarkModuleAsCompleted(ctx context.Context, id string) error {
	_, err := m.progress.UpdateOne(ctx,
		bson.M{"_id": id},
		bson.D{
			{"$set", bson.M{
				"completedat": time.Now(),
			}},
		})

	if err != nil {
		return store.ErrInternal(err)
	}
	return nil
}

func (m *MongoStorage) MarkModulesAsOpen(ctx context.Context, toOpenIDs []string) error {
	now := time.Now()
	_, err := m.progress.UpdateMany(ctx,
		bson.M{"_id": bson.M{"$in": toOpenIDs}},
		bson.D{
			{"$set", bson.M{
				"openedat": now,
			}},
		},
		nil)
	if err != nil {
		return store.ErrInternal(err)
	}
	return nil
}
