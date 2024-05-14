package api

import (
	"dp-portal-acl/internal/model"
	"strings"
	"time"

	"dp-portal-acl/config"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/golang-jwt/jwt/v5"

	jwtware "github.com/gofiber/contrib/jwt"
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
	aclRouter.Post("/", r.CreateACL).Name(CreateACLAction.String())
	aclRouter.Get("/list", r.GetAclList).Name(ListACLAction.String())
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

	app.Get("/auth/token", func(c *fiber.Ctx) error {
		secret := c.Query("secret_key", "")
		if secret == "" {
			return fiber.ErrBadRequest
		}
		claims := jwt.MapClaims{
			"permission": []string{CreateACLAction.String(), ListACLAction.String()},
			"exp":        time.Now().Add(time.Hour * 24 * time.Duration(r.config.ServerConfig.TokenExpireDay)).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		// Generate encoded token and send it as response.
		t, err := token.SignedString([]byte(secret))
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		return c.JSON(fiber.Map{
			"token": t,
		})
	})

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(r.config.ServerConfig.SecretKey)},
		SuccessHandler: func(c *fiber.Ctx) error {
			user := c.Locals("user").(*jwt.Token)
			claims := user.Claims.(jwt.MapClaims)

			permission := claims["permission"].([]interface{})
			next := c.Next()

			currentAction := c.Route().Name
			hasPermission := false
			for _, action := range permission {
				if currentAction == action.(string) {
					hasPermission = true
				}
			}
			if !hasPermission {
				return jwt.ErrTokenMalformed
			}
			return next
		},
	}))

	api := app.Group("/api")
	aclRouter := api.Group("/acl")

	r.registerACLRouter(aclRouter)
	app.Listen(addr)
}
