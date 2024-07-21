package middleware

import (
	"gin-mall/consts"
	"gin-mall/pkg/e"
	"gin-mall/pkg/utils/ctl"
	"gin-mall/pkg/utils/jwt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//使用context上下文
		var code int
		code = e.SUCCESS
		//查询当前请求的header-Token
		accesstoken := c.GetHeader("access_token")
		refreshtoken := c.GetHeader("refresh_token")
		if accesstoken == "" { //token为空的情况
			code = http.StatusNotFound
			c.JSON(e.InvalidParams, gin.H{
				"status": code,
				"msg":    e.GetMsg(code),
				"data":   "缺少Token",
			})
			c.Abort() //跳出请求
			return
		}
		//解析token验证旧的access和refresh是否正确
		newAccesstoken, newRefreshtoken, err := jwt.ParseRefreshToken(accesstoken, refreshtoken)
		if err != nil { //token解析错误
			code = e.ErrorAuthCheckTokenFail
			//token解析超时,expireat为系统内置定义时间,这里已经在解析的函数内部判断了是否超时，超时则重新获取token
		} //else if time.Now().Unix() > claims.ExpiresAt {
		//code = e.ErrorAuthCheckTokenTimeout
		//}
		//解析失败
		if code != e.SUCCESS {
			c.JSON(e.InvalidParams, gin.H{
				"status": code,
				"msg":    e.GetMsg(code),
				"data":   "鉴权失败",
			})
			c.Abort()
			return
		}
		claims, err := jwt.ParseToken(newAccesstoken)
		if err != nil {
			code = e.ErrorAuthCheckTokenFail
			c.JSON(e.InvalidParams, gin.H{
				"status": code,
				"msg":    e.GetMsg(code),
				"data":   err.Error(),
			})
			c.Abort()
			return
		}
		SetToken(c, newAccesstoken, newRefreshtoken)
		//c.Request.Context获取当前请求的上下文
		//&ctl.UserInfo{Id: claims.ID}创建新的userInfo并将当前token声明的id传进去
		c.Request = c.Request.WithContext(ctl.NewContext(c.Request.Context(), &ctl.UserInfo{Id: claims.ID}))
		//ctl.InitUserInfo(c.Request.Context())
		c.Next()
		////Token解析成功,,将用户id绑定至context
		//c.Request = c.Request.WithContext(ctl.NewContext(c.Request.Context(), &ctl.UserInfo{Id: claims.Id}))
		//c.Next()
	}
}

// 设置http头和cookies，通过http响应将两个token传回到客户端，然后设置客户端的cookie也保存两个token
func SetToken(c *gin.Context, accessToken, refreshToken string) {
	secure := IsHttps(c)
	//保存在响应头中
	c.Header(consts.AccessTokenHeader, accessToken)
	c.Header(consts.RefreshTokenHeader, refreshToken)
	//保存在cookie中
	c.SetCookie(consts.AccessTokenHeader, accessToken, consts.MaxAge, "/", "", secure, true)
	c.SetCookie(consts.RefreshTokenHeader, refreshToken, consts.MaxAge, "/", "", secure, true)
}

// 判断是否https
func IsHttps(c *gin.Context) bool {
	if c.GetHeader(consts.HeaderForwardedProto) == "https" || c.Request.TLS != nil {
		return true
	}
	return false
}
