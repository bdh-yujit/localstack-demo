package adapter

import (
	"fmt"
	"strconv"
	"time"

	"github.com/test-lambda/application/domain/model"
	"github.com/test-lambda/application/domain/repository"
	"github.com/test-lambda/infrastructure/dynamodb"
	"github.com/test-lambda/infrastructure/mysql"
)

type userRepositoryAdapter struct {
	userDynamoDbDao dynamodb.UserDao
	userMysqlDao    mysql.UserDao
}

func NewUserRepository() repository.UserRepository {
	userDynamoDbDao, err := dynamodb.NewDynamoDbUserDao()
	if err != nil {
		panic(err)
	}
	userMysqlDao := mysql.NewMySqlUserDao()
	return &userRepositoryAdapter{
		userDynamoDbDao: userDynamoDbDao,
		userMysqlDao:    userMysqlDao,
	}
}

func (u *userRepositoryAdapter) Save(user model.User) error {
	return u.userDynamoDbDao.Save(user)
}

func (u *userRepositoryAdapter) Get(id int64) (*model.User, error) {
	user, err := u.userMysqlDao.Get(id)
	if err != nil {
		return nil, err
	}
	return &model.User{
		ID:        strconv.Itoa(int(user.UserID)),
		Name:      fmt.Sprintf("%s %s", user.Firstname, user.Lastname),
		BirthDate: time.Now(),
	}, nil
}
