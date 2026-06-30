<template>
  <div class="messages-page">
    <el-row :gutter="16">
      <el-col :span="9">
        <el-card class="conversation-card" shadow="never">
          <div class="card-header">
            <h3 class="card-title">申请沟通</h3>
            <el-input
              v-model="searchKeyword"
              placeholder="搜索候选人/职位"
              size="small"
              clearable
              style="width: 200px"
              :prefix-icon="Search"
            />
          </div>
          <el-scrollbar class="conversation-scroll">
            <div
              v-for="app in filteredApplications"
              :key="app.id"
              class="conversation-item"
              :class="{ active: currentApplication?.id === app.id }"
              @click="selectApplication(app)"
            >
              <el-avatar :size="42" class="avatar">{{ app.resume.candidateName.charAt(0) }}</el-avatar>
              <div class="conversation-body">
                <div class="conversation-top">
                  <span class="candidate-name">{{ app.resume.candidateName }}</span>
                  <el-tag :type="getStatusTagType(app.status)" size="small" effect="light">
                    {{ getStatusText(app.status) }}
                  </el-tag>
                </div>
                <div class="conversation-job">{{ app.jobTitle }}</div>
                <div class="conversation-meta">
                  <span class="last-msg">{{ getLastMessageTip(app) }}</span>
                  <span class="last-time">{{ formatLastTime(app) }}</span>
                </div>
              </div>
              <el-badge v-if="app.unreadCount > 0" :value="app.unreadCount" class="unread-badge" />
            </div>
            <el-empty v-if="filteredApplications.length === 0" description="暂无沟通记录" :image-size="100" />
          </el-scrollbar>
        </el-card>
      </el-col>

      <el-col :span="15">
        <el-card v-if="currentApplication" class="chat-card" shadow="never">
          <div class="chat-header">
            <div class="chat-header-left">
              <el-avatar :size="42" class="avatar">{{ currentApplication.resume.candidateName.charAt(0) }}</el-avatar>
              <div>
                <div class="chat-title">
                  <span>{{ currentApplication.resume.candidateName }}</span>
                  <el-tag :type="getStatusTagType(currentApplication.status)" size="small" effect="light" style="margin-left: 8px">
                    {{ getStatusText(currentApplication.status) }}
                  </el-tag>
                </div>
                <div class="chat-subtitle">
                  应聘：{{ currentApplication.jobTitle }} · 联系电话：{{ currentApplication.resume.contact }}
                </div>
              </div>
            </div>
            <div class="chat-header-right">
              <el-select
                v-model="currentApplication.status"
                size="small"
                style="width: 130px"
                @change="onStatusChange"
              >
                <el-option
                  v-for="item in statusOptions"
                  :key="item.value"
                  :label="item.label"
                  :value="item.value"
                />
              </el-select>
            </div>
          </div>

          <el-divider style="margin: 10px 0" />

          <el-scrollbar ref="messageScrollRef" class="message-scroll" @scroll="handleScroll">
            <div class="message-list">
              <div v-if="messageList.length === 0" class="empty-tip">
                <el-empty description="还没有消息，发送第一条消息开始沟通吧" :image-size="80" />
              </div>
              <div
                v-for="msg in messageList"
                :key="msg.id"
                class="message-item"
                :class="msg.senderRole === 'recruiter' ? 'right' : 'left'"
              >
                <div v-if="msg.senderRole === 'candidate'" class="msg-avatar-wrap">
                  <el-avatar :size="36" class="msg-avatar candidate">{{ msg.senderName.charAt(0) }}</el-avatar>
                </div>
                <div class="msg-content-wrap">
                  <div class="msg-meta">
                    <span class="msg-sender">{{ msg.senderName }}</span>
                    <el-tag v-if="msg.senderRole === 'recruiter'" size="small" type="success">招聘方</el-tag>
                    <el-tag v-else size="small" type="info">候选人</el-tag>
                    <span class="msg-time">{{ formatMsgTime(msg.sentAt) }}</span>
                  </div>
                  <div class="msg-bubble" :class="msg.senderRole === 'recruiter' ? 'recruiter' : 'candidate'">
                    {{ msg.content }}
                  </div>
                </div>
                <div v-if="msg.senderRole === 'recruiter'" class="msg-avatar-wrap">
                  <el-avatar :size="36" class="msg-avatar recruiter">招</el-avatar>
                </div>
              </div>
            </div>
          </el-scrollbar>

          <el-divider style="margin: 10px 0" />

          <div class="chat-input-area">
            <div class="sender-switch">
              <span class="sender-label">当前身份（模拟）：</span>
              <el-radio-group v-model="currentSenderRole" size="small">
                <el-radio-button value="recruiter">招聘方</el-radio-button>
                <el-radio-button value="candidate">候选人</el-radio-button>
              </el-radio-group>
              <el-input
                v-if="currentSenderRole === 'candidate'"
                v-model="candidateSenderName"
                size="small"
                placeholder="输入候选人姓名"
                style="width: 140px; margin-left: 8px"
              />
            </div>
            <el-input
              v-model="messageContent"
              type="textarea"
              :rows="3"
              placeholder="请输入消息内容，Enter发送，Shift+Enter换行"
              resize="none"
              @keydown.enter.prevent.exact="handleSendMessage"
              class="message-input"
            />
            <div class="input-footer">
              <div class="quick-replies">
                <el-tag
                  v-for="tpl in quickReplies"
                  :key="tpl"
                  size="small"
                  class="quick-tag"
                  type="info"
                  effect="plain"
                  @click="messageContent = tpl"
                >
                  {{ tpl }}
                </el-tag>
              </div>
              <el-button type="primary" @click="handleSendMessage" :disabled="!messageContent.trim()" :icon="Promotion">
                发送
              </el-button>
            </div>
          </div>
        </el-card>

        <el-card v-else class="chat-card empty-chat" shadow="never">
          <el-empty description="请从左侧选择一个申请开始沟通" :image-size="160">
            <template #image>
              <el-icon class="empty-chat-icon"><ChatLineSquare /></el-icon>
            </template>
          </el-empty>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, watch, nextTick } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Search, Promotion, ChatLineSquare } from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import { listApplications, listMessages, sendMessage as sendMessageApi, updateApplicationStatus } from '@/api'

