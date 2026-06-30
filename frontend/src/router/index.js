import { createRouter, createWebHashHistory } from 'vue-router'
import Layout from '@/layout/Layout.vue'

const routes = [
  {
    path: '/',
    component: Layout,
    redirect: '/stats',
    children: [
      {
        path: 'stats',
        name: 'Stats',
        component: () => import('@/views/Stats.vue'),
        meta: { title: '统计概览', icon: 'DataAnalysis' }
      },
      {
        path: 'jobs',
        name: 'Jobs',
        component: () => import('@/views/Jobs.vue'),
        meta: { title: '职位管理', icon: 'Briefcase' }
      },
      {
        path: 'applications',
        name: 'Applications',
        component: () => import('@/views/Applications.vue'),
        meta: { title: '申请管理', icon: 'Document' }
      },
      {
        path: 'messages',
        name: 'Messages',
        component: () => import('@/views/Messages.vue'),
        meta: { title: '沟通中心', icon: 'ChatDotRound' }
      },
      {
        path: 'candidates',
        name: 'Candidates',
        component: () => import('@/views/Candidates.vue'),
        meta: { title: '候选人管理', icon: 'UserFilled' }
      },
      {
        path: 'candidates/detail',
        name: 'CandidateDetail',
        component: () => import('@/views/CandidateDetail.vue'),
        meta: { title: '候选人详情', icon: 'UserFilled', hidden: true }
      },
      {
        path: 'interviews',
        name: 'Interviews',
        component: () => import('@/views/Interviews.vue'),
        meta: { title: '面试安排', icon: 'Calendar' }
      },
      {
        path: 'offers',
        name: 'Offers',
        component: () => import('@/views/Offers.vue'),
        meta: { title: 'Offer 管理', icon: 'Tickets' }
      }
    ]
  }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

export default router
