package controller

import (
	"bfJenkinsApi/service"
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"net/http"
)

var Namespaces namespaces

type namespaces struct{}

// 获取namespaces列表
func (*namespaces) GetNamespaces(ctx *gin.Context) {
	params := new(struct {
		FilterName string `form:"filter_name"`
		Page       int    `form:"page"`
		Limit      int    `form:"limit"`
		K8sRegion  string `form:"k8s_region"`
	})

	if err := ctx.ShouldBind(params); err != nil {
		logger.Error("绑定请求参数失败," + err.Error())
		ctx.JSON(http.StatusOK, gin.H{
			"code": 90400,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	//业务逻辑处理
	data, err := service.Namespaces.GetNamespaces(ctx, params.K8sRegion, params.FilterName, params.Limit, params.Page)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 90500,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "获取Namespaces列表成功",
		"data": data,
	})
	return
}

// 获取namespaces列表，仅返回名称
func (*namespaces) GetNamespacesOnlyName(ctx *gin.Context) {
	params := new(struct {
		FilterName string `form:"filter_name"`
		Page       int    `form:"page"`
		Limit      int    `form:"limit"`
		K8sRegion  string `form:"k8s_region"`
	})

	if err := ctx.ShouldBind(params); err != nil {
		logger.Error("绑定请求参数失败," + err.Error())
		ctx.JSON(http.StatusOK, gin.H{
			"code": 90400,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	//业务逻辑处理
	data, err := service.Namespaces.GetNamespacesOnlyName(ctx, params.K8sRegion, params.FilterName, params.Limit, params.Page)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 90500,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "获取Namespaces列表成功",
		"data": data,
	})
	return
}
