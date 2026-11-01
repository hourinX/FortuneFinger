<script setup lang="ts">
import AppTabBar from '@/components/app-tabbar/index.vue'
import BrandHeader from '@/views/checkin/components/BrandHeader.vue'

defineProps<{
  active: 'lottery' | 'profile'
  eyebrow: string
  title: string
  description: string
}>()
</script>

<template>
  <main class="placeholder-page">
    <header>
      <BrandHeader company="FortuneFinger" caption="PARTICIPANT CENTER" />
    </header>
    <section>
      <div class="placeholder-page__symbol">
        <i></i>
      </div>
      <p>{{ eyebrow }}</p>
      <h1>{{ title }}</h1>
      <span>{{ description }}</span>
      <small>功能正在准备中</small>
    </section>
    <AppTabBar :active="active" />
  </main>
</template>

<style scoped>
.placeholder-page { display: grid; width: min(100%, 440px); height: 100dvh; margin: 0 auto; grid-template-rows: auto minmax(0, 1fr) auto; overflow: hidden; background: linear-gradient(180deg, #f8f4ed, #f2ece4); }
.placeholder-page > header { padding: max(24px, env(safe-area-inset-top)) 18px 12px; }
.placeholder-page > section { display: flex; margin: 18px 16px 24px; padding: 32px 28px; flex-direction: column; align-items: center; justify-content: center; text-align: center; background: rgb(255 255 255 / 86%); border: 1px solid rgb(255 255 255 / 74%); border-radius: 22px; box-shadow: 0 24px 70px rgb(81 61 45 / 9%); }
.placeholder-page__symbol { position: relative; width: 70px; height: 70px; margin-bottom: 22px; background: #fff1ee; border-radius: 22px; transform: rotate(45deg); }
.placeholder-page__symbol::before, .placeholder-page__symbol::after, .placeholder-page__symbol i { position: absolute; background: #df443d; border-radius: 50%; content: ''; }
.placeholder-page__symbol::before { top: 18px; left: 18px; width: 9px; height: 9px; }
.placeholder-page__symbol::after { right: 18px; bottom: 18px; width: 13px; height: 13px; opacity: .55; }
.placeholder-page__symbol i { top: 31px; left: 31px; width: 7px; height: 7px; opacity: .75; }
.placeholder-page section > p { margin: 0 0 8px; color: #dc433c; font-size: 8px; font-weight: 750; letter-spacing: .2em; }
.placeholder-page h1 { margin: 0; color: #302c29; font-size: 25px; }
.placeholder-page section > span { max-width: 270px; margin-top: 11px; color: #8f867f; font-size: 12px; line-height: 20px; }
.placeholder-page section > small { margin-top: 24px; padding: 6px 11px; color: #a09992; font-size: 9px; background: #f5f1eb; border-radius: 12px; }
</style>
