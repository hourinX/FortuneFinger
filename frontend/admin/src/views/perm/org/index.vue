<template>
  <section class="org-page">
    <a-form class="org-toolbar" layout="inline" :model="query" @finish="search">
      <a-form-item>
        <a-input v-model:value="query.name" allow-clear placeholder="组织名称">
          <template #prefix><search-outlined /></template>
        </a-input>
      </a-form-item>
      <a-form-item
        ><a-input v-model:value="query.code" allow-clear placeholder="组织编码"
      /></a-form-item>
      <a-form-item>
        <a-select
          v-model:value="query.types"
          allow-clear
          placeholder="组织类型"
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
      <a-form-item class="org-toolbar__create">
        <a-button type="primary" @click="formRef?.openCreate()"><plus-outlined />新增</a-button>
      </a-form-item>
    </a-form>

    <a-table
      class="org-table"
      :columns="columns"
      :data-source="visibleRows"
      :loading="loading"
      :pagination="false"
      row-key="id"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'name'">
          <div class="org-name" :style="{ paddingLeft: `${record.depth * 22}px` }">
            <cluster-outlined v-if="record.types === 1" class="org-icon org-icon--group" />
            <bank-outlined v-else-if="record.types === 2" class="org-icon org-icon--company" />
            <apartment-outlined v-else class="org-icon org-icon--department" />
            <span class="org-title">
              <template v-for="(part, index) in getNameParts(record)" :key="index">
                <mark v-if="part.matched" class="org-title-highlight">{{ part.text }}</mark>
                <span v-else>{{ part.text }}</span>
              </template>
            </span>
            <button
              v-if="record.hasChildren"
              type="button"
              class="org-expand-button"
              :class="{ 'org-expand-button--open': !collapsedIds.has(record.id) }"
              :aria-label="collapsedIds.has(record.id) ? '展开' : '收起'"
              @click="toggleRow(record.id)"
            >
              <right-outlined />
            </button>
          </div>
        </template>
        <template v-else-if="column.key === 'code'">
          <span class="org-code">{{ record.code }}</span>
        </template>
        <template v-else-if="column.key === 'types'">
          <span class="type-label" :class="`type-label--${record.types}`">
            {{ typeLabel(record.types) }}
          </span>
        </template>
        <template v-else-if="column.key === 'status'">
          <StatusPulse :status="record.status" />
        </template>
        <template v-else-if="column.key === 'action'">
          <a-space :size="14">
            <a v-if="record.types !== 3" class="action-link" @click="formRef?.openCreate(record)">
              <plus-outlined />添加
            </a>
            <a class="action-link" @click="formRef?.openEdit(record)"><edit-outlined />编辑</a>
            <a class="danger-link" @click="removeOrg(record)"><delete-outlined />删除</a>
          </a-space>
        </template>
      </template>
      <template #emptyText><a-empty description="暂无组织数据" /></template>
    </a-table>

    <OrgForm ref="formRef" :organizations="rows" @success="loadOrganizations" />
  </section>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import {
  ApartmentOutlined,
  BankOutlined,
  ClusterOutlined,
  DeleteOutlined,
  EditOutlined,
  PlusOutlined,
  RightOutlined,
  SearchOutlined,
} from '@ant-design/icons-vue'
import { Modal, message } from 'ant-design-vue'
import StatusPulse from '@/components/StatusPulse.vue'
import { deleteOrg, getOrgList } from '@/api/org'
import OrgForm from './components/form.vue'

const columns = [
  { title: '组织名称', dataIndex: 'name', key: 'name' },
  { title: '组织编码', dataIndex: 'code', key: 'code', width: 210 },
  { title: '类型', dataIndex: 'types', key: 'types', width: 120 },
  { title: '排序', dataIndex: 'sort', key: 'sort', width: 100 },
  { title: '状态', dataIndex: 'status', key: 'status', width: 120 },
  { title: '操作', key: 'action', width: 250 },
]
const typeOptions = [
  { label: '集团', value: 1 },
  { label: '公司', value: 2 },
  { label: '部门', value: 3 },
]
const statusOptions = [
  { label: '启用', value: 1 },
  { label: '禁用', value: 2 },
]
const query = reactive({ name: '', code: '', types: undefined, status: undefined })
const activeQuery = reactive({ name: '', code: '', types: undefined, status: undefined })
const rows = ref([])
const loading = ref(false)
const formRef = ref()
const collapsedIds = ref(new Set())
const rowMap = computed(() => new Map(rows.value.map((item) => [item.id, item])))
const visibleRows = computed(() =>
  rows.value.filter((row) => {
    let parentID = row.parentId
    while (parentID !== '0') {
      if (collapsedIds.value.has(parentID)) return false
      parentID = rowMap.value.get(parentID)?.parentId || '0'
    }
    return true
  }),
)
let initialized = false

function flattenOrganizations(nodes, depth = 0) {
  return nodes.flatMap((node) => {
    const children = node.children || []
    const row = {
      ...node,
      parentId: node.parent_id,
      depth,
      hasChildren: children.length > 0,
    }
    delete row.children
    return [row, ...flattenOrganizations(children, depth + 1)]
  })
}

function toggleRow(id) {
  const next = new Set(collapsedIds.value)
  if (next.has(id)) next.delete(id)
  else next.add(id)
  collapsedIds.value = next
}

function typeLabel(types) {
  return typeOptions.find((item) => item.value === types)?.label || '组织'
}

