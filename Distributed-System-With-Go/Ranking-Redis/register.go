package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type userRequest struct {
	Name string `json:"name"`
}

func register(c *gin.Context) {
	var newUser userRequest
	if err := c.BindJSON(&newUser); err != nil {
		return
	}
	print("err")
	user, err := addUser(newUser.Name)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		return
	}
	c.IndentedJSON(http.StatusCreated, user)
}
