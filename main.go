package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

var client *mongo.Client

func main() {
	// Use context with timeout for MongoDB connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// MongoDB URI (Ensure to replace `<username>` and `<password>` with your actual credentials)
	mongoURI := "mongodb+srv://anshulagnihotri008:LQ9NDKyg3Uv59Utk@cluster0.yjrxe.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"

	// Connect to MongoDB
	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}

	// Verify MongoDB connection
	if err = client.Ping(ctx, nil); err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB!")

	// Initialize Gin router
	router := gin.Default()

	// Middleware and routes
	router.Use(gin.Logger(), gin.Recovery())
	router.POST("/signup", Signup)
	router.POST("/login", Login)

	// Start the server
	router.Run(":8080")
}

// Signup handler
func Signup(c *gin.Context) {
	type User struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var user User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash the password before storing it in the database
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Store the user in the MongoDB collection
	collection := client.Database("mydb").Collection("users")
	_, err = collection.InsertOne(c, bson.M{"username": user.Username, "password": string(hashedPassword)})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

// Login handler
func Login(c *gin.Context) {
	type User struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var user User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Retrieve user from the database
	collection := client.Database("mydb").Collection("users")
	var storedUser bson.M
	err := collection.FindOne(c, bson.M{"username": user.Username}).Decode(&storedUser)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Compare the stored hashed password with the provided password
	err = bcrypt.CompareHashAndPassword([]byte(storedUser["password"].(string)), []byte(user.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}
