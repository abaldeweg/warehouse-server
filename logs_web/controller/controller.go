package controller

import (
	"net/http"

	"github.com/abaldeweg/warehouse-server/logs_web/db"
	"github.com/gin-gonic/gin"
)

// GetEvents handles the GET request to retrieve logs.
func GetEvents(c *gin.Context) {
	var filter map[string]interface{}
	if err := c.ShouldBindJSON(&filter); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	h, _ := db.NewDBHandler()
	d, _ := h.FindDemanded(filter)
	defer h.Close()

	c.JSON(http.StatusOK, d)
}
