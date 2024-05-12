package main

import (
	"dp-portal-acl/config"
)

func main() {
	config := config.NewConfig()
	httpServer := NewServer(config)
	httpServer.Start()
}
