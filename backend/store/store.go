package store

import (
	"fmt"
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
	interviews   []models.Interview
	offers       []models.Offer
	requirements []models.HiringRequirement
	reqSeqNo     int
	mu           sync.RWMutex
)

func genReqNo() string {
	reqSeqNo++
	return fmt.Sprintf("REQ-%04d", reqSeqNo)
}

func InitMockData() {
	mu.Lock()
	defer mu.Unlock()

	now := time.Now()

	reqSeqNo = 0

	requirements = []models.HiringRequirement{
		{
			ID:               uuid.New().String(),
			ReqNo:            genReqNo(),
			Department:       "技术部",
			JobTitle:         "高级后端开发工程师",
			HeadcountType:    models.HeadcountRegular,
			Headcount:        3,
			ExpectedOnboard:  now.AddDate(0, 0, 20).Format("2006-01-02"),
			SalaryMin:        25000,
			SalaryMax:        45000,
			WorkLocation:     "北京",
			JobDescription:   "负责公司核心业务系统的后端架构设计与开发",
			Requirements:     "5年以上Go/Java开发经验，熟悉微服务架构",
			Reason:           "业务扩张",
			Status:           models.ReqStatusApproved,
			Initiator:        "张总监",
			InitiatorContact: "13800000001",
			ApprovalNodes:    nil,
			OperationLogs:    nil,
			CreatedAt:        now.AddDate(0, 0, -30).Format(time.RFC3339),
			UpdatedAt:        now.AddDate(0, 0, -25).Format(time.RFC3339),
		},
		{
			ID:               uuid.New().String(),
			ReqNo:            genReqNo(),
			Department:       "技术部",
			JobTitle:         "前端开发工程师",
			HeadcountType:    models.HeadcountRegular,
			Headcount:        2,
			ExpectedOnboard:  now.AddDate(0, 0, 15).Format("2006-01-02"),
			SalaryMin:        20000,
			SalaryMax:        35000,
			WorkLocation:     "上海",
			JobDescription:   "负责公司产品的前端开发与优化",
			Requirements:     "3年以上Vue/React开发经验",
			Reason:           "新项目启动",
			Status:           models.ReqStatusPending,
			Initiator:        "李经理",
			InitiatorContact: "13800000002",
			ApprovalNodes:    nil,
			OperationLogs:    nil,
			CreatedAt:        now.AddDate(0, 0, -8).Format(time.RFC3339),
			UpdatedAt:        now.AddDate(0, 0, -8).Format(time.RFC3339),
		},
		{
			ID:               uuid.New().String(),
			ReqNo:            genReqNo(),
			Department:       "产品部",
			JobTitle:         "产品经理",
			HeadcountType:    models.HeadcountRegular,
			Headcount:        1,
			ExpectedOnboard:  now.AddDate(0, 0, 10).Format("2006-01-02"),
			SalaryMin:        22000,
			SalaryMax:        40000,
			WorkLocation:     "深圳",
			JobDescription:   "负责B端产品的需求分析与迭代规划",
			Requirements:     "3年以上B端产品经验，有SaaS经验优先",
			Reason:           "产品线扩张",
			Status:           models.ReqStatusReturned,
			Initiator:        "王主管",
			InitiatorContact: "13800000003",
			ApprovalNodes:    nil,
			OperationLogs:    nil,
			LatestOpinion:    "请补充岗位KPI指标",
			CreatedAt:        now.AddDate(0, 0, -12).Format(time.RFC3339),
			UpdatedAt:        now.AddDate(0, 0, -5).Format(time.RFC3339),
		},
		{
			ID:               uuid.New().String(),
			ReqNo:            genReqNo(),
			Department:       "设计部",
			JobTitle:         "UI设计师",
			HeadcountType:    models.HeadcountContract,
			Headcount:        1,
			ExpectedOnboard:  now.AddDate(0, 0, 5).Format("2006-01-02"),
			SalaryMin:        15000,
			SalaryMax:        28000,
			WorkLocation:     "北京",
			JobDescription:   "负责产品界面设计与用户体验优化",
			Requirements:     "2年以上UI设计经验，熟练使用Figma",
			Reason:           "外包项目需求",
			Status:           models.ReqStatusRejected,
			Initiator:        "陈设计",
			InitiatorContact: "13800000004",
			ApprovalNodes:    nil,
			OperationLogs:    nil,
			LatestOpinion:    "预算不足，暂不审批",
			CreatedAt:        now.AddDate(0, 0, -20).Format(time.RFC3339),
			UpdatedAt:        now.AddDate(0, 0, -15).Format(time.RFC3339),
		},
		{
			ID:               uuid.New().String(),
			ReqNo:            genReqNo(),
			Department:       "数据部",
			JobTitle:         "数据分析师",
			HeadcountType:    models.HeadcountRegular,
			Headcount:        2,
			ExpectedOnboard:  now.AddDate(0, 0, 3).Format("2006-01-02"),
			SalaryMin:        18000,
			SalaryMax:        32000,
			WorkLocation:     "杭州",
			JobDescription:   "负责业务数据分析与报表输出",
			Requirements:     "3年以上数据分析经验，熟练SQL/Python",
			Reason:           "数据团队扩充",
			Status:           models.ReqStatusConverted,
			Initiator:        "刘数据",
			InitiatorContact: "13800000005",
			ApprovalNodes:    nil,
			OperationLogs:    nil,
			CreatedAt:        now.AddDate(0, 0, -40).Format(time.RFC3339),
			UpdatedAt:        now.AddDate(0, 0, -14).Format(time.RFC3339),
		},
		{
			ID:               uuid.New().String(),
			ReqNo:            genReqNo(),
			Department:       "技术部",
			JobTitle:         "测试工程师",
			HeadcountType:    models.HeadcountRegular,
			Headcount:        2,
			ExpectedOnboard:  now.AddDate(0, 0, 7).Format("2006-01-02"),
			SalaryMin:        18000,
			SalaryMax:        30000,
			WorkLocation:     "北京",
			JobDescription:   "负责产品功能测试与自动化测试",
			Requirements:     "3年以上测试经验，熟悉自动化测试框架",
			Reason:           "质量保障体系建设",
			Status:           models.ReqStatusPending,
			Initiator:        "赵测试",
			InitiatorContact: "13800000006",
			ApprovalNodes:    nil,
			OperationLogs:    nil,
			CreatedAt:        now.AddDate(0, 0, -3).Format(time.RFC3339),
			UpdatedAt:        now.AddDate(0, 0, -3).Format(time.RFC3339),
		},
		{
			ID:               uuid.New().String(),
			ReqNo:            genReqNo(),
			Department:       "运营部",
			JobTitle:         "运营专员",
			HeadcountType:    models.HeadcountIntern,
			Headcount:        3,
			ExpectedOnboard:  now.AddDate(0, 0, 4).Format("2006-01-02"),
			SalaryMin:        8000,
			SalaryMax:        12000,
			WorkLocation:     "广州",
			JobDescription:   "负责日常运营活动策划与执行",
			Requirements:     "有运营实习经验优先，应届生也可",
			Reason:           "暑期实习生",
			Status:           models.ReqStatusApproved,
			Initiator:        "孙运营",
			InitiatorContact: "13800000007",
			ApprovalNodes:    nil,
			OperationLogs:    nil,
			CreatedAt:        now.AddDate(0, 0, -10).Format(time.RFC3339),
			UpdatedAt:        now.AddDate(0, 0, -8).Format(time.RFC3339),
		},
		{
			ID:               uuid.New().String(),
			ReqNo:            genReqNo(),
			Department:       "人事部",
			JobTitle:         "HRBP",
			HeadcountType:    models.HeadcountReplacement,
			Headcount:        1,
			ExpectedOnboard:  now.AddDate(0, 0, 18).Format("2006-01-02"),
			SalaryMin:        20000,
			SalaryMax:        35000,
			WorkLocation:     "北京",
			JobDescription:   "负责技术部门HRBP工作",
			Requirements:     "5年以上HRBP经验，有互联网公司背景优先",
			Reason:           "人员替换",
			Status:           models.ReqStatusDraft,
			Initiator:        "周HR",
			InitiatorContact: "13800000008",
			ApprovalNodes:    nil,
			OperationLogs:    nil,
			CreatedAt:        now.AddDate(0, 0, -2).Format(time.RFC3339),
			UpdatedAt:        now.AddDate(0, 0, -2).Format(time.RFC3339),
		},
	}

	jobs = []models.Job{
		{
			ID:            uuid.New().String(),
			Title:         "高级后端开发工程师",
			Department:    "技术部",
			Location:      "北京",
			SalaryMin:     25000,
			SalaryMax:     45000,
			Status:        models.JobStatusOpen,
			Headcount:     3,
			PublishedAt:   now.AddDate(0, 0, -7).Format("2006-01-02"),
			Description:   "负责公司核心业务系统的后端架构设计与开发",
			Requirements:  "5年以上Go/Java开发经验，熟悉微服务架构",
			RequirementId: requirements[0].ID,
		},
		{
			ID:            uuid.New().String(),
			Title:         "前端开发工程师",
			Department:    "技术部",
			Location:      "上海",
			SalaryMin:     20000,
			SalaryMax:     35000,
			Status:        models.JobStatusOpen,
			Headcount:     2,
			PublishedAt:   now.AddDate(0, 0, -5).Format("2006-01-02"),
			Description:   "负责公司产品的前端开发与优化",
			Requirements:  "3年以上Vue/React开发经验",
			RequirementId: requirements[1].ID,
		},
		{
			ID:            uuid.New().String(),
			Title:         "产品经理",
			Department:    "产品部",
			Location:      "深圳",
			SalaryMin:     22000,
			SalaryMax:     40000,
			Status:        models.JobStatusOpen,
			Headcount:     1,
			PublishedAt:   now.AddDate(0, 0, -3).Format("2006-01-02"),
			Description:   "负责B端产品的需求分析与迭代规划",
			Requirements:  "3年以上B端产品经验，有SaaS经验优先",
			RequirementId: requirements[2].ID,
		},
		{
			ID:            uuid.New().String(),
			Title:         "UI设计师",
			Department:    "设计部",
			Location:      "北京",
			SalaryMin:     15000,
			SalaryMax:     28000,
			Status:        models.JobStatusPaused,
			Headcount:     1,
			PublishedAt:   now.AddDate(0, 0, -10).Format("2006-01-02"),
			Description:   "负责产品界面设计与用户体验优化",
			Requirements:  "2年以上UI设计经验，熟练使用Figma",
			RequirementId: requirements[3].ID,
		},
		{
			ID:            uuid.New().String(),
			Title:         "数据分析师",
			Department:    "数据部",
			Location:      "杭州",
			SalaryMin:     18000,
			SalaryMax:     32000,
			Status:        models.JobStatusClosed,
			Headcount:     2,
			PublishedAt:   now.AddDate(0, 0, -14).Format("2006-01-02"),
			Description:   "负责业务数据分析与报表输出",
			Requirements:  "3年以上数据分析经验，熟练SQL/Python",
			RequirementId: requirements[4].ID,
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

	interviews = []models.Interview{
		{
			ID:            uuid.New().String(),
			ApplicationID: applications[0].ID,
			JobID:         jobs[0].ID,
			JobTitle:      jobs[0].Title,
			CandidateName: "张三",
			Contact:       "13800138001",
			Round:         1,
			InterviewTime: now.AddDate(0, 0, 1).Add(10 * time.Hour).Format(time.RFC3339),
			Method:        models.InterviewMethodOnline,
			Interviewer:   "技术总监-王工",
			Location:      "腾讯会议",
			Description:   "技术一面，重点考察 Go 基础、微服务架构和项目经验",
			Status:        models.InterviewStatusScheduled,
			LatestNote:    "已发送面试邀请，候选人已确认",
			CreatedAt:     now.AddDate(0, 0, -1).Format(time.RFC3339),
			UpdatedAt:     now.AddDate(0, 0, -1).Format(time.RFC3339),
		},
		{
			ID:            uuid.New().String(),
			ApplicationID: applications[1].ID,
			JobID:         jobs[1].ID,
			JobTitle:      jobs[1].Title,
			CandidateName: "李四",
			Contact:       "13800138002",
			Round:         1,
			InterviewTime: now.AddDate(0, 0, 2).Add(14 * time.Hour).Format(time.RFC3339),
			Method:        models.InterviewMethodOnsite,
			Interviewer:   "前端负责人-李工",
			Location:      "北京市朝阳区望京SOHO T3 15层",
			Description:   "技术一面，考察 Vue 生态和前端工程化经验",
			Status:        models.InterviewStatusScheduled,
			LatestNote:    "已安排会议室，等待候选人确认",
			CreatedAt:     now.AddDate(0, 0, -0).Format(time.RFC3339),
			UpdatedAt:     now.AddDate(0, 0, -0).Format(time.RFC3339),
		},
		{
			ID:            uuid.New().String(),
			ApplicationID: applications[2].ID,
			JobID:         jobs[0].ID,
			JobTitle:      jobs[0].Title,
			CandidateName: "王五",
			Contact:       "13800138003",
			Round:         1,
			InterviewTime: now.AddDate(0, 0, -1).Add(15 * time.Hour).Format(time.RFC3339),
			Method:        models.InterviewMethodPhone,
			Interviewer:   "招聘专员-刘经理",
			Location:      "电话初筛",
			Description:   "简历初筛，工作经验与岗位要求不匹配",
			Status:        models.InterviewStatusCancelled,
			Feedback:      "候选人工作经验2年，岗位要求5年以上，不符合要求",
			CreatedAt:     now.AddDate(0, 0, -5).Format(time.RFC3339),
			UpdatedAt:     now.AddDate(0, 0, -4).Format(time.RFC3339),
		},
	}

	offers = []models.Offer{
		{
			ID:            uuid.New().String(),
			ApplicationID: applications[0].ID,
			JobID:         jobs[0].ID,
			JobTitle:      jobs[0].Title,
			CandidateName: "张三",
			Contact:       "13800138001",
			SalaryMin:     28000,
			SalaryMax:     35000,
			StartDate:     now.AddDate(0, 0, 15).Format("2006-01-02"),
			Owner:         "招聘专员-刘经理",
			Description:   "高级后端工程师 Offer，含年终 3 个月",
			Status:        models.OfferStatusPending,
			CreatedAt:     now.AddDate(0, 0, -1).Format(time.RFC3339),
			UpdatedAt:     now.AddDate(0, 0, -1).Format(time.RFC3339),
		},
	}

	yearMonth := now.Format("200601")
	reqSeqNo = 0
	generateReqNo := func() string {
		reqSeqNo++
		return fmt.Sprintf("REQ-%s-%04d", yearMonth, reqSeqNo)
	}

	nowStr := now.Format(time.RFC3339)
	sevenDaysAgo := now.AddDate(0, 0, -7).Format(time.RFC3339)
	fiveDaysAgo := now.AddDate(0, 0, -5).Format(time.RFC3339)
	fourDaysAgo := now.AddDate(0, 0, -4).Format(time.RFC3339)
	threeDaysAgo := now.AddDate(0, 0, -3).Format(time.RFC3339)
	twoDaysAgo := now.AddDate(0, 0, -2).Format(time.RFC3339)
	oneDayAgo := now.AddDate(0, 0, -1).Format(time.RFC3339)

	tenDaysLater := now.AddDate(0, 0, 10).Format("2006-01-02")
	twentyDaysLater := now.AddDate(0, 0, 20).Format("2006-01-02")
	thirtyDaysLater := now.AddDate(0, 0, 30).Format("2006-01-02")
	fortyFiveDaysLater := now.AddDate(0, 0, 45).Format("2006-01-02")
	sevenDaysLater := now.AddDate(0, 0, 7).Format("2006-01-02")
	fifteenDaysLater := now.AddDate(0, 0, 15).Format("2006-01-02")

	approvedNodes := []models.RequirementApprovalNode{
		{
			ID:        uuid.New().String(),
			NodeName:  "部门总监审核",
			Approver:  "技术总监-王工",
			Status:    models.ReqStatusApproved,
			Opinion:   "业务发展需要，同意招聘",
			HandledAt: fiveDaysAgo,
			CreatedAt: sevenDaysAgo,
		},
		{
			ID:        uuid.New().String(),
			NodeName:  "招聘负责人审核",
			Approver:  "招聘经理-李总",
			Status:    models.ReqStatusApproved,
			Opinion:   "编制内，同意启动招聘",
			HandledAt: threeDaysAgo,
			CreatedAt: sevenDaysAgo,
		},
	}

	pendingNodes := []models.RequirementApprovalNode{
		{
			ID:        uuid.New().String(),
			NodeName:  "部门总监审核",
			Approver:  "产品总监-陈总",
			Status:    models.ReqStatusApproved,
			Opinion:   "团队扩张需要，同意",
			HandledAt: twoDaysAgo,
			CreatedAt: threeDaysAgo,
		},
		{
			ID:        uuid.New().String(),
			NodeName:  "招聘负责人审核",
			Approver:  "招聘经理-李总",
			Status:    models.ReqStatusPending,
			CreatedAt: threeDaysAgo,
		},
	}

	returnedNodes := []models.RequirementApprovalNode{
		{
			ID:        uuid.New().String(),
			NodeName:  "部门总监审核",
			Approver:  "设计总监-周总",
			Status:    models.ReqStatusReturned,
			Opinion:   "请补充岗位详细职责和任职要求",
			HandledAt: oneDayAgo,
			CreatedAt: twoDaysAgo,
		},
		{
			ID:        uuid.New().String(),
			NodeName:  "招聘负责人审核",
			Approver:  "招聘经理-李总",
			Status:    models.ReqStatusDraft,
			CreatedAt: twoDaysAgo,
		},
	}

	rejectedNodes := []models.RequirementApprovalNode{
		{
			ID:        uuid.New().String(),
			NodeName:  "部门总监审核",
			Approver:  "市场总监-孙总",
			Status:    models.ReqStatusApproved,
			Opinion:   "同意",
			HandledAt: fiveDaysAgo,
			CreatedAt: sevenDaysAgo,
		},
		{
			ID:        uuid.New().String(),
			NodeName:  "招聘负责人审核",
			Approver:  "招聘经理-李总",
			Status:    models.ReqStatusRejected,
			Opinion:   "当前季度编制已满，建议下个季度再申请",
			HandledAt: threeDaysAgo,
			CreatedAt: sevenDaysAgo,
		},
	}

	convertedNodes := []models.RequirementApprovalNode{
		{
			ID:        uuid.New().String(),
			NodeName:  "部门总监审核",
			Approver:  "运营总监-吴总",
			Status:    models.ReqStatusApproved,
			Opinion:   "急需人手，同意",
			HandledAt: fiveDaysAgo,
			CreatedAt: sevenDaysAgo,
		},
		{
			ID:        uuid.New().String(),
			NodeName:  "招聘负责人审核",
			Approver:  "招聘经理-李总",
			Status:    models.ReqStatusApproved,
			Opinion:   "同意，尽快发布",
			HandledAt: fourDaysAgo,
			CreatedAt: sevenDaysAgo,
		},
	}

	requirements = []models.HiringRequirement{
		{
			ID:               uuid.New().String(),
			ReqNo:            generateReqNo(),
			Department:       "技术部",
			JobTitle:         "资深Go开发工程师",
			HeadcountType:    models.HeadcountRegular,
			Headcount:        2,
			ExpectedOnboard:  tenDaysLater,
			SalaryMin:        30000,
			SalaryMax:        50000,
			WorkLocation:     "北京",
			JobDescription:   "负责核心支付系统的设计与开发，参与架构优化",
			Requirements:     "5年以上Go开发经验，有分布式系统、高并发经验",
			Reason:           "支付业务线扩张，现有人员不足",
			Status:           models.ReqStatusDraft,
			Initiator:        "技术经理-赵工",
			InitiatorContact: "13900139001",
			ApprovalNodes:    nil,
			OperationLogs: []models.RequirementOperationLog{
				{
					ID:        uuid.New().String(),
					Action:    "创建",
					Operator:  "技术经理-赵工",
					Detail:    "创建用人需求草稿",
					CreatedAt: nowStr,
				},
			},
			CreatedAt: nowStr,
			UpdatedAt: nowStr,
		},
		{
			ID:               uuid.New().String(),
			ReqNo:            generateReqNo(),
			Department:       "产品部",
			JobTitle:         "高级产品经理",
			HeadcountType:    models.HeadcountRegular,
			Headcount:        1,
			ExpectedOnboard:  twentyDaysLater,
			SalaryMin:        25000,
			SalaryMax:        45000,
			WorkLocation:     "深圳",
			JobDescription:   "负责C端产品线的规划与迭代，推动产品落地",
			Requirements:     "5年以上C端产品经验，有电商/社交产品经验优先",
			Reason:           "新产品线立项",
			Status:           models.ReqStatusPending,
			Initiator:        "产品经理-林经理",
			InitiatorContact: "13900139002",
			ApprovalNodes:    pendingNodes,
			OperationLogs: []models.RequirementOperationLog{
				{
					ID:        uuid.New().String(),
					Action:    "创建",
					Operator:  "产品经理-林经理",
					Detail:    "创建用人需求",
					CreatedAt: threeDaysAgo,
				},
				{
					ID:        uuid.New().String(),
					Action:    "提交审批",
					Operator:  "产品经理-林经理",
					Detail:    "提交部门总监和招聘负责人审批",
					CreatedAt: threeDaysAgo,
				},
				{
					ID:        uuid.New().String(),
					Action:    "通过",
					Operator:  "产品总监-陈总",
					Detail:    "部门总监审核通过",
					CreatedAt: twoDaysAgo,
				},
			},
			LatestOpinion: "等待招聘负责人审核",
			CreatedAt:     threeDaysAgo,
			UpdatedAt:     twoDaysAgo,
		},
		{
			ID:               uuid.New().String(),
			ReqNo:            generateReqNo(),
			Department:       "技术部",
			JobTitle:         "数据平台架构师",
			HeadcountType:    models.HeadcountRegular,
			Headcount:        1,
			ExpectedOnboard:  thirtyDaysLater,
			SalaryMin:        40000,
			SalaryMax:        65000,
			WorkLocation:     "杭州",
			JobDescription:   "负责大数据平台整体架构设计，带领团队完成技术攻关",
			Requirements:     "8年以上经验，精通Hadoop/Spark/Flink，有PB级数据处理经验",
			Reason:           "数据平台升级改造",
			Status:           models.ReqStatusApproved,
			Initiator:        "数据负责人-郑工",
			InitiatorContact: "13900139003",
			ApprovalNodes:    approvedNodes,
			OperationLogs: []models.RequirementOperationLog{
				{
					ID:        uuid.New().String(),
					Action:    "创建",
					Operator:  "数据负责人-郑工",
					Detail:    "创建用人需求",
					CreatedAt: sevenDaysAgo,
				},
				{
					ID:        uuid.New().String(),
					Action:    "提交审批",
					Operator:  "数据负责人-郑工",
					Detail:    "提交部门总监和招聘负责人审批",
					CreatedAt: sevenDaysAgo,
				},
				{
					ID:        uuid.New().String(),
					Action:    "通过",
					Operator:  "技术总监-王工",
					Detail:    "部门总监审核通过",
					CreatedAt: fiveDaysAgo,
				},
				{
					ID:        uuid.New().String(),
					Action:    "通过",
					Operator:  "招聘经理-李总",
					Detail:    "招聘负责人审核通过",
					CreatedAt: threeDaysAgo,
				},
			},
			LatestOpinion: "审批通过，等待发布职位",
			CreatedAt:     sevenDaysAgo,
			UpdatedAt:     threeDaysAgo,
		},
		{
			ID:               uuid.New().String(),
			ReqNo:            generateReqNo(),
			Department:       "设计部",
			JobTitle:         "资深视觉设计师",
			HeadcountType:    models.HeadcountRegular,
			Headcount:        1,
			ExpectedOnboard:  fortyFiveDaysLater,
			SalaryMin:        20000,
			SalaryMax:        35000,
			WorkLocation:     "北京",
			JobDescription:   "负责品牌视觉设计，参与活动页面设计",
			Requirements:     "待补充",
			Status:           models.ReqStatusReturned,
			Initiator:        "设计师-黄工",
			InitiatorContact: "13900139004",
			ApprovalNodes:    returnedNodes,
			OperationLogs: []models.RequirementOperationLog{
				{
					ID:        uuid.New().String(),
					Action:    "创建",
					Operator:  "设计师-黄工",
					Detail:    "创建用人需求",
					CreatedAt: twoDaysAgo,
				},
				{
					ID:        uuid.New().String(),
					Action:    "提交审批",
					Operator:  "设计师-黄工",
					Detail:    "提交部门总监和招聘负责人审批",
					CreatedAt: twoDaysAgo,
				},
				{
					ID:        uuid.New().String(),
					Action:    "退回",
					Operator:  "设计总监-周总",
					Detail:    "请补充岗位详细职责和任职要求",
					CreatedAt: oneDayAgo,
				},
			},
			LatestOpinion: "请补充岗位详细职责和任职要求",
			CreatedAt:     twoDaysAgo,
			UpdatedAt:     oneDayAgo,
		},
		{
			ID:               uuid.New().String(),
			ReqNo:            generateReqNo(),
			Department:       "市场部",
			JobTitle:         "市场推广专员",
			HeadcountType:    models.HeadcountContract,
			Headcount:        3,
			ExpectedOnboard:  sevenDaysLater,
			SalaryMin:        8000,
			SalaryMax:        15000,
			WorkLocation:     "上海",
			JobDescription:   "负责线上线下推广活动执行",
			Requirements:     "1年以上市场推广经验",
			Reason:           "618大促活动需要临时人员",
			Status:           models.ReqStatusRejected,
			Initiator:        "市场经理-徐经理",
			InitiatorContact: "13900139005",
			ApprovalNodes:    rejectedNodes,
			OperationLogs: []models.RequirementOperationLog{
				{
					ID:        uuid.New().String(),
					Action:    "创建",
					Operator:  "市场经理-徐经理",
					Detail:    "创建用人需求",
					CreatedAt: sevenDaysAgo,
				},
				{
					ID:        uuid.New().String(),
					Action:    "提交审批",
					Operator:  "市场经理-徐经理",
					Detail:    "提交部门总监和招聘负责人审批",
					CreatedAt: sevenDaysAgo,
				},
				{
					ID:        uuid.New().String(),
					Action:    "通过",
					Operator:  "市场总监-孙总",
					Detail:    "部门总监审核通过",
					CreatedAt: fiveDaysAgo,
				},
				{
					ID:        uuid.New().String(),
					Action:    "拒绝",
					Operator:  "招聘经理-李总",
					Detail:    "当前季度编制已满，建议下个季度再申请",
					CreatedAt: threeDaysAgo,
				},
			},
			LatestOpinion: "当前季度编制已满，建议下个季度再申请",
			CreatedAt:     sevenDaysAgo,
			UpdatedAt:     threeDaysAgo,
		},
		{
			ID:               uuid.New().String(),
			ReqNo:            generateReqNo(),
			Department:       "运营部",
			JobTitle:         "用户运营主管",
			HeadcountType:    models.HeadcountReplacement,
			Headcount:        1,
			ExpectedOnboard:  fifteenDaysLater,
			SalaryMin:        18000,
			SalaryMax:        30000,
			WorkLocation:     "广州",
			JobDescription:   "负责用户增长体系搭建，提升用户活跃度和留存率",
			Requirements:     "3年以上用户运营经验，有数据驱动增长经验",
			Reason:           "原运营主管离职，需补位",
			Status:           models.ReqStatusConverted,
			Initiator:        "运营经理-何经理",
			InitiatorContact: "13900139006",
			ApprovalNodes:    convertedNodes,
			RelatedJobID:     jobs[0].ID,
			OperationLogs: []models.RequirementOperationLog{
				{
					ID:        uuid.New().String(),
					Action:    "创建",
					Operator:  "运营经理-何经理",
					Detail:    "创建用人需求",
					CreatedAt: sevenDaysAgo,
				},
				{
					ID:        uuid.New().String(),
					Action:    "提交审批",
					Operator:  "运营经理-何经理",
					Detail:    "提交部门总监和招聘负责人审批",
					CreatedAt: sevenDaysAgo,
				},
				{
					ID:        uuid.New().String(),
					Action:    "通过",
					Operator:  "运营总监-吴总",
					Detail:    "部门总监审核通过",
					CreatedAt: fiveDaysAgo,
				},
				{
					ID:        uuid.New().String(),
					Action:    "通过",
					Operator:  "招聘经理-李总",
					Detail:    "招聘负责人审核通过",
					CreatedAt: fourDaysAgo,
				},
				{
					ID:        uuid.New().String(),
					Action:    "转职位",
					Operator:  "招聘专员-刘经理",
					Detail:    "已生成招聘职位",
					CreatedAt: threeDaysAgo,
				},
			},
			LatestOpinion: "已生成招聘职位",
			CreatedAt:     sevenDaysAgo,
			UpdatedAt:     threeDaysAgo,
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
		InterviewsByStatus:   make(map[string]int),
		OffersByStatus:       make(map[string]int),
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

	for _, iv := range interviews {
		stats.TotalInterviews++
		stats.InterviewsByStatus[string(iv.Status)]++
		if iv.Status == models.InterviewStatusScheduled {
			stats.PendingInterviews++
		}
	}

	for _, offer := range offers {
		stats.TotalOffers++
		stats.OffersByStatus[string(offer.Status)]++
		if offer.Status == models.OfferStatusPending || offer.Status == models.OfferStatusSent {
			stats.PendingOffers++
		}
		if offer.Status == models.OfferStatusAccepted {
			stats.AcceptedOffers++
		}
	}

	reqStats := models.RequirementStats{
		ByDepartment: make(map[string]int),
		ByMonth:      make(map[string]int),
	}
	today := time.Now()
	tenDaysLater := today.AddDate(0, 0, 10)
	for _, req := range requirements {
		switch req.Status {
		case models.ReqStatusPending:
			reqStats.PendingCount++
		case models.ReqStatusApproved:
			reqStats.ApprovedCount++
		case models.ReqStatusReturned:
			reqStats.ReturnedCount++
		case models.ReqStatusRejected:
			reqStats.RejectedCount++
		case models.ReqStatusConverted:
			reqStats.ConvertedCount++
		case models.ReqStatusDraft:
			reqStats.DraftCount++
		}
		reqStats.ByDepartment[req.Department]++
		if req.CreatedAt != "" {
			ct, err := time.Parse(time.RFC3339, req.CreatedAt)
			if err == nil {
				monthKey := ct.Format("2006-01")
				reqStats.ByMonth[monthKey]++
			}
		}
		if req.ExpectedOnboard != "" {
			et, err := time.Parse("2006-01-02", req.ExpectedOnboard)
			if err == nil {
				if (et.Equal(today) || et.After(today)) && et.Before(tenDaysLater) {
					reqStats.UrgentCount++
				}
			}
		}
	}
	stats.RequirementStats = reqStats

	return stats
}

func ListInterviews(jobId, status, method string) []models.Interview {
	mu.RLock()
	defer mu.RUnlock()

	result := make([]models.Interview, 0)
	for _, iv := range interviews {
		if jobId != "" && iv.JobID != jobId {
			continue
		}
		if status != "" && string(iv.Status) != status {
			continue
		}
		if method != "" && string(iv.Method) != method {
			continue
		}
		result = append(result, iv)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].InterviewTime > result[j].InterviewTime
	})
	return result
}

func GetInterview(id string) *models.Interview {
	mu.RLock()
	defer mu.RUnlock()

	for i := range interviews {
		if interviews[i].ID == id {
			return &interviews[i]
		}
	}
	return nil
}

func CreateInterview(iv *models.Interview) *models.Interview {
	mu.Lock()
	defer mu.Unlock()

	var app *models.Application
	for i := range applications {
		if applications[i].ID == iv.ApplicationID {
			app = &applications[i]
			break
		}
	}
	if app == nil {
		return nil
	}

	iv.ID = uuid.New().String()
	iv.JobID = app.JobID
	iv.JobTitle = app.JobTitle
	iv.CandidateName = app.Resume.CandidateName
	iv.Contact = app.Resume.Contact
	iv.Status = models.InterviewStatusScheduled
	now := time.Now().Format(time.RFC3339)
	iv.CreatedAt = now
	iv.UpdatedAt = now

	interviews = append(interviews, *iv)

	for i := range applications {
		if applications[i].ID == iv.ApplicationID {
			applications[i].Status = models.AppStatusInterview
			break
		}
	}

	return &interviews[len(interviews)-1]
}

func RescheduleInterview(id, newTime string) *models.Interview {
	mu.Lock()
	defer mu.Unlock()

	for i := range interviews {
		if interviews[i].ID == id {
			interviews[i].InterviewTime = newTime
			interviews[i].UpdatedAt = time.Now().Format(time.RFC3339)
			return &interviews[i]
		}
	}
	return nil
}

func CancelInterview(id string) *models.Interview {
	mu.Lock()
	defer mu.Unlock()

	for i := range interviews {
		if interviews[i].ID == id {
			interviews[i].Status = models.InterviewStatusCancelled
			interviews[i].UpdatedAt = time.Now().Format(time.RFC3339)

			for j := range applications {
				if applications[j].ID == interviews[i].ApplicationID {
					applications[j].Status = models.AppStatusPending
					break
				}
			}
			return &interviews[i]
		}
	}
	return nil
}

type CompleteInterviewRequest struct {
	Feedback   string
	Note       string
	Conclusion string
	Rating     int
	Strengths  string
	Risks      string
	NextSteps  string
}

func CompleteInterview(id string, req CompleteInterviewRequest) *models.Interview {
	mu.Lock()
	defer mu.Unlock()

	for i := range interviews {
		if interviews[i].ID == id {
			interviews[i].Status = models.InterviewStatusCompleted
			interviews[i].Feedback = req.Feedback
			if req.Conclusion != "" {
				interviews[i].Conclusion = models.InterviewFeedbackConclusion(req.Conclusion)
			}
			if req.Rating > 0 {
				interviews[i].Rating = req.Rating
			}
			interviews[i].Strengths = req.Strengths
			interviews[i].Risks = req.Risks
			interviews[i].NextSteps = req.NextSteps
			if req.Note != "" {
				interviews[i].LatestNote = req.Note
			}
			interviews[i].UpdatedAt = time.Now().Format(time.RFC3339)
			return &interviews[i]
		}
	}
	return nil
}

func AddInterviewNote(id, note string) *models.Interview {
	mu.Lock()
	defer mu.Unlock()

	for i := range interviews {
		if interviews[i].ID == id {
			interviews[i].LatestNote = note
			interviews[i].UpdatedAt = time.Now().Format(time.RFC3339)
			return &interviews[i]
		}
	}
	return nil
}

func ListOffers(jobId, status string) []models.Offer {
	mu.RLock()
	defer mu.RUnlock()

	result := make([]models.Offer, 0)
	for _, offer := range offers {
		if jobId != "" && offer.JobID != jobId {
			continue
		}
		if status != "" && string(offer.Status) != status {
			continue
		}
		result = append(result, offer)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt > result[j].CreatedAt
	})
	return result
}

func GetOffer(id string) *models.Offer {
	mu.RLock()
	defer mu.RUnlock()

	for i := range offers {
		if offers[i].ID == id {
			return &offers[i]
		}
	}
	return nil
}

func CreateOffer(offer *models.Offer) *models.Offer {
	mu.Lock()
	defer mu.Unlock()

	var app *models.Application
	for i := range applications {
		if applications[i].ID == offer.ApplicationID {
			app = &applications[i]
			break
		}
	}
	if app == nil {
		return nil
	}

	offer.ID = uuid.New().String()
	offer.JobID = app.JobID
	offer.JobTitle = app.JobTitle
	offer.CandidateName = app.Resume.CandidateName
	offer.Contact = app.Resume.Contact
	offer.Status = models.OfferStatusPending
	now := time.Now().Format(time.RFC3339)
	offer.CreatedAt = now
	offer.UpdatedAt = now

	offers = append(offers, *offer)
	return &offers[len(offers)-1]
}

func SendOffer(id string) *models.Offer {
	mu.Lock()
	defer mu.Unlock()

	now := time.Now().Format(time.RFC3339)
	for i := range offers {
		if offers[i].ID == id {
			offers[i].Status = models.OfferStatusSent
			offers[i].SentAt = now
			offers[i].UpdatedAt = now
			return &offers[i]
		}
	}
	return nil
}

type OfferReplyRequest struct {
	Accepted bool
	Feedback string
}

func ReplyOffer(id string, req OfferReplyRequest) *models.Offer {
	mu.Lock()
	defer mu.Unlock()

	now := time.Now().Format(time.RFC3339)
	for i := range offers {
		if offers[i].ID == id {
			if req.Accepted {
				offers[i].Status = models.OfferStatusAccepted
			} else {
				offers[i].Status = models.OfferStatusRejected
			}
			offers[i].Feedback = req.Feedback
			offers[i].RepliedAt = now
			offers[i].UpdatedAt = now

			for j := range applications {
				if applications[j].ID == offers[i].ApplicationID {
					if req.Accepted {
						applications[j].Status = models.AppStatusHired
					}
					break
				}
			}
			return &offers[i]
		}
	}
	return nil
}

func WithdrawOffer(id, reason string) *models.Offer {
	mu.Lock()
	defer mu.Unlock()

	now := time.Now().Format(time.RFC3339)
	for i := range offers {
		if offers[i].ID == id {
			offers[i].Status = models.OfferStatusWithdrawn
			offers[i].Feedback = reason
			offers[i].UpdatedAt = now
			return &offers[i]
		}
	}
	return nil
}

func UpdateOfferNote(id, note string) *models.Offer {
	mu.Lock()
	defer mu.Unlock()

	for i := range offers {
		if offers[i].ID == id {
			offers[i].Description = note
			offers[i].UpdatedAt = time.Now().Format(time.RFC3339)
			return &offers[i]
		}
	}
	return nil
}

func ListRequirements(department, status, keyword string) []models.HiringRequirement {
	mu.RLock()
	defer mu.RUnlock()

	result := make([]models.HiringRequirement, 0)
	for _, req := range requirements {
		if department != "" && req.Department != department {
			continue
		}
		if status != "" && string(req.Status) != status {
			continue
		}
		if keyword != "" {
			kw := strings.ToLower(keyword)
			if !strings.Contains(strings.ToLower(req.JobTitle), kw) &&
				!strings.Contains(strings.ToLower(req.Department), kw) &&
				!strings.Contains(strings.ToLower(req.ReqNo), kw) &&
				!strings.Contains(strings.ToLower(req.Initiator), kw) &&
				!strings.Contains(strings.ToLower(req.JobDescription), kw) {
				continue
			}
		}
		result = append(result, req)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].UpdatedAt > result[j].UpdatedAt
	})
	return result
}

