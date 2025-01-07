package handlers

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"go_api/models"
	"go_api/utils"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	// Normalize the username for case-insensitive comparison
	normalizedUsername := strings.ToLower(req.Username)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	_, err := models.DB.Exec(context.Background(), "INSERT INTO users2 (username, password_hash) VALUES ($1, $2)", normalizedUsername, string(hashedPassword))
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			c.JSON(http.StatusConflict, gin.H{"message": "Username already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "An error occurred"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func Login(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	var passwordHash string
	err := models.DB.QueryRow(context.Background(), "SELECT password_hash FROM users2 WHERE username=$1", req.Username).Scan(&passwordHash)
	if err == pgx.ErrNoRows || bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid username or password"})
		return
	}

	token, _ := utils.GenerateJWT(req.Username)
	c.JSON(http.StatusOK, gin.H{"token": token})
}
