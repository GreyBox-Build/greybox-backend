package controllers

import (
	"backend/models"
	"github.com/gin-gonic/gin"
	"strconv"
)

func ListHurupayRequest(c *gin.Context) {
	requests, err := models.GetHurupayRequest()
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	serializedData := models.ConvertToSerializer(requests)
	c.JSON(200, gin.H{"status": "fetched hurupay requests", "data": serializedData, "errors": false})
}

func GetHurupayRequest(c *gin.Context) {
	id := c.Param("id")
	intId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	request, err := models.GetHurupayRequestById(int(intId))
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	serializedData := models.ConvertToSerializer(request)
	c.JSON(200, gin.H{"status": "fetched hurupay request", "data": serializedData, "errors": false})
}

func GetHurupayStats(c *gin.Context) {
	stats, err := models.GetHurupayStats()
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "fetched hurupay stats", "data": stats, "errors": false})
}
