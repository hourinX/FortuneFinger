import request from '@/utils/request'
import { encodeParam } from '@/utils/codec'

const API_PREFIX = '/api/user'

export function getUserList(params) {
  return request.get(`${API_PREFIX}/list`, { params })
}

export function getUserDetail(id) {
  return request.get(`${API_PREFIX}/detail/${encodeParam(id)}`)
}

export function createUser(data) {
  return request.post(`${API_PREFIX}/create`, data)
}

export function updateUser(id, data) {
  return request.put(`${API_PREFIX}/update/${encodeParam(id)}`, data)
}

export function deleteUser(id) {
  return request.delete(`${API_PREFIX}/${encodeParam(id)}`)
}

export function assignUserRole(id, data) {
  return request.put(`${API_PREFIX}/role/${encodeParam(id)}`, data)
}

export function getUserInfo() {
  return request.get(`${API_PREFIX}/getUserInfo`)
}
