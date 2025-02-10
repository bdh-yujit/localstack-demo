package usecase

import (
	"github.com/test-async-lambda/adapter"
	"github.com/test-async-lambda/application/domain/model"
	"github.com/test-async-lambda/application/domain/repository"
)

type Usecase struct {
	userRepository repository.UserRepository
}

func NewUsecase() (*Usecase, error) {
	repo, err := adapter.NewUserRepositoryAdapter()
	if err != nil {
		return nil, err
	}
	return &Usecase{
		userRepository: repo,
	}, nil
}

type ExecuteInputUser struct {
	Name      string
	BirthDate string
	ID        string
}

type ExecuteInput struct {
	Users []ExecuteInputUser
}

func (u *Usecase) Execute(input ExecuteInput) error {
	users := make([]model.User, 0, len(input.Users))
	for _, user := range input.Users {
		user := model.User{
			ID:        user.ID,
			Name:      user.Name,
			BirthDate: user.BirthDate,
		}
		users = append(users, user)
	}

	return u.userRepository.SaveUsers(users)
}
