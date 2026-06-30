<template>
  <el-container class="layout-container">
    <el-aside width="220px" class="layout-aside">
      <div class="logo">
        <el-icon class="logo-icon"><OfficeBuilding /></el-icon>
        <span class="logo-text">招聘工作台</span>
      </div>
      <el-menu
        :default-active="activeMenu"
        router
        class="layout-menu"
        background-color="#001529"
        text-color="#ffffffa6"
        active-text-color="#ffffff"
      >
        <el-menu-item v-for="item in menuItems" :key="item.path" :index="item.path">
          <el-icon><component :is="item.icon" /></el-icon>
          <span>{{ item.title }}</span>
        </el-menu-item>
      </el-menu>
    </el-aside>

    <el-container>
      <el-header class="layout-header">
        <div class="header-left">
          <el-icon class="user-icon"><User /></el-icon>
          <span class="user-name">招聘管理员</span>
          <el-tag size="small" type="success" class="role-tag">招聘方</el-tag>
        </div>
        <div class="header-right">
          <el-tag size="small" type="info" round>
            今日 {{ currentDate }}</el-tag>
        </div>
      </el-header>

      <el-main class="layout-main">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { OfficeBuilding, User, UserFilled, DataAnalysis, PieChart, Briefcase, Document, ChatDotRound, Calendar, Tickets, Monitor, DocumentAdd } from '@element-plus/icons-vue'
import dayjs from 'dayjs'

const route = useRoute()

const activeMenu = computed(() => route.path)

const menuItems = [
  { path: '/dashboard', title: '工作台', icon: Monitor },
  { path: '/stats', title: '统计看板', icon: PieChart },
  { path: '/jobs', title: '职位管理', icon: Briefcase },
  { path: '/requirements', title: '用人需求', icon: DocumentAdd },
  { path: '/candidates', title: '候选人管理', icon: UserFilled },
  { path: '/applications', title: '申请管理', icon: Document },
  { path: '/interviews', title: '面试安排', icon: Calendar },
  { path: '/offers', title: 'Offer 管理', icon: Tickets },
  { path: '/messages', title: '沟通中心', icon: ChatDotRound }
]

const currentDate = dayjs().format('YYYY年MM月DD日')
</script>

<style scoped>
.layout-container {
  height: 100%;
}

.layout-aside {
  background-color: #001529;
  overflow: hidden;
  flex-shrink: 0;
}

.logo {
  height: 60px;
  display: flex;
  align-items: center;
  padding: 0 20px;
  color: #fff;
  border-bottom: 1px solid #ffffff1a;
}

.logo-icon {
  font-size: 24px;
  color: #409eff;
  margin-right: 10px;
}

.logo-text {
  font-size: 18px;
  font-weight: 600;
  color: #fff;
}

.layout-menu {
  border-right: none;
  height: calc(100vh - 60px);
  overflow-y: auto;
}

.layout-menu :deep(.el-menu-item) {
  height: 50px;
  line-height: 50px;
}

.layout-menu :deep(.el-menu-item:hover) {
  background-color: #ffffff1a !important;
}

.layout-menu :deep(.el-menu-item.is-active) {
  background-color: #409eff !important;
}

.layout-header {
  background: #fff;
  border-bottom: 1px solid #e6e6e6;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  height: 60px;
  flex-shrink: 0;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 10px;
}

.user-icon {
  font-size: 20px;
  color: #666;
}

.user-name {
  font-size: 14px;
  color: #333;
  font-weight: 500;
}

.role-tag {
  margin-left: 4px;
}

.header-right {
  display: flex;
  align-items: center;
}

.layout-main {
  background-color: #f0f2f5;
  padding: 20px;
  overflow-y: auto;
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
