package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	mongoDbStorage := NewMongoDbStorage()
	handler := NewHandler(mongoDbStorage)

	router := gin.Default()

	router.Use(cors.Default())

	router.GET("/api/v1/persons/:id", handler.GetPerson)
	router.POST("/api/v1/persons", handler.CreatePerson)
	router.GET("/api/v1/persons", handler.GetPersons)
	router.PATCH("/api/v1/persons/:id", handler.UpdatePerson)
	router.DELETE("/api/v1/persons/:id", handler.DeletePerson)

	router.Run()
}
