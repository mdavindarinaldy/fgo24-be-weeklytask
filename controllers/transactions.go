package controllers

import (
	"backend3/models"
	"backend3/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func TopUp(c *gin.Context) {
	topUp := models.TopUpRequest{}
	c.ShouldBind(&topUp)
	userId, _ := c.Get("userId")
	err := models.HandleTopUp(topUp, int(userId.(float64)))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.Response{
			Success: false,
			Message: "Top up failed! Please try again",
		})
		return
	}
	c.JSON(http.StatusOK, utils.Response{
		Success: true,
		Message: "Top up success",
	})
}

func Transfer(c *gin.Context) {
	transfer := models.TransferRequest{}
	c.ShouldBind(&transfer)
	userId, _ := c.Get("userId")
	err := models.HandleTransfer(transfer, int(userId.(float64)))
	if err != nil {
		if err.Error() == "insufficient balance" {
			c.JSON(http.StatusBadRequest, utils.Response{
				Success: false,
				Message: err.Error(),
			})
			return
		}
		c.JSON(http.StatusBadRequest, utils.Response{
			Success: false,
			Message: "Transfer failed! Please try again",
		})
		return
	}
	c.JSON(http.StatusOK, utils.Response{
		Success: true,
		Message: "Transfer success",
	})
}
