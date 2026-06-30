<template>
  <div class="requirements-page">
    <el-card class="filter-card" shadow="never">
      <el-form :inline="true" :model="filters" class="filter-form">
        <el-form-item label="部门">
          <el-select v-model="filters.department" placeholder="全部" clearable style="width: 160px">
            <el-option v-for="dept in departments" :key="dept" :label="dept" :value="dept" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="filters.status" placeholder="全部" clearable style="width: 140px">
            <el-option label="草稿" value="draft" />
            <el-option label="待审批" value="pending" />
            <el-option label="已通过" value="approved" />
            <el-option label="已退回" value="returned" />
            <el-option label="已驳回" value="rejected" />
            <el-option label="已转职位" value="converted" />
          </el-select>
        </el-form-item>
        <el-form-item label="关键词">
          <el-input
            v-model="filters.keyword"
            placeholder="需求编号/岗位名称/发起人"
            clearable
            style="width: 240px"
            :prefix-icon="Search"
            @keyup.enter="loadRequirements"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadRequirements" :icon="Search">搜索</el-button>
          <el-button @click="resetFilters">重置</el-button>
          <el-button type="success" @click="openDialog()" :icon="Plus">新建需求</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card class="list-card" shadow="never">
      <div class="card-header">
        <h3 class="card-title">用人需求列表 <span class="count-tip">共 {{ requirementList.length }} 条</span></h3>
      </div>
      <el-empty v-if="!loading && requirementList.length === 0" description="暂无需求数据">
        <el-button type="primary" @click="loadRequirements">重新加载</el-button>
      </el-empty>
      <el-table v-else :data="requirementList" v-loading="loading" stripe style="width: 100%">
        <el-table-column prop="reqNo" label="需求编号" width="160" />
        <el-table-column prop="department" label="部门" width="100" />
        <el-table-column prop="jobTitle" label="岗位名称" min-width="160">
          <template #default="{ row }">
            <span class="position-text" @click="viewDetail(row)">{{ row.jobTitle }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="headcountType" label="编制类型" width="110" align="center">
          <template #default="{ row }">
            <el-tag :type="getHeadcountTagType(row.headcountType)" size="small" effect="plain">
              {{ getHeadcountText(row.headcountType) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="headcount" label="需求人数" width="100" align="center">
          <template #default="{ row }">
            <el-tag type="info" size="small">{{ row.headcount }}人</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="expectedOnboard" label="期望到岗" width="120" />
        <el-table-column label="预算薪资" width="160" align="center">
          <template #default="{ row }">
            <span class="salary-text">¥{{ formatSalary(row.salaryMin) }} - ¥{{ formatSalary(row.salaryMax) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="审批状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="getStatusTagType(row.status)" size="small" effect="light">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="initiator" label="发起人" width="130" />
        <el-table-column prop="updatedAt" label="最近处理时间" width="180">
          <template #default="{ row }">
            {{ formatTime(row.updatedAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="280" fixed="right" align="center">
          <template #default="{ row }">
            <el-button link type="primary" size="small" @click="viewDetail(row)">查看</el-button>
            <el-button
              link
              type="primary"
              size="small"
              @click="openDialog(row)"
              :disabled="!canEdit(row.status)"
            >编辑</el-button>
            <el-button
              link
              type="success"
              size="small"
              @click="handleSubmitApproval(row)"
              :disabled="!canSubmit(row.status)"
            >提交审批</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <RequirementDialog
      v-model="dialogVisible"
      :editing-row="editingRow"
      :existing-requirements="requirementList"
      @success="loadRequirements"
    />
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Plus } from '@element-plus/icons-vue'
import { listRequirements, submitApproval } from '@/api'
import RequirementDialog from './RequirementDialog.vue'

const router = useRouter()
const loading = ref(false)
const requirementList = ref([])
const departments = ['技术部', '产品部', '市场部', '运营部', '人力资源部', '财务部', '设计部']

const filters = reactive({
  department: '',
  status: '',
  keyword: ''
})

const dialogVisible = ref(false)
const editingRow = ref(null)

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
  return opt ? opt.label : status
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
  return opt ? opt.label : t
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

const formatSalary = (s) => {
  if (!s) return '0'
  if (s >= 1000) return (s / 1000) + 'K'
  return s
}

const formatTime = (t) => {
  if (!t) return '-'
  return t.replace('T', ' ').slice(0, 16)
}

const canEdit = (status) => status === 'draft' || status === 'returned'
const canSubmit = (status) => status === 'draft' || status === 'returned'

const loadRequirements = async () => {
  loading.value = true
  try {
    const res = await listRequirements(filters)
    requirementList.value = res.data || []
  } finally {
    loading.value = false
  }
}

const resetFilters = () => {
  filters.department = ''
  filters.status = ''
  filters.keyword = ''
  loadRequirements()
}

const openDialog = (row) => {
  editingRow.value = row || null
  dialogVisible.value = true
}

const viewDetail = (row) => {
  router.push({ path: '/requirements/detail', query: { id: row.id } })
}

const handleSubmitApproval = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确定要提交「${row.jobTitle}」用人需求审批吗？`,
      '提交确认',
      { type: 'warning' }
    )
  } catch {
    return
  }
  try {
    await submitApproval(row.id)
    ElMessage.success('已提交审批')
    loadRequirements()
  } catch (e) {}
}

onMounted(loadRequirements)
</script>

<style scoped>
.requirements-page {
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

.position-text {
  color: #409eff;
  cursor: pointer;
  font-weight: 500;
}

.position-text:hover {
  text-decoration: underline;
}

.salary-text {
  color: #f56c6c;
  font-weight: 500;
}
</style>
