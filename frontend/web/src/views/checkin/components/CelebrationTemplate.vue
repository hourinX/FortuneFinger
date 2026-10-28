<script setup lang="ts">
import EventBrand from '@/components/EventBrand.vue'
import CelebrationParticles from './CelebrationParticles.vue'
import { createCheckinTemplateState, type CheckinTemplateProps } from '../useCheckin'
import QrCodeImage from './QrCodeImage.vue'
const props = defineProps<CheckinTemplateProps>()
const qr = props.qr
const state = createCheckinTemplateState(props)
const {
  event,
  closed,
  refreshing,
  refreshSeconds,
  refreshProgress,
  attendance,
  countdown,
  deadlineLabel,
  hasTotal,
  qrActive,
  status,
} = state
</script>

<template>
  <main class="celebration">
    <div
      class="stage-light"
      aria-hidden="true"
    ></div>
    <CelebrationParticles />
    <header>
      <EventBrand />
      <span>ANNUAL CELEBRATION · 年度盛典</span>
    </header>
    <section class="headline">
      <p>今夜，相聚于此</p>
      <h1>
        一起，
        <em>让精彩发生</em>
      </h1>
      <div class="headline-rule">
        <span></span>
        WELCOME TO THE CELEBRATION
        <span></span>
      </div>
      <p class="description">{{ event.description }}</p>
    </section>
    <div class="stage">
      <section class="attendance">
        <span class="label">相聚的人，让今夜闪耀</span>
        <div class="number">
          {{ event.checkedIn }}
          <small>位</small>
        </div>
        <p>
          已抵达现场
          <span v-if="hasTotal">/ 共 {{ event.total }} 位伙伴</span>
        </p>
        <div
          v-if="hasTotal"
          class="progress"
          role="progressbar"
          aria-label="签到进度"
          :aria-valuenow="attendance"
          aria-valuemin="0"
          aria-valuemax="100"
        >
          <i :style="{ width: `${attendance}%` }"></i>
        </div>
        <span v-if="hasTotal" class="ratio">{{ attendance }}% 已签到</span>
      </section>
      <section class="scan">
        <div class="arch">
          <div
            class="seal-frame"
            aria-hidden="true"
          >
            <i></i>
            <i></i>
            <i></i>
            <i></i>
          </div>
          <span class="arch-label">盛会 · 邀君入席</span>
          <div
            class="qr"
            :class="{ refreshing }"
          >
            <QrCodeImage
              v-if="qrActive && qr.value"
              :value="qr.value"
            />
            <strong v-else>{{ status === 'pending' ? '签到尚未开始' : closed ? '已停止签到' : '二维码加载中' }}</strong>
          </div>
          <div class="refresh">
            <span>{{ qrActive ? '30 秒自动更新' : '签到二维码暂不可用' }}</span>
            <b v-if="qrActive">{{ refreshSeconds }}s</b>
          </div>
          <div class="refresh-line">
            <i :style="{ width: `${qrActive ? refreshProgress : 0}%` }"></i>
          </div>
          <p>打开手机扫一扫，开启今夜精彩</p>
        </div>
        <div
          class="pedestal"
          aria-hidden="true"
        ></div>
      </section>
      <section class="deadline">
        <span class="label">{{ closed ? '签到已结束' : '距离签到结束' }}</span>
        <div class="clock">
          <template
            v-for="(part, i) in countdown"
            :key="i"
          >
            <b v-if="i">:</b>
            <span>
              {{ part }}
              <small>{{ ['小时', '分钟', '秒'][i] }}</small>
            </span>
          </template>
        </div>
        <p>{{ deadlineLabel }} 截止</p>
        <div class="status">
          <i></i>
          {{ closed ? '期待下次相聚' : '欢迎入场，席位为你而留' }}
        </div>
      </section>
    </div>
    <footer>
      <span>欢聚盛典</span>
      <span>让每一次相聚，都成为值得珍藏的记忆</span>
      <span>{{ event.company }}</span>
    </footer>
  </main>
</template>

