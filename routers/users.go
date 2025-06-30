package routers

import (
	"backend3/controllers"
	"backend3/middlewares"

	"github.com/gin-gonic/gin"
)

func userRouter(r *gin.RouterGroup) {
	r.Use(middlewares.VerifyToken())
	r.PUT("profile", controllers.UpdateProfile)
	r.GET("users", controllers.GetAllUsers)
	r.GET("balance", controllers.GetLatestBalance)
	r.GET("income", controllers.GetTotalIncome)
	r.GET("expense", controllers.GetTotalExpense)
	r.POST("logout", controllers.Logout)
	// r.GET("/:id", controllers.GetUser)
}
