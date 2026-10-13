import request from '@/utils/request'

const API_PREFIX = '/api/auth'
const PUBLIC_OPTIONS = {
  skipAccessToken: true,
  skipAuthRefresh: true,
}

export function login(data) {
  return request.post(`${API_PREFIX}/login`, data, PUBLIC_OPTIONS)
}

export function logout() {
  return request.post(`${API_PREFIX}/logout`, {})
}

export function getCaptcha() {
	return request.get(`${API_PREFIX}/getCaptcha`, PUBLIC_OPTIONS)
}

export function createEnhancedVerification(data) {
  return request.post(`${API_PREFIX}/verification/create`, data, PUBLIC_OPTIONS)
}

export function verifyEnhancedVerification(data) {
  return request.post(`${API_PREFIX}/verification/verify`, data, PUBLIC_OPTIONS)
}
