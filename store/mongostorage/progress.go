package mongostorage

import (
	"context"
	"fmt"
	"strconv"

	"github.com/CyrilKuzmin/itpath69/models"
	"github.com/CyrilKuzmin/itpath69/store"
	"go.mongodb.org/mongo-driver/bson"
)

func (m *MongoStorage) UpdateProgress(ctx context.Context,
	username string, progress map[int]models.ModuleProgress) error {
	newProgress := make(bson.M)
	for i, pr := range progress {
		newProgress[strconv.Itoa(i)] = pr
	}
	fmt.Println(newProgress)
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

// func (m *MongoStorage) OpenModules(ctx context.Context, username string, amount int) error {
// 	// get current user
// 	var user *models.User
// 	err := m.users.FindOne(ctx, bson.D{
// 		{"username", username},
// 	}).Decode(&user)
// 	if err != nil {
// 		return store.ErrInternal(err)
// 	}
// 	// set current time to progress
// 	currTime := time.Now()
// 	newProgress := make(bson.M)
// 	for i, pr := range user.Modules {
// 		newProgress[strconv.Itoa(i)] = pr
// 	}
// 	for i := user.ModulesOpened + 1; i <= user.ModulesOpened+amount; i++ {
// 		newProgress[strconv.Itoa(i)] = models.ModuleProgress{CreatedAt: currTime}
// 	}
// 	_, err = m.users.UpdateOne(ctx,
// 		bson.M{"username": username},
// 		bson.D{
// 			{"$set", bson.M{"modules": newProgress}},
// 		})
// 	if err != nil {
// 		return store.ErrInternal(err)
// 	}
// 	return nil
// }

// func (m *MongoStorage) CompleteModule(ctx context.Context, username string, module int) error {
// 	// get current user
// 	var user *models.User
// 	err := m.users.FindOne(ctx, bson.D{
// 		{"username", username},
// 	}).Decode(&user)
// 	if err != nil {
// 		return store.ErrInternal(err)
// 	}
// 	// set current time as completion time
// 	user.Modules[module] = models.ModuleProgress{
// 		CreatedAt:   user.Modules[module].CreatedAt,
// 		CompletedAt: time.Now(),
// 	}
// 	newProgress := make(bson.M)
// 	for i, pr := range user.Modules {
// 		newProgress[strconv.Itoa(i)] = pr
// 	}
// 	_, err = m.users.UpdateOne(ctx,
// 		bson.M{"username": username},
// 		bson.D{
// 			{"$set", bson.M{"modules": newProgress}},
// 		})
// 	if err != nil {
// 		return store.ErrInternal(err)
// 	}
// 	// if 3 of 4 last modules are completed open new ones
// 	completedOnStage := 0
// 	for i := len(user.Modules) - 4; i < len(user.Modules); i++ {
// 		if !user.Modules[i].CompletedAt.IsZero() {
// 			completedOnStage++
// 		}
// 	}
// 	if completedOnStage > 2 {
// 		return m.OpenModules(ctx, username, 4)
// 	}
// 	return nil
// }
