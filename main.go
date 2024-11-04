package main

import (
	"rest-api-practice/database"
	"rest-api-practice/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// fmt.Println("Whatsup bitches")
	database.InitDB()
	server := gin.Default()

	routes.RegisterRoutes(server)

	server.Run(":8080") // localhost:8080
}
