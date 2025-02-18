package repository

import "github.com/test-lambda/application/domain/model"

type UserRepository interface {
	Save(user model.User) error
	Get(id int64) (*model.User, error)
}
