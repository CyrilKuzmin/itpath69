package mongostorage

import (
	"context"

	"github.com/CyrilKuzmin/itpath69/internal/domain/course"
	"go.mongodb.org/mongo-driver/bson"
)

func (m *MongoStorage) CreateCourse(ctx context.Context, c *course.Course) error {
	return nil
}
func (m *MongoStorage) UpdateCourse(ctx context.Context, c *course.Course) error {
	return nil
}
func (m *MongoStorage) DeleteCourse(ctx context.Context, id string) error {
	return nil
}

func (m *MongoStorage) GetCourseByID(ctx context.Context, id string) (*course.Course, error) {
	var res course.Course
	err := m.courses.FindOne(ctx, bson.D{
		{"_id", id},
	}).Decode(&res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (m *MongoStorage) ListCourses(ctx context.Context) ([]string, error) {
	return nil, nil
}
func (m *MongoStorage) ListCoursesByOwner(ctx context.Context, userId string) ([]string, error) {
	return nil, nil
}
func (m *MongoStorage) MakePrivate(ctx context.Context, id string) error {
	return nil
}
func (m *MongoStorage) MakePublic(ctx context.Context, id string) error {
	return nil
}
func (m *MongoStorage) Publish(ctx context.Context, id string) error {
	return nil
}
func (m *MongoStorage) AddOwner(ctx context.Context, id, userId string) error {
	return nil
}
