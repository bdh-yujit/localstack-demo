package mysql

import (
	"github.com/test-lambda/infrastructure/mysql/gorm/model"
	"gorm.io/gorm"
)

type UserDao interface {
	Get(id int64) (model.User, error)
}

type userDao struct {
	db *gorm.DB
}

func NewMySqlUserDao() UserDao {
	return &userDao{
		db: readerConn,
	}
}

func (u *userDao) Get(id int64) (model.User, error) {
	var user model.User
	result := u.db.Where("user_id = ?", id).First(&user)
	if result.Error != nil {
		return model.User{}, result.Error
	}
	return user, nil
}
