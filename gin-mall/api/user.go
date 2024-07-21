package api

import (
	"fmt"
	"gin-mall/consts"
	"gin-mall/pkg/utils/ctl"
	"gin-mall/pkg/utils/log"
	"gin-mall/service"
	"gin-mall/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UserRegisterHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.UserServiceReq
		//参数绑定
		if err := c.ShouldBind(&req); err != nil {
			log.LogrusObj.Infoln(err)
			c.JSON(http.StatusBadRequest, err)
			return
		}
		if req.NickName == "" {
			req.NickName = consts.DefaultNickName // 可以改成你希望的默认昵称
		}
		l := service.GetUserSrv()
		resp, err := l.UserRegister(c.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Infoln(err)
			c.JSON(http.StatusBadRequest, err)
			return
		}
		c.JSON(http.StatusOK, ctl.RespSuccess(c, resp))
	}
}
func UserLoginHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.UserServiceReq
		//参数绑定
		if err := c.ShouldBind(&req); err != nil {
			log.LogrusObj.Infoln(err)
			c.JSON(http.StatusBadRequest, err)
			return
		}
		l := service.GetUserSrv()
		resp, err := l.UserLogin(c.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Infoln(err)
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, ctl.RespSuccess(c, resp))
	}
}

func UserUpdateHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.UserInfoUpdateReq
		if err := c.ShouldBind(&req); err != nil {
			log.LogrusObj.Infoln(err)
			c.JSON(http.StatusBadRequest, err)
			return
		}
		l := service.GetUserSrv()
		resp, err := l.UserInfoUpdate(c.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Infoln(err)
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, ctl.RespSuccess(c, resp))
	}
}

func ShowUserInfoHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.UserInfoShowReq
		if err := c.ShouldBind(&req); err != nil {
			log.LogrusObj.Infoln(err)
			c.JSON(http.StatusBadRequest, err)
			return
		}
		l := service.GetUserSrv()
		resp, err := l.UserInfoShow(c.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Infoln(err)
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, ctl.RespSuccess(c, resp))
	}
}

func SendEmailHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.SendEmailServiceReq
		if err := c.ShouldBind(&req); err != nil {
			log.LogrusObj.Infoln(err)
			c.JSON(http.StatusBadRequest, err)
			return
		}
		l := service.GetUserSrv()
		resp, err := l.SendEmail(c.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Infoln(err)
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, ctl.RespSuccess(c, resp))
	}
}

func ValidEmailHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println(1)
		var req types.ValidEmailServiceReq
		if err := c.ShouldBind(&req); err != nil {
			log.LogrusObj.Infoln(err)
			c.JSON(http.StatusBadRequest, err)
			return
		}
		l := service.GetUserSrv()
		resp, err := l.Valid(c.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Infoln(err)
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, ctl.RespSuccess(c, resp))
	}
}

func UserFollowingHnadler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.UserFollowingReq
		if err := c.ShouldBind(&req); err != nil {
			log.LogrusObj.Infoln(err)
			c.JSON(http.StatusBadRequest, err)
			return
		}
		l := service.GetUserSrv()
		resp, err := l.UserFollow(c.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Error(err)
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, ctl.RespSuccess(c, resp))
	}
}
func UserUnFollowingHnadler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.UserUnFollowingReq
		if err := c.ShouldBind(&req); err != nil {
			log.LogrusObj.Infoln(err)
			c.JSON(http.StatusBadRequest, err)
			return
		}
		l := service.GetUserSrv()
		resp, err := l.UserUnFollow(c.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Error(err)
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, ctl.RespSuccess(c, resp))
	}
}
