package main

import (
	"gin-distributed/recipes"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.POST("/recipes", recipes.NewRecipeHandler)
	router.Run()
}
