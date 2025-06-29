package middlewares

import (
	"backend3/models"
	"backend3/utils"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func VerifyToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
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

		isBlacklisted, err := models.IsTokenBlacklisted(token[1])
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, utils.Response{
				Success: false,
				Message: "Internal server error",
			})
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if isBlacklisted {
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
