import { useAppStore } from '@/store'

// TODO: 接入 getPermissions，并在权限预加载完成后支持超级管理员通配符。
export default {
  mounted(el, binding) {
    const required = Array.isArray(binding.value) ? binding.value : [binding.value]
    const codes = required.filter((code) => typeof code === 'string' && code)
    const allowed = codes.length > 0 && codes.some((code) => useAppStore().hasPermission(code))
    if (!allowed && el.parentNode) el.parentNode.removeChild(el)
  },
}