func GetRequirement(id string) *models.HiringRequirement {
	mu.RLock()
	defer mu.RUnlock()

	for i := range requirements {
		if requirements[i].ID == id {
			return &requirements[i]
		}
	}
	return nil
}

func CreateRequirement(req *models.HiringRequirement) *models.HiringRequirement {
	mu.Lock()
	defer mu.Unlock()

	now := time.Now()
	yearMonth := now.Format("200601")
	reqSeqNo++
	req.ID = uuid.New().String()
	req.ReqNo = fmt.Sprintf("REQ-%s-%04d", yearMonth, reqSeqNo)
	req.Status = models.ReqStatusDraft
	req.ApprovalNodes = []models.RequirementApprovalNode{}
	req.OperationLogs = []models.RequirementOperationLog{
		{
			ID:        uuid.New().String(),
			Action:    "创建",
			Operator:  req.Initiator,
			Detail:    "创建用人需求草稿",
			CreatedAt: now.Format(time.RFC3339),
		},
	}
	nowStr := now.Format(time.RFC3339)
	req.CreatedAt = nowStr
	req.UpdatedAt = nowStr

	requirements = append(requirements, *req)
	return &requirements[len(requirements)-1]
}

func UpdateRequirement(id string, req *models.HiringRequirement) *models.HiringRequirement {
	mu.Lock()
	defer mu.Unlock()

	for i := range requirements {
		if requirements[i].ID == id {
			if requirements[i].Status != models.ReqStatusDraft && requirements[i].Status != models.ReqStatusReturned {
				return nil
			}
			nowStr := time.Now().Format(time.RFC3339)
			requirements[i].Department = req.Department
			requirements[i].JobTitle = req.JobTitle
			requirements[i].HeadcountType = req.HeadcountType
			requirements[i].Headcount = req.Headcount
			requirements[i].ExpectedOnboard = req.ExpectedOnboard
			requirements[i].SalaryMin = req.SalaryMin
			requirements[i].SalaryMax = req.SalaryMax
			requirements[i].WorkLocation = req.WorkLocation
			requirements[i].JobDescription = req.JobDescription
			requirements[i].Requirements = req.Requirements
			requirements[i].Reason = req.Reason
			if req.Initiator != "" {
				requirements[i].Initiator = req.Initiator
			}
			if req.InitiatorContact != "" {
				requirements[i].InitiatorContact = req.InitiatorContact
			}
			requirements[i].UpdatedAt = nowStr
			requirements[i].OperationLogs = append(requirements[i].OperationLogs, models.RequirementOperationLog{
				ID:        uuid.New().String(),
				Action:    "编辑",
				Operator:  requirements[i].Initiator,
				Detail:    "编辑用人需求内容",
				CreatedAt: nowStr,
			})
			return &requirements[i]
		}
	}
	return nil
}

