<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { computed, nextTick, ref, watch } from 'vue'
import { useCheckinEventsStore } from '@/stores/checkin-events'

const props = withDefaults(defineProps<{
  stage?: 'checkin' | 'waiting' | 'lottery' | 'result'
  side?: 'left' | 'right'
  limit?: number
}>(), { stage: 'checkin', side: 'right', limit: 100 })

const { recent } = storeToRefs(useCheckinEventsStore())
const listRef = ref<HTMLElement>()
const visibleEvents = computed(() => recent.value.slice(0, props.limit).reverse())

function formatTime(value: string) {
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return ''
  return new Intl.DateTimeFormat('zh-CN', {
    hour: '2-digit', minute: '2-digit', second: '2-digit', hour12: false,
  }).format(date)
}

watch(
  () => visibleEvents.value.length,
  async () => {
    await nextTick()
    listRef.value?.scrollTo({ top: listRef.value.scrollHeight, behavior: 'smooth' })
  },
)
</script>

<template>
  <aside
    v-if="visibleEvents.length"
    class="checkin-activity"
    :class="[`checkin-activity--${side}`, { 'checkin-activity--compact': stage !== 'checkin' }]"
    aria-label="签到动态"
  >
    <header class="checkin-activity__header">
      <span><i></i>签到动态</span>
      <small>{{ visibleEvents.length }}</small>
    </header>
    <div ref="listRef" class="checkin-activity__list">
      <TransitionGroup name="activity">
        <article
          v-for="(item, index) in visibleEvents"
          :key="item.eventId"
          class="checkin-activity__item"
          :class="{ 'checkin-activity__item--latest': index === visibleEvents.length - 1 }"
        >
          <div class="checkin-activity__avatar">
            <img v-if="item.avatarUrl" :src="item.avatarUrl" alt="" />
            <span v-else>{{ item.displayName.slice(0, 1) }}</span>
          </div>
          <div class="checkin-activity__summary">
            <strong>{{ item.displayName }}</strong>
            <span>完成了签到</span>
          </div>
          <time :datetime="item.checkedInAt">{{ formatTime(item.checkedInAt) }}</time>
        </article>
      </TransitionGroup>
    </div>
  </aside>
</template>

<style scoped>
.checkin-activity {
  position: fixed;
  z-index: 125;
  top: 50%;
  display: flex;
  width: clamp(300px, 19vw, 370px);
  height: 68vh;
  padding: 15px 12px 10px;
  flex-direction: column;
  overflow: hidden;
  color: #fff;
  background: linear-gradient(145deg, rgb(24 26 31 / 68%), rgb(13 15 19 / 52%));
  border: 1px solid rgb(255 255 255 / 13%);
  border-radius: 16px;
  box-shadow: 0 18px 55px rgb(0 0 0 / 12%), inset 0 1px rgb(255 255 255 / 5%);
  transform: translateY(-50%);
}
.checkin-activity--right { right: 1.6vw; }
.checkin-activity--left { left: 1.6vw; }

.checkin-activity__header {
  display: flex;
  height: 38px;
  padding: 0 7px 12px;
  flex: 0 0 auto;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1px solid rgb(255 255 255 / 8%);
}
.checkin-activity__header > span { display: flex; align-items: center; gap: 9px; font-size: 14px; font-weight: 650; letter-spacing: .08em; }
.checkin-activity__header i { width: 6px; height: 6px; background: #f05a55; border-radius: 50%; box-shadow: 0 0 9px rgb(240 90 85 / 75%); animation: live-pulse 1.8s ease-in-out infinite; }
.checkin-activity__header small { min-width: 24px; padding: 2px 7px; color: #fff; font-size: 10px; text-align: center; background: #d94343; border-radius: 10px; }

.checkin-activity__list {
  display: flex;
  min-height: 0;
  padding: 12px 5px 26px;
  flex: 1;
  flex-direction: column;
  overflow-y: auto;
  overscroll-behavior: contain;
  scrollbar-width: none;
  mask-image: linear-gradient(to bottom, transparent 0%, #000 12%, #000 100%);
}
.checkin-activity__list::-webkit-scrollbar { display: none; }

.checkin-activity__item {
  display: flex;
  width: min(100%, 305px);
  min-height: 56px;
  margin: 0 auto;
  padding: 10px 7px;
  flex: 0 0 auto;
  align-items: center;
  gap: 10px;
  border-bottom: 1px solid rgb(255 255 255 / 6%);
  border-radius: 9px;
}
.checkin-activity__item--latest { animation: latest-glow 1s ease-out; }
.checkin-activity__item:first-child { margin-top: auto; }
.checkin-activity__avatar { display: grid; width: 32px; height: 32px; flex: 0 0 32px; overflow: hidden; place-items: center; color: #fff; font-size: 11px; font-weight: 650; background: linear-gradient(145deg, #e76e68, #b9343c); border: 1px solid rgb(255 255 255 / 24%); border-radius: 50%; }
.checkin-activity__avatar img { width: 100%; height: 100%; object-fit: cover; }
.checkin-activity__summary { display: flex; min-width: 0; flex: 1; align-items: baseline; justify-content: center; gap: 6px; white-space: nowrap; }
.checkin-activity__summary strong { overflow: hidden; color: rgb(255 255 255 / 94%); font-size: 13px; font-weight: 650; text-overflow: ellipsis; }
.checkin-activity__summary span { color: rgb(255 255 255 / 58%); font-size: 11px; }
.checkin-activity__item time { flex: 0 0 auto; color: rgb(255 255 255 / 34%); font-size: 9px; font-variant-numeric: tabular-nums; }
.checkin-activity--compact { width: clamp(250px, 16vw, 300px); height: 48vh; opacity: .76; }
.activity-enter-active, .activity-move { transition: opacity .4s ease, transform .4s ease; }
.activity-enter-from { opacity: 0; transform: translateY(18px); }

@keyframes latest-glow { 0% { background: rgb(238 190 105 / 18%); box-shadow: inset 0 0 0 1px rgb(245 211 149 / 24%); } 100% { background: transparent; box-shadow: none; } }
@keyframes live-pulse { 50% { opacity: .45; transform: scale(.78); } }

@media (max-width: 900px) {
  .checkin-activity { top: auto; right: 14px; bottom: 14px; left: 14px; width: auto; height: 42vh; transform: none; }
}
</style>
