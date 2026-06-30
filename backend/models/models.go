package models

type JobStatus string

const (
	JobStatusOpen   JobStatus = "open"
	JobStatusPaused JobStatus = "paused"
	JobStatusClosed JobStatus = "closed"
)

type ApplicationStatus string

const (
	AppStatusSubmitted ApplicationStatus = "submitted"
	AppStatusPending   ApplicationStatus = "pending"
	AppStatusInterview ApplicationStatus = "interview"
	AppStatusRejected  ApplicationStatus = "rejected"
	AppStatusHired     ApplicationStatus = "hired"
)

type SenderRole string

const (
	RoleRecruiter SenderRole = "recruiter"
	RoleCandidate SenderRole = "candidate"
)

type NoteType string

const (
	NoteTypeInterview NoteType = "interview"
	NoteTypeScreen    NoteType = "screen"
)

type InterviewStatus string

const (
	InterviewStatusScheduled InterviewStatus = "scheduled"
	InterviewStatusCompleted InterviewStatus = "completed"
	InterviewStatusCancelled InterviewStatus = "cancelled"
)

type InterviewMethod string

const (
	InterviewMethodOnsite InterviewMethod = "onsite"
	InterviewMethodOnline InterviewMethod = "online"
	InterviewMethodPhone  InterviewMethod = "phone"
)

type OfferStatus string

const (
	OfferStatusPending   OfferStatus = "pending"
	OfferStatusSent      OfferStatus = "sent"
	OfferStatusAccepted  OfferStatus = "accepted"
	OfferStatusRejected  OfferStatus = "rejected"
	OfferStatusWithdrawn OfferStatus = "withdrawn"
)

type Offer struct {
	ID            string      `json:"id"`
	ApplicationID string      `json:"applicationId"`
	JobID         string      `json:"jobId"`
	JobTitle      string      `json:"jobTitle"`
	CandidateName string      `json:"candidateName"`
	Contact       string      `json:"contact"`
	SalaryMin     int         `json:"salaryMin"`
	SalaryMax     int         `json:"salaryMax"`
	StartDate     string      `json:"startDate"`
	Owner         string      `json:"owner"`
	Description   string      `json:"description"`
	Status        OfferStatus `json:"status"`
	Feedback      string      `json:"feedback,omitempty"`
	CreatedAt     string      `json:"createdAt"`
	UpdatedAt     string      `json:"updatedAt"`
	SentAt        string      `json:"sentAt,omitempty"`
	RepliedAt     string      `json:"repliedAt,omitempty"`
}

type InterviewFeedbackConclusion string

const (
	FeedbackPass    InterviewFeedbackConclusion = "pass"
	FeedbackPending InterviewFeedbackConclusion = "pending"
	FeedbackFail    InterviewFeedbackConclusion = "fail"
)

type Interview struct {
	ID            string                      `json:"id"`
	ApplicationID string                      `json:"applicationId"`
	JobID         string                      `json:"jobId"`
	JobTitle      string                      `json:"jobTitle"`
	CandidateName string                      `json:"candidateName"`
	Contact       string                      `json:"contact"`
	Round         int                         `json:"round"`
	InterviewTime string                      `json:"interviewTime"`
	Method        InterviewMethod             `json:"method"`
	Interviewer   string                      `json:"interviewer"`
	Location      string                      `json:"location"`
	Description   string                      `json:"description"`
	Status        InterviewStatus             `json:"status"`
	Feedback      string                      `json:"feedback,omitempty"`
	Conclusion    InterviewFeedbackConclusion `json:"conclusion,omitempty"`
	Rating        int                         `json:"rating,omitempty"`
	Strengths     string                      `json:"strengths,omitempty"`
	Risks         string                      `json:"risks,omitempty"`
	NextSteps     string                      `json:"nextSteps,omitempty"`
	LatestNote    string                      `json:"latestNote,omitempty"`
	CreatedAt     string                      `json:"createdAt"`
	UpdatedAt     string                      `json:"updatedAt"`
}

type Job struct {
	ID            string    `json:"id"`
	Title         string    `json:"title"`
	Department    string    `json:"department"`
	Location      string    `json:"location"`
	SalaryMin     int       `json:"salaryMin"`
	SalaryMax     int       `json:"salaryMax"`
	Status        JobStatus `json:"status"`
	Headcount     int       `json:"headcount"`
	PublishedAt   string    `json:"publishedAt"`
	Description   string    `json:"description"`
	Requirements  string    `json:"requirements"`
	RequirementId string    `json:"requirementId,omitempty"`
}

type Resume struct {
	CandidateName  string   `json:"candidateName"`
	Contact        string   `json:"contact"`
	TargetPosition string   `json:"targetPosition"`
	YearsOfExp     int      `json:"yearsOfExp"`
	Skills         []string `json:"skills"`
	Summary        string   `json:"summary"`
}

type Application struct {
	ID              string            `json:"id"`
	JobID           string            `json:"jobId"`
	JobTitle        string            `json:"jobTitle"`
	Status          ApplicationStatus `json:"status"`
	SubmittedAt     string            `json:"submittedAt"`
	Resume          Resume            `json:"resume"`
	LastMessageTime *string           `json:"lastMessageTime,omitempty"`
	UnreadCount     int               `json:"unreadCount"`
}

type Message struct {
	ID            string     `json:"id"`
	ApplicationID string     `json:"applicationId"`
	SenderName    string     `json:"senderName"`
	SenderRole    SenderRole `json:"senderRole"`
	Content       string     `json:"content"`
	SentAt        string     `json:"sentAt"`
}

