<script setup lang="ts">
import BrandHeader from '../BrandHeader.vue'

defineProps<{ company?: string }>()
</script>

<template>
  <div class="expired-page">
    <BrandHeader :company="company || 'FortuneFinger'" />
    <section>
      <span class="expired-page__icon">×</span>
      <p>CHECK-IN LINK EXPIRED</p>
      <h1>暂时无法签到</h1>
      <small>二维码已过期</small>
    </section>
  </div>
</template>

<style scoped>
.expired-page { width: min(100% - 32px, 440px); margin: 0 auto; padding: max(24px, env(safe-area-inset-top)) 0 max(24px, env(safe-area-inset-bottom)); }
.expired-page section { display: flex; min-height: 390px; margin-top: 28px; padding: 36px 28px; box-sizing: border-box; flex-direction: column; align-items: center; justify-content: center; text-align: center; background: rgb(255 255 255 / 90%); border: 1px solid rgb(255 255 255 / 75%); border-radius: 22px; box-shadow: 0 24px 70px rgb(81 61 45 / 11%); }
.expired-page__icon { display: grid; width: 64px; height: 64px; margin-bottom: 23px; place-items: center; color: #fff; font-size: 27px; font-weight: 300; background: linear-gradient(145deg, #ef4a42, #d73331); border-radius: 50%; box-shadow: 0 15px 34px rgb(220 54 49 / 22%); }
.expired-page section p { margin: 0 0 9px; color: #d9423b; font-size: 8px; font-weight: 750; letter-spacing: .19em; }
.expired-page h1 { margin: 0; color: #302c29; font-size: 25px; }
.expired-page small { margin-top: 10px; color: #9b928b; font-size: 12px; }
</style>
