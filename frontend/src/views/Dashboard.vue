<template>
  <div class="dashboard">
    <el-row :gutter="16" class="overview-row">
      <el-col :span="3" v-for="card in summaryCards" :key="card.key">
        <div class="summary-card" :style="{ background: card.bg }">
          <div class="sc-top">
            <div class="sc-count" :style="{ color: card.color }">{{ getCount(card.key) }}</div>
            <el-icon class="sc-icon" :style="{ color: card.color }"><component :is="card.icon" /></el-icon>
          </div>
          <div class="sc-label" style="color: #fff">{{ card.label }}</div>
          <div class="sc-sub" style="color: rgba(255,255,255,0.8)">{{ card.sub() }}</div>
        </div>
      </el-col>
    </el-row>

    <el-row :gutter="16" class="action-row">
      <el-col :span="12">
        <el-card class="block-card" shadow="never">
          <div class="block-header">
            <div class="block-title">
              <el-icon :size="18" color="#409eff"><Clock /></el-icon>
              <h3>近期面试安排</h3>
              <el-tag size="small" type="primary" effect="plain">{{ upcomingInterviews.length }} 场</el-tag>
            </div>
            <el-button type="primary" link @click="$router.push('/interviews')">全部面试</el-button>
          </div>
          <el-empty v-if="!loading.upcoming && upcomingInterviews.length === 0" description="近期暂无面试安排" />
          <div v-else v-loading="loading.upcoming" class="interview-list">
            <div v-for="iv in upcomingInterviews" :key="iv.id" class="interview-item" @click="openInterviewDetail(iv)">
              <div class="iv-avatar" :class="iv.status">
                {{ iv.candidateName?.charAt(0) || '-' }}
              </div>
              <div class="iv-content">
                <div class="iv-row1">
                  <span class="iv-name">{{ iv.candidateName }}</span>
                  <el-tag size="small" :type="getIvStatusTag(iv.status)">{{ getIvStatusText(iv.status) }}</el-tag>
                  <el-tag size="small" effect="plain" type="info">{{ getIvMethodText(iv.method) }}</el-tag>
                </div>
                <div class="iv-row2">
                  <span class="iv-job">{{ iv.jobTitle }}</span>
                  <span class="iv-time">
                    <el-icon><Calendar /></el-icon>
                    {{ formatTime(iv.interviewTime) }}
                  </span>
                </div>
                <div class="iv-row3">
                  面试官：{{ iv.interviewer }}
                  <span v-if="iv.location" class="iv-location">｜{{ iv.location }}</span>
                </div>
              </div>
              <div class="iv-actions">
                <el-button v-if="iv.status === 'scheduled'" size="small" type="success" @click.stop="quickReschedule(iv)">改期</el-button>
                <el-button v-if="iv.status === 'scheduled'" size="small" type="warning" @click.stop="quickComplete(iv)">完成</el-button>
                <el-button v-if="iv.status === 'completed'" size="small" type="primary" @click.stop="quickCreateOffer(iv)">发 Offer</el-button>
              </div>
            </div>
          </div>
        </el-card>
      </el-col>

      <el-col :span="12">
        <el-card class="block-card" shadow="never">
          <div class="block-header">
            <div class="block-title">
              <el-icon :size="18" color="#e6a23c"><Tickets /></el-icon>
              <h3>Offer 待办</h3>
              <el-tag size="small" type="warning" effect="plain">{{ pendingOffers.length }} 份</el-tag>
            </div>
            <el-button type="primary" link @click="$router.push('/offers')">全部 Offer</el-button>
          </div>
          <el-empty v-if="!loading.offers && pendingOffers.length === 0" description="暂无待处理 Offer" />
          <div v-else v-loading="loading.offers" class="offer-list">
            <div v-for="o in pendingOffers" :key="o.id" class="offer-item" @click="openOfferDetail(o)">
              <div class="of-avatar">{{ o.candidateName?.charAt(0) || '-' }}</div>
              <div class="of-content">
                <div class="of-row1">
                  <span class="of-name">{{ o.candidateName }}</span>
                  <el-tag size="small" :type="getOfferStatusTag(o.status)">{{ getOfferStatusText(o.status) }}</el-tag>
                </div>
                <div class="of-row2">
                  <span class="of-job">{{ o.jobTitle }}</span>
                  <span class="of-salary">¥{{ formatK(o.salaryMin) }}K - ¥{{ formatK(o.salaryMax) }}K</span>
                </div>
                <div class="of-row3">入职：{{ o.startDate }}｜负责人：{{ o.owner || '未指定' }}</div>
              </div>
              <div class="of-actions">
                <el-button v-if="o.status === 'pending'" size="small" type="primary" @click.stop="handleSendOffer(o)">发送</el-button>
                <el-button v-if="o.status === 'sent'" size="small" type="success" @click.stop="handleReplyOffer(o, true)">接受</el-button>
                <el-button v-if="o.status === 'sent'" size="small" type="warning" @click.stop="handleReplyOffer(o, false)">拒绝</el-button>
                <el-button v-if="!['accepted','rejected','withdrawn'].includes(o.status)" size="small" type="danger" link @click.stop="handleWithdrawOffer(o)">撤回</el-button>
              </div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="16" class="action-row">
      <el-col :span="12">
        <el-card class="block-card" shadow="never">
          <div class="block-header">
            <div class="block-title">
              <el-icon :size="18" color="#67c23a"><Document /></el-icon>
              <h3>待筛选申请</h3>
              <el-tag size="small" type="success" effect="plain">{{ pendingApps.length }} 份</el-tag>
            </div>
            <el-button type="primary" link @click="$router.push('/applications')">全部申请</el-button>
          </div>
          <el-empty v-if="!loading.apps && pendingApps.length === 0" description="所有申请已处理完成 🎉" />
          <div v-else v-loading="loading.apps" class="app-list">
            <div v-for="app in pendingApps" :key="app.id" class="app-item">
              <el-avatar :size="42" class="app-avatar">{{ app.resume?.candidateName?.charAt(0) || '-' }}</el-avatar>
              <div class="app-content">
                <div class="ap-row1">
                  <span class="ap-name">{{ app.resume?.candidateName }}</span>
                  <el-tag size="small" :type="getAppStatusTag(app.status)">{{ getAppStatusText(app.status) }}</el-tag>
                </div>
                <div class="ap-row2">
                  <span class="ap-job">{{ app.jobTitle }}</span>
                  <span class="ap-time">{{ formatTime(app.submittedAt) }}</span>
                </div>
                <div class="ap-row3">
                  <span>{{ app.resume?.contact }}</span>
                  <span v-if="app.resume?.yearsOfExperience">｜{{ app.resume.yearsOfExperience }} 年</span>
                  <span v-if="app.resume?.education">｜{{ app.resume.education }}</span>
                </div>
              </div>
              <div class="ap-actions">
                <el-button size="small" type="primary" @click="goCandidate(app)">候选人</el-button>
                <el-button size="small" type="success" @click="goMessage(app)">沟通</el-button>
                <el-dropdown size="small" trigger="click" @command="(v) => changeAppStatus(app, v)">
                  <el-button size="small">调状态<el-icon class="el-icon--right"><ArrowDown /></el-icon></el-button>
                  <template #dropdown>
                    <el-dropdown-menu>
                      <el-dropdown-item command="pending">待沟通</el-dropdown-item>
                      <el-dropdown-item command="interview">面试中</el-dropdown-item>
                      <el-dropdown-item command="rejected">已拒绝</el-dropdown-item>
                      <el-dropdown-item command="hired">已录用</el-dropdown-item>
                    </el-dropdown-menu>
                  </template>
                </el-dropdown>
              </div>
            </div>
          </div>
        </el-card>
      </el-col>

      <el-col :span="12">
        <el-card class="block-card" shadow="never">
          <div class="block-header">
            <div class="block-title">
              <el-icon :size="18" color="#f56c6c"><ChatDotRound /></el-icon>
              <h3>沟通提醒</h3>
              <el-tag size="small" type="danger" effect="plain">{{ msgReminders.length }} 条</el-tag>
            </div>
            <el-button type="primary" link @click="$router.push('/messages')">沟通中心</el-button>
          </div>
          <el-empty v-if="!loading.candidates && msgReminders.length === 0" description="暂无新沟通提醒" />
          <div v-else class="msg-list">
            <div v-for="m in msgReminders" :key="m.id" class="msg-item" @click="goMessage(m)">
              <el-avatar :size="40" class="msg-avatar">{{ m.candidateName?.charAt(0) || '-' }}</el-avatar>
              <div class="msg-content">
                <div class="msg-row1">
                  <span class="msg-name">{{ m.candidateName }}</span>
                  <span class="msg-time">{{ formatTime(m.lastTime) }}</span>
                </div>
                <div class="msg-row2">{{ m.jobTitle }}</div>
                <div class="msg-row3">{{ m.lastMessage || '点击进入沟通' }}</div>
              </div>
              <el-icon class="msg-arrow"><ArrowRight /></el-icon>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="16" class="action-row">
      <el-col :span="12">
        <el-card class="block-card" shadow="never">
          <div class="block-header">
            <div class="block-title">
              <el-icon :size="18" color="#409eff"><Briefcase /></el-icon>
              <h3>招聘中职位</h3>
              <el-tag size="small" type="primary" effect="plain">{{ openJobs.length }} 个</el-tag>
            </div>
            <el-button type="primary" link @click="$router.push('/jobs')">职位管理</el-button>
          </div>
          <el-empty v-if="!loading.jobs && openJobs.length === 0" description="暂无招聘中职位" />
          <div v-else v-loading="loading.jobs" class="job-list">
            <div v-for="j in openJobs" :key="j.id" class="job-item">
              <div class="jb-main">
                <div class="jb-title">{{ j.title }}</div>
                <div class="jb-meta">
                  <span>{{ j.department }}</span>
                  <span>｜{{ j.location }}</span>
                  <span>｜{{ j.type }}</span>
                  <span v-if="j.salaryMin">｜¥{{ formatK(j.salaryMin) }}K+</span>
                </div>
              </div>
              <div class="jb-right">
                <div class="jb-count">
                  <span class="jc-num">{{ getAppCountForJob(j.id) }}</span>
                  <span class="jc-lbl">份申请</span>
                </div>
                <el-divider direction="vertical" />
                <el-button type="primary" link @click="$router.push('/jobs')">查看</el-button>
              </div>
            </div>
          </div>
        </el-card>
      </el-col>

      <el-col :span="12">
        <el-card class="block-card" shadow="never">
          <div class="block-header">
            <div class="block-title">
              <el-icon :size="18" color="#67c23a"><UserFilled /></el-icon>
              <h3>候选人进展动态</h3>
            </div>
            <el-button type="primary" link @click="$router.push('/candidates')">全部候选人</el-button>
          </div>
          <el-empty v-if="!loading.candidates && recentCandidates.length === 0" description="暂无候选人动态" />
          <div v-else class="timeline-wrap">
            <el-timeline>
              <el-timeline-item
                v-for="(t, i) in recentCandidates.slice(0, 8)"
                :key="i"
                :timestamp="t.time"
                :type="t.type"
                placement="top"
              >
                <div class="tl-content">
                  <span class="tl-name">{{ t.name }}</span>
                  <span class="tl-action">{{ t.action }}</span>
                  <span class="tl-job">「{{ t.job }}」</span>
                </div>
              </el-timeline-item>
            </el-timeline>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-dialog v-model="interviewDetailVisible" title="面试详情" width="640px" destroy-on-close>
      <InterviewDetailContent v-if="currentInterview" :interview="currentInterview" @updated="reloadAll" />
    </el-dialog>

    <el-dialog v-model="offerDetailVisible" title="Offer 详情" width="640px" destroy-on-close>
      <OfferDetailContent v-if="currentOffer" :offer="currentOffer" @updated="reloadAll" />
    </el-dialog>

    <el-dialog v-model="completeVisible" title="标记面试完成" width="520px" destroy-on-close>
      <el-form :model="completeForm" label-width="90px">
        <el-form-item label="面试结果">
          <el-radio-group v-model="completeForm.passed">
            <el-radio :value="true">通过</el-radio>
            <el-radio :value="false">未通过</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="整体评价">
          <el-rate v-model="completeForm.rating" />
        </el-form-item>
        <el-form-item label="面试反馈">
          <el-input
            v-model="completeForm.feedback"
            type="textarea"
            :rows="4"
            placeholder="请填写面试表现、优缺点、录用建议等"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="completeVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="confirmComplete">确认提交</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="rescheduleVisible" title="面试改期" width="520px" destroy-on-close>
      <el-form :model="rescheduleForm" label-width="90px">
        <el-form-item label="新面试时间">
          <el-date-picker
            v-model="rescheduleForm.interviewTime"
            type="datetime"
            style="width: 100%"
            format="YYYY-MM-DD HH:mm"
            value-format="YYYY-MM-DDTHH:mm:ssZ"
          />
        </el-form-item>
        <el-form-item label="改期原因">
          <el-input v-model="rescheduleForm.reason" type="textarea" :rows="2" placeholder="请输入改期原因" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="rescheduleVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="confirmReschedule">确认改期</el-button>
      </template>
    </el-dialog>

    <OfferDialog
      v-model="offerDialogVisible"
      :applications="offerApplications"
      :candidate-name="offerCandidateName"
      :pre-select-application-id="offerPreSelectAppId"
      @success="handleOfferSuccess"
    />
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, defineAsyncComponent, computed, markRaw } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Clock, Calendar, Tickets, Document, ChatDotRound, Briefcase, UserFilled,
  ArrowDown, ArrowRight, DataAnalysis, TrendCharts, Bell, Medal
} from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import {
  getStats, listJobs, listApplications, listInterviews, listOffers, listCandidates,
  getInterview, completeInterview, rescheduleInterview, sendOffer, replyOffer, withdrawOffer,
  updateApplicationStatus, listMessages, getApplication
} from '@/api'

