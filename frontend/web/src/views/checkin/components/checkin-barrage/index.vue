<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { onBeforeUnmount, onMounted, ref, watch } from 'vue'
import {
  useCheckinEventsStore,
  type CheckinArrivalEvent,
} from '@/stores/checkin-events'

interface BarrageItem extends CheckinArrivalEvent {
  instanceId: number
  lane: number
  duration: number
}

const props = withDefaults(defineProps<{
  active?: boolean
  laneCount?: number
  maxVisible?: number
}>(), {
  active: true,
  laneCount: 6,
  maxVisible: 9,
})

const store = useCheckinEventsStore()
const { latest } = storeToRefs(store)
const items = ref<BarrageItem[]>([])
const removalTimers = new Map<number, number>()
const pendingEvents: CheckinArrivalEvent[] = []
let instanceSequence = 0
let nextLane = 0
let pumpTimer: number | undefined

function removeItem(instanceId: number) {
  items.value = items.value.filter(item => item.instanceId !== instanceId)
  const timer = removalTimers.get(instanceId)
  if (timer !== undefined) window.clearTimeout(timer)
  removalTimers.delete(instanceId)
}

function clearItems() {
  removalTimers.forEach(timer => window.clearTimeout(timer))
  removalTimers.clear()
  pendingEvents.length = 0
  items.value = []
}

function enqueue(event: CheckinArrivalEvent) {
  const laneCount = Math.max(1, props.laneCount)
  const instanceId = ++instanceSequence
  const duration = 7.8 + ((instanceId * 17) % 48) / 10
  const item: BarrageItem = {
    ...event,
    instanceId,
    lane: nextLane % laneCount,
    duration,
  }
  nextLane += 1
  items.value = [...items.value.slice(-(props.maxVisible - 1)), item]
  removalTimers.set(instanceId, window.setTimeout(() => removeItem(instanceId), duration * 1000))
}

function pumpQueue() {
  if (props.active && pendingEvents.length && items.value.length < props.maxVisible) {
    const event = pendingEvents.shift()
    if (event) enqueue(event)
  }
  pumpTimer = window.setTimeout(pumpQueue, 180 + Math.floor(Math.random() * 1450))
}

function itemStyle(item: BarrageItem) {
  const top = 3 + (item.lane * 15) / Math.max(1, props.laneCount - 1)
  return {
    '--barrage-top': `${top}%`,
    '--barrage-duration': `${item.duration}s`,
  }
}

watch(
  () => latest.value?.sequence,
  () => {
    if (!props.active || !latest.value) return
    pendingEvents.push(latest.value.event)
  },
)

watch(
  () => props.active,
  (active) => {
    if (!active) clearItems()
  },
)

onMounted(() => {
  pumpQueue()
})

onBeforeUnmount(() => {
  window.clearTimeout(pumpTimer)
  clearItems()
})
</script>

<template>
  <div
    v-if="active"
    class="checkin-barrage"
    aria-hidden="true"
  >
    <TransitionGroup name="barrage">
      <div
        v-for="item in items"
        :key="item.instanceId"
        class="checkin-barrage__item"
        :style="itemStyle(item)"
      >
        <span class="checkin-barrage__avatar">
          <img
            v-if="item.avatarUrl"
            :src="item.avatarUrl"
            alt=""
          />
          <span v-else>{{ item.displayName.slice(0, 1) }}</span>
        </span>
        <span class="checkin-barrage__text">
          <strong>{{ item.displayName }}：</strong>
          <span>{{ item.message || '很高兴和大家相聚在这里' }}</span>
        </span>
      </div>
    </TransitionGroup>
  </div>
</template>

<style scoped>
.checkin-barrage {
  position: fixed;
  z-index: 120;
  inset: 0;
  overflow: hidden;
  pointer-events: none;
}

.checkin-barrage__item {
  position: absolute;
  top: var(--barrage-top);
  left: 100%;
  display: flex;
  max-width: min(560px, 52vw);
  min-height: 32px;
  padding: 4px 8px 4px 4px;
  align-items: center;
  gap: 8px;
  color: #fffdf6;
  white-space: nowrap;
  filter: drop-shadow(0 2px 5px rgb(0 0 0 / 35%));
  animation: barrage-flow var(--barrage-duration) linear forwards;
  will-change: transform;
}

.checkin-barrage__avatar {
  display: grid;
  width: 24px;
  height: 24px;
  flex: 0 0 24px;
  overflow: hidden;
  place-items: center;
  color: #fff;
  font-size: 10px;
  font-weight: 650;
  background: linear-gradient(145deg, #ef6c67, #d9363e);
  border: 1px solid rgb(255 255 255 / 40%);
  border-radius: 50%;
}

.checkin-barrage__avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.checkin-barrage__text {
  display: flex;
  min-width: 0;
  align-items: baseline;
  gap: 3px;
  text-shadow: 0 1px 4px rgb(0 0 0 / 76%), 0 0 12px rgb(255 232 185 / 16%);
}

.checkin-barrage__text strong {
  color: #f0c879;
  font-size: 12px;
  font-weight: 700;
}

.checkin-barrage__text > span {
  overflow: hidden;
  color: #fff;
  font-size: clamp(14px, 1vw, 18px);
  font-weight: 500;
  text-overflow: ellipsis;
  letter-spacing: 0.02em;
}

@keyframes barrage-flow {
  from {
    transform: translate3d(5vw, -50%, 0);
  }
  to {
    transform: translate3d(calc(-100vw - 120%), -50%, 0);
  }
}

.barrage-enter-active {
  transition: opacity 0.28s ease;
}

.barrage-enter-from {
  opacity: 0;
}

@media (max-width: 700px) {
  .checkin-barrage__item {
    max-width: 76vw;
  }
}

@media (prefers-reduced-motion: reduce) {
  .checkin-barrage__item {
    animation-duration: calc(var(--barrage-duration) * 1.6);
  }
}
</style>
