<script setup lang="ts">
import { computed, onMounted, onBeforeUnmount, ref, watch } from "vue";
import { useRoute } from "vue-router";
import gsap from "gsap";
import EventBrand from "@/components/EventBrand.vue";
import QrCodeImage from "../QrCodeImage.vue";
import Atmosphere from "./Atmosphere.vue";
import {
  createCheckinTemplateState,
  type CheckinTemplateProps,
} from "../../useCheckin";

const props = defineProps<CheckinTemplateProps>();
const {
  event,
  countdown,
  closed,
  qrActive,
  status,
  refreshSeconds,
  refreshProgress,
  attendance,
  hasTotal,
} = createCheckinTemplateState(props);
const route = useRoute();
const preview = computed(() => route.query.preview === "1");
const root = ref<HTMLElement>();
const count = ref<HTMLElement>();
const accepting = computed(() => qrActive.value && !closed.value);
const stateLabel = computed(() =>
  status.value === "pending"
    ? "静候开场"
    : closed.value
      ? "签到已结束"
      : "正在签到 · 欢迎入场",
);
const year = computed(() => new Date(event.deadline).getFullYear());
let context: gsap.Context | undefined;
let countTween: gsap.core.Tween | undefined;

onMounted(() => {
  if (matchMedia("(prefers-reduced-motion: reduce)").matches) return;
  context = gsap.context(() => {
    gsap.from("[data-enter]", {
      opacity: 0,
      y: 24,
      duration: 1.4,
      stagger: 0.14,
      ease: "power3.out",
      clearProps: "transform,opacity",
    });
  }, root.value);
});
watch(
  () => event.checkedIn,
  (value, oldValue) => {
    if (
      value <= oldValue ||
      !count.value ||
      matchMedia("(prefers-reduced-motion: reduce)").matches
    )
      return;
    countTween?.kill();
    countTween = gsap.fromTo(
      count.value,
      { color: "#fff7de" },
      { color: "#edc990", duration: 1.5 },
    );
  },
);
onBeforeUnmount(() => {
  countTween?.kill();
  context?.revert();
});
</script>

