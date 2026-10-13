<template>
  <a-modal
    :open="open"
    :footer="null"
    :closable="!loading"
    :mask-closable="false"
    width="372px"
    centered
    :title="title"
    @cancel="close"
  >
    <div class="enhanced-verification">
      <div class="verification-tip">
        <span>{{ verificationTip }}</span>
        <button type="button" :disabled="loading" @click="load">换一张</button>
      </div>
      <div v-if="isClickMode" class="verification-prompt">
        <span>{{ verificationData.mode === 'image-select' ? '选择顺序' : '点击顺序' }}</span>
        <img :src="verificationData.prompt" alt="需要依次点击的文字" draggable="false" />
      </div>
	  <div
		v-if="verificationData?.mode === 'behavior'"
		class="behavior-pad"
		:class="{ active: behaviorTracking }"
		@pointerdown="startBehavior"
		@pointermove="recordBehavior"
		@pointerup="finishBehavior"
		@pointercancel="cancelBehavior"
	  >
		<span>{{ behaviorTracking ? '请继续移动' : '按住并自然移动指针' }}</span>
	  </div>
      <a-spin :spinning="loading">
        <div
		  v-if="verificationData && verificationData.mode !== 'behavior'"
          ref="canvasRef"
          class="verification-canvas"
		  :class="`mode-${verificationData.mode}`"
          :style="canvasStyle"
          @click="selectPoint"
        >
          <img
            class="verification-image"
            :src="verificationData.image"
            alt="登录安全校验图片"
            draggable="false"
            @load="onImageLoad"
          />
		  <template v-if="isClickMode">
			<span
            v-for="(point, index) in points"
            :key="`${point.x}-${point.y}-${index}`"
            class="verification-marker"
            :style="markerStyle(point)"
          >
            {{ index + 1 }}
			</span>
		  </template>
		  <img
			v-if="verificationData.mode === 'slide'"
			class="verification-thumb"
			:class="{ dragging }"
			:src="verificationData.thumb"
			:style="thumbStyle"
			alt="可拖动拼图"
			draggable="false"
			@pointerdown="startDrag"
			@pointermove="moveDrag"
			@pointerup="finishDrag"
			@pointercancel="cancelDrag"
		  />
		  <img
			v-if="verificationData.mode === 'rotate'"
			class="verification-rotate-thumb"
			:src="verificationData.thumb"
			:style="{ transform: `translate(-50%, -50%) rotate(${rotation}deg)` }"
			alt="待旋转图片"
			draggable="false"
		  />
        </div>
      </a-spin>
	  <input
		v-if="verificationData?.mode === 'rotate'"
		v-model.number="rotation"
		class="rotation-control"
		type="range"
		min="0"
		max="360"
		@change="finishRotate"
	  />
      <div class="verification-footer">
		<span v-if="isClickMode">
		  已选择 {{ points.length }}/{{ requiredPoints }} 个位置
		</span>
		<span v-else-if="verificationData?.mode === 'slide'">按住拼图块左右拖动，松开后自动校验</span>
		<span v-else-if="verificationData?.mode === 'rotate'">旋转图片至正确方向，松开后自动校验</span>
		<span v-else>系统会根据移动轨迹判断是否为真实操作</span>
        <button v-if="isClickMode && points.length" type="button" :disabled="loading" @click="resetPoints">
          重新选择
        </button>
      </div>
    </div>
  </a-modal>
</template>

<script setup>
import { computed, reactive, ref, watch } from 'vue'
import { message } from 'ant-design-vue'
import { createEnhancedVerification, verifyEnhancedVerification } from '@/api/auth'