const InterviewDetailContent = defineAsyncComponent(() => import('./InterviewDetailContent.vue'))
const OfferDetailContent = defineAsyncComponent(() => import('./OfferDetailContent.vue'))
const OfferDialog = defineAsyncComponent(() => import('./OfferDialog.vue'))

const router = useRouter()

const loading = reactive({
  jobs: false, apps: false, upcoming: false, offers: false, candidates: false
})
const submitting = ref(false)

const stats = ref({
  totalJobs: 0, openJobs: 0, totalApplications: 0, totalCandidates: 0,
  totalInterviews: 0, pendingInterviews: 0, newThisWeek: 0,
  totalOffers: 0, pendingOffers: 0, acceptedOffers: 0
})
const openJobs = ref([])
const allApps = ref([])
const allInterviews = ref([])
const allOffers = ref([])
const allCandidates = ref([])

const summaryCards = [
  { key: 'openJobs', label: '在招职位', icon: markRaw(Briefcase), color: '#409eff', bg: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)', sub: () => `共 ${stats.value.totalJobs||0} 个职位` },
  { key: 'totalApplications', label: '累计申请', icon: markRaw(Document), color: '#fff', bg: 'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)', sub: () => `本周新增 ${stats.value.newThisWeek||0}` },
  { key: 'totalCandidates', label: '候选人池', icon: markRaw(UserFilled), color: '#fff', bg: 'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)', sub: () => '含历史数据' },
  { key: 'pendingInterviews', label: '待面试', icon: markRaw(Calendar), color: '#fff', bg: 'linear-gradient(135deg, #fa709a 0%, #fee140 100%)', sub: () => `累计 ${stats.value.totalInterviews||0} 场` },
  { key: 'pendingOffers', label: '待发 Offer', icon: markRaw(Tickets), color: '#fff', bg: 'linear-gradient(135deg, #a8edea 0%, #fed6e3 100%)', sub: () => `累计 Offer ${stats.value.totalOffers||0} 份` },
  { key: 'acceptedOffers', label: '成功录用', icon: markRaw(Medal), color: '#fff', bg: 'linear-gradient(135deg, #43e97b 0%, #38f9d7 100%)', sub: () => '入职成功' },
  { key: 'newThisWeek', label: '本周新增', icon: markRaw(TrendCharts), color: '#fff', bg: 'linear-gradient(135deg, #30cfd0 0%, #330867 100%)', sub: () => '最近 7 天' },
  { key: 'totalInterviews', label: '累计面试', icon: markRaw(Bell), color: '#fff', bg: 'linear-gradient(135deg, #ff9a9e 0%, #fecfef 100%)', sub: () => `待面试 ${stats.value.pendingInterviews||0}` }
]
const getCount = (k) => stats.value[k] || 0

