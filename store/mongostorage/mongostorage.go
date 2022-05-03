package mongostorage

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type MongoStorage struct {
	client   *mongo.Client
	ctx      context.Context
	log      *zap.Logger
	users    *mongo.Collection
	modules  *mongo.Collection
	comments *mongo.Collection
	// Sessions used as a session store, should be imported
	Sessions *mongo.Collection
}

func NewMongo(log *zap.Logger, uri, database string) *MongoStorage {
	clientOptions := options.Client().ApplyURI(uri)
	ctx := context.Background()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("cannot init mongo connection", zap.Error(err))
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("cannot ping mongo server", zap.Error(err))
	}
	users := client.Database(database).Collection("users")
	modules := client.Database(database).Collection("modules")
	comments := client.Database(database).Collection("comments")
	sess := client.Database(database).Collection("sessions")
	return &MongoStorage{
		client:   client,
		ctx:      ctx,
		log:      log,
		users:    users,
		comments: comments,
		modules:  modules,
		Sessions: sess,
	}
}

func (m *MongoStorage) Close(ctx context.Context) {
	m.client.Disconnect(ctx)
}
