package mongo

import (
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log/slog"
	"testing"
	"time"
)

func TestMongoInsert(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	monitor := &event.CommandMonitor{
		// 查询执行之前
		Started: func(ctx context.Context, startedEvent *event.CommandStartedEvent) {
			slog.Info("Started startedEvent Command", "cmd", startedEvent.Command)
		},
		Succeeded: func(ctx context.Context, succeededEvent *event.CommandSucceededEvent) {

		},
		Failed: func(ctx context.Context, failedEvent *event.CommandFailedEvent) {

		},
	}
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://root:root@localhost:27017").SetMonitor(monitor))
	//client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017").SetMonitor(monitor))
	assert.NoError(t, err)

	mgDb := client.Database("go-backend")
	coll := mgDb.Collection("user")

	res, err := coll.InsertOne(ctx, &User{
		ID:   1,
		Name: "Jannan",
		Age:  18,
	})
	slog.Info("InsertOne", "res's InsertedID", res.InsertedID)
}

func TestMongoFind(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	monitor := &event.CommandMonitor{
		// 查询执行之前
		Started: func(ctx context.Context, startedEvent *event.CommandStartedEvent) {
			slog.Info("Started startedEvent Command", "cmd", startedEvent.Command)
		},
		Succeeded: func(ctx context.Context, succeededEvent *event.CommandSucceededEvent) {

		},
		Failed: func(ctx context.Context, failedEvent *event.CommandFailedEvent) {

		},
	}
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://root:root@localhost:27017").SetMonitor(monitor))
	//client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017").SetMonitor(monitor))
	assert.NoError(t, err)

	mgDb := client.Database("go-backend")
	coll := mgDb.Collection("user")

	filter := bson.D{bson.E{Key: "id", Value: 1}} // id=1
	var user User
	err = coll.FindOne(ctx, filter).Decode(&user)
	require.NoError(t, err)
	fmt.Printf("%+v\n", user)

	user = User{}
	err = coll.FindOne(ctx, User{
		ID: 1, // 注意设置omitempty
	}).Decode(&user)
	if errors.Is(err, mongo.ErrNoDocuments) {
		assert.NoError(t, err)
	} else {
		require.NoError(t, err)
	}
	fmt.Printf("%+v\n", user)
}

func TestMongoUpdate(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	monitor := &event.CommandMonitor{
		// 查询执行之前
		Started: func(ctx context.Context, startedEvent *event.CommandStartedEvent) {
			slog.Info("Started startedEvent Command", "cmd", startedEvent.Command)
		},
		Succeeded: func(ctx context.Context, succeededEvent *event.CommandSucceededEvent) {

		},
		Failed: func(ctx context.Context, failedEvent *event.CommandFailedEvent) {

		},
	}
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://root:root@localhost:27017").SetMonitor(monitor))
	//client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017").SetMonitor(monitor))
	assert.NoError(t, err)

	mgDb := client.Database("go-backend")
	coll := mgDb.Collection("user")
	filter := bson.M{"id": 1}
	//filter := bson.D{bson.E{Key: "id", Value: 1}} // id=1
	//set := bson.D{bson.E{
	//	Key:   "$set",
	//	Value: bson.E{Key: "name", Value: "newName"},
	//}}
	set := bson.M{"$set": bson.M{"name": "newName"}}

	res, err := coll.UpdateMany(ctx, filter, set)
	require.NoError(t, err)
	fmt.Printf("%+v", res)

}

func TestMongoDel(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	monitor := &event.CommandMonitor{
		// 查询执行之前
		Started: func(ctx context.Context, startedEvent *event.CommandStartedEvent) {
			slog.Info("Started startedEvent Command", "cmd", startedEvent.Command)
		},
		Succeeded: func(ctx context.Context, succeededEvent *event.CommandSucceededEvent) {

		},
		Failed: func(ctx context.Context, failedEvent *event.CommandFailedEvent) {

		},
	}
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://root:root@localhost:27017").SetMonitor(monitor))
	//client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017").SetMonitor(monitor))
	assert.NoError(t, err)

	mgDb := client.Database("go-backend")
	coll := mgDb.Collection("user")
	filter := bson.M{"id": 1}
	_, err = coll.DeleteMany(ctx, filter)
	require.NoError(t, err)
}

type User struct {
	ID   int64  `bson:"id,omitempty"`
	Name string `bson:"name,omitempty"`
	Age  int    `bson:"age,omitempty"`
}
