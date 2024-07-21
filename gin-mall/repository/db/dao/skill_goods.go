package dao

import (
	"context"
	"gin-mall/consts"
	"gin-mall/repository/db/model"
	"gorm.io/gorm"
)

type SkillGoodsDao struct {
	*gorm.DB
}

func NewSkillGoodsDao(ctx context.Context) *SkillGoodsDao {
	return &SkillGoodsDao{NewDBClient(ctx)}
}

func (dao *SkillGoodsDao) Create(in *model.SkillProduct) error {
	return dao.Model(&model.SkillProduct{}).Create(&in).Error
}

// 批量创建 SkillProduct 记录。
func (dao *SkillGoodsDao) BatchCreate(in []*model.SkillProduct) error {
	return dao.Model(&model.SkillProduct{}).
		CreateInBatches(&in, consts.ProductBatchCreate).Error
}

// 批量创建 SkillProduct 记录
// Create 方法批量插入 in 中的所有记录
func (dao *SkillGoodsDao) CreateByList(in []*model.SkillProduct) error {
	return dao.Model(&model.SkillProduct{}).Create(&in).Error
}

func (dao *SkillGoodsDao) ListSkillGoods() (resp []*model.SkillProduct, err error) {
	err = dao.Model(&model.SkillProduct{}).
		Where("num > 0").Find(&resp).Error

	return
}
