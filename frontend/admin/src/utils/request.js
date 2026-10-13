import axios from 'axios'
import { clearAuthStorage, getToken, setAuthTokens } from './auth'
import { encodeQuery } from './codec'

const AUTH_MESSAGE_KEY = 'fortune_admin_auth_message'
let refreshPromise = null

const request = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '',
  timeout: 15000,
  withCredentials: true,
  paramsSerializer: {
    serialize: encodeQuery,
  },
})

function getErrMsg(data, fallback = '请求失败') {
  return data?.message || data?.msg || fallback
}

function redirectToLogin(text) {
  clearAuthStorage()
  if (window.location.pathname === '/login') return
  try {
    sessionStorage.setItem(AUTH_MESSAGE_KEY, text)
  } catch {
    console.warn('sessionStorage is unavailable')
  }
  window.location.replace('/login')
}

function refreshToken() {
  return request.post(
    '/api/auth/refresh',
    {},
    { skipAccessToken: true, skipAuthRefresh: true },
  )
}

export function refreshAuthTokens() {
  if (refreshPromise) return refreshPromise

  refreshPromise = refreshToken()
    .then((data) => {
      if (!data?.access_token) throw new Error('刷新登录状态失败')
      setAuthTokens(data)
      return data
    })
    .finally(() => {
      refreshPromise = null
    })
  return refreshPromise
}

async function retryAfterRefresh(config, text) {
  if (!config || config.skipAuthRefresh || config.authRetried) {
    redirectToLogin(text)
    throw new Error(text)
  }
  try {
    await refreshAuthTokens()
    config.authRetried = true
    return request(config)
  } catch (error) {
    redirectToLogin(error.message || text)
    throw error
  }
}

request.interceptors.request.use((config) => {
  const token = getToken()
  if (token && !config.skipAccessToken) config.headers.Authorization = `Bearer ${token}`
  return config
})

request.interceptors.response.use(
  (response) => {
    const data = response.data
    if (Number(data?.code) === 2) {
      return retryAfterRefresh(response.config, getErrMsg(data, '登录状态已失效，请重新登录'))
    }
    if (data?.code !== undefined && Number(data.code) !== 0) {
      const error = new Error(getErrMsg(data))
      error.code = data.code
      error.data = data.data
      return Promise.reject(error)
    }
    return data?.data !== undefined ? data.data : data
  },
  (error) => {
    const data = error.response?.data
    const text = getErrMsg(data, error.message || '网络请求失败')
    if (error.response?.status === 401) return retryAfterRefresh(error.config, text)
    error.message = text
    if (data?.code !== undefined) error.code = data.code
    return Promise.reject(error)
  },
)

export { AUTH_MESSAGE_KEY }
export default request
