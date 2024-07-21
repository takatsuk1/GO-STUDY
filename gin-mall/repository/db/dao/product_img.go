package dao

import (
	"context"
	"gin-mall/repository/db/model"
	"gin-mall/types"
	"gorm.io/gorm"
)

type ProductImgDao struct {
	*gorm.DB
}

func NewProductImgDao(ctx context.Context) *ProductImgDao {
	return &ProductImgDao{NewDBClient(ctx)}
}

func NewProductImgDaoByDB(db *gorm.DB) *ProductImgDao {
	return &ProductImgDao{db}
}

func (dao *ProductImgDao) CreateProductImg(productImg *model.ProductImg) (err error) {
	err = dao.DB.Model(&model.ProductImg{}).Create(&productImg).Error
	return
}

// 根据商品ID获取商品图片
func (dao *ProductImgDao) ListProductImgByProductId(pId uint) (r []*types.ProductImgResp, err error) {
	err = dao.DB.Model(&model.ProductImg{}).Where("product_id = ?", pId).Find(&r).Error
	return
}
