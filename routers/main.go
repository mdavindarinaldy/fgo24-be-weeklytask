package routers

import "github.com/gin-gonic/gin"

func CombineRouter(r *gin.Engine) {
	userRouter(r.Group("/"))
	authRouter(r.Group("/auth"))
	transactionRouter(r.Group("/transactions"))
}
