import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { createCheckinSession, type CheckinSessionResp } from '@/api/checkin'
import { getOrgCompanies, type OrgCompanyOption } from '@/api/org'

export interface CheckinFormValue {
  name: string
  companyId: string
  companyName: string
  phone: string
  avatar?: File
}

const previewSession: CheckinSessionResp = {
  session_id: 'preview-session',
  expires_at: new Date(Date.now() + 5 * 60_000).toISOString(),
  activity: {
    name: '2026 FortuneFinger 年度盛典',
    company: 'FortuneFinger',
    description: '欢迎来到年度相聚现场，完成签到后即可参与今晚的活动。',
  },
  config: {
    mode: 1,
    require_name: true,
    require_company: true,
    allow_avatar: true,
    allow_message: false,
  },
}

export const useParticipantCheckinStore = defineStore('participant-checkin', () => {
  const session = ref<CheckinSessionResp>()
  const companies = ref<OrgCompanyOption[]>([])
  const loading = ref(false)
  const error = ref('')
  const ready = computed(() => Boolean(session.value) && !loading.value)

  async function load(token: string, preview = false) {
    loading.value = true
    error.value = ''
    try {
      session.value = preview ? previewSession : await createCheckinSession(token)
    } catch (reason) {
      error.value = reason instanceof Error ? reason.message : '签到入口加载失败'
      loading.value = false
      return
    }

    try {
      const resp = await getOrgCompanies()
      companies.value = resp.list
    } catch {
      companies.value = []
    } finally {
      loading.value = false
    }
  }

  return { session, companies, loading, error, ready, load }
})
