<script setup lang="ts">
const letterModules = import.meta.glob<string>(
  '/src/assets/branding/loading-wordmark/*.png',
  {
    eager: true,
    import: 'default',
    query: '?url',
  },
)
const letterSlices = Object.entries(letterModules)
  .sort(([left], [right]) => left.localeCompare(right))
  .map(([, source]) => source)
</script>

<template>
  <div
    class="loading-screen"
    role="status"
    aria-live="polite"
    aria-label="签到大屏正在载入"
  >
    <div class="loading-screen__content">
      <div
        class="loading-screen__brand"
        aria-hidden="true"
      >
        <img
          v-for="(source, index) in letterSlices"
          :key="source"
          :src="source"
          alt=""
          :style="{ '--letter-index': index }"
        />
      </div>
      <div
        class="loading-screen__line"
        aria-hidden="true"
      >
        <i></i>
      </div>
      <p>
        活动现场正在准备
        <span
          class="loading-screen__dots"
          aria-hidden="true"
        >
          <i></i>
          <i></i>
          <i></i>
        </span>
      </p>
    </div>
  </div>
</template>

<style scoped>
.loading-screen {
  position: fixed;
  z-index: 1000;
  inset: 0;
  display: grid;
  place-items: center;
  overflow: hidden;
  color: #292733;
  background:
    radial-gradient(circle at 50% 44%, rgb(239 63 58 / 6%), transparent 27%),
    #f7f6f3;
}

.loading-screen::before,
.loading-screen::after {
  position: absolute;
  width: min(46vw, 720px);
  aspect-ratio: 1;
  border: 1px solid rgb(239 63 58 / 6%);
  border-radius: 50%;
  content: '';
}

.loading-screen::before {
  transform: scale(0.72);
}

.loading-screen::after {
  transform: scale(1.08);
}

.loading-screen__content {
  position: relative;
  z-index: 1;
  display: flex;
  align-items: center;
  flex-direction: column;
}

.loading-screen__brand {
  display: flex;
  width: clamp(250px, 28vw, 500px);
  aspect-ratio: 1864 / 294;
  align-items: stretch;
  overflow: visible;
}

.loading-screen__brand img {
  display: block;
  width: auto;
  min-width: 0;
  height: 100%;
  opacity: 0;
  transform: translateY(115%) scale(0.82);
  animation: letter-arrive 0.62s cubic-bezier(0.2, 0.9, 0.28, 1.2) forwards;
  animation-delay: calc(var(--letter-index) * 65ms + 80ms);
}

.loading-screen__line {
  width: min(230px, 48vw);
  height: 1px;
  margin-top: 24px;
  overflow: hidden;
  background: rgb(44 40 51 / 10%);
}

.loading-screen__line i {
  display: block;
  width: 42%;
  height: 100%;
  background: linear-gradient(90deg, transparent, #ef3f3a, transparent);
  animation: line-travel 1.6s ease-in-out infinite 0.9s;
}

.loading-screen p {
  display: flex;
  align-items: center;
  margin: 14px 0 0;
  color: #9c999d;
  font-size: 10px;
  letter-spacing: 0.42em;
  opacity: 0;
  animation: caption-arrive 0.5s ease forwards 0.95s;
}

.loading-screen__dots {
  display: inline-flex;
  width: 24px;
  align-items: center;
  justify-content: flex-start;
  gap: 3px;
  margin-left: 2px;
}

.loading-screen__dots i {
  width: 3px;
  height: 3px;
  border-radius: 50%;
  background: #ef3f3a;
  opacity: 0.28;
  animation: dot-bounce 1.2s ease-in-out infinite;
}

.loading-screen__dots i:nth-child(2) {
  animation-delay: 0.16s;
}

.loading-screen__dots i:nth-child(3) {
  animation-delay: 0.32s;
}

@keyframes letter-arrive {
  0% {
    opacity: 0;
    transform: translateY(115%) scale(0.82);
  }

  68% {
    opacity: 1;
    transform: translateY(-9%) scale(1.04);
  }

  100% {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

@keyframes line-travel {
  from {
    transform: translateX(-130%);
  }

  to {
    transform: translateX(340%);
  }
}

@keyframes caption-arrive {
  to {
    opacity: 1;
  }
}

@keyframes dot-bounce {
  0%,
  60%,
  100% {
    opacity: 0.28;
    transform: translateY(0);
  }

  30% {
    opacity: 1;
    transform: translateY(-4px);
  }
}

@media (prefers-reduced-motion: reduce) {
  .loading-screen__brand img,
  .loading-screen p {
    opacity: 1;
    transform: none;
    animation: none;
  }

  .loading-screen__line i {
    width: 100%;
    animation: none;
  }

  .loading-screen__dots i {
    opacity: 0.65;
    animation: none;
  }
}
</style>
