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
	InterviewMethodOnsite  InterviewMethod = "onsite"
	InterviewMethodOnline  InterviewMethod = "online"
	InterviewMethodPhone   InterviewMethod = "phone"
)

type Interview struct {
	ID              string          `json:"id"`
	ApplicationID   string          `json:"applicationId"`
	JobID           string          `json:"jobId"`
	JobTitle        string          `json:"jobTitle"`
	CandidateName   string          `json:"candidateName"`
	Contact         string          `json:"contact"`
	Round           int             `json:"round"`
	InterviewTime   string          `json:"interviewTime"`
	Method          InterviewMethod `json:"method"`
	Interviewer     string          `json:"interviewer"`
	Location        string          `json:"location"`
	Description     string          `json:"description"`
	Status          InterviewStatus `json:"status"`
	Feedback        string          `json:"feedback,omitempty"`
	LatestNote      string          `json:"latestNote,omitempty"`
	CreatedAt       string          `json:"createdAt"`
	UpdatedAt       string          `json:"updatedAt"`
}

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
	ID              string             `json:"id"`
	JobID           string             `json:"jobId"`
	JobTitle        string             `json:"jobTitle"`
	Status          ApplicationStatus  `json:"status"`
	SubmittedAt     string             `json:"submittedAt"`
	Resume          Resume             `json:"resume"`
	LastMessageTime *string            `json:"lastMessageTime,omitempty"`
	UnreadCount     int                `json:"unreadCount"`
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
	TotalJobs             int            `json:"totalJobs"`
	OpenJobs              int            `json:"openJobs"`
	TotalApplications     int            `json:"totalApplications"`
	TotalCandidates       int            `json:"totalCandidates"`
	TotalInterviews       int            `json:"totalInterviews"`
	PendingInterviews     int            `json:"pendingInterviews"`
	NewThisWeek           int            `json:"newThisWeek"`
	ApplicationsByStatus  map[string]int `json:"applicationsByStatus"`
	InterviewsByStatus    map[string]int `json:"interviewsByStatus"`
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
