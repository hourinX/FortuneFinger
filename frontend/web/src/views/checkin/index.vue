<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, reactive, ref, watch } from 'vue'
import {
  createCheckinStateSocket,
  getCheckinQRCode,
  getCheckinState,
  getCheckinTemplate,
  type CheckinStateMessage,
  type CheckinStateResp,
} from '@/api/checkin'
import LoadingScreen from './components/LoadingScreen.vue'
import { defaultTemplateId, resolveTemplate } from './templates'
import { useCheckin, type CheckinContent, type CheckinSource } from './useCheckin'

const loadingCycleDuration = 1600
const requestRetryDelay = 1500
const source = reactive<CheckinSource>({
  qrCode: '',
  status: 'pending',
  checkedIn: 0,
  total: 0,
  deadline: Date.now(),
  serverTimeOffset: 0,
})
const content: CheckinContent = {
  company: 'FortuneFinger',
  title: '此刻相聚，共赴新程。',
  description: '欢迎来到我们的年度聚会。一起回顾值得珍藏的瞬间，迎接下一段精彩。',
}
const checkin = useCheckin(source)
const loading = ref(true)
const configuredTemplate = ref(defaultTemplateId)
const activeTemplate = computed(() => resolveTemplate(configuredTemplate.value))
let mounted = true
let initialReady = false
let cycleTimer: number | undefined
let retryTimer: number | undefined
let socketRetryTimer: number | undefined
let requestController: AbortController | undefined
let stateSocket: WebSocket | undefined
let qrRetryTimer: number | undefined

function waitBeforeRetry() {
  return new Promise<void>((resolve) => {
    retryTimer = window.setTimeout(resolve, requestRetryDelay)
  })
}

function applyCheckinState(state: CheckinStateResp) {
  const wasActive = source.status === 'active'
  source.status = state.status
  source.serverTimeOffset = Date.parse(state.server_time) - Date.now()
  source.deadline = Date.parse(state.checkin_end_time)
  checkin.setAttendance(state.checked_in, state.total)
  if (state.status !== 'active') {
    checkin.setQRCode('')
  } else if (!wasActive && initialReady) {
    void refreshQRCode()
  }
}

async function refreshQRCode(signal?: AbortSignal) {
  window.clearTimeout(qrRetryTimer)
  if (source.status !== 'active') {
    checkin.setQRCode('')
    return
  }
  try {
    const qr = await getCheckinQRCode(signal)
    if (!mounted || source.status !== 'active') return
    checkin.setQRCode(qr.qrcode, qr.expires_at)
  } catch {
    if (!mounted || source.status !== 'active') return
    checkin.setQRCode('')
    qrRetryTimer = window.setTimeout(() => void refreshQRCode(), 2000)
  }
}

async function loadInitialData() {
  while (mounted && !initialReady) {
    requestController = new AbortController()
    try {
      const [template, state] = await Promise.all([
        getCheckinTemplate(requestController.signal),
        getCheckinState(requestController.signal),
      ])
      if (!mounted) return
      configuredTemplate.value = template.template
      applyCheckinState(state)
      if (state.status === 'active') {
        await refreshQRCode(requestController.signal)
      }
      if (!mounted) return
      initialReady = true
      connectStateSocket()
    } catch {
      if (!mounted) return
      await waitBeforeRetry()
    }
  }
}

async function refreshStateBeforeReconnect() {
  try {
    const state = await getCheckinState()
    if (mounted) applyCheckinState(state)
  } catch {
    // WebSocket 重连会继续尝试，当前页面保留最后一次有效快照。
  }
}

function scheduleSocketReconnect() {
  window.clearTimeout(socketRetryTimer)
  socketRetryTimer = window.setTimeout(async () => {
    if (!mounted) return
    await refreshStateBeforeReconnect()
    if (mounted) connectStateSocket()
  }, 2000)
}

function connectStateSocket() {
  stateSocket?.close()
  const socket = createCheckinStateSocket()
  stateSocket = socket
  socket.addEventListener('message', (event) => {
    try {
      const message = JSON.parse(event.data) as CheckinStateMessage
      if (message.key === 'checkin.state' && message.data) {
        applyCheckinState(message.data)
      }
    } catch {
      // 忽略无法识别的非业务消息，保持当前有效状态。
    }
  })
  socket.addEventListener('close', () => {
    if (stateSocket === socket) stateSocket = undefined
    if (mounted) scheduleSocketReconnect()
  })
}

function playLoadingCycle() {
  cycleTimer = window.setTimeout(() => {
    if (!mounted) return
    if (initialReady) {
      loading.value = false
      return
    }
    playLoadingCycle()
  }, loadingCycleDuration)
}

onMounted(() => {
  void loadInitialData()
  playLoadingCycle()
})

watch(
  () => checkin.qr.refreshVersion,
  () => void refreshQRCode(),
)

onBeforeUnmount(() => {
  mounted = false
  requestController?.abort()
  stateSocket?.close()
  window.clearTimeout(cycleTimer)
  window.clearTimeout(retryTimer)
  window.clearTimeout(socketRetryTimer)
  window.clearTimeout(qrRetryTimer)
})
</script>

<template>
  <div class="checkin-entry">
    <Transition name="loading-fade">
      <LoadingScreen
        v-if="loading"
      />
    </Transition>
    <component
      v-if="!loading"
      :is="activeTemplate.component"
      :content="content"
      :qr="checkin.qr"
      :attendance="checkin.attendance"
      :deadline="checkin.deadline"
    />
  </div>
</template>

<style scoped>
.checkin-entry {
  width: 100%;
  height: 100%;
}
.loading-fade-leave-active {
  transition: opacity 0.38s ease;
}
.loading-fade-leave-to {
  opacity: 0;
}
</style>
