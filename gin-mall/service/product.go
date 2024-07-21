package service

import (
	"context"
	"gin-mall/config"
	"gin-mall/consts"
	"gin-mall/pkg/utils/ctl"
	"gin-mall/pkg/utils/log"
	"gin-mall/pkg/utils/upload"
	"gin-mall/repository/db/dao"
	"gin-mall/repository/db/model"
	"gin-mall/types"
	"mime/multipart"
	"strconv"
	"sync"
)

var ProductSrvIns *ProductSrv
var ProductSrvOnce sync.Once

type ProductSrv struct{}

func GetProductSrv() *ProductSrv {
	ProductSrvOnce.Do(func() {
		ProductSrvIns = &ProductSrv{}
	})
	return ProductSrvIns
}

func (s *ProductSrv) ProductCreate(c context.Context, files []*multipart.FileHeader, req *types.ProductCreateReq) (resp interface{}, err error) {
	u, err := ctl.GetUserInfo(c)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	uId := u.Id
	//获取设置商品的商家信息
	boss, _ := dao.NewUserDao(c).GetUserById(uId)
	//以第一张作为封面图，获取第一个文件
	tmp, _ := files[0].Open()
	var path string
	//将文件上传至本地静态文件夹
	if config.Config.System.UploadModel == consts.UploadModelLocal {
		path, err = upload.ProductUploadToLocalStatic(tmp, uId, req.Name, "0")
	}
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	//创建新商品实例
	product := &model.Product{
		Name:          req.Name,
		CategoryID:    req.CategoryID,
		Title:         req.Title,
		Info:          req.Info,
		ImgPath:       path,
		Price:         req.Price,
		DiscountPrice: req.DiscountPrice,
		Num:           req.Num,
		OnSale:        true,
		BossID:        uId,
		BossName:      boss.UserName,
	}
	productDao := dao.NewProductDao(c)
	//创建商品行
	err = productDao.CreateProduct(product)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}

	//处理并发操作，waitgroup可以确保所有的并发操作完成后再继续后续处理
	wg := new(sync.WaitGroup)
	wg.Add(len(files) - 1)
	for index := 1; index <= len(files)-1; index++ {
		file := files[index]
		num := strconv.Itoa(index)
		tmp, _ := file.Open()
		if config.Config.System.UploadModel == consts.UploadModelLocal {
			path, err = upload.ProductUploadToLocalStatic(tmp, uId, req.Name, num)
		}
		if err != nil {
			log.LogrusObj.Error(err)
			return nil, err
		}
		//这里处理的是一个商品的多个图片
		productImg := &model.ProductImg{
			ProductID: product.ID,
			ImgPath:   path,
		}
		err = dao.NewProductImgDaoByDB(productDao.DB).CreateProductImg(productImg)
		if err != nil {
			log.LogrusObj.Error(err)
			return nil, err
		}
		//减少waitgourp的计数，表示一个并发操作已经完成
		wg.Done()
	}
	//阻塞当前的goroutine，直到waitgroup的计数为0
	wg.Wait()
	return
}

func (s *ProductSrv) ProductUpdate(c context.Context, req *types.ProductUpdateReq) (resp interface{}, err error) {
	product := &model.Product{
		Name:       req.Name,
		CategoryID: req.CategoryID,
		Title:      req.Title,
		Info:       req.Info,
		//ImgPath: req.ImgPath,
		Price:         req.Price,
		DiscountPrice: req.DiscountPrice,
		OnSale:        req.OnSale,
		Num:           req.Num,
	}
	err = dao.NewProductDao(c).UpdateProduct(req.ID, product)
	if err != nil {
		log.LogrusObj.Error(err)
		return
	}
	return
}

func (s *ProductSrv) ProductDelete(c context.Context, req *types.ProductDeleteReq) (resp interface{}, err error) {
	u, _ := ctl.GetUserInfo(c)
	err = dao.NewProductDao(c).DeleteProduct(req.ID, u.Id)
	if err != nil {
		log.LogrusObj.Error(err)
		return
	}
	return
}

