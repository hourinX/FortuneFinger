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
  const app = useAppStore(pinia)
  const requiresBootstrap =
    !getToken() || !admin.profile || !app.routersLoaded || !app.dictionariesLoaded
  if (requiresBootstrap) app.startFullscreenLoading()
  if (!getToken()) {
    try {
      await admin.restoreSession()
    } catch {
      app.stopFullscreenLoading()
      return { path: '/login', query: { redirect: to.fullPath } }
    }
  } else if (!admin.profile) {
    try {
      await admin.loadProfile()
    } catch {
      app.stopFullscreenLoading()
      return { path: '/login', query: { redirect: to.fullPath } }
    }
  }

  if (!app.routersLoaded) {
    await Promise.all([app.loadRouters(router), app.loadDictionaries()])
    return to.fullPath
  }
  if (!app.dictionariesLoaded) await app.loadDictionaries()
  return true
})

router.afterEach((to, _from, failure) => {
  const app = useAppStore(pinia)
  app.addVisitedRouter(to)
  if (!failure && app.fullscreenLoading) void app.finishFullscreenLoading()
})

router.onError(() => {
  useAppStore(pinia).stopFullscreenLoading()
})

export default router
