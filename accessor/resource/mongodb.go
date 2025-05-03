package resource

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"llm_inference_web/accessor/confparser"
)

var MongoDB *mongo.Database

func initMongoDB(ctx context.Context) {
	uri := fmt.Sprintf(
		"mongodb://%s:%s@%s:%d/",
		confparser.ResourceConfig.MongoDB.User, confparser.ResourceConfig.MongoDB.Password,
		confparser.ResourceConfig.MongoDB.Host, confparser.ResourceConfig.MongoDB.Port,
	)
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		panic(err)
	}
	MongoDB = client.Database(confparser.ResourceConfig.MongoDB.DBName)
	if MongoDB == nil {
		panic("MongoDB is nil")
	}
}
