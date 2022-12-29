package userHandler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mikethai/just-have-time/database"
	"github.com/mikethai/just-have-time/internal/model"
)

type Handler struct {
	repository Repository
}

func NewHandler() *Handler {
	r := NewRepository(database.DB)

	return &Handler{repository: r}
}

type UserFollow struct {
	Follower int64
	Followee int64
}

func (h *Handler) CreateUserFollow(c *fiber.Ctx) error {
	userFollow := new(UserFollow)

	if err := c.BodyParser(&userFollow); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})
	}

	followUser := &CreateUserParameter{
		Msno: userFollow.Follower,
	}
	h.repository.Create(followUser)

	followeeUser := &CreateUserParameter{
		Msno: userFollow.Followee,
	}
	h.repository.Create(followeeUser)

	newFollowModel := model.Follow{
		FollowerID: userFollow.Follower,
		FolloweeID: userFollow.Followee,
	}

	followParameter := &FollowParameter{
		followModel: newFollowModel,
	}

	h.repository.Follow(followParameter)

	// Return the user follow
	return c.JSON(fiber.Map{"status": "success", "message": "Created Follow", "data": userFollow})
}
