import request from '@/utils/request'

export interface OrgCompanyOption {
  id: string
  name: string
}

interface GetOrgCompaniesResp {
  list: OrgCompanyOption[]
}

export function getOrgCompanies(signal?: AbortSignal) {
  return request.get<GetOrgCompaniesResp>('/api/org/companies', { signal })
}
