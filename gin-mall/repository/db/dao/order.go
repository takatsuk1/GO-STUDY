package dao

import (
	"context"
	"gin-mall/repository/db/model"
	"gin-mall/types"
	"gorm.io/gorm"
)

type OrderDao struct {
	*gorm.DB
}

func NewOrderDao(ctx context.Context) *OrderDao {
	return &OrderDao{NewDBClient(ctx)}
}

func NewOrderDaoByDB(db *gorm.DB) *OrderDao {
	return &OrderDao{db}
}

// CreateOrder 创建订单
func (dao *OrderDao) CreateOrder(order *model.Order) error {
	return dao.DB.Create(&order).Error
}

func (dao *OrderDao) ListOrderByCondition(uId uint, req *types.OrderListReq) (r []*types.OrderListResp, count int64, err error) {
	//商城算是一个TOC的应用，TOC的应该是不允许join操作的，看看后续怎么改走缓存，比如走缓存，找找免费的CDN之类的
	//TOC:To Consumer，面向消费者）的应用，如电子商务平台，通常需要处理大量的用户请求和数据交互。
	//为了提高性能和用户体验，这类应用一般会避免使用复杂的数据库操作（如 JOIN 操作），而是依赖缓存和 CDN（内容分发网络）来加速数据访问。
	//通过将静态资源分发到全球各地的 CDN 节点，用户可以从最近的节点获取资源，从而减少延迟，提高加载速度。
	d := dao.DB.Model(&model.Order{}).
		Where("user_id = ?", uId)
	if req.Type != 0 {
		d.Where("type = ?", req.Type)
	}
	d.Count(&count) //

	db := dao.DB.Model(&model.Order{}).
		Joins("AS o LEFT JOIN product AS p ON p.id = o.product_id").
		Joins("LEFT JOIN address AS a ON a.id = o.address_id").
		Where("o.user_id = ?", uId)
	if req.Type != 0 {
		db.Where("type = ?", req.Type)
	}
	db.Offset((req.PageNum - 1) * req.PageSize).
		Limit(req.PageSize).Order("created_at desc").
		Select("o.id AS id," +
			"o.order_num AS order_num," +
			"UNIX_TIMESTAMP(o.created_at) AS created_at," +
			"UNIX_TIMESTAMP(o.updated_at) AS updated_at," +
			"o.user_id AS user_id," +
			"o.product_id AS product_id," +
			"o.boss_id AS boss_id," +
			"o.num AS num," +
			"o.type AS type," +
			"p.name AS name," +
			"p.discount_price AS discount_price," +
			"p.img_path AS img_path," +
			"a.name AS address_name," +
			"a.phone AS address_phone," +
			"a.address AS address").
		Find(&r)

	return
}

func (dao *OrderDao) ShowOrderById(id, uId uint) (r *types.OrderListResp, err error) {
	err = dao.DB.Model(&model.Order{}).
		Joins("AS o LEFT JOIN product AS p ON p.id = o.product_id").
		Joins("LEFT JOIN address AS a ON a.id = o.address_id").
		Where("o.id = ? AND o.user_id = ?", id, uId).
		Select("o.id AS id," +
			"o.order_num AS order_num," +
			"UNIX_TIMESTAMP(o.created_at) AS created_at," +
			"UNIX_TIMESTAMP(o.updated_at) AS updated_at," +
			"o.user_id AS user_id," +
			"o.product_id AS product_id," +
			"o.boss_id AS boss_id," +
			"o.num AS num," +
			"o.type AS type," +
			"p.name AS name," +
			"p.discount_price AS discount_price," +
			"p.img_path AS img_path," +
			"a.name AS address_name," +
			"a.phone AS address_phone," +
			"a.address AS address").
		Find(&r).Error

	return
}

func (dao *OrderDao) DeleteOrderById(id, uId uint) error {
	return dao.DB.Model(&model.Order{}).
		Where("id = ? AND user_id = ?", id, uId).
		Delete(&model.Order{}).Error
}

func (dao *OrderDao) GetOrderById(id, uId uint) (r *model.Order, err error) {
	err = dao.DB.Model(&model.Order{}).
		Where("id = ? AND user_id = ?", id, uId).
		First(&r).Error

	return
}

// UpdateOrderById 更新订单详情
func (dao *OrderDao) UpdateOrderById(id, uId uint, order *model.Order) error {
	return dao.DB.Where("id = ? AND user_id = ?", id, uId).
		Updates(order).Error
}
