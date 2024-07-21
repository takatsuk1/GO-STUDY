package dao

import (
	"context"
	"fmt"
	"gin-mall/pkg/utils/log"
	"gin-mall/repository/db/model"
	"gorm.io/gorm"
)

type UserDao struct {
	*gorm.DB
}

func NewUserDao(ctx context.Context) *UserDao {
	return &UserDao{
		NewDBClient(ctx),
	}
}
func NewUserDaoByDB(db *gorm.DB) *UserDao {
	return &UserDao{db}
}

// CreateUser 创建用户
func (dao *UserDao) CreateUser(user *model.User) error {
	return dao.DB.Model(&model.User{}).Create(&user).Error
}
func (dao *UserDao) ExistOrNotByUserName(userName string) (user *model.User, exist bool, err error) {
	err = dao.DB.Model(&model.User{}).Where("user_name=?", userName).First(&user).Error
	if err != nil {
		return user, false, nil
	}
	return user, true, err
}
func (dao *UserDao) GetUserById(id uint) (user *model.User, err error) {
	err = dao.DB.Model(&model.User{}).Where("id = ?", id).First(&user).Error
	return
}

func (dao *UserDao) UpdateUserById(id uint, user *model.User) (err error) {
	err = dao.DB.Model(&model.User{}).Where("id = ?", id).Updates(&user).Error
	fmt.Println(err)
	return
}

func (dao *UserDao) FollowUser(uId, followerId uint) (err error) {
	u, f := new(model.User), new(model.User)
	dao.DB.Model(&model.User{}).Where(`id = ?`, uId).First(&u)
	dao.DB.Model(&model.User{}).Where(`id = ?`, followerId).First(&f)
	err = dao.DB.Model(&f).Where(`id = ?`, followerId).Association(`Relations`).
		Append([]model.User{*u})
	if err != nil {
		log.LogrusObj.Error(err)
		return err
	}

	return
}
func (dao *UserDao) UnFollowUser(uId, followerId uint) (err error) {
	u, f := new(model.User), new(model.User)
	dao.DB.Model(&model.User{}).Where(`id = ?`, uId).First(&u)
	dao.DB.Model(&model.User{}).Where(`id = ?`, followerId).First(&f)
	err = dao.DB.Model(&u).Association(`Relations`).
		Delete(f)
	if err != nil {
		log.LogrusObj.Error(err)
		return err
	}
	return
}
