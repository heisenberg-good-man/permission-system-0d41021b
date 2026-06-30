package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"recruitment-platform/models"
	"recruitment-platform/store"
)

type createJobRequest struct {
	Title        string          `json:"title" binding:"required"`
	Department   string          `json:"department" binding:"required"`
	Location     string          `json:"location" binding:"required"`
	SalaryMin    int             `json:"salaryMin" binding:"required,min=0"`
	SalaryMax    int             `json:"salaryMax" binding:"required,min=0"`
	Status       models.JobStatus `json:"status" binding:"required"`
	Headcount    int             `json:"headcount" binding:"required,min=1"`
	Description  string          `json:"description"`
	Requirements string          `json:"requirements"`
}

func ListJobs(c *gin.Context) {
	keyword := c.Query("keyword")
	location := c.Query("location")
	status := c.Query("status")

	jobs := store.ListJobs(keyword, location, status)
	c.JSON(http.StatusOK, gin.H{"data": jobs})
}

func GetJob(c *gin.Context) {
	id := c.Param("id")
	job := store.GetJob(id)
	if job == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "职位不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": job})
}

func CreateJob(c *gin.Context) {
	var req createJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	job := &models.Job{
		Title:        req.Title,
		Department:   req.Department,
		Location:     req.Location,
		SalaryMin:    req.SalaryMin,
		SalaryMax:    req.SalaryMax,
		Status:       req.Status,
		Headcount:    req.Headcount,
		Description:  req.Description,
		Requirements: req.Requirements,
	}

	created := store.CreateJob(job)
	c.JSON(http.StatusCreated, gin.H{"data": created})
}

func UpdateJob(c *gin.Context) {
	id := c.Param("id")

	var req createJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	job := &models.Job{
		Title:        req.Title,
		Department:   req.Department,
		Location:     req.Location,
		SalaryMin:    req.SalaryMin,
		SalaryMax:    req.SalaryMax,
		Status:       req.Status,
		Headcount:    req.Headcount,
		Description:  req.Description,
		Requirements: req.Requirements,
	}

	updated := store.UpdateJob(id, job)
	if updated == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "职位不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": updated})
}

type createApplicationRequest struct {
	JobID  string       `json:"jobId" binding:"required"`
	Resume models.Resume `json:"resume" binding:"required"`
}

func ListApplications(c *gin.Context) {
	status := c.Query("status")
	jobId := c.Query("jobId")

	apps := store.ListApplications(status, jobId)
	c.JSON(http.StatusOK, gin.H{"data": apps})
}

func GetApplication(c *gin.Context) {
	id := c.Param("id")
	app := store.GetApplication(id)
	if app == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "申请不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": app})
}

func CreateApplication(c *gin.Context) {
	var req createApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Resume.CandidateName == "" || req.Resume.Contact == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "候选人姓名和联系方式不能为空"})
		return
	}

	app := &models.Application{
		Resume: req.Resume,
	}

	created := store.CreateApplication(app, req.JobID)
	if created == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "职位不存在"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": created})
}

type updateStatusRequest struct {
	Status models.ApplicationStatus `json:"status" binding:"required"`
}

func UpdateApplicationStatus(c *gin.Context) {
	id := c.Param("id")

	var req updateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validStatuses := map[models.ApplicationStatus]bool{
		models.AppStatusSubmitted: true,
		models.AppStatusPending:   true,
		models.AppStatusInterview: true,
		models.AppStatusRejected:  true,
		models.AppStatusHired:     true,
	}

	if !validStatuses[req.Status] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的申请状态"})
		return
	}

	updated := store.UpdateApplicationStatus(id, req.Status)
	if updated == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "申请不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": updated})
}

type sendMessageRequest struct {
	ApplicationID string            `json:"applicationId" binding:"required"`
	SenderName    string            `json:"senderName" binding:"required"`
	SenderRole    models.SenderRole `json:"senderRole" binding:"required"`
	Content       string            `json:"content" binding:"required"`
}

func ListMessages(c *gin.Context) {
	applicationId := c.Param("applicationId")
	msgs := store.ListMessages(applicationId)
	c.JSON(http.StatusOK, gin.H{"data": msgs})
}

func SendMessage(c *gin.Context) {
	var req sendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.SenderRole != models.RoleRecruiter && req.SenderRole != models.RoleCandidate {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的发送者角色"})
		return
	}

	msg := &models.Message{
		ApplicationID: req.ApplicationID,
		SenderName:    req.SenderName,
		SenderRole:    req.SenderRole,
		Content:       req.Content,
	}

	created := store.SendMessage(msg)
	c.JSON(http.StatusCreated, gin.H{"data": created})
}

func GetStats(c *gin.Context) {
	stats := store.GetStats()
	c.JSON(http.StatusOK, gin.H{"data": stats})
}
