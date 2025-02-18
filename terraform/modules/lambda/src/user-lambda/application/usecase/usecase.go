package usecase

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/test-lambda/application/domain/model"
	"github.com/test-lambda/application/domain/repository"
	"github.com/test-lambda/application/domain/service"
)

type Usecase struct {
	UserRepository          repository.UserRepository
	UserAsyncProcessService service.UserProcessService
}

func NewUsecase(
	userRepo repository.UserRepository,
	service service.UserProcessService,
) *Usecase {
	return &Usecase{
		UserRepository:          userRepo,
		UserAsyncProcessService: service,
	}
}

type CreateUserInput struct {
	Name      string
	BirthDate string
}

type CreateUserOutput struct {
	ID string
}

func (u *Usecase) CreateUser(input CreateUserInput) (*CreateUserOutput, error) {
	birthDate, err := time.Parse("2006-01-02", input.BirthDate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse birth date, %w", err)
	}

	id := uuid.New().String()

	user := model.User{
		ID:        id,
		Name:      input.Name,
		BirthDate: birthDate,
	}

	err = u.UserRepository.Save(user)

	if err != nil {
		return nil, fmt.Errorf("failed to save user, %w", err)
	}

	err = u.UserAsyncProcessService.ProcessUser(user)
	if err != nil {
		return nil, fmt.Errorf("failed to process user, %w", err)
	}

	return &CreateUserOutput{
		ID: id,
	}, nil
}

func (u *Usecase) GetUser(id int64) (*model.User, error) {
	user, err := u.UserRepository.Get(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user, %w", err)
	}

	return user, nil
}
