import { defineStore } from 'pinia'
import { getDictOptions } from '@/api/dict'
import { getRouters } from '@/api/menu'
import { registerRouters, resetDynamicRoutes } from '@/router/dynamic'

export const useAppStore = defineStore('app', {
  state: () => ({
    collapsed: false,
    permissions: [],
    dictionaries: {},
    routers: [],
    visitedRouters: [],
    routersLoaded: false,
    dictionariesLoaded: false,
  }),
  getters: {
    hasPermission: (state) => (code) => state.permissions.includes(code),
    dictOptions: (state) => (code) => state.dictionaries[code] || [],
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
    async loadDictionaries(force = false) {
      if (this.dictionariesLoaded && !force) return
      this.dictionaries = (await getDictOptions()) || {}
      this.dictionariesLoaded = true
    },
    resetRouters() {
      resetDynamicRoutes()
      this.routers = []
      this.visitedRouters = []
      this.permissions = []
      this.dictionaries = {}
      this.routersLoaded = false
      this.dictionariesLoaded = false
    },
  },
})
