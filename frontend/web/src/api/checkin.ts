import request from '@/utils/request'

export interface CheckinTemplateResp {
  template: string
}

export function getCheckinTemplate(signal?: AbortSignal) {
  return request.get<CheckinTemplateResp>('/api/checkin/template', { signal })
}
