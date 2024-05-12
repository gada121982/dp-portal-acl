package api

import (
	"dp-portal-acl/internal/model"
	"strings"

	"dp-portal-acl/config"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type Router struct {
	config *config.Config
	model  *model.Model
}

func NewRouter(model *model.Model, config *config.Config) *Router {
	return &Router{
		model:  model,
		config: config,
	}
}

func (r *Router) registerACLRouter(aclRouter fiber.Router) {
	aclRouter.Post("/", r.CreateACL)
	aclRouter.Get("/list", r.GetAclList)
}

func (r *Router) Authentication(c *fiber.Ctx) error {
	token := ""
	bearerToken := c.Get("Authorization", "")

	if val := strings.Split(bearerToken, " "); len(val) > 1 {
		token = val[1]
	}

	if token != r.config.ServerConfig.SecretKey {
		return fiber.ErrUnauthorized
	}
	return c.Next()
}

func (r *Router) Start(addr string) {
	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	app.Use(helmet.New())
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))

	app.Use(r.Authentication)

	// app.Use(auth.Authentication

	api := app.Group("/api")
	aclRouter := api.Group("/acl")
	r.registerACLRouter(aclRouter)

	app.Listen(addr)
}
