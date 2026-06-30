import { ref } from 'vue'
import request from './request'

export const backendHealth = ref({
  ok: false,
  checking: false,
  lastCheck: null,
  error: null
})

let checkTimer = null

export const checkBackendHealth = async () => {
  backendHealth.value.checking = true
  try {
    const res = await request.get('/health')
    backendHealth.value.ok = res.status === 'ok'
    backendHealth.value.error = null
    backendHealth.value.lastCheck = Date.now()
    return true
  } catch (e) {
    backendHealth.value.ok = false
    backendHealth.value.error = e.message || '后端服务不可用'
    backendHealth.value.lastCheck = Date.now()
    return false
  } finally {
    backendHealth.value.checking = false
  }
}

export const startHealthCheck = (intervalMs = 5000) => {
  stopHealthCheck()
  checkBackendHealth()
  checkTimer = setInterval(checkBackendHealth, intervalMs)
}

export const stopHealthCheck = () => {
  if (checkTimer) {
    clearInterval(checkTimer)
    checkTimer = null
  }
}
