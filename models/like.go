package models

import (
	"net/http"
	"strconv"

	"github.com/beego/beego/v2/client/orm"
	"github.com/gin-gonic/gin"
)

type Like struct {
	Id int `orm:"auto;pk"`
	User *User `orm:"rel(fk)"`
	Post *Post `orm:"rel(fk);on_delete(cascade)"`
}

func LikePost(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}
	postID, err := strconv.Atoi(c.Param("post_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}
	var like Like
	like.User = (&(User{Id: userID.(int)}))
	like.Post = (&(Post{ID: postID}))
	o := orm.NewOrm()

	if _, err := o.Insert(&like); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "couldn't upload like"})
		c.Abort()
		return
	}
	// Update like count
	_, err = o.QueryTable("post").Filter("id", postID).Update(orm.Params{"likes_count": orm.ColValue(orm.ColAdd, 1)})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update like count"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Liked Successfully"})
}

func UnLikePost(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}

	postID, err := strconv.Atoi(c.Param("post_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	o := orm.NewOrm()
	
	_, err = o.QueryTable("like").Filter("user_id", userID.(int)).Filter("post_id", postID).Delete()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't Delete Like"})
	}

	_, err = o.QueryTable("post").Filter("id", postID).Update(orm.Params{"likes_count": orm.ColValue(orm.ColMinus, 1)})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update like count"})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Like deleted successfully"})
}
