<template>
  <div class="candidate-detail-page">
    <div class="page-header">
      <el-button @click="goBack" :icon="ArrowLeft">返回列表</el-button>
      <h2 class="page-title">候选人详情</h2>
    </div>

    <el-row v-loading="loading" :gutter="16">
      <el-col :span="9">
        <el-card class="profile-card" shadow="never">
          <div class="profile-header" v-if="candidate">
            <el-avatar :size="72" class="avatar">{{ candidate.candidateName.charAt(0) }}</el-avatar>
            <div class="profile-info">
              <div class="profile-name">
                {{ candidate.candidateName }}
                <el-tag :type="getStatusTagType(candidate.latestStatus)" size="large" effect="light" style="margin-left: 8px">
                  {{ getStatusText(candidate.latestStatus) }}
                </el-tag>
              </div>
              <div class="profile-contact">
                <el-icon><Phone /></el-icon>
                <span>{{ candidate.contact }}</span>
              </div>
              <div class="profile-position">
                <el-icon><Briefcase /></el-icon>
                <span>期望：{{ candidate.targetPosition }} · {{ candidate.yearsOfExp }}年经验</span>
              </div>
            </div>
          </div>

          <el-divider />

          <div class="profile-section" v-if="candidate">
            <h4 class="section-title">技能标签</h4>
            <div class="skill-tags">
              <el-tag
                v-for="skill in candidate.skills"
                :key="skill"
                size="small"
                type="primary"
                effect="plain"
                style="margin: 4px 6px 4px 0"
              >
                {{ skill }}
              </el-tag>
            </div>
          </div>

          <el-divider />

          <div class="profile-section" v-if="candidate">
            <h4 class="section-title">简历摘要</h4>
            <div class="summary-text">{{ candidate.summary }}</div>
          </div>

          <el-divider />

          <div class="profile-section">
            <h4 class="section-title">状态变更</h4>
            <div class="status-buttons">
              <el-button
                v-for="item in changeableStatusOptions"
                :key="item.value"
                :type="candidate?.latestStatus === item.value ? 'primary' : 'default'"
                size="default"
                @click="changeStatus(item.value)"
              >
                {{ item.label }}
              </el-button>
            </div>
          </div>

          <el-divider />

          <div class="profile-section">
            <h4 class="section-title">快捷操作</h4>
            <div class="quick-actions">
              <el-button type="success" @click="goToMessage" :icon="ChatDotRound">
                发起沟通
              </el-button>
              <el-button type="warning" @click="openInterviewDialog" :icon="Calendar">
                安排面试
              </el-button>
            </div>
          </div>
        </el-card>
      </el-col>

      <el-col :span="15">
        <el-card class="timeline-card" shadow="never">
          <div class="card-header">
            <h3 class="card-title">投递记录</h3>
          </div>
          <el-timeline v-if="candidate && candidate.applications.length > 0">
            <el-timeline-item
              v-for="app in candidate.applications"
              :key="app.id"
              :timestamp="formatTime(app.submittedAt)"
              placement="top"
            >
              <el-card shadow="never" class="timeline-item-card">
                <div class="timeline-item-header">
                  <el-tag type="info" size="small" effect="plain">{{ app.jobTitle }}</el-tag>
                  <el-tag :type="getStatusTagType(app.status)" size="small" effect="light">
                    {{ getStatusText(app.status) }}
                  </el-tag>
                </div>
                <div class="timeline-item-actions">
                  <el-button link type="primary" size="small" @click="goToApplication(app.id)">
                    查看申请详情
                  </el-button>
                  <el-button link type="success" size="small" @click="goToMessageById(app.id)">
                    去沟通
                  </el-button>
                  <el-button link type="warning" size="small" @click="openInterviewDialog(app)" :disabled="app.status === 'rejected' || app.status === 'hired'">
                    安排面试
                  </el-button>
                </div>
              </el-card>
            </el-timeline-item>
          </el-timeline>
          <el-empty v-else description="暂无投递记录" :image-size="100" />
        </el-card>

        <el-card class="timeline-card" shadow="never" style="margin-top: 16px">
          <div class="card-header">
            <h3 class="card-title">沟通记录</h3>
          </div>
          <div v-if="candidate && candidate.messages && candidate.messages.length > 0" class="message-simple-list">
            <div
              v-for="msg in candidate.messages"
              :key="msg.id"
              class="message-simple-item"
              :class="msg.senderRole === 'recruiter' ? 'right' : 'left'"
            >
              <div class="msg-header">
                <span class="msg-sender">{{ msg.senderName }}</span>
                <el-tag v-if="msg.senderRole === 'recruiter'" size="small" type="success">招聘方</el-tag>
                <el-tag v-else size="small" type="info">候选人</el-tag>
                <span class="msg-time">{{ formatTime(msg.sentAt) }}</span>
              </div>
              <div class="msg-content">{{ msg.content }}</div>
            </div>
          </div>
          <el-empty v-else description="暂无沟通记录" :image-size="100" />
        </el-card>

        <el-card class="notes-card" shadow="never" style="margin-top: 16px">
          <div class="card-header">
            <h3 class="card-title">备注记录</h3>
          </div>

          <el-form :model="noteForm" class="note-form">
            <el-row :gutter="12">
              <el-col :span="8">
                <el-form-item label="备注类型">
                  <el-select v-model="noteForm.type" style="width: 100%">
                    <el-option label="面试备注" value="interview" />
                    <el-option label="筛选意见" value="screen" />
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col :span="16">
                <el-form-item label="备注人">
                  <el-input v-model="noteForm.createdBy" placeholder="招聘专员-刘经理" />
                </el-form-item>
              </el-col>
            </el-row>
            <el-form-item>
              <el-input
                v-model="noteForm.content"
                type="textarea"
                :rows="3"
                placeholder="请输入备注内容，例如面试评价、筛选意见等"
              />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="submitNote" :loading="submittingNote" :icon="Plus">
                添加备注
              </el-button>
            </el-form-item>
          </el-form>

          <el-divider />

          <el-timeline v-if="candidate && candidate.notes && candidate.notes.length > 0">
            <el-timeline-item
              v-for="note in candidate.notes"
              :key="note.id"
              :timestamp="formatTime(note.createdAt)"
              placement="top"
            >
              <el-card shadow="never" class="note-item-card">
                <div class="note-header">
                  <el-tag :type="note.type === 'interview' ? 'primary' : 'warning'" size="small" effect="light">
                    {{ note.type === 'interview' ? '面试备注' : '筛选意见' }}
                  </el-tag>
                  <span class="note-creator">{{ note.createdBy }}</span>
                </div>
                <div class="note-content">{{ note.content }}</div>
              </el-card>
            </el-timeline-item>
          </el-timeline>
          <el-empty v-else description="暂无备注记录" :image-size="100" />
        </el-card>
      </el-col>
    </el-row>
    <el-dialog v-model="interviewVisible" title="安排面试" width="520px" destroy-on-close>
      <el-form :model="interviewForm" label-width="100px" :rules="interviewRules" ref="interviewFormRef">
        <el-form-item label="选择申请" prop="applicationId">
          <el-select v-model="interviewForm.applicationId" placeholder="请选择投递申请" filterable style="width: 100%">
            <el-option
              v-for="app in candidate?.applications || []"
              :key="app.id"
              :label="`${app.jobTitle}（${getStatusText(app.status)}）`"
              :value="app.id"
              :disabled="app.status === 'rejected' || app.status === 'hired'"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="面试轮次" prop="round">
          <el-input-number v-model="interviewForm.round" :min="1" :max="10" style="width: 100%" />
        </el-form-item>
        <el-form-item label="面试时间" prop="interviewTime">
          <el-date-picker
            v-model="interviewForm.interviewTime"
            type="datetime"
            placeholder="请选择面试时间"
            style="width: 100%"
            format="YYYY-MM-DD HH:mm"
            value-format="YYYY-MM-DDTHH:mm:ssZ"
          />
        </el-form-item>
        <el-form-item label="面试方式" prop="method">
          <el-radio-group v-model="interviewForm.method">
            <el-radio value="onsite">现场</el-radio>
            <el-radio value="online">视频</el-radio>
            <el-radio value="phone">电话</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="面试官" prop="interviewer">
          <el-input v-model="interviewForm.interviewer" placeholder="请输入面试官姓名" />
        </el-form-item>
        <el-form-item label="地点/链接">
          <el-input v-model="interviewForm.location" placeholder="现场地址或会议链接" />
        </el-form-item>
        <el-form-item label="面试说明">
          <el-input
            v-model="interviewForm.description"
            type="textarea"
            :rows="2"
            placeholder="可填写面试内容要点、注意事项等"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="interviewVisible = false">取消</el-button>
        <el-button type="primary" :loading="submittingInterview" @click="submitInterview">
          确认安排
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { ArrowLeft, Phone, Briefcase, ChatDotRound, Plus, Calendar } from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import { getCandidate, updateCandidateStatus, createNote, createInterview } from '@/api'

