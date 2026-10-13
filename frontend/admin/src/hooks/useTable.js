import { reactive, ref } from 'vue'

// TODO: 等普通分页页面稳定后，按 page、pageSize、list、total 约定重写并统一接入。
export function useTable(fetchList) {
  const loading = ref(false)
  const list = ref([])
  const page = reactive({ current: 1, size: 10, total: 0 })

  async function loadData() {
    loading.value = true
    try {
      const data = (await fetchList({ ...page })) || {}
      list.value = data.list || data.records || []
      page.total = Number(data.total) || 0
    } finally {
      loading.value = false
    }
  }

  function changePage(current, size = page.size) {
    page.current = current
    page.size = size
    return loadData()
  }

  return { loading, list, page, loadData, changePage }
}