<template>
  <main ref="root" class="crimson-gala">
    <Atmosphere />
    <header class="masthead" data-enter>
      <EventBrand />
      <span class="edition">{{ event.company }} · ANNUAL GATHERING</span>
      <span class="live-state"
        ><i :class="{ active: accepting }"></i>{{ stateLabel }}</span
      >
    </header>

    <div class="stage">
      <section class="invitation">
        <div class="eyebrow" data-enter>
          <span></span>{{ year }} · 年 度 盛 典
        </div>
        <h1 data-enter>此刻相聚<br /><span>共赴新程</span><i>。</i></h1>
        <div class="dedication" data-enter>
          <span class="seal">盛会<br />有你</span>
          <div>
            <p>{{ event.title }}</p>
            <span>{{ event.description }}</span>
          </div>
        </div>

        <div class="gathering" data-enter>
          <section class="attendance">
            <p class="label">已赴此约 <span>GUESTS ARRIVED</span></p>
            <div class="count">
              <strong ref="count">{{ event.checkedIn }}</strong
              ><span
                >位<span v-if="hasTotal"> / {{ event.total }} 位</span></span
              >
            </div>
            <div
              v-if="hasTotal"
              class="attendance-line"
              role="progressbar"
              aria-label="签到进度"
              :aria-valuenow="attendance"
              aria-valuemin="0"
              aria-valuemax="100"
            >
              <i :style="{ width: `${attendance}%` }"></i>
            </div>
            <p class="caption">
              {{
                hasTotal
                  ? `${attendance}% 的伙伴已抵达现场`
                  : "每一份到来，都让今夜更完整"
              }}
            </p>
          </section>
          <section class="deadline">
            <p class="label">
              {{
                status === "pending"
                  ? "签到截止倒计时"
                  : closed
                    ? "本场签到已结束"
                    : "距离签到结束"
              }}
            </p>
            <div class="clock" aria-label="签到截止倒计时">
              <template v-for="(part, index) in countdown" :key="index">
                <i v-if="index">:</i>
                <span
                  >{{ part
                  }}<small>{{ ["时", "分", "秒"][index] }}</small></span
                >
              </template>
            </div>
            <p class="caption">
              {{ closed ? "请留步，精彩即将继续" : "留一点期待，赴一场相聚" }}
            </p>
          </section>
        </div>
      </section>

      <section class="admission" data-enter aria-label="扫码签到入场券">
        <div class="ticket-top">
          <span>INVITATION</span><span>入 场 凭 证</span>
        </div>
        <div class="ticket-body">
          <span class="ticket-kicker">席位已备 · 恭迎莅临</span>
          <h2>{{ closed ? "相聚仍在继续" : "扫码，入席" }}</h2>
          <div class="qr-frame">
            <svg class="qr-progress" viewBox="0 0 300 300" aria-hidden="true">
              <rect class="track" x="2" y="2" width="296" height="296" rx="4" />
              <rect
                class="elapsed"
                x="2"
                y="2"
                width="296"
                height="296"
                rx="4"
                pathLength="100"
                :stroke-dasharray="`${accepting ? refreshProgress : 0} 100`"
              />
            </svg>
            <div class="qr-paper">
              <QrCodeImage
                v-if="accepting && props.qr.value"
                :value="props.qr.value"
              />
              <div v-else class="qr-placeholder">
                <svg
                  v-if="accepting"
                  viewBox="0 0 64 64"
                  fill="none"
                  aria-hidden="true"
                >
                  <path
                    d="M22 8H8v14M42 8h14v14M8 42v14h14M56 42v14H42M25 25h14v14H25z"
                    stroke="currentColor"
                    stroke-width="2"
                  />
                </svg>
                <strong>{{
                  status === "pending"
                    ? "签到尚未开始"
                    : closed
                      ? "已停止签到"
                      : preview
                        ? "签到二维码"
                        : "二维码加载中"
                }}</strong>
                <small v-if="preview && accepting">演示占位 · 不可扫描</small>
              </div>
            </div>
          </div>
          <div
            class="refresh-status"
            role="progressbar"
            aria-label="二维码刷新剩余进度"
            :aria-valuenow="accepting ? Math.round(refreshProgress) : 0"
            aria-valuemin="0"
            aria-valuemax="100"
          >
            <span>{{ accepting ? "签到码自动更新" : "签到码暂不可用" }}</span>
            <strong v-if="accepting"
              >{{ refreshSeconds }}<small> 秒</small></strong
            >
          </div>
        </div>
        <div class="ticket-stub">
          <span class="ticket-flower" aria-hidden="true">✳</span>
          <div>
            <strong>{{
              closed ? "感谢每一位如约而至的你" : "打开手机扫一扫，留下此刻"
            }}</strong>
            <p>一场盛会，因你而完整</p>
          </div>
        </div>
      </section>
    </div>

    <footer class="colophon" data-enter>
      <span>GOOD COMPANY. GREAT MOMENTS.</span>
      <span class="footer-rule"></span>
      <span>相聚有时 · 精彩无尽</span>
    </footer>
  </main>
</template>

