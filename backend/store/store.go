package store

import (
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"

	"recruitment-platform/models"
)

var (
	jobs         []models.Job
	applications []models.Application
	messages     []models.Message
	notes        []models.Note
	mu           sync.RWMutex
)

func InitMockData() {
	mu.Lock()
	defer mu.Unlock()

	now := time.Now()

	jobs = []models.Job{
		{
			ID:           uuid.New().String(),
			Title:        "高级后端开发工程师",
			Department:   "技术部",
			Location:     "北京",
			SalaryMin:    25000,
			SalaryMax:    45000,
			Status:       models.JobStatusOpen,
			Headcount:    3,
			PublishedAt:  now.AddDate(0, 0, -7).Format("2006-01-02"),
			Description:  "负责公司核心业务系统的后端架构设计与开发",
			Requirements: "5年以上Go/Java开发经验，熟悉微服务架构",
		},
		{
			ID:           uuid.New().String(),
			Title:        "前端开发工程师",
			Department:   "技术部",
			Location:     "上海",
			SalaryMin:    20000,
			SalaryMax:    35000,
			Status:       models.JobStatusOpen,
			Headcount:    2,
			PublishedAt:  now.AddDate(0, 0, -5).Format("2006-01-02"),
			Description:  "负责公司产品的前端开发与优化",
			Requirements: "3年以上Vue/React开发经验",
		},
		{
			ID:           uuid.New().String(),
			Title:        "产品经理",
			Department:   "产品部",
			Location:     "深圳",
			SalaryMin:    22000,
			SalaryMax:    40000,
			Status:       models.JobStatusOpen,
			Headcount:    1,
			PublishedAt:  now.AddDate(0, 0, -3).Format("2006-01-02"),
			Description:  "负责B端产品的需求分析与迭代规划",
			Requirements: "3年以上B端产品经验，有SaaS经验优先",
		},
		{
			ID:           uuid.New().String(),
			Title:        "UI设计师",
			Department:   "设计部",
			Location:     "北京",
			SalaryMin:    15000,
			SalaryMax:    28000,
			Status:       models.JobStatusPaused,
			Headcount:    1,
			PublishedAt:  now.AddDate(0, 0, -10).Format("2006-01-02"),
			Description:  "负责产品界面设计与用户体验优化",
			Requirements: "2年以上UI设计经验，熟练使用Figma",
		},
		{
			ID:           uuid.New().String(),
			Title:        "数据分析师",
			Department:   "数据部",
			Location:     "杭州",
			SalaryMin:    18000,
			SalaryMax:    32000,
			Status:       models.JobStatusClosed,
			Headcount:    2,
			PublishedAt:  now.AddDate(0, 0, -14).Format("2006-01-02"),
			Description:  "负责业务数据分析与报表输出",
			Requirements: "3年以上数据分析经验，熟练SQL/Python",
		},
	}

	submittedTime1 := now.AddDate(0, 0, -4).Format(time.RFC3339)
	submittedTime2 := now.AddDate(0, 0, -2).Format(time.RFC3339)
	submittedTime3 := now.AddDate(0, 0, -6).Format(time.RFC3339)
	submittedTime4 := now.AddDate(0, 0, -1).Format(time.RFC3339)

	msgTime1 := now.AddDate(0, 0, -3).Format(time.RFC3339)
	msgTime2 := now.AddDate(0, 0, -3).Add(2 * time.Hour).Format(time.RFC3339)
	msgTime3 := now.AddDate(0, 0, -1).Format(time.RFC3339)

	applications = []models.Application{
		{
			ID:          uuid.New().String(),
			JobID:       jobs[0].ID,
			JobTitle:    jobs[0].Title,
			Status:      models.AppStatusInterview,
			SubmittedAt: submittedTime1,
			Resume: models.Resume{
				CandidateName:  "张三",
				Contact:        "13800138001",
				TargetPosition: "高级后端开发工程师",
				YearsOfExp:     6,
				Skills:         []string{"Go", "MySQL", "Redis", "Kubernetes", "微服务"},
				Summary:        "6年后端开发经验，曾主导过3个大型分布式系统的设计与实现，熟悉高并发、高可用架构。",
			},
			LastMessageTime: &msgTime2,
			UnreadCount:     1,
		},
		{
			ID:          uuid.New().String(),
			JobID:       jobs[1].ID,
			JobTitle:    jobs[1].Title,
			Status:      models.AppStatusPending,
			SubmittedAt: submittedTime2,
			Resume: models.Resume{
				CandidateName:  "李四",
				Contact:        "13800138002",
				TargetPosition: "前端开发工程师",
				YearsOfExp:     4,
				Skills:         []string{"Vue3", "TypeScript", "Vite", "Element Plus", "Pinia"},
				Summary:        "4年前端开发经验，精通Vue生态，有丰富的中后台系统开发经验。",
			},
			UnreadCount: 0,
		},
		{
			ID:          uuid.New().String(),
			JobID:       jobs[0].ID,
			JobTitle:    jobs[0].Title,
			Status:      models.AppStatusRejected,
			SubmittedAt: submittedTime3,
			Resume: models.Resume{
				CandidateName:  "王五",
				Contact:        "13800138003",
				TargetPosition: "高级后端开发工程师",
				YearsOfExp:     2,
				Skills:         []string{"Java", "Spring Boot", "MySQL"},
				Summary:        "2年Java开发经验，有一定的分布式系统开发基础。",
			},
			UnreadCount: 0,
		},
		{
			ID:          uuid.New().String(),
			JobID:       jobs[2].ID,
			JobTitle:    jobs[2].Title,
			Status:      models.AppStatusSubmitted,
			SubmittedAt: submittedTime4,
			Resume: models.Resume{
				CandidateName:  "赵六",
				Contact:        "13800138004",
				TargetPosition: "产品经理",
				YearsOfExp:     5,
				Skills:         []string{"需求分析", "原型设计", "Axure", "数据分析"},
				Summary:        "5年B端产品经理经验，主导过多个SaaS产品从0到1的设计与落地。",
			},
			UnreadCount: 0,
		},
		{
			ID:          uuid.New().String(),
			JobID:       jobs[1].ID,
			JobTitle:    jobs[1].Title,
			Status:      models.AppStatusHired,
			SubmittedAt: now.AddDate(0, 0, -20).Format(time.RFC3339),
			Resume: models.Resume{
				CandidateName:  "钱七",
				Contact:        "13800138005",
				TargetPosition: "前端开发工程师",
				YearsOfExp:     5,
				Skills:         []string{"React", "Vue", "Node.js", "Webpack"},
				Summary:        "5年前端开发经验，全栈能力，熟悉前端工程化。",
			},
			UnreadCount: 0,
		},
	}

	messages = []models.Message{
		{
			ID:            uuid.New().String(),
			ApplicationID: applications[0].ID,
			SenderName:    "招聘专员-刘经理",
			SenderRole:    models.RoleRecruiter,
			Content:       "您好，看到您的简历非常优秀，方便约个时间沟通一下吗？",
			SentAt:        msgTime1,
		},
		{
			ID:            uuid.New().String(),
			ApplicationID: applications[0].ID,
			SenderName:    "张三",
			SenderRole:    models.RoleCandidate,
			Content:       "您好刘经理，非常感谢关注！我本周三下午和周四上午都有空，您看哪个时间合适？",
			SentAt:        msgTime2,
		},
		{
			ID:            uuid.New().String(),
			ApplicationID: applications[2].ID,
			SenderName:    "招聘专员-刘经理",
			SenderRole:    models.RoleRecruiter,
			Content:       "感谢您投递我们公司的职位，经过综合评估，您的工作经验与当前岗位要求不太匹配，祝您早日找到理想的工作！",
			SentAt:        msgTime3,
		},
	}

	notes = []models.Note{
		{
			ID:        uuid.New().String(),
			Candidate: "张三",
			Contact:   "13800138001",
			Type:      models.NoteTypeInterview,
			Content:   "技术一面评价：基础扎实，对 Go 并发模型理解深刻，项目经验丰富。建议进行二面。",
			CreatedBy: "招聘专员-刘经理",
			CreatedAt: now.AddDate(0, 0, -2).Format(time.RFC3339),
		},
		{
			ID:        uuid.New().String(),
			Candidate: "张三",
			Contact:   "13800138001",
			Type:      models.NoteTypeScreen,
			Content:   "初筛意见：学历背景好，大厂工作经历，期望薪资在预算范围内，可推进。",
			CreatedBy: "招聘专员-刘经理",
			CreatedAt: now.AddDate(0, 0, -4).Format(time.RFC3339),
		},
		{
			ID:        uuid.New().String(),
			Candidate: "李四",
			Contact:   "13800138002",
			Type:      models.NoteTypeScreen,
			Content:   "简历亮点：主导过中后台系统重构，性能优化提升40%。建议邀请面试。",
			CreatedBy: "招聘专员-刘经理",
			CreatedAt: now.AddDate(0, 0, -1).Format(time.RFC3339),
		},
	}
}

