<template>
  <div class="interviews-page">
    <el-card class="filter-card" shadow="never">
      <el-form :inline="true" :model="filters" class="filter-form">
        <el-form-item label="投递职位">
          <el-select v-model="filters.jobId" placeholder="全部职位" clearable filterable style="width: 200px" @change="loadInterviews">
            <el-option
              v-for="job in allJobs"
              :key="job.id"
              :label="job.title"
              :value="job.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="面试状态">
          <el-select v-model="filters.status" placeholder="全部" clearable style="width: 140px" @change="loadInterviews">
            <el-option
              v-for="item in statusOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="面试方式">
          <el-select v-model="filters.method" placeholder="全部" clearable style="width: 140px" @change="loadInterviews">
            <el-option
              v-for="item in methodOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadInterviews" :icon="Refresh">刷新</el-button>
          <el-button @click="resetFilters" :icon="RefreshLeft">重置</el-button>
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
        <h3 class="card-title">面试安排 <span class="count-tip">共 {{ interviewList.length }} 条</span></h3>
      </div>
      <el-table :data="interviewList" v-loading="loading" stripe style="width: 100%">
        <el-table-column label="候选人" min-width="130">
          <template #default="{ row }">
            <div class="candidate-info">
              <el-avatar :size="38" class="avatar">{{ row.candidateName.charAt(0) }}</el-avatar>
              <div class="candidate-detail">
                <div class="candidate-name">{{ row.candidateName }}</div>
                <div class="candidate-contact">{{ row.contact }}</div>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="jobTitle" label="投递职位" min-width="160">
          <template #default="{ row }">
            <el-tag type="info" size="small" effect="plain">{{ row.jobTitle }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="面试轮次" width="100" align="center">
          <template #default="{ row }">第{{ row.round }}轮</template>
        </el-table-column>
        <el-table-column label="面试时间" width="180">
          <template #default="{ row }">{{ formatTime(row.interviewTime) }}</template>
        </el-table-column>
        <el-table-column label="面试方式" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="getMethodTagType(row.method)" size="small" effect="light">
              {{ getMethodText(row.method) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="interviewer" label="面试官" min-width="130" />
        <el-table-column label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="getStatusTagType(row.status)" size="small" effect="light">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="最近备注" min-width="180" show-overflow-tooltip>
          <template #default="{ row }">{{ row.latestNote || '-' }}</template>
        </el-table-column>
        <el-table-column label="操作" width="280" fixed="right" align="center">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="showInterviewDetail(row)">详情</el-button>
            <el-button type="success" link size="small" @click="openRescheduleDialog(row)" :disabled="row.status !== 'scheduled'">改期</el-button>
            <el-popconfirm
              title="确认取消该面试安排？"
              confirm-button-text="确认取消"
              cancel-button-text="再想想"
              @confirm="handleCancel(row)"
            >
              <template #reference>
                <el-button type="warning" link size="small" :disabled="row.status !== 'scheduled'">取消</el-button>
              </template>
            </el-popconfirm>
            <el-button type="info" link size="small" @click="openCompleteDialog(row)" :disabled="row.status !== 'scheduled'">完成</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="detailVisible" title="面试详情" width="720px" destroy-on-close>
      <InterviewDetailContent v-if="currentInterview" :interview="currentInterview" @updated="loadInterviews" />
    </el-dialog>

    <el-dialog v-model="rescheduleVisible" title="面试改期" width="480px" destroy-on-close>
      <el-form :model="rescheduleForm" label-width="100px" :rules="rescheduleRules" ref="rescheduleFormRef">
        <el-form-item label="当前时间">
          <span>{{ formatTime(currentInterview?.interviewTime) }}</span>
        </el-form-item>
        <el-form-item label="新面试时间" prop="interviewTime">
          <el-date-picker
            v-model="rescheduleForm.interviewTime"
            type="datetime"
            placeholder="请选择面试时间"
            style="width: 100%"
            format="YYYY-MM-DD HH:mm"
            value-format="YYYY-MM-DDTHH:mm:ssZ"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="rescheduleVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="confirmReschedule">确认改期</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="completeVisible" title="标记面试完成" width="520px" destroy-on-close>
      <el-form :model="completeForm" label-width="100px">
        <el-form-item label="面试反馈">
          <el-input
            v-model="completeForm.feedback"
            type="textarea"
            :rows="4"
            placeholder="请输入面试整体评价和反馈意见，为后续面试提供参考"
          />
        </el-form-item>
        <el-form-item label="备注">
          <el-input
            v-model="completeForm.note"
            type="textarea"
            :rows="2"
            placeholder="可填写后续安排或其他备注"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="completeVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="confirmComplete">确认完成</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Refresh, RefreshLeft } from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import { listInterviews, listJobs, rescheduleInterview, cancelInterview, completeInterview } from '@/api'
import InterviewDetailContent from './InterviewDetailContent.vue'

const router = useRouter()
const loading = ref(false)
const submitting = ref(false)
const allJobs = ref([])
const interviewList = ref([])

const filters = reactive({
  jobId: '',
  status: '',
  method: ''
})

const statusOptions = [
  { value: 'scheduled', label: '待面试', color: '#409eff' },
  { value: 'completed', label: '已完成', color: '#67c23a' },
  { value: 'cancelled', label: '已取消', color: '#909399' }
]

const methodOptions = [
  { value: 'onsite', label: '现场' },
  { value: 'online', label: '视频' },
  { value: 'phone', label: '电话' }
]

const getStatusCount = (status) => {
  if (!interviewList.value) return 0
  return interviewList.value.filter(i => i.status === status).length
}

const toggleStatusFilter = (status) => {
  filters.status = filters.status === status ? '' : status
  loadInterviews()
}

const resetFilters = () => {
  filters.jobId = ''
  filters.status = ''
  filters.method = ''
  loadInterviews()
}

const formatTime = (t) => {
  if (!t) return '-'
  return dayjs(t).format('YYYY-MM-DD HH:mm')
}

const getStatusText = (s) => {
  const map = { scheduled: '待面试', completed: '已完成', cancelled: '已取消' }
  return map[s] || s
}

const getStatusTagType = (s) => {
  const map = { scheduled: 'primary', completed: 'success', cancelled: 'info' }
  return map[s] || 'info'
}

const getMethodText = (m) => {
  const map = { onsite: '现场', online: '视频', phone: '电话' }
  return map[m] || m
}

const getMethodTagType = (m) => {
  const map = { onsite: 'warning', online: 'primary', phone: 'info' }
  return map[m] || 'info'
}

const loadJobs = async () => {
  try {
    const res = await listJobs()
    allJobs.value = res?.data || []
  } catch (e) {
    console.error('加载职位列表失败:', e)
  }
}

const loadInterviews = async () => {
  loading.value = true
  try {
    const params = {}
    if (filters.jobId) params.jobId = filters.jobId
    if (filters.status) params.status = filters.status
    if (filters.method) params.method = filters.method
    const res = await listInterviews(params)
    interviewList.value = res?.data || []
  } catch (e) {
    console.error('加载面试列表失败:', e)
    ElMessage.error('加载面试列表失败')
  } finally {
    loading.value = false
  }
}

const currentInterview = ref(null)
const detailVisible = ref(false)
const showInterviewDetail = (row) => {
  currentInterview.value = row
  detailVisible.value = true
}

const rescheduleVisible = ref(false)
const rescheduleForm = reactive({ interviewTime: '' })
const rescheduleFormRef = ref(null)
const rescheduleRules = {
  interviewTime: [{ required: true, message: '请选择新的面试时间', trigger: 'change' }]
}
const openRescheduleDialog = (row) => {
  currentInterview.value = row
  rescheduleForm.interviewTime = ''
  rescheduleVisible.value = true
}

const confirmReschedule = async () => {
  await rescheduleFormRef.value?.validate()
  if (!rescheduleForm.interviewTime) return
  submitting.value = true
  try {
    await rescheduleInterview(currentInterview.value.id, { interviewTime: rescheduleForm.interviewTime })
    ElMessage.success('面试改期成功')
    rescheduleVisible.value = false
    loadInterviews()
  } catch (e) {
    ElMessage.error('改期失败：' + (e?.response?.data?.error || e.message))
  } finally {
    submitting.value = false
  }
}

const handleCancel = async () => {
  try {
    await cancelInterview(currentInterview.value.id)
    ElMessage.success('已取消面试安排')
    loadInterviews()
  } catch (e) {
    ElMessage.error('取消失败：' + (e?.response?.data?.error || e.message))
  }
}

const completeVisible = ref(false)
const completeForm = reactive({ feedback: '', note: '' })
const openCompleteDialog = (row) => {
  currentInterview.value = row
  completeForm.feedback = ''
  completeForm.note = ''
  completeVisible.value = true
}

const confirmComplete = async () => {
  submitting.value = true
  try {
    await completeInterview(currentInterview.value.id, completeForm)
    ElMessage.success('面试已标记完成')
    completeVisible.value = false
    loadInterviews()
  } catch (e) {
    ElMessage.error('操作失败：' + (e?.response?.data?.error || e.message))
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  loadJobs()
  loadInterviews()
})
</script>

<style scoped>
.interviews-page {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.filter-card :deep(.el-form-item) {
  margin-bottom: 0;
}

.stats-card {
  padding: 12px 20px;
}

.status-bar {
  display: flex;
  gap: 4px;
}

.status-item {
  flex: 1;
  padding: 12px 16px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
  text-align: center;
}

.status-item:hover {
  background: #f5f7fa;
}

.status-item.active {
  background: #ecf5ff;
}

.status-count {
  font-size: 24px;
  font-weight: 600;
  margin-bottom: 4px;
}

.status-label {
  font-size: 13px;
  color: #666;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}

.card-title {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.count-tip {
  color: #909399;
  font-size: 13px;
  font-weight: 400;
  margin-left: 8px;
}

.candidate-info {
  display: flex;
  align-items: center;
  gap: 10px;
}

.avatar {
  background: linear-gradient(135deg, #409eff, #66b1ff);
}

.candidate-name {
  font-size: 14px;
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
</style>
