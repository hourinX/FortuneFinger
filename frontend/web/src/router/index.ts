import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      redirect: '/checkin',
    },
    {
      path: '/checkin/:eventId?',
      name: 'checkin-screen',
      component: () => import('@/views/checkin/index.vue'),
      meta: { title: '签到大屏' },
    },
    {
      path: '/lottery/:eventId?',
      name: 'lottery-screen',
      component: () => import('@/views/lottery/index.vue'),
      meta: { title: '抽奖大屏' },
    },
    {
      path: '/:pathMatch(.*)*',
      redirect: '/checkin',
    },
  ],
})

router.afterEach((to) => {
  document.title = `${String(to.meta.title || '活动大屏')} - FortuneFinger`
})

export default router
