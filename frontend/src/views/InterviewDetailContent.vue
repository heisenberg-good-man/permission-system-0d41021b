<template>
  <div class="interview-detail">
    <div class="detail-header">
      <div class="candidate-block">
        <el-avatar :size="52" class="avatar">{{ interview.candidateName.charAt(0) }}</el-avatar>
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
      <el-descriptions-item v-if="interview.feedback" label="面试反馈" :span="2">
        <div class="feedback-box">{{ interview.feedback }}</div>
      </el-descriptions-item>
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
        <el-button type="primary" :icon="Check" @click="openComplete">标记已完成</el-button>
      </div>
    </div>

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

    <el-dialog v-model="completeVisible" title="标记面试完成" width="480px" destroy-on-close>
      <el-form :model="completeForm" label-width="90px">
        <el-form-item label="面试反馈">
          <el-input
            v-model="completeForm.feedback"
            type="textarea"
            :rows="4"
            placeholder="请输入面试整体评价和反馈意见"
          />
        </el-form-item>
        <el-form-item label="备注">
          <el-input
            v-model="completeForm.note"
            type="textarea"
            :rows="2"
            placeholder="可填写后续安排"
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
import { ref, reactive } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Phone, Briefcase, Calendar, Clock, CloseBold, Check, ChatDotRound } from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import { rescheduleInterview, cancelInterview, completeInterview, addInterviewNote } from '@/api'

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

const completeVisible = ref(false)
const completeForm = reactive({ feedback: '', note: '' })
const openComplete = () => {
  completeForm.feedback = ''
  completeForm.note = ''
  completeVisible.value = true
}
const confirmComplete = async () => {
  submitting.value = true
  try {
    await completeInterview(props.interview.id, completeForm)
    ElMessage.success('面试已标记完成')
    completeVisible.value = false
    emit('updated')
  } catch (e) {
    ElMessage.error('操作失败：' + (e?.response?.data?.error || e.message))
  } finally {
    submitting.value = false
  }
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

.feedback-box {
  padding: 10px 12px;
  background: #f0f9ff;
  border-radius: 6px;
  border-left: 3px solid #409eff;
  line-height: 1.6;
  color: #303133;
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
