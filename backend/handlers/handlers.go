package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"recruitment-platform/models"
	"recruitment-platform/store"
)

type createJobRequest struct {
	Title        string           `json:"title" binding:"required"`
	Department   string           `json:"department" binding:"required"`
	Location     string           `json:"location" binding:"required"`
	SalaryMin    int              `json:"salaryMin" binding:"required,min=0"`
	SalaryMax    int              `json:"salaryMax" binding:"required,min=0"`
	Status       models.JobStatus `json:"status" binding:"required"`
	Headcount    int              `json:"headcount" binding:"required,min=1"`
	Description  string           `json:"description"`
	Requirements string           `json:"requirements"`
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
	JobID  string        `json:"jobId" binding:"required"`
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
	Contact string                   `json:"contact" binding:"required"`
	Name    string                   `json:"name" binding:"required"`
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
	Candidate string          `json:"candidate" binding:"required"`
	Contact   string          `json:"contact" binding:"required"`
	Type      models.NoteType `json:"type" binding:"required"`
	Content   string          `json:"content" binding:"required"`
	CreatedBy string          `json:"createdBy"`
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
	ApplicationID string                 `json:"applicationId" binding:"required"`
	Round         int                    `json:"round" binding:"required,min=1"`
	InterviewTime string                 `json:"interviewTime" binding:"required"`
	Method        models.InterviewMethod `json:"method" binding:"required"`
	Interviewer   string                 `json:"interviewer" binding:"required"`
	Location      string                 `json:"location"`
	Description   string                 `json:"description"`
}

