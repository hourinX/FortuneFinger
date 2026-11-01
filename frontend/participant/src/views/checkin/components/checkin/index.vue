<script setup lang="ts">
import { useRouter } from 'vue-router'
import type { CheckinSessionResp } from '@/api/checkin'
import type { OrgCompanyOption } from '@/api/org'
import type { CheckinFormValue } from '@/stores/checkin'
import BrandHeader from '../BrandHeader.vue'
import CheckinForm from '../CheckinForm.vue'

defineProps<{
  session: CheckinSessionResp
  companies: OrgCompanyOption[]
}>()

const router = useRouter()

function enterChat(_value: CheckinFormValue) {
  void router.push({ name: 'participant-home', query: { preview: '1' } })
}
</script>

<template>
  <div class="checkin-content">
    <BrandHeader :company="session.activity.company" />
    <section class="checkin-card">
      <div class="checkin-card__status">
        <p>完成签到后即可参与今晚的抽奖活动</p>
        <span><i>✓</i>现场二维码验证成功</span>
      </div>
      <CheckinForm :session="session" :companies="companies" @submit="enterChat" />
    </section>
    <footer>© 2026 FortuneFinger · 让每一次相聚都有回应</footer>
  </div>
</template>

<style scoped>
.checkin-content { width: min(100% - 32px, 440px); margin: 0 auto; padding: max(24px, env(safe-area-inset-top)) 0 max(24px, env(safe-area-inset-bottom)); }
.checkin-card { margin-top: 28px; overflow: hidden; background: rgb(255 255 255 / 90%); border: 1px solid rgb(255 255 255 / 75%); border-radius: 22px; box-shadow: 0 24px 70px rgb(81 61 45 / 11%); }
.checkin-card__status { display: flex; min-height: 58px; padding: 0 22px; align-items: center; justify-content: space-between; gap: 14px; background: linear-gradient(155deg, #fff, #fbf7f1); border-bottom: 1px solid #f0ebe4; }
.checkin-card__status p { margin: 0; color: #5e5752; font-size: 11px; line-height: 17px; }
.checkin-card__status span { display: flex; flex: 0 0 auto; align-items: center; gap: 5px; color: #568465; font-size: 9px; white-space: nowrap; }
.checkin-card__status i { display: grid; width: 14px; height: 14px; place-items: center; color: #fff; font-size: 8px; font-style: normal; background: #62a276; border-radius: 50%; }
.checkin-content footer { padding: 24px 0 4px; color: #b1a8a0; font-size: 9px; text-align: center; }
</style>
