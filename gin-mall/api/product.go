package api

import (
	"errors"
	"fmt"
	"gin-mall/consts"
	"gin-mall/pkg/utils/ctl"
	"gin-mall/pkg/utils/log"
	"gin-mall/service"
	"gin-mall/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateProductHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.ProductCreateReq
		if err := c.ShouldBind(&req); err != nil {
			log.LogrusObj.Infoln(err)
			c.JSON(http.StatusBadRequest, err)
			return
		}
		form, _ := c.MultipartForm()
		fmt.Println(form)
		files := form.File["image"]
		l := service.GetProductSrv()
		resp, err := l.ProductCreate(c.Request.Context(), files, &req)
		if err != nil {
			log.LogrusObj.Error(err)
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, ctl.RespSuccess(c, resp))
	}
}

func UpdateProductHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.ProductUpdateReq
		if err := c.ShouldBind(&req); err != nil {
			log.LogrusObj.Infoln(err)
			c.JSON(http.StatusBadRequest, err)
			return
		}
		l := service.GetProductSrv()
		resp, err := l.ProductUpdate(c.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Error(err)
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, ctl.RespSuccess(c, resp))
	}
}

func DeleteProductHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.ProductDeleteReq
		if err := c.ShouldBind(&req); err != nil {
			log.LogrusObj.Infoln(err)
			c.JSON(http.StatusBadRequest, err)
			return
		}
		l := service.GetProductSrv()
		resp, err := l.ProductDelete(c.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Error(err)
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, ctl.RespSuccess(c, resp))
	}
}

func ListProductHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.ProductListReq
		if err := c.ShouldBind(&req); err != nil {
			log.LogrusObj.Infoln(err)
			c.JSON(http.StatusBadRequest, err)
			return
		}
		if req.PageSize == 0 {
			req.PageSize = consts.BaseProductPageSize
		}
		l := service.GetProductSrv()
		resp, err := l.ProductList(c.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Error(err)
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, ctl.RespSuccess(c, resp))
	}
}

func ShowProductHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.ProductShowReq
		if err := c.ShouldBind(&req); err != nil {
			log.LogrusObj.Infoln(err)
			c.JSON(http.StatusBadRequest, err)
			return
		}
		l := service.GetProductSrv()
		resp, err := l.ProductShow(c.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Error(err)
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, ctl.RespSuccess(c, resp))
	}
}
func SearchProductsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.ProductSearchReq
		if err := c.ShouldBind(&req); err != nil {
			log.LogrusObj.Infoln(err)
			c.JSON(http.StatusBadRequest, err)
			return
		}
		l := service.GetProductSrv()
		resp, err := l.ProductSearch(c.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Error(err)
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, ctl.RespSuccess(c, resp))
	}
}

func ListProductImgHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.ListProductImgReq
		if err := c.ShouldBind(&req); err != nil {
			log.LogrusObj.Infoln(err)
			c.JSON(http.StatusBadRequest, err)
			return
		}
		if req.ID == 0 {
			err := errors.New("参数错误,id不能为空")
			c.JSON(http.StatusBadRequest, err)
			return
		}
		l := service.GetProductSrv()
		resp, err := l.ProductImgList(c.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Error(err)
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, ctl.RespSuccess(c, resp))
	}
}
