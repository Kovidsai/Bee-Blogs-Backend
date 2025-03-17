package routes

import (
	"MINIPROJECT/controllers"
	"MINIPROJECT/middleware"
	"MINIPROJECT/models"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetUpRoutes(r *gin.Engine) {

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Frontend origin (replace with actual frontend URL)
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge: 12*time.Hour, // this cache preflight response
		/* 
			browser sends a preflight(options) request to check whether 
			our backend supports that port or not. this result will be
			cached with MaxAge to stop repeated requests from brower 
		 */
	}))

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to simple blog API"})
	})

	authRoutes := r.Group("/auth")
	{
		authRoutes.POST("/register", controllers.Register)
		authRoutes.POST("/login", controllers.Login)
	}
	protected := r.Group("/api")
	protected.Use(middleware.JWTAuthMiddleware())
	{
		//User-section
		protected.GET("/profile", models.ViewUser)
		protected.PUT("/update-profile", models.UpdateUser) // only user_name can be updated

		//Blogs-section
		protected.GET("/blogs", models.LoadBlogs)
		protected.POST("/uploadblog", models.UploadBlog)
		protected.DELETE("/deleteblog/:postId", models.DeleteBlog)
		protected.GET("/blog/:postId", models.LoadBlogWithId)
		protected.PUT("/updateblog", models.UpdateBlog)
		//comments-section
		protected.POST("/:post_id/comment", models.AddComment)
		protected.GET("/latestcomments/:post_id", models.ViewLatestComments)
		protected.GET("/allcomment/post_id", models.ViewAllComments)
		protected.DELETE("/delete-comment/:post_id/:comment_id", models.DeleteComment)
		//likes-section
		protected.POST("/like/:post_id", models.LikePost)
		protected.DELETE("/unlike/:post_id", models.UnLikePost)

	}
}
