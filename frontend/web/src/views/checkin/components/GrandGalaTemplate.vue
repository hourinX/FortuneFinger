<script setup lang="ts">
import EventBrand from '@/components/EventBrand.vue'
import { createCheckinTemplateState, type CheckinTemplateProps } from '../useCheckin'
import CelebrationParticles from './CelebrationParticles.vue'
const props = defineProps<CheckinTemplateProps>()
const { event, countdown, refreshing, closed, refreshSeconds } = createCheckinTemplateState(props)
const year = new Date(event.deadline).getFullYear()
// 确定性的建筑轮廓，避免组件重渲染时城市跳变。
const buildings = Array.from({ length: 66 }, (_, i) => ({
  x: i * 30 - 20,
  width: 16 + ((i * 7) % 17),
  height: 35 + ((i * 47) % 125) + (i % 9 === 0 ? 65 : 0),
}))
</script>

<template>
  <main class="gala">
    <div
      class="artwork"
      aria-hidden="true"
    >
      <svg
        class="scenery"
        viewBox="0 0 1920 1080"
        preserveAspectRatio="xMidYMid slice"
      >
        <defs>
          <linearGradient
            id="gala-silk"
            x1="0"
            y1="0"
            x2="0.8"
            y2="1"
          >
            <stop stop-color="#300105" />
            <stop
              offset=".28"
              stop-color="#9f0e12"
            />
            <stop
              offset=".45"
              stop-color="#f83e21"
            />
            <stop
              offset=".51"
              stop-color="#ff9b51"
            />
            <stop
              offset=".56"
              stop-color="#ba1712"
            />
            <stop
              offset=".8"
              stop-color="#5d030c"
            />
            <stop
              offset="1"
              stop-color="#e9351a"
            />
          </linearGradient>
          <linearGradient id="gala-fold">
            <stop stop-color="#38030a" />
            <stop
              offset=".35"
              stop-color="#a20c11"
            />
            <stop
              offset=".48"
              stop-color="#f44b28"
            />
            <stop
              offset=".55"
              stop-color="#8d0a10"
            />
            <stop
              offset="1"
              stop-color="#350306"
            />
          </linearGradient>
          <linearGradient
            id="gala-city"
            x1="0"
            y1="1"
            x2="0"
            y2="0"
          >
            <stop stop-color="#f69e4350" />
            <stop
              offset="1"
              stop-color="#5a170b"
            />
          </linearGradient>
          <linearGradient
            id="gala-road"
            x1="0"
            y1="0"
            x2="0"
            y2="1"
          >
            <stop
              stop-color="#ffd89b"
              stop-opacity=".03"
            />
            <stop
              offset="1"
              stop-color="#ffac55"
              stop-opacity=".7"
            />
          </linearGradient>
          <radialGradient id="gala-glow">
            <stop
              stop-color="#ffe6a6"
              stop-opacity=".7"
            />
            <stop
              offset=".25"
              stop-color="#ff9e39"
              stop-opacity=".3"
            />
            <stop
              offset="1"
              stop-color="#ef390c"
              stop-opacity="0"
            />
          </radialGradient>
          <pattern
            id="gala-windows"
            width="9"
            height="12"
            patternUnits="userSpaceOnUse"
          >
            <rect
              x="2"
              y="3"
              width="2"
              height="5"
              fill="#ffb968"
              opacity=".5"
            />
          </pattern>
          <filter id="gala-soft"><feGaussianBlur stdDeviation="3" /></filter>
        </defs>
        <ellipse
          cx="960"
          cy="805"
          rx="850"
          ry="230"
          fill="url(#gala-glow)"
        />
        <g class="silk silk-left">
          <path
            d="M-140 120 C160 135 130 350 405 413 C570 455 470 565 285 590 C100 615 50 760 -120 720 L-90 610 C70 655 90 535 290 520 C440 505 490 490 338 459 C80 411 30 260 -140 315Z"
            fill="url(#gala-silk)"
          />
          <path
            d="M-140 120 C160 135 130 350 405 413 C570 455 470 565 285 590"
            fill="none"
            stroke="#ffb264"
            stroke-opacity=".4"
            stroke-width="1.5"
          />
          <path
            d="M-100 620 C90 520 160 440 280 448 C160 460 130 565 -100 725Z"
            fill="url(#gala-fold)"
          />
        </g>
        <g class="silk silk-right">
          <path
            d="M1960 -50 C1660 80 1870 230 1625 332 C1470 400 1550 530 1740 550 C1900 575 1810 725 1990 695 L2040 565 C1890 600 1920 490 1750 466 C1580 442 1660 411 1765 370 C2000 278 1790 120 2030 80Z"
            fill="url(#gala-silk)"
          />
          <path
            d="M1960 -50 C1660 80 1870 230 1625 332 C1470 400 1550 530 1740 550"
            fill="none"
            stroke="#ffb264"
            stroke-opacity=".5"
            stroke-width="2"
          />
          <path
            d="M1990 490 C1710 555 1695 648 1460 655 C1610 715 1770 637 1950 610Z"
            fill="url(#gala-fold)"
          />
        </g>
        <g opacity=".7">
          <g
            v-for="(building, i) in buildings"
            :key="i"
          >
            <rect
              :x="building.x"
              :y="805 - building.height"
              :width="building.width"
              :height="building.height"
              fill="url(#gala-city)"
              stroke="#e99d48"
              stroke-opacity=".55"
              stroke-width=".7"
            />
            <rect
              :x="building.x + 2"
              :y="808 - building.height"
              :width="building.width - 4"
              :height="building.height - 3"
              fill="url(#gala-windows)"
            />
            <path
              v-if="i % 9 === 0"
              :d="`M${building.x + building.width / 2} ${805 - building.height}v-28`"
              stroke="#ffd39a"
              stroke-width="1"
            />
          </g>
        </g>
        <g
          transform="translate(0 1610) scale(1 -1)"
          opacity=".15"
          filter="url(#gala-soft)"
        >
          <rect
            v-for="(building, i) in buildings"
            :key="i"
            :x="building.x"
            :y="805 - building.height"
            :width="building.width"
            :height="building.height"
            fill="url(#gala-windows)"
          />
        </g>
        <path
          d="M1450 807 Q1680 610 1900 807 M1450 807H1920 M1510 760V807 M1580 713V807 M1650 695V807 M1720 700V807 M1790 725V807 M1860 770V807"
          fill="none"
          stroke="#f2b46c"
          stroke-width="1.5"
          opacity=".5"
        />
        <g
          stroke="url(#gala-road)"
          fill="none"
        >
          <path
            v-for="n in 17"
            :key="n"
            :d="`M960 805 L${(n - 1) * 200 - 640} 1110`"
            :stroke-width="n % 3 === 0 ? 2 : 0.6"
          />
          <path
            d="M0 853H1920 M0 929H1920 M0 1050H1920"
            opacity=".25"
          />
        </g>
        <path
          class="horizon"
          d="M0 810 Q960 795 1920 810"
          fill="none"
          stroke="#ffd99d"
          stroke-width="2"
        />
        <ellipse
          cx="960"
          cy="809"
          rx="450"
          ry="3"
          fill="#ffe1a0"
          opacity=".45"
          filter="url(#gala-soft)"
        />
      </svg>
      <CelebrationParticles />
    </div>
    <div class="gala-brand"><EventBrand /></div>
    <section class="composition">
      <div class="year">
        <i></i>
        {{ year }} 年度盛典
        <i></i>
      </div>
      <h1>
        同心致远
        <span>见未来</span>
      </h1>
      <p class="event-name">{{ year }} 年度公司年会</p>
      <p class="motto">
        新起点
        <i>/</i>
        新征程
        <i>/</i>
        新未来
      </p>
      <div class="scan-area">
        <div
          class="qr"
          :class="{ refreshing }"
        >
          <svg
            viewBox="0 0 64 64"
            fill="none"
            aria-hidden="true"
          >
            <path
              d="M8 24V8h16M40 8h16v16M56 40v16H40M24 56H8V40M25 25h14v14H25Z"
              stroke="currentColor"
              stroke-width="3"
            />
          </svg>
          <strong>{{ closed ? '签到已截止' : refreshing ? '二维码已刷新' : '签到二维码' }}</strong>
          <small>演示占位 · 不可扫描</small>
        </div>
        <h2>扫码签到</h2>
        <p class="invitation">— 期待与您共赴美好 —</p>
        <span class="refresh-label">
          {{ closed ? '签到通道已关闭' : `二维码每 30 秒更新 · ${refreshSeconds}s` }}
        </span>
      </div>
      <footer>
        <span>{{ event.company }}</span>
        <span>
          已签到
          <b>{{ event.checkedIn }}</b>
          / {{ event.total }} 人
        </span>
        <span>
          距签到截止
          <b>{{ countdown.join(' : ') }}</b>
        </span>
      </footer>
    </section>
  </main>
