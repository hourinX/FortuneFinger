import { defineStore } from 'pinia'
import { login as loginApi, logout as logoutApi } from '@/api/auth'
import { getUserInfo } from '@/api/user'
import { clearAuthTokens, getToken, getTokenExpiresAt, setAuthTokens } from '@/utils/auth'
import { refreshAuthTokens } from '@/utils/request'

export const useAdminStore = defineStore('admin', {
  state: () => ({
    token: getToken(),
    expiresAt: getTokenExpiresAt(),
    profile: null,
  }),
  getters: {
    isLoggedIn: (state) => Boolean(state.token),
  },
  actions: {
    async login(payload) {
      const data = await loginApi(payload)
      setAuthTokens(data)
      this.token = getToken()
      this.expiresAt = getTokenExpiresAt()
      await this.loadProfile()
    },
    async restoreSession() {
      const data = await refreshAuthTokens()
      this.token = data.access_token
      this.expiresAt = Number(data.expires_at) || 0
      await this.loadProfile()
    },
    async loadProfile() {
      this.profile = await getUserInfo()
    },
    async logout() {
      try {
        if (getToken()) await logoutApi()
      } catch {
        // 服务端不可用时仍需清理当前页面的登录状态。
      } finally {
        clearAuthTokens()
        this.token = ''
        this.expiresAt = 0
        this.profile = null
      }
    },
  },
})
