package main

import (
	// "fmt"
	// "net/http"
	"net/http"

	"github.com/gin-gonic/gin"
	// "go.mongodb.org/mongo-driver/bson/primitive"
	// "errors"
)

// import "go.mongodb.org/mongo-driver/bson/primitive"

type user struct {
    ID           int				`json:"id"`
    Email        string             `json:"email"`
    PreferredName string            `json:"preferredName"`
}

var users = []user {
	{ID: 1, Email: "test@example.com", PreferredName: "Test"},
}

func getUsers(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, users)
}

func main() {
	router := gin.Default()
	router.GET("/users", getUsers)
	router.Run("localhost:8080")
}
