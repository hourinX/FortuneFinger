import { markRaw } from 'vue'
import GrandGalaTemplate from './components/GrandGalaTemplate.vue'
import InvitationTemplate from './components/InvitationTemplate.vue'
import CelebrationTemplate from './components/CelebrationTemplate.vue'
import LuxeTemplate from './components/LuxeTemplate.vue'

import GlassTemplate from './components/GlassTemplate.vue'
import SlotTemplate from './components/SlotTemplate.vue'
import OrbitTemplate from './components/OrbitTemplate.vue'

// 管理端保存稳定的 id，大屏通过此注册表选择对应组件。
export const checkinTemplates = [
  { id: 'invitation', name: '相聚邀请函', color: '#e5e5d3', component: markRaw(InvitationTemplate) },
  { id: 'celebration', name: '欢聚盛典', color: '#ab252d', component: markRaw(CelebrationTemplate) },
  { id: 'luxe', name: '静雅之夜', color: '#bba57e', component: markRaw(LuxeTemplate) },
  { id: 'glass', name: '幸运制造所', color: '#adbdff', component: markRaw(GlassTemplate) },
  { id: 'slot', name: '鎏金游乐场', color: '#d4b27f', component: markRaw(SlotTemplate) },
  { id: 'orbit', name: '同频星系', color: '#74d6d0', component: markRaw(OrbitTemplate) },
  { id: 'grand-gala', name: '同心致远', color: '#e72813', component: markRaw(GrandGalaTemplate) },
] as const

export const defaultTemplateId = 'invitation'
export function resolveTemplate(id: unknown) {
  return checkinTemplates.find(template => template.id === id) ?? checkinTemplates[0]
}