const route = useRoute()
const loading = ref(false)
const applications = ref([])
const searchKeyword = ref('')
const currentApplication = ref(null)
const messageList = ref([])
const messageScrollRef = ref(null)

const messageContent = ref('')
const currentSenderRole = ref('recruiter')
const candidateSenderName = ref('候选人')

const statusOptions = [
  { value: 'submitted', label: '已投递' },
  { value: 'pending', label: '待沟通' },
  { value: 'interview', label: '面试中' },
  { value: 'rejected', label: '已拒绝' },
  { value: 'hired', label: '已录用' }
]

const quickReplies = [
  '您好，方便约个时间沟通一下吗？',
  '感谢投递，我们正在评估您的简历',
  '请问您方便什么时间参加面试？',
  '期待收到您的回复，谢谢！'
]

const getStatusText = (status) => {
  const opt = statusOptions.find(s => s.value === status)
  return opt ? opt.label : status
}

const getStatusTagType = (status) => {
  const map = { submitted: 'info', pending: 'warning', interview: 'primary', rejected: 'danger', hired: 'success' }
  return map[status] || 'info'
}

const filteredApplications = computed(() => {
  const kw = searchKeyword.value.trim().toLowerCase()
  if (!kw) return applications.value
  return applications.value.filter(a =>
    a.resume.candidateName.toLowerCase().includes(kw) ||
    a.jobTitle.toLowerCase().includes(kw)
  )
})

const loadApplications = async () => {
  loading.value = true
  try {
    const res = await listApplications({})
    applications.value = (res.data || []).sort((a, b) => {
      const ta = a.lastMessageTime || a.submittedAt
      const tb = b.lastMessageTime || b.submittedAt
      return new Date(tb) - new Date(ta)
    })
  } finally {
    loading.value = false
  }
}

const selectApplication = async (app) => {
  currentApplication.value = app
  candidateSenderName.value = app.resume.candidateName
  messageList.value = []
  const res = await listMessages(app.id)
  messageList.value = res.data || []
  await nextTick()
  scrollToBottom()
}

const getLastMessageTip = (app) => {
  if (app.unreadCount > 0) return `${app.unreadCount}条新消息`
  if (app.lastMessageTime) return '已有沟通记录'
  return '暂无沟通'
}

const formatLastTime = (app) => {
  const t = app.lastMessageTime || app.submittedAt
  const d = dayjs(t)
  const now = dayjs()
  if (d.isSame(now, 'day')) return d.format('HH:mm')
  if (d.isSame(now.subtract(1, 'day'), 'day')) return '昨天'
  if (now.diff(d, 'day') < 7) return `${now.diff(d, 'day')}天前`
  return d.format('MM-DD')
}

const formatMsgTime = (t) => dayjs(t).format('YYYY-MM-DD HH:mm:ss')

const scrollToBottom = () => {
  if (!messageScrollRef.value) return
  const wrap = messageScrollRef.value.$el.querySelector('.el-scrollbar__wrap')
  if (wrap) wrap.scrollTop = wrap.scrollHeight
}

const handleScroll = () => {}

const handleSendMessage = async () => {
  if (!currentApplication.value || !messageContent.value.trim()) return
  const senderName = currentSenderRole.value === 'recruiter'
    ? '招聘专员-刘经理'
    : candidateSenderName.value || '候选人'
  try {
    const res = await sendMessageApi({
      applicationId: currentApplication.value.id,
      senderName,
      senderRole: currentSenderRole.value,
      content: messageContent.value.trim()
    })
    messageList.value.push(res.data)
    messageContent.value = ''
    await nextTick()
    scrollToBottom()
    if (currentApplication.value) currentApplication.value.lastMessageTime = res.data.sentAt
  } catch (e) {}
}

