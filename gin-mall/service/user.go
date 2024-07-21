package service

import (
	"context"
	"errors"
	"fmt"
	"gin-mall/config"
	"gin-mall/consts"
	"gin-mall/pkg/utils/ctl"
	"gin-mall/pkg/utils/email"
	"gin-mall/pkg/utils/jwt"
	"gin-mall/pkg/utils/log"
	"gin-mall/repository/db/dao"
	"gin-mall/repository/db/model"
	"gin-mall/types"
	"sync"
)

// 单例模式，懒汉模式，即每次在使用的时候才创建一个实例，为了避免懒汉模式的多线程安全问题，这里
// 使用sync.once保证每次只创建一个实例。
// 相对应的饿汉模式即在每次类加载时就创建实例，由于golang中没有类的概念，则是对实例声明时就创建
// 这里就是在var UserSrvIns的时候就创建就是饿汉模式。懒汉模式就是在api接口中实际使用的时候调用
// GetUserSrv()创建一个实例
var UserSrvIns *UserSrv
var UserSrvOnce sync.Once

// 空结构体是一个实例
type UserSrv struct {
}

func GetUserSrv() *UserSrv {
	//使用srvonce.Do()传递一个空函数保证只创建一个，即在多线程的情况下，如果再次调用这个函数则直接返回UserSrvIns实例
	UserSrvOnce.Do(func() {
		UserSrvIns = &UserSrv{}
	})
	return UserSrvIns
}

func (s *UserSrv) UserRegister(ctx context.Context, req *types.UserServiceReq) (resp interface{}, err error) {
	userDao := dao.NewUserDao(ctx)
	_, exist, err := userDao.ExistOrNotByUserName(req.UserName)
	if err != nil {
		log.LogrusObj.Error(err)
		return
	}
	if exist {
		err = errors.New("用户已存在，请更换用户名")
		return
	}
	user := &model.User{
		NickName: req.NickName,
		UserName: req.UserName,
		Status:   model.Active,
		Money:    "100",
	}
	err = user.SetPassword(req.Password)
	if err != nil {
		log.LogrusObj.Error(err)
		return
	}
	money, err := user.EncryptMoney(req.Key)
	if err != nil {
		log.LogrusObj.Error(err)
		return
	}
	user.Money = money
	err = userDao.CreateUser(user)
	if err != nil {
		return
	}
	return
}

func (s *UserSrv) UserLogin(ctx context.Context, req *types.UserServiceReq) (resp interface{}, err error) {
	var user *model.User
	userDao := dao.NewUserDao(ctx)
	user, exist, err := userDao.ExistOrNotByUserName(req.UserName)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	if !exist {
		err = errors.New("用户已存在，请更换用户名")
		log.LogrusObj.Error(err)
		return
	}
	//检验密码
	if !user.CheckPassword(req.Password) {
		err = errors.New("密码错误")
		log.LogrusObj.Error(err)
		return
	}
	accessToken, refreshToken, err := jwt.GenerateToken(user.ID, req.UserName)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	userResp := &types.UserInfoResp{
		ID:       user.ID,
		NickName: user.NickName,
		UserName: user.UserName,
		Email:    user.Email,
		Status:   user.Status,
		CreateAt: user.CreatedAt.Unix(),
	}
	resp = &types.UserTokenData{
		User:         userResp,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return
}

// 更改用户nick_name,后续增加更改其他信息
func (s *UserSrv) UserInfoUpdate(c context.Context, req *types.UserInfoUpdateReq) (resp interface{}, err error) {
	u, err := ctl.GetUserInfo(c)
	if err != nil {
		log.LogrusObj.Error(err)
		return
	}
	userDao := dao.NewUserDao(c)
	user, err := userDao.GetUserById(u.Id)
	if err != nil {
		fmt.Println(1)
		log.LogrusObj.Error(err)
		return nil, err
	}
	//不够标准化
	if req.NickName != "" {
		user.NickName = req.NickName
	}
	err = userDao.UpdateUserById(u.Id, user)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	return
}

func (s *UserSrv) UserInfoShow(c context.Context, req *types.UserInfoShowReq) (resp interface{}, err error) {
	userDao := dao.NewUserDao(c)
	u, err := ctl.GetUserInfo(c)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	user, err := userDao.GetUserById(u.Id)
	if err != nil {
		log.LogrusObj.Error(err)
		return
	}
	resp = &types.UserInfoResp{
		ID:       user.ID,
		NickName: user.NickName,
		UserName: user.UserName,
		Email:    user.Email,
		Status:   user.Status,
		CreateAt: user.CreatedAt.Unix(),
	}
	return
}

func (s *UserSrv) SendEmail(c context.Context, req *types.SendEmailServiceReq) (resp interface{}, err error) {
	u, err := ctl.GetUserInfo(c)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	var address string
	token, err := jwt.GenerateEmailToken(u.Id, req.OperationType, req.Email, req.Password)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	sender := email.NewEmailSender()
	address = config.Config.Email.ValidEmail + token
	//连接字符串
	mailText := fmt.Sprintf(consts.EmailOperationMap[req.OperationType], address)
	if err = sender.Send(mailText, req.Email, "商城验证码"); err != nil {
		log.LogrusObj.Error(err)
		return
	}
	return
}

func (s *UserSrv) Valid(c context.Context, req *types.ValidEmailServiceReq) (resp interface{}, err error) {
	var userId uint
	var email string
	var password string
	var operationType uint
	//验证token
	if req.Token == "" {
		err = errors.New("token不存在")
		log.LogrusObj.Error(err)
		return
	}
	claims, err := jwt.ParseEmailToken(req.Token)
	if err != nil {
		log.LogrusObj.Error(err)
		return
	} else {
		userId = claims.UserId
		email = claims.Email
		password = claims.Password
		operationType = claims.OperationType
	}
	//获取用户信息
	userDao := dao.NewUserDao(c)
	user, err := userDao.GetUserById(userId)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	switch operationType {
	case consts.EmailOperationBinding:
		user.Email = email
	case consts.EmailOperationNoBinding:
		user.Email = ""
	case consts.EmailOperationUpdatePassword:
		err = user.SetPassword(password)
		if err != nil {
			log.LogrusObj.Error(err)
			return nil, err
		}
	default:
		return nil, errors.New("不支持该操作")
	}
	err = userDao.UpdateUserById(userId, user)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	resp = &types.UserInfoResp{
		ID:       user.ID,
		NickName: user.NickName,
		UserName: user.UserName,
		Email:    user.Email,
		Status:   user.Status,
		CreateAt: user.CreatedAt.Unix(),
	}
	return
}

func (s *UserSrv) UserFollow(c context.Context, req *types.UserFollowingReq) (resp interface{}, err error) {
	u, err := ctl.GetUserInfo(c)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	err = dao.NewUserDao(c).FollowUser(u.Id, req.Id)
	if err != nil {
		log.LogrusObj.Error(err)
		return
	}
	return
}
func (s *UserSrv) UserUnFollow(c context.Context, req *types.UserUnFollowingReq) (resp interface{}, err error) {
	u, err := ctl.GetUserInfo(c)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	err = dao.NewUserDao(c).UnFollowUser(u.Id, req.Id)
	if err != nil {
		log.LogrusObj.Error(err)
		return
	}
	return
}
