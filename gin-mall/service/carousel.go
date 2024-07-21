package service

import (
	"context"
	"gin-mall/pkg/utils/log"
	"gin-mall/repository/db/dao"
	"gin-mall/types"
	"sync"
)

var CarouselSrvIns *CarouselSrv
var CarouselSrvOnce sync.Once

type CarouselSrv struct {
}

func GetCarouselSrv() *CarouselSrv {
	CarouselSrvOnce.Do(func() {
		CarouselSrvIns = &CarouselSrv{}
	})
	return CarouselSrvIns
}

func (s *CarouselSrv) ListCarousel(c context.Context, req *types.ListCarouselReq) (resp interface{}, err error) {
	carousels, err := dao.NewCarouselDao(c).ListCarousel()
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	resp = &types.DataListResp{
		Item:  carousels,
		Total: int64(len(carousels)),
	}
	return

}
