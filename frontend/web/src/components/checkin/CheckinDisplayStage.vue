<template>
  <section class="display-stage">
    <div class="display-stage__heading">
      <div>
        <span>CHECK-IN MOMENTS</span>
        <h2>每一次到场，都值得被看见</h2>
      </div>
      <p>动态照片墙 · 实时呈现</p>
    </div>

    <div class="display-stage__field">
      <div class="display-stage__orbit display-stage__orbit--one"></div>
      <div class="display-stage__orbit display-stage__orbit--two"></div>
      <div class="display-stage__platform"></div>
      <div
        v-for="avatar in avatars"
        :key="avatar.label"
        class="display-stage__avatar"
        :style="avatar.style"
      >
        {{ avatar.label }}
      </div>
      <div class="display-stage__waiting">
        <span class="display-stage__pulse"></span>
        等待新的签到人员进入
      </div>
    </div>

    <footer class="display-stage__footer">
      <span>新签到人员将在此处自然入场</span>
      <span>展示效果将在接入活动配置后切换</span>
    </footer>
  </section>
</template>

<script setup lang="ts">
const avatars = [
  { label: '林', style: { '--x': '11%', '--y': '25%', '--size': '94px', '--delay': '-2s', '--rotate': '-7deg' } },
  { label: '周', style: { '--x': '34%', '--y': '10%', '--size': '72px', '--delay': '-4s', '--rotate': '5deg' } },
  { label: '许', style: { '--x': '63%', '--y': '21%', '--size': '108px', '--delay': '-1s', '--rotate': '8deg' } },
  { label: '陈', style: { '--x': '80%', '--y': '50%', '--size': '78px', '--delay': '-5s', '--rotate': '-6deg' } },
  { label: '沈', style: { '--x': '47%', '--y': '55%', '--size': '92px', '--delay': '-3s', '--rotate': '-4deg' } },
  { label: '顾', style: { '--x': '18%', '--y': '64%', '--size': '70px', '--delay': '-6s', '--rotate': '7deg' } },
]
</script>

<style scoped>
.display-stage {
  display: flex;
  min-width: 0;
  height: 100%;
  padding: clamp(12px, 1vw, 20px) clamp(8px, 1vw, 18px) 0 clamp(24px, 2.1vw, 40px);
  flex-direction: column;
}

.display-stage__heading,
.display-stage__footer {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
}

.display-stage__heading span {
  color: var(--screen-accent);
  font-size: clamp(9px, 0.58vw, 12px);
  font-weight: 700;
  letter-spacing: 0.19em;
}

.display-stage__heading h2 {
  margin: 7px 0 0;
  color: #f0eadf;
  font-size: clamp(21px, 1.55vw, 32px);
  font-weight: 510;
}

.display-stage__heading p,
.display-stage__footer {
  margin: 0;
  color: rgba(255, 255, 255, 0.42);
  font-size: clamp(10px, 0.68vw, 13px);
}

.display-stage__field {
  position: relative;
  min-height: 0;
  margin: clamp(18px, 1.8vw, 34px) 0;
  flex: 1;
  overflow: hidden;
  background:
    radial-gradient(ellipse at 50% 83%, rgba(174, 137, 75, 0.1), transparent 32%),
    radial-gradient(circle at 50% 45%, rgba(255, 247, 229, 0.025), transparent 42%);
}

.display-stage__orbit {
  position: absolute;
  top: 50%;
  left: 50%;
  border: 1px solid rgba(201, 171, 115, 0.09);
  border-radius: 50%;
  transform: translate(-50%, -50%);
}

.display-stage__orbit--one {
  width: 55%;
  aspect-ratio: 1;
}

.display-stage__orbit--two {
  width: 82%;
  aspect-ratio: 1;
}

.display-stage__platform {
  position: absolute;
  bottom: 2%;
  left: 50%;
  width: 58%;
  height: 11%;
  background: linear-gradient(180deg, rgba(202, 170, 110, 0.13), rgba(55, 45, 29, 0.08));
  border: 1px solid rgba(211, 182, 129, 0.18);
  border-radius: 50%;
  box-shadow:
    0 0 34px rgba(157, 123, 65, 0.1),
    inset 0 0 28px rgba(221, 196, 150, 0.06);
  transform: translateX(-50%);
}

.display-stage__avatar {
  position: absolute;
  top: var(--y);
  left: var(--x);
  display: grid;
  width: var(--size);
  aspect-ratio: 1;
  place-items: center;
  color: #e8dfd0;
  font-size: calc(var(--size) * 0.3);
  font-weight: 560;
  background:
    linear-gradient(145deg, rgba(206, 174, 112, 0.1), rgba(255, 255, 255, 0.025)),
    rgba(14, 13, 11, 0.9);
  border: 1px solid rgba(194, 162, 103, 0.24);
  border-radius: 20px;
  box-shadow:
    0 20px 46px rgba(0, 0, 0, 0.34),
    inset 0 1px rgba(238, 219, 180, 0.1);
  transform: rotate(var(--rotate));
  animation: avatar-float 6s ease-in-out var(--delay) infinite alternate;
}

.display-stage__waiting {
  position: absolute;
  bottom: 22px;
  left: 50%;
  display: flex;
  align-items: center;
  gap: 9px;
  color: rgba(255, 255, 255, 0.46);
  font-size: clamp(10px, 0.7vw, 14px);
  transform: translateX(-50%);
  white-space: nowrap;
}

.display-stage__pulse {
  width: 7px;
  height: 7px;
  background: var(--screen-accent);
  border-radius: 50%;
  box-shadow: 0 0 0 0 rgba(199, 163, 101, 0.42);
  animation: pulse 1.8s infinite;
}

.display-stage__footer {
  align-items: center;
}

@keyframes avatar-float {
  to {
    transform: translate3d(0, -13px, 0) rotate(var(--rotate));
  }
}

@keyframes pulse {
  70% {
    box-shadow: 0 0 0 8px rgba(199, 163, 101, 0);
  }
}

@media (prefers-reduced-motion: reduce) {
  .display-stage__avatar,
  .display-stage__pulse {
    animation: none;
  }
}
</style>
