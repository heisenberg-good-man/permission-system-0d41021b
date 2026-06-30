<template>
  <div class="offer-detail">
    <div class="detail-header">
      <div class="candidate-block">
        <el-avatar :size="52" class="avatar">{{ offer.candidateName?.charAt(0) || '-' }}</el-avatar>
        <div class="candidate-info">
          <div class="candidate-name-row">
            <span class="name">{{ offer.candidateName }}</span>
            <el-tag :type="getStatusTagType(offer.status)" size="small" effect="light" style="margin-left: 8px">
              {{ getStatusText(offer.status) }}
            </el-tag>
          </div>
          <div class="candidate-contact">
            <el-icon><Phone /></el-icon>
            {{ offer.contact }}
          </div>
          <div class="candidate-job">
            <el-icon><Briefcase /></el-icon>
            {{ offer.jobTitle }}
          </div>
        </div>
      </div>
    </div>

    <el-descriptions :column="2" border size="default" class="info-desc">
      <el-descriptions-item label="薪资范围">
        <span class="salary-text">¥{{ formatSalary(offer.salaryMin) }}K - ¥{{ formatSalary(offer.salaryMax) }}K</span>
      </el-descriptions-item>
      <el-descriptions-item label="入职日期">{{ offer.startDate || '-' }}</el-descriptions-item>
      <el-descriptions-item label="负责人">{{ offer.owner || '-' }}</el-descriptions-item>
      <el-descriptions-item label="创建时间">{{ formatTime(offer.createdAt) }}</el-descriptions-item>
      <el-descriptions-item v-if="offer.sentAt" label="发送时间">{{ formatTime(offer.sentAt) }}</el-descriptions-item>
      <el-descriptions-item v-if="offer.repliedAt" label="回复时间">{{ formatTime(offer.repliedAt) }}</el-descriptions-item>
      <el-descriptions-item label="补充说明" :span="2">{{ offer.description || '暂无说明' }}</el-descriptions-item>
      <el-descriptions-item v-if="offer.feedback" label="反馈/原因" :span="2">
        <div class="feedback-box">{{ offer.feedback }}</div>
      </el-descriptions-item>
    </el-descriptions>

    <div class="action-section" v-if="offer.status === 'pending' || offer.status === 'sent'">
      <el-divider content-position="left">快捷操作</el-divider>
      <div class="action-buttons">
        <el-button v-if="offer.status === 'pending'" type="primary" :icon="Promotion" @click="handleSend">发送 Offer</el-button>
        <el-button v-if="offer.status === 'sent'" type="success" :icon="Check" @click="handleReply(true)">候选人接受</el-button>
        <el-button v-if="offer.status === 'sent'" type="warning" :icon="CloseBold" @click="handleReply(false)">候选人拒绝</el-button>
        <el-button v-if="offer.status !== 'accepted' && offer.status !== 'rejected' && offer.status !== 'withdrawn'" type="danger" :icon="Close" @click="handleWithdraw">撤回 Offer</el-button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Phone, Briefcase, Promotion, Check, Close, CloseBold } from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import { sendOffer, replyOffer, withdrawOffer } from '@/api'

const props = defineProps({
  offer: { type: Object, required: true }
})
const emit = defineEmits(['updated'])

const submitting = ref(false)

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

const handleSend = async () => {
  submitting.value = true
  try {
    await sendOffer(props.offer.id)
    ElMessage.success('Offer 已发送')
    emit('updated')
  } catch (e) {
    ElMessage.error('发送失败：' + (e?.response?.data?.error || e.message))
  } finally {
    submitting.value = false
  }
}

const handleReply = async (accepted) => {
  const action = accepted ? '接受' : '拒绝'
  try {
    await ElMessageBox.confirm(
      `确认候选人已${action}该 Offer？${accepted ? '接受后申请状态将变更为已录用。' : ''}`,
      `${action} Offer`,
      { type: 'warning' }
    )
    submitting.value = true
    await replyOffer(props.offer.id, { accepted, feedback: `候选人已${action}Offer` })
    ElMessage.success(`候选人已${action}`)
    emit('updated')
  } catch (e) {
    if (e !== 'cancel') {
      ElMessage.error('操作失败：' + (e?.response?.data?.error || e.message))
    }
  } finally {
    submitting.value = false
  }
}

const handleWithdraw = async () => {
  try {
    const { value } = await ElMessageBox.prompt('请输入撤回原因', '撤回 Offer', {
      inputPlaceholder: '请说明撤回原因',
      confirmButtonText: '确认撤回',
      cancelButtonText: '取消',
      type: 'warning'
    })
    submitting.value = true
    await withdrawOffer(props.offer.id, { reason: value || '公司撤回' })
    ElMessage.success('Offer 已撤回')
    emit('updated')
  } catch (e) {
    if (e !== 'cancel') {
      ElMessage.error('撤回失败：' + (e?.response?.data?.error || e.message))
    }
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.offer-detail {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.detail-header {
  padding: 4px 0;
}

.candidate-block {
  display: flex;
  align-items: center;
  gap: 16px;
}

.avatar {
  background: linear-gradient(135deg, #67c23a, #85ce61);
  font-size: 20px;
}

.candidate-info {
  flex: 1;
}

.candidate-name-row {
  display: flex;
  align-items: center;
  font-size: 18px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 6px;
}

.candidate-contact,
.candidate-job {
  font-size: 13px;
  color: #606266;
  display: flex;
  align-items: center;
  gap: 4px;
  margin-top: 3px;
}

.info-desc {
  margin-top: 8px;
}

.salary-text {
  font-weight: 600;
  color: #f56c6c;
}

.feedback-box {
  padding: 10px 12px;
  background: #fef0f0;
  border-radius: 6px;
  border-left: 3px solid #f56c6c;
  line-height: 1.6;
  color: #303133;
}

.action-section {
  margin-top: 4px;
}

.action-buttons {
  display: flex;
  gap: 10px;
}
</style>
