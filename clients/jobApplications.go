package clients

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type JobApplication struct {
	ID        string `json:"id"`
	Applicant string `json:"author"`
	Cost      int    `json:"cost"`
	Location  string `json:"location"`
	Status    string `json:"status"`
}

var jobApplications []JobApplication = []JobApplication{}

func ApplyToJob(c *gin.Context) {
	var newJobApplication JobApplication

	if err := c.BindJSON(&newJobApplication); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"Message": "Wrong Input"})
	}

	newJobApplication.Status = "pending"
	jobApplications = append(jobApplications, newJobApplication)
	c.IndentedJSON(http.StatusCreated, gin.H{"Message": "Job application created"})

}

func UpdateApplicationStatus(c *gin.Context) {
	jobId := c.Param("id")
	jobStatus := c.Param("status")

	for i, app := range jobApplications {
		if app.ID == jobId {
			if jobStatus == "accept" {
				jobApplications[i].Status = "accepted"
			} else {
				jobApplications[i].Status = "denied"
			}
			c.JSON(http.StatusOK, "Application updated")
			return
		}
	}
	c.JSON(http.StatusBadRequest, "bad request")
}

func UpdateJobApplication(c *gin.Context) {
	jobAppId := c.Param("id")
	var updateData map[string]interface{}

	if err := c.BindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	for i, app := range jobApplications {
		if app.ID == jobAppId {
			if applicant, exists := updateData["applicant"].(string); exists {
				jobApplications[i].Applicant = applicant
			}
			if cost, exists := updateData["cost"].(float64); exists {
				jobApplications[i].Cost = int(cost)
			}
			if location, exists := updateData["location"].(string); exists {
				jobApplications[i].Location = location
			}
			if status, exists := updateData["status"].(string); exists {
				jobApplications[i].Status = status
			}

			c.JSON(http.StatusOK, gin.H{"message": "Job application updated", "application": jobApplications[i]})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"message": "Job application not found"})
}

func GetJobApplicationById(id string) (*JobApplication, error) {
	for i, _ := range jobApplications {
		if jobApplications[i].ID == id {
			return &jobApplications[i], nil
		}
	}
	return nil, errors.New("job application not found")

}
