package controllers

import (
	"backend3/models"
	"backend3/utils"
	"net/http"

	"github.com/gin-gonic/gin"
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
