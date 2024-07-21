package dao

import (
	"context"
	"gin-mall/repository/db/model"
	"gin-mall/types"
	"gorm.io/gorm"
)

type FavoritesDao struct {
	*gorm.DB
}

func NewFavoritesDao(ctx context.Context) *FavoritesDao {
	return &FavoritesDao{NewDBClient(ctx)}
}

func NewFavoritesDaoByDB(db *gorm.DB) *FavoritesDao {
	return &FavoritesDao{db}
}

func (dao *FavoritesDao) CreateFavorite(favorite *model.Favorite) (err error) {
	err = dao.Create(favorite).Error
	return
}

func (dao *FavoritesDao) FavoriteExistOrNot(pId, uId uint) (exist bool, err error) {
	var count int64
	err = dao.DB.Model(&model.Favorite{}).
		Where("product_id=? AND user_id=?", pId, uId).Count(&count).Error
	if count == 0 || err != nil {
		return false, err
	}
	return true, err
}

func (dao *FavoritesDao) ListFavoriteByUserId(uId uint, pageSize, pageNum int) (r []*types.FavoriteListResp, total int64, err error) {
	err = dao.DB.Model(&model.Favorite{}).Preload("User").
		Where("user_id = ?", uId).Count(&total).Error
	if err != nil {
		return
	}
	err = dao.DB.Model(&model.Favorite{}).
		Joins("AS f LEFT JOIN user AS u on u.id = f.boss_id").
		Joins("LEFT JOIN product AS p ON p.id = f.product_id").
		Joins("LEFT JOIN category AS c ON c.id = p.category_id").
		Where("f.user_id = ? ", uId).
		Offset((pageNum - 1) * pageSize).Limit(pageSize).
		Select("f.user_id AS user_id," +
			"f.product_id AS product_id," +
			"UNIX_TIMESTAMP(f.created_at) AS created_at," +
			"p.title AS title," +
			"p.info AS info," +
			"p.name AS name," +
			"c.id AS category_id," +
			"c.category_name AS category_name," +
			"u.id AS boss_id," +
			"u.user_name AS boss_name," +
			"p.price AS price," +
			"p.img_path AS img_path," +
			"p.discount_price AS discount_price," +
			"p.num AS num," +
			"p.on_sale AS on_sale").
		Find(&r).Error

	return
}

func (dao *FavoritesDao) DeleteFavoriteById(fId uint) error {
	return dao.DB.Where("id=?", fId).Delete(&model.Favorite{}).Error
}
