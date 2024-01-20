package main

import (
	// "fmt"
	// "net/http"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	// "errors"
)

type table struct {
	ID          int					`json:"id"`
	Name        string            	`json:"name"`
	City	 	string				`json:"city"`
}

type man struct {
    ID           int				`json:"id"`
    Email        string             `json:"email"`
    PreferredName string            `json:"preferredName"`
}

type event struct {
	ID           int				`json:"id"`
	Name         string             `json:"name"`
	Location     string             `json:"location"`
	StartTime    string             `json:"startTime"`
	TableID	 	 int				`json:"tableID"`
}

var men = []man {
	{ID: 1, Email: "test@example.com", PreferredName: "Test"},
}

func getTables(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, men)
}

func getMen(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, men)
}

func main() {
	router := gin.Default()

	router.GET("/men", getMen)
	router.GET("/tables", getTables)
	router.Run("localhost:8080")
}