const pendingApps = computed(() => allApps.value.filter(a => a.status === 'submitted' || a.status === 'pending').slice(0, 6))
const upcomingInterviews = computed(() => {
  const now = dayjs()
  const list = allInterviews.value.filter(iv => {
    if (iv.status === 'scheduled') return true
    return now.diff(dayjs(iv.interviewTime), 'day') <= 7
  })
  return list.sort((a, b) => new Date(a.interviewTime) - new Date(b.interviewTime)).slice(0, 6)
})
const pendingOffers = computed(() => {
  return allOffers.value.filter(o => ['pending','sent'].includes(o.status)).slice(0, 6)
})
const msgReminders = computed(() => {
  const arr = []
  const seen = new Set()
  allApps.value.forEach(a => {
    const key = a.resume?.contact + '|' + a.id
    if (seen.has(key)) return
    seen.add(key)
    arr.push({
      id: a.id,
      applicationId: a.id,
      candidateName: a.resume?.candidateName,
      contact: a.resume?.contact,
      jobTitle: a.jobTitle,
      lastTime: a.submittedAt,
      lastMessage: a.latestMessage || '点击进入查看沟通记录'
    })
  })
  return arr.slice(0, 6)
})
const recentCandidates = computed(() => {
  const arr = []
  allApps.value.slice().sort((a, b) => new Date(b.submittedAt) - new Date(a.submittedAt)).slice(0, 12).forEach(a => {
    const typeMap = { submitted: 'primary', pending: 'warning', interview: 'success', rejected: 'danger', hired: 'success' }
    const actionMap = { submitted: '投递申请', pending: '进入待沟通', interview: '进入面试', rejected: '已拒绝', hired: '已录用' }
    arr.push({
      time: formatTime(a.submittedAt),
      type: typeMap[a.status] || 'primary',
      name: a.resume?.candidateName || '',
      job: a.jobTitle,
      action: actionMap[a.status] || '状态更新'
    })
  })
  return arr
})
const getAppCountForJob = (jid) => allApps.value.filter(a => a.jobId === jid).length

