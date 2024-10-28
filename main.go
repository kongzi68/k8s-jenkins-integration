package main

import (
	"bfJenkinsApi/config"
	"bfJenkinsApi/controller"
	"bfJenkinsApi/service"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	for _, val := range config.K8SRegion {
		service.K8s.Init(val)
	}
	controller.Router.InitApiRouter(r)
	err := r.Run(config.ListenAddr)
	if err != nil {
		return
	} // listen and serve on iamIPaddr:8080 (for windows "localhost:8080")
}
