<template>
  <span v-if="readonly" class="icon-display" :style="iconStyle">
    <component :is="currentIcon" />
  </span>
  <a-popover
    v-else
    v-model:open="open"
    trigger="click"
    placement="bottomLeft"
    overlay-class-name="icon-picker-popover"
  >
    <template #content>
      <div class="icon-panel">
        <a-input v-model:value="keyword" allow-clear placeholder="Search icon code">
          <template #prefix><search-outlined /></template>
        </a-input>
        <div v-if="filteredIcons.length" class="icon-grid">
          <button
            v-for="item in filteredIcons"
            :key="item.code"
            type="button"
            class="icon-option"
            :class="{ 'icon-option--selected': item.code === value }"
            :title="item.code"
            @click="selectIcon(item.code)"
          >
            <component :is="item.component" class="icon-option__icon" />
            <span class="icon-option__code">{{ item.code }}</span>
          </button>
        </div>
        <a-empty v-else :image="simpleImage" description="No matching icons" />
      </div>
    </template>
    <a-input class="picker-trigger" :value="currentCode" readonly placeholder="Select icon">
      <template #prefix><component :is="currentIcon" /></template>
      <template #suffix><down-outlined /></template>
    </a-input>
  </a-popover>
</template>

<script setup>
import { computed, ref } from 'vue'
import { Empty } from 'ant-design-vue'
import * as AntIcons from '@ant-design/icons-vue'

const props = defineProps({
  value: { type: String, default: 'appstore' },
  readonly: { type: Boolean, default: false },
  size: { type: Number, default: 16 },
})
const emit = defineEmits(['update:value'])

const DownOutlined = AntIcons.DownOutlined
const SearchOutlined = AntIcons.SearchOutlined
const iconNames = [
  'Home',
  'Dashboard',
  'Appstore',
  'Menu',
  'Bars',
  'UnorderedList',
  'OrderedList',
  'Profile',

  'Setting',
  'Tool',
  'Control',
  'Sliders',
  'Build',
  'Skin',

  'User',
  'Team',
  'Contacts',
  'Idcard',
  'Solution',
  'Apartment',
  'Cluster',

  'Safety',
  'SafetyCertificate',
  'Lock',
  'Unlock',
  'Key',
  'Audit',
  'Exception',

  'Project',
  'Schedule',
  'Calendar',
  'ClockCircle',
  'History',

  'Folder',
  'FolderOpen',
  'File',
  'FileText',
  'FileAdd',
  'Database',
  'Hdd',
  'Cloud',
  'CloudServer',
  'Inbox',

  'Shop',
  'Shopping',
  'ShoppingCart',
  'Bank',
  'Wallet',
  'PayCircle',
  'CreditCard',
  'AccountBook',
  'Transaction',

  'Gift',
  'Trophy',
  'Crown',
  'Star',
  'Heart',
  'Fire',
  'Thunderbolt',
  'Rocket',
  'Flag',
  'Gold',

  'Tag',
  'Tags',
  'Bell',
  'Notification',
  'Message',
  'Mail',
  'Phone',
  'Mobile',

  'Global',
  'Environment',
  'Compass',
  'Car',

  'DeploymentUnit',
  'Branches',
  'Link',
  'Api',
  'Gateway',
  'Container',
  'Code',
  'Bug',
  'Experiment',

  'Picture',
  'Camera',
  'VideoCamera',
  'PlayCircle',
  'CustomerService',
  'Read',
  'Book',

  'Fund',
  'LineChart',
  'BarChart',
  'PieChart',
  'AreaChart',
  'RadarChart',
]
const iconOptions = [
  { code: 'company', component: AntIcons.BankOutlined },
  ...iconNames
    .map((name) => ({
      code: name.replace(/([a-z0-9])([A-Z])/g, '$1-$2').toLowerCase(),
      component: AntIcons[`${name}Outlined`],
    }))
    .filter((item) => item.component),
]

const open = ref(false)
const keyword = ref('')
const simpleImage = Empty.PRESENTED_IMAGE_SIMPLE
const current = computed(
  () => iconOptions.find((item) => item.code === props.value) || iconOptions[0],
)
const currentCode = computed(() => current.value.code)
const currentIcon = computed(() => current.value.component)
const iconStyle = computed(() => ({ fontSize: `${props.size}px` }))
const filteredIcons = computed(() => {
  const text = keyword.value.trim().toLowerCase()
  if (!text) return iconOptions
  return iconOptions.filter((item) => item.code.toLowerCase().includes(text))
})

function selectIcon(code) {
  emit('update:value', code)
  open.value = false
  keyword.value = ''
}
</script>

<style scoped>
.icon-display {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  color: #777;
}

.picker-trigger {
  width: 260px;
  height: 38px;
  cursor: pointer;
}

.picker-trigger :deep(input) {
  color: #4f596b;
  font-size: 14px;
  cursor: pointer;
}

.picker-trigger :deep(.ant-input-prefix) {
  width: 26px;
  height: 26px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  flex: none;
  margin-inline-end: 10px;
  color: var(--color-primary);
  background: var(--color-primary-light);
  border-radius: 7px;
}
</style>

<style>
.icon-picker-popover .ant-popover-inner {
  padding: 14px;
  border-radius: 10px;
}

.icon-picker-popover .icon-panel {
  width: 392px;
}

.icon-picker-popover .icon-grid {
  max-height: 310px;
  display: grid;
  grid-template-columns: repeat(5, minmax(0, 1fr));
  gap: 7px;
  margin-top: 12px;
  padding-right: 3px;
  overflow-y: auto;
}

.icon-picker-popover .icon-option {
  min-width: 0;
  height: 62px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: column;
  gap: 6px;
  padding: 6px 3px;
  color: #5e687b;
  background: #fff;
  border: 1px solid #edf0f5;
  border-radius: 9px;
  cursor: pointer;
  transition:
    color 150ms ease,
    border-color 150ms ease,
    background-color 150ms ease,
    transform 150ms ease;
}

.icon-picker-popover .icon-option__icon {
  flex: 0 0 auto;
  color: inherit;
  font-size: 18px;
}

.icon-picker-popover .icon-option__code {
  width: 100%;
  overflow: hidden;
  color: #657086;
  font-size: 10px;
  line-height: 12px;
  text-align: center;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.icon-picker-popover .icon-option:hover {
  color: var(--color-primary);
  background: var(--color-primary-light);
  border-color: var(--color-primary-border);
  transform: translateY(-1px);
}

.icon-picker-popover .icon-option--selected {
  color: var(--color-primary);
  background: var(--color-primary-light);
  border-color: var(--color-primary);
  box-shadow: 0 0 0 2px rgb(228 61 54 / 8%);
}

.icon-picker-popover .ant-empty {
  margin: 34px 0 18px;
}
</style>
