<template>
  <section class="menu-page">
    <a-form class="toolbar" layout="inline" :model="query" @finish="search">
      <a-form-item>
        <a-input
          v-model:value="query.keyword"
          class="keyword-input"
          allow-clear
          placeholder="菜单名称"
          @press-enter="search"
        >
          <template #prefix><search-outlined /></template>
        </a-input>
      </a-form-item>
      <a-form-item>
        <a-select
          v-model:value="query.types"
          allow-clear
          placeholder="节点类型"
          :options="typeOptions"
        />
      </a-form-item>
      <a-form-item>
        <a-select
          v-model:value="query.status"
          allow-clear
          placeholder="状态"
          :options="statusOptions"
        />
      </a-form-item>
      <a-form-item>
        <a-space :size="8">
          <a-button type="primary" html-type="submit">查询</a-button>
          <a-button @click="reset">重置</a-button>
        </a-space>
      </a-form-item>
      <a-form-item class="create-action">
        <a-button type="primary" @click="createMenu"><plus-outlined />新增</a-button>
      </a-form-item>
    </a-form>

    <a-table
      class="menu-table"
      :columns="columns"
      :data-source="visibleRows"
      :loading="loading"
      :pagination="false"
      row-key="id"
      :scroll="{ x: 'max-content' }"
    >
      <template #headerCell="{ column }">
        <div class="table-header">
          <span>{{ column.title }}</span>
          <span
            v-if="column.resizable !== false"
            class="resize-handle"
            @mousedown.stop.prevent="startResize($event, column.key)"
          />
        </div>
      </template>

      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'title'">
          <div class="menu-name" :style="{ paddingLeft: `${record.depth * 22}px` }">
            <IconPicker :value="record.icon" readonly :size="16" />
            <span class="menu-title">
              <template v-for="(part, index) in getTitleParts(record.title)" :key="index">
                <mark v-if="part.matched" class="title-highlight">{{ part.text }}</mark>
                <span v-else>{{ part.text }}</span>
              </template>
            </span>
            <button
              v-if="record.hasChildren"
              type="button"
              class="expand-button"
              :class="{ 'expand-button--open': !collapsedIds.has(record.id) }"
              :aria-label="collapsedIds.has(record.id) ? '展开' : '收起'"
              @click="toggleRow(record.id)"
            >
              <right-outlined />
            </button>
          </div>
        </template>
        <template v-else-if="column.key === 'types'">
          <a-tag class="node-type-tag" :style="nodeTypeTagStyle(record.types)">
            {{ typeOptionMap.get(record.types)?.label || '-' }}
          </a-tag>
        </template>
        <template v-else-if="column.key === 'route' || column.key === 'code'">
          <span class="path-text">{{ record[column.key] || '-' }}</span>
        </template>
        <template v-else-if="column.key === 'keepAlive'">
          <span v-if="record.types === 2">{{ record.keepAlive === 1 ? '缓存' : '不缓存' }}</span>
          <span v-else class="empty-text">-</span>
        </template>
        <template v-else-if="column.key === 'status'">
          <StatusPulse :status="record.status" />
        </template>
        <template v-else-if="column.key === 'action'">
          <a-space :size="14">
            <a v-if="record.types !== 3" class="action-link" @click="createChild(record)">
              <plus-outlined />添加
            </a>
            <a class="action-link" @click="editMenu(record)"><edit-outlined />编辑</a>
            <a class="danger-link" @click="removeMenu(record)"><delete-outlined />删除</a>
          </a-space>
        </template>
      </template>
    </a-table>
    <MenuForm ref="formRef" :menus="rows" @success="loadMenus" />
  </section>
</template>

<script setup>
import { computed, onBeforeUnmount, onMounted, reactive, ref } from 'vue'
import {
  DeleteOutlined,
  EditOutlined,
  PlusOutlined,
  RightOutlined,
  SearchOutlined,
} from '@ant-design/icons-vue'
import { Modal, message } from 'ant-design-vue'
import IconPicker from '@/components/IconPicker.vue'
import StatusPulse from '@/components/StatusPulse.vue'
import { useDictionaryOptions } from '@/hooks/useDictionary'
import { deleteMenu, getMenuList } from '@/api/menu'
import MenuForm from './components/form.vue'

