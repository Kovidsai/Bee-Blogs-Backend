package routes

import (
	"MINIPROJECT/controllers"
	"MINIPROJECT/middleware"
	"MINIPROJECT/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetUpRoutes(r *gin.Engine) {
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
		protected.GET("/profile", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Secured profile"})
		})
		protected.POST("/uploadblog", models.UploadBlog)
		protected.DELETE("/deleteblog/:Id", models.DeleteBlog)
		protected.PUT("/updateblog", models.UpdateBlog)

	}
}
