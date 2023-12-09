package controllers

import (
	"context"
	"gin-mongo-api/src/database"
	"gin-mongo-api/src/models"
	"gin-mongo-api/src/res"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

var userCollection *mongo.Collection = database.GetCollection(database.DB, "users")
var validate = validator.New()

func Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var user models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, res.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if validationErr := validate.Struct(&user); validationErr != nil {
			c.JSON(http.StatusBadRequest, res.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		if user.Password != user.CPassword {
			c.JSON(http.StatusBadRequest, res.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": "Password and Confirm Password do not match"}})
			return
		}

		existingUser := models.User{}
		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&existingUser)
		if err == nil {
			c.JSON(http.StatusConflict, res.UserResponse{Status: http.StatusConflict, Message: "error", Data: map[string]interface{}{"data": "User with this email already exists"}})
			return
		} else if err != mongo.ErrNoDocuments {
			c.JSON(http.StatusInternalServerError, res.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, res.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": "Failed to hash password"}})
			return
		}
		newUser := models.User{
			ID:        primitive.NewObjectID(),
			Name:      user.Name,
			Email:     user.Email,
			Password:  string(hashedPassword),
			CPassword: string(hashedPassword),
		}

		result, err := userCollection.InsertOne(ctx, newUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, res.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusCreated, res.UserResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}})
	}
}

func Greeting() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusCreated, res.UserResponse{Status: http.StatusAccepted, Data: map[string]interface{}{"success": true}})
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var loginData models.LoginData
		if err := c.ShouldBindJSON(&loginData); err != nil {
			c.JSON(http.StatusBadRequest, res.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if validationErr := validate.Struct(&loginData); validationErr != nil {
			c.JSON(http.StatusBadRequest, res.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		existingUser := models.User{}
		err := userCollection.FindOne(ctx, bson.M{"email": loginData.Email}).Decode(&existingUser)
		if err == mongo.ErrNoDocuments {

			c.JSON(http.StatusNotFound, res.UserResponse{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": "User not found"}})
			return
		} else if err != nil {

			c.JSON(http.StatusInternalServerError, res.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(loginData.Password))
		if err != nil {

			c.JSON(http.StatusUnauthorized, res.UserResponse{Status: http.StatusUnauthorized, Message: "error", Data: map[string]interface{}{"data": "Invalid credentials"}})
			return
		}

		c.JSON(http.StatusOK, res.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": "Login successful"}})
	}
}
