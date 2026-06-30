<template>
  <div class="not-found">
    <el-result icon="warning" title="404 页面不存在" sub-title="您访问的页面已搬家或从未存在">
      <template #extra>
        <el-button type="primary" @click="goHome">返回工作台</el-button>
        <el-button @click="goBack">返回上一页</el-button>
      </template>
    </el-result>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'

const router = useRouter()
const goHome = () => router.push('/dashboard')
const goBack = () => router.go(-1)
</script>

<style scoped>
.not-found {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f0f2f5;
}
</style>
