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

// CreateAddressHandler 新增收货地址
func CreateAddressHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.AddressCreateReq
		if err := ctx.ShouldBind(&req); err != nil {
			// 参数校验
			log.LogrusObj.Infoln(err)
			ctx.JSON(http.StatusBadRequest, err)
			return
		}

		l := service.GetAddressSrv()
		resp, err := l.AddressCreate(ctx.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Infoln(err)
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}
		ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, resp))
	}
}

// ShowAddressHandler 展示某个收货地址
func ShowAddressHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.AddressGetReq
		if err := ctx.ShouldBind(&req); err != nil {
			// 参数校验
			log.LogrusObj.Infoln(err)
			ctx.JSON(http.StatusBadRequest, err)
			return
		}

		l := service.GetAddressSrv()
		resp, err := l.AddressShow(ctx.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Infoln(err)
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}
		ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, resp))
	}
}

// ListAddressHandler 展示收货地址
func ListAddressHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.AddressListReq
		if err := ctx.ShouldBind(&req); err != nil {
			// 参数校验
			log.LogrusObj.Infoln(err)
			ctx.JSON(http.StatusBadRequest, err)
			return
		}
		if req.PageSize == 0 {
			req.PageSize = consts.BasePageSize
		}

		l := service.GetAddressSrv()
		resp, err := l.AddressList(ctx.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Infoln(err)
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}
		ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, resp))
	}
}

// UpdateAddressHandler 修改收货地址
func UpdateAddressHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.AddressServiceReq
		if err := ctx.ShouldBind(&req); err != nil {
			// 参数校验
			log.LogrusObj.Infoln(err)
			ctx.JSON(http.StatusBadRequest, err)
			return
		}

		l := service.GetAddressSrv()
		resp, err := l.AddressUpdate(ctx.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Infoln(err)
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}
		ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, resp))
	}
}

// DeleteAddressHandler 删除收获地址
func DeleteAddressHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.AddressDeleteReq
		if err := ctx.ShouldBind(&req); err != nil {
			// 参数校验
			log.LogrusObj.Infoln(err)
			ctx.JSON(http.StatusBadRequest, err)
			return
		}

		l := service.GetAddressSrv()
		resp, err := l.AddressDelete(ctx.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Infoln(err)
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}
		ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, resp))
	}
}
