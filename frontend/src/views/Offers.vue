<template>
  <div class="offers-page">
    <el-card class="filter-card" shadow="never">
      <el-form :inline="true" :model="filters" class="filter-form">
        <el-form-item label="投递职位">
          <el-select v-model="filters.jobId" placeholder="全部职位" clearable filterable style="width: 200px" @change="loadOffers">
            <el-option
              v-for="job in allJobs"
              :key="job.id"
              :label="job.title"
              :value="job.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="Offer 状态">
          <el-select v-model="filters.status" placeholder="全部" clearable style="width: 140px" @change="loadOffers">
            <el-option
              v-for="item in statusOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadOffers" :icon="Refresh">刷新</el-button>
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
        <h3 class="card-title">Offer 列表 <span class="count-tip">共 {{ offerList.length }} 条</span></h3>
      </div>
      <el-empty v-if="!loading && offerList.length === 0" description="暂无 Offer 记录，面试通过的候选人可发起 Offer">
        <el-button type="primary" @click="loadOffers">重新加载</el-button>
      </el-empty>
      <el-table :data="offerList" v-loading="loading" stripe style="width: 100%">
        <el-table-column label="候选人" min-width="130">
          <template #default="{ row }">
            <div class="candidate-info">
              <el-avatar :size="38" class="avatar">{{ row.candidateName?.charAt(0) || '-' }}</el-avatar>
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
        <el-table-column label="薪资范围" min-width="140">
          <template #default="{ row }">
            <span class="salary-text">¥{{ formatSalary(row.salaryMin) }}K - ¥{{ formatSalary(row.salaryMax) }}K</span>
          </template>
        </el-table-column>
        <el-table-column prop="startDate" label="入职日期" width="120" />
        <el-table-column prop="owner" label="负责人" min-width="130" />
        <el-table-column label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="getStatusTagType(row.status)" size="small" effect="light">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="创建时间" width="180">
          <template #default="{ row }">{{ formatTime(row.createdAt) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="360" fixed="right" align="center">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="showOfferDetail(row)">详情</el-button>
            <el-button type="success" link size="small" @click="handleSend(row)" :disabled="row.status !== 'pending'">发送</el-button>
            <el-button type="success" link size="small" @click="handleReply(row, true)" :disabled="row.status !== 'sent'">接受</el-button>
            <el-button type="warning" link size="small" @click="handleReply(row, false)" :disabled="row.status !== 'sent'">拒绝</el-button>
            <el-popconfirm
              title="确认撤回该 Offer？"
              confirm-button-text="确认撤回"
              cancel-button-text="再想想"
              @confirm="handleWithdraw(row)"
            >
              <template #reference>
                <el-button type="danger" link size="small" :disabled="row.status === 'accepted' || row.status === 'rejected' || row.status === 'withdrawn'">撤回</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="detailVisible" title="Offer 详情" width="640px" destroy-on-close>
      <OfferDetailContent v-if="currentOffer" :offer="currentOffer" @updated="loadOffers" />
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, defineAsyncComponent } from 'vue'
import { ElMessage } from 'element-plus'
import { Refresh, RefreshLeft } from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import { listOffers, listJobs, sendOffer, replyOffer, withdrawOffer } from '@/api'

const OfferDetailContent = defineAsyncComponent(() => import('./OfferDetailContent.vue'))

const loading = ref(false)
const allJobs = ref([])
const offerList = ref([])

const filters = reactive({
  jobId: '',
  status: ''
})

const statusOptions = [
  { value: 'pending', label: '待发送', color: '#e6a23c' },
  { value: 'sent', label: '已发送', color: '#409eff' },
  { value: 'accepted', label: '已接受', color: '#67c23a' },
  { value: 'rejected', label: '已拒绝', color: '#f56c6c' },
  { value: 'withdrawn', label: '已撤回', color: '#909399' }
]

const getStatusCount = (status) => {
  if (!offerList.value) return 0
  return offerList.value.filter(o => o.status === status).length
}

const toggleStatusFilter = (status) => {
  filters.status = filters.status === status ? '' : status
  loadOffers()
}

const resetFilters = () => {
  filters.jobId = ''
  filters.status = ''
  loadOffers()
}

const formatTime = (t) => {
  if (!t) return '-'
  return dayjs(t).format('YYYY-MM-DD HH:mm')
}

const formatSalary = (s) => {
  if (!s) return 0
  return (s / 1000).toFixed(0)
}

const getStatusText = (s) => {
  const map = { pending: '待发送', sent: '已发送', accepted: '已接受', rejected: '已拒绝', withdrawn: '已撤回' }
  return map[s] || s
}

const getStatusTagType = (s) => {
  const map = { pending: 'warning', sent: 'primary', accepted: 'success', rejected: 'danger', withdrawn: 'info' }
  return map[s] || 'info'
}

const loadJobs = async () => {
  try {
    const res = await listJobs()
    allJobs.value = res?.data || []
  } catch (e) {
    console.error('加载职位列表失败:', e)
  }
}

const loadOffers = async () => {
  loading.value = true
  try {
    const params = {}
    if (filters.jobId) params.jobId = filters.jobId
    if (filters.status) params.status = filters.status
    const res = await listOffers(params)
    offerList.value = res?.data || []
  } catch (e) {
    console.error('加载 Offer 列表失败:', e)
    ElMessage.error('加载 Offer 列表失败')
  } finally {
    loading.value = false
  }
}

const currentOffer = ref(null)
const detailVisible = ref(false)
const showOfferDetail = (row) => {
  currentOffer.value = row
  detailVisible.value = true
}

const handleSend = async (row) => {
  try {
    await sendOffer(row.id)
    ElMessage.success('Offer 已发送')
    loadOffers()
  } catch (e) {
    ElMessage.error('发送失败：' + (e?.response?.data?.error || e.message))
  }
}

const handleReply = async (row, accepted) => {
  const action = accepted ? '接受' : '拒绝'
  try {
    await replyOffer(row.id, { accepted, feedback: `候选人已${action}Offer` })
    ElMessage.success(`候选人已${action}`)
    loadOffers()
  } catch (e) {
    ElMessage.error(`操作失败：` + (e?.response?.data?.error || e.message))
  }
}

const handleWithdraw = async (row) => {
  try {
    await withdrawOffer(row.id, { reason: '公司撤回' })
    ElMessage.success('Offer 已撤回')
    loadOffers()
  } catch (e) {
    ElMessage.error('撤回失败：' + (e?.response?.data?.error || e.message))
  }
}

onMounted(() => {
  loadJobs()
  loadOffers()
})
</script>

<style scoped>
.offers-page {
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
  background: linear-gradient(135deg, #67c23a, #85ce61);
}

.candidate-name {
  font-size: 14px;
  font-weight: 500;
  color: #303133;
}

.candidate-contact {
  font-size: 12px;
  color: #909399;
  margin-top: 2px;
}

.salary-text {
  font-weight: 600;
  color: #f56c6c;
}
</style>
