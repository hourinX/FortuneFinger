<template>
  <section class="settings-page">
    <SettingsNav v-model="activeSection" :items="sections" />

    <main class="settings-content">
      <header class="settings-content__header">
        <div class="settings-content__title">
          <strong>{{ currentSection.label }}</strong>
          <a-popover :content="currentSection.desc" :trigger="['hover', 'click']" placement="right">
            <button
              type="button"
              class="settings-content__help"
              :aria-label="`${currentSection.label}说明`"
            >
              <exclamation-circle-outlined />
            </button>
          </a-popover>
        </div>
        <a-button
          v-permission="'system:config:save'"
          type="primary"
          :loading="saving"
          :disabled="loading"
          @click="saveCurrent"
        >
          保存设置
        </a-button>
      </header>

      <a-spin :spinning="loading">
        <component
          :is="currentSection.component"
          v-model="settings[activeSection]"
          :test-email="testEmail"
          :test-loading="testingMail"
          @update:test-email="testEmail = $event"
          @send-test="sendMailTest"
        />
      </a-spin>
    </main>
  </section>
</template>

<script setup>
import { computed, onActivated, reactive, ref, toRaw, watch } from 'vue'
import {
  CloudUploadOutlined,
  ExclamationCircleOutlined,
  MailOutlined,
  QrcodeOutlined,
  SettingOutlined,
  TrophyOutlined,
} from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { useRoute, useRouter } from 'vue-router'
import { getConfig, saveConfig, sendTestMail } from '@/api/config'
import { useAppStore } from '@/store'
import { DAY_MILLISECONDS, getDateTimeDifference } from '@/utils/datetime'
import BasicSettings from './components/basic.vue'
import LotteryRuleSettings from './components/lottery-rule.vue'
import MailSettings from './components/mail.vue'
import SettingsNav from './components/settings-nav.vue'
import SignSettings from './components/sign.vue'
import UploadSettings from './components/upload.vue'

const sections = [
  {
    code: 'basic',
    label: '基础设置',
    desc: '设置系统名称和基础展示信息',
    permission: 'system:config:info',
    icon: SettingOutlined,
    component: BasicSettings,
  },
  {
    code: 'lottery-rule',
    label: '抽奖规则',
    desc: '设置全局中奖限制和默认处理规则',
    permission: 'system:config:info',
    icon: TrophyOutlined,
    component: LotteryRuleSettings,
  },
  {
    code: 'sign',
    label: '签到设置',
    desc: '设置扫码签到时的信息收集与校验方式',
    permission: 'system:config:info',
    icon: QrcodeOutlined,
    component: SignSettings,
  },
  {
    code: 'upload',
    label: '上传设置',
    desc: '设置管理端允许上传的文件类型和大小',
    permission: 'system:config:info',
    icon: CloudUploadOutlined,
    component: UploadSettings,
  },
  {
    code: 'mail',
    label: '邮件设置',
    desc: '设置系统邮件发送服务和默认发件人信息',
    permission: 'system:config:info',
    icon: MailOutlined,
    component: MailSettings,
  },
]
const route = useRoute()
const router = useRouter()
const app = useAppStore()
const defaultCode = sections[0].code
const findSection = (code) => sections.find((item) => item.code === code)
const configRoutePath = route.path
const configRouterID = String(route.meta.routerId || '')
const routerState = app.getVisitedRouterState(configRouterID)
const activeSection = ref(
  findSection(route.query.code)?.code || findSection(routerState?.code)?.code || defaultCode,
)
const currentSection = computed(() => findSection(activeSection.value))
const initialSettings = {
  basic: {
    name: 'FortuneFinger',
    title: 'FortuneFinger 管理端',
    shortName: 'FortuneFinger',
    language: 'zh-CN',
    footer: 'FortuneFinger 企业活动抽奖系统',
    maintenance: false,
  },
  'lottery-rule': {
    allowRepeat: false,
    allowSameLevel: false,
    maxWins: 1,
    excludeWinners: true,
    overflowMode: 'random',
  },
  sign: {
    mode: 1,
    checkin_start_time: '',
    checkin_end_time: '',
    requireName: true,
    requireCompany: true,
    allowAvatar: true,
    duplicateMode: 'show',
    enterPhotoWall: true,
    photoWallNameMode: 'name_company',
    screen: {
      template: 'invitation',
    },
  },
  upload: {
    imageMaxSize: 10,
    videoMaxSize: 500,
    imageFormats: ['JPG', 'PNG', 'WEBP'],
    videoFormats: ['MP4', 'WEBM'],
    keepOriginalName: true,
  },
  mail: {
    enabled: 2,
    host: 'smtp-relay.brevo.com',
    port: 587,
    security: '1',
    username: '',
    password: '',
    senderName: 'FortuneFinger',
    senderEmail: '',
    replyTo: '',
    timeout: 10,
  },
}
const settings = reactive(structuredClone(initialSettings))
const configIDs = reactive(Object.fromEntries(sections.map(({ code }) => [code, '0'])))
const loading = ref(false)
const saving = ref(false)
const testingMail = ref(false)
const testEmail = ref('')
let requestSequence = 0