func SubmitApproval(id string) *models.HiringRequirement {
	mu.Lock()
	defer mu.Unlock()

	nowStr := time.Now().Format(time.RFC3339)
	for i := range requirements {
		if requirements[i].ID == id {
			if requirements[i].Status != models.ReqStatusDraft && requirements[i].Status != models.ReqStatusReturned {
				return nil
			}
			requirements[i].Status = models.ReqStatusPending
			requirements[i].ApprovalNodes = []models.RequirementApprovalNode{
				{
					ID:        uuid.New().String(),
					NodeName:  "部门总监审核",
					Approver:  "部门总监",
					Status:    models.ReqStatusPending,
					CreatedAt: nowStr,
				},
				{
					ID:        uuid.New().String(),
					NodeName:  "招聘负责人审核",
					Approver:  "招聘负责人",
					Status:    models.ReqStatusPending,
					CreatedAt: nowStr,
				},
			}
			requirements[i].UpdatedAt = nowStr
			requirements[i].LatestOpinion = ""
			requirements[i].OperationLogs = append(requirements[i].OperationLogs, models.RequirementOperationLog{
				ID:        uuid.New().String(),
				Action:    "提交审批",
				Operator:  requirements[i].Initiator,
				Detail:    "提交部门总监和招聘负责人审批",
				CreatedAt: nowStr,
			})
			return &requirements[i]
		}
	}
	return nil
}

