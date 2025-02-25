package models

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/gin-gonic/gin"
)

type Post struct {
	ID        int       `json:"-" orm:"auto;pk"`
	Title     string    `json:"title" orm:"size(200)"`
	Content   string    `json:"content" orm:"type(text)"`
	Author    *User     `json:"-" orm:"rel(fk)"` // Foreign key to User
	CreatedAt time.Time `json:"-" orm:"auto_now_add;type(datetime)"`
}

func UploadBlog(c *gin.Context) {
	/*
		here we do not have a method to stop posting duplicate blogs i.e user could post multiple posts with same content
		--> have to change it
	*/

	// Retrieve `userID` from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Title/Content/Author Missing"})
		c.Abort()
		return
	}
	var post Post
	post.Author = (&(User{Id: userID.(int)}))
	if err := json.Unmarshal(body, &post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Some thing is missing"})
		c.Abort()
		return
	}
	o := orm.NewOrm()

	if _, err := o.Insert(&post); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "couldn't upload post"})
		c.Abort()
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Uploaded Successfully"})
}

func UpdateBlog(c *gin.Context) {

}

// func GetPostById(id int) (*Post, error) {
// 	o := orm.NewOrm()
// 	post := Post{}

// 	err := o.QueryTable("post").Filter("id", id).One(&post)
// 	if err == orm.ErrNoRows {
// 		return nil, errors.New("post not found")
// 	}
// 	return &post, err
// }

func DeleteBlog(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("Id")) // Get Id from URL
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Id"})
		c.Abort()
		return
	}

	// _, err := GetPostById(id)
	o := orm.NewOrm()
	if numRows, err := o.QueryTable("post").Filter("id", id).Delete(); err != nil || numRows == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Post not found"})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted the Post"})
}
