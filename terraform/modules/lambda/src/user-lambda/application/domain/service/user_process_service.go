package service

import "github.com/test-lambda/application/domain/model"

type UserProcessService interface {
	ProcessUser(user model.User) error
}

type userProcessService struct{}
