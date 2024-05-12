package main

import (
	"dp-portal-acl/config"
	"dp-portal-acl/internal/api"
	"dp-portal-acl/internal/db"
	"dp-portal-acl/internal/model"
	"fmt"
)

type Server struct {
	config *config.Config
	router *api.Router
}

func NewServer(config *config.Config) *Server {
	db, err := db.NewDatabase(config.MongoConfig)
	if err != nil {
		panic(err)
	}
	db.CreateIndex()

	model := model.NewModel(db)
	router := api.NewRouter(model, config)
	return &Server{config: config, router: router}
}

func (s *Server) Start() {
	ip := s.config.ServerConfig.IpAddr
	port := s.config.ServerConfig.Port

	addr := fmt.Sprintf(`%s:%d`, ip, port)
	s.router.Start(addr)
}
