package models

type JobStatus string

const (
	JobStatusOpen     JobStatus = "open"
	JobStatusPaused   JobStatus = "paused"
	JobStatusClosed   JobStatus = "closed"
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

type Job struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	Department   string    `json:"department"`
	Location     string    `json:"location"`
	SalaryMin    int       `json:"salaryMin"`
	SalaryMax    int       `json:"salaryMax"`
	Status       JobStatus `json:"status"`
	Headcount    int       `json:"headcount"`
	PublishedAt  string    `json:"publishedAt"`
	Description  string    `json:"description"`
	Requirements string    `json:"requirements"`
}

type Resume struct {
	CandidateName string   `json:"candidateName"`
	Contact       string   `json:"contact"`
	TargetPosition string  `json:"targetPosition"`
	YearsOfExp    int      `json:"yearsOfExp"`
	Skills        []string `json:"skills"`
	Summary       string   `json:"summary"`
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
	ID         string     `json:"id"`
	ApplicationID string  `json:"applicationId"`
	SenderName string     `json:"senderName"`
	SenderRole SenderRole `json:"senderRole"`
	Content    string     `json:"content"`
	SentAt     string     `json:"sentAt"`
}

type Stats struct {
	TotalJobs          int            `json:"totalJobs"`
	OpenJobs           int            `json:"openJobs"`
	TotalApplications  int            `json:"totalApplications"`
	ApplicationsByStatus map[string]int `json:"applicationsByStatus"`
	NewThisWeek        int            `json:"newThisWeek"`
}
