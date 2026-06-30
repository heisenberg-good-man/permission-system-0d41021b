<template>
  <div class="stats-page">
    <el-row :gutter="16">
      <el-col :span="3" v-for="card in statCards" :key="card.key">
        <el-card class="stat-card" shadow="never" :body-style="{ padding: '16px' }">
          <div class="stat-card-inner">
            <div class="stat-icon" :style="{ background: card.bg }">
              <el-icon :size="22"><component :is="card.icon" /></el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-label">{{ card.label }}</div>
              <div class="stat-value" :style="{ color: card.color }">{{ loading ? '--' : statData[card.key] }}</div>
              <div v-if="card.sub" class="stat-sub">{{ card.sub(statData) }}</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="16" style="margin-top: 16px">
      <el-col :span="10">
        <el-card class="chart-card" shadow="never">
          <div class="card-header">
            <h3 class="card-title">申请状态分布</h3>
          </div>
          <div ref="pieChartRef" class="chart-box"></div>
        </el-card>
      </el-col>

      <el-col :span="14">
        <el-card class="chart-card" shadow="never">
          <div class="card-header">
            <h3 class="card-title">各职位申请数 Top5</h3>
          </div>
          <div ref="barChartRef" class="chart-box"></div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="16" style="margin-top: 16px">
      <el-col :span="12">
        <el-card class="chart-card" shadow="never">
          <div class="card-header">
            <h3 class="card-title">最新申请动态</h3>
          </div>
          <el-table :data="recentApplications" stripe style="width: 100%">
            <el-table-column label="候选人" width="130">
              <template #default="{ row }">
                <div class="mini-candidate">
                  <el-avatar :size="32">{{ row.resume.candidateName.charAt(0) }}</el-avatar>
                  <span>{{ row.resume.candidateName }}</span>
                </div>
              </template>
            </el-table-column>
            <el-table-column prop="jobTitle" label="应聘职位" min-width="150" />
            <el-table-column label="状态" width="100" align="center">
              <template #default="{ row }">
                <el-tag :type="getStatusTagType(row.status)" size="small">{{ getStatusText(row.status) }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column label="时间" width="160">
              <template #default="{ row }">{{ formatTime(row.submittedAt) }}</template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>

      <el-col :span="12">
        <el-card class="chart-card" shadow="never">
          <div class="card-header">
            <h3 class="card-title">招聘职位速览</h3>
          </div>
          <el-table :data="hotJobs" stripe style="width: 100%">
            <el-table-column prop="title" label="岗位名称" min-width="160" />
            <el-table-column prop="department" label="部门" width="100" />
            <el-table-column prop="location" label="地点" width="80" />
            <el-table-column label="薪资" width="140">
              <template #default="{ row }">
                <span style="color: #f56c6c">¥{{ row.salaryMin }}-{{ row.salaryMax }}</span>
              </template>
            </el-table-column>
            <el-table-column label="状态" width="90" align="center">
              <template #default="{ row }">
                <el-tag :type="getJobStatusTagType(row.status)" size="small">{{ getJobStatusText(row.status) }}</el-tag>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, nextTick, markRaw } from 'vue'
import * as echarts from 'echarts'
import dayjs from 'dayjs'
import {
  DataAnalysis, Briefcase, User, TrendCharts, Bell, UserFilled, Calendar, Tickets, Medal
} from '@element-plus/icons-vue'
import { getStats, listApplications, listJobs } from '@/api'

const loading = ref(true)
const pieChartRef = ref(null)
const barChartRef = ref(null)
let pieChart = null
let barChart = null

const statData = reactive({
  totalJobs: 0,
  openJobs: 0,
  totalApplications: 0,
  totalCandidates: 0,
  totalInterviews: 0,
  pendingInterviews: 0,
  newThisWeek: 0,
  totalOffers: 0,
  pendingOffers: 0,
  acceptedOffers: 0
})

