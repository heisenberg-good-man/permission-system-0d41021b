<template>
  <el-dialog
    v-model="visibleLocal"
    :title="editingRow ? '编辑用人需求' : '新建用人需求'"
    width="720px"
    destroy-on-close
    @close="handleClose"
  >
    <el-form
      ref="formRef"
      :model="form"
      :rules="formRules"
      label-width="110px"
      @validate="onFieldValidate"
    >
      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item label="部门" prop="department">
            <el-select v-model="form.department" placeholder="请选择部门" style="width: 100%" clearable>
              <el-option v-for="dept in departments" :key="dept" :label="dept" :value="dept" />
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="岗位名称" prop="jobTitle">
            <el-input v-model="form.jobTitle" placeholder="请输入岗位名称" clearable />
          </el-form-item>
        </el-col>
      </el-row>

      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item label="编制类型" prop="headcountType">
            <el-radio-group v-model="form.headcountType">
              <el-radio value="regular">正式编制</el-radio>
              <el-radio value="contract">合同工</el-radio>
              <el-radio value="intern">实习生</el-radio>
              <el-radio value="replacement">替补编制</el-radio>
            </el-radio-group>
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="需求人数" prop="headcount">
            <el-input-number v-model="form.headcount" :min="1" :max="100" style="width: 100%" />
          </el-form-item>
        </el-col>
      </el-row>

      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item label="期望到岗日期" prop="expectedOnboard">
            <el-date-picker
              v-model="form.expectedOnboard"
              type="date"
              placeholder="请选择期望到岗日期"
              style="width: 100%"
              format="YYYY-MM-DD"
              value-format="YYYY-MM-DD"
              :disabled-date="disabledDate"
            />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="发起人" prop="initiator">
            <el-input v-model="form.initiator" placeholder="默认当前登录用户" />
          </el-form-item>
        </el-col>
      </el-row>

      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item label="工作地点" prop="workLocation">
            <el-input v-model="form.workLocation" placeholder="请输入工作地点（如：北京）" clearable />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="发起人联系方式" prop="initiatorContact">
            <el-input v-model="form.initiatorContact" placeholder="请输入手机号" clearable />
          </el-form-item>
        </el-col>
      </el-row>

      <el-form-item label="预算薪资范围" required>
        <div style="display: flex; align-items: center; gap: 12px; width: 100%">
          <el-form-item prop="salaryMin" style="flex: 1; margin-bottom: 0">
            <el-input-number v-model="form.salaryMin" :min="0" :step="1" style="width: 100%" placeholder="最低薪资(K)" />
          </el-form-item>
          <span style="color: #909399">至</span>
          <el-form-item prop="salaryMax" style="flex: 1; margin-bottom: 0">
            <el-input-number v-model="form.salaryMax" :min="0" :step="1" style="width: 100%" placeholder="最高薪资(K)" />
          </el-form-item>
          <span style="color: #909399">K</span>
        </div>
      </el-form-item>

      <el-form-item label="招聘原因" prop="reason">
        <el-input v-model="form.reason" type="textarea" :rows="2" placeholder="请说明招聘原因，如业务扩张、补位等" />
      </el-form-item>

      <el-form-item label="岗位说明" prop="jobDescription">
        <el-input
          v-model="form.jobDescription"
          type="textarea"
          :rows="3"
          placeholder="请输入岗位说明、主要职责等"
        />
      </el-form-item>

      <el-form-item label="能力要求" prop="requirements">
        <el-input
          v-model="form.requirements"
          type="textarea"
          :rows="3"
          placeholder="请输入任职要求、技能要求等"
        />
      </el-form-item>

      <el-alert
        v-if="duplicateWarning"
        type="warning"
        :closable="false"
        show-icon
        style="margin-bottom: 20px"
      >
        <template #title>
          <span>检测到同部门已存在相同岗位的需求，请确认是否重复提交</span>
        </template>
      </el-alert>
    </el-form>

    <template #footer>
      <el-button @click="visibleLocal = false">取消</el-button>
      <el-button type="primary" :loading="submitting" @click="handleSubmit">
        {{ editingRow ? '保存修改' : '保存需求' }}
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, reactive, computed, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import dayjs from 'dayjs'
import { createRequirement, updateRequirement } from '@/api'

const props = defineProps({
  modelValue: { type: Boolean, default: false },
  editingRow: { type: Object, default: null },
  existingRequirements: { type: Array, default: () => [] }
})
const emit = defineEmits(['update:modelValue', 'success'])

const visibleLocal = computed({
  get: () => props.modelValue,
  set: (v) => emit('update:modelValue', v)
})

const departments = ['技术部', '产品部', '市场部', '运营部', '人力资源部', '财务部', '设计部']

const formRef = ref(null)
const submitting = ref(false)

