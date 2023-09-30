package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func getBooks(c *gin.Context){

	var bookList []book

    cursor, err := collection.Find(context.Background(), bson.M{})
    if err != nil {
        c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Error fetching books"})
        return
    }
    defer cursor.Close(context.Background())

    for cursor.Next(context.Background()) {
        var b book
        if err := cursor.Decode(&b); err != nil {
            c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Error decoding books"})
            return
        }
        bookList = append(bookList, b)
    }

    c.IndentedJSON(http.StatusOK, bookList)
}

func bookByID(c *gin.Context) {
    id := c.Param("id")
    var b book

    err := collection.FindOne(context.Background(), bson.M{"id": id}).Decode(&b)
    if err != nil {
        c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found"})
        return
    }

    c.IndentedJSON(http.StatusOK, b)
}
func createBook( c *gin.Context){

var newBook book
if err := c.BindJSON(&newBook); err != nil {
	c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid JSON"})
	return
}

_, err := collection.InsertOne(context.Background(), newBook)
if err != nil {
	c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Error creating book"})
	return
}

c.IndentedJSON(http.StatusCreated, newBook)
}
func removeBook(c *gin.Context){
	id,ok :=c.GetQuery("id")
	if(!ok){
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "missing id query parameters"})
		return
	}
	_, err := collection.DeleteOne(context.Background(), bson.M{"id": id})
	if err != nil{
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "book not found"})
		return
	}
	
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Book removed from library"})
}
func checkoutBook(c *gin.Context) {
    id := c.Query("id") 

    filter := bson.M{"id": id}
    update := bson.M{"$inc": bson.M{"quantity": -1}}

    result, err := collection.UpdateOne(context.Background(), filter, update)
    if err != nil {
        c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found or error updating quantity"})
        return
    }

    if result.ModifiedCount == 0 {
        c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found or quantity already at zero"})
        return
    }

    c.IndentedJSON(http.StatusOK, gin.H{"message": "book sold"})
}

func returnBook(c *gin.Context) {
    id := c.Query("id") 


    filter := bson.M{"id": id}
    update := bson.M{"$inc": bson.M{"quantity": 1}} 

    _, err := collection.UpdateOne(context.Background(), filter, update)
    if err != nil {
        c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found or error updating quantity"})
        return
    }

  
    var updatedBook book
    err = collection.FindOne(context.Background(), filter).Decode(&updatedBook)
    if err != nil {
        c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found or error fetching updated book"})
        return
    }

    c.IndentedJSON(http.StatusOK, updatedBook)
}