const statCards = [
  {
    key: 'totalJobs',
    label: '职位总数',
    icon: markRaw(Briefcase),
    color: '#409eff',
    bg: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
    sub: (d) => `招聘中 ${d.openJobs} 个`
  },
  {
    key: 'openJobs',
    label: '在招职位',
    icon: markRaw(Bell),
    color: '#67c23a',
    bg: 'linear-gradient(135deg, #11998e 0%, #38ef7d 100%)',
    sub: (d) => `占比 ${d.totalJobs ? Math.round(d.openJobs / d.totalJobs * 100) : 0}%`
  },
  {
    key: 'totalCandidates',
    label: '候选人总数',
    icon: markRaw(UserFilled),
    color: '#e6a23c',
    bg: 'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)',
    sub: (d) => `累计申请 ${d.totalApplications} 份`
  },
  {
    key: 'pendingInterviews',
    label: '待面试',
    icon: markRaw(Calendar),
    color: '#f56c6c',
    bg: 'linear-gradient(135deg, #fa709a 0%, #fee140 100%)',
    sub: (d) => `累计面试 ${d.totalInterviews} 场`
  },
  {
    key: 'totalApplications',
    label: '累计申请',
    icon: markRaw(User),
    color: '#409eff',
    bg: 'linear-gradient(135deg, #2193b0 0%, #6dd5ed 100%)',
    sub: () => '含全状态'
  },
  {
    key: 'newThisWeek',
    label: '本周新增',
    icon: markRaw(TrendCharts),
    color: '#67c23a',
    bg: 'linear-gradient(135deg, #11998e 0%, #38ef7d 100%)',
    sub: () => '最近 7 天'
  },
  {
    key: 'pendingOffers',
    label: '待发 Offer',
    icon: markRaw(Tickets),
    color: '#e6a23c',
    bg: 'linear-gradient(135deg, #f6d365 0%, #fda085 100%)',
    sub: (d) => `累计 Offer ${d.totalOffers} 份`
  },
  {
    key: 'acceptedOffers',
    label: '已接受 Offer',
    icon: markRaw(Medal),
    color: '#67c23a',
    bg: 'linear-gradient(135deg, #43e97b 0%, #38f9d7 100%)',
    sub: () => '成功入职'
  }
]

const recentApplications = ref([])
const hotJobs = ref([])

const statusOptions = [
  { value: 'submitted', label: '已投递' },
  { value: 'pending', label: '待沟通' },
  { value: 'interview', label: '面试中' },
  { value: 'rejected', label: '已拒绝' },
  { value: 'hired', label: '已录用' }
]

const getStatusText = (s) => (statusOptions.find(x => x.value === s) || {}).label || s
const getStatusTagType = (s) => {
  const m = { submitted: 'info', pending: 'warning', interview: 'primary', rejected: 'danger', hired: 'success' }
  return m[s] || 'info'
}
const getJobStatusText = (s) => ({ open: '招聘中', paused: '已暂停', closed: '已关闭' }[s] || s)
const getJobStatusTagType = (s) => ({ open: 'success', paused: 'warning', closed: 'info' }[s] || 'info')
const formatTime = (t) => dayjs(t).format('YYYY-MM-DD HH:mm')

const initCharts = () => {
  if (pieChartRef.value) {
    pieChart = echarts.init(pieChartRef.value)
  }
  if (barChartRef.value) {
    barChart = echarts.init(barChartRef.value)
  }
  window.addEventListener('resize', handleResize)
}

const handleResize = () => {
  pieChart && pieChart.resize()
  barChart && barChart.resize()
}

const updatePieChart = (byStatus) => {
  if (!pieChart) return
  const pieColors = ['#909399', '#e6a23c', '#409eff', '#f56c6c', '#67c23a']
  const data = statusOptions.map((opt, i) => ({
    value: byStatus[opt.value] || 0,
    name: opt.label,
    itemStyle: { color: pieColors[i] }
  })).filter(d => d.value > 0)

  pieChart.setOption({
    tooltip: { trigger: 'item', formatter: '{b}: {c} ({d}%)' },
    legend: { orient: 'vertical', right: 10, top: 'center', itemWidth: 10, itemHeight: 10 },
    series: [{
      type: 'pie',
      radius: ['45%', '70%'],
      center: ['40%', '50%'],
      avoidLabelOverlap: true,
      label: { show: false },
      emphasis: { label: { show: true, fontSize: 14, fontWeight: 'bold' } },
      data
    }]
  })
}

