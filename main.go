package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func main() {
	var err error
	client, err = mongo.Connect(nil, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()

	router.POST("/signup", Signup)
	router.POST("/login", Login)

	router.Run(":8080")
}

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

	collection := client.Database("mydb").Collection("users")
	_, err := collection.InsertOne(c, bson.M{"username": user.Username, "password": user.Password})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

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

	collection := client.Database("mydb").Collection("users")
	var result bson.M
	err := collection.FindOne(c, bson.M{"username": user.Username, "password": user.Password}).Decode(&result)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}
