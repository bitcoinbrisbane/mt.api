package main

import (
	// "fmt"
	// "net/http"
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	// "errors"
)

type MensTableWorker struct {
	client *mongo.Client
}

func NewMensTableWorker(client *mongo.Client) *MensTableWorker {
	return &MensTableWorker{client: client}
}

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

	collection := ConnectDB()

	cursor, err := collection.Find(context.TODO(), bson.M{})
    if err != nil {
        return
    }

    defer cursor.Close(context.TODO())

	c.IndentedJSON(http.StatusOK, men)
}

func ConnectDB() *mongo.Collection {
    clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
    client, err := mongo.Connect(context.TODO(), clientOptions)
    if err != nil {
        log.Fatal(err)
    }

    collection := client.Database("your_db_name").Collection("users")
    return collection
}

func main() {
	router := gin.Default()

	router.GET("/men", getMen)
	router.GET("/tables", getTables)
	router.Run("localhost:8080")
}
