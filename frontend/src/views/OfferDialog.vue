<template>
  <el-dialog v-model="visibleLocal" title="发起 Offer" width="560px" destroy-on-close @close="handleClose">
    <el-form :model="form" label-width="100px" :rules="rules" ref="formRef">
      <el-form-item label="候选人">
        <span>{{ candidateName }}</span>
      </el-form-item>
      <el-form-item label="投递职位" prop="applicationId">
        <el-select v-model="form.applicationId" placeholder="请选择投递申请" filterable style="width: 100%">
          <el-option
            v-for="app in applications"
            :key="app.id"
            :label="app.jobTitle"
            :value="app.id"
          />
        </el-select>
      </el-form-item>
      <el-form-item label="薪资范围">
        <div style="display: flex; align-items: center; gap: 8px; width: 100%">
          <el-input-number v-model="form.salaryMin" :min="0" :step="1000" style="flex: 1" placeholder="最低薪资" />
          <span>至</span>
          <el-input-number v-model="form.salaryMax" :min="0" :step="1000" style="flex: 1" placeholder="最高薪资" />
        </div>
      </el-form-item>
      <el-form-item label="入职日期" prop="startDate">
        <el-date-picker
          v-model="form.startDate"
          type="date"
          placeholder="请选择入职日期"
          style="width: 100%"
          format="YYYY-MM-DD"
          value-format="YYYY-MM-DD"
        />
      </el-form-item>
      <el-form-item label="负责人">
        <el-input v-model="form.owner" placeholder="请输入负责人，默认招聘专员-刘经理" />
      </el-form-item>
      <el-form-item label="补充说明">
        <el-input
          v-model="form.description"
          type="textarea"
          :rows="3"
          placeholder="可填写薪资结构、奖金、试用期等说明"
        />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="visibleLocal = false">取消</el-button>
      <el-button type="primary" :loading="submitting" @click="handleSubmit">确认发起</el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, reactive, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { createOffer } from '@/api'

const props = defineProps({
  modelValue: { type: Boolean, default: false },
  applications: { type: Array, default: () => [] },
  candidateName: { type: String, default: '' },
  preSelectApplicationId: { type: String, default: '' }
})
const emit = defineEmits(['update:modelValue', 'success'])

const visibleLocal = computed({
  get: () => props.modelValue,
  set: (v) => emit('update:modelValue', v)
})

const formRef = ref(null)
const submitting = ref(false)
const form = reactive({
  applicationId: '',
  salaryMin: 20000,
  salaryMax: 30000,
  startDate: '',
  owner: '招聘专员-刘经理',
  description: ''
})

const rules = {
  applicationId: [{ required: true, message: '请选择投递申请', trigger: 'change' }],
  startDate: [{ required: true, message: '请选择入职日期', trigger: 'change' }]
}

watch(() => props.modelValue, (v) => {
  if (v) {
    form.applicationId = props.preSelectApplicationId || (props.applications[0]?.id || '')
    form.salaryMin = 20000
    form.salaryMax = 30000
    form.startDate = ''
    form.owner = '招聘专员-刘经理'
    form.description = ''
  }
})

const handleClose = () => {
  formRef.value?.resetFields()
}

const handleSubmit = async () => {
  await formRef.value?.validate()
  if (!form.applicationId || !form.startDate) return
  if (form.salaryMax < form.salaryMin) {
    ElMessage.warning('最高薪资不能低于最低薪资')
    return
  }
  submitting.value = true
  try {
    await createOffer({ ...form })
    ElMessage.success('Offer 已创建')
    visibleLocal.value = false
    emit('success')
  } catch (e) {
    ElMessage.error('创建失败：' + (e?.response?.data?.error || e.message))
  } finally {
    submitting.value = false
  }
}
</script>