const formatTime = (t) => {
  if (!t) return '-'
  return dayjs(t).format('YYYY-MM-DD HH:mm')
}
const formatK = (n) => n ? (n / 1000).toFixed(0) : 0

const getIvStatusText = s => ({ scheduled: '待进行', completed: '已完成', cancelled: '已取消' }[s] || s)
const getIvStatusTag = s => ({ scheduled: 'warning', completed: 'success', cancelled: 'info' }[s] || '')
const getIvMethodText = m => ({ onsite: '现场', online: '视频', phone: '电话' }[m] || m)
const getOfferStatusText = s => ({ pending: '待发送', sent: '已发送', accepted: '已接受', rejected: '已拒绝', withdrawn: '已撤回' }[s] || s)
const getOfferStatusTag = s => ({ pending: 'warning', sent: 'primary', accepted: 'success', rejected: 'danger', withdrawn: 'info' }[s] || '')
const getAppStatusText = s => ({ submitted: '已投递', pending: '待沟通', interview: '面试中', rejected: '已拒绝', hired: '已录用' }[s] || s)
const getAppStatusTag = s => ({ submitted: '', pending: 'warning', interview: 'success', rejected: 'danger', hired: 'success' }[s] || 'info')

const reloadAll = () => {
  loadStats()
  loadJobs()
  loadApps()
  loadInterviews()
  loadOffers()
  loadCandidates()
}

