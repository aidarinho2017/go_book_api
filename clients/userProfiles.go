package clients

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"go_api/models"
	"net/http"
)

type User struct {
	ID         int      `json:"id"`
	Name       string   `json:"name"`
	Email      string   `json:"email"`
	Phone      string   `json:"phone"` // Renamed from PhoneNumber to match DB
	PassportID string   `json:"passport_id"`
	Reviews    []Review `json:"reviews"`
	AvgRating  float64  `json:"avg_rating"`
}

type Review struct {
	ReviewerID string  `json:"reviewer_id"`
	Comment    string  `json:"comment"`
	Rating     float64 `json:"rating"`
}

func AddReview(c *gin.Context) {
	var req struct {
		ReviewedID int    `json:"reviewed_id"`
		JobID      int    `json:"job_id"`
		Rating     int    `json:"rating"`
		Comment    string `json:"comment"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	// Check if the job was completed by both users
	var count int
	err := models.DB.QueryRow(context.Background(), `
        SELECT COUNT(*) FROM job_applications
        WHERE id = $1 AND status = 'completed'
        AND (applicant_id = $2 OR owner_id = $2)`,
		req.JobID, req.ReviewedID).Scan(&count)

	if err != nil || count == 0 {
		c.JSON(http.StatusForbidden, gin.H{"message": "Review not allowed"})
		return
	}

	// Insert review
	_, err = models.DB.Exec(context.Background(), `
        INSERT INTO reviews (reviewer_id, reviewed_id, job_id, rating, comment)
        VALUES ($1, $2, $3, $4, $5)`,
		c.GetInt("userID"), req.ReviewedID, req.JobID, req.Rating, req.Comment)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error saving review"})
		return
	}

	// Update avg_rating
	_, err = models.DB.Exec(context.Background(), `
        UPDATE user_profiles SET avg_rating = (
            SELECT COALESCE(AVG(rating), 0) FROM reviews WHERE reviewed_id = $1
        ) WHERE user_id = $1`, req.ReviewedID)

	c.JSON(http.StatusCreated, gin.H{"message": "Review submitted successfully"})
}

func GetUserProfile(c *gin.Context) {
	userID := c.Param("id")

	var profile struct {
		ID        int     `json:"id"`
		Username  string  `json:"username"`
		Email     string  `json:"email"`
		Phone     string  `json:"phone"`
		AvgRating float64 `json:"avg_rating"`
		Reviews   []Review
	}

	err := models.DB.QueryRow(context.Background(), `
        SELECT u.id, u.username, u.email, u.phone, up.avg_rating
        FROM users2 u
        LEFT JOIN user_profiles up ON u.id = up.user_id
        WHERE u.id = $1`, userID).
		Scan(&profile.ID, &profile.Username, &profile.Email, &profile.Phone, &profile.AvgRating)

	if err == pgx.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	rows, _ := models.DB.Query(context.Background(), `
        SELECT r.rating, r.comment FROM reviews r WHERE r.reviewed_id = $1`, userID)
	defer rows.Close()

	for rows.Next() {
		var review Review
		rows.Scan(&review.Rating, &review.Comment)
		profile.Reviews = append(profile.Reviews, review)
	}

	c.JSON(http.StatusOK, profile)
}
