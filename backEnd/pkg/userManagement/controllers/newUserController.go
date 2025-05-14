package controller

import (
	"database/sql"
	userModel "social-network/pkg/userManagement/models"
)

type UserController struct {
	um *userModel.UserModel
}

func NewUserController(dbMain *sql.DB) *UserController {
	return &UserController{
		um: userModel.NewUserModel(dbMain),
	}
}