const loadStats = async () => {
  try {
    const res = await getStats()
    stats.value = res?.data || {}
  } catch (e) { console.error(e) }
}
const loadJobs = async () => {
  loading.jobs = true
  try {
    const res = await listJobs()
    const list = res?.data || []
    openJobs.value = list.filter(j => j.status === 'open').slice(0, 6)
  } finally { loading.jobs = false }
}
const loadApps = async () => {
  loading.apps = true
  try {
    const res = await listApplications()
    allApps.value = res?.data || []
  } finally { loading.apps = false }
}
const loadInterviews = async () => {
  loading.upcoming = true
  try {
    const res = await listInterviews()
    allInterviews.value = res?.data || []
  } finally { loading.upcoming = false }
}
const loadOffers = async () => {
  loading.offers = true
  try {
    const res = await listOffers()
    allOffers.value = res?.data || []
  } finally { loading.offers = false }
}
const loadCandidates = async () => {
  loading.candidates = true
  try {
    const res = await listCandidates()
    allCandidates.value = res?.data || []
  } finally { loading.candidates = false }
}

const interviewDetailVisible = ref(false)
const currentInterview = ref(null)
const openInterviewDetail = async (iv) => {
  try {
    const res = await getInterview(iv.id)
    currentInterview.value = res?.data || iv
    interviewDetailVisible.value = true
  } catch (e) {
    currentInterview.value = iv
    interviewDetailVisible.value = true
  }
}

