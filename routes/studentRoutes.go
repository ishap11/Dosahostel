package routes

import (
	"github.com/adityjoshi/Dosahostel/controllers"
	kafkamanager "github.com/adityjoshi/Dosahostel/kafka/manager"
	"github.com/adityjoshi/Dosahostel/middleware"

	"github.com/gin-gonic/gin"
)

func StudentRoutes(incomingRoutes *gin.Engine, km *kafkamanager.KafkaManager) {
	incomingRoutes.POST("/student/register", controllers.BusinessAdminReg)
	incomingRoutes.GET("/student/login", controllers.StudentLogin)
	incomingRoutes.POST("/student/complaint", middleware.AuthorizeStudent(), controllers.PostInventory)
	incomingRoutes.POST("/student/bulk", func(c *gin.Context) {
		c.Set("km", km)
		controllers.PostComplaintKafka(c)
	})
}
