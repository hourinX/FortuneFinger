<script setup lang="ts">
import EventBrand from '@/components/EventBrand.vue'
import { createCheckinTemplateState, type CheckinTemplateProps } from '../useCheckin'
import QrCodeImage from './QrCodeImage.vue'

const props = defineProps<CheckinTemplateProps & { theme: 'invitation' | 'celebration' | 'luxe' }>()
const qr = props.qr
const state = createCheckinTemplateState(props)
const {
  event,
  refreshing,
  refreshVersion,
  closed,
  countdown,
  refreshSeconds,
  refreshProgress,
  attendance,
  deadlineLabel,
  hasTotal,
  qrActive,
  status,
} = state
</script>

<template>
  <main
    class="checkin-page"
    :class="theme"
  >
    <header class="masthead">
      <EventBrand tone="light" />
      <span class="event-tag">
        年度相聚
        <span>/</span>
        ANNUAL GATHERING
      </span>
    </header>

    <div class="content">
      <section
        class="welcome"
        aria-labelledby="event-title"
      >
        <div class="eyebrow">
          <span></span>
          很高兴，在这里见到你
        </div>
        <h1
          id="event-title"
          class="invitation-title"
        >
          <span class="ink-line">{{ event.title.split('，')[0] }}，</span>
          <span class="clay-line">
            {{ event.title.split('，').slice(1).join('，') }}
            <i
              class="invitation-seal"
              aria-hidden="true"
            >
              相聚
            </i>
          </span>
        </h1>
        <div class="invitation-dedication">
          <span>一封写给你的邀请函</span>
          <i></i>
          <span>WITH YOU, ALWAYS.</span>
        </div>
        <p class="description">{{ event.description }}</p>
        <div class="attendance">
          <div class="attendance-heading">
            <span>现场签到</span>
            <span class="live">
              <i></i>
              {{ closed ? '签到已结束' : '入场进行中' }}
            </span>
          </div>
          <div class="attendance-numbers">
            <strong>{{ event.checkedIn }}</strong>
            <span v-if="hasTotal">
              / {{ event.total }}
              <small>位伙伴</small>
            </span>
          </div>
          <div
            v-if="hasTotal"
            class="track"
            role="progressbar"
            aria-label="签到进度"
            :aria-valuenow="attendance"
            aria-valuemin="0"
            aria-valuemax="100"
          >
            <span :style="{ width: `${attendance}%` }"></span>
          </div>
          <div v-if="hasTotal" class="progress-caption">
            <span>每一份到场，都让相聚更完整。</span>
            <strong>{{ attendance }}%</strong>
          </div>
        </div>
      </section>

      <section
        class="ticket"
        aria-labelledby="scan-title"
      >
        <div class="ticket-main">
          <div class="ticket-top">
            <span>YOUR INVITATION</span>
            <span class="ticket-number">01 / CHECK IN</span>
          </div>
          <h2 id="scan-title">{{ closed ? '本次签到已结束' : '赴一场相聚之约' }}</h2>
          <p>打开手机相机或扫一扫，完成活动签到</p>
          <div
            class="qr-frame"
            :class="{ refreshing, closed }"
          >
            <div
              :key="refreshVersion"
              class="qr-placeholder"
            >
              <span class="corner corner-tl"></span>
              <span class="corner corner-tr"></span>
              <span class="corner corner-bl"></span>
              <span class="corner corner-br"></span>
              <QrCodeImage
                v-if="qrActive && qr.value"
                :value="qr.value"
              />
              <strong v-else>{{ status === 'pending' ? '签到尚未开始' : closed ? '已停止签到' : '二维码加载中' }}</strong>
            </div>
            <div
              v-if="refreshing && qrActive"
              class="refresh-feedback"
              role="status"
            >
              <span>↻</span>
              二维码已刷新
            </div>
          </div>
          <div class="refresh-caption">
            <svg
              :class="{ spinning: refreshing }"
              viewBox="0 0 24 24"
              fill="none"
              aria-hidden="true"
            >
              <path
                d="M20 8a8 8 0 1 0 0 8M20 3v5h-5"
                stroke="currentColor"
                stroke-width="1.6"
                stroke-linecap="round"
                stroke-linejoin="round"
              />
            </svg>
            <span>{{ qrActive ? '二维码每 30 秒自动更新' : '签到二维码暂不可用' }}</span>
            <strong v-if="qrActive">{{ refreshSeconds }}s</strong>
          </div>
          <div class="refresh-track">
            <span :style="{ width: `${qrActive ? refreshProgress : 0}%` }"></span>
          </div>
        </div>
        <div class="ticket-bottom">
          <div class="deadline-label">
            <span>{{ closed ? '期待下次相聚' : '距离签到截止' }}</span>
            <small>{{ deadlineLabel }}</small>
          </div>
          <div
            class="countdown"
            aria-label="距离签到截止的时分秒"
          >
            <template
              v-for="(part, index) in countdown"
              :key="index"
            >
              <b v-if="index > 0">:</b>
              <span>
                {{ part }}
                <small>{{ ['时', '分', '秒'][index] }}</small>
              </span>
            </template>
          </div>
        </div>
      </section>
    </div>

    <footer class="page-footer">
      <span>GOOD PEOPLE. GREAT MOMENTS.</span>
      <span>
        由 {{ event.company }} 陪伴每一次相聚
        <i>↗</i>
      </span>
    </footer>
  </main>
