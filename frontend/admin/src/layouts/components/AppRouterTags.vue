<template>
  <div v-if="app.visitedRouters.length" class="router-tags-shell">
    <button
      v-show="hasOverflow"
      type="button"
      class="router-tags__nav"
      :disabled="!canScrollLeft"
      aria-label="向左查看标签"
      @click="scroll(-1)"
    >
      <left-outlined />
    </button>
    <div ref="trackRef" class="router-tags" @scroll="updateScrollState" @wheel="handleWheel">
      <button
        v-for="item in app.visitedRouters"
        :key="item.id"
        type="button"
        class="router-tag"
        :class="{ 'router-tag--active': route.path === item.route }"
        @click="open(item.route)"
      >
        <span>{{ item.title }}</span>
        <close-outlined class="router-tag__close" @click.stop="close(item)" />
      </button>
    </div>
    <button
      v-show="hasOverflow"
      type="button"
      class="router-tags__nav"
      :disabled="!canScrollRight"
      aria-label="向右查看标签"
      @click="scroll(1)"
    >
      <right-outlined />
    </button>
  </div>
</template>

<script setup>
import { nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { CloseOutlined, LeftOutlined, RightOutlined } from '@ant-design/icons-vue'
import { useRoute, useRouter } from 'vue-router'
import { useAppStore } from '@/store'

const app = useAppStore()
const route = useRoute()
const router = useRouter()
const trackRef = ref()
const hasOverflow = ref(false)
const canScrollLeft = ref(false)
const canScrollRight = ref(false)
let resizeObserver
let revealFrame

function open(path) {
  if (path !== route.path) router.push(path)
}

function close(item) {
  const index = app.visitedRouters.findIndex((routerItem) => routerItem.id === item.id)
  const isActive = route.path === item.route
  app.removeVisitedRouter(item.id)
  if (!isActive || !app.visitedRouters.length) return
  const target = app.visitedRouters[Math.max(0, index - 1)]
  router.push(target.route)
}

function updateScrollState() {
  const track = trackRef.value
  if (!track) return
  hasOverflow.value = track.scrollWidth > track.clientWidth + 1
  canScrollLeft.value = track.scrollLeft > 1
  canScrollRight.value = track.scrollLeft + track.clientWidth < track.scrollWidth - 1
}

function scroll(direction) {
  trackRef.value?.scrollBy({ left: direction * 180, behavior: 'smooth' })
}

function handleWheel(event) {
  if (!hasOverflow.value) return
  event.preventDefault()
  trackRef.value.scrollLeft += event.deltaY || event.deltaX
}

async function revealActiveTag() {
  await nextTick()
  const track = trackRef.value
  const active = track?.querySelector('.router-tag--active')
  if (track && active) {
    const left = active.offsetLeft
    const right = left + active.offsetWidth
    if (left < track.scrollLeft) track.scrollLeft = left
    else if (right > track.scrollLeft + track.clientWidth) {
      track.scrollLeft = right - track.clientWidth
    }
  }
  updateScrollState()
}

function scheduleRevealActiveTag() {
  cancelAnimationFrame(revealFrame)
  revealFrame = requestAnimationFrame(revealActiveTag)
}

watch([() => route.path, () => app.visitedRouters.length], scheduleRevealActiveTag, {
  flush: 'post',
})

onMounted(() => {
  resizeObserver = new ResizeObserver(updateScrollState)
  if (trackRef.value) resizeObserver.observe(trackRef.value)
  scheduleRevealActiveTag()
})

onBeforeUnmount(() => {
  cancelAnimationFrame(revealFrame)
  resizeObserver?.disconnect()
})
</script>

<style scoped>
.router-tags-shell {
  display: flex;
  align-items: center;
  gap: 5px;
  min-width: 0;
  min-height: 26px;
  margin: 5px 20px 0;
}
.router-tags {
  display: flex;
  flex: 1;
  align-items: center;
  gap: 6px;
  min-width: 0;
  overflow-x: auto;
  scrollbar-width: none;
}
.router-tags::-webkit-scrollbar {
  display: none;
}
.router-tags__nav {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  flex: 0 0 24px;
  width: 24px;
  height: 26px;
  padding: 0;
  color: #737880;
  font-size: 10px;
  background: #fff;
  border: 1px solid #eceef1;
  border-radius: 4px;
  cursor: pointer;
}
.router-tags__nav:hover:not(:disabled) {
  color: var(--color-primary);
  border-color: var(--color-primary-border);
}
.router-tags__nav:disabled {
  color: #c7c9cd;
  cursor: default;
}
.router-tag {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  flex: 0 0 auto;
  height: 26px;
  padding: 0 7px 0 10px;
  color: #6f747c;
  font-size: 12px;
  line-height: 24px;
  background: #fff;
  border: 1px solid #eceef1;
  border-radius: 4px;
  cursor: pointer;
  transition:
    color 160ms ease,
    border-color 160ms ease,
    background-color 160ms ease;
}
.router-tag:hover {
  color: var(--color-primary);
  border-color: var(--color-primary-border);
}
.router-tag--active {
  color: var(--color-primary);
  background: #fff8f7;
  border-color: var(--color-primary-border);
}
.router-tag__close {
  color: #a4a8ae;
  font-size: 9px;
}
.router-tag__close:hover {
  color: var(--color-primary);
}
</style>
