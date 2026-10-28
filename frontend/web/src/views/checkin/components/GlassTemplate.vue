<script setup lang="ts">
import EventBrand from '@/components/EventBrand.vue'
import MachineAtmosphere from './MachineAtmosphere.vue'
import GlassMachine from './GlassMachine.vue'
import PreviewQr from './PreviewQr.vue'
import { createCheckinTemplateState, type CheckinTemplateProps } from '../useCheckin'
const props = defineProps<CheckinTemplateProps>()
const { event, countdown, attendance } = createCheckinTemplateState(props)
</script>
<template>
  <main class="glass-page">
    <MachineAtmosphere />
    <header>
      <EventBrand />
      <small>LUCKY ATELIER / 幸运制造所</small>
    </header>
    <div class="main">
      <aside>
        <span class="eyebrow">WELCOME ABOARD</span>
        <h1>
          让相聚
          <br />
          多一点
          <span>幸运。</span>
        </h1>
        <p>{{ event.description }}</p>
        <div class="qr">
          <PreviewQr
            :qr="qr"
            :deadline="deadline"
          />
        </div>
      </aside>
      <section class="chamber">
        <div class="chamber-title">
          <span>THE LUCKY MACHINE</span>
          <small>每一次到场，都是一颗幸运星</small>
        </div>
        <GlassMachine :count="event.checkedIn" />
      </section>
    </div>
    <footer>
      <div>
        <small>已到场 / 总人数</small>
        <strong>
          {{ event.checkedIn }}
          <em v-if="event.total !== null">/ {{ event.total }}</em>
        </strong>
      </div>
      <div v-if="event.total !== null" class="progress">
        <small>相聚进度 {{ attendance }}%</small>
        <progress
          :value="attendance"
          max="100"
          aria-label="签到进度"
        />
      </div>
      <div>
        <small>签到截止倒计时</small>
        <strong>{{ countdown.join(' : ') }}</strong>
      </div>
    </footer>
  </main>
</template>
<style scoped>
.glass-page {
  height: 100dvh;
  overflow: auto;
  background: radial-gradient(ellipse at 70% 45%, #273d65, transparent 58%), #0b1429;
  color: #e5edff;
  padding: 30px 5vw 85px;
  display: flex;
  flex-direction: column;
  font-family: 'Sora Variable', 'Microsoft YaHei', sans-serif;
}
header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
header > span {
  font-size: 22px;
}
header small {
  font-size: 10px;
  letter-spacing: 3px;
  color: #95abcf;
}
.main {
  display: grid;
  grid-template-columns: 1fr 1.8fr;
  gap: 5vw;
  flex: 1;
  min-height: 550px;
  padding-top: 35px;
}
aside {
  max-width: 380px;
}
.eyebrow {
  font-size: 10px;
  letter-spacing: 4px;
  color: #a79adb;
}
h1 {
  font-size: clamp(34px, 3.2vw, 60px);
  line-height: 1.4;
  font-weight: 450;
  margin: 20px 0;
}
h1 span {
  color: #baa5ff;
}
aside > p {
  font-size: 12px;
  line-height: 1.9;
  color: #8d9fbd;
}
.qr {
  width: 190px;
  margin-top: 28px;
}
.chamber {
  display: flex;
  flex-direction: column;
  min-width: 0;
}
.chamber-title {
  display: flex;
  justify-content: space-between;
  color: #92acd1;
  font-size: 10px;
  letter-spacing: 2px;
}
.chamber-title small {
  letter-spacing: 0;
}
.chamber :deep(.machine) {
  flex: 1;
  height: auto;
}
footer {
  display: grid;
  grid-template-columns: 1fr 1.2fr 1fr;
  gap: 7vw;
  padding: 25px 0 0;
  border-top: 1px solid #9ab8ff20;
}
footer small {
  display: block;
  color: #8097b8;
  font-size: 11px;
}
footer strong {
  display: block;
  font-size: 30px;
  font-weight: 350;
  margin-top: 10px;
  font-variant-numeric: tabular-nums;
}
footer em {
  font-size: 15px;
  font-style: normal;
  color: #8097b8;
}
progress {
  width: 100%;
  height: 4px;
  accent-color: #b5a1e8;
  margin-top: 28px;
}
@media (max-width: 800px) {
  header small {
    display: none;
  }
  .main {
    grid-template-columns: 1fr;
    gap: 30px;
  }
  aside {
    max-width: none;
  }
  .qr {
    width: 190px;
  }
  .chamber {
    height: 440px;
  }
  .chamber-title {
    font-size: 8px;
  }
  footer {
    grid-template-columns: 1fr 1fr;
    gap: 25px;
  }
  .progress {
    display: none;
  }
  footer strong {
    font-size: 22px;
  }
}

.glass-page {
  background: radial-gradient(ellipse at 70% 45%, #2a3159, transparent 60%), #101427;
}
.main {
  grid-template-columns: 0.85fr 1.8fr;
  gap: 4vw;
}
h1 {
  font-family: 'STXingkai', '华文行楷', 'KaiTi', '楷体', serif;
  font-size: clamp(40px, 4.3vw, 78px);
  font-weight: 400;
  line-height: 1.25;
  letter-spacing: 0.03em;
  color: #dce6f7;
}
h1 span {
  display: block;
  margin-left: 0.6em;
  background: linear-gradient(125deg, #d9c5ff, #a2baff 55%, #e1d6fc);
  background-clip: text;
  -webkit-text-fill-color: transparent;
}
.eyebrow {
  font-size: 9px;
  letter-spacing: 5px;
  color: #968bb8;
}
aside > p {
  font-family: 'Songti SC', 'SimSun', serif;
  font-size: 13px;
  color: #8d9abb;
  letter-spacing: 1px;
}
.qr {
  width: clamp(180px, 12vw, 225px);
  padding: 10px;
  border: 1px solid #b5b5ed35;
  border-radius: 22px;
  background: #c4c6ff05;
}
.qr :deep(.placeholder) {
  border-radius: 13px;
  background: #eef0f7;
  color: #5b577b;
}
.qr :deep(.scan-block p) {
  font-size: 9px;
  margin: 13px 0 3px;
}
.chamber-title {
  padding-bottom: 15px;
  border-bottom: 1px solid #b4b8ec15;
  letter-spacing: 3px;
}
.chamber-title small {
  color: #98a6c2;
}
footer strong {
  font-family: Georgia, serif;
  font-size: 36px;
  color: #c9c7eb;
}
footer small {
  letter-spacing: 2px;
}
@media (max-width: 800px) {
  .main {
    grid-template-columns: 1fr;
  }
  h1 {
    font-size: 48px;
  }
  .qr {
    width: 190px;
  }
  footer strong {
    font-size: 26px;
  }
}

.glass-page {
  position: relative;
  isolation: isolate;
  background: radial-gradient(ellipse at 70% 58%, #273252, #080e20 70%);
}
.chamber {
  position: relative;
}
.chamber::after {
  content: '';
  position: absolute;
  bottom: 3%;
  left: 10%;
  right: 10%;
  height: 8%;
  border-radius: 50%;
  border: 1px solid #99b8ff38;
  background: radial-gradient(ellipse, #8ca7ed20, transparent 70%);
  box-shadow: 0 0 40px #788dd815;
  pointer-events: none;
}
</style>
