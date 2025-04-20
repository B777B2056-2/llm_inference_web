package dao

import (
	"gorm.io/gorm"
	"llm_online_inference/usercenter/resource"
)

type User struct {
	gorm.Model
	Name     string `gorm:"column:name"`
	Password string `gorm:"column:password"`
}

func (User) TableName() string {
	return "user"
}

type UserDao struct {
	db *gorm.DB
}

func NewUserDao() *UserDao {
	return &UserDao{db: resource.DB}
}

func NewUserDaoWithTX(tx *gorm.DB) *UserDao {
	return &UserDao{db: tx}
}

func (u *UserDao) Create(name, pwd string) error {
	return u.db.Create(&User{Name: name, Password: pwd}).Error
}

func (u *UserDao) GetByID(id uint) (*User, error) {
	var user User
	err := u.db.Where("id = ?", id).First(&user).Error
	return &user, err
}

func (u *UserDao) GetByName(name string) (*User, error) {
	var user User
	err := u.db.Where("name = ?", name).First(&user).Error
	return &user, err
}