const defaultForm = {
  department: '',
  jobTitle: '',
  headcountType: 'regular',
  headcount: 1,
  expectedOnboard: '',
  initiator: '部门负责人',
  initiatorContact: '13800138000',
  workLocation: '北京',
  salaryMin: 15,
  salaryMax: 25,
  reason: '',
  jobDescription: '',
  requirements: ''
}

const form = reactive({ ...defaultForm })

const validateSalaryRange = (rule, value, callback) => {
  if (value === undefined || value === null || value === '') {
    callback(new Error('请填写薪资范围'))
    return
  }
  if (form.salaryMax < form.salaryMin) {
    callback(new Error('最高薪资不能低于最低薪资'))
    return
  }
  callback()
}

const validateExpectedDate = (rule, value, callback) => {
  if (!value) {
    callback(new Error('请选择期望到岗日期'))
    return
  }
  const today = dayjs().startOf('day')
  const selected = dayjs(value).startOf('day')
  if (selected.isBefore(today)) {
    callback(new Error('期望到岗日期不能早于今天'))
    return
  }
  callback()
}

const formRules = {
  department: [{ required: true, message: '请选择部门', trigger: 'change' }],
  jobTitle: [
    { required: true, message: '请输入岗位名称', trigger: 'blur' },
    { min: 2, max: 50, message: '岗位名称长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  headcountType: [{ required: true, message: '请选择编制类型', trigger: 'change' }],
  headcount: [
    { required: true, message: '请输入需求人数', trigger: 'blur' },
    {
      validator: (rule, value, callback) => {
        if (value < 1) {
          callback(new Error('需求人数至少为 1 人'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ],
  expectedOnboard: [
    { required: true, validator: validateExpectedDate, trigger: 'change' }
  ],
  salaryMin: [
    { required: true, message: '请输入最低薪资', trigger: 'blur' },
    { validator: validateSalaryRange, trigger: 'blur' }
  ],
  salaryMax: [
    { required: true, message: '请输入最高薪资', trigger: 'blur' },
    { validator: validateSalaryRange, trigger: 'blur' }
  ],
  jobDescription: [
    { required: true, message: '请输入岗位说明', trigger: 'blur' }
  ],
  requirements: [
    { required: true, message: '请输入能力要求', trigger: 'blur' }
  ]
}

const disabledDate = (time) => {
  const today = dayjs().startOf('day')
  return time.getTime() < today.valueOf()
}

const duplicateWarning = computed(() => {
  if (!form.department || !form.jobTitle) return false
  return props.existingRequirements.some(req => {
    if (props.editingRow && req.id === props.editingRow.id) return false
    if (req.status === 'rejected' || req.status === 'converted') return false
    return req.department === form.department && req.jobTitle === form.jobTitle
  })
})

const onFieldValidate = (prop, isValid, message) => {
}

const kToYuan = (k) => Math.round(Number(k) * 1000)
const yuanToK = (y) => y >= 1000 ? Number(y) / 1000 : Number(y)

watch(() => props.modelValue, (v) => {
  if (v) {
    if (props.editingRow) {
      const row = { ...props.editingRow }
      if (row.jobTitle) row.jobTitle = row.jobTitle
      row.salaryMin = yuanToK(row.salaryMin)
      row.salaryMax = yuanToK(row.salaryMax)
      Object.assign(form, row)
    } else {
      Object.assign(form, { ...defaultForm })
    }
  }
})

const handleClose = () => {
  formRef.value?.resetFields()
}

const handleSubmit = async () => {
  let valid = true
  try {
    await formRef.value?.validate()
  } catch (e) {
    valid = false
  }
  if (!valid) {
    ElMessage.warning('请检查表单填写是否正确')
    return
  }

  if (form.salaryMax < form.salaryMin) {
    ElMessage.warning('最高薪资不能低于最低薪资')
    return
  }

  if (form.expectedOnboard) {
    const today = dayjs().startOf('day')
    const selected = dayjs(form.expectedOnboard).startOf('day')
    if (selected.isBefore(today)) {
      ElMessage.warning('期望到岗日期不能早于今天')
      return
    }
  }

  if (duplicateWarning.value) {
    try {
      await ElMessageBox.confirm(
        '同部门已存在相同岗位的需求，是否仍然继续提交？',
        '重复提示',
        { type: 'warning', confirmButtonText: '继续提交', cancelButtonText: '取消' }
      )
    } catch {
      return
    }
  }

  submitting.value = true
  try {
    const payload = { ...form }
    payload.salaryMin = kToYuan(payload.salaryMin)
    payload.salaryMax = kToYuan(payload.salaryMax)
    if (props.editingRow) {
      await updateRequirement(props.editingRow.id, payload)
      ElMessage.success('需求更新成功')
    } else {
      await createRequirement(payload)
      ElMessage.success('需求创建成功')
    }
    visibleLocal.value = false
    emit('success')
  } catch (e) {
    ElMessage.error('保存失败：' + (e?.response?.data?.error || e.message))
  } finally {
    submitting.value = false
  }
}
</script>