<style scoped>
.crimson-gala {
  position: relative;
  isolation: isolate;
  display: flex;
  flex-direction: column;
  height: 100dvh;
  overflow-x: hidden;
  overflow-y: auto;
  padding: 3.8vh 5vw 3vh;
  color: #f5dfbb;
  background: #200407;
  font-family: "Sora Variable", "Microsoft YaHei", sans-serif;
}
.masthead,
.stage,
.colophon {
  position: relative;
  z-index: 1;
}
.masthead {
  display: flex;
  align-items: center;
  gap: 28px;
}
.edition {
  padding-left: 26px;
  border-left: 1px solid #d4a36a35;
  color: #c6937f;
  font-size: clamp(9px, 0.65vw, 14px);
  letter-spacing: 0.2em;
}
.live-state {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-left: auto;
  color: #dab5a2;
  font-size: clamp(11px, 0.75vw, 16px);
}
.live-state i {
  width: 5px;
  height: 5px;
  border-radius: 50%;
  background: #957260;
}
.live-state i.active {
  background: #f4cc82;
  box-shadow: 0 0 12px #f4cc8270;
}
.stage {
  flex: 1;
  display: grid;
  grid-template-columns: 1.55fr 1fr;
  gap: 7vw;
  align-items: center;
  padding: 6vh 5vw 12vh;
}
.eyebrow {
  display: flex;
  align-items: center;
  gap: 18px;
  font-size: clamp(13px, 1vw, 22px);
  color: #ddb887;
  letter-spacing: 0.24em;
}
.eyebrow > span {
  width: 36px;
  height: 1px;
  background: #d8ad62;
}
h1 {
  position: relative;
  width: fit-content;
  margin: 3vh 0 4vh;
  font-family: "STXingkai", "华文行楷", "KaiTi", "楷体", serif;
  font-size: clamp(68px, 7.4vw, 175px);
  line-height: 1.18;
  font-weight: 400;
  letter-spacing: 0.05em;
  color: #f3d7a0;
  text-shadow: 0 4px 24px #26000480;
}
h1 > span {
  background: linear-gradient(
    125deg,
    #ad7136 0%,
    #f7e1b8 28%,
    #d7a15b 55%,
    #ffe9bd 78%,
    #c49250 100%
  );
  background-clip: text;
  -webkit-text-fill-color: transparent;
}
h1 > i {
  position: absolute;
  bottom: 0.1em;
  right: -0.6em;
  font-style: normal;
  font-size: 0.65em;
  color: #cb9360;
}
.dedication {
  display: flex;
  align-items: center;
  gap: 20px;
}
.seal {
  flex-shrink: 0;
  border: 1px solid #d4a36a70;
  padding: 6px;
  font-family: "KaiTi", serif;
  font-size: 15px;
  line-height: 1.2;
  color: #e5b67d;
}
.dedication p {
  margin: 0 0 8px;
  font-size: clamp(16px, 1.2vw, 25px);
  letter-spacing: 0.16em;
}
.dedication div > span {
  display: block;
  max-width: 36em;
  font-size: clamp(11px, 0.75vw, 16px);
  line-height: 1.9;
  color: #c19788;
}
.gathering {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 2.5vw;
  margin-top: 5vh;
  padding-top: 3vh;
  border-top: 1px solid #e6bb792b;
}
.label {
  margin: 0 0 10px;
  font-size: clamp(12px, 0.85vw, 18px);
  color: #dfb8a0;
}
.label > span {
  margin-left: 8px;
  font-size: 8px;
  letter-spacing: 0.12em;
  opacity: 0.5;
}
.count {
  display: flex;
  align-items: baseline;
  gap: 14px;
}
.count > strong {
  color: #edc990;
  font-size: clamp(38px, 3.4vw, 80px);
  font-weight: 400;
  line-height: 1.3;
  font-variant-numeric: tabular-nums;
}
.count > span {
  font-size: 12px;
  color: #c39782;
}
.attendance-line {
  height: 2px;
  margin-top: 12px;
  background: #e4bd8720;
}
.attendance-line > i {
  display: block;
  height: 100%;
  background: linear-gradient(90deg, #a76b34, #f3d7a0);
  transition: width 0.8s ease;
}
.caption {
  margin: 12px 0 0;
  font-size: clamp(10px, 0.65vw, 14px);
  color: #b7897b;
}
.deadline {
  border-left: 1px solid #e6bb792b;
  padding-left: 2.5vw;
}
.clock {
  display: flex;
  align-items: baseline;
  gap: 9px;
  padding-top: 8px;
  font-size: clamp(22px, 2.05vw, 48px);
  font-variant-numeric: tabular-nums;
  font-weight: 300;
}
.clock i {
  font-style: normal;
  color: #af785a;
}
.clock small {
  margin-left: 4px;
  font-size: 9px;
  color: #b7897b;
}
.clock > span {
  white-space: nowrap;
}
.admission {
  justify-self: center;
  width: min(100%, 380px);
  color: #593123;
  background: #f5ead6;
  box-shadow:
    0 25px 80px #18000266,
    0 0 70px #e8a04a0a;
  transform: rotate(2deg);
}
.ticket-top {
  display: flex;
  justify-content: space-between;
  padding: 18px 24px;
  background: #e3c79b;
  font-size: 10px;
  letter-spacing: 0.16em;
}
.ticket-body {
  padding: 28px 34px 22px;
  text-align: center;
  background: radial-gradient(ellipse at 50% 0%, #fffaf0, transparent 75%);
}
.ticket-kicker {
  font-size: 11px;
  color: #a18059;
  letter-spacing: 0.22em;
}
h2 {
  margin: 12px 0 28px;
  font-family: "Songti SC", "SimSun", serif;
  font-size: clamp(25px, 1.8vw, 38px);
  font-weight: 400;
  letter-spacing: 0.16em;
}
.qr-frame {
  position: relative;
  padding: 12px;
  aspect-ratio: 1;
}
.qr-progress {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
  fill: none;
  stroke-width: 1.5;
  pointer-events: none;
}
.track {
  stroke: #b68a4d30;
}
.elapsed {
  stroke: #a87839;
  transition: stroke-dasharray 0.25s linear;
}
.qr-paper {
  display: grid;
  place-items: center;
  width: 100%;
  height: 100%;
  background: #fff;
}
.qr-placeholder {
  display: flex;
  align-items: center;
  flex-direction: column;
  gap: 12px;
  color: #ad8955;
}
.qr-placeholder svg {
  width: 62px;
  margin-bottom: 8px;
}
.qr-placeholder strong {
  font-size: 14px;
  font-weight: 400;
}
.qr-placeholder small {
  font-size: 10px;
  color: #b2a594;
}
.refresh-status {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  margin-top: 16px;
  font-size: 11px;
  color: #997957;
}
.refresh-status strong {
  color: #704527;
  font-weight: 400;
  font-size: 20px;
  font-variant-numeric: tabular-nums;
}
.refresh-status small {
  font-size: 10px;
}
.ticket-stub {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 16px;
  padding: 22px 16px;
  border-top: 1px dashed #bc9a6b70;
  background: #eadbc1;
}
.ticket-flower {
  font-size: 32px;
  color: #a88044;
}
.ticket-stub strong {
  font-size: 12px;
  font-weight: 400;
}
.ticket-stub p {
  margin: 6px 0 0;
  font-size: 10px;
  letter-spacing: 0.16em;
  color: #a0805f;
}
.colophon {
  display: flex;
  align-items: center;
  gap: 24px;
  color: #c29276;
  font-size: 10px;
  letter-spacing: 0.2em;
}
.footer-rule {
  height: 1px;
  flex: 1;
  background: #e3bc7625;
}
@media (min-width: 2000px) {
  .admission {
    width: min(100%, 20vw);
  }
}
@media (max-width: 850px) {
  .crimson-gala {
    padding: 25px 24px;
  }
  .edition {
    display: none;
  }
  .stage {
    grid-template-columns: 1fr;
    padding: 55px 0 65px;
    gap: 45px;
  }
  .invitation {
    max-width: 550px;
    margin: auto;
  }
  h1 {
    font-size: clamp(62px, 13vw, 100px);
  }
  .gathering {
    gap: 20px;
  }
  .deadline {
    padding-left: 20px;
  }
  .clock {
    gap: 5px;
    font-size: 21px;
  }
  .clock small {
    display: block;
    margin-top: 5px;
    margin-left: 0;
  }
  .label > span {
    display: none;
  }
  .admission {
    width: min(100%, 340px);
    transform: none;
  }
  .colophon {
    font-size: 8px;
    gap: 12px;
    letter-spacing: 0.08em;
  }
}
</style>
