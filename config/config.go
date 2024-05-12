package config

import (
	"github.com/Netflix/go-env"
	"github.com/joho/godotenv"
)

type ServerConfig struct {
	IpAddr    string `env:"IP_ADDR"`
	Port      int    `env:"PORT"`
	SecretKey string `env:"SECRET_KEY"`
}

type MongoConfig struct {
	URI      string `env:"MONGO_URI"`      // mongodb://<hostname>:<port>
	Username string `env:"MONGO_USERNAME"` // dp-portal
	Password string `env:"MONGO_PASSWORD"` // ********
	Database string `env:"MONGO_DATABASE"` // dp-portal
}

type Config struct {
	MongoConfig  *MongoConfig
	ServerConfig *ServerConfig
	Extras       env.EnvSet
}

func NewConfig() *Config {
	godotenv.Load()
	var err error
	var mongoConfig MongoConfig
	var serverConfig ServerConfig

	_, err = env.UnmarshalFromEnviron(&mongoConfig)
	if err != nil {
		panic(err)
	}

	_, err = env.UnmarshalFromEnviron(&serverConfig)
	if err != nil {
		panic(err)
	}

	return &Config{
		MongoConfig:  &mongoConfig,
		ServerConfig: &serverConfig,
	}
}
