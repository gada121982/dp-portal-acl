package model

import (
	"dp-portal-acl/internal/db"
)

type Model struct {
	db *db.Database
}

func NewModel(db *db.Database) *Model {
	return &Model{
		db: db,
	}
}
