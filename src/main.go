package main

import (
	"fmt"
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

func connectDB() *mongo.Collection {
	clientOptions := options.Client().ApplyURI("mongodb://root:example@localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	collection := client.Database("example").Collection("men")
	return collection
}

type MensTableWorker struct {
	client *mongo.Client
}

func NewMensTableWorker(client *mongo.Client) *MensTableWorker {
	return &MensTableWorker{client: client}
}

type Table struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	City string `json:"city"`
}

type Man struct {
	ID            int    `json:"id"`
	Email         string `json:"email"`
	PreferredName string `json:"preferredName"`
	TableID       int    `json:"tableID"`
}

type Event struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Location  string `json:"location"`
	StartTime string `json:"startTime"`
	TableID   int    `json:"tableID"`
}

var men = []Man{
}

var tables = []Table{
	{ID: 1, Name: "Table 1", City: "City 1"},
}

func getTables(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, tables)
}

func addMan(c *gin.Context) {
    var man Man

    // Decode the request body into the user struct
    if err := c.BindJSON(&man); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

	collection := connectDB()

    // Insert the user into the database
    result, err := collection.InsertOne(c, man)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Return the ID of the inserted document
    c.JSON(http.StatusOK, gin.H{"id": result.InsertedID})
}

func getMen(c *gin.Context) {
	collection := connectDB()

	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return
	}

	for cursor.Next(context.TODO()) {
		var man Man
		cursor.Decode(&man)
		men = append(men, man)
	}

	defer cursor.Close(context.TODO())

	c.IndentedJSON(http.StatusOK, men)
}

func main() {
	router := gin.Default()

	router.GET("/men", getMen)
	router.POST("/men", addMan)
	router.GET("/tables", getTables)
	router.Run("localhost:8080")
}
