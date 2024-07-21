package service

import (
	"context"
	"encoding/json"
	"fmt"
	"gin-mall/pkg/utils/log"
	"gin-mall/repository/cache"
	"gin-mall/repository/db/dao"
	"gin-mall/repository/db/model"
	"gin-mall/types"
	"sync"
)

var SkillProductSrvIns *SkillProductSrv
var SkillProductSrvOnce sync.Once

type SkillProductSrv struct {
}

func GetSkillProductSrv() *SkillProductSrv {
	SkillProductSrvOnce.Do(func() {
		SkillProductSrvIns = &SkillProductSrv{}
	})
	return SkillProductSrvIns
}

func (s *SkillProductSrv) InitSkillGoods(c context.Context) (resp interface{}, err error) {
	spList := make([]*model.SkillProduct, 0)
	for i := 1; i < 10; i++ {
		spList = append(spList, &model.SkillProduct{
			ProductId: uint(i),
			BossId:    2,
			Title:     "秒杀商品测试使用",
			Money:     200,
			Num:       10,
		})
	}
	err = dao.NewSkillGoodsDao(c).BatchCreate(spList)
	if err != nil {
		log.LogrusObj.Infoln(err)
		return
	}
	for i := range spList {
		jsonBytes, errx := json.Marshal(spList[i])
		if errx != nil {
			log.LogrusObj.Infoln(errx)
			return
		}
		jsonString := string(jsonBytes)
		_, errx = cache.RedisClient.LPush(c, cache.SkillProductListKey, jsonString).Result()
		if errx != nil {
			log.LogrusObj.Infoln(errx)
			return nil, errx
		}
	}

	return
}

// ListSkillGoods 列表展示
func (s *SkillProductSrv) ListSkillGoods(ctx context.Context) (resp interface{}, err error) {
	// 读缓存
	rc := cache.RedisClient
	// 获取列表
	skillProductList, err := rc.LRange(ctx, cache.SkillProductListKey, 0, -1).Result()
	if err != nil {
		log.LogrusObj.Infoln(err)
		return
	}

	if len(skillProductList) == 0 {
		skill, errx := dao.NewSkillGoodsDao(ctx).ListSkillGoods()
		if errx != nil {
			log.LogrusObj.Infoln(errx)
			return nil, errx
		}

		for i := range skill {
			// 将结构体转换为JSON格式的字符串
			jsonBytes, errx := json.Marshal(skill[i])
			if errx != nil {
				log.LogrusObj.Infoln(errx)
				return
			}
			// 将字节数组转换为字符串
			jsonString := string(jsonBytes)
			_, errx = rc.LPush(ctx, cache.SkillProductListKey, jsonString).Result()
			if errx != nil {
				log.LogrusObj.Infoln(errx)
				return nil, errx
			}
		}
		resp = skill
	} else {
		resp = skillProductList
	}

	return
}

// GetSkillGoods 详情展示
func (s *SkillProductSrv) GetSkillGoods(ctx context.Context, req *types.GetSkillProductReq) (resp interface{}, err error) {
	// 读缓存
	rc := cache.RedisClient
	// 获取列表
	resp, err = rc.Get(ctx,
		fmt.Sprintf(cache.SkillProductKey, req.ProductId)).Result()
	if err != nil {
		log.LogrusObj.Infoln(err)
		return
	}

	return
}

// SkillProduct 秒杀商品
func (s *SkillProductSrv) SkillProduct(ctx context.Context, req *types.SkillProductReq) (resp interface{}, err error) {
	// 读缓存
	rc := cache.RedisClient
	// 获取数据
	resp, err = rc.Get(ctx,
		fmt.Sprintf(cache.SkillProductKey, req.ProductId)).Result()
	if err != nil {
		log.LogrusObj.Infoln(err)
		return
	}

	return
}
