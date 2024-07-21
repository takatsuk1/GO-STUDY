package ctl

import (
	"context"
	"errors"
)

type key int

var userKey key

type UserInfo struct {
	Id uint `json:"id"`
}

func GetUserInfo(c context.Context) (*UserInfo, error) {
	user, ok := FromContext(c)
	if !ok {
		return nil, errors.New("获取用户信息错误")
	}
	return user, nil
}

func NewContext(c context.Context, u *UserInfo) context.Context {
	return context.WithValue(c, userKey, u)
}

func FromContext(c context.Context) (*UserInfo, bool) {
	u, ok := c.Value(userKey).(*UserInfo)
	return u, ok
}
