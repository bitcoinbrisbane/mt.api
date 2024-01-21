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

var jwtKey = []byte("your_secret_key") // Ideally, this should be loaded from a secure source like environment variables


// GenerateToken generates a new JWT token for a given username
func generateToken(username string) (string, error) {
	expirationTime := time.Now().Add(12 * time.Hour)
	claims := &models.JWTClaims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	return tokenString, err
}

// ValidateToken validates the JWT token and returns the user claims
func validateToken(tokenString string) (*models.JWTClaims, error) {
	claims := &models.JWTClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, err //jwt.NewValidationError("Invalid token")
	}

	return claims, nil
}

var tables = []models.Table{
	{ID: 1, Name: "Table 1", City: "City 1"},
}

func getTables(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, tables)
}

func authenticate(c *gin.Context) {
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

	returnedToken, err := generateToken(login.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": returnedToken})
}

func init() {
	// Load environment variables
	initialisers.LoadEnvVariables()
}

func main() {
	router := gin.Default()

	fmt.Println("Starting the application...")

	router.POST("/authenticate", authenticate)
	router.GET("/men", controllers.GetMen)
	router.GET("/men/:email", controllers.GetManByEmail)
	router.POST("/men", controllers.AddMan)
	router.GET("/tables", getTables)
	router.Run("localhost:8080")
}
