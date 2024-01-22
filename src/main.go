package main

import (
	// "context"
	"fmt"
	"net/http"

	"example.com/mt/src/controllers"
	"example.com/mt/src/initialisers"
	"example.com/mt/src/models"
	"github.com/gin-gonic/gin"
	// "github.com/graphql-go/graphql"
	// "github.com/graphql-go/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

var tables = []models.Table{
	{ID: 1, Name: "Table 1", City: "City 1"},
}

// Defining the Playground handler
func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/query")
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
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

	router.GET("/", playgroundHandler())
	// router.POST("/graphql", http.GraphqlHandler())

	router.POST("/authenticate", controllers.Authenticate)
	router.GET("/men", controllers.GetMen)
	router.GET("/men/:email", controllers.GetManByEmail)
	router.POST("/men", controllers.AddMan)
	router.GET("/tables", getTables)

	router.Run("localhost:8080")
}
