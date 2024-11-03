package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func JSONResponse(c *gin.Context, status int, payload interface{}) {
	c.JSON(status, payload)
}

func Success(c *gin.Context, data interface{}) {
	JSONResponse(c, http.StatusOK, gin.H{"success": true, "data": data})
}

func Error(c *gin.Context, status int, message string) {
	JSONResponse(c, status, gin.H{"success": false, "error": message})
}