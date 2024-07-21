package api

import (
	"gin-mall/consts"
	"gin-mall/pkg/utils/ctl"
	"gin-mall/pkg/utils/log"
	"gin-mall/service"
	"gin-mall/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateOrderHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.OrderCreateReq
		if err := c.ShouldBind(&req); err != nil {
			log.LogrusObj.Infoln(err)
			c.JSON(http.StatusBadRequest, err)
			return
		}
		l := service.GetOrderSrv()
		resp, err := l.OrderCreate(c.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Error(err)
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, ctl.RespSuccess(c, resp))
	}
}

func ListOrdersHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.OrderListReq
		if err := c.ShouldBind(&req); err != nil {
			log.LogrusObj.Infoln(err)
			c.JSON(http.StatusBadRequest, err)
			return
		}
		if req.PageSize == 0 {
			req.PageSize = consts.BasePageSize
		}
		l := service.GetOrderSrv()
		resp, err := l.OrderList(c.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Error(err)
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, ctl.RespSuccess(c, resp))
	}
}

func ShowOrderHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.OrderShowReq
		if err := c.ShouldBind(&req); err != nil {
			log.LogrusObj.Infoln(err)
			c.JSON(http.StatusBadRequest, err)
			return
		}
		l := service.GetOrderSrv()
		resp, err := l.OrderShow(c.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Error(err)
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, ctl.RespSuccess(c, resp))
	}
}

func DeleteOrderHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.OrderDeleteReq
		if err := c.ShouldBind(&req); err != nil {
			log.LogrusObj.Infoln(err)
			c.JSON(http.StatusBadRequest, err)
			return
		}
		l := service.GetOrderSrv()
		resp, err := l.OrderDelete(c.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Error(err)
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, ctl.RespSuccess(c, resp))
	}
}
