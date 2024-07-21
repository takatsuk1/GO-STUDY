package api

import (
	"gin-mall/pkg/utils/ctl"
	"gin-mall/pkg/utils/log"
	"gin-mall/service"
	"gin-mall/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ListCategoryHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.ListCategoryReq
		if err := c.ShouldBind(&req); err != nil {
			log.LogrusObj.Infoln(err)
			c.JSON(http.StatusBadRequest, err)
			return
		}
		l := service.GetCategorySrv()
		resp, err := l.CategoryList(c.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Error(err)
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, ctl.RespSuccess(c, resp))
	}
}
