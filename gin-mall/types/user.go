package types

type UserServiceReq struct {
	NickName string `form:"nick_name" json:"nick_name"`
	UserName string `form:"user_name" json:"user_name"`
	Password string `form:"password" json:"password"`
	Key      string `form:"key" json:"key"` // 前端进行判断
}

type UserInfoResp struct {
	ID       uint   `json:"id"`
	NickName string `json:"nick_name"`
	UserName string `json:"user_name"`
	Type     int    `json:"type"`
	Email    string `json:"email"`
	Status   string `json:"status"`
	CreateAt int64  `json:"create_at"`
}

type UserTokenData struct {
	User         interface{} `json:"user"`
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
}

type UserInfoUpdateReq struct {
	NickName string `form:"nick_name" json:"nick_name"`
}

type UserInfoShowReq struct {
}

type SendEmailServiceReq struct {
	Email    string `form:"email" json:"email"`
	Password string `form:"password" json:"password"`
	// OpertionType 1:绑定邮箱 2：解绑邮箱 3：改密码
	OperationType uint `form:"operation_type" json:"operation_type"`
}

type ValidEmailServiceReq struct {
	Token string `json:"token" form:"token"`
}

type UserFollowingReq struct {
	Id uint `json:"id" form:"id"`
}
type UserUnFollowingReq struct {
	Id uint `json:"id" form:"id"`
}
