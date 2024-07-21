package consts

import "time"

const (
	AccessTokenExpireDuration  = 24 * time.Hour
	RefreshTokenExpireDuration = 10 * 24 * time.Hour
)

const (
	DefaultNickName = "默认用户"
)
const (
	AccessTokenHeader    = "access_token"
	RefreshTokenHeader   = "refresh_token"
	HeaderForwardedProto = "X-Forwarded-Proto"
	MaxAge               = 3600 * 24
)

// 预定义每个常量的键 uint类型
const (
	EmailOperationBinding = iota + 1
	EmailOperationNoBinding
	EmailOperationUpdatePassword
)

var EmailOperationMap = map[uint]string{
	EmailOperationBinding:        "您正在绑定邮箱, 请点击链接确定身份 %s",
	EmailOperationNoBinding:      "您正在解邦邮箱, 请点击链接确定身份 %s",
	EmailOperationUpdatePassword: "您正在修改密码, 请点击链接校验身份 %s",
}
