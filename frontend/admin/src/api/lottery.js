import request from '@/utils/request'
import { encodeParam } from '@/utils/codec'

const ACTIVITY_API_PREFIX = '/api/lottery/activity'
const ROUND_API_PREFIX = '/api/lottery/round'

export function getActivityList(params) {
  return request.get(`${ACTIVITY_API_PREFIX}/list`, { params })
}

export function getActivityDetail(id) {
  return request.get(`${ACTIVITY_API_PREFIX}/detail/${encodeParam(id)}`)
}

export function createActivity(data) {
  return request.post(`${ACTIVITY_API_PREFIX}/create`, data)
}

export function updateActivity(id, data) {
  return request.put(`${ACTIVITY_API_PREFIX}/update/${encodeParam(id)}`, data)
}

export function updateActivityCheckinTemplate(id, data) {
  return request.put(`${ACTIVITY_API_PREFIX}/template/update/${encodeParam(id)}`, data)
}

export function deleteActivity(id) {
  return request.delete(`${ACTIVITY_API_PREFIX}/${encodeParam(id)}`)
}

export function getRoundList(params) {
  return request.get(`${ROUND_API_PREFIX}/list`, { params })
}

export function getRoundDetail(id) {
  return request.get(`${ROUND_API_PREFIX}/detail/${encodeParam(id)}`)
}

export function createRound(data) {
  return request.post(`${ROUND_API_PREFIX}/create`, data)
}

export function updateRound(id, data) {
  return request.put(`${ROUND_API_PREFIX}/update/${encodeParam(id)}`, data)
}

export function deleteRound(id) {
  return request.delete(`${ROUND_API_PREFIX}/${encodeParam(id)}`)
}
