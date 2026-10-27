<script setup lang="ts">
import CosmicBackground from './CosmicBackground.vue'
import EventBrand from '@/components/EventBrand.vue'
import PreviewQr from './PreviewQr.vue'
import { createCheckinTemplateState, type CheckinTemplateProps } from '../useCheckin'
const props = defineProps<CheckinTemplateProps>()
const { event, attendance, countdown, arrival } = createCheckinTemplateState(props)
</script>
<template>
  <main class="orbit-page">
    <CosmicBackground />
    <div
      class="grid"
      aria-hidden="true"
    ></div>
    <header>
      <EventBrand />
      <span>CONNECTED / ANNUAL EVENT</span>
    </header>
    <div class="headline">
      <small>每一个你，点亮我们的星系</small>
      <h1>此刻，同频共振。</h1>
      <p>{{ event.description }}</p>
    </div>
    <section class="command">
      <aside>
        <span class="label">01 / ARRIVALS</span>
        <h2>相聚坐标</h2>
        <strong>
          {{ event.checkedIn }}
          <small>/ {{ event.total }}</small>
        </strong>
        <div
          class="segments"
          role="progressbar"
          aria-label="签到进度"
          :aria-valuenow="attendance"
          aria-valuemin="0"
          aria-valuemax="100"
        >
          <i
            v-for="n in 20"
            :key="n"
            :class="{ lit: n <= attendance / 5 }"
          ></i>
        </div>
        <p>{{ attendance }}% 的伙伴已抵达</p>
      </aside>
      <div class="portal">
        <div
          class="orbit orbit-a"
          aria-hidden="true"
        ></div>
        <div
          class="orbit orbit-b"
          aria-hidden="true"
        ></div>
        <div
          class="orbit orbit-c"
          aria-hidden="true"
        ></div>
        <span class="portal-label">CHECK-IN PORTAL</span>
        <div class="scan">
          <PreviewQr
            :qr="qr"
            :deadline="deadline"
          />
        </div>
        <div
          :key="arrival"
          class="arrival"
          role="status"
        >
          {{ arrival ? '新的连接已建立 · 欢迎到场' : '等待与你连接' }}
        </div>
      </div>
      <aside class="time">
        <span class="label">02 / TIME WINDOW</span>
        <h2>入场时间窗</h2>
        <div class="clock">
          <div
            v-for="(part, i) in countdown"
            :key="i"
          >
            <strong>{{ part }}</strong>
            <small>{{ ['HOURS', 'MINUTES', 'SECONDS'][i] }}</small>
          </div>
        </div>
        <p>请在时间窗关闭前完成签到</p>
      </aside>
    </section>
    <footer>
      <span>TOGETHER WE GO FURTHER</span>
      <div>
        <i></i>
        签到交互预览
      </div>
      <span>FF / CONNECTION HUB</span>
    </footer>
  </main>
