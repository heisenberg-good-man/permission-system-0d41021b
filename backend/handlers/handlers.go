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

func ListCandidates(c *gin.Context) {
	keyword := c.Query("keyword")
	jobId := c.Query("jobId")
	status := c.Query("status")

	candidates := store.ListCandidates(keyword, jobId, status)
	c.JSON(http.StatusOK, gin.H{"data": candidates})
}

func GetCandidate(c *gin.Context) {
	contact := c.Query("contact")
	name := c.Query("name")

	if contact == "" || name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "联系方式和姓名不能为空"})
		return
	}

	candidate := store.GetCandidate(contact, name)
	if candidate == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "候选人不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": candidate})
}

type updateCandidateStatusRequest struct {
	Contact string                     `json:"contact" binding:"required"`
	Name    string                     `json:"name" binding:"required"`
	Status  models.ApplicationStatus `json:"status" binding:"required"`
}

func UpdateCandidateStatus(c *gin.Context) {
	var req updateCandidateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validStatuses := map[models.ApplicationStatus]bool{
		models.AppStatusPending:   true,
		models.AppStatusInterview: true,
		models.AppStatusRejected:  true,
		models.AppStatusHired:     true,
	}

	if !validStatuses[req.Status] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的候选人状态"})
		return
	}

	err := store.UpdateCandidateStatus(req.Contact, req.Name, req.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "状态更新成功"})
}

func ListNotes(c *gin.Context) {
	contact := c.Query("contact")
	name := c.Query("name")

	notes := store.ListNotes(contact, name)
	c.JSON(http.StatusOK, gin.H{"data": notes})
}

type createNoteRequest struct {
	Candidate string           `json:"candidate" binding:"required"`
	Contact   string           `json:"contact" binding:"required"`
	Type      models.NoteType `json:"type" binding:"required"`
	Content   string           `json:"content" binding:"required"`
	CreatedBy string           `json:"createdBy"`
}

func CreateNote(c *gin.Context) {
	var req createNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validTypes := map[models.NoteType]bool{
		models.NoteTypeInterview: true,
		models.NoteTypeScreen:    true,
	}

	if !validTypes[req.Type] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的备注类型"})
		return
	}

	note := &models.Note{
		Candidate: req.Candidate,
		Contact:   req.Contact,
		Type:      req.Type,
		Content:   req.Content,
		CreatedBy: req.CreatedBy,
	}

	created := store.CreateNote(note)
	c.JSON(http.StatusCreated, gin.H{"data": created})
}

func ListInterviews(c *gin.Context) {
	jobId := c.Query("jobId")
	status := c.Query("status")
	method := c.Query("method")

	interviews := store.ListInterviews(jobId, status, method)
	c.JSON(http.StatusOK, gin.H{"data": interviews})
}

func GetInterview(c *gin.Context) {
	id := c.Param("id")
	iv := store.GetInterview(id)
	if iv == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "面试安排不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": iv})
}

type createInterviewRequest struct {
	ApplicationID string                `json:"applicationId" binding:"required"`
	Round         int                   `json:"round" binding:"required,min=1"`
	InterviewTime string                `json:"interviewTime" binding:"required"`
	Method        models.InterviewMethod `json:"method" binding:"required"`
	Interviewer   string                `json:"interviewer" binding:"required"`
	Location      string                `json:"location"`
	Description   string                `json:"description"`
}

func CreateInterview(c *gin.Context) {
	var req createInterviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validMethods := map[models.InterviewMethod]bool{
		models.InterviewMethodOnsite: true,
		models.InterviewMethodOnline:  true,
		models.InterviewMethodPhone:   true,
	}
	if !validMethods[req.Method] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的面试方式"})
		return
	}

	iv := &models.Interview{
		ApplicationID: req.ApplicationID,
		Round:         req.Round,
		InterviewTime: req.InterviewTime,
		Method:        req.Method,
		Interviewer:   req.Interviewer,
		Location:      req.Location,
		Description:   req.Description,
	}

	created := store.CreateInterview(iv)
	if created == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "申请不存在，无法创建面试"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": created})
}

type rescheduleInterviewRequest struct {
	InterviewTime string `json:"interviewTime" binding:"required"`
}

func RescheduleInterview(c *gin.Context) {
	id := c.Param("id")

	var req rescheduleInterviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updated := store.RescheduleInterview(id, req.InterviewTime)
	if updated == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "面试安排不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": updated})
}

func CancelInterview(c *gin.Context) {
	id := c.Param("id")

	updated := store.CancelInterview(id)
	if updated == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "面试安排不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": updated})
}

type completeInterviewRequest struct {
	Feedback string `json:"feedback"`
	Note     string `json:"note"`
}

func CompleteInterview(c *gin.Context) {
	id := c.Param("id")

	var req completeInterviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updated := store.CompleteInterview(id, store.CompleteInterviewRequest{
		Feedback: req.Feedback,
		Note:     req.Note,
	})
	if updated == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "面试安排不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": updated})
}

type addInterviewNoteRequest struct {
	Note string `json:"note" binding:"required"`
}

func AddInterviewNote(c *gin.Context) {
	id := c.Param("id")

	var req addInterviewNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updated := store.AddInterviewNote(id, req.Note)
	if updated == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "面试安排不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": updated})
}
