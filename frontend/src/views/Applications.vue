<template>
  <div class="applications-page">
    <el-card class="filter-card" shadow="never">
      <el-form :inline="true" :model="filters" class="filter-form">
        <el-form-item label="申请状态">
          <el-select v-model="filters.status" placeholder="全部" clearable style="width: 160px" @change="loadApplications">
            <el-option
              v-for="item in statusOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="关联职位">
          <el-select v-model="filters.jobId" placeholder="全部职位" clearable filterable style="width: 200px" @change="loadApplications">
            <el-option
              v-for="job in allJobs"
              :key="job.id"
              :label="job.title"
              :value="job.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadApplications" :icon="Refresh">刷新</el-button>
          <el-button type="success" @click="openApplyDialog()" :icon="Plus">模拟投递简历</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card class="stats-card" shadow="never">
      <div class="status-bar">
        <div
          v-for="item in statusOptions"
          :key="item.value"
          class="status-item"
          :class="{ active: filters.status === item.value }"
          @click="toggleStatusFilter(item.value)"
        >
          <div class="status-count" :style="{ color: item.color }">
            {{ getStatusCount(item.value) }}
          </div>
          <div class="status-label">{{ item.label }}</div>
        </div>
      </div>
    </el-card>

    <el-card class="list-card" shadow="never">
      <div class="card-header">
        <h3 class="card-title">申请列表 <span class="count-tip">共 {{ applicationList.length }} 条</span></h3>
      </div>
      <el-table :data="applicationList" v-loading="loading" stripe style="width: 100%">
        <el-table-column label="候选人" min-width="140">
          <template #default="{ row }">
            <div class="candidate-info">
              <el-avatar :size="38" class="avatar">{{ row.resume.candidateName.charAt(0) }}</el-avatar>
              <div class="candidate-detail">
                <div class="candidate-name" @click="showApplicationDetail(row)">{{ row.resume.candidateName }}</div>
                <div class="candidate-contact">{{ row.resume.Contact || row.resume.contact }}</div>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="jobTitle" label="应聘职位" min-width="170">
          <template #default="{ row }">
            <el-tag type="info" size="small" effect="plain">{{ row.jobTitle }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="期望岗位" min-width="140">
          <template #default="{ row }">{{ row.resume.targetPosition }}</template>
        </el-table-column>
        <el-table-column label="工作年限" width="100" align="center">
          <template #default="{ row }">{{ row.resume.yearsOfExp }}年</template>
        </el-table-column>
        <el-table-column label="技能标签" min-width="220">
          <template #default="{ row }">
            <el-tag
              v-for="skill in (row.resume.skills || []).slice(0, 4)"
              :key="skill"
              size="small"
              type="primary"
              effect="plain"
              style="margin: 2px 4px 2px 0"
            >
              {{ skill }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="120" align="center">
          <template #default="{ row }">
            <el-tag :type="getStatusTagType(row.status)" size="small" effect="light">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="投递时间" width="180">
          <template #default="{ row }">{{ formatTime(row.submittedAt) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="280" fixed="right" align="center">
          <template #default="{ row }">
            <el-button link type="primary" size="small" @click="showApplicationDetail(row)">查看</el-button>
            <el-button link type="success" size="small" @click="goMessage(row)">沟通</el-button>
            <el-dropdown trigger="click" @command="(val) => changeStatus(row, val)">
              <el-button link type="warning" size="small">改状态<el-icon class="el-icon--right"><ArrowDown /></el-icon></el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item
                    v-for="item in statusOptions"
                    :key="item.value"
                    :command="item.value"
                  >
                    {{ item.label }}
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog
      v-model="applyDialogVisible"
      :title="applyJob ? `投递简历 - ${applyJob.title}` : '投递简历'"
      width="600px"
      destroy-on-close
    >
      <el-form ref="applyFormRef" :model="applyForm" :rules="applyFormRules" label-width="100px">
        <el-form-item label="应聘职位" prop="jobId">
          <el-select v-model="applyForm.jobId" placeholder="请选择职位" filterable style="width: 100%">
            <el-option
              v-for="job in allJobs.filter(j => j.status === 'open')"
              :key="job.id"
              :label="`${job.title} - ${job.location}`"
              :value="job.id"
            />
          </el-select>
        </el-form-item>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="候选人姓名" prop="candidateName">
              <el-input v-model="applyForm.resume.candidateName" placeholder="请输入姓名" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="联系方式" prop="contact">
              <el-input v-model="applyForm.resume.contact" placeholder="请输入手机号" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="期望岗位" prop="targetPosition">
              <el-input v-model="applyForm.resume.targetPosition" placeholder="请输入期望岗位" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="工作年限" prop="yearsOfExp">
              <el-input-number v-model="applyForm.resume.yearsOfExp" :min="0" :max="40" style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="技能标签" prop="skills">
          <el-select
            v-model="applyForm.resume.skills"
            multiple
            filterable
            allow-create
            default-first-option
            placeholder="输入后回车添加技能标签"
            style="width: 100%"
          >
            <el-option v-for="s in commonSkills" :key="s" :label="s" :value="s" />
          </el-select>
        </el-form-item>
        <el-form-item label="简历摘要" prop="summary">
          <el-input v-model="applyForm.resume.summary" type="textarea" :rows="4" placeholder="请输入简历摘要，包括主要工作经历和项目经验" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="applyDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitApply">提交申请</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="detailVisible" title="申请详情" width="700px" destroy-on-close>
      <div v-if="currentApplication" class="application-detail">
        <div class="detail-status-bar">
          <div>
            <span class="detail-job-title">{{ currentApplication.jobTitle }}</span>
            <span class="detail-apply-time">投递于 {{ formatTime(currentApplication.submittedAt) }}</span>
          </div>
          <el-tag :type="getStatusTagType(currentApplication.status)" size="large" effect="light">
            {{ getStatusText(currentApplication.status) }}
          </el-tag>
        </div>
        <el-descriptions :column="2" border size="default" class="detail-desc">
          <el-descriptions-item label="候选人">
            <strong>{{ currentApplication.resume.candidateName }}</strong>
          </el-descriptions-item>
          <el-descriptions-item label="联系方式">{{ currentApplication.resume.contact }}</el-descriptions-item>
          <el-descriptions-item label="期望岗位">{{ currentApplication.resume.targetPosition }}</el-descriptions-item>
          <el-descriptions-item label="工作年限">{{ currentApplication.resume.yearsOfExp }}年</el-descriptions-item>
          <el-descriptions-item label="技能标签" :span="2">
            <el-tag
              v-for="skill in (currentApplication.resume.skills || [])"
              :key="skill"
              size="small"
              type="primary"
              effect="plain"
              style="margin: 2px 4px 2px 0"
            >
              {{ skill }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="简历摘要" :span="2">
            <div style="line-height: 1.7; white-space: pre-wrap;">{{ currentApplication.resume.summary }}</div>
          </el-descriptions-item>
        </el-descriptions>

        <el-divider />

        <div class="status-change-section">
          <h4 class="section-title">变更申请状态</h4>
          <div class="status-btns">
            <el-radio-group v-model="tempStatus" size="default">
              <el-radio-button
                v-for="item in statusOptions"
                :key="item.value"
                :value="item.value"
              >
                {{ item.label }}
              </el-radio-button>
            </el-radio-group>
            <el-button type="primary" @click="confirmChangeStatus">确认变更</el-button>
          </div>
        </div>
      </div>
      <template #footer>
        <el-button @click="detailVisible = false">关闭</el-button>
        <el-button type="success" @click="goMessage(currentApplication)">去沟通</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Refresh, Plus, ArrowDown } from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import {
  listApplications, createApplication, updateApplicationStatus, listJobs
} from '@/api'

const route = useRoute()
const router = useRouter()

const loading = ref(false)
const applicationList = ref([])
const allJobs = ref([])

const filters = reactive({
  status: '',
  jobId: ''
})

const statusOptions = [
  { value: 'submitted', label: '已投递', color: '#909399' },
  { value: 'pending', label: '待沟通', color: '#e6a23c' },
  { value: 'interview', label: '面试中', color: '#409eff' },
  { value: 'rejected', label: '已拒绝', color: '#f56c6c' },
  { value: 'hired', label: '已录用', color: '#67c23a' }
]

const commonSkills = ['Go', 'Java', 'Python', 'JavaScript', 'TypeScript', 'Vue', 'React', 'Node.js', 'MySQL', 'Redis', 'Kubernetes', '微服务', '需求分析', '产品设计', '数据分析']

const applyDialogVisible = ref(false)
const applyFormRef = ref(null)
const applyJob = ref(null)
const applyForm = reactive({
  jobId: '',
  resume: {
    candidateName: '',
    contact: '',
    targetPosition: '',
    yearsOfExp: 3,
    skills: [],
    summary: ''
  }
})

const applyFormRules = {
  jobId: [{ required: true, message: '请选择应聘职位', trigger: 'change' }],
  candidateName: [{ required: true, message: '请输入姓名', trigger: 'blur' }],
  contact: [{ required: true, message: '请输入联系方式', trigger: 'blur' }],
  targetPosition: [{ required: true, message: '请输入期望岗位', trigger: 'blur' }],
  yearsOfExp: [{ required: true, message: '请选择工作年限', trigger: 'blur' }],
  summary: [{ required: true, message: '请输入简历摘要', trigger: 'blur' }]
}

const detailVisible = ref(false)
const currentApplication = ref(null)
const tempStatus = ref('')

const getStatusText = (status) => {
  const opt = statusOptions.find(s => s.value === status)
  return opt ? opt.label : status
}

const getStatusTagType = (status) => {
  const map = { submitted: 'info', pending: 'warning', interview: 'primary', rejected: 'danger', hired: 'success' }
  return map[status] || 'info'
}

const formatTime = (t) => dayjs(t).format('YYYY-MM-DD HH:mm')

const loadApplications = async () => {
  loading.value = true
  try {
    const res = await listApplications(filters)
    applicationList.value = res.data || []
  } finally {
    loading.value = false
  }
}

const loadJobs = async () => {
  const res = await listJobs({})
  allJobs.value = res.data || []
}

const getStatusCount = (status) => {
  return applicationList.value.filter(a => a.status === status).length
}

const toggleStatusFilter = (status) => {
  filters.status = filters.status === status ? '' : status
  loadApplications()
}

const openApplyDialog = (job) => {
  applyJob.value = job || null
  applyForm.jobId = job ? job.id : ''
  applyForm.resume = { candidateName: '', contact: '', targetPosition: job ? job.title : '', yearsOfExp: 3, skills: [], summary: '' }
  applyDialogVisible.value = true
}

const submitApply = async () => {
  if (!applyFormRef.value) return
  try {
    await applyFormRef.value.validate()
  } catch (e) { return }

  const payload = {
    jobId: applyForm.jobId,
    resume: applyForm.resume
  }
  try {
    await createApplication(payload)
    ElMessage.success('简历投递成功')
    applyDialogVisible.value = false
    loadApplications()
  } catch (e) {}
}

const showApplicationDetail = (row) => {
  currentApplication.value = row
  tempStatus.value = row.status
  detailVisible.value = true
}

const goMessage = (row) => {
  detailVisible.value = false
  router.push({ path: '/messages', query: { applicationId: row.id } })
}

const changeStatus = async (row, status) => {
  try {
    await updateApplicationStatus(row.id, status)
    ElMessage.success(`状态已变更为「${getStatusText(status)}」`)
    loadApplications()
  } catch (e) {}
}

const confirmChangeStatus = async () => {
  if (!currentApplication.value) return
  if (tempStatus.value === currentApplication.value.status) {
    ElMessage.info('状态未变更')
    return
  }
  try {
    await updateApplicationStatus(currentApplication.value.id, tempStatus.value)
    ElMessage.success(`状态已变更为「${getStatusText(tempStatus.value)}」`)
    currentApplication.value.status = tempStatus.value
    loadApplications()
  } catch (e) {}
}

watch(
  () => route.query.applyJobId,
  (val) => {
    if (val) {
      const job = allJobs.value.find(j => j.id === val)
      if (job) {
        openApplyDialog(job)
      } else {
        loadJobs().then(() => {
          const j = allJobs.value.find(x => x.id === val)
          if (j) openApplyDialog(j)
        })
      }
    }
  },
  { immediate: true }
)

onMounted(() => {
  loadApplications()
  loadJobs()
})
</script>

<style scoped>
.applications-page {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.filter-card,
.list-card,
.stats-card {
  border-radius: 8px;
}

.card-header {
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

.status-bar {
  display: flex;
  gap: 16px;
  flex-wrap: wrap;
}

.status-item {
  flex: 1;
  min-width: 110px;
  padding: 12px 16px;
  border-radius: 8px;
  cursor: pointer;
  border: 1px solid #ebeef5;
  transition: all 0.2s;
  text-align: center;
}

.status-item:hover {
  border-color: #409eff;
  transform: translateY(-2px);
}

.status-item.active {
  background: #ecf5ff;
  border-color: #409eff;
}

.status-count {
  font-size: 28px;
  font-weight: 700;
  line-height: 1.2;
}

.status-label {
  font-size: 13px;
  color: #606266;
  margin-top: 4px;
}

.candidate-info {
  display: flex;
  align-items: center;
  gap: 10px;
}

.avatar {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: #fff;
  font-weight: 600;
}

.candidate-detail {
  display: flex;
  flex-direction: column;
}

.candidate-name {
  font-weight: 500;
  color: #303133;
  cursor: pointer;
}

.candidate-name:hover {
  color: #409eff;
}

.candidate-contact {
  font-size: 12px;
  color: #909399;
  margin-top: 2px;
}

.application-detail .detail-status-bar {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 20px;
}

.detail-job-title {
  font-size: 20px;
  font-weight: 600;
  color: #1f2937;
  margin-right: 12px;
}

.detail-apply-time {
  font-size: 13px;
  color: #909399;
}

.detail-desc {
  margin-top: 10px;
}

.status-change-section .section-title {
  font-size: 14px;
  font-weight: 600;
  color: #303133;
  margin: 0 0 12px;
}

.status-btns {
  display: flex;
  align-items: center;
  gap: 16px;
  flex-wrap: wrap;
}
</style>
