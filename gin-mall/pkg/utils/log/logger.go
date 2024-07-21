package log

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"path"
	"time"
)

var LogrusObj *logrus.Logger

func InitLog() {
	//判断当前logrusObj是否存在
	if LogrusObj != nil {
		src, _ := setOutputFile()
		LogrusObj.Out = src
		return
	}
	//实例化一个新的日志记录器logger，并设置输出文件
	logger := logrus.New()
	src, _ := setOutputFile()
	//设置输出
	logger.Out = src
	//设置输出等级
	logger.SetLevel(logrus.DebugLevel)
	//设置输出时间格式
	logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	LogrusObj = logger
}
func setOutputFile() (*os.File, error) {
	//设置当前时间为now
	now := time.Now()
	//使用os.Getwd()获取当前目录，并将日志目录设置为当前目录
	logFilePath := ""
	if dir, err := os.Getwd(); err == nil {
		logFilePath = dir + "/logs/"
	}
	fmt.Println(logFilePath)
	//使用os.Stat()检查当前目录是否存在日志文件夹，不存在则使用MkdirAll创建日志文件夹
	_, err := os.Stat(logFilePath)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(logFilePath, 0777); err != nil {
			log.Println(err.Error())
			return nil, err
		}
	}
	//使用Format设置日志输出文件名字的格式
	logFileName := now.Format("2006-01-02") + ".log"
	fileName := path.Join(logFilePath, logFileName)
	//使用Stat检查是否存在日志文件，不存在则使用os.Create()创建
	if _, err := os.Stat(fileName); err != nil {
		if _, err := os.Create(fileName); err != nil {
			log.Println(err.Error())
			return nil, err
		}
	}
	//使用os.OpenFile()打开日志文件
	src, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return src, nil
}
