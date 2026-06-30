<template>
  <div class="interview-detail">
    <div class="detail-header">
      <div class="candidate-block">
        <el-avatar :size="52" class="avatar">{{ interview.candidateName?.charAt(0) || '-' }}</el-avatar>
        <div class="candidate-info">
          <div class="candidate-name-row">
            <span class="name">{{ interview.candidateName }}</span>
            <el-tag :type="getStatusTagType(interview.status)" size="small" effect="light" style="margin-left: 8px">
              {{ getStatusText(interview.status) }}
            </el-tag>
            <el-tag :type="getMethodTagType(interview.method)" size="small" effect="plain" style="margin-left: 6px">
              {{ getMethodText(interview.method) }}
            </el-tag>
          </div>
          <div class="candidate-contact">
            <el-icon><Phone /></el-icon>
            {{ interview.contact }}
          </div>
          <div class="candidate-job">
            <el-icon><Briefcase /></el-icon>
            {{ interview.jobTitle }} · 第{{ interview.round }}轮
          </div>
        </div>
      </div>
    </div>

    <el-descriptions :column="2" border size="default" class="info-desc">
      <el-descriptions-item label="面试时间">
        <span class="time-value">
          <el-icon><Calendar /></el-icon>
          {{ formatTime(interview.interviewTime) }}
        </span>
      </el-descriptions-item>
      <el-descriptions-item label="面试官">{{ interview.interviewer || '-' }}</el-descriptions-item>
      <el-descriptions-item label="面试地点/方式">{{ interview.location || '-' }}</el-descriptions-item>
      <el-descriptions-item label="创建时间">{{ formatTime(interview.createdAt) }}</el-descriptions-item>
      <el-descriptions-item label="面试说明" :span="2">{{ interview.description || '暂无说明' }}</el-descriptions-item>
      <el-descriptions-item v-if="interview.latestNote" label="最近备注" :span="2">
        <div class="note-box">
          <el-icon style="color: #e6a23c"><ChatDotRound /></el-icon>
          {{ interview.latestNote }}
        </div>
      </el-descriptions-item>
    </el-descriptions>

    <div class="action-section" v-if="interview.status === 'scheduled'">
      <el-divider content-position="left">快捷操作</el-divider>
      <div class="action-buttons">
        <el-button type="success" :icon="Clock" @click="openReschedule">改期</el-button>
        <el-button type="warning" :icon="CloseBold" @click="handleCancel">取消面试</el-button>
        <el-button type="primary" :icon="Check" @click="openFeedback">填写面试反馈</el-button>
      </div>
    </div>

    <div v-if="interview.status === 'completed' && hasFeedback" class="feedback-section">
      <el-divider content-position="left">面试反馈</el-divider>
      <div class="feedback-card">
        <div class="fb-row">
          <span class="fb-label">反馈结论</span>
          <el-tag :type="getConclusionTagType(interview.conclusion)" size="default" effect="dark">
            {{ getConclusionText(interview.conclusion) }}
          </el-tag>
        </div>
        <div class="fb-row" v-if="interview.rating">
          <span class="fb-label">评分</span>
          <el-rate v-model="interview.rating" disabled />
        </div>
        <div class="fb-row" v-if="interview.feedback">
          <span class="fb-label">整体评价</span>
          <div class="fb-content-box">{{ interview.feedback }}</div>
        </div>
        <div class="fb-row" v-if="interview.strengths">
          <span class="fb-label">优势</span>
          <div class="fb-content-box strengths">{{ interview.strengths }}</div>
        </div>
        <div class="fb-row" v-if="interview.risks">
          <span class="fb-label">风险点</span>
          <div class="fb-content-box risks">{{ interview.risks }}</div>
        </div>
        <div class="fb-row" v-if="interview.nextSteps">
          <span class="fb-label">下一步建议</span>
          <div class="fb-content-box next-steps">{{ interview.nextSteps }}</div>
        </div>
      </div>
      <div class="action-section" v-if="canCreateOffer" style="margin-top: 12px">
        <el-button type="primary" :icon="Promotion" @click="openOfferDialog">发起 Offer</el-button>
      </div>
    </div>

    <div v-if="interview.status === 'completed' && !hasFeedback" class="feedback-section">
      <el-divider content-position="left">面试反馈</el-divider>
      <el-empty description="该面试尚未填写结构化反馈" :image-size="60">
        <el-button type="primary" @click="openFeedback">补充反馈</el-button>
      </el-empty>
      <div class="action-section" style="margin-top: 12px">
        <el-button type="primary" :icon="Promotion" @click="openOfferDialog">发起 Offer</el-button>
      </div>
    </div>

    <OfferDialog
      v-model="offerVisible"
      :applications="offerApplications"
      :candidate-name="offerCandidateName"
      :pre-select-application-id="offerPreSelectId"
      @success="handleOfferSuccess"
    />

    <el-divider content-position="left">补充备注</el-divider>
    <div class="note-section">
      <el-input
        v-model="noteForm.note"
        type="textarea"
        :rows="3"
        placeholder="可补充面试备注、筛选意见等，保存后立即显示"
      />
      <div class="note-submit">
        <el-button type="primary" :loading="submitting" @click="submitNote">保存备注</el-button>
      </div>
    </div>

    <el-dialog v-model="rescheduleVisible" title="面试改期" width="440px" destroy-on-close>
      <el-form :model="rescheduleForm" label-width="100px" :rules="rescheduleRules" ref="rescheduleFormRef">
        <el-form-item label="当前时间">
          <span>{{ formatTime(interview.interviewTime) }}</span>
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

    <el-dialog v-model="feedbackVisible" title="面试反馈" width="560px" destroy-on-close>
      <el-form :model="feedbackForm" label-width="100px" :rules="feedbackRules" ref="feedbackFormRef">
        <el-form-item label="反馈结论" prop="conclusion">
          <el-radio-group v-model="feedbackForm.conclusion">
            <el-radio value="pass">通过</el-radio>
            <el-radio value="pending">待定</el-radio>
            <el-radio value="fail">不通过</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="评分">
          <el-rate v-model="feedbackForm.rating" show-text :texts="['很差','差','一般','良好','优秀']" />
        </el-form-item>
        <el-form-item label="整体评价">
          <el-input
            v-model="feedbackForm.feedback"
            type="textarea"
            :rows="3"
            placeholder="请填写面试整体评价和反馈意见"
          />
        </el-form-item>
        <el-form-item label="优势">
          <el-input
            v-model="feedbackForm.strengths"
            type="textarea"
            :rows="2"
            placeholder="候选人的突出优势、亮点等"
          />
        </el-form-item>
        <el-form-item label="风险点">
          <el-input
            v-model="feedbackForm.risks"
            type="textarea"
            :rows="2"
            placeholder="候选人的不足、潜在风险等"
          />
        </el-form-item>
        <el-form-item label="下一步建议">
          <el-input
            v-model="feedbackForm.nextSteps"
            type="textarea"
            :rows="2"
            placeholder="建议后续安排，如推进二面、发Offer等"
          />
        </el-form-item>
        <el-form-item label="备注">
          <el-input
            v-model="feedbackForm.note"
            type="textarea"
            :rows="2"
            placeholder="可填写其他备注"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="feedbackVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="confirmFeedback">提交反馈</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Phone, Briefcase, Calendar, Clock, CloseBold, Check, ChatDotRound, Promotion } from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import { rescheduleInterview, cancelInterview, completeInterview, addInterviewNote, getApplication } from '@/api'
