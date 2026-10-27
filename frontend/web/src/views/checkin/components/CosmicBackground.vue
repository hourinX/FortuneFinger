<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'
const host = ref<HTMLCanvasElement>()
const reduced = window.matchMedia('(prefers-reduced-motion: reduce)')
let frame = 0,
  previous = 0,
  time = 0,
  width = 0,
  height = 0
let observer: ResizeObserver | undefined
// 星云预绘制到离屏纹理；运行时只旋转合成，不逐帧生成噪声。
const galaxy = document.createElement('canvas')
galaxy.width = galaxy.height = 1400
const dust = galaxy.getContext('2d')!
let seed = 7341
function random() {
  seed = (seed * 16807) % 2147483647
  return seed / 2147483647
}
for (let i = 0; i < 16000; i++) {
  const r = Math.pow(random(), 0.65) * 640
  const arm = ((i % 4) * Math.PI) / 2
  const angle = arm + r * 0.009 + (random() - 0.5) * (0.28 + r / 1100)
  const x = 700 + Math.cos(angle) * r,
    y = 700 + Math.sin(angle) * r
  const size = 2 + random() * 17
  const glow = dust.createRadialGradient(x, y, 0, x, y, size)
  const color = i % 3 === 0 ? '158,116,215' : '94,153,215'
  glow.addColorStop(0, `rgba(${color},.055)`)
  glow.addColorStop(1, `rgba(${color},0)`)
  dust.fillStyle = glow
  dust.fillRect(x - size, y - size, size * 2, size * 2)
  if (i % 5 === 0) {
    dust.fillStyle = `rgba(216,224,255,${0.15 + random() * 0.65})`
    dust.beginPath()
    dust.arc(x, y, 0.3 + random() * 1.1, 0, Math.PI * 2)
    dust.fill()
  }
}
const core = dust.createRadialGradient(700, 700, 0, 700, 700, 150)
core.addColorStop(0, '#fff0d9d9')
core.addColorStop(0.12, '#e1d5ff88')
core.addColorStop(0.4, '#a6a0ec35')
core.addColorStop(1, '#7975cf00')
dust.fillStyle = core
dust.fillRect(550, 550, 300, 300)
const stars = Array.from({ length: 380 }, (_, i) => ({
  angle: i * 2.39996,
  radius: Math.sqrt((i + 0.5) / 380),
  depth: 0.25 + ((i * 0.618) % 1) * 0.75,
  phase: i * 1.71,
}))
function draw(stamp: number) {
  const ctx = host.value?.getContext('2d')
  if (!ctx) return
  if (previous) time += Math.min((stamp - previous) / 1000, 0.05)
  previous = stamp
  ctx.clearRect(0, 0, width, height)
  ctx.save()
  ctx.translate(width * 0.56, height * 0.55)
  ctx.rotate(-0.34)
  ctx.scale(1, 0.5)
  ctx.rotate(time * 0.008)
  const extent = Math.max(width * 1.05, height * 1.4)
  ctx.globalAlpha = 0.85
  ctx.drawImage(galaxy, -extent / 2, -extent / 2, extent, extent)
  ctx.restore()
  // 慢速倾斜星盘：不同深度的星点以不同速度公转。
  ctx.save()
  ctx.translate(width * 0.5, height * 0.55)
  ctx.rotate(-0.25)
  for (const star of stars) {
    const a = star.angle + time * 0.012 * star.depth
    const r = star.radius * Math.max(width, height) * 0.77
    const x = Math.cos(a) * r,
      y = Math.sin(a) * r * 0.48
    const size = 0.5 + star.depth * 1.1
    ctx.globalAlpha = 0.2 + star.depth * 0.45 + Math.sin(time * 0.5 + star.phase) * 0.12
    ctx.fillStyle = star.depth > 0.8 ? '#d5d9ff' : '#789cc9'
    ctx.beginPath()
    ctx.arc(x, y, size, 0, Math.PI * 2)
    ctx.fill()
    if (star.depth > 0.96) {
      ctx.globalAlpha *= 0.3
      ctx.fillRect(x - 4, y - 0.4, 8, 0.8)
      ctx.fillRect(x - 0.4, y - 4, 0.8, 8)
    }
  }
  ctx.restore()
  if (!document.hidden && !reduced.matches) frame = requestAnimationFrame(draw)
}
function resume() {
  cancelAnimationFrame(frame)
  previous = 0
  if (!document.hidden) frame = requestAnimationFrame(draw)
}
onMounted(() => {
  observer = new ResizeObserver(() => {
    const el = host.value!,
      box = el.getBoundingClientRect(),
      dpr = Math.min(devicePixelRatio, 1.5)
    width = box.width
    height = box.height
    el.width = width * dpr
    el.height = height * dpr
    el.getContext('2d')?.setTransform(dpr, 0, 0, dpr, 0, 0)
    resume()
  })
  observer.observe(host.value!)
  document.addEventListener('visibilitychange', resume)
  reduced.addEventListener('change', resume)
})
onUnmounted(() => {
  cancelAnimationFrame(frame)
  observer?.disconnect()
  document.removeEventListener('visibilitychange', resume)
  reduced.removeEventListener('change', resume)
})
</script>
<template>
  <div
    class="cosmos"
    aria-hidden="true"
  >
    <div class="nebula"></div>
    <div class="planet"></div>
    <div class="sun"></div>
    <canvas ref="host"></canvas>
  </div>
</template>
<style scoped>
.cosmos {
  position: absolute;
  inset: 0;
  pointer-events: none;
  overflow: hidden;
  z-index: -1;
  background: #060a17;
}
.cosmos canvas {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
}
.nebula {
  position: absolute;
  inset: -15%;
  background:
    radial-gradient(ellipse 35% 14% at 40% 55%, #5663b432, transparent 85%),
    radial-gradient(ellipse 28% 20% at 65% 48%, #683cae26, transparent 85%),
    radial-gradient(ellipse 55% 24% at 50% 55%, #23537b26, transparent 80%);
  transform: rotate(-17deg);
}
.planet {
  position: absolute;
  right: -8%;
  top: 5%;
  width: 31vw;
  aspect-ratio: 1;
  border-radius: 50%;
  background: radial-gradient(circle at 22% 30%, #2e4860, #111b31 45%, #050917 70%);
  box-shadow:
    inset 3px 2px 3px #9cbedb42,
    inset 12px 7px 25px #729ddd12;
  opacity: 0.6;
}

.nebula {
  opacity: 0.8;
  animation: nebula-breathe 34s ease-in-out infinite alternate;
  background:
    radial-gradient(ellipse 30% 20% at 30% 65%, #56309560, transparent 80%),
    radial-gradient(ellipse 25% 18% at 75% 35%, #246b9955, transparent 85%),
    radial-gradient(ellipse 40% 10% at 50% 55%, #7a629f30, transparent 80%);
}
.sun {
  position: absolute;
  right: 15%;
  top: 20%;
  width: 5px;
  height: 5px;
  border-radius: 50%;
  background: #e9f1ff;
  box-shadow:
    0 0 8px 2px #fff,
    0 0 35px 10px #9ebcff88,
    0 0 100px 35px #668ae52a;
}
.sun::after {
  content: '';
  position: absolute;
  left: -50px;
  top: 2px;
  width: 105px;
  height: 1px;
  background: linear-gradient(90deg, transparent, #c4daff, transparent);
}
@keyframes nebula-breathe {
  to {
    transform: rotate(-9deg) scale(1.08);
    opacity: 0.55;
  }
}
@media (prefers-reduced-motion: reduce) {
  .nebula {
    animation: none;
  }
}
</style>
