package mongodb

import (
	"GO-DFS/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
	"time"
)

func initClient(config *config.GlobalConfig) *mongo.Client {
	clientOptions := options.Client().ApplyURI(config.MongoDBAddr)
	var ctx = context.TODO()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		config.Logger.Error("MongoDB connect failed")
	}
	return client
}

func ExistDataBase(dateBaseName string, config *config.GlobalConfig) bool {

	client := initClient(config)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	databases, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		panic(err)
	}
	// 检查数据库是否存在
	dbName := dateBaseName
	exists := false
	for _, db := range databases {
		if db == dbName {
			exists = true
			break
		}
	}

	return exists

}
