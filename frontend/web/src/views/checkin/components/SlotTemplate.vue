<script setup lang="ts">
import EventBrand from '@/components/EventBrand.vue'
import { onUnmounted, ref, watch } from 'vue'
import PreviewQr from './PreviewQr.vue'
import { createCheckinTemplateState, type CheckinTemplateProps } from '../useCheckin'
const props = defineProps<CheckinTemplateProps>()
const { event, arrival, countdown } = createCheckinTemplateState(props)
const rolling = ref(false)
let timer: ReturnType<typeof setTimeout> | undefined
watch(arrival, () => {
  rolling.value = true
  clearTimeout(timer)
  timer = setTimeout(() => {
    rolling.value = false
  }, 1800)
})
onUnmounted(() => clearTimeout(timer))
</script>
<template>
  <main class="slot-page">
    <header>
      <EventBrand />
      <span>THE GOLDEN NIGHT</span>
    </header>
    <section class="title">
      <small>相聚即是好运</small>
      <h1>今夜，好运登场</h1>
      <p>{{ event.description }}</p>
    </section>
    <section class="cabinet">
      <div class="marquee">
        <span
          v-for="n in 14"
          :key="n"
        ></span>
      </div>
      <div class="cabinet-label">WELCOME · 幸运签到站</div>
      <div
        class="reels"
        :class="{ rolling }"
      >
        <div
          v-for="(symbol, index) in ['✦', '♛', '✦']"
          :key="index"
          class="reel"
        >
          <div
            class="strip"
            :style="{ '--delay': `${index * 0.12}s` }"
          >
            <span>♛</span>
            <span>✦</span>
            <span>{{ symbol }}</span>
            <span>♛</span>
            <span>✦</span>
          </div>
        </div>
      </div>
      <div
        class="result"
        role="status"
      >
        {{
          rolling ? '好运正在汇聚…' : arrival ? '欢迎新伙伴加入现场' : '每一位到场，都是今晚的主角'
        }}
      </div>
      <div class="console">
        <div class="stats">
          <small>已到场</small>
          <strong>
            {{ event.checkedIn }}
            <em>/ {{ event.total }}</em>
          </strong>
          <span>位伙伴已加入</span>
        </div>
        <div class="scan">
          <PreviewQr
            :qr="qr"
            :deadline="deadline"
          />
        </div>
        <div class="time">
          <small>距离签到结束</small>
          <strong>{{ countdown.join(':') }}</strong>
          <span>请在截止前完成签到</span>
        </div>
      </div>
      <div class="base">FORTUNE FAVORS TOGETHERNESS</div>
    </section>
    <footer>签到互动演示 · 转轮仅用于欢迎反馈，不产生中奖结果</footer>
  </main>