function isMatchedOrganization(record) {
  const hasQuery =
    activeQuery.name ||
    activeQuery.code ||
    activeQuery.types !== undefined ||
    activeQuery.status !== undefined
  if (!hasQuery) return false

  return (
    (!activeQuery.name ||
      record.name.toLocaleLowerCase().includes(activeQuery.name.toLocaleLowerCase())) &&
    (!activeQuery.code || record.code === activeQuery.code) &&
    (activeQuery.types === undefined || record.types === activeQuery.types) &&
    (activeQuery.status === undefined || record.status === activeQuery.status)
  )
}

function getNameParts(record) {
  if (!isMatchedOrganization(record)) return [{ text: record.name, matched: false }]
  const keyword = activeQuery.name
  if (!keyword) return [{ text: record.name, matched: true }]

  const parts = []
  const source = record.name.toLocaleLowerCase()
  const target = keyword.toLocaleLowerCase()
  let start = 0
  let index = source.indexOf(target, start)
  while (index !== -1) {
    if (index > start) parts.push({ text: record.name.slice(start, index), matched: false })
    parts.push({ text: record.name.slice(index, index + keyword.length), matched: true })
    start = index + keyword.length
    index = source.indexOf(target, start)
  }
  if (start < record.name.length) {
    parts.push({ text: record.name.slice(start), matched: false })
  }
  return parts
}

function search() {
  loadOrganizations()
}

function reset() {
  Object.assign(query, { name: '', code: '', types: undefined, status: undefined })
  loadOrganizations()
}

async function loadOrganizations() {
  loading.value = true
  try {
    const data = await getOrgList({
      name: query.name || undefined,
      code: query.code || undefined,
      types: query.types,
      status: query.status,
    })
    rows.value = flattenOrganizations(data?.list || [])
    Object.assign(activeQuery, query)
    if (!initialized) {
      collapsedIds.value = new Set(rows.value.filter((row) => row.hasChildren).map((row) => row.id))
      initialized = true
    } else if (hasQueryConditions()) {
      const next = new Set(collapsedIds.value)
      rows.value.filter((row) => row.hasChildren).forEach((row) => next.delete(row.id))
      collapsedIds.value = next
    }
  } catch (error) {
    message.error(error.message)
  } finally {
    loading.value = false
  }
}

function hasQueryConditions() {
  return Boolean(
    query.name || query.code || query.types !== undefined || query.status !== undefined,
  )
}

function removeOrg(record) {
  Modal.confirm({
    title: '确认删除组织？',
    content: `删除“${record.name}”后无法恢复。`,
    okText: '删除',
    okType: 'danger',
    cancelText: '取消',
    async onOk() {
      try {
        await deleteOrg(record.id)
        message.success('组织删除成功')
        await loadOrganizations()
      } catch (error) {
        message.error(error.message)
        return Promise.reject(error)
      }
    },
  })
}

onMounted(loadOrganizations)
</script>

<style scoped lang="less">
.org-page {
  min-width: 0;
}

.org-toolbar {
  min-height: 58px;
  align-items: center;
  padding: 0 0 18px;

  :deep(.ant-form-item) {
    margin-bottom: 0;
  }

  :deep(.ant-input) {
    width: 190px;
  }

  :deep(.ant-select) {
    width: 132px;
  }
}

.org-toolbar__create {
  margin-right: 0 !important;
  margin-left: auto !important;
}

.org-table {
  width: 100%;

  :deep(.ant-table-thead > tr > th) {
    height: 48px;
    padding: 0 16px;
    color: #494949;
    font-size: 13px;
    font-weight: 600;
    background: #fafafa;
    border-bottom: 1px solid #e8e8e8;
  }

  :deep(.ant-table-tbody > tr > td) {
    height: 52px;
    padding: 10px 16px;
    color: #3b3b3b;
    border-bottom: 1px solid #f0f0f0;
  }

  :deep(.ant-table-thead > tr > th:first-child),
  :deep(.ant-table-tbody > tr > td:first-child) {
    padding-left: 24px;
  }

  :deep(.ant-table-thead > tr > th:last-child),
  :deep(.ant-table-tbody > tr > td:last-child) {
    padding-right: 16px;
  }

  :deep(.ant-table-tbody > tr:hover > td) {
    background: #fafafa;
  }
}

.org-name {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 0;
  font-weight: 500;
}

.org-title {
  min-width: 0;
}

.org-title-highlight {
  padding: 1px 2px;
  color: inherit;
  background: #fde8e6;
  border-radius: 3px;
}

.org-icon {
  font-size: 16px;
}

.org-icon--group {
  color: #d14b43;
}

.org-icon--company {
  color: #5f76a4;
}

.org-icon--department {
  color: #758272;
}

.org-expand-button {
  width: 22px;
  height: 22px;
  flex: 0 0 auto;
}

.org-expand-button {
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

.org-expand-button:hover {
  color: #333;
  background: transparent;
}

.org-expand-button :deep(.anticon) {
  transition: transform 160ms ease;
}

.org-expand-button--open :deep(.anticon) {
  transform: rotate(90deg);
}

.type-label {
  display: inline-flex;
  height: 24px;
  align-items: center;
  padding: 0 9px;
  color: #555;
  font-size: 12px;
  background: #f5f5f5;
  border: 1px solid #e6e6e6;
  border-radius: 4px;
}

.type-label--2 {
  color: #246b45;
  background: #f1f8f4;
  border-color: #d5eadc;
}

.type-label--3 {
  color: #8a5a18;
  background: #fbf7ee;
  border-color: #eee1c8;
}

.org-code {
  color: #666;
  font-family: ui-monospace, SFMono-Regular, Consolas, monospace;
  font-size: 12px;
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
</style>
