import { defineStore } from 'pinia'
import { computed, ref } from 'vue'

export const useScreenStore = defineStore('screen', () => {
  const serverTimeOffset = ref(0)
  const connected = ref(true)
  const now = computed(() => new Date(Date.now() + serverTimeOffset.value))

  function syncServerTime(serverTimestamp: number) {
    serverTimeOffset.value = serverTimestamp - Date.now()
  }

  function setConnected(value: boolean) {
    connected.value = value
  }

  return {
    connected,
    now,
    serverTimeOffset,
    setConnected,
    syncServerTime,
  }
})
