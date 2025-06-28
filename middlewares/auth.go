package middlewares

import (
	"backend3/utils"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func VerifyToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		godotenv.Load() // ini baiknya dipindah ke main aja biar ga overhead
		secretKey := os.Getenv("APP_SECRET")
		token := strings.Split(ctx.GetHeader("Authorization"), "Bearer ")

		if len(token) < 2 {
			ctx.JSON(http.StatusUnauthorized, utils.Response{
				Success: false,
				Message: "Unauthorized",
			})
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		rawToken, err := jwt.Parse(token[1], func(t *jwt.Token) (any, error) {
			return []byte(secretKey), nil
		})

		if err != nil {
			ctx.JSON(http.StatusUnauthorized, utils.Response{
				Success: false,
				Message: "Invalid token",
			})
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		userId := rawToken.Claims.(jwt.MapClaims)["userId"]
		ctx.Set("userId", userId)
		ctx.Next()
	}
}
