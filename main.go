package main

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type book struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Quantity int    `json:"quantity"`
}

var books = []book{
	{ID: "1", Title: "In Search of Lost Time", Author: "Marcel Proust", Quantity: 2},
	{ID: "2", Title: "The Great Gatsby", Author: "F. Scott Fitzgerald", Quantity: 5},
	{ID: "3", Title: "War and Peace", Author: "Leo Tolstoy", Quantity: 6},
}

func getBooks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, books)
}

func createBooks(c *gin.Context) {
	var newBooks []book
	if err := c.BindJSON(&newBooks); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	books = append(books, newBooks...)
	c.IndentedJSON(http.StatusCreated, newBooks)
}

func bookById(c *gin.Context) {
	id := c.Param("id")
	book, err := getBookById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."})
		return
	}
	c.IndentedJSON(http.StatusOK, book)
}

func getBookById(id string) (*book, error) {
	for i, b := range books {
		if b.ID == id {
			return &books[i], nil
		}
	}
	return nil, errors.New("Book not found")
}

func checkoutBook(c *gin.Context) {
	id, ok := c.GetQuery("id")
	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid Book ID."})
		return
	}
	book, err := getBookById(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."})
		return
	}
	if book.Quantity <= 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid Book Quantity."})
		return
	}
	book.Quantity--
	c.IndentedJSON(http.StatusOK, book)
}

func returnBook(c *gin.Context) {
	id, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id query parameter."})
		return
	}

	book, err := getBookById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."})
		return
	}

	book.Quantity += 1
	c.IndentedJSON(http.StatusOK, book)
}

func deleteBook(c *gin.Context) {
	id := c.Param("id")
	for i, b := range books {
		if b.ID == id {
			books = append(books[:i], books[i+1:]...) // Remove the book
			c.IndentedJSON(http.StatusOK, gin.H{"message": "Book deleted."})
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."})
}

func searchBooks(c *gin.Context) {
	query := c.Query("q")
	var results []book

	for _, b := range books {
		if containsIgnoreCase(b.Title, query) || containsIgnoreCase(b.Author, query) {
			results = append(results, b)
		}
	}

	if len(results) == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "No books match your query."})
		return
	}

	c.IndentedJSON(http.StatusOK, results)
}

func containsIgnoreCase(source, target string) bool {
	return strings.Contains(strings.ToLower(source), strings.ToLower(target))
}

func main() {
	router := gin.Default()
	router.GET("/books", getBooks)
	router.GET("/books/:id", bookById)
	router.POST("/books", createBooks)
	router.PATCH("/checkout", checkoutBook)
	router.PATCH("/return", returnBook)
	router.DELETE("/delete/:id", deleteBook)
	router.GET("/search", searchBooks)

	err := router.Run("localhost:8080")
	if err != nil {
		return
	}
}
