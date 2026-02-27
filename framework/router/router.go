package router

import (
	"github.com/gin-gonic/gin"
)

// Engine creates a gin engine with CORS and sets it to release mode.
func Engine() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	return r
}
