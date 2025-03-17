package models

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       int    `json:"-" orm:"auto"`
	Name     string `json:"name" orm:"size(30)"`
	Email    string `json:"email" orm:"size(100);unique"`
	Password string `json:"password" orm:"size(100)"`
}

// Hash the password using bcrypt to protect user privacy
func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func ViewUser(c *gin.Context){
	c.JSON(http.StatusTemporaryRedirect, gin.H{"message": "This feature is currently Under Development :)"})
}

func UpdateUser(c *gin.Context){
	c.JSON(http.StatusTemporaryRedirect, gin.H{"message": "This feature is currently Under Development :)"})
}