type Note struct {
	ID        string   `json:"id"`
	Candidate string   `json:"candidate"`
	Contact   string   `json:"contact"`
	Type      NoteType `json:"type"`
	Content   string   `json:"content"`
	CreatedBy string   `json:"createdBy"`
	CreatedAt string   `json:"createdAt"`
}

type Stats struct {
	TotalJobs            int              `json:"totalJobs"`
	OpenJobs             int              `json:"openJobs"`
	TotalApplications    int              `json:"totalApplications"`
	TotalCandidates      int              `json:"totalCandidates"`
	TotalInterviews      int              `json:"totalInterviews"`
	PendingInterviews    int              `json:"pendingInterviews"`
	TotalOffers          int              `json:"totalOffers"`
	PendingOffers        int              `json:"pendingOffers"`
	AcceptedOffers       int              `json:"acceptedOffers"`
	NewThisWeek          int              `json:"newThisWeek"`
	ApplicationsByStatus map[string]int   `json:"applicationsByStatus"`
	InterviewsByStatus   map[string]int   `json:"interviewsByStatus"`
	OffersByStatus       map[string]int   `json:"offersByStatus"`
	RequirementStats     RequirementStats `json:"requirementStats"`
}

type ApplicationSimple struct {
	ID          string            `json:"id"`
	JobID       string            `json:"jobId"`
	JobTitle    string            `json:"jobTitle"`
	Status      ApplicationStatus `json:"status"`
	SubmittedAt string            `json:"submittedAt"`
}

type Candidate struct {
	CandidateName   string              `json:"candidateName"`
	Contact         string              `json:"contact"`
	TargetPosition  string              `json:"targetPosition"`
	YearsOfExp      int                 `json:"yearsOfExp"`
	Skills          []string            `json:"skills"`
	Summary         string              `json:"summary"`
	Applications    []ApplicationSimple `json:"applications"`
	LatestStatus    ApplicationStatus   `json:"latestStatus"`
	LatestJobTitle  string              `json:"latestJobTitle"`
	LastMessageTime *string             `json:"lastMessageTime,omitempty"`
	AppliedAt       string              `json:"appliedAt"`
	NoteCount       int                 `json:"noteCount"`
}

type CandidateDetail struct {
	Candidate Candidate `json:"candidate"`
	Messages  []Message `json:"messages"`
	Notes     []Note    `json:"notes"`
}

type RequirementStatus string

const (
	ReqStatusDraft     RequirementStatus = "draft"
	ReqStatusPending   RequirementStatus = "pending"
	ReqStatusApproved  RequirementStatus = "approved"
	ReqStatusReturned  RequirementStatus = "returned"
	ReqStatusRejected  RequirementStatus = "rejected"
	ReqStatusConverted RequirementStatus = "converted"
)

type HeadcountType string

const (
	HeadcountRegular     HeadcountType = "regular"
	HeadcountContract    HeadcountType = "contract"
	HeadcountIntern      HeadcountType = "intern"
	HeadcountReplacement HeadcountType = "replacement"
)

type RequirementApprovalNode struct {
	ID        string            `json:"id"`
	NodeName  string            `json:"nodeName"`
	Approver  string            `json:"approver"`
	Status    RequirementStatus `json:"status"`
	Opinion   string            `json:"opinion,omitempty"`
	HandledAt string            `json:"handledAt,omitempty"`
	CreatedAt string            `json:"createdAt"`
}

type RequirementOperationLog struct {
	ID        string `json:"id"`
	Action    string `json:"action"`
	Operator  string `json:"operator"`
	Detail    string `json:"detail,omitempty"`
	CreatedAt string `json:"createdAt"`
}

type HiringRequirement struct {
	ID               string                    `json:"id"`
	ReqNo            string                    `json:"reqNo"`
	Department       string                    `json:"department"`
	JobTitle         string                    `json:"jobTitle"`
	HeadcountType    HeadcountType             `json:"headcountType"`
	Headcount        int                       `json:"headcount"`
	ExpectedOnboard  string                    `json:"expectedOnboard"`
	SalaryMin        int                       `json:"salaryMin"`
	SalaryMax        int                       `json:"salaryMax"`
	WorkLocation     string                    `json:"workLocation"`
	JobDescription   string                    `json:"jobDescription"`
	Requirements     string                    `json:"requirements"`
	Reason           string                    `json:"reason,omitempty"`
	Status           RequirementStatus         `json:"status"`
	Initiator        string                    `json:"initiator"`
	InitiatorContact string                    `json:"initiatorContact"`
	ApprovalNodes    []RequirementApprovalNode `json:"approvalNodes"`
	OperationLogs    []RequirementOperationLog `json:"operationLogs"`
	RelatedJobID     string                    `json:"relatedJobId,omitempty"`
	LatestOpinion    string                    `json:"latestOpinion,omitempty"`
	CreatedAt        string                    `json:"createdAt"`
	UpdatedAt        string                    `json:"updatedAt"`
}

type RequirementStats struct {
	PendingCount   int            `json:"pendingCount"`
	ApprovedCount  int            `json:"approvedCount"`
	ReturnedCount  int            `json:"returnedCount"`
	RejectedCount  int            `json:"rejectedCount"`
	ConvertedCount int            `json:"convertedCount"`
	DraftCount     int            `json:"draftCount"`
	UrgentCount    int            `json:"urgentCount"`
	ByDepartment   map[string]int `json:"byDepartment"`
	ByMonth        map[string]int `json:"byMonth"`
}
