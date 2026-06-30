<template>
  <div class="candidates-page">
    <el-card class="filter-card" shadow="never">
      <el-form :inline="true" :model="filters" class="filter-form">
        <el-form-item label="关键词">
          <el-input
            v-model="filters.keyword"
            placeholder="姓名/手机号/期望岗位/技能"
            clearable
            style="width: 260px"
            @keyup.enter="loadCandidates"
          />
        </el-form-item>
        <el-form-item label="投递职位">
          <el-select v-model="filters.jobId" placeholder="全部职位" clearable filterable style="width: 200px" @change="loadCandidates">
            <el-option
              v-for="job in allJobs"
              :key="job.id"
              :label="job.title"
              :value="job.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="当前状态">
          <el-select v-model="filters.status" placeholder="全部状态" clearable style="width: 140px" @change="loadCandidates">
            <el-option
              v-for="item in statusOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadCandidates" :icon="Search">搜索</el-button>
          <el-button @click="resetFilters" :icon="Refresh">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card class="list-card" shadow="never">
      <div class="card-header">
        <h3 class="card-title">候选人列表 <span class="count-tip">共 {{ candidateList.length }} 位</span></h3>
      </div>
      <el-table :data="candidateList" v-loading="loading" stripe style="width: 100%">
        <el-table-column label="候选人" min-width="160">
          <template #default="{ row }">
            <div class="candidate-info">
              <el-avatar :size="40" class="avatar">{{ row.candidateName.charAt(0) }}</el-avatar>
              <div class="candidate-detail">
                <div class="candidate-name" @click="goToDetail(row)">
                  {{ row.candidateName }}
                  <el-tag v-if="row.noteCount > 0" type="warning" size="small" class="note-tag">
                    {{ row.noteCount }}条备注
                  </el-tag>
                </div>
                <div class="candidate-contact">{{ row.contact }}</div>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="期望岗位" min-width="150">
          <template #default="{ row }">{{ row.targetPosition }}</template>
        </el-table-column>
        <el-table-column label="投递职位" min-width="170">
          <template #default="{ row }">
            <div class="job-tags">
              <el-tag
                v-for="app in row.applications.slice(0, 2)"
                :key="app.id"
                type="info"
                size="small"
                effect="plain"
                style="margin: 2px 4px 2px 0"
              >
                {{ app.jobTitle }}
              </el-tag>
              <el-tag v-if="row.applications.length > 2" size="small" effect="plain">
                +{{ row.applications.length - 2 }}
              </el-tag>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="工作年限" width="100" align="center">
          <template #default="{ row }">{{ row.yearsOfExp }}年</template>
        </el-table-column>
        <el-table-column label="技能标签" min-width="220">
          <template #default="{ row }">
            <el-tag
              v-for="skill in (row.skills || []).slice(0, 4)"
              :key="skill"
              size="small"
              type="primary"
              effect="plain"
              style="margin: 2px 4px 2px 0"
            >
              {{ skill }}
            </el-tag>
            <el-tag v-if="row.skills && row.skills.length > 4" size="small" effect="plain">
              +{{ row.skills.length - 4 }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="当前状态" width="120" align="center">
          <template #default="{ row }">
            <el-tag :type="getStatusTagType(row.latestStatus)" size="small" effect="light">
              {{ getStatusText(row.latestStatus) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="最近沟通" width="160">
          <template #default="{ row }">
            {{ row.lastMessageTime ? formatTime(row.lastMessageTime) : '-' }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180" fixed="right" align="center">
          <template #default="{ row }">
            <el-button link type="primary" size="small" @click="goToDetail(row)">查看详情</el-button>
            <el-dropdown trigger="click" @command="(val) => changeStatus(row, val)">
              <el-button link type="warning" size="small">改状态<el-icon class="el-icon--right"><ArrowDown /></el-icon></el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item
                    v-for="item in changeableStatusOptions"
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
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Search, Refresh, ArrowDown } from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import { listCandidates, updateCandidateStatus, listJobs } from '@/api'

const router = useRouter()

const loading = ref(false)
const candidateList = ref([])
const allJobs = ref([])

const filters = reactive({
  keyword: '',
  jobId: '',
  status: ''
})

const statusOptions = [
  { value: 'submitted', label: '已投递', color: '#909399' },
  { value: 'pending', label: '待沟通', color: '#e6a23c' },
  { value: 'interview', label: '面试中', color: '#409eff' },
  { value: 'rejected', label: '已拒绝', color: '#f56c6c' },
  { value: 'hired', label: '已录用', color: '#67c23a' }
]

const changeableStatusOptions = statusOptions.filter(s => s.value !== 'submitted')

const getStatusText = (status) => {
  const opt = statusOptions.find(s => s.value === status)
  return opt ? opt.label : status
}

const getStatusTagType = (status) => {
  const map = { submitted: 'info', pending: 'warning', interview: 'primary', rejected: 'danger', hired: 'success' }
  return map[status] || 'info'
}

const formatTime = (t) => dayjs(t).format('YYYY-MM-DD HH:mm')

const loadCandidates = async () => {
  loading.value = true
  try {
    const res = await listCandidates(filters)
    candidateList.value = res.data || []
  } finally {
    loading.value = false
  }
}

const loadJobs = async () => {
  const res = await listJobs({})
  allJobs.value = res.data || []
}

const resetFilters = () => {
  filters.keyword = ''
  filters.jobId = ''
  filters.status = ''
  loadCandidates()
}

const goToDetail = (row) => {
  router.push({
    path: '/candidates/detail',
    query: { contact: row.contact, name: row.candidateName }
  })
}

const changeStatus = async (row, status) => {
  try {
    await updateCandidateStatus({
      contact: row.contact,
      name: row.candidateName,
      status
    })
    ElMessage.success(`状态已变更为「${getStatusText(status)}」`)
    loadCandidates()
  } catch (e) {}
}

onMounted(() => {
  loadCandidates()
  loadJobs()
})
</script>

<style scoped>
.candidates-page {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.filter-card,
.list-card {
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
  display: flex;
  align-items: center;
  gap: 6px;
}

.candidate-name:hover {
  color: #409eff;
}

.note-tag {
  margin-left: 4px;
}

.candidate-contact {
  font-size: 12px;
  color: #909399;
  margin-top: 2px;
}

.job-tags {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
}
</style>
