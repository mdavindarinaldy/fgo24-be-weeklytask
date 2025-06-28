package routers

import (
	"backend3/controllers"

	"github.com/gin-gonic/gin"
)

func authRouter(r *gin.RouterGroup) {
	r.POST("/register", controllers.AuthRegister)
	r.POST("/login", controllers.AuthLogin)
}