import OfferDialog from './OfferDialog.vue'

const props = defineProps({
  interview: { type: Object, required: true }
})
const emit = defineEmits(['updated'])

const submitting = ref(false)

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
const getConclusionText = (c) => {
  const map = { pass: '通过', pending: '待定', fail: '不通过' }
  return map[c] || c
}
const getConclusionTagType = (c) => {
  const map = { pass: 'success', pending: 'warning', fail: 'danger' }
  return map[c] || 'info'
}

const hasFeedback = computed(() => {
  return props.interview.conclusion || props.interview.rating || props.interview.feedback
})

const canCreateOffer = computed(() => {
  return props.interview.status === 'completed' &&
    (props.interview.conclusion === 'pass' || props.interview.conclusion === 'pending')
})

const noteForm = reactive({ note: '' })
const submitNote = async () => {
  if (!noteForm.note.trim()) {
    ElMessage.warning('请输入备注内容')
    return
  }
  submitting.value = true
  try {
    await addInterviewNote(props.interview.id, { note: noteForm.note.trim() })
    ElMessage.success('备注已保存')
    noteForm.note = ''
    emit('updated')
  } catch (e) {
    ElMessage.error('保存失败：' + (e?.response?.data?.error || e.message))
  } finally {
    submitting.value = false
  }
}

const rescheduleVisible = ref(false)
const rescheduleForm = reactive({ interviewTime: '' })
const rescheduleFormRef = ref(null)
const rescheduleRules = {
  interviewTime: [{ required: true, message: '请选择新的面试时间', trigger: 'change' }]
}
const openReschedule = () => {
  rescheduleForm.interviewTime = ''
  rescheduleVisible.value = true
}
const confirmReschedule = async () => {
  await rescheduleFormRef.value?.validate()
  if (!rescheduleForm.interviewTime) return
  submitting.value = true
  try {
    await rescheduleInterview(props.interview.id, { interviewTime: rescheduleForm.interviewTime })
    ElMessage.success('面试改期成功')
    rescheduleVisible.value = false
    emit('updated')
  } catch (e) {
    ElMessage.error('改期失败：' + (e?.response?.data?.error || e.message))
  } finally {
    submitting.value = false
  }
}

