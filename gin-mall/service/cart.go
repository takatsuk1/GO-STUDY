package service

import (
	"context"
	"errors"
	"gin-mall/config"
	"gin-mall/consts"
	"gin-mall/pkg/e"
	"gin-mall/pkg/utils/ctl"
	"gin-mall/pkg/utils/log"
	"gin-mall/repository/db/dao"
	"gin-mall/types"
	"sync"
)

var CartSrvIns *CartSrv
var CartSrvOnce sync.Once

type CartSrv struct{}

func GetCartSrv() *CartSrv {
	CartSrvOnce.Do(func() {
		CartSrvIns = &CartSrv{}
	})
	return CartSrvIns
}

func (s *CartSrv) CartCreate(c context.Context, req *types.CartCreateReq) (resp interface{}, err error) {
	u, err := ctl.GetUserInfo(c)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	//判断有没有这个商品
	_, err = dao.NewProductDao(c).GetProductById(req.ProductId)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	//创建购物车
	cartDao := dao.NewCartDao(c)
	_, status, err := cartDao.CreateCart(req.ProductId, u.Id, req.BossID)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	if status == e.ErrorProductMoreCart {
		err = errors.New(e.GetMsg(status))
		return
	}
	return
}

func (s *CartSrv) CartList(c context.Context, req *types.CartListReq) (resp interface{}, err error) {
	u, err := ctl.GetUserInfo(c)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	carts, err := dao.NewCartDao(c).ListCartByUserId(u.Id, req.PageNum, req.PageSize)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	for i := range carts {
		if config.Config.System.UploadModel == consts.UploadModelLocal {
			carts[i].ImgPath = config.Config.PhotoPath.PhotoHost + config.Config.System.HttpPort + config.Config.PhotoPath.ProductPath + carts[i].ImgPath
		}
	}
	resp = &types.DataListResp{
		Item:  carts,
		Total: int64(len(carts)),
	}
	return
}

func (s *CartSrv) CartUpdate(c context.Context, req *types.UpdateCartServiceReq) (resp interface{}, err error) {
	u, err := ctl.GetUserInfo(c)
	if err != nil {
		log.LogrusObj.Error(err)
		return
	}
	err = dao.NewCartDao(c).UpdateCartNumById(req.Id, u.Id, req.Num)
	if err != nil {
		log.LogrusObj.Error(err)
		return
	}
	return
}

func (s *CartSrv) CartDelete(c context.Context, req *types.CartDeleteReq) (resp interface{}, err error) {
	u, err := ctl.GetUserInfo(c)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	err = dao.NewCartDao(c).DeleteCartById(req.Id, u.Id)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	return
}
