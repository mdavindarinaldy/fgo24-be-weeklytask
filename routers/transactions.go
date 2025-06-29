package routers

import (
	"backend3/controllers"
	"backend3/middlewares"

	"github.com/gin-gonic/gin"
)

func transactionRouter(r *gin.RouterGroup) {
	r.Use(middlewares.VerifyToken())
	r.POST("/top-up", controllers.TopUp)
	r.POST("/transfer", controllers.Transfer)
	r.GET("/history", controllers.HistoryTransaction)
}
