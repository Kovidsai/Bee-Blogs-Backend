package main

import (
	"MINIPROJECT/config"
	"MINIPROJECT/routes"

	"github.com/gin-gonic/gin"
)

func main() {

	config.ConnectDB()

	r := gin.Default()
	routes.SetUpRoutes(r)
	r.Run(":8080")
}
