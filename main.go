package main

import (
	"go_api/clients"
	"go_api/handlers"
	"go_api/middleware"
	"go_api/models"

	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	models.ConnectDatabase()
	router := gin.Default()
	router.POST("/register", handlers.Register)
	router.POST("/login", handlers.Login)
	protected := router.Group("/")
	protected.Use(middleware.JWTAuthMiddleware())
	protected.GET("/protected", func(c *gin.Context) {
		username := c.MustGet("username").(string)
		c.JSON(http.StatusOK, gin.H{"message": "Welcome " + username})
	})

	router.GET("/jobs", clients.GetJobs)
	router.GET("/jobs/:location/:mincost/:maxcost")
	router.POST("/apply", clients.ApplyToJob)
	router.PUT("/application/:id/:status", clients.UpdateApplicationStatus) // Accept or Deny
	router.GET("/jobs", clients.GetJobs)
	router.POST("/jobs", clients.CreateJobs)
	router.GET("/job/:id", clients.JobByID)

	//router.GET("/books", clients.getBooks)
	//router.GET("/books/:id", bookById)
	//router.POST("/books", createBooks)
	//router.PATCH("/checkout", checkoutBook)
	//router.PATCH("/return", returnBook)
	//router.DELETE("/delete/:id", deleteBook)
	//router.GET("/search", searchBooks)

	err := router.Run("localhost:8080")
	if err != nil {
		return
	}
}
