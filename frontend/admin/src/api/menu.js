import request from '@/utils/request'
import { encodeParam } from '@/utils/codec'

const API_PREFIX = '/api/menu'

export function getMenuList(params) {
  return request.get(`${API_PREFIX}/list`, { params })
}

export function getMenuDetail(id) {
  return request.get(`${API_PREFIX}/detail/${encodeParam(id)}`)
}

export function createMenu(data) {
  return request.post(`${API_PREFIX}/create`, data)
}

export function updateMenu(id, data) {
  return request.put(`${API_PREFIX}/update/${encodeParam(id)}`, data)
}

export function deleteMenu(id) {
  return request.delete(`${API_PREFIX}/${encodeParam(id)}`)
}

export function getRouters() {
  return request.get(`${API_PREFIX}/getRouters`)
}

export function getPermission() {
  return request.get(`${API_PREFIX}/getPermission`)
}

export function getPermissionTree(roleId) {
  return request.get(`${API_PREFIX}/permission/${encodeParam(roleId)}`)
}
