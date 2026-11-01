<script setup lang="ts">
import { onMounted, onBeforeUnmount, ref } from "vue";
import * as THREE from "three";

const host = ref<HTMLDivElement>();
let renderer: THREE.WebGLRenderer | undefined;
let observer: ResizeObserver | undefined;
let frame = 0;
let disposeScene: (() => void) | undefined;
let resume: (() => void) | undefined;

onMounted(() => {
  if (!host.value) return;
  const reduced = matchMedia("(prefers-reduced-motion: reduce)").matches;
  const low = innerWidth < 850 || navigator.hardwareConcurrency <= 4;
  try {
    renderer = new THREE.WebGLRenderer({
      alpha: true,
      antialias: !low,
      powerPreference: "low-power",
    });
  } catch {
    return; // CSS retains the stage lighting when WebGL is unavailable.
  }
  const canvasRenderer = renderer;
  canvasRenderer.setPixelRatio(Math.min(devicePixelRatio, low ? 1 : 1.5));
  canvasRenderer.setClearColor(0x000000, 0);
  host.value.appendChild(canvasRenderer.domElement);
  const scene = new THREE.Scene();
  const camera = new THREE.OrthographicCamera(-1, 1, 1, -1, 0.1, 10);
  camera.position.z = 2;
  const uniforms = { uTime: { value: 0 } };
  const geometry = new THREE.PlaneGeometry(2, 2);
  // Three overlapping silk surfaces. Each fold has a broad diffuse face,
  // a narrower grazing reflection, and a dark reverse side.
  const material = new THREE.ShaderMaterial({
    transparent: true,
    depthWrite: false,
    depthTest: false,
    uniforms,
    vertexShader: `
      varying vec2 vUv;
      void main() {
        vUv = uv;
        gl_Position = projectionMatrix * modelViewMatrix * vec4(position, 1.0);
      }
    `,
    fragmentShader: `
      varying vec2 vUv;
      uniform float uTime;
      void main() {
        vec3 color = vec3(0.0);
        float alpha = 0.0;
        for (int i = 0; i < 3; i++) {
          float layer = float(i);
          float x = vUv.x;
          float phase = layer * 1.9;
          float bend = sin(x * 5.4 + phase);
          float path = 0.07 + layer * 0.025 + bend * 0.035
            + 0.13 * pow(abs(x - 0.52) * 1.8, 2.0);
          float width = 0.025 + 0.032 * (0.5 + 0.5 * sin(x * 7.0 + phase));
          float across = (vUv.y - path) / width;
          float face = 1.0 - smoothstep(0.72, 1.0, abs(across));
          float fold = sin(x * 7.0 + phase + across * 1.6);
          float reflection = pow(max(0.0, fold), 9.0);
          float flow = 0.72 + 0.28 * sin(x * 9.0 - uTime * 0.23 + phase);
          vec3 silk = mix(vec3(0.16, 0.012, 0.02), vec3(0.57, 0.065, 0.045), fold * 0.5 + 0.5);
          silk += vec3(0.58, 0.31, 0.12) * reflection * flow;
          float edge = exp(-pow((across - 0.88) * 55.0, 2.0));
          float glint = pow(0.5 + 0.5 * sin(x * 12.0 + phase - uTime * 0.15), 12.0);
          silk += vec3(0.9, 0.65, 0.33) * edge * (0.12 + glint * 0.8);
          float opacity = face * (0.55 + layer * 0.12);
          color = mix(color, silk, opacity);
          alpha = opacity + alpha * (1.0 - opacity);
          float halo = exp(-pow((across - 0.85) * 5.0, 2.0)) * glint * 0.04;
          color += vec3(0.9, 0.5, 0.2) * halo;
          alpha = max(alpha, halo);
        }
        gl_FragColor = vec4(color, alpha);
      }
    `,
  });
  scene.add(new THREE.Mesh(geometry, material));

  const amount = low ? 45 : 100;
  const positions = new Float32Array(amount * 3);
  const seeds = new Float32Array(amount);
  for (let i = 0; i < amount; i++) {
    // Stable placement makes previews consistent and avoids full-screen noise.
    const seed = ((i * 73 + 19) % 101) / 101;
    positions[i * 3] = ((i * 0.618034) % 1) * 2 - 1;
    positions[i * 3 + 1] =
      i < amount * 0.65 ? -0.98 + seed * 0.48 : seed * 2 - 1;
    seeds[i] = seed;
  }
  const particles = new THREE.BufferGeometry();
  particles.setAttribute("position", new THREE.BufferAttribute(positions, 3));
  particles.setAttribute("aSeed", new THREE.BufferAttribute(seeds, 1));
  const pointMaterial = new THREE.ShaderMaterial({
    uniforms: { ...uniforms, uDpr: { value: canvasRenderer.getPixelRatio() } },
    transparent: true,
    depthWrite: false,
    depthTest: false,
    blending: THREE.AdditiveBlending,
    vertexShader: `
      attribute float aSeed;
      uniform float uTime;
      uniform float uDpr;
      varying float vLight;
      void main() {
        vec3 p = position;
        p.x += sin(uTime * 0.06 + aSeed * 24.0) * 0.018;
        p.y += sin(uTime * 0.09 + aSeed * 17.0) * 0.025;
        vLight = (0.25 + aSeed * 0.55) * (0.65 + 0.35 * sin(uTime * 0.5 + aSeed * 30.0));
        gl_PointSize = (1.3 + pow(aSeed, 5.0) * 5.0) * uDpr;
        gl_Position = projectionMatrix * modelViewMatrix * vec4(p, 1.0);
      }
    `,
    fragmentShader: `
      varying float vLight;
      void main() {
        float d = length(gl_PointCoord - 0.5) * 2.0;
        float glow = exp(-d * d * 6.0) * (1.0 - smoothstep(0.6, 1.0, d));
        gl_FragColor = vec4(1.0, 0.72, 0.36, glow * vLight);
      }
    `,
  });
  scene.add(new THREE.Points(particles, pointMaterial));
  let last = performance.now();
  function render(now: number) {
    const elapsed = Math.min((now - last) / 1000, 0.05);
    last = now;
    uniforms.uTime.value += elapsed;
    canvasRenderer.render(scene, camera);
    if (!reduced && !document.hidden) frame = requestAnimationFrame(render);
  }
  resume = () => {
    cancelAnimationFrame(frame);
    if (document.hidden) return;
    last = performance.now();
    frame = requestAnimationFrame(render);
  };
  observer = new ResizeObserver(([entry]) => {
    if (!entry) return;
    canvasRenderer.setSize(entry.contentRect.width, entry.contentRect.height);
    canvasRenderer.render(scene, camera);
  });
  observer.observe(host.value);
  document.addEventListener("visibilitychange", resume);
  resume();
  disposeScene = () => {
    geometry.dispose();
    material.dispose();
    particles.dispose();
    pointMaterial.dispose();
  };
});

