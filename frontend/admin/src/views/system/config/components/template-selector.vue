<template>
	<div class="template-selector">
		<a-radio-group v-model:value="selectedTemplate" class="template-selector__grid">
			<a-radio
				v-for="item in templates"
				:key="item.code"
				:value="item.code"
				class="template-selector__item"
			>
				<div class="template-selector__cover">
					<img :src="item.cover" :alt="`${item.name}模板缩略图`" />
				</div>
				<div class="template-selector__meta">
					<div class="template-selector__text">
						<strong>{{ item.name }}</strong>
						<span>{{ item.description }}</span>
					</div>
					<button type="button" @click.stop.prevent="openPreview(item)">
						<eye-outlined />
						预览
					</button>
				</div>
			</a-radio>
		</a-radio-group>

		<a-modal
			v-model:open="previewVisible"
			wrap-class-name="checkin-template-preview-modal"
			:width="'min(94vw, 1500px)'"
			:footer="null"
			:destroy-on-close="true"
			:centered="true"
		>
			<template #title>{{ previewTemplate?.name }} · 实时预览</template>
			<div class="template-selector__preview">
				<iframe
					v-if="previewTemplate"
					:src="previewUrl"
					:aria-label="`${previewTemplate.name}实时预览`"
				></iframe>
			</div>
		</a-modal>
	</div>
</template>

<script setup>
import { computed, ref } from 'vue'
import { EyeOutlined } from '@ant-design/icons-vue'
import celebrationCover from '@/assets/thumbnail/celebration.png'
import glassCover from '@/assets/thumbnail/glass.png'
import grandGalaCover from '@/assets/thumbnail/grand-gala.png'
import invitationCover from '@/assets/thumbnail/invitation.png'
import luxeCover from '@/assets/thumbnail/luxe.png'
import orbitCover from '@/assets/thumbnail/orbit.png'
import slotCover from '@/assets/thumbnail/slot.png'

const selectedTemplate = defineModel({ type: String, required: true })
const previewVisible = ref(false)
const previewTemplate = ref(null)

const templates = [
	{
		code: 'invitation',
		name: '相聚邀请函',
		description: '温暖克制的邀请函式签到大屏',
		cover: invitationCover,
	},
	{
		code: 'celebration',
		name: '欢聚盛典',
		description: '喜庆热烈的年度盛典入场舞台',
		cover: celebrationCover,
	},
	{
		code: 'luxe',
		name: '静雅之夜',
		description: '月色、星光与丝带构成的静谧会场',
		cover: luxeCover,
	},
	{
		code: 'glass',
		name: '幸运制造所',
		description: '带有物理球仓互动的未来实验室',
		cover: glassCover,
	},
	{
		code: 'slot',
		name: '鎏金游乐场',
		description: '剧场灯光与老虎机结合的欢庆空间',
		cover: slotCover,
	},
	{
		code: 'orbit',
		name: '同频星系',
		description: '持续流动的银河与星际签到入口',
		cover: orbitCover,
	},
	{
		code: 'grand-gala',
		name: '同心致远',
		description: '红金主视觉与城市盛典舞台',
		cover: grandGalaCover,
	},
]

const previewUrl = computed(() => {
	if (!previewTemplate.value) return ''
	const configuredUrl = import.meta.env.VITE_WEB_PREVIEW_URL
	const baseUrl = configuredUrl || (import.meta.env.DEV ? `${location.protocol}//${location.hostname}:5174/checkin` : '/checkin')
	const url = new URL(baseUrl, location.origin)
	url.searchParams.set('preview', '1')
	url.searchParams.set('template', previewTemplate.value.code)
	return url.toString()
})

function openPreview(item) {
	previewTemplate.value = item
	previewVisible.value = true
}
</script>

<style scoped lang="less">
	.template-selector {
		padding-top: 20px;

		&__grid {
			display: grid;
			grid-template-columns: repeat(3, minmax(0, 1fr));
			gap: 16px;
		}

		&__item {
			position: relative;
			display: block;
			margin: 0;
			padding: 6px 6px 14px;
			overflow: hidden;
			background: #fff;
			border: 1px solid #e8e8e8;
			border-radius: 9px;
			transition: border-color 0.2s, box-shadow 0.2s;

			&:hover {
				border-color: #efaaa7;
			}

			&:has(.ant-radio-checked) {
				border-color: var(--color-primary);
				box-shadow: 0 0 0 2px rgb(239 63 58 / 8%);
			}

			:deep(.ant-radio) {
				position: absolute;
				top: 14px;
				right: 14px;
				z-index: 2;
				padding: 4px;
				background: rgb(255 255 255 / 90%);
				border-radius: 50%;
			}

			:deep(.ant-radio + span) {
				display: block;
				padding-inline-start: 0;
			}
		}

		&__cover {
			position: relative;
			aspect-ratio: 16 / 9;
			overflow: hidden;
			background: #151515;
			border-radius: 5px;

			img {
				display: block;
				width: 100%;
				height: 100%;
				object-fit: cover;
			}

		}

		&__meta {
			display: flex;
			padding: 12px 6px 0;
			align-items: flex-end;
			justify-content: space-between;
			gap: 12px;
			white-space: normal;

			button {
				display: inline-flex;
				padding: 4px 0;
				flex: 0 0 auto;
				align-items: center;
				gap: 5px;
				color: #777;
				font-size: 12px;
				background: transparent;
				border: 0;
				cursor: pointer;

				&:hover {
					color: var(--color-primary);
				}
			}
		}

		&__text {
			display: flex;
			min-width: 0;
			flex-direction: column;
			gap: 4px;

			strong {
				color: #333;
				font-size: 14px;
			}

			span {
				color: #999;
				font-size: 12px;
				line-height: 18px;
			}
		}

		&__preview {
			aspect-ratio: 16 / 9;
			overflow: hidden;
			background: #080808;
			border-radius: 6px;

			iframe {
				display: block;
				width: 100%;
				height: 100%;
				border: 0;
			}
		}
	}

	@media (max-width: 1100px) {
		.template-selector__grid {
			grid-template-columns: repeat(2, minmax(0, 1fr));
		}
	}

	@media (max-width: 720px) {
		.template-selector__grid {
			grid-template-columns: 1fr;
		}
	}
</style>
