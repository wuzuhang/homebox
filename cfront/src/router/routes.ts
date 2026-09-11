import type { RouteRecordRaw } from 'vue-router'
import Layout from '@/layout/index.vue'

const routes: Array<RouteRecordRaw> = [
  {
    path: '/auth',
    name: 'Auth',
    component: () => import('@/views/auth/index.vue'),
    meta: {
      title: '登录 · Home Box',
      noCache: true,
    },
  },
  {
    path: '/',
    name: 'root',
    component: Layout,
    redirect: { name: 'Auth' },
    meta: {
      requiresAuth: true,
    },
    children: [
      {
        path: 'home',
        name: 'Home',
        component: () => import('@/views/demo/index.vue'),
        meta: {
          title: '主页',
        },
      },
      {
        path: 'profile',
        name: 'ProfileDetail',
        component: () => import('@/views/profile-detail/index.vue'),
        meta: {
          title: '个人信息',
          noCache: true,
          hideTabbar: true,
        },
      },
      {
        path: 'profile/edit',
        name: 'ProfileEdit',
        component: () => import('@/views/profile-edit/index.vue'),
        meta: {
          title: '编辑基础信息',
          noCache: true,
          hideTabbar: true,
        },
      },
      {
        path: 'medicines',
        name: 'Medicines',
        component: () => import('@/views/medicines/index.vue'),
        meta: {
          title: '我的药品',
          noCache: true,
          hideTabbar: true,
        },
      },
      {
        path: 'medicines/new',
        name: 'MedicineCreate',
        component: () => import('@/views/medicine-form/index.vue'),
        meta: {
          title: '新增药品',
          noCache: true,
          hideTabbar: true,
        },
      },
      {
        path: 'medicines/:id/edit',
        name: 'MedicineEdit',
        component: () => import('@/views/medicine-form/index.vue'),
        meta: {
          title: '编辑药品',
          noCache: true,
          hideTabbar: true,
        },
      },
      {
        path: 'medicines/:id',
        name: 'MedicineDetail',
        component: () => import('@/views/medicine-detail/index.vue'),
        meta: {
          title: '药品详情',
          noCache: true,
          hideTabbar: true,
        },
      },
      {
        path: 'tools',
        name: 'Tools',
        component: () => import('@/views/tools/index.vue'),
        meta: {
          title: '工具',
        },
      },
      {
        path: 'about',
        name: 'About',
        component: () => import('@/views/about/index.vue'),
        meta: {
          title: '关于',
          noCache: true,
        },
      },
    ],
  },
  // 404 页面
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/404.vue'),
    meta: {
      title: '页面未找到',
      noCache: true,
    },
  },
]

export default routes
