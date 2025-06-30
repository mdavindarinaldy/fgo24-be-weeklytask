package controllers

import (
	"backend3/models"
	"backend3/utils"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func UpdateProfile(c *gin.Context) {
	user := models.User{}
	c.ShouldBind(&user)
	userId, _ := c.Get("userId")
	err := models.HandleUpdate(user, int(userId.(float64)))
	if err != nil {
		if err.Error() == "email already used by another user" || err.Error() == "user data should not be empty" {
			c.JSON(http.StatusBadRequest, utils.Response{
				Success: false,
				Message: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, utils.Response{
			Success: false,
			Message: "Internal server error",
		})
		return
	}
	c.JSON(http.StatusOK, utils.Response{
		Success: true,
		Message: "Update profile success",
	})
}

// func GetUser(c *gin.Context) {
// 	id, _ := strconv.Atoi(c.Param("id"))
// 	user, err := models.GetDetailUser(id)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, utils.Response{
// 			Success: false,
// 			Message: "Failed to get detail user",
// 		})
// 		return
// 	}
// 	c.JSON(http.StatusOK, utils.Response{
// 		Success: true,
// 		Message: "Success to get detail user",
// 		Result: struct {
// 			Id          int
// 			Name        string
// 			Email       string
// 			PhoneNumber string
// 		}{
// 			Id:          user.Id,
// 			Name:        user.Name,
// 			Email:       user.Email,
// 			PhoneNumber: user.PhoneNumber,
// 		},
// 	})
// }

func GetAllUsers(c *gin.Context) {
	search := strings.ToLower(c.DefaultQuery("search", "a"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	users, pageData, err := models.GetAllUsers(search, page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Success: false,
			Message: "Internal server error",
		})
		return
	}
	c.JSON(http.StatusOK, utils.Response{
		Success: true,
		Message: "Success to get users",
		PageInfo: models.PageData{
			CurrentPage: pageData.CurrentPage,
			TotalPage:   pageData.TotalPage,
			TotalData:   pageData.TotalData,
		},
		Result: users,
	})
}

func GetLatestBalance(c *gin.Context) {
	userId, _ := c.Get("userId")
	balance := models.GetLatestBalance(int(userId.(float64)))
	c.JSON(http.StatusOK, utils.Response{
		Success: true,
		Message: "Success to get user's balance",
		Result: struct {
			Balance float64 `json:"balance"`
		}{
			Balance: balance,
		},
	})
}

func GetTotalIncome(c *gin.Context) {
	userId, _ := c.Get("userId")
	income, now, duration := models.GetTotalIncome(int(userId.(float64)))
	c.JSON(http.StatusOK, utils.Response{
		Success: true,
		Message: "Success to get user's income",
		Result: struct {
			Income   float64 `json:"income"`
			Duration any     `json:"duration"`
		}{
			Income: income,
			Duration: struct {
				TimeStart time.Time `json:"timeStart"`
				TimeEnd   time.Time `json:"timeEnd"`
			}{
				TimeStart: duration,
				TimeEnd:   now,
			},
		},
	})
}

func GetTotalExpense(c *gin.Context) {
	userId, _ := c.Get("userId")
	expense, now, duration := models.GetTotalExpense(int(userId.(float64)))
	c.JSON(http.StatusOK, utils.Response{
		Success: true,
		Message: "Success to get user's expense",
		Result: struct {
			Expense  float64 `json:"expense"`
			Duration any     `json:"duration"`
		}{
			Expense: expense,
			Duration: struct {
				TimeStart time.Time `json:"timeStart"`
				TimeEnd   time.Time `json:"timeEnd"`
			}{
				TimeStart: duration,
				TimeEnd:   now,
			},
		},
	})
}

func Logout(c *gin.Context) {
	secretKey := os.Getenv("APP_SECRET")
	token := strings.Split(c.GetHeader("Authorization"), "Bearer ")
	rawToken, _ := jwt.Parse(token[1], func(t *jwt.Token) (any, error) {
		return []byte(secretKey), nil
	})
	var expiresAt time.Time = time.Unix(int64(rawToken.Claims.(jwt.MapClaims)["exp"].(float64)), 0)

	err := models.AddToBlacklist(token[1], expiresAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Success: false,
			Message: "Failed to blacklist token",
		})
		return
	}

	c.JSON(http.StatusOK, utils.Response{
		Success: true,
		Message: "Logout successful",
	})
}
