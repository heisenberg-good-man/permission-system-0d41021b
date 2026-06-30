<template>
  <div class="service-down">
    <div class="error-card">
      <el-icon class="error-icon"><WarningFilled /></el-icon>
      <h2 class="error-title">{{ title }}</h2>
      <p class="error-desc">{{ description }}</p>
      <div class="error-actions">
        <el-button type="primary" @click="handleRetry" :loading="loading" :icon="Refresh">
          重新连接
        </el-button>
        <el-button @click="goHome">返回首页</el-button>
      </div>
      <div class="error-tips">
        <h4>排查建议：</h4>
        <ul>
          <li>确认后端服务已启动：<code>cd backend && go run main.go</code></li>
          <li>确认后端监听端口为 <code>8080</code></li>
          <li>检查防火墙是否拦截本地连接</li>
          <li>查看后端控制台日志，确认无报错</li>
        </ul>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { WarningFilled, Refresh } from '@element-plus/icons-vue'
import { checkBackendHealth } from '@/api/health'

const props = defineProps({
  title: { type: String, default: '后端服务暂不可用' },
  description: { type: String, default: '无法连接到招聘平台后端 API 服务，请确认后端已启动并监听 8080 端口。' }
})

const router = useRouter()
const loading = ref(false)

const handleRetry = async () => {
  loading.value = true
  try {
    const ok = await checkBackendHealth()
    if (ok) {
      ElMessage.success('连接恢复，正在刷新页面...')
      setTimeout(() => window.location.reload(), 800)
    } else {
      ElMessage.error('仍然无法连接到后端服务')
    }
  } finally {
    loading.value = false
  }
}

const goHome = () => {
  router.push('/')
}
</script>

<style scoped>
.service-down {
  min-height: 60vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
}

.error-card {
  max-width: 560px;
  width: 100%;
  padding: 40px;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  text-align: center;
}

.error-icon {
  font-size: 64px;
  color: #f56c6c;
  margin-bottom: 16px;
}

.error-title {
  font-size: 24px;
  font-weight: 600;
  color: #303133;
  margin: 0 0 12px;
}

.error-desc {
  font-size: 14px;
  color: #606266;
  line-height: 1.6;
  margin: 0 0 24px;
}

.error-actions {
  display: flex;
  gap: 12px;
  justify-content: center;
  margin-bottom: 28px;
}

.error-tips {
  text-align: left;
  background: #f5f7fa;
  border-radius: 8px;
  padding: 16px 20px;
}

.error-tips h4 {
  font-size: 14px;
  font-weight: 600;
  color: #303133;
  margin: 0 0 10px;
}

.error-tips ul {
  margin: 0;
  padding-left: 20px;
}

.error-tips li {
  font-size: 13px;
  color: #606266;
  line-height: 1.8;
}

.error-tips code {
  background: #e4e7ed;
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 12px;
  color: #606266;
}
</style>