const props = defineProps({
  open: Boolean,
	scene: {
		type: String,
		default: 'login',
	},
  subject: {
    type: String,
    default: '',
  },
	title: {
		type: String,
		default: '登录安全校验',
	},
	purpose: {
		type: String,
		default: '为了保障登录安全',
	},
})
const emit = defineEmits(['close', 'verified'])
const requiredPoints = ref(0)
const loading = ref(false)
const dragging = ref(false)
const behaviorTracking = ref(false)
const canvasRef = ref()
const verificationID = ref('')
const verificationData = ref(null)
const imageSize = reactive({ width: 300, height: 220 })
const points = ref([])
const rotation = ref(0)
const slidePosition = reactive({ x: 0, y: 0 })
const dragStart = reactive({ clientX: 0, x: 0 })
const canvasStyle = computed(() => ({
  aspectRatio: `${imageSize.width} / ${imageSize.height}`,
}))
const verificationTip = computed(() =>
	`${props.purpose}${
		({
			'click-text': '，请依次点击下图中的文字',
			'click-icon': '，请依次点击下图中的图标',
			'image-select': '，请依次选择指定图片',
			rotate: '，请将图片旋转至正确方向',
			behavior: '，请完成一次自然的指针移动',
			slide: '，请将拼图拖到缺口位置',
		}[verificationData.value?.mode] || '，请完成校验')
	}`,
)
const isClickMode = computed(() =>
	['click-text', 'click-icon', 'image-select'].includes(verificationData.value?.mode),
)
const thumbStyle = computed(() => ({
  width: `${(verificationData.value?.thumbWidth / imageSize.width) * 100}%`,
  height: `${(verificationData.value?.thumbHeight / imageSize.height) * 100}%`,
  left: `${(slidePosition.x / imageSize.width) * 100}%`,
  top: `${(slidePosition.y / imageSize.height) * 100}%`,
}))

function imageSource(value, mime) {
  if (!value || value.startsWith('data:')) return value
  return `data:${mime};base64,${value}`
}

async function load() {
  if (!props.subject) return
  loading.value = true
	dragging.value = false
  points.value = []
  verificationData.value = null
  try {
	const data = await createEnhancedVerification({
		scene: props.scene,
		subject: props.subject,
	})
    verificationID.value = data.verification_id
    verificationData.value = {
	  mode: data.mode,
      image: imageSource(data.image, 'image/jpeg'),
      prompt: imageSource(data.prompt, 'image/png'),
	  thumb: imageSource(data.thumb, 'image/png'),
	  thumbX: data.thumb_x,
	  thumbY: data.thumb_y,
	  thumbWidth: data.thumb_width,
	  thumbHeight: data.thumb_height,
    }
	requiredPoints.value = data.required_points
	slidePosition.x = data.thumb_x
	slidePosition.y = data.thumb_y
	rotation.value = 0
  } catch (error) {
    message.error(error.message || '登录安全校验加载失败')
    close()
  } finally {
    loading.value = false
  }
}

function onImageLoad(event) {
  imageSize.width = event.target.naturalWidth || 300
  imageSize.height = event.target.naturalHeight || 220
}

async function selectPoint(event) {
	if (!isClickMode.value) return
	if (loading.value || points.value.length >= requiredPoints.value || !canvasRef.value) return
  const rect = canvasRef.value.getBoundingClientRect()
  const point = {
    x: Math.round((event.clientX - rect.left) * (imageSize.width / rect.width)),
    y: Math.round((event.clientY - rect.top) * (imageSize.height / rect.height)),
  }
  points.value.push(point)
  if (points.value.length === requiredPoints.value) await verify()
}

async function finishRotate() {
	points.value = [{ x: Math.round(rotation.value), y: 0 }]
	await verify()
}

function startBehavior(event) {
	if (loading.value) return
	behaviorTracking.value = true
	points.value = []
	event.currentTarget.setPointerCapture(event.pointerId)
	recordBehavior(event)
}

function recordBehavior(event) {
	if (!behaviorTracking.value) return
	const rect = event.currentTarget.getBoundingClientRect()
	points.value.push({
		x: Math.max(0, Math.round(event.clientX - rect.left)),
		y: Math.max(0, Math.round(event.clientY - rect.top)),
		t: Date.now(),
	})
}

async function finishBehavior(event) {
	if (!behaviorTracking.value) return
	recordBehavior(event)
	behaviorTracking.value = false
	event.currentTarget.releasePointerCapture(event.pointerId)
	await verify()
}

function cancelBehavior() {
	behaviorTracking.value = false
	points.value = []
}

function startDrag(event) {
	if (loading.value) return
	dragging.value = true
	dragStart.clientX = event.clientX
	dragStart.x = slidePosition.x
	event.currentTarget.setPointerCapture(event.pointerId)
}

function moveDrag(event) {
	if (!dragging.value || !canvasRef.value) return
	const rect = canvasRef.value.getBoundingClientRect()
	const scale = imageSize.width / rect.width
	const maxX = imageSize.width - verificationData.value.thumbWidth
	slidePosition.x = Math.min(
		maxX,
		Math.max(0, dragStart.x + (event.clientX - dragStart.clientX) * scale),
	)
}

