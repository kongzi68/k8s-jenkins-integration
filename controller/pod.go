package controller

import (
	"bfJenkinsApi/service"
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"net/http"
)

var Pod pod

type pod struct{}

// 获取pod列表
func (*pod) GetPods(ctx *gin.Context) {
	params := new(struct {
		FilterName string `form:"filter_name"`
		Namespace  string `form:"namespace"`
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
	data, err := service.Pod.GetPods(ctx, params.K8sRegion, params.FilterName, params.Namespace, params.Limit, params.Page)
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
		"msg":  "获取Pod列表成功",
		"data": data,
	})
	return
}

// 获取pod中容器的镜像列表
func (*pod) GetPodsContainersImages(ctx *gin.Context) {
	params := new(struct {
		FilterName string `form:"filter_name"`
		Page       int    `form:"page"`
		Limit      int    `form:"limit"`
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
	data, err := service.Pod.GetPodsContainersImages(ctx, params.FilterName, params.Limit, params.Page)
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
		"msg":  "获取Pod中容器的镜像列表成功",
		"data": data,
	})
	return
}
