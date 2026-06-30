# 招聘平台 - 前后端分离原型

基于 Go + Gin + Vue 3 + Element Plus 构建的招聘平台可运行原型，使用内存 Mock 数据，不依赖真实数据库。

## 功能模块

### 工作台（5大入口）
- **📊 统计概览**：职位总数、在招职位、累计申请、本周新增、申请状态分布饼图、各职位申请数 Top5 柱状图、最新申请动态、招聘职位速览
- **💼 职位管理**：关键词/地点/状态筛选、新增/编辑职位、职位详情、快捷投递入口
- **👥 候选人管理**：关键词/职位/状态筛选、候选人列表（姓名、联系方式、期望岗位、投递职位、工作年限、技能标签、当前状态、最近沟通、备注数）、候选人详情（简历摘要、投递记录时间线、沟通记录、状态操作、备注区）
- **📝 申请管理**：5状态（已投递/待沟通/面试中/已拒绝/已录用）彩色统计卡片、简历投递、申请详情、状态变更
- **💬 沟通中心**：左侧会话列表、右侧气泡聊天、招聘方/候选人身份切换、消息发送

### 核心特性
- ✅ **数据一致性**：候选人通过申请聚合生成，状态变更同步影响申请列表和统计数字
- ✅ **前端稳定性**：后端健康检查（5秒轮询）、服务断开友好提示页（含排查建议和重试按钮）、服务恢复通知提醒
- ✅ **完整业务流程**：新增职位 → 投递简历 → 候选人筛选 → 状态变更 → 备注记录 → 双方沟通 → 统计变化
- ✅ **5种申请状态**：submitted（已投递）→ pending（待沟通）→ interview（面试中）→ rejected（已拒绝）→ hired（已录用）
- ✅ **备注系统**：支持面试备注、筛选意见两种类型，保存后立即显示在详情时间线

## 技术栈

### 后端
- **语言**：Go 1.21+
- **框架**：Gin v1.9.1
- **跨域**：gin-contrib/cors v1.5.0
- **ID生成**：google/uuid v1.5.0
- **并发安全**：sync.RWMutex
- **数据存储**：内存存储 + Mock 数据

### 前端
- **框架**：Vue 3（Composition API）
- **构建工具**：Vite 5
- **路由**：Vue Router 4
- **UI组件**：Element Plus 2.4
- **图表**：ECharts 5
- **HTTP客户端**：Axios 1.6
- **日期处理**：Day.js 1.11

## 快速开始

### 环境要求
- Go 1.21+
- Node.js 18+

### 启动后端服务

**方式一：使用启动脚本**
```powershell
cd backend
.\start-server.ps1
```

**方式二：手动启动**
```powershell
cd backend
go build -mod=mod -o server.exe .
.\server.exe
```

**方式三：开发模式**
```powershell
cd backend
go run main.go
```

后端服务启动后访问：http://localhost:8080

### 启动前端服务

**方式一：使用启动脚本**
```powershell
cd frontend
.\start-dev.ps1
```

**方式二：手动启动**
```powershell
cd frontend
npm install
npm run dev
```

前端服务启动后访问：http://localhost:5173

### 健康检查
- 后端健康检查：http://localhost:8080/api/health
- 预期返回：`{"status":"ok"}`

## API 接口列表

### 职位管理
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/jobs | 职位列表（支持 keyword/location/status 筛选） |
| GET | /api/jobs/:id | 职位详情 |
| POST | /api/jobs | 新增职位 |
| PUT | /api/jobs/:id | 编辑职位 |

### 申请管理
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/applications | 申请列表（支持 status/jobId 筛选） |
| GET | /api/applications/:id | 申请详情 |
| POST | /api/applications | 提交简历申请 |
| PUT | /api/applications/:id/status | 更新申请状态 |

### 候选人管理
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/candidates | 候选人列表（支持 keyword/jobId/status 筛选） |
| GET | /api/candidates/detail | 候选人详情（contact + name 参数） |
| PUT | /api/candidates/status | 更新候选人所有申请状态 |

### 沟通消息
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/messages/application/:applicationId | 获取申请的消息列表 |
| POST | /api/messages | 发送消息 |

