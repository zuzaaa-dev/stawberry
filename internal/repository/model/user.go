package model

import (
	"github.com/zuzaaa-dev/stawberry/internal/domain/entity"
	"github.com/zuzaaa-dev/stawberry/internal/domain/service/user"
)

type User struct {
	ID       uint   `gorm:"column:id"`
	Name     string `gorm:"column:name"`
	Phone    string `gorm:"column:phone_number"`
	Password string `gorm:"column:password_hash"`
	Email    string `gorm:"column:email"`
	IsStore  bool   `gorm:"column:is_store"`
}

func ConvertUserFromSvc(u user.User) User {
	return User{
		Name:     u.Name,
		Email:    u.Email,
		Phone:    u.Phone,
		Password: u.Password,
		IsStore:  u.IsStore,
	}
}

func ConvertUserToEntity(u User) entity.User {
	return entity.User{
		ID:       u.ID,
		Name:     u.Name,
		Email:    u.Email,
		Phone:    u.Phone,
		Password: u.Password,
		IsStore:  u.IsStore,
	}
}