const offerDetailVisible = ref(false)
const currentOffer = ref(null)
const openOfferDetail = (o) => {
  currentOffer.value = o
  offerDetailVisible.value = true
}

const completeVisible = ref(false)
const completeForm = reactive({ rating: 4, passed: true, feedback: '' })
let completingIv = null
const quickComplete = (iv) => {
  completingIv = iv
  completeForm.rating = 4
  completeForm.passed = true
  completeForm.feedback = ''
  completeVisible.value = true
}
const confirmComplete = async () => {
  if (!completingIv) return
  submitting.value = true
  try {
    await completeInterview(completingIv.id, { ...completeForm })
    ElMessage.success('面试已标记完成')
    completeVisible.value = false
    reloadAll()
  } catch (e) {
    ElMessage.error('操作失败：' + (e?.response?.data?.error || e.message))
  } finally { submitting.value = false }
}

const rescheduleVisible = ref(false)
const rescheduleForm = reactive({ interviewTime: '', reason: '' })
let reschedulingIv = null
const quickReschedule = (iv) => {
  reschedulingIv = iv
  rescheduleForm.interviewTime = ''
  rescheduleForm.reason = ''
  rescheduleVisible.value = true
}
const confirmReschedule = async () => {
  if (!reschedulingIv || !rescheduleForm.interviewTime) {
    ElMessage.warning('请选择新的面试时间')
    return
  }
  submitting.value = true
  try {
    await rescheduleInterview(reschedulingIv.id, rescheduleForm)
    ElMessage.success('面试改期成功')
    rescheduleVisible.value = false
    reloadAll()
  } catch (e) {
    ElMessage.error('操作失败：' + (e?.response?.data?.error || e.message))
  } finally { submitting.value = false }
}

