package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

var Router router

type router struct{}

func (*router) InitApiRouter(router *gin.Engine) {
	router.GET("/testapi", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "testapi success",
			"data": nil,
		})
	})
	// pod操作
	router.GET("/api/k8s/pods", Pod.GetPods)
	router.GET("/api/k8s/pods/containers/all-images", Pod.GetPodsContainersImages)
	// namespaces操作
	router.GET("/api/k8s/namespaces", Namespaces.GetNamespaces)
	router.GET("/api/k8s/namespaces/only-name", Namespaces.GetNamespacesOnlyName)
	// deployments操作
	router.GET("/api/k8s/deployments/only-name", Deployments.GetDeploymentsOnlyName)
	// statefulsets操作
	router.GET("/api/k8s/statefulsets/only-name", Statefulsets.GetStatefulsetsOnlyName)
}
