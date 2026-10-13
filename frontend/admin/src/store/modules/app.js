import { defineStore } from 'pinia'
import { getRouters } from '@/api/menu'
import { registerRouters, resetDynamicRoutes } from '@/router/dynamic'

export const useAppStore = defineStore('app', {
  state: () => ({
    collapsed: false,
    permissions: [],
    routers: [],
    visitedRouters: [],
    routersLoaded: false,
  }),
  getters: {
    hasPermission: (state) => (code) => state.permissions.includes(code),
  },
  actions: {
    toggleCollapsed() {
      this.collapsed = !this.collapsed
    },
    addVisitedRouter(route) {
      if (!route.meta?.routerId) return
      if (this.visitedRouters.some((item) => item.id === route.meta.routerId)) return
      this.visitedRouters.push({
        id: route.meta.routerId,
        title: route.meta.title,
        route: route.path,
      })
    },
    removeVisitedRouter(id) {
      const index = this.visitedRouters.findIndex((item) => item.id === id)
      if (index !== -1) this.visitedRouters.splice(index, 1)
    },
    async loadRouters(router) {
      if (this.routersLoaded) return
      const data = await getRouters()
      const routers = Array.isArray(data?.list) ? data.list : []
      registerRouters(router, routers)
      this.routers = routers
      this.routersLoaded = true
    },
    resetRouters() {
      resetDynamicRoutes()
      this.routers = []
      this.visitedRouters = []
      this.permissions = []
      this.routersLoaded = false
    },
  },
})
