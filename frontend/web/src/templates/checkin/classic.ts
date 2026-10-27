import type { CheckinTemplateDefinition } from './types'

export const classicCheckinTemplate: CheckinTemplateDefinition = {
  id: 'classic-dark',
  name: '经典深色',
  defaults: {
    accentColor: '#ef4444',
    background: '#080a11',
    displayMode: 'photo-wall',
    fontFamily: 'Inter, "PingFang SC", "Microsoft YaHei", sans-serif',
    layout: 'split',
    showLogo: true,
    subtitle: '签到进行中 · 欢迎入场',
    title: '2026 FortuneFinger 年度盛典',
  },
  supportedDisplayModes: ['lottery-machine', 'floating-orbs', 'photo-wall', 'recent-list'],
}
