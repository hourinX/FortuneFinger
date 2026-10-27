export type CheckinDisplayMode =
  | 'lottery-machine'
  | 'floating-orbs'
  | 'photo-wall'
  | 'recent-list'

export interface CheckinTemplateDefaults {
  accentColor: string
  background: string
  displayMode: CheckinDisplayMode
  fontFamily: string
  layout: 'split'
  showLogo: boolean
  subtitle: string
  title: string
}

export interface CheckinTemplateDefinition {
  id: string
  name: string
  preview?: string
  defaults: CheckinTemplateDefaults
  supportedDisplayModes: CheckinDisplayMode[]
}

export interface CheckinScreenConfig extends CheckinTemplateDefaults {
  templateId: string
  activityLogo?: string
  countdownTarget?: string
}
