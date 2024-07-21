package dao

import (
	"context"
	"gin-mall/repository/db/model"
	"gin-mall/types"

	"gorm.io/gorm"
)

type ProductDao struct {
	*gorm.DB
}

// 从上下文中获取配置信息并创建新的数据库连接的场景
func NewProductDao(c context.Context) *ProductDao {
	return &ProductDao{NewDBClient(c)}
}

// 已经有一个现成的数据库实例，并希望直接使用它的场景。
func NewProductDaoByDB(db *gorm.DB) *ProductDao {
	return &ProductDao{db}
}

func (dao *ProductDao) CreateProduct(product *model.Product) error {
	return dao.DB.Model(&model.Product{}).Create(&product).Error
}

func (dao *ProductDao) UpdateProduct(pId uint, product *model.Product) error {
	return dao.DB.Model(&model.Product{}).Where("id = ?", pId).Updates(product).Error
}

func (dao *ProductDao) DeleteProduct(pId uint, uId uint) error {
	return dao.DB.Model(&model.Product{}).Where("id = ? AND boss_id = ?", pId, uId).Delete(&model.Product{}).Error
}

func (dao *ProductDao) ListProductByCondition(condition map[string]interface{}, page types.BasePage) (products []*model.Product, err error) {
	err = dao.DB.Where(condition).Offset((page.PageNum - 1) * page.PageSize).
		Limit(page.PageSize).
		Find(&products).Error
	return
}

// 根据情况获取商品的数量
func (dao *ProductDao) CountProductByCondition(condition map[string]interface{}) (total int64, err error) {
	err = dao.DB.Model(&model.Product{}).Where(condition).Count(&total).Error
	return
}

func (dao *ProductDao) ShowProductById(id uint) (product *model.Product, err error) {
	err = dao.DB.Model(&model.Product{}).Where("id = ?", id).First(&product).Error
	return
}

func (dao *ProductDao) SearchProduct(name, info string, page types.BasePage) (products []*model.Product, count int64, err error) {
	err = dao.DB.Model(&model.Product{}).
		Where("name LIKE ? OR info LIKE ?", "%"+name+"%", "%"+info+"%").
		Offset((page.PageNum - 1) * page.PageSize).
		Find(&products).Error
	if err != nil {
		return
	}
	err = dao.DB.Model(&model.Product{}).
		Where("name LIKE ? OR info LIKE ?", "%"+name+"%", "%"+info+"%").
		Count(&count).
		Error
	return
}

func (dao *ProductDao) GetProductById(id uint) (product *model.Product, err error) {
	err = dao.DB.Model(&model.Product{}).Where("id = ?", id).First(&product).Error
	return
}
