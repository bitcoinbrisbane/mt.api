package controllers

import (
	"context"
	"net/http"

	"example.com/mt/src/initialisers"
	"example.com/mt/src/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/mongo"
)

func AddMan(c *gin.Context) {
	var man models.Man

	// Decode the request body into the user struct
	if c.BindJSON(&man) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	var password, err = bcrypt.GenerateFromPassword([]byte(man.Password), 10)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}

	man.Password = string(password)
	collection := initialisers.ConnectDB("men")

	// Insert the user into the database
	result, err := collection.InsertOne(c, man)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the ID of the inserted document
	c.JSON(http.StatusOK, gin.H{"id": result.InsertedID})
}

func GetMen(c *gin.Context) {
	collection := initialisers.ConnectDB("men")

	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return
	}

	var men = []models.Man{}
	for cursor.Next(context.TODO()) {
		var man models.Man
		cursor.Decode(&man)
		men = append(men, man)
	}

	defer cursor.Close(context.TODO())

	c.IndentedJSON(http.StatusOK, men)
}

func GetManByEmail(c *gin.Context) {
	email := c.Param("email") // Get the email from the URL parameter
	collection := initialisers.ConnectDB("men")

	var man models.Man
	err := collection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&man)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Man not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, man)
}
