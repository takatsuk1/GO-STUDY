package service

import (
	"context"
	"errors"
	"fmt"
	"gin-mall/consts"
	"gin-mall/pkg/utils/ctl"
	"gin-mall/pkg/utils/log"
	"gin-mall/repository/db/dao"
	"gin-mall/repository/db/model"
	"gin-mall/types"
	"gorm.io/gorm"
	"sync"
)

var PaymentSrvIns *PaymentSrv
var PaymentSrvOnce sync.Once

type PaymentSrv struct {
}

func GetPaymentSrv() *PaymentSrv {
	PaymentSrvOnce.Do(func() {
		PaymentSrvIns = &PaymentSrv{}
	})
	return PaymentSrvIns
}

func (s *PaymentSrv) PayDown(ctx context.Context, req *types.PaymentDownReq) (resp interface{}, err error) {
	u, err := ctl.GetUserInfo(ctx)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	//transactiono事务处理方法，接受一个函数作文参数，
	//事务用于确保数据的完整性，通过将一组操作作为一个单元来执行，
	//如果其中任何一部分失败，
	//则可以回滚到事务开始前的状态，以保证数据一致性。
	err = dao.NewOrderDao(ctx).Transaction(func(tx *gorm.DB) error {
		uId := u.Id
		//获取订单信息
		payment, err := dao.NewOrderDaoByDB(tx).GetOrderById(req.OrderId, uId)
		if err != nil {
			log.LogrusObj.Error(err)
			return err
		}
		//获取订单价格与数量
		money := payment.Money
		num := payment.Num
		money = money * float64(num)
		//获取顾客信息
		userDao := dao.NewUserDaoByDB(tx)
		user, err := userDao.GetUserById(uId)
		if err != nil {
			log.LogrusObj.Error(err)
			return err
		}

		// 对钱进行解密。减去订单。再进行加密。
		moneyFloat, err := user.DecryptMoney(req.Key)
		if err != nil {
			log.LogrusObj.Error(err)
			return err
		}
		if moneyFloat-money < 0.0 { // 金额不足进行回滚
			log.LogrusObj.Error(err)
			return errors.New("金币不足")
		}
		//获得剩余金额
		finMoney := fmt.Sprintf("%f", moneyFloat-money)
		user.Money = finMoney
		//重新对钱进行加密
		user.Money, err = user.EncryptMoney(req.Key)
		if err != nil {
			log.LogrusObj.Error(err)
			return err
		}

		err = userDao.UpdateUserById(uId, user)
		if err != nil { // 更新用户金额失败，回滚
			log.LogrusObj.Error(err)
			return err
		}
		//获取商家金额
		boss, err := userDao.GetUserById(uint(req.BossID))
		if err != nil {
			log.LogrusObj.Error(err)
			return err
		}
		//对商家金额进行增加
		moneyFloat, _ = boss.DecryptMoney(req.Key)
		finMoney = fmt.Sprintf("%f", moneyFloat+money)
		boss.Money = finMoney
		boss.Money, err = boss.EncryptMoney(req.Key)
		if err != nil {
			log.LogrusObj.Error(err)
			return err
		}

		err = userDao.UpdateUserById(uint(req.BossID), boss)
		if err != nil { // 更新boss金额失败，回滚
			log.LogrusObj.Error(err)
			return err
		}

		productDao := dao.NewProductDaoByDB(tx)
		product, err := productDao.GetProductById(uint(req.ProductID))
		if err != nil {
			log.LogrusObj.Error(err)
			return err
		}
		product.Num -= num
		err = productDao.UpdateProduct(uint(req.ProductID), product)
		if err != nil { // 更新商品数量减少失败，回滚
			log.LogrusObj.Error(err)
			return err
		}

		// 更新订单状态
		payment.Type = consts.OrderTypePendingShipping
		err = dao.NewOrderDaoByDB(tx).UpdateOrderById(req.OrderId, uId, payment)
		if err != nil { // 更新订单失败，回滚
			log.LogrusObj.Error(err)
			return err
		}

		productUser := model.Product{
			Name:          product.Name,
			CategoryID:    product.CategoryID,
			Title:         product.Title,
			Info:          product.Info,
			ImgPath:       product.ImgPath,
			Price:         product.Price,
			DiscountPrice: product.DiscountPrice,
			Num:           num,
			OnSale:        false,
			BossID:        uId,
			BossName:      user.UserName,
			//BossAvatar:    user.Avatar,
		}

		err = productDao.CreateProduct(&productUser)
		if err != nil { // 买完商品后创建成了自己的商品失败。订单失败，回滚
			log.LogrusObj.Error(err)
			return err
		}

		return nil

	})

	if err != nil {
		log.LogrusObj.Error(err)
		return
	}
	return
}