</template>
<style scoped>
.orbit-page {
  height: 100dvh;
  overflow: auto;
  position: relative;
  isolation: isolate;
  background: radial-gradient(ellipse at 50% 55%, #13384c, transparent 60%), #060f1c;
  color: #d9f4fa;
  padding: 30px 5vw 85px;
  display: flex;
  flex-direction: column;
  font-family: 'Sora Variable', 'Microsoft YaHei', sans-serif;
}
.grid {
  position: absolute;
  inset: 0;
  z-index: -1;
  pointer-events: none;
  background-image:
    linear-gradient(#63b9d305 1px, transparent 1px),
    linear-gradient(90deg, #63b9d305 1px, transparent 1px);
  background-size: 60px 60px;
}
header,
footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 20px;
}
header b {
  font-size: 20px;
}
header > span {
  font-size: 9px;
  letter-spacing: 3px;
  color: #598b9c;
}
.headline {
  text-align: center;
  margin: 35px 0 10px;
}
.headline > small {
  font-size: 11px;
  letter-spacing: 4px;
  color: #68cad0;
}
h1 {
  font-size: clamp(34px, 3.5vw, 64px);
  font-weight: 350;
  letter-spacing: 4px;
  margin: 20px 0;
}
.headline p {
  font-size: 11px;
  color: #7b9da8;
}
.command {
  display: grid;
  grid-template-columns: 1fr 1.8fr 1fr;
  align-items: center;
  gap: 25px;
  flex: 1;
  min-height: 480px;
}
.label {
  font-size: 9px;
  letter-spacing: 3px;
  color: #508495;
}
h2 {
  font-size: 15px;
  font-weight: 400;
  margin: 20px 0;
}
aside > strong {
  font-size: 60px;
  font-weight: 300;
}
aside > strong small {
  font-size: 14px;
  color: #5d8f9e;
}
aside p {
  font-size: 10px;
  color: #648995;
  margin-top: 20px;
}
.segments {
  display: flex;
  gap: 4px;
  height: 15px;
  margin-top: 22px;
}
.segments i {
  flex: 1;
  background: #153241;
}
.segments i.lit {
  background: #74d6d0;
  box-shadow: 0 0 10px #66ccc822;
}
.portal {
  position: relative;
  min-height: 450px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
}
.orbit {
  position: absolute;
  width: 100%;
  aspect-ratio: 1;
  border: 1px solid #5fbac52e;
  border-radius: 50%;
  pointer-events: none;
  max-width: 470px;
}
.orbit-a {
  border-top: 2px solid #7dd7d5;
  border-bottom: 2px solid #7dd7d5;
  animation: rotate 55s linear infinite;
  box-shadow: 0 0 35px #4bc7da0b;
}
.orbit-b {
  width: 88%;
  border-style: dashed;
  animation: rotate 90s linear infinite reverse;
}
.orbit-c {
  width: 110%;
  border-color: #4b899312;
}
.portal-label {
  font-size: 9px;
  letter-spacing: 4px;
  color: #71b6bf;
  margin-bottom: 22px;
}
.scan {
  width: 185px;
  z-index: 1;
}
.scan :deep(.placeholder) {
  border-radius: 2px;
  background: #daedf0;
}
.arrival {
  font-size: 10px;
  color: #69c7c4;
  margin-top: 25px;
  animation: appear 1s ease;
}
.clock {
  display: flex;
  gap: 17px;
}
.clock > div > strong {
  font-size: 32px;
  font-weight: 300;
}
.clock small {
  display: block;
  font-size: 7px;
  letter-spacing: 1px;
  color: #5c8997;
  margin-top: 8px;
}
.time {
  padding-left: 25px;
}
footer {
  border-top: 1px solid #518b9725;
  padding-top: 22px;
  font-size: 8px;
  letter-spacing: 2px;
  color: #487280;
}
footer i {
  display: inline-block;
  width: 5px;
  height: 5px;
  background: #65d8c0;
  border-radius: 50%;
  margin-right: 8px;
}
@keyframes rotate {
  to {
    transform: rotate(360deg);
  }
}
@keyframes appear {
  from {
    opacity: 0;
    transform: translateY(8px);
  }
  to {
    opacity: 1;
    transform: none;
  }
}
@media (max-width: 850px) {
  header > span {
    display: none;
  }
  .headline p {
    line-height: 1.8;
  }
  .command {
    grid-template-columns: 1fr 1fr;
  }
  .portal {
    grid-row: 1;
    grid-column: 1/-1;
    min-height: 420px;
  }
  .orbit {
    max-width: 350px;
  }
  .time {
    padding: 0;
  }
  .clock {
    gap: 8px;
  }
  .clock > div > strong {
    font-size: 23px;
  }
  aside > strong {
    font-size: 45px;
  }
  .label {
    font-size: 7px;
    letter-spacing: 1px;
  }
  footer {
    margin-top: 30px;
  }
  footer > span:last-child {
    display: none;
  }
}
.orbit-page {
  background: #060a17;
}
.grid {
  display: none;
}
.headline {
  margin-top: 30px;
}
.headline h1 {
  font-family: 'STXingkai', '华文行楷', 'KaiTi', '楷体', serif;
  font-size: clamp(42px, 5vw, 90px);
  font-weight: 400;
  letter-spacing: 0.1em;
  background: linear-gradient(130deg, #f1edff, #97acd9 42%, #c2b1ea 70%, #d8f0fa);
  background-clip: text;
  -webkit-text-fill-color: transparent;
  filter: drop-shadow(0 0 16px #6a7be425);
}
.headline > small {
  color: #9a96c4;
  letter-spacing: 5px;
  font-size: 10px;
}
.headline p {
  font-family: 'Songti SC', 'SimSun', serif;
  letter-spacing: 2px;
  color: #8993ad;
}
.portal {
  perspective: 900px;
}
.orbit-a {
  border-color: #7a8fcf25;
  border-top-color: #b9b8f0;
  border-bottom-color: #7eb5d6;
  box-shadow: 0 0 35px #7970be12;
}
.orbit-b {
  border-color: #a59ac82a;
  border-style: solid;
  transform: rotateX(60deg) rotateZ(-20deg);
  animation: galactic-ring 45s linear infinite;
}
.orbit-c {
  border-color: #9ba6da15;
  transform: rotateY(50deg);
  animation: galactic-outer 60s linear infinite;
}
.scan {
  width: 205px;
  padding: 10px;
  border: 1px solid #b0bce745;
  background: #6b6d9a0d;
  box-shadow: 0 0 35px #9197d012;
}
.scan :deep(.placeholder) {
  background: #e8ecf5;
  color: #4c5575;
  border-radius: 0;
}
.scan :deep(.scan-block p) {
  font-size: 9px;
}
.portal-label {
  color: #a6abd2;
  letter-spacing: 5px;
}
.arrival {
  color: #a9b2da;
}
.label {
  color: #8a8eb4;
}
h2 {
  font-family: 'Songti SC', 'SimSun', serif;
  font-size: 19px;
  letter-spacing: 4px;
  color: #c8cee4;
}
aside > strong {
  font-family: Georgia, serif;
  color: #c6c6e4;
}
.clock > div > strong {
  font-family: Georgia, serif;
  color: #c6c6e4;
}
.segments i.lit {
  background: #a2aad5;
  box-shadow: 0 0 8px #a2aad522;
}
.segments i {
  background: #232b45;
}
.command aside {
  border-top: 1px solid #a4a9cb25;
  padding-top: 22px;
}
.time {
  padding-left: 0;
}
@keyframes galactic-ring {
  to {
    transform: rotateX(60deg) rotateZ(340deg);
  }
}
@keyframes galactic-outer {
  to {
    transform: rotateY(50deg) rotateZ(360deg);
  }
}
@media (max-width: 850px) {
  .headline h1 {
    font-size: clamp(32px, 7vw, 54px);
  }
  .scan {
    width: 180px;
  }
  .portal {
    min-height: 430px;
  }
  h2 {
    font-size: 15px;
    letter-spacing: 2px;
  }
}

@media (prefers-reduced-motion: reduce) {
  * {
    animation: none !important;
  }
}
</style>