### 备注管理
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/notes | 备注列表（支持 contact/name 筛选） |
| POST | /api/notes | 新增备注 |

### 统计数据
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/stats | 统计面板数据 |
| GET | /api/health | 健康检查 |

## 数据结构说明

### 核心模型关联
```
职位 (Job)
  ↓ 1:N
申请 (Application) → 包含简历信息
  ↓ 1:N
消息 (Message)
  ↓ 1:N
备注 (Note)
```

### 候选人聚合规则
- 候选人不单独存储，通过 `contact + candidateName` 组合键从申请表聚合
- 同一候选人投递多个职位时，在列表页合并展示，详情页显示所有投递记录
- 更新候选人状态时，同步更新该候选人所有申请的状态

### Mock 数据（预置）
- 5 个职位（高级后端、前端、产品经理、UI设计师、数据分析师）
- 5 份申请（覆盖全部5种状态）
- 3 条消息（2条面试沟通、1条拒绝通知）
- 3 条备注（2条张三、1条李四）

## 完整验证路径

### 路径一：候选人管理全流程
1. 首页自动进入「统计概览」，查看整体数据
2. 点击侧边栏「候选人管理」，查看5位候选人列表
3. **筛选验证**：关键词输入「前端」，点击搜索，列表应只显示李四、钱七
4. 点击「重置」恢复全部列表
5. 点击候选人「张三」的「查看详情」
6. **状态变更验证**：点击「待沟通」按钮，确认状态变为待沟通
7. 点击侧边栏「申请管理」，确认张三的申请状态已同步变为「待沟通」
8. 点击侧边栏「统计概览」，确认最新申请动态中张三状态已更新
9. 回到候选人详情页，**备注验证**：输入备注内容，点击「添加备注」，确认新备注显示在时间线顶部
10. **沟通验证**：点击「发起沟通」，在沟通中心发送消息

### 路径二：招聘主流程
1. 「职位管理」→「新增职位」→ 填写信息提交
2. 职位列表中点击「投递」→ 填写简历信息提交
3. 「申请管理」查看新申请 → 改状态为「待沟通」
4. 「沟通中心」与候选人互发消息
5. 「统计概览」查看数据变化

## 目录结构

```
permission-system-0d41021b/
├── backend/                    # Go 后端
│   ├── handlers/               # API 处理器
│   │   └── handlers.go
│   ├── models/                 # 数据模型
│   │   └── models.go
│   ├── store/                  # 内存存储 + Mock 数据
│   │   └── store.go
│   ├── main.go                 # 应用入口
│   ├── go.mod
│   ├── go.sum
│   └── start-server.ps1        # 启动脚本
├── frontend/                   # Vue 3 前端
│   ├── src/
│   │   ├── api/                # API 封装
│   │   │   ├── index.js        # 业务接口
│   │   │   ├── health.js       # 健康检查
│   │   │   └── request.js      # Axios 实例
│   │   ├── components/         # 公共组件
│   │   │   └── ServiceDown.vue # 服务降级提示
│   │   ├── layout/             # 布局组件
│   │   │   └── Layout.vue      # 侧边栏 + 顶栏
│   │   ├── router/             # 路由配置
│   │   │   └── index.js
│   │   ├── views/              # 页面组件
│   │   │   ├── Stats.vue       # 统计概览
│   │   │   ├── Jobs.vue        # 职位管理
│   │   │   ├── Candidates.vue  # 候选人列表
│   │   │   ├── CandidateDetail.vue  # 候选人详情
│   │   │   ├── Applications.vue    # 申请管理
│   │   │   └── Messages.vue    # 沟通中心
│   │   ├── App.vue             # 根组件
│   │   └── main.js             # 入口文件
│   ├── package.json
│   ├── vite.config.js
│   └── start-dev.ps1           # 启动脚本
└── README.md
```

## 后续扩展方向

- 🔐 登录注册、角色权限（招聘方/应聘方/管理员）
- 📄 简历评分、智能匹配候选人
- 💳 实名认证、在线支付、平台担保
- 🤖 接入 LLM（如 DeepSeek）进行智能简历分析和岗位匹配
- 📊 更丰富的数据分析看板和报表导出
- 💾 接入真实数据库（MySQL/PostgreSQL）
- 🐳 Docker 容器化部署
