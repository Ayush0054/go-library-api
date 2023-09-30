package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)
type book struct {
	ID string     `json:"id"`
	Title string	`json:"title"` 
	Author string    `json:"author"`
	Quantity int     `json:"quantity"`
}
var client *mongo.Client
var collection *mongo.Collection




func main(){

	clientOptions := options.Client().ApplyURI("db uRI")
    client, err := mongo.Connect(context.Background(), clientOptions)
    if err != nil {
        log.Fatal(err)
    }

   
    err = client.Ping(context.Background(), nil)
    if err != nil {
        log.Fatal(err)
    }

    collection = client.Database("library").Collection("books")
	router := gin.Default()
     router.GET("/books", getBooks)
	 router.GET("/books/:id", bookByID)
	 router.POST("/books", createBook)
	 router.DELETE("/remove", removeBook)
	 router.PATCH("/checkout", checkoutBook)
	 router.PATCH("/return", returnBook)
	 router.Run("localhost:8080")
}