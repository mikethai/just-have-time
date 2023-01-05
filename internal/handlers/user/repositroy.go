package userHandler

import (
	"github.com/mikethai/just-have-time/internal/model"

	"gorm.io/gorm"
)

type Repository interface {
	Create(param *CreateUserParameter) (*model.User, error)
	Fetch(param *FetchUserParameter) (*model.User, error)
	Follow(param *FollowParameter) (*model.Follow, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{
		db: db,
	}
}

type CreateUserParameter struct {
	Msno    string
	MsnoInt int64
}

type FetchUserParameter struct {
	Msno string
}

type FollowParameter struct {
	followModel model.Follow
}

type User struct {
	Msno string
}

func (r *repository) Create(param *CreateUserParameter) (*model.User, error) {

	newUser := model.User{
		Msno:    param.Msno,
		MsnoInt: param.MsnoInt,
	}

	if err := r.db.Where("msno = ?", newUser.Msno).First(&newUser).Error; err != nil {
		r.db.Create(&newUser)
	}

	return &newUser, nil
}

func (r *repository) Fetch(param *FetchUserParameter) (*model.User, error) {

	newUser := model.User{
		Msno: param.Msno,
	}

	if err := r.db.Where("msno = ?", newUser.Msno).First(&newUser).Error; err != nil {
		return nil, err
	}

	return &newUser, nil
}

func (r *repository) Follow(param *FollowParameter) (*model.Follow, error) {

	newFollowModel := model.Follow{
		FollowerID: param.followModel.FollowerID,
		FolloweeID: param.followModel.FolloweeID,
	}

	r.db.Create(&newFollowModel)

	return &newFollowModel, nil
}
