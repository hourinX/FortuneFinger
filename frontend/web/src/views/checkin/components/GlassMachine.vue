<script setup lang="ts">
import { onMounted, onUnmounted, ref, watch } from 'vue'
import * as THREE from 'three'
import * as CANNON from 'cannon-es'
const props = defineProps<{ count: number }>()
const host = ref<HTMLElement>()
const failed = ref(false)
const visible = ref(0)
// 保守限制可见容量，球体运动交由刚体世界；不设置不可见的半高反弹面。
const radius = 0.22,
  boundary = 2.05 - radius,
  capacity = 32
let renderer: THREE.WebGLRenderer | undefined
let observer: ResizeObserver | undefined
let frame = 0,
  last = 0,
  time = 0,
  nextEntry = 0
const scene = new THREE.Scene()
const camera = new THREE.PerspectiveCamera(35, 1, 0.1, 50)
camera.position.set(0, 1.1, 9.5)
camera.lookAt(0, -0.25, 0)
const world = new CANNON.World({ gravity: new CANNON.Vec3(0, -3.8, 0) })
world.defaultContactMaterial.friction = 0.18
world.defaultContactMaterial.restitution = 0.48
world.allowSleep = false
const force = new CANNON.Vec3()
// 球形容器使用向内的切平面围合，避免凹面三角网格的漏碰撞。
for (let row = 1; row < 20; row++) {
  const phi = (Math.PI * row) / 20
  for (let col = 0; col < 40; col++) {
    const theta = (Math.PI * 2 * col) / 40
    const n = new CANNON.Vec3(
      Math.sin(phi) * Math.cos(theta),
      Math.cos(phi),
      Math.sin(phi) * Math.sin(theta),
    )
    const wall = new CANNON.Body({ mass: 0, shape: new CANNON.Plane() })
    wall.position.copy(n.scale(2.05))
    wall.quaternion.setFromVectors(new CANNON.Vec3(0, 0, 1), n.negate())
    world.addBody(wall)
  }
}
for (const y of [-1, 1]) {
  const wall = new CANNON.Body({ mass: 0, shape: new CANNON.Plane() })
  wall.position.set(0, y * 2.05, 0)
  wall.quaternion.setFromVectors(new CANNON.Vec3(0, 0, 1), new CANNON.Vec3(0, -y, 0))
  world.addBody(wall)
}
const balls: { mesh: THREE.Group; body: CANNON.Body }[] = []
const resources: (THREE.BufferGeometry | THREE.Material)[] = []
const reduced = window.matchMedia('(prefers-reduced-motion: reduce)')
function geometry<T extends THREE.BufferGeometry>(item: T) {
  resources.push(item)
  return item
}
function material<T extends THREE.Material>(item: T) {
  resources.push(item)
  return item
}
const ballGeo = geometry(new THREE.SphereGeometry(radius, 24, 16))
const eyeGeo = geometry(new THREE.SphereGeometry(0.033, 8, 6))
const eyeMat = material(new THREE.MeshBasicMaterial({ color: '#26364e' }))
const colors = ['#ff9fbd', '#9edbff', '#ffe196', '#c4a9ff', '#92ead3']
const ballMats = colors.map((color) =>
  material(new THREE.MeshStandardMaterial({ color, roughness: 0.36, metalness: 0.05 })),
)
function addBall(initial = false) {
  if (balls.length >= capacity) return
  const group = new THREE.Group()
  group.add(new THREE.Mesh(ballGeo, ballMats[balls.length % colors.length]))
  for (const x of [-0.065, 0.065]) {
    const eye = new THREE.Mesh(eyeGeo, eyeMat)
    eye.position.set(x, 0.035, radius * 0.94)
    group.add(eye)
  }
  if (initial) {
    let valid = false
    for (let attempt = 0; attempt < 300 && !valid; attempt++) {
      group.position.set(
        (Math.random() - 0.5) * 3,
        -radius - Math.random() * 1.4,
        (Math.random() - 0.5) * 2.3,
      )
      valid =
        group.position.length() < boundary &&
        balls.every((b) => b.mesh.position.distanceTo(group.position) > radius * 2.05)
    }
  } else group.position.set((Math.random() - 0.5) * 0.25, 1.75, 0.15)
  group.userData.entered = initial
  scene.add(group)
  const body = new CANNON.Body({
    mass: 1,
    shape: new CANNON.Sphere(radius),
    linearDamping: 0.22,
    angularDamping: 0.3,
  })
  body.position.set(group.position.x, group.position.y, group.position.z)
  body.velocity.set((Math.random() - 0.5) * 0.5, -0.2, (Math.random() - 0.5) * 0.5)
  world.addBody(body)
  balls.push({ mesh: group, body })
  visible.value = balls.length
}
function syncCount() {
  while (balls.length < Math.min(capacity, props.count)) addBall(true)
}
watch(
  () => props.count,
  (count, old) => {
    if (count < old) {
      while (balls.length > Math.min(capacity, count)) {
        const ball = balls.pop()!
        scene.remove(ball.mesh)
        world.removeBody(ball.body)
      }
      visible.value = balls.length
    }
    // 新签到由渲染循环依次入球，人数先更新，展示队列不阻塞签到。
    if (reduced.matches) {
      syncCount()
      renderer?.render(scene, camera)
    }
  },
)
function render(stamp: number) {
  if (!renderer) return
  const dt = last ? Math.min((stamp - last) / 1000, 0.025) : 0.016
  last = stamp
  time += dt
  if (time >= nextEntry && balls.length < Math.min(capacity, props.count)) {
    addBall(reduced.matches)
    nextEntry = time + 0.65
  }
  if (!reduced.matches) {
    for (const [i, ball] of balls.entries()) {
      const p = ball.body.position
      // 底部偏心气流，离开喷口后自然回落，碰撞产生滚动与旋转。
      const jet = Math.max(0, 1 - Math.hypot(p.x + 0.45, p.z) / 0.9)
      force.set(-p.z * 0.65 + Math.sin(time + i) * 0.35, p.y < -1 ? jet * 9 : 0, p.x * 0.65)
      ball.body.applyForce(force)
    }
    world.step(1 / 60, dt, 3)
    for (const ball of balls) {
      ball.mesh.position.copy(ball.body.position)
      ball.mesh.quaternion.copy(ball.body.quaternion)
    }
  }

  renderer.render(scene, camera)
  if (!document.hidden && !reduced.matches) frame = requestAnimationFrame(render)
}
function resume() {
  cancelAnimationFrame(frame)
  last = 0
  if (reduced.matches) {
    syncCount()
    for (const ball of balls) {
      if (ball.mesh.position.y > -radius) ball.mesh.position.set(0, -radius, 0)
      ball.mesh.userData.entered = true
    }
  }
  if (!document.hidden) frame = requestAnimationFrame(render)
}
onMounted(() => {
  try {
    renderer = new THREE.WebGLRenderer({ alpha: true, antialias: true })
  } catch {
    failed.value = true
    return
  }
  renderer.setPixelRatio(Math.min(devicePixelRatio, 1.5))
  host.value!.appendChild(renderer.domElement)
  scene.add(new THREE.HemisphereLight('#f4faff', '#5b4590', 3))
  const light = new THREE.DirectionalLight('#ffffff', 4)
  light.position.set(-3, 5, 4)
  scene.add(light)
  const glass = new THREE.Mesh(
    geometry(new THREE.SphereGeometry(2.05, 64, 40)),
    material(
      new THREE.MeshPhongMaterial({
        color: '#c4f2ff',
        transparent: true,
        opacity: 0.11,
        shininess: 130,
        side: THREE.FrontSide,
        depthWrite: false,
      }),
    ),
  )
  glass.renderOrder = 2
  scene.add(glass)
  const ring = new THREE.Mesh(
    geometry(new THREE.TorusGeometry(2.06, 0.018, 8, 100)),
    material(new THREE.MeshBasicMaterial({ color: '#b4d9ec', transparent: true, opacity: 0.45 })),
  )
  scene.add(ring)
  const base = new THREE.Mesh(
    geometry(new THREE.CylinderGeometry(1.35, 1.55, 0.32, 64)),
    material(
      new THREE.MeshStandardMaterial({ color: '#5b6586', metalness: 0.65, roughness: 0.23 }),
    ),
  )
  base.position.y = -2.15
  scene.add(base)
  const arc = new THREE.Mesh(
    geometry(new THREE.TorusGeometry(1.98, 0.027, 8, 70, 0.9)),
    material(new THREE.MeshBasicMaterial({ color: '#ffffff', transparent: true, opacity: 0.8 })),
  )
  arc.rotation.z = 1.75
  arc.position.z = 0.35
  scene.add(arc)
  const collar = new THREE.Mesh(
    geometry(new THREE.TorusGeometry(1.32, 0.04, 12, 80)),
    material(new THREE.MeshBasicMaterial({ color: '#adcfff' })),
  )
  collar.rotation.x = Math.PI / 2
  collar.position.y = -1.98
  scene.add(collar)
  const foot = new THREE.Mesh(
    geometry(new THREE.CylinderGeometry(1.62, 1.8, 0.12, 64)),
    material(new THREE.MeshStandardMaterial({ color: '#171e3b', metalness: 0.8, roughness: 0.22 })),
  )
  foot.position.y = -2.38
  scene.add(foot)
  const rimLight = new THREE.PointLight('#a49bff', 20, 10)
  rimLight.position.set(2, 1, -2)
  scene.add(rimLight)
  syncCount()
  observer = new ResizeObserver(() => {
    const box = host.value!.getBoundingClientRect()
    renderer!.setSize(box.width, box.height)
    camera.aspect = box.width / Math.max(1, box.height)
    camera.updateProjectionMatrix()
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
  world.bodies.slice().forEach((body) => world.removeBody(body))
  resources.forEach((r) => r.dispose())
  renderer?.dispose()
  renderer?.domElement.remove()
})
</script>
<template>
  <div class="machine">
    <div
      ref="host"
      class="canvas"
    ></div>
    <p v-if="failed">当前设备不支持 3D 球仓</p>
    <span class="caption">
      {{ visible >= capacity ? '球仓展示已满 · 签到仍持续记录' : `球仓内 ${visible} 颗幸运星` }}
    </span>
  </div>
</template>
<style scoped>
.machine {
  position: relative;
  width: 100%;
  height: 100%;
  min-height: 340px;
}
.canvas {
  width: 100%;
  height: 100%;
  min-height: 340px;
}
.caption {
  position: absolute;
  bottom: 4%;
  left: 0;
  width: 100%;
  text-align: center;
  color: #91a3bb;
  font-size: 11px;
}
.machine p {
  position: absolute;
  inset: 40% 0;
  text-align: center;
}
</style>