func ListJobs(keyword, location, status string) []models.Job {
	mu.RLock()
	defer mu.RUnlock()

	result := make([]models.Job, 0)
	for _, job := range jobs {
		if keyword != "" && !containsKeyword(job, keyword) {
			continue
		}
		if location != "" && job.Location != location {
			continue
		}
		if status != "" && string(job.Status) != status {
			continue
		}
		result = append(result, job)
	}
	return result
}

func containsKeyword(job models.Job, keyword string) bool {
	return containsIgnoreCase(job.Title, keyword) ||
		containsIgnoreCase(job.Department, keyword) ||
		containsIgnoreCase(job.Description, keyword) ||
		containsIgnoreCase(job.Requirements, keyword)
}

func containsIgnoreCase(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	lowS := strings.ToLower(s)
	lowSub := strings.ToLower(substr)
	for i := 0; i <= len(lowS)-len(lowSub); i++ {
		if lowS[i:i+len(lowSub)] == lowSub {
			return true
		}
	}
	return false
}

func GetJob(id string) *models.Job {
	mu.RLock()
	defer mu.RUnlock()

	for i := range jobs {
		if jobs[i].ID == id {
			return &jobs[i]
		}
	}
	return nil
}

func CreateJob(job *models.Job) models.Job {
	mu.Lock()
	defer mu.Unlock()

	job.ID = uuid.New().String()
	if job.PublishedAt == "" {
		job.PublishedAt = time.Now().Format("2006-01-02")
	}
	jobs = append(jobs, *job)
	return *job
}

