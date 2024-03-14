package routes

import (
	"github.com/dzulfiikar/middle-backend-programmer-test/cmd/controllers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.RouterGroup) {
	router.POST("/png-to-jpeg", controllers.PngToJpegConverter)
	router.POST("/resize-image", controllers.ResizeImage)
}