const baseColumns = [
  { title: '菜单名称', dataIndex: 'title', key: 'title', minWidth: 150 },
  { title: '类型', dataIndex: 'types', key: 'types', minWidth: 70 },
  { title: '路由', dataIndex: 'route', key: 'route', minWidth: 110 },
  { title: '权限编码', dataIndex: 'code', key: 'code', minWidth: 120 },
  { title: '排序', dataIndex: 'sort', key: 'sort', width: 70, minWidth: 70, align: 'center' },
  {
    title: '页面缓存',
    dataIndex: 'keepAlive',
    key: 'keepAlive',
    width: 90,
    minWidth: 90,
    align: 'center',
  },
  { title: '状态', dataIndex: 'status', key: 'status', minWidth: 78 },
  {
    title: '操作',
    key: 'action',
    resizable: false,
  },
]

const columnWidths = reactive({})
const columns = computed(() =>
  baseColumns.map((column) =>
    columnWidths[column.key] ? { ...column, width: columnWidths[column.key] } : column,
  ),
)

const typeOptions = useDictionaryOptions('node_types', true)
const statusOptions = useDictionaryOptions('sys_status', true)
const typeOptionMap = computed(() => new Map(typeOptions.value.map((item) => [item.value, item])))

const nodeTypeColors = {
  default: { color: '#555555', backgroundColor: '#f5f5f5', borderColor: '#e6e6e6' },
  green: { color: '#246b45', backgroundColor: '#f1f8f4', borderColor: '#d5eadc' },
  orange: { color: '#8a5a18', backgroundColor: '#fbf7ee', borderColor: '#eee1c8' },
}

function nodeTypeTagStyle(types) {
  const color = typeOptionMap.value.get(types)?.color || 'default'
  return nodeTypeColors[color] || { color, borderColor: color }
}

const query = reactive({ keyword: '', types: undefined, status: undefined })
const rows = ref([])
const activeKeyword = ref('')
const loading = ref(false)
const formRef = ref()
const collapsedIds = ref(new Set())
let initialized = false
const rowMap = computed(() => new Map(rows.value.map((row) => [row.id, row])))

const visibleRows = computed(() => {
  return rows.value.filter((row) => {
    let parentID = row.pid
    while (parentID !== '0') {
      if (collapsedIds.value.has(parentID)) return false
      parentID = rowMap.value.get(parentID)?.pid || '0'
    }
    return true
  })
})
function toggleRow(id) {
  const next = new Set(collapsedIds.value)
  if (next.has(id)) next.delete(id)
  else next.add(id)
  collapsedIds.value = next
}

function getTitleParts(title) {
  const keyword = activeKeyword.value
  if (!keyword) return [{ text: title, matched: false }]

  const parts = []
  const source = title.toLocaleLowerCase()
  const target = keyword.toLocaleLowerCase()
  let start = 0
  let index = source.indexOf(target, start)
  while (index !== -1) {
    if (index > start) parts.push({ text: title.slice(start, index), matched: false })
    parts.push({ text: title.slice(index, index + keyword.length), matched: true })
    start = index + keyword.length
    index = source.indexOf(target, start)
  }
  if (start < title.length) parts.push({ text: title.slice(start), matched: false })
  return parts.length ? parts : [{ text: title, matched: false }]
}

let resizeState

function startResize(event, key) {
  const column = baseColumns.find((item) => item.key === key)
  const currentWidth = event.currentTarget.closest('th')?.getBoundingClientRect().width
  resizeState = {
    key,
    startX: event.clientX,
    startWidth: currentWidth || column.minWidth,
    minWidth: column.minWidth,
  }
  document.body.classList.add('is-resizing-column')
  window.addEventListener('mousemove', resizeColumn)
  window.addEventListener('mouseup', stopResize)
}

function resizeColumn(event) {
  if (!resizeState) return
  columnWidths[resizeState.key] = Math.max(
    resizeState.minWidth,
    resizeState.startWidth + event.clientX - resizeState.startX,
  )
}

function stopResize() {
  resizeState = undefined
  document.body.classList.remove('is-resizing-column')
  window.removeEventListener('mousemove', resizeColumn)
  window.removeEventListener('mouseup', stopResize)
}

function search() {
  loadMenus()
}

function reset() {
  Object.assign(query, { keyword: '', types: undefined, status: undefined })
  loadMenus()
}

function createMenu() {
  formRef.value?.openCreate()
}

function createChild(parent) {
  formRef.value?.openCreate(parent)
}

function editMenu(record) {
  formRef.value?.openEdit(record.id)
}

function flattenMenus(menus, depth = 0) {
  return menus.flatMap((menu) => {
    const children = menu.children || []
    const row = {
      ...menu,
      depth,
      hasChildren: children.length > 0,
      keepAlive: menu.keep_alive,
    }
    delete row.children
    return [row, ...flattenMenus(children, depth + 1)]
  })
}

