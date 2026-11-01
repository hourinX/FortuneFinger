function normalizeTagStyle(style) {
  if (!style) return undefined
  return {
    color: style.color,
    backgroundColor: style.backgroundColor || style.background,
    borderColor: style.borderColor || style.border,
  }
}

const dictionaryTagPalettes = {
  node_types: {
    default: { color: '#555555', backgroundColor: '#f5f5f5', borderColor: '#e6e6e6' },
    green: { color: '#246b45', backgroundColor: '#f1f8f4', borderColor: '#d5eadc' },
    orange: { color: '#8a5a18', backgroundColor: '#fbf7ee', borderColor: '#eee1c8' },
  },
  prize_level: {
    grand: { color: '#B7791F', background: '#FFF7E6', border: '#F5DBA4' },
    first: { color: '#8B5CF6', background: '#F5F0FF', border: '#DDD0FF' },
    second: { color: '#2563EB', background: '#EFF6FF', border: '#BFDBFE' },
    third: { color: '#0F9F78', background: '#EAFBF3', border: '#B8E8D3' },
    participation: { color: '#64748B', background: '#F1F5F9', border: '#DCE3EB' },
  },
  activities_status: {
    pending: { color: '#D97706', background: '#FFF7E6', border: '#FFD591' },
    active: { color: '#E53935', background: '#FFF1F0', border: '#FFA39E' },
    finished: { color: '#16A34A', background: '#F0FDF4', border: '#BBF7D0' },
  },
}

export function dictionaryTagStyle(option, styles = {}) {
  if (!option) return {}

  const customStyle = styles[option.value] || styles[option.color]
  if (customStyle) return normalizeTagStyle(customStyle)

  if (!option.color || option.color === 'default') return {}
  return {
    color: option.color,
    borderColor: option.color,
  }
}

export function dictionaryTagStyleByCode(code, option) {
  return dictionaryTagStyle(option, dictionaryTagPalettes[code])
}
