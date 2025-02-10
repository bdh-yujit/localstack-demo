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

type ExecuteInput struct {
	Name      string
	BirthDate string
}

type ExecuteOutput struct {
	ID string
}

func (u *Usecase) Execute(input ExecuteInput) (*ExecuteOutput, error) {
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

	return &ExecuteOutput{
		ID: id,
	}, nil
}
