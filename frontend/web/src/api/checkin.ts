import request from '@/utils/request'

export interface CheckinTemplateResp {
  template: string
}

export type CheckinStatus = 'pending' | 'active' | 'ended'

export interface CheckinStateResp {
  status: CheckinStatus
  server_time: string
  activity_start_time: string
  activity_end_time: string
  checkin_start_time: string
  checkin_end_time: string
  checked_in: number
  total: number | null
}

export interface CheckinQRCodeResp {
  token: string
  qrcode: string
  expires_at: string
  refresh_after: number
}

export interface CheckinStateMessage {
  key: 'checkin.state'
  data: CheckinStateResp
}

export function getCheckinTemplate(code: string, signal?: AbortSignal) {
  return request.get<CheckinTemplateResp>('/api/checkin/template', { params: { code }, signal })
}

export function getCheckinState(code: string, signal?: AbortSignal) {
  return request.get<CheckinStateResp>('/api/checkin/state', { params: { code }, signal })
}

export function getCheckinQRCode(code: string, signal?: AbortSignal) {
  return request.get<CheckinQRCodeResp>('/api/checkin/qrcode', { params: { code }, signal })
}

export function createCheckinStateSocket(code: string) {
  const baseURL = import.meta.env.VITE_API_BASE_URL || window.location.origin
  const url = new URL('/api/ws', baseURL)
  url.protocol = url.protocol === 'https:' ? 'wss:' : 'ws:'
  url.searchParams.set('key', 'checkin.state')
  url.searchParams.set('code', code)
  return new WebSocket(url)
}
