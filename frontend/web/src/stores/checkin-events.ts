import { defineStore } from 'pinia'
import { computed, ref, shallowRef } from 'vue'

export interface CheckinArrivalEvent {
  eventId: string
  checkinId: string
  displayName: string
  companyName?: string
  avatarUrl?: string
  message?: string
  checkedInAt: string
}

interface CheckinArrivalSnapshot {
  sequence: number
  event: CheckinArrivalEvent
}

const recentLimit = 120

export const useCheckinEventsStore = defineStore('checkin-events', () => {
  const recent = ref<CheckinArrivalEvent[]>([])
  const latest = shallowRef<CheckinArrivalSnapshot>()
  const sequence = ref(0)
  const knownEventIds = new Set<string>()

  const hasEvents = computed(() => recent.value.length > 0)

  function pushArrival(event: CheckinArrivalEvent) {
    if (!event.eventId || knownEventIds.has(event.eventId)) return

    knownEventIds.add(event.eventId)
    recent.value = [event, ...recent.value].slice(0, recentLimit)
    sequence.value += 1
    latest.value = { sequence: sequence.value, event }

    if (knownEventIds.size > recentLimit * 3) {
      const retained = new Set(recent.value.map(item => item.eventId))
      knownEventIds.forEach((id) => {
        if (!retained.has(id)) knownEventIds.delete(id)
      })
    }
  }

  function replaceRecent(events: CheckinArrivalEvent[]) {
    recent.value = []
    knownEventIds.clear()
    events.slice(0, recentLimit).reverse().forEach(pushArrival)
  }

  function clear() {
    recent.value = []
    latest.value = undefined
    knownEventIds.clear()
  }

  return {
    recent,
    latest,
    hasEvents,
    pushArrival,
    replaceRecent,
    clear,
  }
})