</template>

<style scoped>
.checkin-page {
  --ink: #233c34;
  --muted: #7b8279;
  --accent: #d86942;
  width: 100%;
  height: 100dvh;
  overflow: auto;
  color: var(--ink);
  background: #f4f2e9;
  padding: clamp(24px, 3.5vw, 64px) clamp(24px, 6vw, 120px) 24px;
  display: flex;
  flex-direction: column;
  font-family: 'Sora Variable', 'PingFang SC', 'Microsoft YaHei', sans-serif;
}
.masthead,
.brand,
.ticket-top,
.attendance-heading,
.progress-caption,
.page-footer,
.refresh-caption {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.masthead {
  gap: 20px;
}
.brand {
  gap: 12px;
  color: var(--ink);
  text-decoration: none;
  font-size: clamp(19px, 1.6vw, 29px);
  font-weight: 700;
  letter-spacing: -1px;
}
.brand svg {
  width: 40px;
  height: 40px;
  background: var(--ink);
  color: #f4f2e9;
  padding: 6px;
  border-radius: 11px;
}
.event-tag {
  color: var(--muted);
  font-size: 10px;
  letter-spacing: 1.5px;
}
.event-tag span {
  margin: 0 12px;
  color: #b9bbaf;
}
.content {
  display: grid;
  grid-template-columns: 1.15fr 1fr;
  gap: clamp(48px, 9vw, 180px);
  align-items: center;
  flex: 1;
  width: 100%;
  max-width: 1600px;
  margin: 0 auto;
  padding: 48px 0;
}
.eyebrow {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 13px;
  letter-spacing: 2px;
}
.eyebrow > span {
  width: 28px;
  height: 2px;
  background: var(--accent);
}
h1 {
  position: relative;
  width: fit-content;
  margin: 30px 0 24px;
  font-size: clamp(40px, 4.7vw, 90px);
  line-height: 1.35;
  font-weight: 550;
  letter-spacing: -3px;
}
h1 > span {
  color: var(--accent);
}
.spark {
  position: absolute;
  width: clamp(34px, 4vw, 70px);
  height: auto;
  right: -60px;
  top: -15px;
  color: var(--accent);
  transform: rotate(12deg);
}
.description {
  max-width: 440px;
  font-size: clamp(14px, 1.05vw, 19px);
  line-height: 2;
  color: #747b71;
}
.attendance {
  max-width: 440px;
  margin-top: clamp(36px, 5vw, 80px);
}
.attendance-heading {
  font-size: 13px;
}
.live {
  font-size: 11px;
  color: #667568;
  display: flex;
  align-items: center;
  gap: 7px;
}
.live i {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: #5c8066;
}
.attendance-numbers {
  display: flex;
  align-items: baseline;
  gap: 14px;
  margin: 14px 0 18px;
  font-variant-numeric: tabular-nums;
}
.attendance-numbers strong {
  font-size: clamp(48px, 4.3vw, 74px);
  font-weight: 450;
  letter-spacing: -3px;
}
.attendance-numbers > span {
  font-size: 21px;
  color: #93988d;
}
.attendance-numbers small {
  font-size: 12px;
  margin-left: 8px;
}
.track {
  height: 5px;
  background: #dedfd4;
  border-radius: 10px;
  overflow: hidden;
}
.track > span {
  display: block;
  height: 100%;
  background: var(--ink);
  border-radius: inherit;
}
.progress-caption {
  margin-top: 13px;
  font-size: 11px;
  color: var(--muted);
  gap: 12px;
}
.progress-caption strong {
  font-weight: 500;
  color: var(--ink);
}
.ticket {
  width: 100%;
  max-width: 520px;
  justify-self: end;
  background: #fffefa;
  border-radius: 20px;
  box-shadow:
    0 22px 65px #3644360c,
    0 2px 5px #36443605;
  transform: rotate(1deg);
}
.ticket-main {
  padding: clamp(24px, 2.6vw, 46px);
}
.ticket-top {
  color: #858e7f;
  font-size: 9px;
  letter-spacing: 1.5px;
  gap: 10px;
}
.ticket-number {
  color: #a1a598;
  font-size: 8px;
}
h2 {
  margin: 30px 0 9px;
  font-size: clamp(23px, 1.7vw, 31px);
  font-weight: 550;
  text-align: center;
  letter-spacing: -1px;
}
.ticket-main > p {
  color: var(--muted);
  font-size: 11px;
  text-align: center;
  margin: 0 0 24px;
}
.qr-frame {
  position: relative;
  width: min(100%, 270px);
  aspect-ratio: 1;
  padding: 16px;
  margin: 0 auto;
  background: #f5f5ee;
  border: 1px solid #e6e9df;
  border-radius: 9px;
  transition: box-shadow 0.3s;
}
.qr-placeholder {
  width: 100%;
  height: 100%;
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background: #fffefa;
  color: var(--ink);
}
.qr-placeholder > svg {
  width: 58px;
  margin-bottom: 16px;
  opacity: 0.7;
}
.qr-placeholder strong {
  font-size: 14px;
  font-weight: 500;
}
.qr-placeholder > span:last-child {
  font-size: 10px;
  color: #92998b;
  margin-top: 9px;
}
.corner {
  position: absolute;
  width: 18px;
  height: 18px;
  border: solid #71806a;
}
.corner-tl {
  top: 0;
  left: 0;
  border-width: 2px 0 0 2px;
}
.corner-tr {
  top: 0;
  right: 0;
  border-width: 2px 2px 0 0;
}
.corner-bl {
  bottom: 0;
  left: 0;
  border-width: 0 0 2px 2px;
}
.corner-br {
  bottom: 0;
  right: 0;
  border-width: 0 2px 2px 0;
}
.refreshing {
  box-shadow: 0 0 0 5px #dce6d780;
}
.refresh-feedback {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  background: #eef2e9f5;
  border-radius: inherit;
  font-size: 13px;
}
.refresh-feedback > span,
.spinning {
  animation: spin 0.9s ease;
}
.refresh-feedback > span {
  font-size: 24px;
}
.closed {
  opacity: 0.6;
}
.refresh-caption {
  max-width: 270px;
  margin: 17px auto 9px;
  color: var(--muted);
  gap: 6px;
  font-size: 10px;
}
.refresh-caption svg {
  width: 14px;
  height: 14px;
  flex-shrink: 0;
}
.refresh-caption > span {
  flex: 1;
}
.refresh-caption strong {
  font-size: 11px;
  color: var(--ink);
  min-width: 25px;
  text-align: right;
  font-variant-numeric: tabular-nums;
}
.refresh-track {
  max-width: 270px;
  height: 2px;
  margin: auto;
  background: #edf0e7;
  overflow: hidden;
}
.refresh-track > span {
  display: block;
  height: 100%;
  background: #9daa90;
  transition: width 0.25s linear;
}
.ticket-bottom {
  position: relative;
  padding: 24px clamp(24px, 2.6vw, 46px) 28px;
  border-top: 1px dashed #d9dece;
}
.ticket-bottom::before,
.ticket-bottom::after {
  content: '';
  position: absolute;
  top: -10px;
  width: 20px;
  height: 20px;
  background: #f4f2e9;
  border-radius: 50%;
}
.ticket-bottom::before {
  left: -10px;
}
.ticket-bottom::after {
  right: -10px;
}
.deadline-label {
  display: flex;
  justify-content: space-between;
  font-size: 11px;
  color: var(--muted);
  gap: 10px;
}
.deadline-label small {
  font-size: 10px;
}
.countdown {
  display: flex;
  align-items: baseline;
  justify-content: center;
  gap: 18px;
  margin-top: 15px;
  font-variant-numeric: tabular-nums;
}
.countdown > span {
  font-size: clamp(28px, 2.4vw, 42px);
  font-weight: 450;
  letter-spacing: -1px;
}
.countdown small {
  font-size: 9px;
  color: #9ba18f;
  margin-left: 6px;
  letter-spacing: 0;
}
.countdown b {
  font-size: 24px;
  font-weight: 300;
  color: #bcc2b2;
}
.page-footer {
  border-top: 1px solid #dce0d3;
  padding-top: 20px;
  gap: 20px;
  color: #92998a;
  font-size: 10px;
}
.page-footer > span:first-child {
  letter-spacing: 2px;
}
.page-footer i {
  font-style: normal;
  margin-left: 10px;
  color: var(--ink);
}
@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}
@media (min-width: 2200px) {
  .content {
    max-width: 2100px;
  }
  .ticket {
    max-width: 650px;
  }
  .qr-frame {
    width: 340px;
  }
  .refresh-caption,
  .refresh-track {
    max-width: 340px;
  }
}
@media (max-width: 850px) {
  .checkin-page {
    padding: 24px;
  }
  .event-tag {
    display: none;
  }
  .content {
    grid-template-columns: 1fr;
    gap: 32px;
    padding: 40px 0;
    max-width: 520px;
  }
  .welcome {
    text-align: center;
  }
  .eyebrow {
    justify-content: center;
    font-size: 11px;
  }
  h1 {
    margin: 20px auto;
    font-size: clamp(36px, 7vw, 55px);
    letter-spacing: -2px;
  }
  .spark {
    right: -32px;
    width: 32px;
    top: -8px;
  }
  .description {
    margin: 0 auto;
    font-size: 13px;
  }
  .attendance {
    margin: 28px auto 0;
    text-align: left;
  }
  .attendance-numbers {
    margin: 8px 0 12px;
  }
  .attendance-numbers strong {
    font-size: 44px;
  }
  .ticket {
    justify-self: center;
    transform: none;
  }
  .ticket-main {
    padding: 26px;
  }
  .ticket-main h2 {
    margin-top: 24px;
  }
  .page-footer {
    flex-direction: column;
    gap: 8px;
    text-align: center;
    font-size: 9px;
  }
}