watch(
  activeSection,
  (code) => {
    app.setVisitedRouterState(configRouterID, { code })
    if (route.path === configRoutePath && route.query.code !== code) {
      router.push({ query: { ...route.query, code } })
    }
    loadConfig(code)
  },
  { immediate: true },
)

watch(
  () => route.query.code,
  (code) => {
    if (route.path !== configRoutePath) return
    const section = findSection(code)
    if (section) {
      activeSection.value = section.code
      return
    }
    const state = app.getVisitedRouterState(configRouterID)
    const fallbackCode = findSection(state?.code)?.code || defaultCode
    if (activeSection.value !== fallbackCode) activeSection.value = fallbackCode
    router.replace({ query: { ...route.query, code: fallbackCode } })
  },
)

onActivated(() => {
  const state = app.getVisitedRouterState(configRouterID)
  const code = findSection(route.query.code)?.code || findSection(state?.code)?.code || defaultCode
  if (activeSection.value !== code) activeSection.value = code
  else app.setVisitedRouterState(configRouterID, { code })
  if (route.path === configRoutePath && route.query.code !== code) {
    router.replace({ query: { ...route.query, code } })
  }
})

async function loadConfig(code) {
  if (!app.hasPermission('system:config:info')) return
  const sequence = ++requestSequence
  loading.value = true
  try {
    const resp = await getConfig(code)
    if (sequence !== requestSequence || resp.code !== code) return
    configIDs[code] = resp.id || '0'
    const config = {
      ...structuredClone(initialSettings[code]),
      ...structuredClone(resp.config),
    }
    settings[code] = config
  } catch (error) {
    if (sequence === requestSequence) message.error(error.message || '加载系统配置失败')
  } finally {
    if (sequence === requestSequence) loading.value = false
  }
}

async function saveCurrent() {
  const code = activeSection.value
  if (code === 'sign') {
    const sign = settings.sign
    if (!sign.checkin_start_time || !sign.checkin_end_time) {
      message.warning('请选择允许开始签到时间和签到截止时间')
      return
    }
    const difference = getDateTimeDifference(sign.checkin_start_time, sign.checkin_end_time)
    if (!Number.isFinite(difference)) {
      message.warning('签到时间格式无效，请重新选择')
      return
    }
    if (difference <= 0) {
      message.warning('签到截止时间必须晚于允许开始签到时间')
      return
    }
    if (difference > 7 * DAY_MILLISECONDS) {
      message.warning('签到开放时长不能超过 7 天')
      return
    }
  }
  saving.value = true
  try {
    const config = structuredClone(toRaw(settings[code]))
    if (code === 'mail') delete config.testEmail
    const resp = await saveConfig({
      id: configIDs[code],
      code,
      config,
    })
    if (resp.code === code && resp.id) configIDs[code] = resp.id
    message.success('保存成功')
    await loadConfig(code)
  } catch (error) {
    message.error(error.message || '保存系统配置失败')
  } finally {
    saving.value = false
  }
}

