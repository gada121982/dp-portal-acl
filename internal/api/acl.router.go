package api

import (
	"context"
	"dp-portal-acl/internal/model"

	"github.com/gofiber/fiber/v2"
)

type CreateACLRequest struct {
	UserId     string           `json:"user_id" xml:"user_id" form:"user_id"`
	Permission model.Permission `json:"permission" xml:"permission" form:"permission"`
}

func (r CreateACLRequest) Validate() error {
	if !model.PermissionMap[r.Permission] {
		return fiber.NewError(fiber.StatusBadRequest, "Permission not valid")
	}
	return nil
}

func (r *Router) CreateACL(c *fiber.Ctx) error {
	body := new(CreateACLRequest)
	if err := c.BodyParser(body); err != nil {
		return fiber.ErrBadRequest
	}
	if err := body.Validate(); err != nil {
		return err
	}
	acl := &model.CubeFSAclModel{
		UserId:     body.UserId,
		Permission: body.Permission,
		Type:       model.CubeFsACLType,
	}
	if err := r.model.CreateACL(context.Background(), acl); err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(acl)
}

type GetAclListResponse struct {
	UserId     string           `json:"user_id"`
	Permission model.Permission `json:"permission"`
}

func (r *Router) GetAclList(c *fiber.Ctx) error {
	resp := []GetAclListResponse{}
	aclType := c.QueryInt("type")

	if aclType == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "type not valid")
	}

	data, err := r.model.ListACL(context.TODO(), model.ACLType(aclType))
	if err != nil {
		return fiber.ErrInternalServerError
	}
	for _, val := range data {
		resp = append(resp, GetAclListResponse{
			UserId:     val.UserId,
			Permission: val.Permission,
		})
	}
	return c.Status(fiber.StatusOK).JSON(resp)
}
