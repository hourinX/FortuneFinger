import { createRouter, createWebHistory } from 'vue-router'
import { getToken } from '@/utils/auth'
import { DEFAULT_BRANDING } from '@/constants/branding'
import { pinia, useAdminStore, useAppStore } from '@/store'
import routes from './routes'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
})

router.beforeEach(async (to) => {
  document.title = to.meta.title
    ? `${to.meta.title} - ${DEFAULT_BRANDING.name}`
    : DEFAULT_BRANDING.adminTitle
  if (to.meta.public) return true
  const admin = useAdminStore(pinia)
  if (!getToken()) {
    try {
      await admin.restoreSession()
    } catch {
      return { path: '/login', query: { redirect: to.fullPath } }
    }
  } else if (!admin.profile) {
    try {
      await admin.loadProfile()
    } catch {
      return { path: '/login', query: { redirect: to.fullPath } }
    }
  }

  const app = useAppStore(pinia)
  if (!app.routersLoaded) {
    await app.loadRouters(router)
    return to.fullPath
  }
  return true
})

router.afterEach((to) => {
  const app = useAppStore(pinia)
  app.addVisitedRouter(to)
})

export default router
