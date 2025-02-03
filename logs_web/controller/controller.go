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

	h, err := db.NewDBHandler()
  if err != nil {
    c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
    return
  }
	defer h.Close()

	d, err := h.FindDemanded(filter)
  if err != nil {
    c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
    return
  }

	c.JSON(http.StatusOK, d)
}
