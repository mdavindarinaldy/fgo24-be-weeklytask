package controllers

import (
	"backend3/models"
	"backend3/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthRegister(c *gin.Context) {
	user := models.User{}
	c.ShouldBind(&user)
	err := models.HandleRegister(user)
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
			Message: "Failed to register user",
		})
		return
	}
	c.JSON(http.StatusOK, utils.Response{
		Success: true,
		Message: "Create user success!",
	})
}

func AuthLogin(c *gin.Context) {
	user := struct {
		Email    string `form:"email" json:"email" db:"email" binding:"required,email"`
		Password string `form:"password" json:"password" db:"password" binding:"required"`
		Pin      string `form:"pin" json:"pin" db:"pin" binding:"required"`
	}{}
	c.ShouldBind(&user)
	userData, err := models.GetUserByEmail(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Success: false,
			Message: "Internal server error",
		})
		return
	}
	if userData == (models.User{}) {
		c.JSON(http.StatusBadRequest, utils.Response{
			Success: false,
			Message: "User is not registered",
		})
		return
	}
	if userData.Email == user.Email && userData.Password == user.Password && userData.Pin == user.Pin {
		token, err := GenerateToken(userData)
		if err != nil {
			c.JSON(http.StatusBadRequest, utils.Response{
				Success: false,
				Message: "Failed to generate token",
			})
			return
		}
		c.JSON(http.StatusOK, utils.Response{
			Success: true,
			Message: "Login success!",
			Result:  token,
		})
		return
	} else {
		c.JSON(http.StatusBadRequest, utils.Response{
			Success: false,
			Message: "Password and/or PIN is wrong!",
		})
	}
}