func CreateInterview(c *gin.Context) {
	var req createInterviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validMethods := map[models.InterviewMethod]bool{
		models.InterviewMethodOnsite: true,
		models.InterviewMethodOnline: true,
		models.InterviewMethodPhone:  true,
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
	Feedback   string `json:"feedback"`
	Note       string `json:"note"`
	Conclusion string `json:"conclusion"`
	Rating     int    `json:"rating"`
	Strengths  string `json:"strengths"`
	Risks      string `json:"risks"`
	NextSteps  string `json:"nextSteps"`
}

func CompleteInterview(c *gin.Context) {
	id := c.Param("id")

	var req completeInterviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updated := store.CompleteInterview(id, store.CompleteInterviewRequest{
		Feedback:   req.Feedback,
		Note:       req.Note,
		Conclusion: req.Conclusion,
		Rating:     req.Rating,
		Strengths:  req.Strengths,
		Risks:      req.Risks,
		NextSteps:  req.NextSteps,
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

func ListOffers(c *gin.Context) {
	jobId := c.Query("jobId")
	status := c.Query("status")
	offers := store.ListOffers(jobId, status)
	c.JSON(http.StatusOK, gin.H{"data": offers})
}

func GetOffer(c *gin.Context) {
	id := c.Param("id")
	offer := store.GetOffer(id)
	if offer == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Offer 不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": offer})
}

type createOfferRequest struct {
	ApplicationID string `json:"applicationId" binding:"required"`
	SalaryMin     int    `json:"salaryMin" binding:"required,min=0"`
	SalaryMax     int    `json:"salaryMax" binding:"required,min=0"`
	StartDate     string `json:"startDate" binding:"required"`
	Owner         string `json:"owner"`
	Description   string `json:"description"`
}

func CreateOffer(c *gin.Context) {
	var req createOfferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.SalaryMax < req.SalaryMin {
		c.JSON(http.StatusBadRequest, gin.H{"error": "最高薪资不能低于最低薪资"})
		return
	}
	offer := store.CreateOffer(&models.Offer{
		ApplicationID: req.ApplicationID,
		SalaryMin:     req.SalaryMin,
		SalaryMax:     req.SalaryMax,
		StartDate:     req.StartDate,
		Owner:         req.Owner,
		Description:   req.Description,
	})
	if offer == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "申请不存在，无法创建 Offer"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": offer})
}

func SendOffer(c *gin.Context) {
	id := c.Param("id")
	updated := store.SendOffer(id)
	if updated == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Offer 不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": updated})
}

type replyOfferRequest struct {
	Accepted bool   `json:"accepted" binding:"required"`
	Feedback string `json:"feedback"`
}

func ReplyOffer(c *gin.Context) {
	id := c.Param("id")
	var req replyOfferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updated := store.ReplyOffer(id, store.OfferReplyRequest{
		Accepted: req.Accepted,
		Feedback: req.Feedback,
	})
	if updated == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Offer 不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": updated})
}

type withdrawOfferRequest struct {
	Reason string `json:"reason"`
}

func WithdrawOffer(c *gin.Context) {
	id := c.Param("id")
	var req withdrawOfferRequest
	c.ShouldBindJSON(&req)
	updated := store.WithdrawOffer(id, req.Reason)
	if updated == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Offer 不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": updated})
}

type updateOfferNoteRequest struct {
	Note string `json:"note" binding:"required"`
}

func UpdateOfferNote(c *gin.Context) {
	id := c.Param("id")
	var req updateOfferNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updated := store.UpdateOfferNote(id, req.Note)
	if updated == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Offer 不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": updated})
}

type listRequirementsRequest struct {
	Department string `form:"department"`
	Status     string `form:"status"`
	Keyword    string `form:"keyword"`
}

func ListRequirements(c *gin.Context) {
	var req listRequirementsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result := store.ListRequirements(req.Department, req.Status, req.Keyword)
	c.JSON(http.StatusOK, gin.H{"data": result})
}

func GetRequirement(c *gin.Context) {
	id := c.Param("id")
	req := store.GetRequirement(id)
	if req == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用人需求不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": req})
}

type createRequirementRequest struct {
	Department       string               `json:"department" binding:"required"`
	JobTitle         string               `json:"jobTitle" binding:"required"`
	HeadcountType    models.HeadcountType `json:"headcountType" binding:"required"`
	Headcount        int                  `json:"headcount" binding:"required,min=1"`
	ExpectedOnboard  string               `json:"expectedOnboard" binding:"required"`
	SalaryMin        int                  `json:"salaryMin" binding:"required,min=0"`
	SalaryMax        int                  `json:"salaryMax" binding:"required,min=0"`
	WorkLocation     string               `json:"workLocation" binding:"required"`
	JobDescription   string               `json:"jobDescription" binding:"required"`
	Requirements     string               `json:"requirements" binding:"required"`
	Reason           string               `json:"reason"`
	Initiator        string               `json:"initiator" binding:"required"`
	InitiatorContact string               `json:"initiatorContact" binding:"required"`
}

func CreateRequirement(c *gin.Context) {
	var req createRequirementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.SalaryMin > req.SalaryMax {
		c.JSON(http.StatusBadRequest, gin.H{"error": "最低薪资不能大于最高薪资"})
		return
	}

	hr := &models.HiringRequirement{
		Department:       req.Department,
		JobTitle:         req.JobTitle,
		HeadcountType:    req.HeadcountType,
		Headcount:        req.Headcount,
		ExpectedOnboard:  req.ExpectedOnboard,
		SalaryMin:        req.SalaryMin,
		SalaryMax:        req.SalaryMax,
		WorkLocation:     req.WorkLocation,
		JobDescription:   req.JobDescription,
		Requirements:     req.Requirements,
		Reason:           req.Reason,
		Initiator:        req.Initiator,
		InitiatorContact: req.InitiatorContact,
	}

	created := store.CreateRequirement(hr)
	c.JSON(http.StatusOK, gin.H{"data": created})
}

type updateRequirementRequest struct {
	ID               string               `json:"id" binding:"required"`
	Department       string               `json:"department" binding:"required"`
	JobTitle         string               `json:"jobTitle" binding:"required"`
	HeadcountType    models.HeadcountType `json:"headcountType" binding:"required"`
	Headcount        int                  `json:"headcount" binding:"required,min=1"`
	ExpectedOnboard  string               `json:"expectedOnboard" binding:"required"`
	SalaryMin        int                  `json:"salaryMin" binding:"required,min=0"`
	SalaryMax        int                  `json:"salaryMax" binding:"required,min=0"`
	WorkLocation     string               `json:"workLocation" binding:"required"`
	JobDescription   string               `json:"jobDescription" binding:"required"`
	Requirements     string               `json:"requirements" binding:"required"`
	Reason           string               `json:"reason"`
	Initiator        string               `json:"initiator" binding:"required"`
	InitiatorContact string               `json:"initiatorContact" binding:"required"`
}

func UpdateRequirement(c *gin.Context) {
	id := c.Param("id")

	var req updateRequirementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if id != req.ID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "路径ID与请求体ID不一致"})
		return
	}

	if req.SalaryMin > req.SalaryMax {
		c.JSON(http.StatusBadRequest, gin.H{"error": "最低薪资不能大于最高薪资"})
		return
	}

	existing := store.GetRequirement(id)
	if existing == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用人需求不存在"})
		return
	}

	if existing.Status != models.ReqStatusDraft && existing.Status != models.ReqStatusReturned {
		c.JSON(http.StatusBadRequest, gin.H{"error": "仅草稿或退回状态的需求可以编辑"})
		return
	}

	hr := &models.HiringRequirement{
		Department:       req.Department,
		JobTitle:         req.JobTitle,
		HeadcountType:    req.HeadcountType,
		Headcount:        req.Headcount,
		ExpectedOnboard:  req.ExpectedOnboard,
		SalaryMin:        req.SalaryMin,
		SalaryMax:        req.SalaryMax,
		WorkLocation:     req.WorkLocation,
		JobDescription:   req.JobDescription,
		Requirements:     req.Requirements,
		Reason:           req.Reason,
		Initiator:        req.Initiator,
		InitiatorContact: req.InitiatorContact,
	}

	updated := store.UpdateRequirement(id, hr)
	if updated == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "更新失败，仅草稿或退回状态的需求可以编辑"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": updated})
}

