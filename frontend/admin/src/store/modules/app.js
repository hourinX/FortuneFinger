import { defineStore } from 'pinia'
import { getDictOptions } from '@/api/dict'
import { getRouters } from '@/api/menu'
import { registerRouters, resetDynamicRoutes } from '@/router/dynamic'

const VISITED_ROUTERS_KEY = 'fortunefinger-visited-routers'

function restoreVisitedRouters() {
  try {
    const value = JSON.parse(sessionStorage.getItem(VISITED_ROUTERS_KEY) || '[]')
    if (!Array.isArray(value)) return []
    return value.filter((item) => item?.id && item?.title && item?.route)
  } catch {
    return []
  }
}

export const useAppStore = defineStore('app', {
  state: () => ({
    collapsed: false,
    permissions: [],
    dictionaries: {},
    routers: [],
    visitedRouters: restoreVisitedRouters(),
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
    restoreDefaults() {
      this.collapsed = false
    },
    addVisitedRouter(route) {
      if (!route.meta?.routerId) return
      const id = String(route.meta.routerId)
      if (this.visitedRouters.some((item) => item.id === id)) return
      this.visitedRouters.push({
        id,
        title: route.meta.title,
        route: route.path,
      })
      this.persistVisitedRouters()
    },
    removeVisitedRouter(id) {
      const index = this.visitedRouters.findIndex((item) => item.id === id)
      if (index === -1) return
      this.visitedRouters.splice(index, 1)
      this.persistVisitedRouters()
    },
    persistVisitedRouters() {
      sessionStorage.setItem(VISITED_ROUTERS_KEY, JSON.stringify(this.visitedRouters))
    },
    async loadRouters(router) {
      if (this.routersLoaded) return
      const data = await getRouters()
      const routers = Array.isArray(data?.list) ? data.list : []
      registerRouters(router, routers)
      this.routers = routers
      const menuIDs = new Set()
      const collectMenuIDs = (items) => {
        for (const item of items || []) {
          if (item.types === 2 && item.route) menuIDs.add(String(item.id))
          collectMenuIDs(item.children)
        }
      }
      collectMenuIDs(routers)
      this.visitedRouters = this.visitedRouters.filter((item) => menuIDs.has(String(item.id)))
      this.persistVisitedRouters()
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
      sessionStorage.removeItem(VISITED_ROUTERS_KEY)
      this.permissions = []
      this.dictionaries = {}
      this.routersLoaded = false
      this.dictionariesLoaded = false
    },
  },
})
