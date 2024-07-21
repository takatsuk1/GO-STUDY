package dao

import "gin-mall/repository/db/model"

func migrate() (err error) {
	err = _db.Set("gorm:table_options", "charset=utf8mb4").
		AutoMigrate(&model.User{}, &model.Product{}, &model.ProductImg{}, &model.Category{},
			&model.Carousel{}, &model.Favorite{}, &model.Order{}, &model.Cart{}, &model.Address{},
			&model.SkillProduct{}, &model.SkillProduct2MQ{})

	return
}
