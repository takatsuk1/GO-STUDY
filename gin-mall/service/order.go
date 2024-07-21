package service

import (
	"context"
	"fmt"
	"gin-mall/config"
	"gin-mall/consts"
	"gin-mall/pkg/utils/ctl"
	"gin-mall/pkg/utils/log"
	"gin-mall/repository/cache"
	"gin-mall/repository/db/dao"
	"gin-mall/repository/db/model"
	"gin-mall/types"
	"github.com/redis/go-redis/v9"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

const OrderTimeKey = "OrderTime"

var OrderSrvIns *OrderSrv
var OrderSrvOnce sync.Once

type OrderSrv struct {
}

func GetOrderSrv() *OrderSrv {
	OrderSrvOnce.Do(func() {
		OrderSrvIns = &OrderSrv{}
	})
	return OrderSrvIns
}

func (s *OrderSrv) OrderCreate(c context.Context, req *types.OrderCreateReq) (resp interface{}, err error) {
	u, err := ctl.GetUserInfo(c)
	if err != nil {
		log.LogrusObj.Error(err)
		return
	}
	order := &model.Order{
		UserID:    u.Id,
		ProductID: req.ProductID,
		BossID:    req.BossID,
		Num:       int(req.Num),
		Money:     float64(req.Money),
		Type:      1,
	}
	addressDao := dao.NewAddressDao(c)
	address, err := addressDao.GetAddressByAid(req.AddressID, u.Id)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	order.AddressID = address.ID
	//rand.NewSource(time.Now().UnixNano()) 用当前时间的纳秒级时间戳作为种子，确保每次生成的随机数不同。
	//rand.new()创建随机数生成器
	//Int31n(1000000000) 生成0到 999999999 之间的随机整数。
	//使用 fmt.Sprintf 将生成的随机整数格式化为一个长度为9的字符串
	number := fmt.Sprintf("%09v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000000))
	productNum := strconv.Itoa(int(req.ProductID))
	userNum := strconv.Itoa(int(req.UserID))
	number = number + productNum + userNum
	orderNum, _ := strconv.ParseUint(number, 10, 64)
	order.OrderNum = orderNum

	orderDao := dao.NewOrderDao(c)
	err = orderDao.CreateOrder(order)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	//订单号存入redis
	//redis.Z 是 Go 语言的 Redis 客户端包 github.com/go-redis/redis/v8
	//中定义的一个结构体，
	//用于表示 Redis 有序集合（Sorted Set）中的一个元素。在 Redis 有序集合中，
	//每个元素都有一个成员（member）和一个与之相关的分数（score）。
	data := redis.Z{
		//Score 字段存储排序分数，这里是当前时间加上 15 分钟后的时间戳。
		Score: float64(time.Now().Unix()) + 15*time.Minute.Seconds(),
		//Member 字段存储成员的值，这里是 orderNum，假设是一个订单号。
		Member: orderNum,
	}
	cache.RedisClient.ZAdd(cache.RedisContext, OrderTimeKey, data)

	return
}

func (s *OrderSrv) OrderList(c context.Context, req *types.OrderListReq) (resp interface{}, err error) {
	u, err := ctl.GetUserInfo(c)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	orders, total, err := dao.NewOrderDao(c).ListOrderByCondition(u.Id, req)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	for i := range orders {
		if config.Config.System.UploadModel == consts.UploadModelLocal {
			orders[i].ImgPath = config.Config.PhotoPath.PhotoHost + config.Config.System.HttpPort + config.Config.PhotoPath.ProductPath + orders[i].ImgPath
		}
	}
	resp = types.DataListResp{
		Item:  orders,
		Total: total,
	}
	return
}

func (s *OrderSrv) OrderShow(c context.Context, req *types.OrderShowReq) (resp interface{}, err error) {
	u, err := ctl.GetUserInfo(c)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	order, err := dao.NewOrderDao(c).ShowOrderById(req.OrderId, u.Id)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	if config.Config.System.UploadModel == consts.UploadModelLocal {
		order.ImgPath = config.Config.PhotoPath.PhotoHost + config.Config.System.HttpPort + config.Config.PhotoPath.ProductPath + order.ImgPath
	}
	resp = order
	return
}

func (s *OrderSrv) OrderDelete(c context.Context, req *types.OrderDeleteReq) (resp interface{}, err error) {
	u, err := ctl.GetUserInfo(c)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	err = dao.NewOrderDao(c).DeleteOrderById(req.OrderId, u.Id)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	return
}
