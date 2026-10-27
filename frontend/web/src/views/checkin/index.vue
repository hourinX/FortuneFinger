<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, reactive, ref } from 'vue'
import { getCheckinTemplate } from '@/api/checkin'
import LoadingScreen from './components/LoadingScreen.vue'
import { defaultTemplateId, resolveTemplate } from './templates'
import { useCheckin, type CheckinContent, type CheckinSource } from './useCheckin'

const loadingCycleDuration = 1600
const requestRetryDelay = 1500
const source = reactive<CheckinSource>({
  qrCode: '',
  checkedIn: 128,
  total: 200,
  deadline: Date.now() + 90 * 60 * 1000,
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
let templateReady = false
let cycleTimer: number | undefined
let retryTimer: number | undefined
let requestController: AbortController | undefined

function waitBeforeRetry() {
  return new Promise<void>((resolve) => {
    retryTimer = window.setTimeout(resolve, requestRetryDelay)
  })
}

async function loadTemplate() {
  while (mounted && !templateReady) {
    requestController = new AbortController()
    try {
      const resp = await getCheckinTemplate(requestController.signal)
      if (!mounted) return
      configuredTemplate.value = resp.template
      templateReady = true
    } catch {
      if (!mounted) return
      await waitBeforeRetry()
    }
  }
}

function playLoadingCycle() {
  cycleTimer = window.setTimeout(() => {
    if (!mounted) return
    if (templateReady) {
      loading.value = false
      return
    }
    playLoadingCycle()
  }, loadingCycleDuration)
}

onMounted(() => {
  void loadTemplate()
  playLoadingCycle()
})

onBeforeUnmount(() => {
  mounted = false
  requestController?.abort()
  window.clearTimeout(cycleTimer)
  window.clearTimeout(retryTimer)
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
