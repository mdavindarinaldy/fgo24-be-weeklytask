package routers

import (
	"backend3/controllers"
	"backend3/middlewares"

	"github.com/gin-gonic/gin"
)

func userRouter(r *gin.RouterGroup) {
	r.Use(middlewares.VerifyToken())
	r.PUT("/update-profile", controllers.UpdateProfile)
	r.GET("/:id", controllers.GetUser)
	r.GET("", controllers.GetAllUsers)
	r.GET("/get-balance", controllers.GetLatestBalance)
	r.GET("/get-income", controllers.GetTotalIncome)
	r.GET("/get-expense", controllers.GetTotalExpense)
}