async function sendMailTest(email) {
  testingMail.value = true
  try {
    await sendTestMail({ email })
    message.success('测试邮件发送成功')
  } catch (error) {
    message.error(error.message || '测试邮件发送失败')
  } finally {
    testingMail.value = false
  }
}
</script>

<style scoped lang="less">
.settings-page {
  display: grid;
  height: calc(100vh - 170px);
  min-height: calc(100vh - 170px);
  grid-template-columns: 210px minmax(0, 1fr);
  overflow: hidden;
  border: 1px solid #efeded;
  border-radius: 8px;
}

.settings-content {
  min-width: 0;
  min-height: 0;
  padding: 0 28px 28px;
  overflow-y: auto;
  scrollbar-width: none;
  background: #fff;

  &::-webkit-scrollbar {
    display: none;
  }
}

.settings-content__header {
  position: sticky;
  top: 0;
  z-index: 5;
  display: flex;
  min-height: 76px;
  align-items: center;
  justify-content: space-between;
  gap: 24px;
  background: #fff;
  border-bottom: 1px solid #f0f0f0;
}

.settings-content__title {
  display: flex;
  min-width: 0;
  align-items: center;
  gap: 8px;
}

.settings-content__title strong {
  color: #303030;
  font-size: 18px;
  font-weight: 600;
  line-height: 26px;
}

.settings-content__help {
  display: inline-flex;
  padding: 0;
  align-items: center;
  color: #999;
  font-size: 15px;
  background: transparent;
  border: 0;
  cursor: pointer;

  &:hover,
  &:focus-visible {
    color: var(--color-primary);
  }
}

.settings-form {
  width: min(100%, 900px);
  padding-top: 16px;
}

:deep(.settings-form__grid) {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 0 20px;
}

.settings-form :deep(.ant-input-number),
.settings-form :deep(.ant-select) {
  width: 100%;
}

:deep(.settings-group + .settings-group) {
  margin-top: 36px;
}

:deep(.settings-group__title) {
  position: sticky;
  top: 76px;
  z-index: 4;
  padding: 0 0 14px 12px;
  color: #262626;
  font-size: 16px;
  font-weight: 600;
  line-height: 24px;
  background: #fff;
  border-bottom: 1px solid #f2f2f2;

  &::before {
    position: absolute;
    top: 5px;
    bottom: 19px;
    left: 0;
    width: 3px;
    background: var(--color-primary);
    border-radius: 2px;
    content: '';
  }
}

:deep(.settings-row) {
  display: flex;
  min-height: 76px;
  align-items: center;
  justify-content: space-between;
  gap: 32px;
  border-bottom: 1px solid #f5f5f5;
}

:deep(.settings-row__text) {
  display: flex;
  min-width: 0;
  flex-direction: column;
  gap: 5px;
}

:deep(.settings-row__text strong) {
  color: #333;
  font-size: 14px;
  font-weight: 600;
  line-height: 22px;
}

:deep(.settings-row__text span) {
  color: #8c8c8c;
  font-size: 13px;
  line-height: 20px;
}

:deep(.settings-row__control) {
  flex: 0 0 auto;
}

:deep(.settings-control) {
  width: 220px !important;
}

@media (max-width: 900px) {
  .settings-page {
    height: auto;
    grid-template-columns: 1fr;
    overflow: visible;
  }

  .settings-content {
    padding: 0 18px 24px;
    overflow-y: visible;
  }

  :deep(.settings-form__grid) {
    grid-template-columns: 1fr;
  }
}
</style>