func (s *ProductSrv) ProductList(c context.Context, req *types.ProductListReq) (resp interface{}, err error) {
	var total int64
	//可以在condition和productListReq中增加额外搜索条件
	condition := make(map[string]interface{})
	if req.CategoryID != 0 {
		condition["category_id"] = req.CategoryID
	}
	productDao := dao.NewProductDao(c)
	products, _ := productDao.ListProductByCondition(condition, req.BasePage)
	total, err = productDao.CountProductByCondition(condition)
	if err != nil {
		log.LogrusObj.Error(err)
		return
	}
	pRespList := make([]*types.ProductResp, 0)
	for _, p := range products {
		pResp := &types.ProductResp{
			ID:            p.ID,
			Name:          p.Name,
			CategoryID:    p.CategoryID,
			Title:         p.Title,
			Info:          p.Info,
			ImgPath:       p.ImgPath,
			Price:         p.Price,
			DiscountPrice: p.DiscountPrice,
			//View: p.View(),
			CreatedAt: p.CreatedAt.Unix(),
			Num:       p.Num,
			OnSale:    p.OnSale,
			BossID:    p.BossID,
			BossName:  p.BossName,
		}
		if config.Config.System.UploadModel == consts.UploadModelLocal {
			pResp.ImgPath = config.Config.PhotoPath.PhotoHost + config.Config.System.HttpPort + config.Config.PhotoPath.ProductPath + pResp.ImgPath
		}
		pRespList = append(pRespList, pResp)
	}
	resp = &types.DataListResp{
		Item:  pRespList,
		Total: total,
	}
	return
}
func (s *ProductSrv) ProductShow(c context.Context, req *types.ProductShowReq) (resp interface{}, err error) {
	p, err := dao.NewProductDao(c).ShowProductById(req.ID)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	pResp := &types.ProductResp{
		ID:            p.ID,
		Name:          p.Name,
		CategoryID:    p.CategoryID,
		Title:         p.Title,
		Info:          p.Info,
		ImgPath:       p.ImgPath,
		Price:         p.Price,
		DiscountPrice: p.DiscountPrice,
		CreatedAt:     p.CreatedAt.Unix(),
		Num:           p.Num,
		OnSale:        p.OnSale,
		BossID:        p.BossID,
		BossName:      p.BossName,
		//View : p.View(),
	}
	if config.Config.System.UploadModel == consts.UploadModelLocal {
		pResp.ImgPath = config.Config.PhotoPath.PhotoHost + config.Config.System.HttpPort + config.Config.PhotoPath.ProductPath + pResp.ImgPath
	}
	resp = pResp
	return
}

func (s *ProductSrv) ProductSearch(c context.Context, req *types.ProductSearchReq) (resp interface{}, err error) {
	product, count, err := dao.NewProductDao(c).SearchProduct(req.Name, req.Info, req.BasePage)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	pRespList := make([]*types.ProductResp, 0)
	for _, p := range product {
		pResp := &types.ProductResp{
			ID:            p.ID,
			Name:          p.Name,
			CategoryID:    p.CategoryID,
			Title:         p.Title,
			Info:          p.Info,
			ImgPath:       p.ImgPath,
			Price:         p.Price,
			DiscountPrice: p.DiscountPrice,
			CreatedAt:     p.CreatedAt.Unix(),
			Num:           p.Num,
			OnSale:        p.OnSale,
			BossID:        p.BossID,
			BossName:      p.BossName,
			//View: p.View(),
		}
		if config.Config.System.UploadModel == consts.UploadModelLocal {
			pResp.ImgPath = config.Config.PhotoPath.PhotoHost + config.Config.System.HttpPort + config.Config.PhotoPath.ProductPath + pResp.ImgPath
		}
		pRespList = append(pRespList, pResp)
	}
	resp = &types.DataListResp{
		Item:  pRespList,
		Total: count,
	}
	return
}

func (s *ProductSrv) ProductImgList(c context.Context, req *types.ListProductImgReq) (resp interface{}, err error) {
	productImgs, err := dao.NewProductImgDao(c).ListProductImgByProductId(req.ID)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	for i := range productImgs {
		if config.Config.System.UploadModel == consts.UploadModelLocal {
			productImgs[i].ImgPath = config.Config.PhotoPath.PhotoHost + config.Config.System.HttpPort + config.Config.PhotoPath.ProductPath + productImgs[i].ImgPath
		}
	}
	resp = &types.DataListResp{
		Item:  productImgs,
		Total: int64(len(productImgs)),
	}
	return
}
