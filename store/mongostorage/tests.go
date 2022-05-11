package mongostorage

import (
	"context"

	"github.com/CyrilKuzmin/itpath69/internal/domain/tests"
	"github.com/CyrilKuzmin/itpath69/store"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (m *MongoStorage) GetTestByID(ctx context.Context, id string) (*tests.Test, error) {
	var res *tests.Test
	err := m.tests.FindOne(ctx, bson.D{
		{"_id", id},
	}).Decode(&res)
	if err == mongo.ErrNoDocuments {
		return res, store.ErrTestNotFound(id)
	}
	if err != nil {
		return res, store.ErrInternal(err)
	}
	return res, nil
}
func (m *MongoStorage) GetTestsByUser(ctx context.Context, userId string) ([]*tests.Test, error) {
	res := make([]*tests.Test, 0)
	cur, err := m.tests.Find(ctx, bson.D{{"userid", userId}})
	if err == mongo.ErrNoDocuments {
		return res, nil
	}
	for cur.Next(ctx) {
		var t *tests.Test
		err := cur.Decode(&t)
		if err != nil {
			return nil, store.ErrInternal(err)
		}
		res = append(res, t)
	}
	return res, nil
}
func (m *MongoStorage) SaveTest(ctx context.Context, test *tests.Test) error {
	_, err := m.tests.InsertOne(ctx, test)
	if err != nil {
		return store.ErrInternal(err)
	}
	return nil
}

// Questions
// GetModuleQuestions returns given amount of questions for module. It uses MongoDB aggregation
func (m *MongoStorage) GetModuleQuestions(ctx context.Context, moduleId int, amount int) ([]*tests.Question, error) {
	qs := make([]*tests.Question, 0)
	matchStage := bson.D{{"$match", bson.D{{"moduleid", moduleId}}}}
	amountStage := bson.D{{"$sample", bson.D{{"size", amount}}}}
	cur, err := m.questions.Aggregate(ctx, mongo.Pipeline{matchStage, amountStage})
	defer cur.Close(ctx)
	if err != nil {
		return nil, err
	}
	for cur.Next(ctx) {
		var q *tests.Question
		err := cur.Decode(&q)
		if err != nil {
			return nil, store.ErrInternal(err)
		}
		qs = append(qs, q)
	}
	return qs, nil
}

// Content Manager method
func (m *MongoStorage) SaveQuestions(ctx context.Context, qs []tests.Question) error {
	items := make([]interface{}, 0)
	for _, q := range qs {
		items = append(items, q)
	}
	return m.saveStuff(ctx, m.questions, items)
}
