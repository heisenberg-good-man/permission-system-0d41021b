<template>
  <div class="jobs-page">
    <el-card class="filter-card" shadow="never">
      <el-form :inline="true" :model="filters" class="filter-form">
        <el-form-item label="关键词">
          <el-input
            v-model="filters.keyword"
            placeholder="岗位名称/部门/描述"
            clearable
            style="width: 220px"
            :prefix-icon="Search"
            @keyup.enter="loadJobs"
          />
        </el-form-item>
        <el-form-item label="地点">
          <el-select v-model="filters.location" placeholder="全部" clearable style="width: 140px">
            <el-option v-for="loc in locations" :key="loc" :label="loc" :value="loc" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="filters.status" placeholder="全部" clearable style="width: 140px">
            <el-option label="招聘中" value="open" />
            <el-option label="已暂停" value="paused" />
            <el-option label="已关闭" value="closed" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadJobs" :icon="Search">搜索</el-button>
          <el-button @click="resetFilters">重置</el-button>
          <el-button type="success" @click="openJobDialog()" :icon="Plus">新增职位</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card class="list-card" shadow="never">
      <div class="card-header">
        <h3 class="card-title">职位列表 <span class="count-tip">共 {{ jobList.length }} 条</span></h3>
      </div>
      <el-empty v-if="!loading && jobList.length === 0" description="暂无职位数据，请先新增职位开始招聘">
        <el-button type="primary" @click="loadJobs">重新加载</el-button>
      </el-empty>
      <el-table v-else :data="jobList" v-loading="loading" stripe style="width: 100%">
        <el-table-column prop="title" label="岗位名称" min-width="180">
          <template #default="{ row }">
            <span class="job-title" @click="showJobDetail(row)">{{ row.title }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="department" label="部门" width="120" />
        <el-table-column prop="location" label="地点" width="100" />
        <el-table-column label="薪资范围" width="160">
          <template #default="{ row }">
            <span class="salary-text">¥{{ row.salaryMin }} - ¥{{ row.salaryMax }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="headcount" label="需求人数" width="100" align="center">
          <template #default="{ row }">
            <el-tag type="info" size="small">{{ row.headcount }}人</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="110" align="center">
          <template #default="{ row }">
            <el-tag :type="jobStatusTagType(row.status)" size="small" effect="light">
              {{ jobStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="publishedAt" label="发布时间" width="120" />
        <el-table-column label="操作" width="300" fixed="right" align="center">
          <template #default="{ row }">
            <div class="action-column">
              <div class="action-row-1">
                <el-button link type="primary" size="small" @click="showJobDetail(row)">详情</el-button>
                <el-button link type="primary" size="small" @click="openJobDialog(row)">编辑</el-button>
                <el-button link type="warning" size="small" @click="goApply(row)">投递简历</el-button>
              </div>
              <div class="action-row-2" v-if="row.requirementId">
                <el-tag
                  v-if="reqCache[row.requirementId]"
                  type="primary"
                  size="small"
                  effect="plain"
                  class="req-tag"
                  @click="showRequirementDetail(row.requirementId)"
                >
                  <el-icon><Link /></el-icon>
                  {{ reqCache[row.requirementId].reqNo }}
                </el-tag>
                <el-tag
                  v-else
                  type="info"
                  size="small"
                  effect="plain"
                  class="req-tag loading-tag"
                >
                  <el-icon class="is-loading"><Loading /></el-icon>
                  加载中
                </el-tag>
              </div>
            </div>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog
      v-model="jobDialogVisible"
      :title="editingJob ? '编辑职位' : '新增职位'"
      width="600px"
      destroy-on-close
    >
      <el-form
        ref="jobFormRef"
        :model="jobForm"
        :rules="jobFormRules"
        label-width="100px"
      >
        <el-form-item label="岗位名称" prop="title">
          <el-input v-model="jobForm.title" placeholder="请输入岗位名称" />
        </el-form-item>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="部门" prop="department">
              <el-input v-model="jobForm.department" placeholder="请输入部门" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="地点" prop="location">
              <el-input v-model="jobForm.location" placeholder="请输入工作地点" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="薪资下限" prop="salaryMin">
              <el-input-number v-model="jobForm.salaryMin" :min="0" :step="1000" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="薪资上限" prop="salaryMax">
              <el-input-number v-model="jobForm.salaryMax" :min="0" :step="1000" style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="招聘状态" prop="status">
              <el-select v-model="jobForm.status" style="width: 100%">
                <el-option label="招聘中" value="open" />
                <el-option label="已暂停" value="paused" />
                <el-option label="已关闭" value="closed" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="需求人数" prop="headcount">
              <el-input-number v-model="jobForm.headcount" :min="1" :max="50" style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="职位描述" prop="description">
          <el-input v-model="jobForm.description" type="textarea" :rows="3" placeholder="请输入职位描述" />
        </el-form-item>
        <el-form-item label="任职要求" prop="requirements">
          <el-input v-model="jobForm.requirements" type="textarea" :rows="3" placeholder="请输入任职要求" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="jobDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitJobForm">保存</el-button>
      </template>
    </el-dialog>

    <el-dialog
      v-model="jobDetailVisible"
      title="职位详情"
      width="650px"
      destroy-on-close
    >
      <div v-if="currentJob" class="job-detail">
        <div class="detail-header">
          <h2 class="detail-title">{{ currentJob.title }}</h2>
          <el-tag :type="jobStatusTagType(currentJob.status)" size="large" effect="light">
            {{ jobStatusText(currentJob.status) }}
          </el-tag>
        </div>
        <div class="detail-meta">
          <span><el-icon><OfficeBuilding /></el-icon> {{ currentJob.department }}</span>
          <span><el-icon><LocationFilled /></el-icon> {{ currentJob.location }}</span>
          <span class="salary-highlight">
            <el-icon><Money /></el-icon> ¥{{ currentJob.salaryMin }} - ¥{{ currentJob.salaryMax }}
          </span>
          <span><el-icon><UserFilled /></el-icon> 招{{ currentJob.headcount }}人</span>
          <span><el-icon><Calendar /></el-icon> {{ currentJob.publishedAt }}</span>
        </div>

        <div v-if="currentJob.requirementId" class="detail-requirement">
          <el-divider content-position="left">
            <span class="divider-title"><el-icon><Link /></el-icon> 关联需求</span>
          </el-divider>
          <div v-if="reqLoadingIds[currentJob.requirementId]" class="req-loading">
            <el-icon class="is-loading"><Loading /></el-icon> 需求信息加载中...
          </div>
          <div v-else-if="currentJobReq" class="req-info-card">
            <div class="req-info-row">
              <div class="req-info-item">
                <span class="req-info-label">需求编号</span>
                <el-tag type="primary" size="default">{{ currentJobReq.reqNo }}</el-tag>
              </div>
              <div class="req-info-item">
                <span class="req-info-label">需求部门</span>
                <span class="req-info-value">{{ currentJobReq.department }}</span>
              </div>
              <div class="req-info-item">
                <span class="req-info-label">需求状态</span>
                <el-tag :type="reqStatusTagType(currentJobReq.status)" size="default">
                  {{ reqStatusText(currentJobReq.status) }}
                </el-tag>
              </div>
            </div>
            <div class="req-info-row">
              <div class="req-info-item">
                <span class="req-info-label">招聘人数</span>
                <span class="req-info-value">{{ currentJobReq.headcount }}人</span>
              </div>
              <div class="req-info-item">
                <span class="req-info-label">期望到岗</span>
                <span class="req-info-value">{{ currentJobReq.expectedOnboard }}</span>
              </div>
            </div>
            <div class="req-action-row">
              <el-button type="primary" size="small" @click="openReqDetail(currentJob.requirementId)">
                <el-icon><View /></el-icon> 查看需求详情
              </el-button>
            </div>
          </div>
          <el-empty v-else description="需求信息加载失败" :image-size="60">
            <el-button type="primary" size="small" @click="loadRequirement(currentJob.requirementId)">
              重新加载
            </el-button>
          </el-empty>
        </div>

        <el-divider />
        <div class="detail-section">
          <h4><el-icon><Document /></el-icon> 职位描述</h4>
          <p>{{ currentJob.description }}</p>
        </div>
        <div class="detail-section">
          <h4><el-icon><StarFilled /></el-icon> 任职要求</h4>
          <p>{{ currentJob.requirements }}</p>
        </div>
      </div>
      <template #footer>
        <el-button @click="jobDetailVisible = false">关闭</el-button>
        <el-button type="warning" @click="goApply(currentJob)">投递简历</el-button>
        <el-button type="primary" @click="editCurrentJob">编辑职位</el-button>
      </template>
    </el-dialog>

    <el-dialog
      v-model="reqDetailVisible"
      title="需求详情"
      width="700px"
      destroy-on-close
    >
      <div v-if="currentReq" class="req-detail">
        <div class="req-detail-header">
          <el-tag type="primary" size="large">{{ currentReq.reqNo }}</el-tag>
          <el-tag :type="reqStatusTagType(currentReq.status)" size="large" effect="light">
            {{ reqStatusText(currentReq.status) }}
          </el-tag>
        </div>
        <h3 class="req-detail-title">{{ currentReq.jobTitle }}</h3>
        <div class="req-detail-grid">
          <div class="req-grid-item">
            <span class="req-grid-label">需求部门</span>
            <span class="req-grid-value">{{ currentReq.department }}</span>
          </div>
          <div class="req-grid-item">
            <span class="req-grid-label">工作地点</span>
            <span class="req-grid-value">{{ currentReq.workLocation }}</span>
          </div>
          <div class="req-grid-item">
            <span class="req-grid-label">招聘人数</span>
            <span class="req-grid-value">{{ currentReq.headcount }}人</span>
          </div>
          <div class="req-grid-item">
            <span class="req-grid-label">人员类型</span>
            <span class="req-grid-value">{{ headcountTypeText(currentReq.headcountType) }}</span>
          </div>
          <div class="req-grid-item">
            <span class="req-grid-label">薪资范围</span>
            <span class="req-grid-value salary-highlight">
              ¥{{ currentReq.salaryMin }} - ¥{{ currentReq.salaryMax }}
            </span>
          </div>
          <div class="req-grid-item">
            <span class="req-grid-label">期望到岗</span>
            <span class="req-grid-value">{{ currentReq.expectedOnboard }}</span>
          </div>
          <div class="req-grid-item">
            <span class="req-grid-label">发起人</span>
            <span class="req-grid-value">{{ currentReq.initiator }}</span>
          </div>
          <div class="req-grid-item">
            <span class="req-grid-label">创建时间</span>
            <span class="req-grid-value">{{ formatReqTime(currentReq.createdAt) }}</span>
          </div>
        </div>
        <el-divider />
        <div class="detail-section">
          <h4><el-icon><Document /></el-icon> 职位描述</h4>
          <p>{{ currentReq.jobDescription }}</p>
        </div>
        <div class="detail-section">
          <h4><el-icon><StarFilled /></el-icon> 任职要求</h4>
          <p>{{ currentReq.requirements }}</p>
        </div>
        <div v-if="currentReq.reason" class="detail-section">
          <h4><el-icon><QuestionFilled /></el-icon> 招聘原因</h4>
          <p>{{ currentReq.reason }}</p>
        </div>
        <div v-if="currentReq.latestOpinion" class="detail-section">
          <h4><el-icon><ChatDotRound /></el-icon> 最新审批意见</h4>
          <el-alert :title="currentReq.latestOpinion" type="info" :closable="false" show-icon />
        </div>
      </div>
      <template #footer>
        <el-button @click="reqDetailVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed, watch } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  Search, Plus, OfficeBuilding, LocationFilled, Money, UserFilled, Calendar, Document, StarFilled,
  Link, Loading, View, QuestionFilled, ChatDotRound
} from '@element-plus/icons-vue'
import { listJobs, createJob, updateJob, getRequirement } from '@/api'
import dayjs from 'dayjs'

const router = useRouter()
const loading = ref(false)
const jobList = ref([])
const locations = ['北京', '上海', '深圳', '杭州', '广州', '成都']

const filters = reactive({
  keyword: '',
  location: '',
  status: ''
})

const jobDialogVisible = ref(false)
const editingJob = ref(null)
const jobFormRef = ref(null)
const jobForm = reactive({
  title: '',
  department: '',
  location: '',
  salaryMin: 15000,
  salaryMax: 25000,
  status: 'open',
  headcount: 1,
  description: '',
  requirements: ''
})

const jobFormRules = {
  title: [{ required: true, message: '请输入岗位名称', trigger: 'blur' }],
  department: [{ required: true, message: '请输入部门', trigger: 'blur' }],
  location: [{ required: true, message: '请输入地点', trigger: 'blur' }],
  salaryMin: [{ required: true, message: '请输入薪资下限', trigger: 'blur' }],
  salaryMax: [{ required: true, message: '请输入薪资上限', trigger: 'blur' }],
  status: [{ required: true, message: '请选择状态', trigger: 'change' }],
  headcount: [{ required: true, message: '请输入需求人数', trigger: 'blur' }]
}

const jobDetailVisible = ref(false)
const currentJob = ref(null)
const reqCache = reactive({})
const reqLoadingIds = reactive({})
const reqDetailVisible = ref(false)
const currentReq = ref(null)

const currentJobReq = computed(() => {
  if (!currentJob.value || !currentJob.value.requirementId) return null
  return reqCache[currentJob.value.requirementId] || null
})

const reqStatusMap = {
  draft: { label: '草稿', type: 'info' },
  pending: { label: '待审批', type: 'warning' },
  approved: { label: '已通过', type: 'success' },
  returned: { label: '已退回', type: 'info' },
  rejected: { label: '已驳回', type: 'danger' },
  converted: { label: '已转职位', type: 'primary' }
}

const headcountTypeMap = {
  regular: '正式',
  contract: '外包',
  intern: '实习',
  replacement: '替补'
}

const reqStatusText = (s) => (reqStatusMap[s] || {}).label || s
const reqStatusTagType = (s) => (reqStatusMap[s] || {}).type || 'info'
const headcountTypeText = (t) => headcountTypeMap[t] || t
const formatReqTime = (t) => t ? dayjs(t).format('YYYY-MM-DD HH:mm') : '-'

const loadJobs = async () => {
  loading.value = true
  try {
    const res = await listJobs(filters)
    jobList.value = res.data || []
    jobList.value.forEach(job => {
      if (job.requirementId && !reqCache[job.requirementId] && !reqLoadingIds[job.requirementId]) {
        loadRequirement(job.requirementId)
      }
    })
  } finally {
    loading.value = false
  }
}

const loadRequirement = async (id) => {
  if (!id) return
  reqLoadingIds[id] = true
  try {
    const res = await getRequirement(id)
    if (res.data) {
      reqCache[id] = res.data
    }
  } catch (e) {
    console.error('加载需求信息失败:', e)
  } finally {
    reqLoadingIds[id] = false
  }
}

const showRequirementDetail = (id) => {
  openReqDetail(id)
}

const openReqDetail = async (id) => {
  if (!reqCache[id]) {
    await loadRequirement(id)
  }
  if (reqCache[id]) {
    currentReq.value = reqCache[id]
    reqDetailVisible.value = true
  } else {
    ElMessage.warning('需求信息加载失败，请稍后重试')
  }
}

const resetFilters = () => {
  filters.keyword = ''
  filters.location = ''
  filters.status = ''
  loadJobs()
}

const jobStatusText = (status) => {
  const map = { open: '招聘中', paused: '已暂停', closed: '已关闭' }
  return map[status] || status
}

const jobStatusTagType = (status) => {
  const map = { open: 'success', paused: 'warning', closed: 'info' }
  return map[status] || 'info'
}

const openJobDialog = (row) => {
  editingJob.value = row || null
  if (row) {
    Object.assign(jobForm, row)
  } else {
    Object.assign(jobForm, {
      title: '',
      department: '',
      location: '',
      salaryMin: 15000,
      salaryMax: 25000,
      status: 'open',
      headcount: 1,
      description: '',
      requirements: ''
    })
  }
  jobDialogVisible.value = true
}

const submitJobForm = async () => {
  if (!jobFormRef.value) return
  try {
    await jobFormRef.value.validate()
  } catch (e) {
    return
  }
  try {
    if (editingJob.value) {
      await updateJob(editingJob.value.id, jobForm)
      ElMessage.success('职位更新成功')
    } else {
      await createJob(jobForm)
      ElMessage.success('职位创建成功')
    }
    jobDialogVisible.value = false
    loadJobs()
  } catch (e) {
    // error handled in interceptor
  }
}

const showJobDetail = (row) => {
  currentJob.value = row
  jobDetailVisible.value = true
  if (row.requirementId && !reqCache[row.requirementId] && !reqLoadingIds[row.requirementId]) {
    loadRequirement(row.requirementId)
  }
}

const editCurrentJob = () => {
  jobDetailVisible.value = false
  openJobDialog(currentJob.value)
}

const goApply = (row) => {
  jobDetailVisible.value = false
  router.push({ path: '/applications', query: { applyJobId: row.id } })
}

onMounted(loadJobs)
</script>

<style scoped>
.jobs-page {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.filter-card {
  border-radius: 8px;
}

.filter-form {
  margin: 0;
}

.list-card {
  border-radius: 8px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.card-title {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
  margin: 0;
}

.count-tip {
  font-size: 13px;
  font-weight: normal;
  color: #909399;
  margin-left: 8px;
}

.job-title {
  color: #409eff;
  cursor: pointer;
  font-weight: 500;
}

.job-title:hover {
  text-decoration: underline;
}

.salary-text {
  color: #f56c6c;
  font-weight: 500;
}

.action-column {
  display: flex;
  flex-direction: column;
  gap: 6px;
  align-items: center;
}

.action-row-1 {
  display: flex;
  gap: 4px;
  justify-content: center;
}

.action-row-2 {
  display: flex;
  justify-content: center;
}

.req-tag {
  cursor: pointer;
  transition: all 0.2s;
}

.req-tag:hover {
  opacity: 0.8;
}

.loading-tag {
  cursor: default;
}

.job-detail .detail-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.detail-title {
  font-size: 22px;
  font-weight: 600;
  color: #1f2937;
  margin: 0;
}

.detail-meta {
  margin-top: 16px;
  display: flex;
  flex-wrap: wrap;
  gap: 20px;
  color: #606266;
  font-size: 14px;
}

.detail-meta span {
  display: flex;
  align-items: center;
  gap: 4px;
}

.salary-highlight {
  color: #f56c6c;
  font-weight: 600;
}

.divider-title {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 14px;
  font-weight: 600;
  color: #303133;
}

.detail-requirement {
  margin-top: 8px;
}

.req-loading {
  padding: 16px 0;
  text-align: center;
  color: #909399;
  font-size: 14px;
}

.req-info-card {
  background: #f5f7fa;
  border-radius: 8px;
  padding: 16px;
}

.req-info-row {
  display: flex;
  gap: 24px;
  margin-bottom: 12px;
}

.req-info-row:last-of-type {
  margin-bottom: 0;
}

.req-info-item {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 8px;
}

.req-info-label {
  color: #909399;
  font-size: 13px;
  flex-shrink: 0;
}

.req-info-value {
  color: #303133;
  font-size: 14px;
  font-weight: 500;
}

.req-action-row {
  margin-top: 12px;
  padding-top: 12px;
  border-top: 1px dashed #e4e7ed;
}

.detail-section {
  margin-bottom: 16px;
}

.detail-section h4 {
  font-size: 15px;
  color: #303133;
  margin: 0 0 10px;
  display: flex;
  align-items: center;
  gap: 6px;
}

.detail-section p {
  color: #606266;
  line-height: 1.7;
  margin: 0;
  white-space: pre-wrap;
}

.req-detail-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.req-detail-title {
  font-size: 20px;
  font-weight: 600;
  color: #1f2937;
  margin: 0 0 16px;
}

.req-detail-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
  padding: 16px;
  background: #f5f7fa;
  border-radius: 8px;
}

.req-grid-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.req-grid-label {
  font-size: 12px;
  color: #909399;
}

.req-grid-value {
  font-size: 14px;
  color: #303133;
  font-weight: 500;
}
</style>
