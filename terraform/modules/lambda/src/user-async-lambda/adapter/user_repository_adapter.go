package adapter

import (
	"github.com/test-async-lambda/application/domain/model"
	"github.com/test-async-lambda/application/domain/repository"
	"github.com/test-async-lambda/infrastructure/storage"
)

type UserRepositoryAdapter struct {
	s3Storate storage.S3Storage
}

func NewUserRepositoryAdapter() (repository.UserRepository, error) {
	s, err := storage.NewS3Storage()
	if err != nil {
		return nil, err
	}
	return &UserRepositoryAdapter{
		s3Storate: s,
	}, nil
}

func (u *UserRepositoryAdapter) SaveUsers(users []model.User) error {
	return u.s3Storate.SaveUsers(users)
}
