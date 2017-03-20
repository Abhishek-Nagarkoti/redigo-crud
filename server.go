package main

import (
	"github.com/Abhishek-Nagarkoti/redigo-crud/handlers"
	"github.com/joho/godotenv"
	"gopkg.in/gin-gonic/gin.v1"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	h := handlers.Handler{}
	h.Connect()

	router := gin.Default()
	router.Use(func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE")
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, X-HTTP-Method-Override,Authorization, Content-Type, Accept")
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		if ctx.Request.Method == "OPTIONS" {
			ctx.JSON(200, gin.H{"All": "Good"})
		} else {
			ctx.Next()
		}

	})

	api := router.Group("/")
	{
		api.POST("/", h.Set)
		api.GET("/", h.Get)
		api.PUT("/:id", h.Update)
		api.DELETE("/:id", h.Delete)
		api.GET("/:name", h.Find)
	}

	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	router.Run(":8080")

}