func ApproveRequirement(id string, nodeIndex int, action string, opinion string) *models.HiringRequirement {
	mu.Lock()
	defer mu.Unlock()

	nowStr := time.Now().Format(time.RFC3339)
	for i := range requirements {
		if requirements[i].ID == id {
			if requirements[i].Status != models.ReqStatusPending {
				return nil
			}
			if nodeIndex < 0 || nodeIndex >= len(requirements[i].ApprovalNodes) {
				return nil
			}

			node := &requirements[i].ApprovalNodes[nodeIndex]
			if node.Status != models.ReqStatusPending {
				return nil
			}

			node.Opinion = opinion
			node.HandledAt = nowStr

			switch action {
			case "approve":
				node.Status = models.ReqStatusApproved
				isLast := nodeIndex == len(requirements[i].ApprovalNodes)-1
				allApproved := true
				for _, n := range requirements[i].ApprovalNodes {
					if n.Status != models.ReqStatusApproved {
						allApproved = false
						break
					}
				}
				if isLast && allApproved {
					requirements[i].Status = models.ReqStatusApproved
				}
				requirements[i].LatestOpinion = opinion
				logAction := "通过"
				logDetail := node.NodeName + "审核通过"
				if opinion != "" {
					logDetail += "：" + opinion
				}
				requirements[i].OperationLogs = append(requirements[i].OperationLogs, models.RequirementOperationLog{
					ID:        uuid.New().String(),
					Action:    logAction,
					Operator:  node.Approver,
					Detail:    logDetail,
					CreatedAt: nowStr,
				})
			case "return":
				node.Status = models.ReqStatusReturned
				requirements[i].Status = models.ReqStatusReturned
				requirements[i].LatestOpinion = opinion
				for j := nodeIndex + 1; j < len(requirements[i].ApprovalNodes); j++ {
					requirements[i].ApprovalNodes[j].Status = models.ReqStatusDraft
				}
				logDetail := node.NodeName + "退回修改"
				if opinion != "" {
					logDetail += "：" + opinion
				}
				requirements[i].OperationLogs = append(requirements[i].OperationLogs, models.RequirementOperationLog{
					ID:        uuid.New().String(),
					Action:    "退回",
					Operator:  node.Approver,
					Detail:    logDetail,
					CreatedAt: nowStr,
				})
			case "reject":
				node.Status = models.ReqStatusRejected
				requirements[i].Status = models.ReqStatusRejected
				requirements[i].LatestOpinion = opinion
				logDetail := node.NodeName + "审核拒绝"
				if opinion != "" {
					logDetail += "：" + opinion
				}
				requirements[i].OperationLogs = append(requirements[i].OperationLogs, models.RequirementOperationLog{
					ID:        uuid.New().String(),
					Action:    "拒绝",
					Operator:  node.Approver,
					Detail:    logDetail,
					CreatedAt: nowStr,
				})
			default:
				return nil
			}

			requirements[i].UpdatedAt = nowStr
			return &requirements[i]
		}
	}
	return nil
}

