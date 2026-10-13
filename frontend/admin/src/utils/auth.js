let accessToken = ''
let accessExpiresAt = 0

const legacyStorageKeys = [
  'fortune_admin_access_token',
  'fortune_admin_refresh_token',
  'fortune_admin_token_expires_at',
  'fortune_admin_refresh_expires_at',
]
legacyStorageKeys.forEach((key) => localStorage.removeItem(key))

export const getToken = () => accessToken
export const getTokenExpiresAt = () => accessExpiresAt

export function setAuthTokens(data = {}) {
  accessToken = data.access_token || ''
  accessExpiresAt = Number(data.expires_at) || 0
}

export function clearAuthTokens() {
  accessToken = ''
  accessExpiresAt = 0
}

export const clearAuthStorage = clearAuthTokens
