package service

import (
	"context"
	"errors"
	"gin-mall/config"
	"gin-mall/consts"
	"gin-mall/pkg/utils/ctl"
	"gin-mall/pkg/utils/log"
	"gin-mall/repository/db/dao"
	"gin-mall/repository/db/model"
	"gin-mall/types"
	"sync"
)

var FavoriteSrvIns *FavoriteSrv
var FavoriteSrvOnce sync.Once

type FavoriteSrv struct {
}

func GetFavoriteSrv() *FavoriteSrv {
	FavoriteSrvOnce.Do(func() {
		FavoriteSrvIns = &FavoriteSrv{}
	})
	return FavoriteSrvIns
}

func (s *FavoriteSrv) FavoriteCreate(c context.Context, req *types.FavoriteCreateReq) (resp interface{}, err error) {
	u, err := ctl.GetUserInfo(c)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	fDao := dao.NewFavoritesDao(c)
	exist, _ := fDao.FavoriteExistOrNot(req.ProductId, u.Id)
	if exist {
		err = errors.New("该商品已存在")
		log.LogrusObj.Error(err)
		return
	}
	userDao := dao.NewUserDao(c)
	user, err := userDao.GetUserById(u.Id)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}

	bossDao := dao.NewUserDaoByDB(userDao.DB)
	boss, err := bossDao.GetUserById(req.BossId)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}

	pDao := dao.NewProductDao(c)
	product, err := pDao.GetProductById(req.ProductId)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}

	favorite := &model.Favorite{
		UserID:    u.Id,
		User:      *user,
		Product:   *product,
		ProductID: req.ProductId,
		BossID:    req.BossId,
		Boss:      *boss,
	}
	err = fDao.CreateFavorite(favorite)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	return
}

func (s *FavoriteSrv) FavoriteList(c context.Context, req *types.FavoritesServiceReq) (resp interface{}, err error) {
	u, err := ctl.GetUserInfo(c)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	favorites, total, err := dao.NewFavoritesDao(c).ListFavoriteByUserId(u.Id, req.PageSize, req.PageNum)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	for i := range favorites {
		if config.Config.System.UploadModel == consts.UploadModelLocal {
			favorites[i].ImgPath = config.Config.PhotoPath.PhotoHost + config.Config.System.HttpPort + config.Config.PhotoPath.ProductPath + favorites[i].ImgPath
		}
	}
	resp = &types.DataListResp{
		Item:  favorites,
		Total: total,
	}
	return
}

func (s *FavoriteSrv) FavoriteDelete(c context.Context, req *types.FavoriteDeleteReq) (resp interface{}, err error) {
	favoriteDao := dao.NewFavoritesDao(c)
	err = favoriteDao.DeleteFavoriteById(req.Id)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	return
}
