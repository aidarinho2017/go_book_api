package clients

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type job struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Cost     int    `json:"cost"`
	Location string `json:"location"`
	Status   string `json:"status"`
}

type jobApplication struct {
	ID        string `json:"id"`
	Applicant string `json:"author"`
	Cost      int    `json:"cost"`
	Location  string `json:"location"`
	Status    string `json:"status"`
}

var jobs = []job{
	{ID: "1", Title: "Babysitting", Author: "JohnDoe", Cost: 50, Location: "Almaty", Status: "Open"},
	{ID: "2", Title: "Car Wash", Author: "JaneDoe", Cost: 30, Location: "Astana", Status: "Open"},
	{ID: "3", Title: "House Cleaning", Author: "MikeSmith", Cost: 40, Location: "Shymkent", Status: "Open"},
}

var jobApplications []jobApplication = []jobApplication{}

func GetJobs(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, jobs)
}

func getJobsFiltered(c *gin.Context) {
	location := c.Query("location")
	minCost := c.Query("min_cost")
	maxCost := c.Query("max_cost")
	var filteredJobs []job
	minimumCost, err := strconv.Atoi(minCost)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "wrong Input"})
		return
	}
	maximumCost, err := strconv.Atoi(maxCost)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "wrong Input"})
		return
	}

	for _, j := range jobs {

		if (location == "" || j.Location == location) &&
			(minCost == "" || j.Cost >= minimumCost) &&
			(maxCost == "" || j.Cost <= maximumCost) {
			filteredJobs = append(filteredJobs, j)
		}
	}
	c.JSON(http.StatusOK, filteredJobs)
}

func CreateJobs(c *gin.Context) {
	var newJobs []job
	if err := c.BindJSON(&newJobs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	jobs = append(jobs, newJobs...)
	c.IndentedJSON(http.StatusCreated, newJobs)
}

func GetJobByID(id string) (*job, error) {
	for i, b := range jobs {
		if b.ID == id {
			return &jobs[i], nil
		}
	}
	return nil, errors.New("No book found.")
}

func JobByID(c *gin.Context) {
	c.Param("id")
	job, err := GetJobByID("id")
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"Message": "Not found"})
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": job})
}

func ApplyToJob(c *gin.Context) {
	var newJobApplication jobApplication

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

func UpdateJob(c *gin.Context) {
	jobId := c.Param("id")
	var updateData map[string]interface{}

	if err := c.BindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	for i, job := range jobs {
		if job.ID == jobId {
			// Update only provided fields
			if title, exists := updateData["title"].(string); exists {
				jobs[i].Title = title
			}
			if author, exists := updateData["author"].(string); exists {
				jobs[i].Author = author
			}
			if cost, exists := updateData["cost"].(float64); exists { // JSON numbers default to float64
				jobs[i].Cost = int(cost)
			}
			if location, exists := updateData["location"].(string); exists {
				jobs[i].Location = location
			}
			if status, exists := updateData["status"].(string); exists {
				jobs[i].Status = status
			}

			c.JSON(http.StatusOK, gin.H{"message": "Job updated", "job": jobs[i]})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"message": "Job not found"})
}
