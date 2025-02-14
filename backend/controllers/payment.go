package controllers

import (
	"backend/apis"
	"backend/models"
	"backend/serializers"
	"backend/utils/tokens"

	"github.com/gin-gonic/gin"
)

func BorderLessOnramp(c *gin.Context) {
	userId, err := tokens.ExtractUserID(c)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}
	user, err := models.GetUserByID(userId)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	var input serializers.BorderlessOnramp
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	borderless := apis.NewBorderless()
	response, err := borderless.MakeDeposit(
		input.Amount, input.Asset, input.Country, input.Fiat)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	borderlessRequest := models.BorderlessRequest{}
	borderlessRequest.FiatAmount = input.Amount
	borderlessRequest.Asset = input.Asset     // change this for prod use
	borderlessRequest.Country = input.Country //change this for prod use
	borderlessRequest.UserId = userId
	borderlessRequest.User = user
	borderlessRequest.Status = "Pending"
	borderlessRequest.AccountId = response.Destination.AccountID
	borderlessRequest.TxId = response.ID
	borderlessRequest.FeeAmount = response.FeeAmount
	if err := models.CreateBorderlessRequest(&borderlessRequest); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	transactionInstruction, err := borderless.GetTransaction(response.ID)
	c.JSON(200, gin.H{"data": transactionInstruction, "status": "success", "errors": false})
}