const offerDialogVisible = ref(false)
const offerApplications = ref([])
const offerCandidateName = ref('')
const offerPreSelectAppId = ref('')
const quickCreateOffer = async (iv) => {
  offerCandidateName.value = iv.candidateName
  offerPreSelectAppId.value = iv.applicationId
  offerApplications.value = []
  if (iv.applicationId) {
    try {
      const res = await getApplication(iv.applicationId)
      if (res?.data) offerApplications.value = [res.data]
    } catch (e) { console.error(e) }
  }
  offerDialogVisible.value = true
}
const handleOfferSuccess = () => {
  ElMessage.success('Offer 创建成功')
  reloadAll()
}

const handleSendOffer = async (o) => {
  try {
    await sendOffer(o.id)
    ElMessage.success('Offer 已发送')
    reloadAll()
  } catch (e) {
    ElMessage.error('操作失败：' + (e?.response?.data?.error || e.message))
  }
}
const handleReplyOffer = async (o, accepted) => {
  const action = accepted ? '接受' : '拒绝'
  try {
    await ElMessageBox.confirm(`确认候选人已${action}该 Offer？`, `${action} Offer`, { type: 'warning' })
    await replyOffer(o.id, { accepted, feedback: `候选人已${action}Offer` })
    ElMessage.success(`候选人已${action}`)
    reloadAll()
  } catch (e) {
    if (e !== 'cancel') ElMessage.error('操作失败：' + (e?.response?.data?.error || e.message))
  }
}
const handleWithdrawOffer = async (o) => {
  try {
    const { value } = await ElMessageBox.prompt('请输入撤回原因', '撤回 Offer', {
      inputPlaceholder: '请说明撤回原因',
      confirmButtonText: '确认撤回', cancelButtonText: '取消', type: 'warning'
    })
    await withdrawOffer(o.id, { reason: value || '公司撤回' })
    ElMessage.success('Offer 已撤回')
    reloadAll()
  } catch (e) {
    if (e !== 'cancel') ElMessage.error('操作失败：' + (e?.response?.data?.error || e.message))
  }
}

const goCandidate = (app) => {
  const name = encodeURIComponent(app.resume?.candidateName || '')
  router.push(`/candidates/detail?contact=${app.resume?.contact || ''}&name=${name}`)
}
const goMessage = (row) => {
  const appId = row.applicationId || row.id
  router.push({ path: '/messages', query: { applicationId: appId } })
}
const changeAppStatus = async (app, newStatus) => {
  try {
    await updateApplicationStatus(app.id, newStatus)
    ElMessage.success('状态已更新')
    reloadAll()
  } catch (e) {
    ElMessage.error('操作失败：' + (e?.response?.data?.error || e.message))
  }
}

onMounted(reloadAll)
</script>

