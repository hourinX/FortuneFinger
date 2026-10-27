<script setup>
const letterModules = import.meta.glob('/src/assets/branding/loading-wordmark/*.png', {
	eager: true,
	import: 'default',
	query: '?url',
})
const letterSlices = Object.entries(letterModules)
	.sort(([left], [right]) => left.localeCompare(right))
	.map(([, source]) => source)
</script>

<template>
	<div class="loading-screen" role="status" aria-live="polite" aria-label="页面正在载入">
		<div class="loading-screen__content">
			<div class="loading-screen__brand" aria-hidden="true">
				<img
					v-for="(source, index) in letterSlices"
					:key="source"
					:src="source"
					alt=""
					:style="{ '--letter-index': index }"
				/>
			</div>
			<div class="loading-screen__line" aria-hidden="true">
				<i></i>
			</div>
			<p>页面正在准备</p>
		</div>
	</div>
</template>

<style scoped lang="less">
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

		&::before,
		&::after {
			position: absolute;
			width: min(46vw, 720px);
			aspect-ratio: 1;
			border: 1px solid rgb(239 63 58 / 6%);
			border-radius: 50%;
			content: '';
		}

		&::before {
			transform: scale(0.72);
		}

		&::after {
			transform: scale(1.08);
		}

		&__content {
			position: relative;
			z-index: 1;
			display: flex;
			align-items: center;
			flex-direction: column;
		}

		&__brand {
			display: flex;
			width: clamp(250px, 28vw, 500px);
			aspect-ratio: 1864 / 294;
			align-items: stretch;
			overflow: visible;

			img {
				display: block;
				width: auto;
				min-width: 0;
				height: 100%;
				opacity: 0;
				transform: translateY(115%) scale(0.82);
				animation: letter-arrive 0.62s cubic-bezier(0.2, 0.9, 0.28, 1.2) forwards;
				animation-delay: calc(var(--letter-index) * 65ms + 80ms);
			}
		}

		&__line {
			width: min(230px, 48vw);
			height: 1px;
			margin-top: 24px;
			overflow: hidden;
			background: rgb(44 40 51 / 10%);

			i {
				display: block;
				width: 42%;
				height: 100%;
				background: linear-gradient(90deg, transparent, #ef3f3a, transparent);
				animation: line-travel 1.6s ease-in-out infinite 0.9s;
			}
		}

		p {
			margin: 14px 0 0;
			color: #9c999d;
			font-size: 10px;
			letter-spacing: 0.42em;
			opacity: 0;
			animation: caption-arrive 0.5s ease forwards 0.95s;
		}
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

	@media (prefers-reduced-motion: reduce) {
		.loading-screen {
			&__brand img,
			p {
				opacity: 1;
				transform: none;
				animation: none;
			}

			&__line i {
				width: 100%;
				animation: none;
			}
		}
	}
</style>