const onStatusChange = async (newStatus) => {
  if (!currentApplication.value) return
  try {
    await updateApplicationStatus(currentApplication.value.id, newStatus)
    ElMessage.success(`状态已变更为「${getStatusText(newStatus)}」`)
    loadApplications()
  } catch (e) {
    // 回滚
    loadApplications()
  }
}

watch(
  () => route.query.applicationId,
  async (val) => {
    if (!val) return
    await loadApplications()
    const app = applications.value.find(a => a.id === val)
    if (app) selectApplication(app)
  },
  { immediate: true }
)

onMounted(loadApplications)
</script>

<style scoped>
.messages-page {
  height: 100%;
}

.conversation-card,
.chat-card {
  border-radius: 8px;
  height: calc(100vh - 140px);
  display: flex;
  flex-direction: column;
}

.chat-card :deep(.el-card__body) {
  display: flex;
  flex-direction: column;
  height: 100%;
  padding: 16px 20px;
  overflow: hidden;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  flex-shrink: 0;
}

.card-title {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
  margin: 0;
}

.conversation-scroll {
  flex: 1;
  margin: 0 -8px;
  padding: 0 8px;
}

.conversation-item {
  display: flex;
  align-items: center;
  padding: 12px 10px;
  border-radius: 8px;
  cursor: pointer;
  margin-bottom: 4px;
  position: relative;
  transition: background 0.2s;
}

.conversation-item:hover {
  background: #f5f7fa;
}

.conversation-item.active {
  background: #ecf5ff;
}

.avatar {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: #fff;
  font-weight: 600;
  flex-shrink: 0;
}

.conversation-body {
  flex: 1;
  margin-left: 10px;
  min-width: 0;
}

.conversation-top {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 6px;
}

.candidate-name {
  font-weight: 600;
  color: #303133;
  font-size: 14px;
}

.conversation-job {
  font-size: 12px;
  color: #606266;
  margin-top: 3px;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
}

.conversation-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 4px;
}

.last-msg {
  font-size: 12px;
  color: #909399;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
  flex: 1;
  margin-right: 8px;
}

.last-time {
  font-size: 11px;
  color: #c0c4cc;
  flex-shrink: 0;
}

.unread-badge {
  position: absolute;
  top: 12px;
  right: 12px;
}

.chat-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-shrink: 0;
}

.chat-header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.chat-title {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.chat-subtitle {
  font-size: 13px;
  color: #909399;
  margin-top: 4px;
}

.message-scroll {
  flex: 1;
  margin: 0 -8px;
  padding: 8px;
}

.message-list {
  padding: 8px 0;
}

.empty-tip {
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.message-item {
  display: flex;
  margin-bottom: 18px;
  gap: 10px;
}

.message-item.right {
  flex-direction: row-reverse;
}

.message-item.right .msg-content-wrap {
  align-items: flex-end;
}

.message-item.right .msg-meta {
  flex-direction: row-reverse;
}

.message-item.right .msg-meta .msg-time {
  margin: 0 6px 0 0;
}

.msg-avatar-wrap {
  flex-shrink: 0;
}

.msg-avatar {
  font-weight: 600;
  color: #fff;
}

.msg-avatar.recruiter {
  background: linear-gradient(135deg, #11998e 0%, #38ef7d 100%);
}

.msg-avatar.candidate {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.msg-content-wrap {
  display: flex;
  flex-direction: column;
  max-width: 65%;
}

.msg-meta {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-bottom: 4px;
  font-size: 12px;
}

.msg-sender {
  color: #606266;
  font-weight: 500;
}

.msg-time {
  color: #c0c4cc;
  margin-left: 6px;
}

.msg-bubble {
  padding: 10px 14px;
  border-radius: 10px;
  line-height: 1.6;
  word-wrap: break-word;
  white-space: pre-wrap;
  font-size: 14px;
}

.msg-bubble.recruiter {
  background: #ecf5ff;
  color: #303133;
  border-top-right-radius: 4px;
}

.msg-bubble.candidate {
  background: #fff;
  border: 1px solid #ebeef5;
  color: #303133;
  border-top-left-radius: 4px;
}

.chat-input-area {
  flex-shrink: 0;
}

.sender-switch {
  margin-bottom: 10px;
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 4px;
}

.sender-label {
  font-size: 13px;
  color: #606266;
}

.message-input :deep(.el-textarea__inner) {
  border-radius: 8px;
}

.input-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 10px;
  gap: 12px;
  flex-wrap: wrap;
}

.quick-replies {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  flex: 1;
}

.quick-tag {
  cursor: pointer;
}

.empty-chat {
  display: flex;
  align-items: center;
  justify-content: center;
}

.empty-chat-icon {
  font-size: 140px;
  color: #dcdfe6;
}
</style>