async function finishDrag(event) {
	if (!dragging.value) return
	dragging.value = false
	behaviorTracking.value = false
	event.currentTarget.releasePointerCapture(event.pointerId)
	points.value = [{ x: Math.round(slidePosition.x), y: Math.round(slidePosition.y) }]
	await verify()
}

function cancelDrag() {
	dragging.value = false
	slidePosition.x = verificationData.value?.thumbX || 0
}

function markerStyle(point) {
  return {
    left: `${(point.x / imageSize.width) * 100}%`,
    top: `${(point.y / imageSize.height) * 100}%`,
  }
}

function resetPoints() {
  points.value = []
}

async function verify() {
  loading.value = true
  try {
    const data = await verifyEnhancedVerification({
      verification_id: verificationID.value,
	  scene: props.scene,
      subject: props.subject,
      points: points.value,
    })
    emit('verified', data.verification_token)
  } catch (error) {
	const fallback = verificationData.value?.mode === 'slide'
		? '拼图位置不正确，请重试'
		: '点击位置不正确，请重试'
    message.error(error.message || fallback)
    await load()
  } finally {
    loading.value = false
  }
}

function close() {
  verificationID.value = ''
	dragging.value = false
  verificationData.value = null
  points.value = []
  emit('close')
}

watch(
  () => props.open,
  (visible) => {
    if (visible) load()
  },
)
</script>

<style scoped lang="less">
	.enhanced-verification {
		padding-top: 2px;
	}

	.verification-tip,
	.verification-footer {
		display: flex;
		align-items: center;
		justify-content: space-between;
		color: #666d78;
		font-size: 13px;

		button {
			padding: 0;
			color: var(--color-primary);
			background: transparent;
			border: 0;
			cursor: pointer;
		}
	}

	.verification-tip {
		margin-bottom: 10px;
	}

	.verification-prompt {
		display: flex;
		align-items: center;
		gap: 10px;
		margin-bottom: 10px;
		padding: 6px 10px;
		color: #7b818b;
		font-size: 12px;
		background: #f7f8fa;
		border-radius: 7px;

		img {
			width: 150px;
			height: 40px;
			object-fit: contain;
		}
	}

	.verification-canvas {
		position: relative;
		width: 100%;
		overflow: hidden;
		background: #f4f5f7;
		border: 1px solid #ebecef;
		border-radius: 10px;
		cursor: crosshair;
		user-select: none;
		touch-action: none;

		&.mode-slide {
			cursor: default;
		}
	}

	.verification-image {
		display: block;
		width: 100%;
		height: 100%;
		pointer-events: none;
	}

	.verification-marker {
		position: absolute;
		display: flex;
		align-items: center;
		justify-content: center;
		width: 24px;
		height: 24px;
		color: #fff;
		font-size: 12px;
		font-weight: 600;
		background: var(--color-primary);
		border: 2px solid #fff;
		border-radius: 50%;
		box-shadow: 0 2px 6px rgb(0 0 0 / 24%);
		transform: translate(-50%, -50%);
		pointer-events: none;
	}

	.verification-thumb {
		position: absolute;
		z-index: 2;
		filter: drop-shadow(0 2px 4px rgb(0 0 0 / 25%));
		cursor: ew-resize;
		touch-action: none;

		&.dragging {
			filter: drop-shadow(0 3px 7px rgb(0 0 0 / 35%));
			cursor: grabbing;
		}
	}

	.verification-rotate-thumb {
		position: absolute;
		top: 50%;
		left: 50%;
		width: 54%;
		height: 74%;
		object-fit: contain;
		pointer-events: none;
	}

	.rotation-control {
		width: 100%;
		margin-top: 12px;
		accent-color: var(--color-primary);
	}

	.behavior-pad {
		display: flex;
		align-items: center;
		justify-content: center;
		height: 190px;
		color: #7b818b;
		font-size: 14px;
		background:
			radial-gradient(circle at center, rgb(var(--color-primary-rgb) / 8%) 0 44px, transparent 45px),
			#f7f8fa;
		border: 1px dashed #d9dce2;
		border-radius: 10px;
		cursor: grab;
		touch-action: none;
		user-select: none;

		&.active {
			color: var(--color-primary);
			border-color: var(--color-primary);
			cursor: grabbing;
		}
	}

	.verification-footer {
		min-height: 20px;
		margin-top: 10px;
		color: #999fa8;
		font-size: 12px;
	}
</style>
