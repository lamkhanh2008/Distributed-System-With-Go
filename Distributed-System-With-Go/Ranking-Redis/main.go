package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.DebugMode)
	router := gin.Default()
	router.GET("/", about)
	router.POST("/register", register)

	router.Run("0.0.0.0:8080")
}

func about(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{"application": "Ranking System"})
}
