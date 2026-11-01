import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    { path: '/', redirect: '/checkin' },
    {
      path: '/checkin',
      name: 'participant-checkin',
      component: () => import('@/views/checkin/index.vue'),
      meta: { title: '活动签到' },
    },
    {
      path: '/home',
      alias: '/chat',
      name: 'participant-home',
      component: () => import('@/views/chat/index.vue'),
      meta: { title: '活动现场' },
    },
    {
      path: '/lottery',
      name: 'participant-lottery',
      component: () => import('@/views/placeholder/index.vue'),
      props: {
        active: 'lottery',
        eyebrow: 'LUCKY DRAW',
        title: '抽奖现场',
        description: '抽奖开始后，轮次与中奖结果将在这里同步展示。',
      },
      meta: { title: '抽奖现场' },
    },
    {
      path: '/profile',
      name: 'participant-profile',
      component: () => import('@/views/profile/index.vue'),
      meta: { title: '我的活动' },
    },
    { path: '/:pathMatch(.*)*', redirect: '/checkin' },
  ],
})

router.afterEach((to) => {
  document.title = `${String(to.meta.title || '活动签到')} - FortuneFinger`
})

export default router
