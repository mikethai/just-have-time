package userHandler

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/mikethai/just-have-time/database"
	"github.com/mikethai/just-have-time/internal/model"
)

type Handler struct {
	repository Repository
	httpClient HttpClient
}

func NewHandler() *Handler {
	r := NewRepository(database.DB)
	c := NewHttpClient(&http.Client{})

	return &Handler{
		repository: r,
		httpClient: c,
	}
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

func (h *Handler) SyncUserFollow(c *fiber.Ctx) error {
	userInfo := new(struct {
		Msno int64
		SID  string
	})

	if err := c.BodyParser(&userInfo); err != nil {
		return c.Status(503).JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})
	}

	followee, err := h.httpClient.GetUserFollowing(userInfo.Msno, userInfo.SID)
	if err != nil {
		fmt.Println(err)
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "The request has some wrong.", "data": err.Error()})
	}

	followUser := &CreateUserParameter{
		Msno: userInfo.Msno,
	}
	h.repository.Create(followUser)

	for _, followeeUser := range *followee {

		h.repository.Create(&CreateUserParameter{
			Msno: followeeUser.Msno,
		})

		newFollowModel := model.Follow{
			FollowerID: userInfo.Msno,
			FolloweeID: followeeUser.Msno,
		}

		h.repository.Follow(&FollowParameter{
			followModel: newFollowModel,
		})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Created Follow"})
}

func (h *Handler) CreateUser(msno int64) {
	user := &CreateUserParameter{
		Msno: msno,
	}
	h.repository.Create(user)
}
