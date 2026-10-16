import request from '@/utils/request'

const API_PREFIX = '/api/asset'

export function uploadAsset(file, scene) {
  const data = new FormData()
  data.append('file', file)
  data.append('scene', scene)
  return request.post(`${API_PREFIX}/upload`, data)
}

export function generateAvatar() {
  return request.post(`${API_PREFIX}/generateAvatar`, {})
}
