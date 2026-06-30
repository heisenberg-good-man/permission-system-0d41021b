<template>
  <div class="requirement-detail-page">
    <div class="page-header">
      <el-button @click="goBack" :icon="ArrowLeft">返回列表</el-button>
      <h2 class="page-title">用人需求详情</h2>
      <div class="header-actions">
        <el-tag :type="getStatusTagType(requirement?.status)" size="large" effect="light">
          {{ getStatusText(requirement?.status) }}
        </el-tag>
        <el-button
          v-if="canEdit"
          type="primary"
          plain
          @click="openEditDialog"
        >编辑需求</el-button>
        <el-button
          v-if="canSubmit"
          type="success"
          plain
          @click="handleSubmitApproval"
        >提交审批</el-button>
        <el-button
          v-if="canApprove"
          type="warning"
          @click="openApproveDialog"
        >处理审批</el-button>
        <el-button
          v-if="canConvert"
          type="primary"
          @click="handleConvertToJob"
        >生成招聘职位</el-button>
      </div>
    </div>

    <el-row v-loading="loading" :gutter="16">
      <el-col :span="16">
        <el-card class="info-card" shadow="never">
          <div class="card-header">
            <h3 class="card-title">基本信息</h3>
          </div>
          <el-descriptions :column="3" border v-if="requirement">
            <el-descriptions-item label="需求编号">{{ requirement.reqNo }}</el-descriptions-item>
            <el-descriptions-item label="部门">{{ requirement.department }}</el-descriptions-item>
            <el-descriptions-item label="岗位名称">{{ requirement.jobTitle }}</el-descriptions-item>
            <el-descriptions-item label="编制类型">
              <el-tag :type="getHeadcountTagType(requirement.headcountType)" size="small" effect="plain">
                {{ getHeadcountText(requirement.headcountType) }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="需求人数">{{ requirement.headcount }}人</el-descriptions-item>
            <el-descriptions-item label="期望到岗">{{ requirement.expectedOnboard }}</el-descriptions-item>
            <el-descriptions-item label="预算薪资" :span="3">
              <span class="salary-text">¥{{ formatSalary(requirement.salaryMin) }} - ¥{{ formatSalary(requirement.salaryMax) }}</span>
            </el-descriptions-item>
            <el-descriptions-item label="工作地点">{{ requirement.workLocation }}</el-descriptions-item>
            <el-descriptions-item label="发起人">{{ requirement.initiator }}</el-descriptions-item>
            <el-descriptions-item label="联系方式">{{ requirement.initiatorContact }}</el-descriptions-item>
            <el-descriptions-item label="创建时间">{{ formatTime(requirement.createdAt) }}</el-descriptions-item>
            <el-descriptions-item label="最近更新">{{ formatTime(requirement.updatedAt) }}</el-descriptions-item>
            <el-descriptions-item label="最近处理意见" :span="3">
              {{ requirement.latestOpinion || '-' }}
            </el-descriptions-item>
          </el-descriptions>
        </el-card>

        <el-card class="section-card" shadow="never" style="margin-top: 16px">
          <div class="card-header">
            <h3 class="card-title">
              <el-icon><Document /></el-icon>
              招聘原因
            </h3>
          </div>
          <div class="section-content" v-if="requirement">{{ requirement.reason || '暂无' }}</div>
        </el-card>

        <el-card class="section-card" shadow="never" style="margin-top: 16px">
          <div class="card-header">
            <h3 class="card-title">
              <el-icon><Document /></el-icon>
              岗位说明
            </h3>
          </div>
          <div class="section-content" v-if="requirement">{{ requirement.jobDescription || '暂无岗位说明' }}</div>
        </el-card>

        <el-card class="section-card" shadow="never" style="margin-top: 16px">
          <div class="card-header">
            <h3 class="card-title">
              <el-icon><StarFilled /></el-icon>
              能力要求
            </h3>
          </div>
          <div class="section-content" v-if="requirement">{{ requirement.requirements || '暂无能力要求' }}</div>
        </el-card>

        <el-card class="section-card" shadow="never" style="margin-top: 16px" v-if="requirement?.relatedJobID">
          <div class="card-header">
            <h3 class="card-title">
              <el-icon><Briefcase /></el-icon>
              关联职位
            </h3>
          </div>
          <div class="related-job-content">
            <div class="job-info">
              <div class="job-title">{{ requirement.jobTitle }}</div>
              <div class="job-meta">
                <el-tag type="info" size="small">{{ requirement.department }}</el-tag>
                <el-tag type="info" size="small">{{ requirement.workLocation }}</el-tag>
                <span class="job-salary">¥{{ formatSalary(requirement.salaryMin) }} - ¥{{ formatSalary(requirement.salaryMax) }}</span>
                <span>招{{ requirement.headcount }}人</span>
                <el-tag type="success" size="small" effect="plain">已关联职位ID: {{ requirement.relatedJobID }}</el-tag>
              </div>
            </div>
            <el-button type="primary" @click="goToJob(requirement.relatedJobID)">
              <el-icon><View /></el-icon>
              跳转职位详情
            </el-button>
          </div>
        </el-card>

        <el-card class="section-card" shadow="never" style="margin-top: 16px">
          <div class="card-header">
            <h3 class="card-title">
              <el-icon><Operation /></el-icon>
              操作记录
            </h3>
          </div>
          <el-table v-if="operationLogs && operationLogs.length > 0" :data="operationLogs" stripe>
            <el-table-column prop="operator" label="操作人" width="150" />
            <el-table-column prop="action" label="操作类型" width="110">
              <template #default="{ row }">
                <el-tag size="small" :type="getActionTagType(row.action)">{{ row.action }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="detail" label="操作内容" />
            <el-table-column label="操作时间" width="170">
              <template #default="{ row }">
                {{ formatTime(row.createdAt) }}
              </template>
            </el-table-column>
          </el-table>
          <el-empty v-else description="暂无操作记录" :image-size="100" />
        </el-card>
      </el-col>

      <el-col :span="8">
        <el-card class="timeline-card" shadow="never">
          <div class="card-header">
            <h3 class="card-title">
              <el-icon><Promotion /></el-icon>
              审批节点
            </h3>
          </div>
          <el-steps v-if="approvalNodes && approvalNodes.length > 0" :active="currentStepIndex" finish-status="success" direction="vertical">
            <el-step
              v-for="(step, index) in approvalNodes"
              :key="index"
              :title="step.title"
              :description="step.description"
              :status="step.status"
            >
              <template #icon>
                <div class="step-icon-wrapper" :class="step.status">
                  <el-icon v-if="step.status === 'success'"><CircleCheck /></el-icon>
                  <el-icon v-else-if="step.status === 'error'"><CircleClose /></el-icon>
                  <el-icon v-else-if="step.status === 'process'"><Loading /></el-icon>
                  <el-icon v-else><Clock /></el-icon>
                </div>
              </template>
            </el-step>
          </el-steps>

          <el-timeline v-if="approvalTimeline && approvalTimeline.length > 0" class="approval-timeline">
            <el-timeline-item
              v-for="(item, index) in approvalTimeline"
              :key="index"
              :timestamp="item.time"
              :type="getTimelineType(item.status)"
              placement="top"
            >
              <el-card shadow="never" class="timeline-item-card">
                <div class="timeline-item-header">
                  <span class="approver-name">{{ item.approver }}</span>
                  <el-tag :type="getApprovalTagType(item.status)" size="small" effect="light">
                    {{ item.statusText }}
                  </el-tag>
                </div>
                <div v-if="item.opinion" class="timeline-opinion">
                  <span class="opinion-label">处理意见：</span>
                  <span class="opinion-content">{{ item.opinion }}</span>
                </div>
              </el-card>
            </el-timeline-item>
          </el-timeline>
          <el-empty v-if="(!approvalNodes || approvalNodes.length === 0) && (!approvalTimeline || approvalTimeline.length === 0)" description="暂无审批记录，请先提交审批" :image-size="100" />
        </el-card>
      </el-col>
    </el-row>

    <RequirementDialog
      v-model="editDialogVisible"
      :editing-row="requirement"
      :existing-requirements="[]"
      @success="loadDetail"
    />

    <el-dialog v-model="approveDialogVisible" title="审批处理" width="520px" destroy-on-close>
      <el-form :model="approveForm" label-width="100px" :rules="approveRules" ref="approveFormRef">
        <el-form-item label="处理动作" prop="action">
          <el-radio-group v-model="approveForm.action">
            <el-radio value="approve">通过</el-radio>
            <el-radio value="return">退回修改</el-radio>
            <el-radio value="reject">驳回</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="处理意见" prop="opinion">
          <el-input v-model="approveForm.opinion" type="textarea" :rows="4" placeholder="请填写处理意见" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="approveDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitApproval">提交</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  ArrowLeft,
  Document,
  StarFilled,
  Briefcase,
  Operation,
  Promotion,
  CircleCheck,
  CircleClose,
  Loading,
  Clock,
  View
} from '@element-plus/icons-vue'
import { getRequirement, submitApproval as apiSubmitApproval, approveRequirement, convertToJob } from '@/api'
import RequirementDialog from './RequirementDialog.vue'

const route = useRoute()
const router = useRouter()

const loading = ref(false)
const requirement = ref(null)
const operationLogs = ref([])
const approvalNodes = ref([])
const approvalTimeline = ref([])
const currentStepIndex = ref(0)
const editDialogVisible = ref(false)
const approveDialogVisible = ref(false)
const approveFormRef = ref(null)

const approveForm = reactive({
  action: 'approve',
  opinion: ''
})

const approveRules = {
  opinion: [{ required: true, message: '请填写处理意见', trigger: 'blur' }]
}

const statusOptions = [
  { value: 'draft', label: '草稿' },
  { value: 'pending', label: '待审批' },
  { value: 'approved', label: '已通过' },
  { value: 'returned', label: '已退回' },
  { value: 'rejected', label: '已驳回' },
  { value: 'converted', label: '已转职位' }
]

const headcountOptions = [
  { value: 'regular', label: '正式编制' },
  { value: 'contract', label: '合同工' },
  { value: 'intern', label: '实习生' },
  { value: 'replacement', label: '替补编制' }
]

const getStatusText = (status) => {
  const opt = statusOptions.find(s => s.value === status)
  return opt ? opt.label : status || '-'
}

const getStatusTagType = (status) => {
  const map = {
    draft: 'info',
    pending: 'warning',
    approved: 'success',
    returned: 'warning',
    rejected: 'danger',
    converted: 'success'
  }
  return map[status] || 'info'
}

const getHeadcountText = (t) => {
  const opt = headcountOptions.find(o => o.value === t)
  return opt ? opt.label : t || '-'
}

const getHeadcountTagType = (t) => {
  const map = {
    regular: 'success',
    contract: 'warning',
    intern: 'primary',
    replacement: 'info'
  }
  return map[t] || 'info'
}

const getTimelineType = (status) => {
  const map = {
    approved: 'success',
    rejected: 'danger',
    returned: 'warning',
    pending: 'primary',
    draft: 'info'
  }
  return map[status] || 'primary'
}

const getApprovalTagType = (status) => {
  const map = {
    approved: 'success',
    rejected: 'danger',
    returned: 'warning',
    pending: 'primary',
    draft: 'info'
  }
  return map[status] || 'info'
}

const getActionTagType = (action) => {
  const map = {
    '创建': 'info',
    '提交审批': 'warning',
    '通过': 'success',
    '退回': 'warning',
    '拒绝': 'danger',
    '驳回': 'danger',
    '转职位': 'success',
    '编辑': 'primary'
  }
  return map[action] || 'info'
}

const formatSalary = (s) => {
  if (!s) return '0'
  if (s >= 1000) return (s / 1000) + 'K'
  return s
}

const formatTime = (t) => {
  if (!t) return '-'
  return t.replace('T', ' ').slice(0, 16)
}

const canEdit = computed(() => requirement.value && (requirement.value.status === 'draft' || requirement.value.status === 'returned'))
const canSubmit = computed(() => requirement.value && (requirement.value.status === 'draft' || requirement.value.status === 'returned'))
const canApprove = computed(() => {
  if (!requirement.value || requirement.value.status !== 'pending') return false
  if (!requirement.value.approvalNodes || requirement.value.approvalNodes.length === 0) return false
  return requirement.value.approvalNodes.some(n => n.status === 'pending')
})
const canConvert = computed(() => requirement.value && requirement.value.status === 'approved' && !requirement.value.relatedJobID)

const pendingNodeIndex = computed(() => {
  if (!requirement.value?.approvalNodes) return 0
  const idx = requirement.value.approvalNodes.findIndex(n => n.status === 'pending')
  return idx >= 0 ? idx : 0
})

const buildApprovalData = () => {
  const nodes = requirement.value?.approvalNodes || []
  approvalNodes.value = nodes.map((n, idx) => {
    let status = 'wait'
    if (n.status === 'approved') status = 'success'
    else if (n.status === 'rejected' || n.status === 'returned') status = 'error'
    else if (n.status === 'pending') status = 'process'
    const desc = [
      n.approver ? `审批人: ${n.approver}` : '',
      n.handledAt ? `处理时间: ${formatTime(n.handledAt)}` : ''
    ].filter(Boolean).join('\n')
    return {
      title: n.nodeName || `节点${idx + 1}`,
      description: desc || '等待处理',
      status
    }
  })
  const idx = nodes.findIndex(n => n.status === 'pending')
  currentStepIndex.value = idx >= 0 ? idx : (nodes.every(n => n.status === 'approved') ? nodes.length : 0)
  approvalTimeline.value = nodes.filter(n => n.handledAt || n.status !== 'draft').map(n => ({
    approver: n.approver || n.nodeName,
    status: n.status,
    statusText: getStatusText(n.status) === n.status ? (
      n.status === 'pending' ? '处理中' : n.status === 'approved' ? '通过' : n.status === 'returned' ? '退回' : n.status === 'rejected' ? '驳回' : '待处理'
    ) : getStatusText(n.status),
    opinion: n.opinion,
    time: formatTime(n.handledAt || n.createdAt)
  }))
}

const loadDetail = async () => {
  const { id } = route.query
  if (!id) {
    ElMessage.error('参数错误：缺少需求ID')
    return
  }
  loading.value = true
  try {
    const res = await getRequirement(id)
    const data = res.data || {}
    requirement.value = data
    operationLogs.value = data.operationLogs || []
    buildApprovalData()
  } finally {
    loading.value = false
  }
}

const openEditDialog = () => {
  editDialogVisible.value = true
}

const handleSubmitApproval = async () => {
  try {
    await ElMessageBox.confirm(
      `确定要提交「${requirement.value.jobTitle}」的审批吗？`,
      '提交确认',
      { type: 'warning' }
    )
  } catch { return }
  loading.value = true
  try {
    await apiSubmitApproval(requirement.value.id)
    ElMessage.success('已提交审批')
    await loadDetail()
  } finally {
    loading.value = false
  }
}

const openApproveDialog = () => {
  approveForm.action = 'approve'
  approveForm.opinion = ''
  approveDialogVisible.value = true
}

const submitApproval = async () => {
  let valid = true
  try { await approveFormRef.value?.validate() } catch { valid = false }
  if (!valid) return
  loading.value = true
  try {
    await approveRequirement(requirement.value.id, {
      nodeIndex: pendingNodeIndex.value,
      action: approveForm.action,
      opinion: approveForm.opinion
    })
    ElMessage.success('审批处理完成')
    approveDialogVisible.value = false
    await loadDetail()
  } finally {
    loading.value = false
  }
}

const handleConvertToJob = async () => {
  try {
    await ElMessageBox.confirm(
      `确定要为「${requirement.value.jobTitle}」生成招聘职位吗？`,
      '生成确认',
      { type: 'info' }
    )
  } catch { return }
  loading.value = true
  try {
    const res = await convertToJob(requirement.value.id, {})
    ElMessage.success('已生成招聘职位')
    await loadDetail()
    if (res.data?.id) {
      setTimeout(() => goToJob(res.data.id), 500)
    }
  } finally {
    loading.value = false
  }
}

const goBack = () => {
  router.push({ path: '/requirements' })
}

const goToJob = (jobId) => {
  router.push({ path: '/jobs', query: { jobId } })
}

onMounted(loadDetail)
</script>

<style scoped>
.requirement-detail-page {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.page-header {
  display: flex;
  align-items: center;
  gap: 16px;
}

.page-title {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
  color: #1f2937;
  flex: 1;
}

.header-actions {
  display: flex;
  gap: 10px;
  align-items: center;
}

.info-card,
.section-card,
.timeline-card {
  border-radius: 8px;
}

.card-header {
  margin-bottom: 16px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-title {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
  margin: 0;
  display: flex;
  align-items: center;
  gap: 6px;
}

.salary-text {
  color: #f56c6c;
  font-weight: 600;
}

.section-content {
  font-size: 14px;
  line-height: 1.8;
  color: #606266;
  white-space: pre-wrap;
}

.related-job-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background: #f5f7fa;
  border-radius: 8px;
}

.job-info {
  flex: 1;
}

.job-title {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 8px;
}

.job-meta {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 13px;
  color: #606266;
  flex-wrap: wrap;
}

.job-salary {
  color: #f56c6c;
  font-weight: 500;
}

.approval-timeline {
  margin-top: 24px;
  padding-top: 24px;
  border-top: 1px solid #ebeef5;
}

.step-icon-wrapper {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f0f2f5;
  color: #909399;
  font-size: 16px;
}

.step-icon-wrapper.success {
  background: #f0f9eb;
  color: #67c23a;
}

.step-icon-wrapper.error {
  background: #fef0f0;
  color: #f56c6c;
}

.step-icon-wrapper.process {
  background: #ecf5ff;
  color: #409eff;
}

.timeline-item-card {
  background: #f8f9fa;
  border: none;
}

.timeline-item-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.approver-name {
  font-weight: 500;
  color: #303133;
  font-size: 14px;
}

.timeline-opinion {
  font-size: 13px;
  color: #606266;
  line-height: 1.6;
}

.opinion-label {
  color: #909399;
}

.opinion-content {
  white-space: pre-wrap;
}
</style>
