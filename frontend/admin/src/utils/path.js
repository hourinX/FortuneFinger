export function normalizePath(path = '') {
  if (!path) return ''
  const normalized = `/${String(path).replace(/^\/+|\/+$/g, '')}`
  return normalized === '/' ? '/' : normalized.replace(/\/{2,}/g, '/')
}
