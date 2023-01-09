package userHandler

import (
	"net/http"

	"github.com/go-playground/validator/v10"
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
	Follower string `validate:"required"`
	Followee string `validate:"required"`
	SID      string `validate:"required"`
}

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

var validate = validator.New()

func ValidateStruct(userFollow UserFollow) []*ErrorResponse {
	var errors []*ErrorResponse
	err := validate.Struct(userFollow)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}

func (h *Handler) SyncUserFollow(c *fiber.Ctx) error {
	userInfo := new(struct {
		Msno string
		Name string
		SID  string
	})

	if err := c.BodyParser(&userInfo); err != nil {
		return c.Status(503).JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})
	}
	userProfile, err := h.httpClient.GetUserProfile(userInfo.Msno, userInfo.SID)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "The request has some wrong.", "data": err.Error()})
	}

	followee, err := h.httpClient.GetUserFollowing(userInfo.Msno, userInfo.SID)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "The request has some wrong.", "data": err.Error()})
	}
	followerMsnoInt, _ := h.httpClient.DecryptMsno(userInfo.Msno, userInfo.SID)
	followUser := &CreateUserParameter{
		Msno:     userInfo.Msno,
		MsnoInt:  *followerMsnoInt,
		UserName: userProfile.Name,
	}
	h.repository.Create(followUser)

	for _, followeeUser := range *followee {

		followeeMsnoInt, _ := h.httpClient.DecryptMsno(followeeUser.Msno, userInfo.SID)

		h.repository.Create(&CreateUserParameter{
			Msno:     followeeUser.Msno,
			MsnoInt:  *followeeMsnoInt,
			UserName: followeeUser.Name,
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

func (h *Handler) FetchUser(msno string) (*model.User, error) {
	user := &FetchUserParameter{
		Msno: msno,
	}
	userInfo, err := h.repository.Fetch(user)
	if err != nil {
		return nil, err
	}

	return userInfo, nil
}
