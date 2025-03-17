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

type Comment struct {
	ID        int       `json:"-" orm:"auto;pk"`
	Content   string    `json:"content" orm:"type(text)"`
	Post      *Post     `json:"-" orm:"rel(fk);on_delete(cascade)"` // Foreign key to Post, this will make sure All Comments are deleted when a post is deleted
	Author    *User     `json:"-" orm:"rel(fk)"`                    // Foreign key to User
	CreatedAt time.Time `json:"-" orm:"auto_now_add;type(datetime)"`
}

func AddComment(c *gin.Context) {
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
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Title/Content/Author Missing"})
		c.Abort()
		return
	}
	var comment Comment
	comment.Author = (&(User{Id: userID.(int)}))
	comment.Post = (&(Post{ID: postID}))
	if err := json.Unmarshal(body, &comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Some thing is missing"})
		c.Abort()
		return
	}
	o := orm.NewOrm()

	if _, err := o.Insert(&comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "couldn't upload post"})
		c.Abort()
		return
	}
	// Update comments count
	_, err = o.QueryTable("post").Filter("id", postID).Update(orm.Params{"comments_count": orm.ColValue(orm.ColAdd, 1)})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update comment count"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Uploaded Successfully"})
}

func ViewLatestComments(c *gin.Context) {
	postId, err := strconv.Atoi(c.Param("post_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}
	o := orm.NewOrm()
	var comments []Comment
	_, err = o.QueryTable("comment").Filter("post_id", postId).RelatedSel().Limit(4).OrderBy("-CreatedAt").All(&comments)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Fetch the data"})
	}

	var response []map[string]any
	for _, comment := range comments {
		response = append(response, map[string]any{
			"id":         comment.ID,
			"content":    comment.Content,
			"author":     comment.Author.Name,
			"created_at": comment.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})

}

func ViewAllComments(c *gin.Context) {

	c.JSON(http.StatusPartialContent, gin.H{"message": "This Will be Available Soon"})

	/// change this later
	// o := orm.NewOrm()
	// var comments []Comment
	// // getting info from frondend about how much content to load and from where
	// limit, _ := strconv.Atoi(c.DefaultQuery("limit", "5"))
	// offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	// // fetch from the posts dataset
	// count, err := o.QueryTable("comment").Count()
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count"})
	// }
	// _, err = o.QueryTable("post").RelatedSel().Limit(limit).Offset(offset).All(&posts)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Fetch the data"})
	// }

	// hasMore := (offset + limit) < int(count)

	// var response []map[string]any
	// for _, post := range posts {
	// 	response = append(response, map[string]any{
	// 		"id":             post.ID,
	// 		"title":          post.Title,
	// 		"content":        post.Content,
	// 		"likes_count":    post.LikesCount,
	// 		"comments_count": post.CommentsCount,
	// 		"author":         post.Author.Name,
	// 		"created_at":     post.CreatedAt,
	// 	})
	// }

	// c.JSON(http.StatusOK, gin.H{
	// 	"data":    response,
	// 	"hasMore": hasMore,
	// })
}

func DeleteComment(c *gin.Context) {
	commentId, err := strconv.Atoi(c.Param("comment_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Comment Id"})
		return
	}

	postID, err := strconv.Atoi(c.Param("post_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	o := orm.NewOrm()
	_, err = o.QueryTable("comment").Filter("id", commentId).Delete()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't Delete Comment"})
	}

	_, err = o.QueryTable("post").Filter("id", postID).Update(orm.Params{"comment_count": orm.ColValue(orm.ColMinus, 1)})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update comment count"})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}
