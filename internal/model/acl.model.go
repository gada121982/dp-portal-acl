package model

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/net/context"
)

type ACLType int
type Permission int

var (
	CubeFsACLType ACLType = 1

	UserPermission    Permission = 1
	PartnerPermission Permission = 2

	PermissionMap = map[Permission]bool{
		UserPermission:    true,
		PartnerPermission: true,
	}
)

type AclBaseModel struct {
}

type CubeFSAclModel struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Type       ACLType            `bson:"type,omitempty"`
	UserId     string             `bson:"user_id,omitempty"` // cubefs user id
	Permission Permission         `bson:"permission,omitempty"`
}

func (m *Model) CreateACL(ctx context.Context, acl *CubeFSAclModel) error {
	acl.ID = primitive.NewObjectID()
	_, err := m.db.AclCollection.InsertOne(ctx, acl)
	if err != nil {
		return err
	}
	return nil
}

func (m *Model) ListACL(ctx context.Context, aclType ACLType) ([]CubeFSAclModel, error) {
	var data []CubeFSAclModel

	cursor, err := m.db.AclCollection.Find(ctx, bson.D{{Key: "type", Value: aclType}})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(ctx, &data); err != nil {
		return nil, err
	}
	return data, nil
}