</template>

<style scoped>
.gala {
  position: relative;
  isolation: isolate;
  width: 100%;
  height: 100dvh;
  min-height: 540px;
  overflow: hidden;
  color: #efc18a;
  background: radial-gradient(
    ellipse at 50% 77%,
    #e34c17 0%,
    #a4120c 28%,
    #480704 64%,
    #210302 100%
  );
  font-family: 'Microsoft YaHei', 'PingFang SC', sans-serif;
}
.artwork {
  position: absolute;
  inset: 0;
  pointer-events: none;
  z-index: -1;
}
.scenery {
  width: 100%;
  height: 100%;
}
.composition {
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 5vh 5vw 70px;
  text-align: center;
}
.year {
  display: flex;
  align-items: center;
  gap: 22px;
  font-size: clamp(12px, 1.15vw, 22px);
  letter-spacing: 0.35em;
}
.year i {
  width: 65px;
  height: 1px;
  background: linear-gradient(90deg, transparent, #efc18a70);
}
.year i:last-child {
  transform: rotate(180deg);
}
h1 {
  font-family: 'STXingkai', '华文行楷', 'KaiTi', '楷体', 'Songti SC', serif;
  font-size: clamp(42px, 6.4vw, 125px);
  line-height: 1.15;
  letter-spacing: 0.04em;
  margin: 2.5vh 0 1.2vh;
  font-weight: 700;
  color: #f6d39b;
  background: linear-gradient(
    165deg,
    #fff4dc 10%,
    #efc18a 35%,
    #b67436 46%,
    #ffe2a9 52%,
    #d59955 74%,
    #fff0cf 90%
  );
  background-clip: text;
  -webkit-text-fill-color: transparent;
  filter: drop-shadow(0 3px 1px #4e1307) drop-shadow(0 0 12px #ffae3b20);
}
.event-name {
  font-size: clamp(14px, 1.6vw, 30px);
  letter-spacing: 0.3em;
  color: #fff0d6;
  margin: 0;
}
.motto {
  font-size: clamp(10px, 0.9vw, 17px);
  letter-spacing: 0.2em;
  margin: 1.5vh 0 2vh;
}
.motto i {
  font-style: normal;
  opacity: 0.5;
  margin: 0 15px;
}
.scan-area {
  display: flex;
  flex-direction: column;
  align-items: center;
}
.qr {
  width: clamp(145px, 21vh, 290px);
  aspect-ratio: 1;
  border: 4px solid #efc18a;
  border-radius: 10px;
  background: #fff;
  color: #28211b;
  box-shadow:
    0 0 0 1px #fff1c570,
    0 0 28px #ffa34c15;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;
  transition: box-shadow 0.4s;
}
.qr svg {
  width: 30%;
  color: #6d4e31;
}
.qr strong {
  font-size: 13px;
  font-weight: 400;
}
.qr small {
  font-size: 9px;
  color: #96908a;
}
.refreshing {
  box-shadow: 0 0 35px #ffcf7870;
}
h2 {
  font-size: clamp(18px, 1.6vw, 30px);
  letter-spacing: 0.3em;
  font-weight: 500;
  color: #fff4dc;
  margin: 1.5vh 0 0.7vh;
}
.invitation {
  font-size: clamp(10px, 0.85vw, 16px);
  letter-spacing: 0.15em;
  margin: 0;
}
.refresh-label {
  font-size: 10px;
  opacity: 0.6;
  margin-top: 10px;
}
footer {
  margin-top: auto;
  width: 100%;
  display: flex;
  justify-content: space-between;
  gap: 20px;
  font-size: clamp(10px, 0.85vw, 15px);
  text-shadow: 0 2px 5px #300;
}
footer b {
  font-weight: 400;
  color: #fff0d6;
  font-variant-numeric: tabular-nums;
}
.silk {
  transform-box: fill-box;
  transform-origin: center;
  animation: silk-breathe 14s ease-in-out infinite alternate;
}
.silk-right {
  animation-delay: -7s;
}
.horizon {
  animation: horizon-light 7s ease-in-out infinite alternate;
}
@keyframes silk-breathe {
  to {
    transform: translateY(8px) rotate(0.8deg);
  }
}
@keyframes horizon-light {
  from {
    opacity: 0.35;
  }
  to {
    opacity: 0.85;
  }
}
@media (max-width: 700px) {
  .gala {
    height: 100dvh;
    min-height: 650px;
    overflow: auto;
  }
  .composition {
    padding: 35px 18px 85px;
  }
  h1 {
    font-size: clamp(32px, 8vw, 55px);
    margin: 25px 0 18px;
  }
  .year {
    font-size: 11px;
  }
  .year i {
    width: 25px;
  }
  .event-name {
    font-size: 14px;
  }
  .motto {
    margin: 18px 0 25px;
  }
  .qr {
    width: 180px;
  }
  .scenery {
    opacity: 0.65;
  }
  footer {
    flex-wrap: wrap;
    justify-content: center;
    font-size: 10px;
  }
  footer > span:first-child {
    width: 100%;
  }
}
@media (prefers-reduced-motion: reduce) {
  * {
    animation: none !important;
    transition: none !important;
  }
}
.gala-brand {
  position: absolute;
  top: 25px;
  left: 4vw;
}
@media (max-width: 700px) {
  .gala-brand {
    top: 16px;
    left: 18px;
  }
  .composition {
    padding-top: 76px;
  }
}
</style>
