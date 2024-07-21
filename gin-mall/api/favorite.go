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

func CreateFavoriteHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.FavoriteCreateReq
		if err := c.ShouldBind(&req); err != nil {
			log.LogrusObj.Infoln(err)
			c.JSON(http.StatusBadRequest, err)
			return
		}
		l := service.GetFavoriteSrv()
		resp, err := l.FavoriteCreate(c.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Infoln(err)
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, ctl.RespSuccess(c, resp))
	}
}

// ListFavoritesHandler 收藏夹详情接口
func ListFavoritesHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.FavoritesServiceReq
		if err := ctx.ShouldBind(&req); err != nil {
			// 参数校验
			log.LogrusObj.Infoln(err)
			ctx.JSON(http.StatusBadRequest, err)
			return
		}
		if req.PageSize == 0 {
			req.PageSize = consts.BasePageSize
		}

		l := service.GetFavoriteSrv()
		resp, err := l.FavoriteList(ctx.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Infoln(err)
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}
		ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, resp))
	}
}

func DeleteFavoriteHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.FavoriteDeleteReq
		if err := ctx.ShouldBind(&req); err != nil {
			log.LogrusObj.Infoln(err)
			ctx.JSON(http.StatusBadRequest, err)
			return
		}
		l := service.GetFavoriteSrv()
		resp, err := l.FavoriteDelete(ctx.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Infoln(err)
			ctx.JSON(http.StatusInternalServerError, err)
		}
		ctx.JSON(http.StatusOK, ctl.RespSuccess(nil, resp))
	}
}