onBeforeUnmount(() => {
  cancelAnimationFrame(frame);
  observer?.disconnect();
  if (resume) document.removeEventListener("visibilitychange", resume);
  disposeScene?.();
  renderer?.dispose();
  renderer?.domElement.remove();
});
</script>

<template>
  <div ref="host" class="crimson-atmosphere" aria-hidden="true">
    <div class="spotlight"></div>
  </div>
</template>

<style scoped>
.crimson-atmosphere {
  position: absolute;
  inset: 0;
  overflow: hidden;
  pointer-events: none;
  background:
    radial-gradient(ellipse 60% 30% at 62% 106%, #be391b65, transparent 85%),
    radial-gradient(ellipse 55% 80% at 58% 42%, #8d161e70, transparent 95%),
    linear-gradient(115deg, #210307, #52090d 50%, #250407);
}
.crimson-atmosphere :deep(canvas) {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
}
.spotlight {
  position: absolute;
  top: -35%;
  right: 3%;
  width: 65%;
  height: 140%;
  background: conic-gradient(
    from 157deg at 65% 0%,
    transparent 0deg,
    #efb77607 6deg,
    #efb77613 11deg,
    transparent 17deg,
    #efb77608 23deg,
    transparent 28deg
  );
  filter: blur(24px);
  transform: rotate(-14deg);
}
</style>
