package clients

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
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
