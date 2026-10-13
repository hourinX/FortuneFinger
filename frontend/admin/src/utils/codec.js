const encoder = new TextEncoder()

export function encodeParam(value) {
  const bytes = encoder.encode(String(value ?? ''))
  let out = ''
  for (let i = 0; i < bytes.length; i++) {
    out += `%${bytes[i].toString(16).toUpperCase().padStart(2, '0')}`
  }
  return out
}

export function encodeQuery(params) {
  const pairs = []
  for (const [key, value] of Object.entries(params || {})) {
    if (value === undefined || value === null || value === '') continue
    const values = Array.isArray(value) ? value : [value]
    for (const item of values) {
      pairs.push(`${encodeURIComponent(key)}=${encodeParam(item)}`)
    }
  }
  return pairs.join('&')
}

// TODO: x25519 encryption for sensitive data