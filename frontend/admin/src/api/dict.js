import request from '@/utils/request'
import { encodeParam } from '@/utils/codec'

const API_PREFIX = '/api/dict'

export function getDictTypeList(params) {
  return request.get(`${API_PREFIX}/type/list`, { params })
}

export function getDictTypeDetail(id) {
  return request.get(`${API_PREFIX}/type/detail/${encodeParam(id)}`)
}

export function createDictType(data) {
  return request.post(`${API_PREFIX}/type/create`, data)
}

export function updateDictType(id, data) {
  return request.put(`${API_PREFIX}/type/update/${encodeParam(id)}`, data)
}

export function deleteDictType(id) {
  return request.delete(`${API_PREFIX}/type/${encodeParam(id)}`)
}

export function getDictDataList(params) {
  return request.get(`${API_PREFIX}/data/list`, { params })
}

export function getDictDataDetail(id) {
  return request.get(`${API_PREFIX}/data/detail/${encodeParam(id)}`)
}

export function createDictData(data) {
  return request.post(`${API_PREFIX}/data/create`, data)
}

export function updateDictData(id, data) {
  return request.put(`${API_PREFIX}/data/update/${encodeParam(id)}`, data)
}

export function deleteDictData(id) {
  return request.delete(`${API_PREFIX}/data/${encodeParam(id)}`)
}

export function getDictOptions(codes) {
  return request.get(`${API_PREFIX}/options`, {
    params: { codes: Array.isArray(codes) ? codes.join(',') : codes || undefined },
  })
}
