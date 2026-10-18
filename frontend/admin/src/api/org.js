import request from '@/utils/request'
import { encodeParam } from '@/utils/codec'

const API_PREFIX = '/api/org'

export function getOrgList(params) {
  return request.get(`${API_PREFIX}/list`, { params })
}

export function getOrgDetail(id) {
  return request.get(`${API_PREFIX}/detail/${encodeParam(id)}`)
}

export function createOrg(data) {
  return request.post(`${API_PREFIX}/create`, data)
}

export function updateOrg(id, data) {
  return request.put(`${API_PREFIX}/update/${encodeParam(id)}`, data)
}

export function deleteOrg(id) {
  return request.delete(`${API_PREFIX}/${encodeParam(id)}`)
}

export function getOrgOptions() {
  return request.get(`${API_PREFIX}/options`)
}

export function getOrgCompanies() {
  return request.get(`${API_PREFIX}/companies`)
}
