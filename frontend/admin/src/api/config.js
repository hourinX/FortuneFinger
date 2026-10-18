import request from '@/utils/request'

const API_PREFIX = '/api/config'

export function getConfig(code) {
  return request.get(`${API_PREFIX}/getConfig`, { params: { code } })
}

export function saveConfig(data) {
  return request.put(`${API_PREFIX}/save`, data)
}

export function sendTestMail(data) {
  return request.post(`${API_PREFIX}/testMail`, data)
}
