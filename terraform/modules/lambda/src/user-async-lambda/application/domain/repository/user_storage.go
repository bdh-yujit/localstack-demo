package repository

import "github.com/test-async-lambda/application/domain/model"

type UserRepository interface {
	SaveUsers(users []model.User) error
}