func ConvertToJob(id string) *models.Job {
	mu.Lock()
	defer mu.Unlock()

	nowStr := time.Now().Format(time.RFC3339)
	for i := range requirements {
		if requirements[i].ID == id {
			if requirements[i].Status != models.ReqStatusApproved {
				return nil
			}
			if requirements[i].RelatedJobID != "" {
				return nil
			}

			jobID := uuid.New().String()
			publishedAt := time.Now().Format("2006-01-02")
			newJob := models.Job{
				ID:            jobID,
				Title:         requirements[i].JobTitle,
				Department:    requirements[i].Department,
				Location:      requirements[i].WorkLocation,
				SalaryMin:     requirements[i].SalaryMin,
				SalaryMax:     requirements[i].SalaryMax,
				Status:        models.JobStatusOpen,
				Headcount:     requirements[i].Headcount,
				Description:   requirements[i].JobDescription,
				Requirements:  requirements[i].Requirements,
				PublishedAt:   publishedAt,
				RequirementId: id,
			}
			jobs = append(jobs, newJob)

			requirements[i].RelatedJobID = jobID
			requirements[i].Status = models.ReqStatusConverted
			requirements[i].UpdatedAt = nowStr
			requirements[i].OperationLogs = append(requirements[i].OperationLogs, models.RequirementOperationLog{
				ID:        uuid.New().String(),
				Action:    "转职位",
				Operator:  "招聘专员-刘经理",
				Detail:    "已生成招聘职位：" + newJob.Title,
				CreatedAt: nowStr,
			})

			return &jobs[len(jobs)-1]
		}
	}
	return nil
}

