package mongo

import (
	"context"
	"dp-portal-acl/config"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewClient(config *config.MongoConfig) (*mongo.Client, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	credential := options.Credential{
		AuthSource: config.Database,
		Username:   config.Username,
		Password:   config.Password,
	}
	clientOpts := options.Client().ApplyURI(config.URI).SetAuth(credential)
	client, err := mongo.Connect(ctx, clientOpts)
	return client, err
}
