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
	ID            int       `json:"-" orm:"auto;pk"`
	Title         string    `json:"title" orm:"size(200)"`
	Content       string    `json:"content" orm:"type(text)"`
	LikesCount    int       `json:"likes_count" orm:"default(0)"`
	CommentsCount int       `json:"comments_count" orm:"default(0)"`
	Author        *User     `json:"-" orm:"rel(fk)"` // Foreign key to User
	CreatedAt     time.Time `json:"-" orm:"auto_now_add;type(datetime)"`
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
	postId, err := strconv.Atoi(c.Param("postId")) // Get Id from URL
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Id"})
		c.Abort()
		return
	}

	// _, err := GetPostById(id)
	o := orm.NewOrm()
	if numRows, err := o.QueryTable("post").Filter("id", postId).Delete(); err != nil || numRows == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Post not found"})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted the Post"})
}

func LoadBlogWithId(c *gin.Context) {
	userID, _ := c.Get("userID")

	postId, err := strconv.Atoi(c.Param("postId")) // Get post Id from URL
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Id"})
		c.Abort()
		return
	}

	o := orm.NewOrm()
	post := Post{ID: postId}

	// Try to fetch the blog post
	err = o.Read(&post)
	if err == orm.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Blog post not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	// Check for liked status
	liked := false
	if userID != nil {
		count, err := o.QueryTable("like").Filter("user_id", userID.(int)).Filter("post_id", postId).Count()
		if err != nil {
			c.JSON(http.StatusNoContent, gin.H{"error": "DataBase Error"})
			return
		}
		liked = count > 0
	}

	// Return the blog post as JSON
	c.JSON(http.StatusOK, gin.H{
		"data":  post,
		"liked": liked,
	})
}

func LoadBlogs(c *gin.Context) {
	o := orm.NewOrm()
	var posts []Post
	// getting info from frondend about how much content to load and from where
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "5"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	// fetch from the posts dataset
	count, err := o.QueryTable("post").Count()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count"})
	}
	_, err = o.QueryTable("post").RelatedSel().Limit(limit).Offset(offset).All(&posts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Fetch the data"})
	}

	hasMore := (offset + limit) < int(count)

	var response []map[string]any
	for _, post := range posts {
		response = append(response, map[string]any{
			"id":             post.ID,
			"title":          post.Title,
			"content":        post.Content,
			"likes_count":    post.LikesCount,
			"comments_count": post.CommentsCount,
			"author":         post.Author.Name,
			"created_at":     post.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    response,
		"hasMore": hasMore,
	})
}