func GetRequirementStats() models.RequirementStats {
	mu.RLock()
	defer mu.RUnlock()

	stats := models.RequirementStats{
		ByDepartment: make(map[string]int),
		ByMonth:      make(map[string]int),
	}

	now := time.Now()
	urgentThreshold := now.AddDate(0, 0, 14)

	for _, req := range requirements {
		switch req.Status {
		case models.ReqStatusDraft:
			stats.DraftCount++
		case models.ReqStatusPending:
			stats.PendingCount++
		case models.ReqStatusApproved:
			stats.ApprovedCount++
		case models.ReqStatusReturned:
			stats.ReturnedCount++
		case models.ReqStatusRejected:
			stats.RejectedCount++
		case models.ReqStatusConverted:
			stats.ConvertedCount++
		}

		stats.ByDepartment[req.Department]++

		if req.CreatedAt != "" {
			createdTime, err := time.Parse(time.RFC3339, req.CreatedAt)
			if err == nil {
				monthKey := createdTime.Format("2006-01")
				stats.ByMonth[monthKey]++
			}
		}

		if req.Status != models.ReqStatusRejected && req.Status != models.ReqStatusConverted {
			if req.ExpectedOnboard != "" {
				onboardTime, err := time.Parse("2006-01-02", req.ExpectedOnboard)
				if err == nil {
					if onboardTime.Before(urgentThreshold) || onboardTime.Equal(urgentThreshold) {
						stats.UrgentCount++
					}
				}
			}
		}
	}

	return stats
}
