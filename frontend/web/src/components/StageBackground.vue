<template>
  <div ref="container" class="stage-background" aria-hidden="true"></div>
</template>

<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref } from 'vue'
import { gsap } from 'gsap'
import * as THREE from 'three'

const props = withDefaults(
  defineProps<{
    color?: string
  }>(),
  {
    color: '#ef4444',
  },
)

const container = ref<HTMLDivElement>()
let frameId = 0
let resizeObserver: ResizeObserver | undefined
let renderer: THREE.WebGLRenderer | undefined

onMounted(() => {
  if (!container.value) return

  const scene = new THREE.Scene()
  const camera = new THREE.PerspectiveCamera(45, 1, 0.1, 100)
  camera.position.z = 6

  renderer = new THREE.WebGLRenderer({ alpha: true, antialias: true })
  renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2))
  container.value.appendChild(renderer.domElement)

  const geometry = new THREE.BufferGeometry()
  const positions = new Float32Array(180 * 3)
  for (let index = 0; index < positions.length; index += 3) {
    positions[index] = (Math.random() - 0.5) * 12
    positions[index + 1] = (Math.random() - 0.5) * 7
    positions[index + 2] = (Math.random() - 0.5) * 4
  }
  geometry.setAttribute('position', new THREE.BufferAttribute(positions, 3))

  const material = new THREE.PointsMaterial({
    color: props.color,
    opacity: 0.32,
    size: 0.035,
    transparent: true,
  })
  const particles = new THREE.Points(geometry, material)
  scene.add(particles)

  const entrance = gsap.fromTo(
    material,
    { opacity: 0 },
    { duration: 1.2, opacity: 0.32, ease: 'power2.out' },
  )

  const resize = () => {
    if (!container.value || !renderer) return
    const { clientWidth, clientHeight } = container.value
    renderer.setSize(clientWidth, clientHeight, false)
    camera.aspect = clientWidth / Math.max(clientHeight, 1)
    camera.updateProjectionMatrix()
  }

  resizeObserver = new ResizeObserver(resize)
  resizeObserver.observe(container.value)
  resize()

  const render = () => {
    particles.rotation.y += 0.0006
    particles.rotation.x += 0.00015
    renderer?.render(scene, camera)
    frameId = window.requestAnimationFrame(render)
  }
  render()

  onBeforeUnmount(() => {
    window.cancelAnimationFrame(frameId)
    resizeObserver?.disconnect()
    entrance.kill()
    geometry.dispose()
    material.dispose()
    renderer?.dispose()
    renderer?.domElement.remove()
  })
})
</script>

<style scoped>
.stage-background {
  position: absolute;
  inset: 0;
  overflow: hidden;
  pointer-events: none;
}
</style>
