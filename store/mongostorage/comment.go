package mongostorage

import (
	"context"
	"errors"

	"github.com/CyrilKuzmin/itpath69/internal/domain/comment"
	"github.com/CyrilKuzmin/itpath69/store"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (m *MongoStorage) CreateComment(ctx context.Context, c *comment.Comment) error {
	_, err := m.comments.InsertOne(ctx, c)
	if err != nil {
		return store.ErrInternal(err)
	}
	return nil
}

func (m *MongoStorage) UpdateComment(ctx context.Context, c *comment.Comment) error {
	var check *comment.Comment
	err := m.comments.FindOne(ctx, bson.D{{"_id", c.Id}}).Decode(check)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return m.CreateComment(ctx, c)
	}
	_, err = m.comments.UpdateOne(ctx,
		bson.M{"_id": c.Id},
		bson.D{
			{"$set", bson.M{
				"text":       c.Text,
				"modifiedat": c.ModifiedAt,
			}},
		})
	if err != nil {
		return store.ErrInternal(err)
	}
	return nil
}

func (m *MongoStorage) DeleteCommentByID(ctx context.Context, commentId string) error {
	_, err := m.comments.DeleteOne(ctx, bson.M{"_id": commentId})
	if err != nil {
		return store.ErrInternal(err)
	}
	return nil
}

func (m *MongoStorage) ListCommentsByModule(ctx context.Context, user string, module int) ([]*comment.Comment, error) {
	res := make([]*comment.Comment, 0)
	cur, err := m.comments.Find(ctx, bson.M{
		"user":     user,
		"moduleid": module,
	})
	if err != nil {
		return nil, store.ErrInternal(err)
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var c comment.Comment
		err := cur.Decode(&c)
		if err != nil {
			return nil, store.ErrInternal(err)
		}
		res = append(res, &c)
	}
	if err := cur.Err(); err != nil {
		return nil, store.ErrInternal(err)
	}
	return res, nil
}

func (m *MongoStorage) GetCommentByID(ctx context.Context, user, commentId string) (*comment.Comment, error) {
	res := comment.Comment{}
	err := m.comments.FindOne(ctx, bson.M{
		"_id":  commentId,
		"user": user,
	}).Decode(&res)
	if err != nil {
		return nil, store.ErrInternal(err)
	}
	return &res, nil
}
