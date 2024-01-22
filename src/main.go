package main

import (
	// "context"
	"fmt"
	"net/http"

	"example.com/mt/src/controllers"
	"example.com/mt/src/initialisers"
	"example.com/mt/src/models"
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

var tables = []models.Table{
	{ID: 1, Name: "Table 1", City: "City 1"},
}


var manType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Man",
	Fields: graphql.Fields{
		"email": &graphql.Field{
			Type: graphql.String,
		},
		"perferredname": &graphql.Field{
			Type: graphql.String,
		},
		"id": &graphql.Field{
			Type: graphql.Int,
		},
	},
})

var rootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
	Fields: graphql.Fields{
		"man": &graphql.Field{
			Type:        manType,
			Description: "Get single man",
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				return nil, nil
			},
		},

		"men": &graphql.Field{
			Type:        graphql.NewList(manType),
			Description: "List of men",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return nil, nil
			},
		},
	},
})

// define schema, with our rootQuery
var MenSchema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query: rootQuery,
})

// Defining the Playground handler
func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/query")
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func graphQLHandler() gin.HandlerFunc {
	h := handler.New(&handler.Config{
		Schema:   &MenSchema,
		Pretty:   true,
		GraphiQL: false,
	})
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

	router.GET("/query", playgroundHandler())
	router.POST("/query", graphQLHandler())

	router.POST("/authenticate", controllers.Authenticate)
	router.GET("/men", controllers.GetMen)
	router.GET("/men/:email", controllers.GetManByEmail)
	router.POST("/men", controllers.AddMan)
	router.GET("/tables", getTables)

	router.Run("localhost:8080")
}