<style scoped>
.dashboard {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.overview-row, .action-row { margin-bottom: 0; }

.summary-card {
  border-radius: 10px;
  padding: 16px 18px;
  color: #fff;
  box-shadow: 0 4px 14px rgba(0,0,0,0.06);
}
.sc-top { display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 8px; }
.sc-count { font-size: 28px; font-weight: 700; line-height: 1; }
.sc-icon { font-size: 22px; opacity: 0.9; }
.sc-label { font-size: 13px; margin-bottom: 4px; }
.sc-sub { font-size: 12px; }

.block-card { border-radius: 10px; }
.block-header {
  display: flex; align-items: center; justify-content: space-between; margin-bottom: 16px;
  padding-bottom: 12px; border-bottom: 1px solid #f0f0f0;
}
.block-title { display: flex; align-items: center; gap: 8px; }
.block-title h3 { margin: 0; font-size: 16px; font-weight: 600; color: #303133; }

/* interview list */
.interview-list, .offer-list, .app-list, .msg-list, .job-list { display: flex; flex-direction: column; gap: 10px; }

.interview-item, .offer-item, .app-item, .msg-item {
  display: flex; align-items: center; gap: 12px;
  padding: 12px; border: 1px solid #f0f0f0; border-radius: 8px;
  transition: all 0.2s; cursor: pointer;
}
.interview-item:hover, .offer-item:hover, .app-item:hover, .msg-item:hover {
  border-color: #409eff; background: #fafcff;
}
.iv-avatar, .of-avatar {
  width: 42px; height: 42px; border-radius: 50%;
  background: linear-gradient(135deg, #67c23a, #85ce61);
  color: #fff; display: flex; align-items: center; justify-content: center;
  font-weight: 600; font-size: 16px; flex-shrink: 0;
}
.iv-avatar.completed { background: linear-gradient(135deg, #409eff, #79bbff); }
.iv-avatar.cancelled { background: linear-gradient(135deg, #909399, #c0c4cc); }

.iv-content, .of-content, .app-content, .msg-content { flex: 1; min-width: 0; }
.iv-row1, .of-row1, .ap-row1, .msg-row1 {
  display: flex; align-items: center; gap: 8px; margin-bottom: 4px;
}
.iv-name, .of-name, .ap-name, .msg-name { font-weight: 600; color: #303133; font-size: 14px; }
.iv-row2, .of-row2, .ap-row2, .msg-row2 {
  display: flex; align-items: center; gap: 12px; font-size: 13px; color: #606266; margin-bottom: 4px;
}
.iv-time, .ap-time, .msg-time { display: inline-flex; align-items: center; gap: 4px; color: #909399; font-size: 12px; }
.iv-row3, .of-row3, .ap-row3, .msg-row3 { font-size: 12px; color: #909399; }
.iv-job, .of-job, .ap-job { color: #409eff; font-weight: 500; }
.of-salary { color: #f56c6c; font-weight: 600; }

.iv-actions, .of-actions, .ap-actions { display: flex; gap: 6px; flex-shrink: 0; }

.app-avatar { background: linear-gradient(135deg, #409eff, #79bbff); color: #fff; }
.msg-avatar { background: linear-gradient(135deg, #e6a23c, #f3d19e); color: #fff; }
.msg-arrow { color: #c0c4cc; }
.msg-row3 {
  overflow: hidden; text-overflow: ellipsis; white-space: nowrap; max-width: 260px;
  color: #606266;
}
.msg-row1 { justify-content: space-between; }

.job-item {
  display: flex; align-items: center; justify-content: space-between;
  padding: 12px; border: 1px solid #f0f0f0; border-radius: 8px;
  transition: all 0.2s;
}
.job-item:hover { border-color: #67c23a; background: #f5fff5; }
.jb-title { font-size: 15px; font-weight: 600; color: #303133; margin-bottom: 4px; }
.jb-meta { font-size: 12px; color: #606266; display: flex; gap: 4px; }
.jb-right { display: flex; align-items: center; gap: 10px; }
.jc-num { font-size: 20px; font-weight: 700; color: #e6a23c; margin-right: 4px; }
.jc-lbl { font-size: 12px; color: #909399; }
.jb-count { padding: 0 12px; text-align: center; }

.timeline-wrap { padding: 4px 8px 0; }
.tl-content { font-size: 13px; color: #606266; line-height: 1.6; }
.tl-name { font-weight: 600; color: #303133; margin-right: 4px; }
.tl-action { margin-right: 4px; }
.tl-job { color: #409eff; }
</style>
