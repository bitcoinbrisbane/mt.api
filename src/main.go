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
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
	// "example.com/mt/jwtutils"
	// "errors"
)

var jwtKey = []byte("your_secret_key") // Ideally, this should be loaded from a secure source like environment variables

type JWTClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// GenerateToken generates a new JWT token for a given username
func generateToken(username string) (string, error) {
	expirationTime := time.Now().Add(12 * time.Hour)
	claims := &JWTClaims{
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
func validateToken(tokenString string) (*JWTClaims, error) {
	claims := &JWTClaims{}

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

var men = []Man{}

var tables = []Table{
	{ID: 1, Name: "Table 1", City: "City 1"},
}

func connectDB() *mongo.Collection {
	clientOptions := options.Client().ApplyURI("mongodb://root:example@localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	collection := client.Database("example").Collection("men")
	return collection
}

type Login struct {
	Email    string `json:"name"`
	Password string `json:"password"`
}

type Table struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	City string `json:"city"`
}

type Man struct {
	Email         string `json:"email"`
	PreferredName string `json:"preferredName"`
	Password      string `json:"password"`
	TableID       int    `json:"tableID"`
}

type Event struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Location  string `json:"location"`
	StartTime string `json:"startTime"`
	TableID   int    `json:"tableID"`
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

func getUserByEmail(c *gin.Context) {
	email := c.Param("email") // Get the email from the URL parameter
	collection := connectDB()

	var man Man
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

func authenticate(c *gin.Context) {
	var login Login

	// Decode the request body into the user struct
	if err := c.BindJSON(&login); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	collection := connectDB()

	var man Man
	err := collection.FindOne(context.TODO(), bson.M{"email": login.Email}).Decode(&man)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Man not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	if man.Password != login.Password {
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

func main() {
	router := gin.Default()

	fmt.Println("Starting the application...")

	router.POST("/authenticate", authenticate)
	router.GET("/men", getMen)
	router.GET("/men/:email", getUserByEmail)
	router.POST("/men", addMan)
	router.GET("/tables", getTables)
	router.Run("localhost:8080")
}