</template>
<style scoped>
.slot-page {
  height: 100dvh;
  overflow: auto;
  padding: 30px 5vw 85px;
  color: #f4dba7;
  background: radial-gradient(ellipse at 50% 50%, #40304f, transparent 65%), #130e20;
  font-family: 'Sora Variable', 'Microsoft YaHei', sans-serif;
}
header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
header b {
  font-size: 20px;
}
header > span {
  font-size: 10px;
  letter-spacing: 4px;
  color: #8e7e9a;
}
.title {
  text-align: center;
  margin: 30px 0;
}
.title small {
  letter-spacing: 7px;
  font-size: 11px;
}
h1 {
  font-family: 'SimSun', serif;
  font-size: clamp(32px, 3.5vw, 60px);
  font-weight: 400;
  margin: 15px 0;
}
.title p {
  font-size: 11px;
  color: #a897b3;
}
.cabinet {
  max-width: 850px;
  margin: auto;
  border: 1px solid #d3ab65;
  border-radius: 35px 35px 14px 14px;
  padding: 15px 35px 0;
  background: linear-gradient(120deg, #4a314c, #21182e 55%, #503b48);
  box-shadow:
    0 0 0 7px #ae814c18,
    0 25px 80px #0006;
  position: relative;
}
.marquee {
  display: flex;
  justify-content: space-between;
  padding: 2px 6px;
}
.marquee > span {
  width: 5px;
  height: 5px;
  background: #ffdc96;
  border-radius: 50%;
  box-shadow: 0 0 9px #ffce76;
  animation: bulb 3s ease-in-out infinite alternate;
}
.marquee > span:nth-child(even) {
  animation-delay: -1.5s;
}
.cabinet-label {
  text-align: center;
  letter-spacing: 4px;
  font-size: 11px;
  padding: 22px;
}
.reels {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
  background: #100e17;
  padding: 12px;
  border: 1px solid #ad8b57;
  border-radius: 12px;
}
.reel {
  height: 150px;
  overflow: hidden;
  background: linear-gradient(#a99270, #f9e5bd 30%, #fff0cf 60%, #b69b77);
  border-radius: 7px;
  position: relative;
}
.reel::after {
  content: '';
  position: absolute;
  inset: 0;
  box-shadow:
    inset 0 15px 20px #35202050,
    inset 0 -15px 20px #35202050;
  pointer-events: none;
}
.strip {
  transform: translateY(-300px);
}
.strip span {
  display: block;
  text-align: center;
  height: 150px;
  line-height: 150px;
  font-size: 80px;
  color: #8b4d43;
  text-shadow: 1px 2px #fff9;
}
.rolling .strip {
  animation: roll 0.45s linear infinite;
  animation-delay: var(--delay);
}
.result {
  text-align: center;
  font-size: 12px;
  letter-spacing: 2px;
  padding: 20px;
  color: #dabd89;
}
.console {
  display: grid;
  grid-template-columns: 1fr 180px 1fr;
  gap: 25px;
  align-items: center;
  text-align: center;
}
.console small {
  font-size: 10px;
  color: #bc9fae;
}
.stats strong,
.time strong {
  display: block;
  font-size: 35px;
  font-weight: 400;
  margin: 15px 0;
}
.stats em {
  font-size: 14px;
  font-style: normal;
}
.time strong {
  font-size: 26px;
  font-variant-numeric: tabular-nums;
}
.console > div > span {
  font-size: 10px;
  color: #987f99;
}
.scan :deep(.placeholder) {
  border-radius: 3px;
}
.scan :deep(.placeholder small) {
  font-size: 8px;
}
.scan :deep(.scan-block p) {
  font-size: 8px;
}
.base {
  text-align: center;
  border-top: 1px solid #d3ab6530;
  font-size: 8px;
  letter-spacing: 4px;
  padding: 18px;
  margin-top: 24px;
  color: #9c8097;
}
footer {
  text-align: center;
  font-size: 9px;
  color: #75627f;
  margin-top: 26px;
}
@keyframes bulb {
  to {
    opacity: 0.35;
    box-shadow: none;
  }
}
@keyframes roll {
  from {
    transform: translateY(0);
  }
  to {
    transform: translateY(-300px);
  }
}
@media (max-width: 700px) {
  header > span {
    display: none;
  }
  .cabinet {
    padding: 14px 14px 0;
  }
  .reel {
    height: 100px;
  }
  .strip span {
    height: 100px;
    line-height: 100px;
    font-size: 60px;
  }
  .strip {
    transform: translateY(-200px);
  }
  .console {
    grid-template-columns: 1fr 1fr;
    gap: 20px;
  }
  .scan {
    grid-row: 1;
    grid-column: 1/-1;
    width: 180px;
    justify-self: center;
  }
  .time strong {
    font-size: 20px;
  }
  .title p {
    line-height: 1.7;
  }
  .base {
    letter-spacing: 1px;
  }
  @keyframes roll {
    from {
      transform: translateY(0);
    }
    to {
      transform: translateY(-200px);
    }
  }
}
.slot-page {
  background: radial-gradient(ellipse at 50% 45%, #3e243b, transparent 65%), #140d1d;
}
.title {
  margin: 26px 0 32px;
}
.title h1 {
  font-family: 'STXingkai', '华文行楷', 'KaiTi', '楷体', serif;
  font-size: clamp(40px, 4.6vw, 84px);
  font-weight: 400;
  letter-spacing: 0.08em;
  background: linear-gradient(170deg, #ffefd0, #c89b62 55%, #e9cba1);
  background-clip: text;
  -webkit-text-fill-color: transparent;
  filter: drop-shadow(0 3px 0 #160b13);
}
.title small {
  font-size: 10px;
  letter-spacing: 6px;
  color: #bc958c;
}
.title p {
  font-family: 'Songti SC', 'SimSun', serif;
  font-size: 12px;
  letter-spacing: 2px;
}
.cabinet {
  max-width: 900px;
  border-radius: 80px 80px 12px 12px;
  padding-top: 22px;
  border-color: #b28d59;
  background: linear-gradient(110deg, #49323b, #251b2b 55%, #4b353b);
  box-shadow:
    0 0 0 5px #b6915d0d,
    0 24px 70px #0005;
}
.marquee {
  padding: 0 25px;
}
.cabinet-label {
  font-family: 'Songti SC', 'SimSun', serif;
  letter-spacing: 6px;
  font-size: 15px;
  padding: 20px;
}
.reels {
  border-radius: 5px;
  border-color: #b4905c50;
  background: #120f19;
}
.reel {
  border-radius: 2px;
  background: linear-gradient(#8b7861, #e4d2b1 30%, #f3e6cb 62%, #9b8465);
}
.strip span {
  color: #88624e;
}
.console {
  grid-template-columns: 1fr 195px 1fr;
}
.scan {
  padding: 8px;
  border: 1px solid #cead7440;
}
.scan :deep(.placeholder) {
  background: #eee3cf;
}
.stats strong {
  font-family: Georgia, serif;
  font-size: 48px;
  color: #e0c49b;
}
.time strong {
  font-family: Georgia, serif;
  color: #e0c49b;
}
.result {
  font-family: 'Songti SC', 'SimSun', serif;
  font-size: 14px;
  letter-spacing: 3px;
}
.base {
  font-size: 7px;
  letter-spacing: 5px;
}
@media (max-width: 700px) {
  .console {
    grid-template-columns: 1fr 1fr;
  }
  .title h1 {
    font-size: 36px;
  }
  .cabinet {
    border-radius: 40px 40px 10px 10px;
  }
  .marquee {
    padding: 0 12px;
  }
  .cabinet-label {
    font-size: 12px;
    letter-spacing: 3px;
  }
  .stats strong {
    font-size: 35px;
  }
}

@media (prefers-reduced-motion: reduce) {
  * {
    animation: none !important;
  }
}
</style>
