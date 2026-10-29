package routers

import (
	"mcs_bab_7/controllers"

	"github.com/gin-gonic/gin"
)

func StartServer() *gin.Engine {
	routers := gin.Default()
	routers.POST("/servo/init-proj", controllers.InitProj)
	routers.GET("/servo/status", controllers.GetStatus)
	routers.PUT("/servo/update/:srv_status", controllers.UpdateStatus)
	return routers
}