func UpdateJob(id string, job *models.Job) *models.Job {
	mu.Lock()
	defer mu.Unlock()

	for i := range jobs {
		if jobs[i].ID == id {
			job.ID = id
			if job.PublishedAt == "" {
				job.PublishedAt = jobs[i].PublishedAt
			}
			jobs[i] = *job
			return &jobs[i]
		}
	}
	return nil
}

func ListApplications(status, jobId string) []models.Application {
	mu.RLock()
	defer mu.RUnlock()

	result := make([]models.Application, 0)
	for _, app := range applications {
		if status != "" && string(app.Status) != status {
			continue
		}
		if jobId != "" && app.JobID != jobId {
			continue
		}
		result = append(result, app)
	}
	return result
}

func GetApplication(id string) *models.Application {
	mu.RLock()
	defer mu.RUnlock()

	for i := range applications {
		if applications[i].ID == id {
			return &applications[i]
		}
	}
	return nil
}

func CreateApplication(app *models.Application, jobId string) *models.Application {
	mu.Lock()
	defer mu.Unlock()

	job := GetJob(jobId)
	if job == nil {
		return nil
	}

	app.ID = uuid.New().String()
	app.JobID = jobId
	app.JobTitle = job.Title
	app.Status = models.AppStatusSubmitted
	app.SubmittedAt = time.Now().Format(time.RFC3339)
	app.UnreadCount = 0
	applications = append(applications, *app)
	return &applications[len(applications)-1]
}

