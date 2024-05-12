package db

import (
	"context"
	"dp-portal-acl/config"

	mongoPkg "dp-portal-acl/pkg/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const ACLCollection = "user_acl"

type Database struct {
	AclCollection *mongo.Collection
}

func (db *Database) CreateIndex() {
	db.AclCollection.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys:    bson.D{{Key: "user_id", Value: -1}},
		Options: options.Index().SetUnique(true),
	})

}

func NewDatabase(c *config.MongoConfig) (*Database, error) {
	client, err := mongoPkg.NewClient(c)
	if err != nil {
		return nil, err
	}
	db := client.Database(c.Database)

	if err := db.Client().Ping(context.Background(), nil); err != nil {
		panic(err)
	}
	aclColl := db.Collection(ACLCollection)

	return &Database{
		AclCollection: aclColl,
	}, nil
}