async function loadMenus() {
  loading.value = true
  try {
    const data = await getMenuList({
      title: query.keyword || undefined,
      types: query.types,
      status: query.status,
    })
    rows.value = flattenMenus(data?.list || [])
    activeKeyword.value = query.keyword
    if (!initialized) {
      collapsedIds.value = new Set(rows.value.filter((row) => row.hasChildren).map((row) => row.id))
      initialized = true
    }
  } catch (error) {
    message.error(error.message)
  } finally {
    loading.value = false
  }
}

function removeMenu(record) {
  Modal.confirm({
    title: '确认删除菜单？',
    content: `删除“${record.title}”后无法恢复。`,
    okText: '删除',
    okType: 'danger',
    cancelText: '取消',
    async onOk() {
      try {
        await deleteMenu(record.id)
        message.success('菜单删除成功')
        await loadMenus()
      } catch (error) {
        message.error(error.message)
        return Promise.reject(error)
      }
    },
  })
}

onMounted(loadMenus)
onBeforeUnmount(stopResize)
</script>

<style scoped lang="less">
.menu-page {
  min-width: 0;
}

.toolbar {
  min-height: 58px;
  align-items: center;
  padding: 0 0 18px;
}

.toolbar :deep(.ant-form-item) {
  margin-bottom: 0;
}

.toolbar :deep(.create-action) {
  margin-right: 0;
  margin-left: auto;
}

.toolbar :deep(.ant-select) {
  width: 132px;
}

.keyword-input {
  width: 236px;
}

.menu-table {
  width: 100%;
}

.menu-table :deep(.ant-table-thead > tr > th) {
  height: 48px;
  padding: 0 16px;
  color: #494949;
  font-size: 13px;
  font-weight: 600;
  background: #fafafa;
  border-bottom: 1px solid #e8e8e8;
}

.menu-table :deep(.ant-table-tbody > tr > td) {
  height: 52px;
  padding: 10px 16px;
  color: #3b3b3b;
  border-bottom: 1px solid #f0f0f0;
}

.menu-table :deep(.ant-table-thead > tr > th:first-child),
.menu-table :deep(.ant-table-tbody > tr > td:first-child) {
  padding-left: 24px;
}

.menu-table :deep(.ant-table-thead > tr > th:last-child),
.menu-table :deep(.ant-table-tbody > tr > td:last-child) {
  padding-right: 16px;
}

.menu-table :deep(.ant-table-tbody > tr:hover > td) {
  background: #fafafa;
}

.table-header {
  position: relative;
  display: flex;
  align-items: center;
  width: 100%;
  height: 48px;
}

.resize-handle {
  position: absolute;
  top: 9px;
  right: -16px;
  bottom: 9px;
  width: 9px;
  cursor: col-resize;
  z-index: 2;
}

.resize-handle::after {
  position: absolute;
  top: 0;
  right: 4px;
  bottom: 0;
  width: 1px;
  content: '';
  background: transparent;
}

.resize-handle:hover::after {
  background: #bfbfbf;
}

.menu-name {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 0;
  font-weight: 500;
}

.menu-title {
  min-width: 0;
}

.title-highlight {
  padding: 1px 2px;
  color: inherit;
  background: #fde8e6;
  border-radius: 3px;
}

.expand-button {
  width: 22px;
  height: 22px;
  flex: 0 0 auto;
}

.expand-button {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 0;
  color: #8c8c8c;
  font-size: 10px;
  background: transparent;
  border: 0;
  border-radius: 5px;
  cursor: pointer;
  transition:
    color 160ms ease,
    background-color 160ms ease;
}

.expand-button:hover {
  color: #333;
  background: transparent;
}

.expand-button :deep(.anticon) {
  transition: transform 160ms ease;
}

.expand-button--open :deep(.anticon) {
  transform: rotate(90deg);
}

.path-text {
  color: #666;
  font-family: ui-monospace, SFMono-Regular, Consolas, monospace;
  font-size: 12px;
}

.empty-text {
  color: #bfbfbf;
}

.action-link {
  color: #1677ff;
}

.danger-link {
  color: #cf3f36;
}

.action-link :deep(.anticon),
.danger-link :deep(.anticon) {
  margin-right: 4px;
}

@media (max-width: 900px) {
  .toolbar {
    gap: 10px;
    padding-bottom: 16px;
  }

  .toolbar :deep(.ant-form-item) {
    margin: 0;
  }
}
</style>

<style>
body.is-resizing-column {
  cursor: col-resize;
  user-select: none;
}
</style>
