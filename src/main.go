package main

import (
	"fmt"
	"context"
	"net/http"

	"example.com/mt/src/controllers"
	"example.com/mt/src/initialisers"
	"example.com/mt/src/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"

	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

var tables = []models.Table{
	{ID: 1, Name: "Table 1", City: "City 1"},
}

func getTables(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, tables)
}


func init() {
	// Load environment variables
	initialisers.LoadEnvVariables()
}

func main() {
	router := gin.Default()

	fmt.Println("Starting the application...")

	router.POST("/authenticate", controllers.Authenticate)
	router.GET("/men", controllers.GetMen)
	router.GET("/men/:email", controllers.GetManByEmail)
	router.POST("/men", controllers.AddMan)
	router.GET("/tables", getTables)
	router.Run("localhost:8080")
}