.invitation .eyebrow {
  color: #898471;
  font-size: 11px;
  letter-spacing: 0.2em;
}
.invitation .eyebrow > span {
  height: 1px;
  width: 34px;
  background: #b39772;
}
.invitation .invitation-title {
  font-family: 'STXingkai', '华文行楷', 'KaiTi', '楷体', 'Songti SC', serif;
  font-size: clamp(42px, 4.8vw, 88px);
  font-weight: 400;
  line-height: 1.4;
  letter-spacing: 0.04em;
  margin: 32px 0 22px;
  transform: rotate(-2deg);
}
.invitation-title > span {
  display: block;
  width: fit-content;
}
.invitation-title .ink-line {
  color: #354c41;
  background: linear-gradient(125deg, #263f35, #516450 48%, #344c40 75%);
  background-clip: text;
  -webkit-text-fill-color: transparent;
}
.invitation-title .clay-line {
  position: relative;
  margin-left: 0.4em;
  color: #ae684c;
  background: linear-gradient(110deg, #92553f, #be8260 46%, #a35c43 75%);
  background-clip: text;
  -webkit-text-fill-color: transparent;
}
.invitation-seal {
  position: absolute;
  right: -22px;
  bottom: 12px;
  padding: 5px 3px;
  border: 1px solid #a96e4c80;
  color: #a36b4f;
  -webkit-text-fill-color: #a36b4f;
  font:
    normal 10px 'SimSun',
    serif;
  writing-mode: vertical-rl;
  letter-spacing: 2px;
  transform: rotate(4deg);
  border-radius: 2px;
}
.invitation-dedication {
  display: flex;
  align-items: center;
  gap: 12px;
  margin: 22px 0 18px;
  color: #8f8875;
  font-size: 10px;
  letter-spacing: 1px;
}
.invitation-dedication > i {
  width: 24px;
  height: 1px;
  background: #c9bda5;
}
.invitation-dedication > span:last-child {
  font-size: 8px;
  letter-spacing: 2px;
}
.invitation .description {
  font-family: 'Songti SC', 'SimSun', serif;
  max-width: 410px;
  font-size: clamp(13px, 1vw, 17px);
  line-height: 2.1;
  letter-spacing: 0.06em;
  color: #7f8374;
}
.invitation h2 {
  font-family: 'Songti SC', 'SimSun', serif;
  font-size: clamp(22px, 1.8vw, 32px);
  font-weight: 400;
  letter-spacing: 0.12em;
}
.invitation .ticket-top {
  font-family: Georgia, serif;
  letter-spacing: 0.16em;
}
.invitation .attendance-numbers strong,
.invitation .countdown > span {
  font-family: Georgia, 'Times New Roman', serif;
  font-weight: 400;
}
.invitation .attendance-heading {
  font-family: 'Songti SC', 'SimSun', serif;
  letter-spacing: 1px;
}
@media (max-width: 850px) {
  .invitation .invitation-title {
    font-size: clamp(37px, 7vw, 56px);
    margin: 24px auto;
  }
  .invitation-dedication {
    justify-content: center;
    gap: 8px;
    flex-wrap: wrap;
  }
  .invitation-dedication > span:last-child {
    font-size: 7px;
  }
  .invitation .description {
    font-size: 13px;
  }
  .invitation-seal {
    right: -17px;
    bottom: 8px;
    font-size: 8px;
  }
}
@media (prefers-reduced-motion: reduce) {
  *,
  *::before,
  *::after {
    animation: none !important;
    transition: none !important;
  }
}
</style>
