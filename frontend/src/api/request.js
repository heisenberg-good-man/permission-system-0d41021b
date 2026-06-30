import axios from 'axios'
import { ElMessage } from 'element-plus'

const request = axios.create({
  baseURL: '/api',
  timeout: 10000
})

request.interceptors.response.use(
  response => {
    return response.data
  },
  error => {
    console.error('Request error:', error)
    let msg = '请求失败'
    if (!error.response) {
      if (error.code === 'ECONNABORTED' || error.message.includes('timeout')) {
        msg = '请求超时，请检查后端服务是否正常启动（8080 端口）'
      } else if (error.message) {
        msg = '网络错误：' + error.message
      } else {
        msg = '后端服务不可用，请启动后端后重试'
      }
    } else if (error.response.status === 404) {
      msg = '请求的接口不存在'
    } else if (error.response.status === 500) {
      msg = '后端服务异常，请联系管理员'
    } else {
      msg = error.response?.data?.error || ('HTTP ' + error.response.status + ' 错误')
    }
    ElMessage.error(msg)
    return Promise.reject(error)
  }
)

export default request
