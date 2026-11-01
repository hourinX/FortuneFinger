<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useParticipantCheckinStore } from '@/stores/checkin'
import CheckinPage from './components/checkin/index.vue'
import ExpiredPage from './components/expired/index.vue'

const route = useRoute()
const store = useParticipantCheckinStore()
const { session, companies, loading, error } = storeToRefs(store)
const token = computed(() => String(route.query.token || ''))
const forceExpired = computed(() => route.query.expired === '1')
const preview = computed(() => (
  route.query.preview === '1' || (import.meta.env.DEV && !token.value && !forceExpired.value)
))

onMounted(() => {
  if (forceExpired.value) return
  if (!token.value && !preview.value) return
  void store.load(token.value, preview.value)
})
</script>

<template>
  <main class="checkin-entry">
    <div class="checkin-entry__glow checkin-entry__glow--one"></div>
    <div class="checkin-entry__glow checkin-entry__glow--two"></div>

    <section v-if="loading" class="checkin-loading">
      <span></span>
      <h1>正在进入签到现场</h1>
      <p>正在验证本次签到二维码</p>
    </section>
    <CheckinPage v-else-if="session" :session="session" :companies="companies" />
    <ExpiredPage v-else-if="forceExpired || error || (!token && !preview)" />
  </main>
</template>

<style scoped>
.checkin-entry { position: relative; min-height: 100dvh; overflow: hidden; background: #f6f2eb; }
.checkin-entry__glow { position: fixed; border-radius: 50%; pointer-events: none; }
.checkin-entry__glow--one { top: -180px; right: -170px; width: 390px; height: 390px; background: radial-gradient(circle, rgb(237 65 57 / 15%), transparent 68%); }
.checkin-entry__glow--two { bottom: -210px; left: -180px; width: 430px; height: 430px; background: radial-gradient(circle, rgb(211 164 91 / 14%), transparent 68%); }
.checkin-loading { position: relative; z-index: 1; display: flex; min-height: 100dvh; flex-direction: column; align-items: center; justify-content: center; text-align: center; }
.checkin-loading span { width: 38px; height: 38px; border: 3px solid #eadfd7; border-top-color: #e4433d; border-radius: 50%; animation: spin .8s linear infinite; }
.checkin-loading h1 { margin: 22px 0 7px; color: #302c29; font-size: 20px; }
.checkin-loading p { margin: 0; color: #9b928b; font-size: 12px; }
@keyframes spin { to { transform: rotate(360deg); } }
</style>
