<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'
const canvas = ref<HTMLCanvasElement>()
let frame = 0
let observer: ResizeObserver | undefined
let width = 0,
  height = 0,
  previous = 0,
  elapsed = 0
const reduced = window.matchMedia('(prefers-reduced-motion: reduce)')
const particles = Array.from({ length: 65 }, (_, index) => ({
  x: Math.random(),
  y: Math.random(),
  size: 1 + Math.random() * 3,
  speed: 0.012 + Math.random() * 0.018,
  phase: index * 2.4,
}))
function draw(time: number) {
  const ctx = canvas.value?.getContext('2d')
  if (!ctx) return
  elapsed += previous ? Math.min((time - previous) / 1000, 0.05) : 0
  previous = time
  ctx.clearRect(0, 0, width, height)
  for (const p of particles) {
    const y = ((((p.y - elapsed * p.speed) % 1) + 1) % 1) * height
    const x = (p.x + Math.sin(elapsed * 0.25 + p.phase) * 0.015) * width
    ctx.save()
    ctx.translate(x, y)
    ctx.rotate(p.phase + elapsed * 0.15)
    ctx.globalAlpha = 0.18 + 0.4 * Math.pow(Math.sin(p.phase + elapsed * 0.4), 2)
    ctx.fillStyle = '#ffdaa0'
    ctx.fillRect(-p.size / 2, -p.size, p.size, p.size * 2)
    ctx.restore()
  }
  if (!reduced.matches && !document.hidden) frame = requestAnimationFrame(draw)
}
function resume() {
  cancelAnimationFrame(frame)
  previous = 0
  if (!document.hidden) frame = requestAnimationFrame(draw)
}
onMounted(() => {
  const el = canvas.value!
  observer = new ResizeObserver(() => {
    const bounds = el.getBoundingClientRect()
    width = bounds.width
    height = bounds.height
    const dpr = Math.min(window.devicePixelRatio, 1.5)
    el.width = width * dpr
    el.height = height * dpr
    el.getContext('2d')?.setTransform(dpr, 0, 0, dpr, 0, 0)
    resume()
  })
  observer.observe(el)
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
  <canvas
    ref="canvas"
    class="confetti"
    aria-hidden="true"
  />
</template>
<style scoped>
.confetti {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
}
</style>
