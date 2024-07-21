package service

import (
	"context"
	"gin-mall/pkg/utils/log"
	"gin-mall/repository/db/dao"
	"gin-mall/types"
	"sync"
)

var CategorySrvIns *CategorySrv
var CategorySrvOnce sync.Once

type CategorySrv struct {
}

func GetCategorySrv() *CategorySrv {
	CategorySrvOnce.Do(func() {
		CategorySrvIns = &CategorySrv{}
	})
	return CategorySrvIns
}

func (s *CategorySrv) CategoryList(c context.Context, req *types.ListCategoryReq) (resp interface{}, err error) {
	categories, err := dao.NewCategoryDao(c).ListCategory()
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	cResp := make([]*types.ListCategoryResp, 0)
	for _, v := range categories {
		cResp = append(cResp, &types.ListCategoryResp{
			ID:           v.ID,
			CategoryName: v.CategoryName,
			CreatedAt:    v.CreatedAt.Unix(),
		})
	}
	resp = &types.DataListResp{
		Item:  cResp,
		Total: int64(len(cResp)),
	}
	return
}
