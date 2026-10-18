import { defineStore } from 'pinia'
import { getDictOptions } from '@/api/dict'
import { getPermission, getRouters } from '@/api/menu'
import { registerRouters, resetDynamicRoutes } from '@/router/dynamic'

const VISITED_ROUTERS_KEY = 'fortunefinger-visited-routers'
const VISITED_ROUTER_STATES_KEY = 'fortunefinger-visited-router-states'

function restoreVisitedRouters() {
  try {
    const value = JSON.parse(sessionStorage.getItem(VISITED_ROUTERS_KEY) || '[]')
    if (!Array.isArray(value)) return []
    return value.filter((item) => item?.id && item?.title && item?.route)
  } catch {
    return []
  }
}

function restoreVisitedRouterStates() {
  try {
    const value = JSON.parse(sessionStorage.getItem(VISITED_ROUTER_STATES_KEY) || '{}')
    return value && typeof value === 'object' && !Array.isArray(value) ? value : {}
  } catch {
    return {}
  }
}

export const useAppStore = defineStore('app', {
  state: () => ({
    collapsed: false,
    permissions: [],
    dictionaries: {},
    routers: [],
    visitedRouters: restoreVisitedRouters(),
    visitedRouterStates: restoreVisitedRouterStates(),
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
      delete this.visitedRouterStates[String(id)]
      this.persistVisitedRouters()
      this.persistVisitedRouterStates()
    },
    persistVisitedRouters() {
      sessionStorage.setItem(VISITED_ROUTERS_KEY, JSON.stringify(this.visitedRouters))
    },
    getVisitedRouterState(id) {
      return this.visitedRouterStates[String(id)] || null
    },
    setVisitedRouterState(id, state) {
      if (!id) return
      this.visitedRouterStates[String(id)] = state
      this.persistVisitedRouterStates()
    },
    persistVisitedRouterStates() {
      sessionStorage.setItem(VISITED_ROUTER_STATES_KEY, JSON.stringify(this.visitedRouterStates))
    },
    async loadRouters(router) {
      if (this.routersLoaded) return
      const [routerData, permissionData] = await Promise.all([getRouters(), getPermission()])
      const routers = Array.isArray(routerData?.list) ? routerData.list : []
      this.permissions = Array.isArray(permissionData?.list) ? permissionData.list : []
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
      this.visitedRouterStates = Object.fromEntries(
        Object.entries(this.visitedRouterStates).filter(([id]) => menuIDs.has(String(id))),
      )
      this.persistVisitedRouters()
      this.persistVisitedRouterStates()
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
      this.visitedRouterStates = {}
      sessionStorage.removeItem(VISITED_ROUTER_STATES_KEY)
      this.permissions = []
      this.dictionaries = {}
      this.routersLoaded = false
      this.dictionariesLoaded = false
    },
  },
})
