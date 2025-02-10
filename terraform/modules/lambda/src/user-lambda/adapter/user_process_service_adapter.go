package adapter

import (
	"github.com/test-lambda/application/domain/model"
	"github.com/test-lambda/application/domain/service"
	"github.com/test-lambda/infrastructure/sqs"
)

type UserProcessServiceAdapter struct {
	UserAsyncTask sqs.UserAsyncTask
}

func NewUserProcessServiceAdapter() (service.UserProcessService, error) {
	uat, err := sqs.NewUserAsyncTask()
	if err != nil {
		return nil, err
	}
	return &UserProcessServiceAdapter{
		UserAsyncTask: uat,
	}, nil
}

func (u *UserProcessServiceAdapter) ProcessUser(user model.User) error {
	return u.UserAsyncTask.SaveToS3(user)
}
