package main

import (
	"ginEssential/controller"

	"github.com/gin-gonic/gin"
)

func CollectRoute(r *gin.Engine) *gin.Engine {
	r.POST("/register", controller.Register)

	return r
}
