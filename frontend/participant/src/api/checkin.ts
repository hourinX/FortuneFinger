import request from '@/utils/request'

export interface CheckinSessionResp {
  session_id: string
  expires_at: string
  activity: {
    name: string
    company: string
    description: string
  }
  config: {
    mode: 1 | 2
    require_name: boolean
    require_company: boolean
    allow_avatar: boolean
    allow_message: boolean
  }
}

export function createCheckinSession(token: string, signal?: AbortSignal) {
  return request.post<CheckinSessionResp, { token: string }>(
    '/api/checkin/session',
    { token },
    { signal },
  )
}
