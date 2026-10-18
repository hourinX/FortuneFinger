import request from '@/utils/request'
import { encodeParam } from '@/utils/codec'

const API_PREFIX = '/api/prize/pool'

export function getPrizeList(params) {
  return request.get(`${API_PREFIX}/list`, { params })
}

export function getPrizeDetail(id) {
  return request.get(`${API_PREFIX}/detail/${encodeParam(id)}`)
}

export function createPrize(data) {
  return request.post(`${API_PREFIX}/create`, data)
}

export function updatePrize(id, data) {
  return request.put(`${API_PREFIX}/update/${encodeParam(id)}`, data)
}

export function deletePrize(id) {
  return request.delete(`${API_PREFIX}/${encodeParam(id)}`)
}