func UpdateApplicationStatus(id string, status models.ApplicationStatus) *models.Application {
	mu.Lock()
	defer mu.Unlock()

	for i := range applications {
		if applications[i].ID == id {
			applications[i].Status = status
			return &applications[i]
		}
	}
	return nil
}

func ListMessages(applicationId string) []models.Message {
	mu.RLock()
	defer mu.RUnlock()

	result := make([]models.Message, 0)
	for _, msg := range messages {
		if msg.ApplicationID == applicationId {
			result = append(result, msg)
		}
	}
	return result
}

func SendMessage(msg *models.Message) models.Message {
	mu.Lock()
	defer mu.Unlock()

	msg.ID = uuid.New().String()
	msg.SentAt = time.Now().Format(time.RFC3339)
	messages = append(messages, *msg)

	for i := range applications {
		if applications[i].ID == msg.ApplicationID {
			t := msg.SentAt
			applications[i].LastMessageTime = &t
			if msg.SenderRole == models.RoleCandidate {
				applications[i].UnreadCount++
			}
		}
	}

	return *msg
}

func buildCandidate(apps []models.Application, noteCount int) models.Candidate {
	latest := apps[0]
	appSimples := make([]models.ApplicationSimple, 0, len(apps))
	for _, app := range apps {
		appSimples = append(appSimples, models.ApplicationSimple{
			ID:          app.ID,
			JobID:       app.JobID,
			JobTitle:    app.JobTitle,
			Status:      app.Status,
			SubmittedAt: app.SubmittedAt,
		})
	}

	return models.Candidate{
		CandidateName:   latest.Resume.CandidateName,
		Contact:         latest.Resume.Contact,
		TargetPosition:  latest.Resume.TargetPosition,
		YearsOfExp:      latest.Resume.YearsOfExp,
		Skills:          latest.Resume.Skills,
		Summary:         latest.Resume.Summary,
		Applications:    appSimples,
		LatestStatus:    latest.Status,
		LatestJobTitle:  latest.JobTitle,
		LastMessageTime: latest.LastMessageTime,
		AppliedAt:       latest.SubmittedAt,
		NoteCount:       noteCount,
	}
}

func candidateMatchesFilter(c models.Candidate, keyword, jobId, status string) bool {
	if keyword != "" {
		kw := strings.ToLower(keyword)
		if !strings.Contains(strings.ToLower(c.CandidateName), kw) &&
			!strings.Contains(strings.ToLower(c.Contact), kw) &&
			!strings.Contains(strings.ToLower(c.TargetPosition), kw) &&
			!strings.Contains(strings.ToLower(c.Summary), kw) {
			skillMatch := false
			for _, s := range c.Skills {
				if strings.Contains(strings.ToLower(s), kw) {
					skillMatch = true
					break
				}
			}
			if !skillMatch {
				return false
			}
		}
	}
	if jobId != "" {
		jobMatch := false
		for _, a := range c.Applications {
			if a.JobID == jobId {
				jobMatch = true
				break
			}
		}
		if !jobMatch {
			return false
		}
	}
	if status != "" && string(c.LatestStatus) != status {
		return false
	}
	return true
}

func ListCandidates(keyword, jobId, status string) []models.Candidate {
	mu.RLock()
	defer mu.RUnlock()

	grouped := make(map[string][]models.Application)
	for _, app := range applications {
		key := app.Resume.Contact + "|" + app.Resume.CandidateName
		grouped[key] = append(grouped[key], app)
	}

	noteCountMap := make(map[string]int)
	for _, note := range notes {
		key := note.Contact + "|" + note.Candidate
		noteCountMap[key]++
	}

	result := make([]models.Candidate, 0)
	for key, apps := range grouped {
		sort.Slice(apps, func(i, j int) bool {
			return apps[i].SubmittedAt > apps[j].SubmittedAt
		})

		noteCount := noteCountMap[key]
		candidate := buildCandidate(apps, noteCount)
		if candidateMatchesFilter(candidate, keyword, jobId, status) {
			result = append(result, candidate)
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].AppliedAt > result[j].AppliedAt
	})

	return result
}

