<template>
  <main class="screen-shell">
    <StageBackground :color="accent" />
    <div class="screen-shell__glow"></div>
    <section class="screen-shell__content">
      <header class="screen-shell__header">
        <div class="screen-shell__brand">
          <img :src="brandIcon || brandMark" alt="" />
          <div>
            <strong>{{ brandName }}</strong>
            <span>ANNUAL EVENT EXPERIENCE</span>
          </div>
        </div>
        <div class="screen-shell__status">
          <span :class="{ 'screen-shell__dot--offline': !screen.connected }" class="screen-shell__dot"></span>
          {{ screen.connected ? '服务已连接' : '连接已断开' }}
        </div>
      </header>
      <slot />
    </section>
  </main>
</template>

<script setup lang="ts">
import { useScreenStore } from '@/stores/screen'
import brandMark from '@/assets/branding/fortune-finger.png'
import StageBackground from './StageBackground.vue'

withDefaults(
  defineProps<{
    accent?: string
    brandIcon?: string
    brandName?: string
  }>(),
  {
    accent: '#ef4444',
    brandIcon: '',
    brandName: 'FortuneFinger',
  },
)

const screen = useScreenStore()
</script>

<style scoped>
.screen-shell {
  position: relative;
  width: 100vw;
  height: 100vh;
  overflow: hidden;
  background:
    radial-gradient(ellipse 64% 68% at 57% 47%, rgba(139, 18, 18, 0.68) 0%, rgba(91, 8, 13, 0.76) 46%, rgba(43, 4, 9, 0.92) 74%, rgba(15, 2, 5, 0.98) 100%),
    radial-gradient(circle at 74% 4%, rgba(226, 83, 38, 0.13), transparent 31%),
    #1c0307;
}

.screen-shell::before {
  position: absolute;
  top: -32%;
  right: 7%;
  width: 46%;
  height: 142%;
  background: linear-gradient(112deg, transparent 20%, rgba(255, 181, 136, 0.05) 49%, transparent 72%);
  content: '';
  filter: blur(10px);
  transform: rotate(8deg);
}

.screen-shell::after {
  position: absolute;
  right: -8%;
  bottom: -23%;
  width: 55%;
  height: 46%;
  border: 2px solid rgba(255, 184, 113, 0.16);
  border-radius: 50%;
  box-shadow: 0 0 36px rgba(255, 68, 45, 0.16);
  content: '';
  transform: rotate(-8deg);
}

.screen-shell__glow {
  position: absolute;
  top: -20vh;
  left: 50%;
  width: 54vw;
  height: 54vw;
  background: rgba(233, 75, 49, 0.045);
  border-radius: 50%;
  filter: blur(80px);
  transform: translateX(-50%);
}

.screen-shell__content {
  position: relative;
  z-index: 1;
  display: flex;
  height: 100%;
  padding: clamp(20px, 1.75vw, 34px) clamp(28px, 2.5vw, 48px) clamp(26px, 2.4vw, 46px);
  flex-direction: column;
}

.screen-shell__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.screen-shell__brand {
  display: flex;
  align-items: center;
  gap: 13px;

  > img {
    width: clamp(27px, 1.8vw, 36px);
    height: clamp(27px, 1.8vw, 36px);
    padding: 5px;
    background: rgba(255, 255, 255, 0.94);
    border-radius: 9px;
    object-fit: contain;
  }

  > div {
    display: flex;
    flex-direction: column;
    gap: 4px;
  }

  strong {
    font-family: 'Sora Variable', sans-serif;
    font-size: clamp(15px, 1.05vw, 21px);
    font-weight: 680;
    letter-spacing: -0.045em;
  }

  span {
    color: rgba(255, 255, 255, 0.42);
    font-size: clamp(6px, 0.45vw, 9px);
    letter-spacing: 0.28em;
  }
}

.screen-shell__status {
  display: flex;
  align-items: center;
  gap: 8px;
  color: rgba(255, 255, 255, 0.54);
  font-size: clamp(10px, 0.68vw, 13px);
}

.screen-shell__dot {
  width: 7px;
  height: 7px;
  background: #43c784;
  border-radius: 50%;
  box-shadow: 0 0 12px rgba(67, 199, 132, 0.72);
}

.screen-shell__dot--offline {
  background: #ef4444;
  box-shadow: 0 0 12px rgba(239, 68, 68, 0.72);
}
</style>
