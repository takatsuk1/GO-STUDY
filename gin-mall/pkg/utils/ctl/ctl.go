package ctl

import (
	"gin-mall/pkg/e"
	"github.com/gin-gonic/gin"
)

// Response 基础序列化器
type Response struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
	Msg    string      `json:"msg"`
	Error  string      `json:"error"`
}

// RespSuccess 带data成功返回
func RespSuccess(ctx *gin.Context, data interface{}, code ...int) *Response {
	status := e.SUCCESS
	if code != nil {
		status = code[0]
	}

	if data == nil {
		data = "操作成功"
	}

	r := &Response{
		Status: status,
		Data:   data,
		Msg:    e.GetMsg(status),
	}

	return r
}