const route = useRoute()
const router = useRouter()

const loading = ref(false)
const submittingNote = ref(false)
const candidate = ref(null)

const noteForm = reactive({
  type: 'interview',
  content: '',
  createdBy: '招聘专员-刘经理'
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

const loadCandidate = async () => {
  const { contact, name } = route.query
  if (!contact || !name) {
    ElMessage.error('参数错误')
    return
  }

  loading.value = true
  try {
    const res = await getCandidate({ contact, name })
    const payload = res.data || {}
    candidate.value = {
      ...(payload.candidate || {}),
      messages: payload.messages || [],
      notes: payload.notes || []
    }
  } finally {
    loading.value = false
  }
}

const goBack = () => {
  router.back()
}

const changeStatus = async (status) => {
  if (!candidate.value) return
  try {
    await updateCandidateStatus({
      contact: candidate.value.contact,
      name: candidate.value.candidateName,
      status
    })
    ElMessage.success(`状态已变更为「${getStatusText(status)}」`)
    loadCandidate()
  } catch (e) {}
}

const goToMessage = () => {
  if (!candidate.value || candidate.value.applications.length === 0) {
    ElMessage.warning('该候选人暂无投递记录')
    return
  }
  const latestApp = candidate.value.applications[0]
  router.push({ path: '/messages', query: { applicationId: latestApp.id } })
}

const goToMessageById = (applicationId) => {
  router.push({ path: '/messages', query: { applicationId } })
}

const goToApplication = (applicationId) => {
  router.push({ path: '/applications', query: { applicationId } })
}

const submitNote = async () => {
  if (!candidate.value) return
  if (!noteForm.content.trim()) {
    ElMessage.warning('请输入备注内容')
    return
  }

  submittingNote.value = true
  try {
    await createNote({
      candidate: candidate.value.candidateName,
      contact: candidate.value.contact,
      type: noteForm.type,
      content: noteForm.content.trim(),
      createdBy: noteForm.createdBy.trim() || '招聘专员-刘经理'
    })
    ElMessage.success('备注添加成功')
    noteForm.content = ''
    loadCandidate()
  } finally {
    submittingNote.value = false
  }
}

const interviewVisible = ref(false)
const submittingInterview = ref(false)
const interviewFormRef = ref(null)
const interviewForm = reactive({
  applicationId: '',
  round: 1,
  interviewTime: '',
  method: 'online',
  interviewer: '',
  location: '',
  description: ''
})
const interviewRules = {
  applicationId: [{ required: true, message: '请选择申请', trigger: 'change' }],
  round: [{ required: true, message: '请选择面试轮次', trigger: 'blur' }],
  interviewTime: [{ required: true, message: '请选择面试时间', trigger: 'change' }],
  method: [{ required: true, message: '请选择面试方式', trigger: 'change' }],
  interviewer: [{ required: true, message: '请输入面试官', trigger: 'blur' }]
}

const openInterviewDialog = (app) => {
  if (!candidate.value || candidate.value.applications.length === 0) {
    ElMessage.warning('该候选人暂无有效投递申请')
    return
  }
  interviewForm.applicationId = app?.id || candidate.value.applications[0].id
  interviewForm.round = 1
  interviewForm.interviewTime = ''
  interviewForm.method = 'online'
  interviewForm.interviewer = '招聘专员-刘经理'
  interviewForm.location = ''
  interviewForm.description = ''
  interviewVisible.value = true
}

const submitInterview = async () => {
  await interviewFormRef.value?.validate()
  if (!interviewForm.applicationId || !interviewForm.interviewTime || !interviewForm.interviewer) return
  submittingInterview.value = true
  try {
    await createInterview({ ...interviewForm })
    ElMessage.success('面试已安排，候选人状态已同步为「面试中」')
    interviewVisible.value = false
    loadCandidate()
  } catch (e) {
    ElMessage.error('安排失败：' + (e?.response?.data?.error || e.message))
  } finally {
    submittingInterview.value = false
  }
}

onMounted(() => {
  loadCandidate()
})
</script>

<style scoped>
.candidate-detail-page {
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
}

.profile-card,
.timeline-card,
.notes-card {
  border-radius: 8px;
}

.profile-header {
  display: flex;
  align-items: center;
  gap: 16px;
}

.avatar {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: #fff;
  font-weight: 600;
  font-size: 28px;
}

.profile-info {
  flex: 1;
}

.profile-name {
  font-size: 22px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 8px;
  display: flex;
  align-items: center;
}

.profile-contact,
.profile-position {
  font-size: 14px;
  color: #606266;
  margin-top: 6px;
  display: flex;
  align-items: center;
  gap: 6px;
}

.section-title {
  font-size: 14px;
  font-weight: 600;
  color: #303133;
  margin: 0 0 12px;
}

.skill-tags {
  display: flex;
  flex-wrap: wrap;
}

.summary-text {
  font-size: 14px;
  line-height: 1.8;
  color: #606266;
  white-space: pre-wrap;
}

.status-buttons {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.quick-actions {
  display: flex;
  gap: 8px;
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

.timeline-item-actions {
  display: flex;
  gap: 12px;
}

.message-simple-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.message-simple-item {
  padding: 12px 16px;
  border-radius: 8px;
  background: #f8f9fa;
}

.message-simple-item.right {
  background: #ecf5ff;
}

.msg-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 6px;
}

.msg-sender {
  font-weight: 500;
  color: #303133;
}

.msg-time {
  font-size: 12px;
  color: #909399;
  margin-left: auto;
}

.msg-content {
  font-size: 14px;
  color: #606266;
  line-height: 1.6;
}

.note-form {
  margin-bottom: 16px;
}

.note-item-card {
  background: #fdf6ec;
  border: none;
}

.note-header {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 8px;
}

.note-creator {
  font-size: 13px;
  color: #909399;
}

.note-content {
  font-size: 14px;
  color: #606266;
  line-height: 1.7;
  white-space: pre-wrap;
}
</style>