const handleCancel = () => {
  ElMessageBox.confirm('确认取消该面试安排？取消后申请状态将回退为待沟通。', '取消面试', {
    confirmButtonText: '确认取消',
    cancelButtonText: '再想想',
    type: 'warning'
  }).then(async () => {
    try {
      await cancelInterview(props.interview.id)
      ElMessage.success('已取消面试安排')
      emit('updated')
    } catch (e) {
      ElMessage.error('取消失败：' + (e?.response?.data?.error || e.message))
    }
  }).catch(() => {})
}

const feedbackVisible = ref(false)
const feedbackFormRef = ref(null)
const feedbackForm = reactive({
  conclusion: 'pass',
  rating: 4,
  feedback: '',
  strengths: '',
  risks: '',
  nextSteps: '',
  note: ''
})
const feedbackRules = {
  conclusion: [{ required: true, message: '请选择反馈结论', trigger: 'change' }]
}
const openFeedback = () => {
  feedbackForm.conclusion = props.interview.conclusion || 'pass'
  feedbackForm.rating = props.interview.rating || 4
  feedbackForm.feedback = props.interview.feedback || ''
  feedbackForm.strengths = props.interview.strengths || ''
  feedbackForm.risks = props.interview.risks || ''
  feedbackForm.nextSteps = props.interview.nextSteps || ''
  feedbackForm.note = ''
  feedbackVisible.value = true
}
const confirmFeedback = async () => {
  await feedbackFormRef.value?.validate()
  submitting.value = true
  try {
    await completeInterview(props.interview.id, {
      conclusion: feedbackForm.conclusion,
      rating: feedbackForm.rating,
      feedback: feedbackForm.feedback,
      strengths: feedbackForm.strengths,
      risks: feedbackForm.risks,
      nextSteps: feedbackForm.nextSteps,
      note: feedbackForm.note
    })
    ElMessage.success('面试反馈已提交')
    feedbackVisible.value = false
    emit('updated')
  } catch (e) {
    ElMessage.error('提交失败：' + (e?.response?.data?.error || e.message))
  } finally {
    submitting.value = false
  }
}

const offerVisible = ref(false)
const offerApplications = ref([])
const offerPreSelectId = ref('')
const offerCandidateName = ref('')
const loadOfferApplications = async () => {
  if (!props.interview.applicationId) return
  try {
    const res = await getApplication(props.interview.applicationId)
    if (res?.data) {
      offerApplications.value = [res.data]
      offerPreSelectId.value = res.data.id
      offerCandidateName.value = res.data.resume?.candidateName || props.interview.candidateName
    }
  } catch (e) {
    console.error('加载申请失败:', e)
  }
}
const openOfferDialog = async () => {
  offerCandidateName.value = props.interview.candidateName
  offerPreSelectId.value = props.interview.applicationId || ''
  offerApplications.value = []
  if (props.interview.applicationId) {
    await loadOfferApplications()
  }
  offerVisible.value = true
}
const handleOfferSuccess = () => {
  ElMessage.success('Offer 已创建')
  emit('updated')
}
</script>

<style scoped>
.interview-detail {
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
  background: linear-gradient(135deg, #409eff, #66b1ff);
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

.time-value {
  display: inline-flex;
  align-items: center;
  gap: 4px;
}

.note-box {
  padding: 8px 12px;
  background: #fdf6ec;
  border-radius: 6px;
  display: flex;
  align-items: flex-start;
  gap: 6px;
  color: #606266;
  line-height: 1.5;
}

.feedback-section {
  margin-top: 4px;
}

.feedback-card {
  background: #fafafa;
  border-radius: 10px;
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.fb-row {
  display: flex;
  align-items: flex-start;
  gap: 12px;
}

.fb-label {
  font-size: 14px;
  font-weight: 600;
  color: #606266;
  min-width: 80px;
  flex-shrink: 0;
  padding-top: 4px;
}

.fb-content-box {
  flex: 1;
  padding: 10px 14px;
  border-radius: 6px;
  line-height: 1.7;
  color: #303133;
  font-size: 14px;
}

.fb-content-box.strengths {
  background: #f0f9eb;
  border-left: 3px solid #67c23a;
}

.fb-content-box.risks {
  background: #fef0f0;
  border-left: 3px solid #f56c6c;
}

.fb-content-box.next-steps {
  background: #fdf6ec;
  border-left: 3px solid #e6a23c;
}

.action-section {
  margin-top: 4px;
}

.action-buttons {
  display: flex;
  gap: 10px;
}

.note-section {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.note-submit {
  display: flex;
  justify-content: flex-end;
}
</style>
