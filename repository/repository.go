package repository

import (
	"gorm.io/gorm"
)

type Users struct {
	gorm.Model
	Id    int32
	Name  string
	Email string
	Phone string
}

type UsersRepository struct {
	db *gorm.DB
}

func NewRepository(connect *gorm.DB) *UsersRepository {
	return &UsersRepository{db: connect}
}

func (u *UsersRepository) CreateUser(user *Users) error {
	return u.db.Create(user).Error
}

func (u *UsersRepository) DeleteUser(userID int) error {
	return u.db.Delete(&Users{}, userID).Error
}

func (u *UsersRepository) GetUsers() ([]Users, error) {
	var users []Users
	err := u.db.Find(&users).Error
	return users, err
}