func SubmitApproval(c *gin.Context) {
	id := c.Param("id")

	existing := store.GetRequirement(id)
	if existing == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用人需求不存在"})
		return
	}

	if existing.Status != models.ReqStatusDraft && existing.Status != models.ReqStatusReturned {
		c.JSON(http.StatusBadRequest, gin.H{"error": "仅草稿或退回状态的需求可以提交审批"})
		return
	}

	updated := store.SubmitApproval(id)
	if updated == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "提交审批失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": updated})
}

type approveRequirementRequest struct {
	NodeIndex *int   `json:"nodeIndex" binding:"required"`
	Action    string `json:"action" binding:"required"`
	Opinion   string `json:"opinion"`
}

func ApproveRequirement(c *gin.Context) {
	id := c.Param("id")

	var req approveRequirementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	nodeIndex := *req.NodeIndex

	validActions := map[string]bool{
		"approve": true,
		"return":  true,
		"reject":  true,
	}
	if !validActions[req.Action] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的审批操作，仅支持 approve/return/reject"})
		return
	}

	existing := store.GetRequirement(id)
	if existing == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用人需求不存在"})
		return
	}

	if existing.Status != models.ReqStatusPending {
		c.JSON(http.StatusBadRequest, gin.H{"error": "仅审批中状态的需求可以审批"})
		return
	}

	if nodeIndex < 0 || nodeIndex >= len(existing.ApprovalNodes) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的审批节点索引"})
		return
	}

	updated := store.ApproveRequirement(id, nodeIndex, req.Action, req.Opinion)
	if updated == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "审批失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": updated})
}

func ConvertToJob(c *gin.Context) {
	id := c.Param("id")

	existing := store.GetRequirement(id)
	if existing == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用人需求不存在"})
		return
	}

	if existing.Status != models.ReqStatusApproved {
		c.JSON(http.StatusBadRequest, gin.H{"error": "仅已通过审批的需求可以转换为职位"})
		return
	}

	if existing.RelatedJobID != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该需求已转换为职位"})
		return
	}

	job := store.ConvertToJob(id)
	if job == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "转换为职位失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": job})
}

func GetRequirementStats(c *gin.Context) {
	stats := store.GetRequirementStats()
	c.JSON(http.StatusOK, gin.H{"data": stats})
}