const updateBarChart = (apps, jobs) => {
  if (!barChart) return
  const jobCount = {}
  apps.forEach(a => {
    jobCount[a.jobId] = (jobCount[a.jobId] || 0) + 1
  })
  const sorted = jobs
    .map(j => ({ name: j.title, count: jobCount[j.id] || 0, id: j.id }))
    .sort((a, b) => b.count - a.count)
    .slice(0, 5)

  barChart.setOption({
    tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' } },
    grid: { left: 20, right: 30, bottom: 40, top: 20, containLabel: true },
    xAxis: {
      type: 'category',
      data: sorted.map(s => s.name),
      axisLabel: { interval: 0, rotate: 15, fontSize: 11 }
    },
    yAxis: { type: 'value', minInterval: 1 },
    series: [{
      type: 'bar',
      data: sorted.map(s => s.count),
      itemStyle: {
        borderRadius: [6, 6, 0, 0],
        color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
          { offset: 0, color: '#667eea' },
          { offset: 1, color: '#764ba2' }
        ])
      },
      barWidth: '40%',
      label: { show: true, position: 'top', color: '#666', fontSize: 12, fontWeight: 'bold' }
    }]
  })
}

const loadAll = async () => {
  loading.value = true
  try {
    const [statsRes, appsRes, jobsRes] = await Promise.all([
      getStats(), listApplications({}), listJobs({})
    ])
    const stats = statsRes.data || {}
    statData.totalJobs = stats.totalJobs || 0
    statData.openJobs = stats.openJobs || 0
    statData.totalApplications = stats.totalApplications || 0
    statData.totalCandidates = stats.totalCandidates || 0
    statData.totalInterviews = stats.totalInterviews || 0
    statData.pendingInterviews = stats.pendingInterviews || 0
    statData.newThisWeek = stats.newThisWeek || 0
    statData.totalOffers = stats.totalOffers || 0
    statData.pendingOffers = stats.pendingOffers || 0
    statData.acceptedOffers = stats.acceptedOffers || 0

    const apps = appsRes.data || []
    const jobs = jobsRes.data || []

    recentApplications.value = [...apps]
      .sort((a, b) => new Date(b.submittedAt) - new Date(a.submittedAt))
      .slice(0, 6)

    hotJobs.value = jobs.slice(0, 6)

    await nextTick()
    if (!pieChart || !barChart) initCharts()
    updatePieChart(stats.applicationsByStatus || {})
    updateBarChart(apps, jobs)
  } finally {
    loading.value = false
  }
}

onMounted(loadAll)
</script>

<style scoped>
.stats-page {
  display: flex;
  flex-direction: column;
}

.stat-card {
  border-radius: 8px;
}

.stat-card-inner {
  display: flex;
  align-items: center;
  gap: 16px;
}

.stat-icon {
  width: 56px;
  height: 56px;
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  flex-shrink: 0;
}

.stat-content {
  flex: 1;
  min-width: 0;
}

.stat-label {
  font-size: 13px;
  color: #909399;
}

.stat-value {
  font-size: 28px;
  font-weight: 700;
  margin: 4px 0;
  line-height: 1.2;
}

.stat-sub {
  font-size: 12px;
  color: #c0c4cc;
}

.chart-card {
  border-radius: 8px;
}

.card-header {
  margin-bottom: 16px;
}

.card-title {
  font-size: 15px;
  font-weight: 600;
  color: #1f2937;
  margin: 0;
}

.chart-box {
  width: 100%;
  height: 280px;
}

.mini-candidate {
  display: flex;
  align-items: center;
  gap: 8px;
}

.mini-candidate .el-avatar {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: #fff;
  font-weight: 600;
  font-size: 13px;
}
</style>
