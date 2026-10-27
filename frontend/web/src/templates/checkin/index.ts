import { classicCheckinTemplate } from './classic'
import type { CheckinScreenConfig, CheckinTemplateDefinition } from './types'

const templates: CheckinTemplateDefinition[] = [classicCheckinTemplate]

export function getCheckinTemplate(id: string) {
  return templates.find((template) => template.id === id) || classicCheckinTemplate
}

export function createCheckinScreenConfig(
  templateId: string,
  config: Partial<CheckinScreenConfig> = {},
): CheckinScreenConfig {
  const template = getCheckinTemplate(templateId)
  return {
    templateId: template.id,
    ...template.defaults,
    ...config,
  }
}

export type { CheckinDisplayMode, CheckinScreenConfig, CheckinTemplateDefinition } from './types'
