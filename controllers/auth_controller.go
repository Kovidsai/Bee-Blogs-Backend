package controllers

import (
	"MINIPROJECT/models"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func GenerateJWT(UserId int) (string, error) {
	// when the token is generated it will be there untill it is expired
	// TODO: have to implement mechanism to ensure token is expired once he logs out before that we have to implement logout :)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": UserId,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token valid for 24 hours
	})

	return token.SignedString(jwtSecret)

}

func Register(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Input"})
		return
	}

	var user models.User
	if err := json.Unmarshal(body, &user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Input"})
		return
	}
	// Hash Password
	if err := user.HashPassword(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}
	// Insert into database
	o := orm.NewOrm()
	if _, err := o.Insert(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully!"})
}

func Login(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Fetch user from database
	var input models.User
	if err := json.Unmarshal(body, &input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Input"})
		return
	}
	o := orm.NewOrm()
	user := models.User{Email: input.Email}

	if err := o.QueryTable("user").Filter("Email", input.Email).One(&user); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Check Password
	if !user.CheckPassword(input.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT Token
	token, err := GenerateJWT(user.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})

}
