package controller

import (
	"net/http"
	"strconv"

	"github.com/abaldeweg/warehouse-server/logs/db"
	"github.com/gin-gonic/gin"
)

// GetLogs handles the GET request to retrieve logs.
func GetLogs(c *gin.Context) {
	fromParam := c.Param("from")
	toParam := c.Param("to")

	if len(fromParam) != 8 || len(toParam) != 8 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "'from' and 'to' parameters must be 8 characters long"})
		return
	}

	from, err := strconv.Atoi(fromParam)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid 'from' parameter"})
		return
	}

	to, err := strconv.Atoi(toParam)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid 'to' parameter"})
		return
	}

	h, _ := db.NewDBHandler()
	d, _ := h.Get(from, to)
	defer h.Close()

	c.JSON(http.StatusOK, d)
}
