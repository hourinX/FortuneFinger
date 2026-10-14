import { computed } from 'vue'
import { useAppStore } from '@/store'

export function useDictionaryOptions(code, numeric = false) {
  const app = useAppStore()
  return computed(() =>
    app.dictOptions(code).map((item) => ({
      ...item,
      value: numeric ? Number(item.value) : item.value,
    })),
  )
}
