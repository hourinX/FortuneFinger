import { useAppStore } from '@/store'
export default {
  mounted(el, binding) {
    const required = Array.isArray(binding.value) ? binding.value : [binding.value]
    const codes = required.filter((code) => typeof code === 'string' && code)
    const allowed = codes.length > 0 && codes.some((code) => useAppStore().hasPermission(code))
    if (!allowed && el.parentNode) el.parentNode.removeChild(el)
  },
}
