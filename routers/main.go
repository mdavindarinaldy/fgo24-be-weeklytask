package routers

import "github.com/gin-gonic/gin"

func CombineRouter(r *gin.Engine) {
	userRouter(r.Group("/users"))
	authRouter(r.Group("/auth"))
	transactionRouter(r.Group("/transaction"))
}
