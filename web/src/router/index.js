import { createRouter, createWebHashHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    redirect: '/login'
  },
  {
    path: '/init',
    name: 'Init',
    component: () => import('@/view/init/index.vue')
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/view/login/index.vue')
  },
  {
    path: '/scanUpload',
    name: 'ScanUpload',
    meta: {
      title: '扫码上传',
      client: true
    },
    component: () => import('@/view/example/upload/scanUpload.vue')
  },
  {
    path: '/:catchAll(.*)',
    meta: {
      closeTab: true
    },
    component: () => import('@/view/error/index.vue')
  },
  {
    path: '/home',
    name: 'Home',
    component: () => import('@/view/home/index.vue')
  },
  {
    path: '/my',
    name: 'My',
    component: () => import('@/view/my/index.vue')
  },
  {
    path: '/mustRead',
    name: 'MustRead',
    component: () => import('@/view/mustRead/index.vue')
  },
  {
    path: '/changePassword',
    name: 'ChangePassword',
    component: () => import('@/view/changePassword/index.vue')
  },
  {
    path: '/memberManagement',
    name: 'MemberManagement',
    component: () => import('@/view/memberManagement/index.vue')
  },
  {
    path: '/albumDetail/:id',
    name: 'AlbumDetail',
    component: () => import('@/view/albumDetail/index.vue')
  }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

export default router
