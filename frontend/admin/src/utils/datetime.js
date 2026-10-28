import dayjs from 'dayjs'
import 'dayjs/locale/zh-cn'

export const DATE_FORMAT = 'YYYY-MM-DD'
export const DATE_TIME_FORMAT = 'YYYY-MM-DD HH:mm:ss'
export const DAY_MILLISECONDS = 24 * 60 * 60 * 1000

export function setupDateTimeLocale() {
  dayjs.locale('zh-cn')
}

export function formatDateTime(value, format = DATE_TIME_FORMAT) {
  if (!value) return ''
  const date = dayjs(value)
  return date.isValid() ? date.format(format) : ''
}

export function getDateTimeDifference(start, end) {
  const startTime = dayjs(start)
  const endTime = dayjs(end)
  if (!startTime.isValid() || !endTime.isValid()) return Number.NaN
  return endTime.diff(startTime)
}
