<template>
  <div id="app-root">
    <ServiceDown v-if="!backendHealth.ok && !backendHealth.checking && backendHealth.lastCheck" />
    <div v-else>
      <div v-if="!backendHealth.ok && backendHealth.checking" class="connecting-tip">
        <el-icon class="loading-icon"><Loading /></el-icon>
        <span>正在连接后端服务...</span>
      </div>
      <router-view v-slot="{ Component }">
        <transition name="fade" mode="out-in">
          <component :is="Component" />
        </transition>
      </router-view>
    </div>
  </div>
</template>

<script setup>
import { onMounted, onBeforeUnmount, watch } from 'vue'
import { ElNotification } from 'element-plus'
import { Loading } from '@element-plus/icons-vue'
import ServiceDown from '@/components/ServiceDown.vue'
import { backendHealth, startHealthCheck, stopHealthCheck, checkBackendHealth } from '@/api/health'

let notifiedDown = false
let notifiedRecovered = false

onMounted(async () => {
  const ok = await checkBackendHealth()
  if (!ok) {
    ElNotification.warning({
      title: '后端服务未就绪',
      message: '请确认后端服务已启动并监听 8080 端口',
      duration: 5000,
      position: 'top-right'
    })
  }
  startHealthCheck(5000)

  watch(
    () => backendHealth.value.ok,
    (newVal, oldVal) => {
      if (oldVal === true && newVal === false && !notifiedDown) {
        notifiedDown = true
        notifiedRecovered = false
        ElNotification.error({
          title: '后端服务断开',
          message: '请检查后端服务是否正常运行',
          duration: 0
        })
      }
      if (oldVal === false && newVal === true && !notifiedRecovered) {
        notifiedRecovered = true
        notifiedDown = false
        ElNotification.success({
          title: '后端服务已恢复',
          message: '可以继续操作了',
          duration: 3000
        })
      }
    }
  )
})

onBeforeUnmount(() => {
  stopHealthCheck()
})
</script>

<style>
#app-root {
  height: 100%;
  width: 100%;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
}

.connecting-tip {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 9999;
  background: #fff7e6;
  color: #e6a23c;
  padding: 8px 20px;
  font-size: 13px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.loading-icon {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
