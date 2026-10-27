<script setup lang="ts">
import EventBrand from '@/components/EventBrand.vue'
import { onMounted, onUnmounted, ref } from 'vue'
import { createCheckinTemplateState, type CheckinTemplateProps } from '../useCheckin'
const props = defineProps<CheckinTemplateProps>()
const { event, closed, refreshing, refreshSeconds, attendance, countdown, deadlineLabel } =
  createCheckinTemplateState(props)
const paused = ref(document.hidden)
function visibility() {
  paused.value = document.hidden
}
onMounted(() => document.addEventListener('visibilitychange', visibility))
onUnmounted(() => document.removeEventListener('visibilitychange', visibility))
const stars = Array.from({ length: 78 }, (_, i) => ({
  left: `${(i * 61.803 + 13) % 100}%`,
  top: `${(i * 37.719 + 7) % 91}%`,
  width: `${i % 9 === 0 ? 2.4 : 1.2 + (i % 3) * 0.3}px`,
  height: `${i % 9 === 0 ? 2.4 : 1.2 + (i % 3) * 0.3}px`,
  animationDelay: `${-i * 0.7}s`,
  animationDuration: `${8 + (i % 7)}s`,
}))
</script>

<template>
  <main
    class="night"
    :class="{ paused }"
  >
    <div
      class="atmosphere"
      aria-hidden="true"
    >
      <svg
        class="moon moon-detailed"
        viewBox="0 0 200 200"
      >
        <defs>
          <radialGradient
            id="lunar-body"
            cx=".32"
            cy=".25"
            r=".78"
          >
            <stop stop-color="#f0ecda" />
            <stop
              offset=".55"
              stop-color="#c9ccba"
            />
            <stop
              offset=".83"
              stop-color="#8a9c90"
            />
            <stop
              offset="1"
              stop-color="#405d56"
            />
          </radialGradient>
          <radialGradient id="lunar-crater">
            <stop
              stop-color="#5d756a"
              stop-opacity=".25"
            />
            <stop
              offset=".65"
              stop-color="#819382"
              stop-opacity=".12"
            />
            <stop
              offset=".88"
              stop-color="#e2dfc7"
              stop-opacity=".16"
            />
            <stop
              offset="1"
              stop-color="#c3ccb2"
              stop-opacity="0"
            />
          </radialGradient>
          <linearGradient id="lunar-shadow">
            <stop
              offset=".48"
              stop-color="#15282a"
              stop-opacity="0"
            />
            <stop
              offset=".8"
              stop-color="#15282a"
              stop-opacity=".06"
            />
            <stop
              offset="1"
              stop-color="#122428"
              stop-opacity=".6"
            />
          </linearGradient>
          <clipPath id="lunar-disc">
            <circle
              cx="100"
              cy="100"
              r="94"
            />
          </clipPath>
          <filter
            id="lunar-grain"
            x="0"
            y="0"
            width="100%"
            height="100%"
          >
            <feTurbulence
              type="fractalNoise"
              baseFrequency=".11"
              numOctaves="3"
              seed="12"
            />
            <feColorMatrix
              type="saturate"
              values="0"
            />
            <feComponentTransfer>
              <feFuncA
                type="linear"
                slope=".12"
              />
            </feComponentTransfer>
            <feBlend
              in="SourceGraphic"
              mode="multiply"
            />
          </filter>
        </defs>
        <g clip-path="url(#lunar-disc)">
          <circle
            cx="100"
            cy="100"
            r="94"
            fill="url(#lunar-body)"
          />
          <g
            fill="#6c8173"
            opacity=".19"
            filter="url(#lunar-grain)"
          >
            <path
              d="M44 58Q68 33 91 52T119 86Q102 111 76 97T44 58ZM115 104Q143 74 160 109T142 152Q117 164 115 137Z"
            />
            <ellipse
              cx="65"
              cy="127"
              rx="22"
              ry="16"
            />
          </g>
          <circle
            v-for="n in 18"
            :key="n"
            :cx="28 + ((n * 37) % 143)"
            :cy="23 + ((n * 53) % 150)"
            :r="4 + ((n * 7) % 15)"
            fill="url(#lunar-crater)"
          />
          <circle
            cx="100"
            cy="100"
            r="94"
            fill="transparent"
            filter="url(#lunar-grain)"
          />
          <circle
            cx="100"
            cy="100"
            r="94"
            fill="url(#lunar-shadow)"
          />
        </g>
        <path
          d="M35 168A94 94 0 0 1 150 20"
          fill="none"
          stroke="#f5edce"
          stroke-opacity=".32"
          stroke-width=".7"
        />
      </svg>
      <div class="moon-halo"></div>
      <span
        v-for="(star, i) in stars"
        :key="i"
        class="star"
        :style="star"
      ></span>
      <svg
        class="landscape"
        viewBox="0 0 1920 1080"
        preserveAspectRatio="xMidYMid slice"
      >
        <defs>
          <radialGradient
            id="night-silk"
            cx=".65"
            cy=".25"
            r=".85"
          >
            <stop
              stop-color="#bcccbc"
              stop-opacity=".3"
            />
            <stop
              offset=".45"
              stop-color="#809c90"
              stop-opacity=".21"
            />
            <stop
              offset=".8"
              stop-color="#53776b"
              stop-opacity=".13"
            />
            <stop
              offset="1"
              stop-color="#3a5b53"
              stop-opacity=".06"
            />
          </radialGradient>
          <radialGradient
            id="night-turn"
            cx=".57"
            cy=".18"
            r=".95"
          >
            <stop
              stop-color="#d2d8c0"
              stop-opacity=".58"
            />
            <stop
              offset=".3"
              stop-color="#b0c2ae"
              stop-opacity=".45"
            />
            <stop
              offset=".64"
              stop-color="#6d9282"
              stop-opacity=".22"
            />
            <stop
              offset="1"
              stop-color="#35564c"
              stop-opacity="0"
            />
          </radialGradient>
          <linearGradient id="night-edge">
            <stop
              stop-color="#bdcab5"
              stop-opacity="0"
            />
            <stop
              offset=".38"
              stop-color="#e8e2c9"
              stop-opacity=".32"
            />
            <stop
              offset=".64"
              stop-color="#b4c7b9"
              stop-opacity=".18"
            />
            <stop
              offset="1"
              stop-color="#b6c8ba"
              stop-opacity=".24"
            />
          </linearGradient>
          <radialGradient
            id="night-silk-front"
            cx=".56"
            cy=".35"
            r=".9"
          >
            <stop
              stop-color="#b7cbb7"
              stop-opacity="0.288"
            />
            <stop
              offset=".42"
              stop-color="#90b09c"
              stop-opacity="0.224"
            />
            <stop
              offset=".75"
              stop-color="#698d7c"
              stop-opacity="0.152"
            />
            <stop
              offset="1"
              stop-color="#43685b"
              stop-opacity="0.072"
            />
          </radialGradient>
          <radialGradient
            id="night-sheen"
            cx=".55"
            cy=".25"
            r=".85"
          >
            <stop
              stop-color="#d7dec4"
              stop-opacity=".14"
            />
            <stop
              offset=".5"
              stop-color="#acc7ae"
              stop-opacity=".07"
            />
            <stop
              offset="1"
              stop-color="#95b6a0"
              stop-opacity="0"
            />
          </radialGradient>
          <filter
            id="night-ribbon-haze"
            x="-15%"
            y="-45%"
            width="130%"
            height="190%"
            color-interpolation-filters="sRGB"
          >
            <feGaussianBlur
              in="SourceAlpha"
              stdDeviation="12"
              result="halo"
            />
            <feFlood
              flood-color="#bed8bc"
              flood-opacity=".14"
            />
            <feComposite
              in2="halo"
              operator="in"
            />
            <feMerge>
              <feMergeNode />
              <feMergeNode in="SourceGraphic" />
            </feMerge>
          </filter>
          <filter
            id="night-fold-soften"
            x="-20%"
            y="-30%"
            width="140%"
            height="160%"
          >
            <feGaussianBlur stdDeviation="6" />
          </filter>
          <mask
            id="night-fold-feather"
            maskUnits="userSpaceOnUse"
            x="500"
            y="750"
            width="450"
            height="300"
          >
            <rect
              x="500"
              y="750"
              width="450"
              height="300"
              fill="url(#night-fold-mask)"
            />
          </mask>
          <linearGradient
            id="night-fold-mask"
            gradientUnits="userSpaceOnUse"
            x1="545"
            y1="0"
            x2="900"
            y2="0"
          >
            <stop stop-color="black" />
            <stop
              offset=".32"
              stop-color="white"
            />
            <stop
              offset=".8"
              stop-color="white"
            />
            <stop
              offset="1"
              stop-color="#777"
            />
          </linearGradient>
          <filter
            id="night-contact"
            x="-20%"
            y="-50%"
            width="140%"
            height="200%"
          >
            <feGaussianBlur stdDeviation="5" />
          </filter>
          <linearGradient
            id="night-underlay"
            x1="0"
            y1="0"
            x2="0"
            y2="1"
          >
            <stop
              stop-color="#1c3433"
              stop-opacity="0.342"
            />
            <stop
              offset=".5"
              stop-color="#203b36"
              stop-opacity="0.405"
            />
            <stop
              offset="1"
              stop-color="#142b2a"
              stop-opacity="0.293"
            />
          </linearGradient>
          <radialGradient
            id="night-fold-depth"
            cx=".62"
            cy=".25"
            r=".85"
          >
            <stop
              stop-color="#c9d2bc"
              stop-opacity="0.416"
            />
            <stop
              offset=".3"
              stop-color="#9eb6a2"
              stop-opacity="0.312"
            />
            <stop
              offset=".65"
              stop-color="#4f7765"
              stop-opacity="0.293"
            />
            <stop
              offset="1"
              stop-color="#183a30"
              stop-opacity="0.403"
            />
          </radialGradient>
          <filter
            id="night-distance"
            x="-10%"
            y="-20%"
            width="120%"
            height="140%"
          >
            <feGaussianBlur stdDeviation="1.8" />
          </filter>
          <filter
            id="night-occlusion"
            x="-15%"
            y="-30%"
            width="130%"
            height="160%"
          >
            <feGaussianBlur stdDeviation="7" />
          </filter>
          <radialGradient id="night-reflection">
            <stop
              stop-color="#d3d7be"
              stop-opacity=".12"
            />
            <stop
              offset="1"
              stop-color="#c1cdb7"
              stop-opacity="0"
            />
          </radialGradient>
        </defs>
        <ellipse
          cx="1190"
          cy="885"
          rx="650"
          ry="180"
          fill="url(#night-reflection)"
        />
        <!-- 每次翻折分为背面、转折面与迎光面；局部交叉由绘制顺序和接触阴影表现。 -->
        <g
          class="ribbon ribbon-far"
          opacity=".32"
          filter="url(#night-distance)"
        >
          <path
            d="M-100 752 C130 768 180 1025 470 921 C730 828 1010 731 1245 850 C1480 970 1710 907 2030 720 L2030 786 C1740 1000 1510 1030 1215 929 C990 851 777 915 492 969 C193 1060 80 870 -100 858Z"
            fill="url(#night-silk)"
          />
          <path
            d="M-100 752 C130 768 180 1025 470 921 C730 828 1010 731 1245 850"
            fill="none"
            stroke="url(#night-edge)"
            stroke-width=".8"
          />
        </g>
        <g
          class="ribbon ribbon-near"
          filter="url(#night-ribbon-haze)"
        >
          <path
            d="M540 953 C685 857 718 801 797 801 C862 797 910 856 865 901"
            fill="none"
            stroke="#071c1b"
            stroke-width="16"
            opacity=".29"
            filter="url(#night-occlusion)"
          />
          <!-- 左侧宽幅垂坠，向中央缩窄后反折。 -->
          <path
            d="M-80 804 C180 1040 340 1040 545 936 C647 885 686 791 790 784 C877 779 922 851 866 897 C818 935 755 934 709 947 C1050 1040 1440 1017 2020 782 L2020 935 C1520 1108 1048 1104 707 995 C620 966 682 912 772 878 C838 854 819 817 790 830 C712 863 713 932 571 1008 C320 1145 120 1096 -80 961Z"
            fill="url(#night-underlay)"
          />
          <path
            d="M-80 804 C180 1040 340 1040 545 936 C647 885 686 791 790 784 C877 779 922 851 866 897 C818 935 755 934 709 947 C1050 1040 1440 1017 2020 782 L2020 935 C1520 1108 1048 1104 707 995 C620 966 682 912 772 878 C838 854 819 817 790 830 C712 863 713 932 571 1008 C320 1145 120 1096 -80 961Z"
            fill="url(#night-silk-front)"
          />
          <!-- 翻折内壁较暗，窄颈到宽面的变化构成扭转。 -->
          <path
            d="M545 936 C647 885 686 791 790 784 C877 779 922 851 866 897 C818 935 755 934 709 947 C688 929 735 893 772 878 C838 854 819 817 790 830 C712 863 713 932 571 1008Z"
            fill="url(#night-fold-depth)"
            opacity=".9"
            mask="url(#night-fold-feather)"
          />
          <path
            d="M704 946 C770 937 827 924 866 897"
            fill="none"
            stroke="#0d2525"
            stroke-width="11"
            opacity=".42"
            filter="url(#night-contact)"
          />
          <!-- 月光落在翻折肩部及前景铺展面上。 -->
          <path
            d="M781 788 C839 782 891 818 889 855"
            fill="none"
            stroke="#d4dcc2"
            stroke-width=".8"
            opacity=".2"
          />
          <path
            d="M772 878 C738 893 716 909 709 929"
            fill="none"
            stroke="#0d2925"
            stroke-width="5"
            opacity=".4"
            filter="url(#night-contact)"
          />
          <path
            d="M709 947 C1050 1040 1440 1017 2020 782 L2020 831 C1450 1036 1100 1066 709 947Z"
            fill="url(#night-sheen)"
            filter="url(#night-fold-soften)"
          />
          <path
            d="M-80 804 C180 1040 340 1040 545 936 C647 885 686 791 790 784 C847 781 884 811 890 844"
            fill="none"
            stroke="url(#night-edge)"
            stroke-width="1"
          />
          <path
            d="M709 947 C1050 1040 1440 1017 2020 782"
            fill="none"
            stroke="url(#night-edge)"
            stroke-width="1.2"
          />
        </g>
        <!-- 右侧第二条短带收尖，在主带前方形成一段独立叠层。 -->
        <g
          class="ribbon ribbon-fold"
          filter="url(#night-ribbon-haze)"
        >
          <path
            d="M1410 988 C1560 883 1660 778 1815 804 C1885 820 1900 866 1858 893"
            fill="none"
            stroke="#071c1b"
            stroke-width="14"
            opacity=".3"
            filter="url(#night-occlusion)"
          />
          <path
            d="M1390 982 C1510 864 1660 735 1830 786 C1900 807 1910 866 1860 887 C1815 909 1770 895 1725 914 C1850 913 1950 930 2020 959 L2020 1024 C1900 972 1790 968 1720 952 C1650 936 1680 894 1785 861 C1837 844 1814 813 1780 820 C1630 844 1520 952 1390 982Z"
            fill="url(#night-underlay)"
          />
          <path
            d="M1390 982 C1510 864 1660 735 1830 786 C1900 807 1910 866 1860 887 C1815 909 1770 895 1725 914 C1850 913 1950 930 2020 959 L2020 1024 C1900 972 1790 968 1720 952 C1650 936 1680 894 1785 861 C1837 844 1814 813 1780 820 C1630 844 1520 952 1390 982Z"
            fill="url(#night-fold-depth)"
            opacity=".88"
          />
          <path
            d="M1390 982 C1510 864 1660 735 1830 786 C1881 802 1900 837 1884 865"
            fill="none"
            stroke="url(#night-edge)"
            stroke-width=".8"
          />
        </g>
      </svg>
    </div>
    <header>
      <EventBrand />
      <span>AN EVENING OF QUIET BRILLIANCE</span>
      <span>年度 · 雅集</span>
    </header>
    <section class="stage">
      <div class="poetry">
        <p class="eyebrow">良 辰 · 与 君 共 赴</p>
        <h1>
          <span>静雅</span>
          <span>
            之夜
            <span class="seal">雅集</span>
          </span>
        </h1>
        <p class="subtitle">光华不语，相聚有声。</p>
        <p class="description">{{ event.description }}</p>
        <section
          class="status"
          aria-label="签到状态"
        >
          <div class="guests">
            <span class="label">已赴此约</span>
            <div>
              <strong
                :key="event.checkedIn"
                class="guest-number"
              >
                {{ event.checkedIn }}
              </strong>
              <span>/ {{ event.total }} 位</span>
            </div>
          </div>
          <div
            class="attendance-ring"
            role="progressbar"
            aria-label="签到进度"
            :aria-valuenow="attendance"
            aria-valuemin="0"
            aria-valuemax="100"
          >
            <svg
              viewBox="0 0 64 64"
              aria-hidden="true"
            >
              <circle
                cx="32"
                cy="32"
                r="28"
                class="ring-track"
              />
              <circle
                cx="32"
                cy="32"
                r="28"
                class="ring-value"
                pathLength="100"
                :stroke-dasharray="`${attendance} 100`"
              />
            </svg>
            <span>
              {{ attendance }}
              <small>%</small>
            </span>
            <p>相聚渐满</p>
          </div>
          <div class="time">
            <span class="label">{{ closed ? '签到已截止' : '入席倒计时' }}</span>
            <div class="clock">
              <template
                v-for="(part, i) in countdown"
                :key="i"
              >
                <b v-if="i">:</b>
                <span>
                  {{ part }}
                  <small>{{ ['时', '分', '秒'][i] }}</small>
                </span>
              </template>
            </div>
            <p>{{ deadlineLabel }} 截止</p>
          </div>
        </section>
      </div>
      <section
        class="entry"
        aria-label="扫码签到"
      >
        <span class="entry-eyebrow">YOUR PRESENCE, OUR HONOR</span>
        <div class="qr-shell">
          <div
            class="qr"
            :class="{ refreshing }"
          >
            <svg
              viewBox="0 0 64 64"
              fill="none"
              aria-hidden="true"
            >
              <path
                d="M8 24V8h16M40 8h16v16M56 40v16H40M24 56H8V40M25 25h14v14H25Z"
                stroke="currentColor"
                stroke-width="2"
              />
            </svg>
            <strong>
              {{ closed ? '签到已结束' : refreshing ? '二维码已刷新' : '签到二维码' }}
            </strong>
            <small>演示占位 · 不可扫描</small>
          </div>
        </div>
        <h2>扫码入席</h2>
        <p>席位已备，静候光临</p>
        <span class="refresh">{{ closed ? '已停止刷新' : `${refreshSeconds}s 后更新签到码` }}</span>
      </section>
    </section>
    <footer>
      <span>A NIGHT TO REMEMBER</span>
      <span>同赏今宵 · 共赴新程</span>
    </footer>
  </main>
