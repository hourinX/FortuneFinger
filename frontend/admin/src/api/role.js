import request from '@/utils/request'
import { encodeParam } from '@/utils/codec'

const API_PREFIX = '/api/role'

export function getRoleList(params) {
  return request.get(`${API_PREFIX}/list`, { params })
}

export function getRoleDetail(id) {
  return request.get(`${API_PREFIX}/detail/${encodeParam(id)}`)
}

export function createRole(data) {
  return request.post(`${API_PREFIX}/create`, data)
}

export function updateRole(id, data) {
  return request.put(`${API_PREFIX}/update/${encodeParam(id)}`, data)
}

export function deleteRole(id) {
  return request.delete(`${API_PREFIX}/${encodeParam(id)}`)
}

export function getRoleOptions() {
  return request.get(`${API_PREFIX}/options`)
}

export function updateRolePermission(id, data) {
  return request.put(`${API_PREFIX}/permission/${encodeParam(id)}`, data)
}
