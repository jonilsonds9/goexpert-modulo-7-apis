package database

import "github.com/jonilsonds9/goexpert-modulo-7-apis/internal/entity"

type UserInterface interface {
	Create(user *entity.User) error
	FindByEmail(email string) (*entity.User, error)
}