</template>

<style scoped>
.night {
  position: relative;
  isolation: isolate;
  height: 100dvh;
  overflow: auto;
  display: flex;
  flex-direction: column;
  padding: 30px 6vw 78px;
  background:
    radial-gradient(ellipse at 66% 38%, #233738 0%, transparent 50%),
    linear-gradient(145deg, #10191e, #15272b 65%, #101a21);
  color: #d8d6c9;
  font-family: 'Sora Variable', 'Microsoft YaHei', sans-serif;
}
.atmosphere {
  position: absolute;
  inset: 0;
  overflow: hidden;
  z-index: -1;
  pointer-events: none;
}
.moon {
  position: absolute;
  width: 35vw;
  height: 35vw;
  max-width: 650px;
  max-height: 650px;
  border-radius: 50%;
  top: 8%;
  left: 44%;
  background: radial-gradient(circle at 35% 35%, #e3e2cd09, #b5c7be03 65%, transparent 72%);
  box-shadow: inset -1px 1px 0 #e2d5b645;
  transform: rotate(-25deg);
}
.moon-halo {
  position: absolute;
  inset: 0;
  background: radial-gradient(ellipse at 60% 38%, #c0b38a0c, transparent 45%);
  animation: breath 16s ease-in-out infinite alternate;
}
.landscape {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
}
.water {
  transform-origin: center;
}
.water-back {
  animation: water 24s ease-in-out infinite alternate;
}
.water-front {
  animation: water 19s ease-in-out -9s infinite alternate-reverse;
}
.star {
  position: absolute;
  width: 2px;
  height: 2px;
  border-radius: 50%;
  background: #e2d7bb;
  animation: starlight ease-in-out infinite alternate;
  opacity: 0.2;
}
.star:nth-child(7n) {
  width: 3px;
  height: 3px;
  box-shadow: 0 0 7px #dfcda34a;
}
header,
footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 25px;
}
header {
  font-size: 9px;
  letter-spacing: 3px;
  color: #8e9b96;
}
.brand {
  font-size: 18px;
  letter-spacing: 0;
  color: #c9c9bb;
}
header > span:last-child {
  letter-spacing: 5px;
}
.stage {
  display: grid;
  grid-template-columns: 1.2fr 1fr;
  align-items: center;
  gap: 8vw;
  flex: 1;
  width: 100%;
  max-width: 1500px;
  margin: 0 auto;
  padding: 35px 3vw 35px;
}
.poetry {
  position: relative;
}
.eyebrow {
  font-size: 11px;
  letter-spacing: 4px;
  color: #adac98;
  margin: 0 0 20px;
}
h1 {
  font-family: 'STXingkai', '华文行楷', 'KaiTi', '楷体', 'Songti SC', serif;
  font-size: clamp(58px, 6.4vw, 116px);
  font-weight: 400;
  letter-spacing: 0.16em;
  line-height: 1.12;
  margin: 0;
  transform: rotate(-3deg);
}
h1 > span {
  display: block;
  width: fit-content;
  color: #d5cbb3;
  background: linear-gradient(
    115deg,
    #899e9c 5%,
    #bbbda9 27%,
    #e1d7bf 45%,
    #b5b8a5 53%,
    #93a5a2 78%,
    #d3cbb4
  );
  background-size: 180% 100%;
  background-clip: text;
  -webkit-text-fill-color: transparent;
  filter: drop-shadow(0 2px 1px #081319);
  animation: title-sheen 22s ease-in-out infinite alternate;
}
h1 > span + span {
  margin-left: 1.15em;
  margin-top: 0.12em;
  position: relative;
}
.seal {
  position: absolute;
  right: -15px;
  bottom: 8px;
  font:
    11px 'SimSun',
    serif;
  writing-mode: vertical-rl;
  letter-spacing: 3px;
  color: #b8aa8b;
  -webkit-text-fill-color: #b8aa8b;
  border: 1px solid #9d967d80;
  padding: 7px 4px;
  transform: rotate(3deg);
}
.subtitle {
  font-family: 'Songti SC', 'SimSun', serif;
  font-size: 20px;
  letter-spacing: 5px;
  margin: 30px 0 16px;
  color: #c0c1b0;
}
.description {
  max-width: 360px;
  font-size: 11px;
  line-height: 2;
  color: #819491;
}
.entry {
  text-align: center;
  justify-self: center;
  width: 100%;
  max-width: 330px;
  padding-top: 30px;
}
.entry-eyebrow {
  display: block;
  font-size: 8px;
  letter-spacing: 3px;
  color: #91a19a;
  margin-bottom: 25px;
}
.qr-shell {
  position: relative;
  width: clamp(170px, 17vw, 260px);
  margin: auto;
  padding: 12px;
  border: 1px solid #adad9045;
}
.qr-shell::before,
.qr-shell::after {
  content: '';
  position: absolute;
  width: 28px;
  height: 28px;
  border: solid #bcb39590;
  pointer-events: none;
}
.qr-shell::before {
  top: -1px;
  left: -1px;
  border-width: 1px 0 0 1px;
}
.qr-shell::after {
  bottom: -1px;
  right: -1px;
  border-width: 0 1px 1px 0;
}
.qr {
  aspect-ratio: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 14px;
  background: #e7e6db;
  color: #344744;
  transition: background 0.5s;
}
.qr svg {
  width: 28%;
}
.qr strong {
  font-size: 12px;
  font-weight: 400;
}
.qr small {
  font-size: 9px;
  color: #8d958a;
}
.refreshing {
  background: #f5f1e4;
}
h2 {
  font-family: 'Songti SC', 'SimSun', serif;
  font-size: 23px;
  font-weight: 400;
  letter-spacing: 6px;
  margin: 25px 0 10px;
}
.entry > p {
  font-size: 10px;
  letter-spacing: 3px;
  color: #95a19a;
}
.refresh {
  font-size: 9px;
  color: #657f7b;
  display: block;
  margin-top: 15px;
}
.status {
  display: grid;
  grid-template-columns: 1fr 1.25fr 1.2fr;
  align-items: center;
  gap: 6vw;
  border-top: 1px solid #afbaa020;
  padding: 25px 3vw 20px;
}
.label {
  font-size: 10px;
  letter-spacing: 3px;
  color: #9da79b;
}
.guests > div {
  display: flex;
  align-items: baseline;
  gap: 14px;
  margin-top: 6px;
}
.guest-number {
  font-size: 48px;
  font-weight: 300;
  color: #d2cbb5;
  font-variant-numeric: tabular-nums;
  animation: arrival 0.6s ease;
}
.guests > div > span {
  color: #6d8881;
  font-size: 12px;
}
.time p {
  font-size: 9px;
  color: #657f78;
  margin: 13px 0 0;
}
.time {
  justify-self: end;
}
.clock {
  display: flex;
  align-items: baseline;
  gap: 12px;
  margin-top: 12px;
  font-variant-numeric: tabular-nums;
}
.clock > span {
  font-size: 30px;
  font-weight: 300;
  color: #c7c7b4;
}
.clock small {
  font-size: 8px;
  margin-left: 7px;
  color: #748c80;
}
.clock b {
  font-size: 18px;
  font-weight: 300;
  color: #5e7971;
}
footer {
  font-size: 8px;
  letter-spacing: 2px;
  color: #536e69;
  padding: 12px 3vw 0;
}
.paused * {
  animation-play-state: paused !important;
}
@keyframes title-sheen {
  to {
    background-position: 100% 0;
  }
}
@keyframes breath {
  to {
    opacity: 0.35;
    transform: translateX(20px);
  }
}
@keyframes starlight {
  to {
    opacity: 0.55;
    translate: 0 -9px;
  }
}
@keyframes arrival {
  from {
    opacity: 0.4;
    transform: translateY(5px);
  }
  to {
    opacity: 1;
    transform: none;
  }
}
@media (max-width: 800px) {
  .night {
    padding: 24px 24px 85px;
  }
  header > span:nth-child(2) {
    display: none;
  }
  header > span:last-child {
    font-size: 8px;
    letter-spacing: 2px;
  }
  .stage {
    grid-template-columns: 1fr;
    gap: 20px;
    padding: 40px 0 30px;
  }
  .poetry {
    justify-self: center;
  }
  h1 {
    font-size: 70px;
  }
  .subtitle {
    font-size: 17px;
    letter-spacing: 3px;
  }
  .description {
    font-size: 11px;
  }
  .moon {
    width: 90vw;
    height: 90vw;
    left: 20%;
    top: 15%;
  }
  .entry {
    padding-top: 10px;
  }
  .qr-shell {
    width: 200px;
  }
  .status {
    grid-template-columns: 1fr 1fr;
    gap: 25px;
    padding: 22px 0;
  }
  .time {
    grid-column: 1/-1;
    justify-self: stretch;
    border-top: 1px solid #afbaa020;
    padding-top: 20px;
  }
  .clock {
    gap: 17px;
  }
  .clock > span {
    font-size: 34px;
  }
  .guest-number {
    font-size: 40px;
  }
  .label {
    font-size: 9px;
    letter-spacing: 1px;
  }
  footer {
    padding: 10px 0 0;
    font-size: 7px;
    letter-spacing: 1px;
  }
}
/* 顶部月轮是局部光源，光晕不覆盖整个画面。 */
.moon {
  width: clamp(78px, 8vw, 150px);
  height: auto;
  aspect-ratio: 1;
  left: 62%;
  top: 9%;
  max-width: none;
  max-height: none;
  transform: none;
  background:
    radial-gradient(ellipse at 38% 33%, #929f9230 0%, transparent 28%),
    radial-gradient(ellipse at 68% 62%, #73857a35 0%, transparent 24%),
    radial-gradient(circle at 35% 28%, #f0edda, #cbd1bd 56%, #859b91 91%);
  box-shadow:
    0 0 35px #d3dbc71a,
    0 0 110px #b8cbb31a,
    inset -10px -6px 20px #233d3b50;
  opacity: 0.88;
}
.moon-halo {
  background: radial-gradient(ellipse 35% 65% at 67% 25%, #b8c5aa13, transparent 74%);
}
.ribbon {
  transform-origin: 50% 85%;
  animation: ribbon-drift 20s ease-in-out infinite alternate;
}
.ribbon-far {
  animation-duration: 27s;
  animation-delay: -10s;
}
.ribbon-near {
  animation-delay: -4s;
}
.stage {
  padding-top: 55px;
  padding-bottom: 120px;
  align-items: center;
}
.status {
  grid-template-columns: auto 72px;
  gap: 22px 28px;
  border: 0;
  padding: 28px 0 0;
  max-width: 360px;
  align-items: center;
}
.guests > div {
  margin-top: 8px;
}
.guest-number {
  font-size: 38px;
}
.attendance-ring {
  position: relative;
  width: 64px;
  height: 64px;
  justify-self: end;
}
.attendance-ring svg {
  width: 64px;
  height: 64px;
  transform: rotate(-90deg);
}
.attendance-ring circle {
  fill: none;
  stroke-width: 1.5;
}
.ring-track {
  stroke: #b5c5ac20;
}
.ring-value {
  stroke: #c4c9af;
  stroke-linecap: round;
  transition: stroke-dasharray 0.7s ease;
}
.attendance-ring > span {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
  color: #c4c9af;
  font-variant-numeric: tabular-nums;
}
.attendance-ring small {
  font-size: 8px;
  margin-left: 2px;
}
.attendance-ring p {
  text-align: center;
  margin: 5px -10px 0;
  font-size: 8px;
  color: #819489;
  letter-spacing: 2px;
}
.time {
  grid-column: 1/-1;
  justify-self: stretch;
  border-top: 1px solid #afbaa01c;
  padding-top: 15px;
  display: grid;
  grid-template-columns: auto 1fr;
  align-items: center;
  gap: 8px 16px;
}
.clock {
  margin: 0;
  justify-self: end;
  gap: 8px;
}
.clock > span {
  font-size: 23px;
}
.clock small {
  font-size: 7px;
  margin-left: 4px;
}
.time p {
  grid-column: 1/-1;
  margin: 0;
  text-align: right;
  font-size: 8px;
}
footer {
  position: relative;
  margin-top: auto;
  color: #82968a;
}
@keyframes ribbon-drift {
  to {
    transform: translateY(-6px) rotate(0.3deg);
  }
}
@media (max-width: 800px) {
  .moon {
    left: auto;
    right: 8%;
    top: 95px;
    width: 72px;
  }
  .stage {
    padding-top: 45px;
    padding-bottom: 70px;
  }
  .poetry {
    width: 100%;
    max-width: 360px;
  }
  .status {
    grid-template-columns: 1fr 72px;
    gap: 20px;
    padding-top: 25px;
  }
  .time {
    padding-top: 15px;
  }
  .clock > span {
    font-size: 22px;
  }
  .clock {
    gap: 7px;
  }
  .entry {
    padding-top: 10px;
  }
  .landscape {
    height: 100%;
    opacity: 0.65;
  }
}
.moon.moon-detailed {
  background: none;
  box-shadow: none;
  overflow: visible;
  filter: drop-shadow(0 0 10px #d4dcc024) drop-shadow(0 0 34px #b3c6af14);
  opacity: 0.9;
}
.ribbon-fold {
  animation-duration: 26s;
  animation-delay: -12s;
}
@media (prefers-reduced-motion: reduce) {
  *,
  *::before,
  *::after {
    animation: none !important;
    transition: none !important;
  }
}
</style>
