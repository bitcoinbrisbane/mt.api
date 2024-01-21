package controllers


import (
	"context"
	"net/http"

	"example.com/mt/src/initialisers"
	"example.com/mt/src/models"
	"example.com/mt/src/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/mongo"
)

func Authenticate(c *gin.Context) {
	var login models.Login

	// Decode the request body into the user struct
	if c.BindJSON(&login) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	collection := initialisers.ConnectDB("men")

	var man models.Man
	err := collection.FindOne(context.TODO(), bson.M{"email": login.Email}).Decode(&man)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Man not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(man.Password), []byte(login.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	returnedToken, err := utils.GenerateToken(login.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": returnedToken})
}