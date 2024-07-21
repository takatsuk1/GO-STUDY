package upload

import (
	"fmt"
	"gin-mall/config"
	"gin-mall/pkg/utils/log"
	"io"
	"mime/multipart"
	"os"
	"strconv"
)

func ProductUploadToLocalStatic(file multipart.File, bossId uint, productName string, productNumber string) (filePath string, err error) {
	//将商家的id转为int
	bId := strconv.Itoa(int(bossId))
	//设置静态文件夹目录
	basePath := "." + config.Config.PhotoPath.ProductPath + "boss" + bId + "/" + productName + "/"
	//如果不存在则创建目录
	if !DirExistOrNot(basePath) {
		CreateDir(basePath)
	}
	//设置商品图片名字
	productPath := fmt.Sprintf("%s%s.jpg", productName, productNumber)
	//读取传来的图片文件
	content, err := io.ReadAll(file)
	if err != nil {
		log.LogrusObj.Error(err)
		return "", err
	}
	//写入文件
	err = os.WriteFile(basePath+productPath, content, 0666)
	if err != nil {
		log.LogrusObj.Error(err)
		return "", err
	}
	//返回文件目录
	return fmt.Sprintf("boss%s/%s", bId, productName), err
}
func CreateDir(dirName string) bool {
	//0755代表目录的权限设置，所有者可以读写和执行。其他人可以读和执行
	//MkdirAll递归创建目录
	err := os.MkdirAll(dirName, 0755)
	if err != nil {
		log.LogrusObj.Error(err)
		return false
	}
	return true
}

func DirExistOrNot(fileAddr string) bool {
	//os.Stat获取文件信息
	s, err := os.Stat(fileAddr)
	if err != nil {
		log.LogrusObj.Error(err)
		return false
	}
	//s.IsDir用于判断文件系统中的路径是否为目录。
	return s.IsDir()
}
