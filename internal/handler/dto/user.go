package dto

import (
	"github.com/zuzaaa-dev/stawberry/internal/domain/service/user"
)

type RegistrationUserReq struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
	IsStore  bool   `json:"is_store" binding:"required"`
}

type RegistrationUserResp struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (ru *RegistrationUserReq) ConvertToSvc() user.User {
	return user.User{
		Name:     ru.Name,
		Password: ru.Password,
		Email:    ru.Email,
		IsStore:  ru.IsStore,
	}
}

type LoginUserReq struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginUserResp struct {
	Token string `json:"token"`
}