func GetCandidate(contact, name string) *models.CandidateDetail {
	mu.RLock()
	defer mu.RUnlock()

	var candidateApps []models.Application
	for _, app := range applications {
		if app.Resume.Contact == contact && app.Resume.CandidateName == name {
			candidateApps = append(candidateApps, app)
		}
	}
	if len(candidateApps) == 0 {
		return nil
	}

	sort.Slice(candidateApps, func(i, j int) bool {
		return candidateApps[i].SubmittedAt > candidateApps[j].SubmittedAt
	})

	var candidateNotes []models.Note
	for _, note := range notes {
		if note.Contact == contact && note.Candidate == name {
			candidateNotes = append(candidateNotes, note)
		}
	}
	sort.Slice(candidateNotes, func(i, j int) bool {
		return candidateNotes[i].CreatedAt > candidateNotes[j].CreatedAt
	})

	appIds := make(map[string]bool)
	for _, app := range candidateApps {
		appIds[app.ID] = true
	}

	var candidateMessages []models.Message
	for _, msg := range messages {
		if appIds[msg.ApplicationID] {
			candidateMessages = append(candidateMessages, msg)
		}
	}
	sort.Slice(candidateMessages, func(i, j int) bool {
		return candidateMessages[i].SentAt < candidateMessages[j].SentAt
	})

	noteCount := len(candidateNotes)
	base := buildCandidate(candidateApps, noteCount)

	return &models.CandidateDetail{
		Candidate: base,
		Messages:  candidateMessages,
		Notes:     candidateNotes,
	}
}

func UpdateCandidateStatus(contact, name string, status models.ApplicationStatus) error {
	mu.Lock()
	defer mu.Unlock()

	updated := false
	for i := range applications {
		if applications[i].Resume.Contact == contact && applications[i].Resume.CandidateName == name {
			applications[i].Status = status
			updated = true
		}
	}
	if !updated {
		return nil
	}
	return nil
}

func ListNotes(contact, name string) []models.Note {
	mu.RLock()
	defer mu.RUnlock()

	result := make([]models.Note, 0)
	for _, note := range notes {
		if contact != "" && note.Contact != contact {
			continue
		}
		if name != "" && note.Candidate != name {
			continue
		}
		result = append(result, note)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt > result[j].CreatedAt
	})
	return result
}

func CreateNote(note *models.Note) models.Note {
	mu.Lock()
	defer mu.Unlock()

	note.ID = uuid.New().String()
	note.CreatedAt = time.Now().Format(time.RFC3339)
	if note.CreatedBy == "" {
		note.CreatedBy = "招聘专员-刘经理"
	}
	notes = append(notes, *note)
	return *note
}

func GetStats() models.Stats {
	mu.RLock()
	defer mu.RUnlock()

	stats := models.Stats{
		ApplicationsByStatus: make(map[string]int),
	}

	weekAgo := time.Now().AddDate(0, 0, -7)
	candidateSet := make(map[string]bool)

	for _, job := range jobs {
		stats.TotalJobs++
		if job.Status == models.JobStatusOpen {
			stats.OpenJobs++
		}
	}

	for _, app := range applications {
		stats.TotalApplications++
		stats.ApplicationsByStatus[string(app.Status)]++
		key := app.Resume.Contact + "|" + app.Resume.CandidateName
		candidateSet[key] = true

		submittedTime, err := time.Parse(time.RFC3339, app.SubmittedAt)
		if err == nil && submittedTime.After(weekAgo) {
			stats.NewThisWeek++
		}
	}

	stats.TotalCandidates = len(candidateSet)

	return stats
}