<style scoped>
.celebration {
  position: relative;
  isolation: isolate;
  height: 100dvh;
  overflow: auto;
  padding: 32px 5vw 76px;
  display: flex;
  flex-direction: column;
  color: #ffe6b7;
  background: radial-gradient(ellipse at 50% 65%, #bd332a 0%, #811522 45%, #390b17 100%);
  font-family: 'Sora Variable', 'Microsoft YaHei', sans-serif;
}
.celebration > :not(.confetti):not(.stage-light) {
  position: relative;
}
.stage-light {
  position: absolute;
  inset: 0;
  overflow: hidden;
  pointer-events: none;
  background:
    linear-gradient(115deg, transparent 20%, #ffb95c09 30%, transparent 40%),
    linear-gradient(245deg, transparent 20%, #ffb95c09 30%, transparent 40%);
}
.stage-light::after {
  content: '';
  position: absolute;
  width: 85vw;
  height: 85vw;
  border: 1px solid #ffda9d15;
  border-radius: 50%;
  left: 50%;
  top: 24%;
  transform: translateX(-50%);
  box-shadow:
    0 0 0 35px #ffda9d04,
    0 0 0 70px #ffda9d03;
}
header,
footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 20px;
}
header > span {
  font-size: 10px;
  letter-spacing: 2px;
  opacity: 0.6;
}
.brand {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 21px;
  font-weight: 650;
}
.brand svg {
  width: 30px;
}
.headline {
  text-align: center;
  margin: auto 0 24px;
  padding-top: 34px;
}
.headline > p:first-child {
  font-size: 13px;
  letter-spacing: 8px;
  margin: 0 0 16px;
}
h1 {
  margin: 0;
  font-size: clamp(34px, 4vw, 76px);
  letter-spacing: 3px;
  font-weight: 550;
}
h1 em {
  font-style: normal;
  color: #fff0ce;
}
.headline-rule {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 16px;
  font-size: 9px;
  letter-spacing: 3px;
  margin: 20px 0;
}
.headline-rule span {
  width: 60px;
  height: 1px;
  background: #edc48355;
}
.description {
  color: #e3b8a2;
  font-size: 12px;
}
.stage {
  display: grid;
  grid-template-columns: 1fr 1.1fr 1fr;
  align-items: center;
  gap: 5vw;
  width: 100%;
  max-width: 1500px;
  margin: 10px auto auto;
}
.label {
  font-size: 12px;
  letter-spacing: 2px;
}
.attendance,
.deadline {
  text-align: center;
}
.number {
  font-size: clamp(60px, 7vw, 125px);
  font-weight: 400;
  letter-spacing: -5px;
  margin-top: 10px;
}
.number small {
  font-size: 14px;
  margin-left: 12px;
  letter-spacing: 0;
}
.attendance p,
.deadline p {
  font-size: 12px;
}
.attendance p span {
  opacity: 0.55;
}
.progress {
  height: 3px;
  background: #ffffff15;
  margin: 25px auto 10px;
  max-width: 270px;
}
.progress i,
.refresh-line i {
  display: block;
  height: 100%;
  background: #edc07c;
}
.ratio {
  font-size: 10px;
  opacity: 0.6;
}
.scan {
  text-align: center;
}
.arch {
  position: relative;
  width: min(100%, 345px);
  margin: auto;
  padding: 32px 26px 24px;
  border: 1px solid #ebc38780;
  border-radius: 160px 160px 12px 12px;
  background: linear-gradient(160deg, #f8c78213, #59112155);
  box-shadow:
    inset 0 0 45px #ffbb5c09,
    0 20px 45px #39071330;
}
.arch-label {
  display: block;
  font-size: 12px;
  margin: 15px 0 24px;
  letter-spacing: 4px;
}
.qr {
  width: 100%;
  max-width: 235px;
  aspect-ratio: 1;
  margin: auto;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 13px;
  background: #fff2d8;
  color: #8c2630;
  border: 8px solid #ecd5b1;
  border-radius: 4px;
  transition: box-shadow 0.3s;
}
.qr svg {
  width: 58px;
}
.qr strong {
  font-size: 16px;
  font-weight: 500;
}
.qr small {
  font-size: 9px;
  opacity: 0.6;
}
.refreshing {
  box-shadow: 0 0 35px #ffe1a585;
}
.refresh {
  display: flex;
  justify-content: space-between;
  max-width: 235px;
  margin: 16px auto 8px;
  font-size: 10px;
}
.refresh b {
  font-variant-numeric: tabular-nums;
}
.refresh-line {
  height: 2px;
  max-width: 235px;
  margin: auto;
  background: #ffffff18;
}
.refresh-line i {
  transition: width 0.25s linear;
}
.arch > p {
  font-size: 10px;
  opacity: 0.6;
  margin: 18px 0 0;
}
.pedestal {
  width: 110%;
  height: 24px;
  margin: -3px 0 0 -5%;
  border: 1px solid #eec78a55;
  border-radius: 50%;
  background: linear-gradient(#edb66930, #661422);
  box-shadow: 0 14px 30px #290b1540;
}
.clock {
  display: flex;
  align-items: baseline;
  justify-content: center;
  gap: 12px;
  margin: 30px 0 20px;
  font-variant-numeric: tabular-nums;
}
.clock > span {
  font-size: clamp(28px, 3vw, 54px);
  font-weight: 350;
}
.clock small {
  display: block;
  font-size: 9px;
  margin-top: 8px;
  opacity: 0.5;
}
.clock b {
  font-size: 20px;
  font-weight: 300;
  opacity: 0.4;
}
.deadline p {
  opacity: 0.5;
}
.status {
  font-size: 10px;
  margin-top: 32px;
}
.status i {
  display: inline-block;
  width: 5px;
  height: 5px;
  border-radius: 50%;
  background: #f7d09a;
  margin-right: 8px;
}
footer {
  margin-top: 35px;
  border-top: 1px solid #edc48320;
  padding-top: 18px;
  font-size: 9px;
  color: #dcae99;
}
@media (max-width: 850px) {
  .celebration {
    padding: 24px 24px 85px;
  }
  header > span {
    display: none;
  }
  .headline {
    padding-top: 40px;
  }
  .headline > p:first-child {
    letter-spacing: 4px;
  }
  .description {
    line-height: 1.8;
  }
  .stage {
    grid-template-columns: 1fr 1fr;
    gap: 35px 15px;
  }
  .scan {
    grid-column: 1 / -1;
    grid-row: 1;
  }
  .attendance,
  .deadline {
    align-self: start;
  }
  .number {
    font-size: 60px;
  }
  .label {
    font-size: 10px;
    letter-spacing: 0;
  }
  .clock {
    gap: 5px;
    margin-top: 20px;
  }
  .clock > span {
    font-size: 25px;
  }
  .status {
    line-height: 1.7;
  }
  .attendance p {
    line-height: 1.7;
  }
  footer > span:last-child {
    display: none;
  }
}

.headline {
  padding-top: 25px;
  margin-bottom: 30px;
}
h1 {
  font-family: 'STXingkai', '华文行楷', 'KaiTi', '楷体', serif;
  font-size: clamp(42px, 5.3vw, 96px);
  font-weight: 400;
  letter-spacing: 0.06em;
  color: #f4d3a0;
  filter: drop-shadow(0 3px 2px #43080d);
}
h1 em {
  color: inherit;
  background: linear-gradient(155deg, #fff0ce, #e6ad64 55%, #f9dbab);
  background-clip: text;
  -webkit-text-fill-color: transparent;
}
.headline > p:first-child {
  font-size: 11px;
  letter-spacing: 6px;
  color: #d4ac83;
}
.headline-rule {
  opacity: 0.6;
  font-size: 8px;
}
.headline .description {
  font-family: 'Songti SC', 'SimSun', serif;
  letter-spacing: 2px;
}
.stage {
  grid-template-columns: 1fr 1.2fr 1fr;
}
.scan {
  position: relative;
}
.arch {
  border: 0;
  background: none;
  box-shadow: none;
  border-radius: 0;
  padding: 12px 20px 22px;
  width: min(100%, 350px);
}
.arch-label {
  font-family: 'Songti SC', 'SimSun', serif;
  font-size: 23px;
  letter-spacing: 5px;
  margin: 0 0 26px;
  color: #f3d2a4;
}
.qr {
  max-width: 245px;
  background: #fffaf0;
  border: 0;
  border-radius: 2px;
  outline: 1px solid #efc48a80;
  outline-offset: 10px;
  box-shadow: 0 0 45px #ffbf5a09;
}
.refresh {
  margin-top: 30px;
  max-width: 245px;
  color: #c99f80;
}
.refresh-line {
  max-width: 245px;
}
.arch > p {
  font-family: 'Songti SC', 'SimSun', serif;
  font-size: 12px;
  letter-spacing: 2px;
}
.pedestal {
  display: none;
}
.label {
  font-family: 'Songti SC', 'SimSun', serif;
  font-size: 14px;
  letter-spacing: 2px;
}
.number {
  font-family: Georgia, serif;
  color: #f1d1a3;
}
.clock > span {
  font-family: Georgia, serif;
  color: #f1d1a3;
}
.attendance,
.deadline {
  padding: 25px 0;
  border-top: 1px solid #dfb87b25;
  border-bottom: 1px solid #dfb87b25;
}
@media (max-width: 850px) {
  .stage {
    grid-template-columns: 1fr 1fr;
  }
  .headline h1 {
    font-size: clamp(33px, 7vw, 55px);
  }
  .arch {
    padding: 15px 24px 20px;
  }
  .attendance,
  .deadline {
    padding: 20px 0;
  }
  .label {
    font-size: 11px;
  }
  .headline {
    margin-bottom: 20px;
  }
}

@media (prefers-reduced-motion: reduce) {
  * {
    transition: none !important;
  }
}

/* 金色折屏围合中央入场台，与两侧数据保持独立层次。 */
.arch {
  isolation: isolate;
  width: min(100%, 390px);
  padding: 42px 32px 36px;
  background: linear-gradient(135deg, #e6a24e0b, #4a061e88 55%, #eeac5114);
  border: 1px solid #f2c78550;
  box-shadow:
    0 24px 50px #33061166,
    inset 0 0 50px #e2a5560a;
}
.arch::before,
.arch::after {
  content: '';
  position: absolute;
  top: 18px;
  bottom: 18px;
  width: 30px;
  border: 1px solid #ecc18055;
  background: linear-gradient(90deg, #e8ad4d15, #6d182433);
  z-index: -1;
}
.arch::before {
  left: -24px;
  transform: skewY(-20deg);
}
.arch::after {
  right: -24px;
  transform: skewY(20deg);
}
.seal-frame {
  position: absolute;
  inset: 10px;
  pointer-events: none;
  border: 1px solid #edc78a20;
}
.seal-frame i {
  position: absolute;
  width: 14px;
  height: 14px;
  border: 2px solid #eac189;
  transform: rotate(45deg);
  background: #7d2030;
}
.seal-frame i:nth-child(1) {
  left: -7px;
  top: -7px;
}
.seal-frame i:nth-child(2) {
  right: -7px;
  top: -7px;
}
.seal-frame i:nth-child(3) {
  left: -7px;
  bottom: -7px;
}
.seal-frame i:nth-child(4) {
  right: -7px;
  bottom: -7px;
}
.arch-label {
  font-size: 20px;
  letter-spacing: 6px;
  margin-bottom: 30px;
}
.qr {
  max-width: 220px;
  outline-offset: 7px;
}
.refresh,
.refresh-line {
  max-width: 220px;
}
.pedestal {
  display: block;
  height: 14px;
  width: 110%;
  margin: 14px 0 0 -5%;
  border: 0;
  border-radius: 50%;
  background: radial-gradient(ellipse, #ffd18c45, transparent 70%);
  box-shadow: 0 5px 30px #e9a95818;
}
@media (max-width: 850px) {
  .arch {
    width: 310px;
    padding: 30px 28px;
  }
  .arch::before {
    left: -16px;
    width: 20px;
  }
  .arch::after {
    right: -16px;
    width: 20px;
  }
}
</style>
