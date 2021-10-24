package repository

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Users struct {
	gorm.Model
	Id    int32
	Name  string
	Email string
	Phone string
}

type Log struct {
	DateTime time.Time
	Message  string
}

type UsersRepository struct {
	db  *gorm.DB
	log *gorm.DB
}

func NewRepository(DBConnect *gorm.DB, clickhouseConnect *gorm.DB) *UsersRepository {
	return &UsersRepository{db: DBConnect, log: clickhouseConnect}
}

func (u *UsersRepository) CreateUser(user *Users) error {
	err := u.db.Create(user).Error
	if err == nil {
		var message = fmt.Sprintf("add user : %v", user.Name)
		var log = Log{DateTime: time.Now(), Message: message}
		u.log.Create(log)
	}
	return err
}

func (u *UsersRepository) DeleteUser(userID int) error {
	return u.db.Delete(&Users{}, userID).Error
}

func (u *UsersRepository) GetUsers() ([]Users, error) {
	var users []Users
	err := u.db.Find(&users).Error
	return users, err
}
