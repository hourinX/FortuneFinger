<template>
  <section
    ref="layoutRef"
    class="dict-page"
    :style="{ '--dict-type-panel-width': `${typePanelWidth}px` }"
  >
    <aside class="dict-type-panel">
      <header class="dict-type-panel__header">
        <div>
          <strong>字典类型</strong>
        </div>
        <a-space :size="8">
          <a-button
            class="dict-type-panel__refresh"
            type="text"
            size="small"
            :loading="typeLoading"
            title="刷新当前页"
            @click="loadTypes"
          >
            <reload-outlined />
          </a-button>
          <a-button type="primary" size="small" @click="typeFormRef?.openCreate()">
            <plus-outlined />新增
          </a-button>
        </a-space>
      </header>

      <div class="dict-type-panel__search">
        <a-input
          v-model:value="typeQuery.name"
          allow-clear
          placeholder="搜索字典名称"
          @input="searchTypes"
        >
          <template #prefix><search-outlined /></template>
        </a-input>
      </div>

      <a-spin :spinning="typeLoading" class="dict-type-panel__loading">
        <div class="dict-type-list">
          <button
            v-for="item in typeRows"
            :key="item.id"
            type="button"
            class="dict-type-item"
            :class="{ active: selectedType?.id === item.id }"
            @click="selectType(item)"
          >
            <span class="dict-type-item__icon">
              <database-outlined v-if="item.is_system === 1" />
              <tags-outlined v-else />
            </span>
            <span class="dict-type-item__identity">
              <span class="dict-type-item__title">
                <strong>{{ item.name }}</strong>
                <a-tag v-if="item.is_system === 1" color="blue">系统字典</a-tag>
              </span>
              <code>{{ item.code }}</code>
            </span>
            <StatusPulse :status="item.status" :show-text="false" />
            <a-dropdown :trigger="['click']" @click.stop>
              <a class="dict-type-item__more" @click.prevent><more-outlined /></a>
              <template #overlay>
                <a-menu>
                  <a-menu-item @click="typeFormRef?.openEdit(item.id)">
                    <edit-outlined />编辑
                  </a-menu-item>
                  <a-menu-item danger :disabled="item.is_system === 1" @click="removeType(item)">
                    <delete-outlined />删除
                  </a-menu-item>
                </a-menu>
              </template>
            </a-dropdown>
          </button>
          <a-empty v-if="!typeLoading && !typeRows.length" description="暂无匹配字典" />
        </div>
      </a-spin>

      <div v-if="typeTotal > typePage.pageSize" class="dict-type-panel__pagination">
        <span>共 {{ typeTotal }} 个字典</span>
        <a-pagination
          v-model:current="typePage.current"
          :page-size="typePage.pageSize"
          :total="typeTotal"
          :show-size-changer="false"
          size="small"
          @change="loadTypes"
        />
      </div>
    </aside>

    <div
      class="dict-panel-resizer"
      role="separator"
      aria-label="调整字典类型和字典值区域宽度"
      @pointerdown="startResize"
    >
      <span />
    </div>

    <main class="dict-value-stage">
      <header class="dict-value-stage__header">
        <div>
          <strong>{{ selectedType?.name || '请选择字典类型' }}</strong>
          <code v-if="selectedType">{{ selectedType.code }}</code>
        </div>
        <a-button
          type="primary"
          :disabled="!selectedType"
          @click="dataFormRef?.openCreate(selectedType)"
        >
          <plus-outlined />添加字典值
        </a-button>
      </header>

      <div class="dict-value-toolbar">
        <a-input
          v-model:value="dataQuery.label"
          allow-clear
          placeholder="搜索标签"
          :disabled="!selectedType"
          @press-enter="searchData"
        >
          <template #prefix><search-outlined /></template>
        </a-input>
        <a-select
          v-model:value="dataQuery.status"
          allow-clear
          placeholder="全部状态"
          :options="statusOptions"
          :disabled="!selectedType"
        />
        <a-button type="primary" ghost :disabled="!selectedType" @click="searchData">筛选</a-button>
        <a-button type="text" :disabled="!selectedType" @click="resetData">清除条件</a-button>
        <span class="dict-value-toolbar__total">{{ dataTotal }} 个值</span>
      </div>

      <a-table
        class="dict-value-table"
        row-key="id"
        :columns="dataColumns"
        :data-source="dataRows"
        :loading="dataLoading"
        :pagination="false"
        :locale="{ emptyText: selectedType ? '当前字典暂无数据' : '请先选择一个字典' }"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'label'">
            <a-tag :color="tagColor(record.color)">{{ record.label }}</a-tag>
          </template>
          <template v-else-if="column.key === 'value'">
            <code>{{ record.value }}</code>
          </template>
          <template v-else-if="column.key === 'is_default'">
            <a-tag v-if="record.is_default === 1" color="green">是</a-tag>
            <span v-else>否</span>
          </template>
          <template v-else-if="column.key === 'status'">
            <StatusPulse :status="record.status" />
          </template>
          <template v-else-if="column.key === 'remark'">
            {{ record.remark || '-' }}
          </template>
          <template v-else-if="column.key === 'action'">
            <a-space :size="12">
              <a @click="dataFormRef?.openEdit(selectedType, record.id, record)">
                <edit-outlined />编辑
              </a>
              <a class="danger-link" @click="removeData(record)"> <delete-outlined />删除 </a>
            </a-space>
          </template>
        </template>
      </a-table>

      <a-pagination
        v-if="dataTotal > dataPage.pageSize"
        v-model:current="dataPage.current"
        v-model:page-size="dataPage.pageSize"
        :total="dataTotal"
        show-size-changer
        :show-total="(total) => '共 ' + total + ' 条'"
        @change="changeDataPage"
      />
    </main>

    <TypeForm ref="typeFormRef" @success="refreshTypes" />
    <DataForm ref="dataFormRef" @success="refreshData" />
  </section>
</template>

<script setup>
import { onBeforeUnmount, onMounted, reactive, ref } from 'vue'
import {
  DatabaseOutlined,
  DeleteOutlined,
  EditOutlined,
  MoreOutlined,
  PlusOutlined,
  ReloadOutlined,
  SearchOutlined,
  TagsOutlined,
} from '@ant-design/icons-vue'
import { Modal, message } from 'ant-design-vue'
import { deleteDictData, deleteDictType, getDictDataList, getDictTypeList } from '@/api/dict'
import StatusPulse from '@/components/StatusPulse.vue'
import { useDictionaryOptions } from '@/hooks/useDictionary'
import { useAppStore } from '@/store'
import DataForm from './components/data-form.vue'
import TypeForm from './components/type-form.vue'

const statusOptions = useDictionaryOptions('sys_status', true)
const app = useAppStore()
const dataColumns = [
  { title: '标签', key: 'label', dataIndex: 'label', width: 140 },
  { title: '字典值', key: 'value', dataIndex: 'value', width: 150 },
  { title: '默认', key: 'is_default', dataIndex: 'is_default', width: 80 },
  { title: '排序', key: 'sort', dataIndex: 'sort', width: 80 },
  { title: '状态', key: 'status', dataIndex: 'status', width: 100 },
  { title: '备注', key: 'remark', dataIndex: 'remark' },
  { title: '操作', key: 'action', width: 140 },
]
const typeQuery = reactive({ name: '' })
const typePage = reactive({ current: 1, pageSize: 10 })
const typeRows = ref([])
const typeTotal = ref(0)
const typeLoading = ref(false)
const selectedType = ref()
const typeFormRef = ref()
const dataQuery = reactive({ label: '', status: undefined })
const dataPage = reactive({ current: 1, pageSize: 10 })
const dataRows = ref([])
const dataTotal = ref(0)
const dataLoading = ref(false)
const dataFormRef = ref()
const layoutRef = ref()
const typePanelWidth = ref(340)
let resizeStartX = 0
let resizeStartWidth = 0

function resizePanels(event) {
  const containerWidth = layoutRef.value?.clientWidth || 0
  const maxWidth = Math.max(300, Math.min(520, containerWidth - 540))
  typePanelWidth.value = Math.min(
    maxWidth,
    Math.max(300, resizeStartWidth + event.clientX - resizeStartX),
  )
}

function stopResize() {
  window.removeEventListener('pointermove', resizePanels)
  window.removeEventListener('pointerup', stopResize)
  document.body.style.removeProperty('cursor')
  document.body.style.removeProperty('user-select')
}

function startResize(event) {
  resizeStartX = event.clientX
  resizeStartWidth = typePanelWidth.value
  document.body.style.cursor = 'col-resize'
  document.body.style.userSelect = 'none'
  window.addEventListener('pointermove', resizePanels)
  window.addEventListener('pointerup', stopResize)
}

function tagColor(color) {
  return color === 'default' ? undefined : color
}

function searchTypes() {
  typePage.current = 1
  loadTypes()
}

function selectType(item) {
  if (selectedType.value?.id === item.id) return
  selectedType.value = item
  dataPage.current = 1
  Object.assign(dataQuery, { label: '', status: undefined })
  loadData()
}

async function loadTypes() {
  typeLoading.value = true
  try {
    const resp = await getDictTypeList({
      page: typePage.current,
      pageSize: typePage.pageSize,
      name: typeQuery.name || undefined,
    })
    typeRows.value = resp.list || []
    typeTotal.value = resp.total || 0
    const active = typeRows.value.find((item) => item.id === selectedType.value?.id)
    if (active) selectedType.value = active
    if (!selectedType.value && typeRows.value.length) {
      selectedType.value = typeRows.value[0]
      dataPage.current = 1
    }
    await loadData()
  } catch (error) {
    message.error(error.message)
  } finally {
    typeLoading.value = false
  }
}

function searchData() {
  dataPage.current = 1
  loadData()
}

function resetData() {
  Object.assign(dataQuery, { label: '', status: undefined })
  searchData()
}

function changeDataPage(current, pageSize) {
  dataPage.current = current
  dataPage.pageSize = pageSize
  loadData()
}

async function loadData() {
  if (!selectedType.value) {
    dataRows.value = []
    dataTotal.value = 0
    return
  }
  dataLoading.value = true
  try {
    const resp = await getDictDataList({
      page: dataPage.current,
      pageSize: dataPage.pageSize,
      dict_code: selectedType.value.code,
      label: dataQuery.label || undefined,
      status: dataQuery.status,
    })
    dataRows.value = resp.list || []
    dataTotal.value = resp.total || 0
  } catch (error) {
    message.error(error.message)
  } finally {
    dataLoading.value = false
  }
}

async function refreshTypes() {
  await loadTypes()
  await app.loadDictionaries(true)
}

async function refreshData() {
  await loadData()
  await app.loadDictionaries(true)
}

function removeType(record) {
  Modal.confirm({
    title: '确认删除字典类型？',
    content: `删除“${record.name}”后无法恢复。`,
    okText: '删除',
    okType: 'danger',
    cancelText: '取消',
    async onOk() {
      try {
        await deleteDictType(record.id)
        message.success('字典类型删除成功')
        if (selectedType.value?.id === record.id) selectedType.value = undefined
        if (typeRows.value.length === 1 && typePage.current > 1) typePage.current -= 1
        await refreshTypes()
      } catch (error) {
        message.error(error.message)
        return Promise.reject(error)
      }
    },
  })
}

function removeData(record) {
  Modal.confirm({
    title: '确认删除字典数据？',
    content: `删除“${record.label}”后无法恢复。`,
    okText: '删除',
    okType: 'danger',
    cancelText: '取消',
    async onOk() {
      try {
        await deleteDictData(record.id)
        message.success('字典数据删除成功')
        if (dataRows.value.length === 1 && dataPage.current > 1) dataPage.current -= 1
        await refreshData()
      } catch (error) {
        message.error(error.message)
        return Promise.reject(error)
      }
    },
  })
}

onMounted(loadTypes)
onBeforeUnmount(stopResize)
</script>

<style scoped lang="less">
.dict-page {
  display: grid;
  min-height: calc(100vh - 132px);
  grid-template-columns: var(--dict-type-panel-width) 14px minmax(0, 1fr);
  align-items: stretch;

  .dict-type-panel,
  .dict-value-stage {
    min-width: 0;
    overflow: hidden;
    border: 1px solid #ece9e8;
    border-radius: 12px;
    box-shadow: 0 5px 20px rgba(48, 38, 35, 0.03);
  }

  .dict-type-panel {
    display: flex;
    max-height: calc(100vh - 132px);
    flex-direction: column;
    background: #faf9f9;
  }

  .dict-type-panel__header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 16px;
    background: #fff;
    border-bottom: 1px solid #efedec;

    strong {
      color: #35302e;
      font-size: 15px;
      font-weight: 650;
    }
  }

  .dict-type-panel__search {
    padding: 12px 12px 8px;

    :deep(.ant-input-affix-wrapper) {
      background: #fff;
    }
  }

  .dict-type-panel__loading {
    min-height: 0;
    flex: 1;

    :deep(.ant-spin-container) {
      height: 100%;
    }
  }

  .dict-type-list {
    height: 100%;
    padding: 2px 8px 12px;
    overflow-y: auto;
    scrollbar-color: #d9d3d1 transparent;
    scrollbar-width: thin;
  }

  .dict-type-panel__pagination {
    display: flex;
    flex: 0 0 auto;
    flex-wrap: wrap;
    align-items: center;
    gap: 3px 8px;
    justify-content: flex-end;
    padding: 11px 8px;

    > span {
      flex: 0 0 auto;
      color: #99918e;
      font-size: 12px;
      line-height: 16px;
      white-space: nowrap;
    }

    :deep(.ant-pagination) {
      flex: 0 0 auto;
    }
  }

  .dict-type-item {
    display: grid;
    width: 100%;
    min-width: 0;
    grid-template-columns: 32px minmax(0, 1fr) 10px 18px;
    align-items: center;
    gap: 8px;
    margin: 2px 0;
    padding: 9px 8px;
    color: #77706d;
    text-align: left;
    background: transparent;
    border: 0;
    border-radius: 9px;
    cursor: pointer;
    transition: background-color 150ms ease;

    &:hover {
      background: #f1eeee;
    }

    &.active {
      background: #fff;
      box-shadow: 0 3px 12px rgba(54, 43, 39, 0.07);

      .dict-type-item__icon {
        color: #fff;
        background: var(--color-primary);
      }

      .dict-type-item__identity strong {
        color: var(--color-primary);
      }
    }
  }

  .dict-type-item__icon {
    display: inline-flex;
    width: 32px;
    height: 32px;
    align-items: center;
    justify-content: center;
    color: #817976;
    background: #ebe8e7;
    border-radius: 8px;
  }

  .dict-type-item__identity {
    min-width: 0;

    strong,
    code {
      display: block;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }

    strong {
      color: #49423f;
      font-size: 12px;
      font-weight: 550;
    }

    code {
      margin-top: 2px;
      color: #aaa3a0;
      font-size: 9px;
    }
  }

  .dict-type-item__title {
    display: flex;
    min-width: 0;
    align-items: center;
    gap: 5px;

    strong {
      min-width: 0;
    }

    :deep(.ant-tag) {
      flex: 0 0 auto;
      margin: 0;
      padding: 0 4px;
      font-size: 9px;
      line-height: 16px;
    }
  }

  .dict-type-item__more {
    color: #a49d9a;

    &:hover {
      color: #514946;
    }
  }

  .dict-panel-resizer {
    position: relative;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: col-resize;
    touch-action: none;

    &::before {
      position: absolute;
      inset: 0 -3px;
      content: '';
    }

    span {
      width: 2px;
      height: 46px;
      background: #ddd8d6;
      border-radius: 2px;
      transition: all 150ms ease;
    }

    &:hover span,
    &:active span {
      width: 3px;
      height: 62px;
      background: var(--color-primary);
    }
  }

  .dict-value-stage {
    background: #fff;
  }

  .dict-value-stage__header {
    display: flex;
    min-height: 68px;
    align-items: center;
    justify-content: space-between;
    padding: 14px 16px;
    border-bottom: 1px solid #efedec;

    > div {
      display: flex;
      align-items: baseline;
      gap: 9px;
    }

    strong {
      color: #35302e;
      font-size: 16px;
      font-weight: 650;
    }

    code {
      color: #99918e;
      font-size: 11px;
    }
  }

  .dict-value-toolbar {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 14px 16px;

    :deep(.ant-input-affix-wrapper) {
      width: 220px;
    }

    :deep(.ant-select) {
      width: 130px;
    }
  }

  .dict-value-toolbar__total {
    margin-left: auto;
    color: #a19a97;
    font-size: 12px;
  }

  .dict-value-table {
    padding: 0 16px;

    code {
      color: #5d5552;
      font-size: 12px;
    }
  }

  .dict-value-stage > .ant-pagination {
    display: flex;
    justify-content: flex-end;
    margin: 16px;
  }

  .danger-link {
    color: #cf3f36;
  }

  .dict-value-table a:not(.danger-link) {
    color: #1677ff;
  }
}

@media (max-width: 900px) {
  .dict-page {
    grid-template-columns: minmax(280px, var(--dict-type-panel-width)) 14px minmax(420px, 1fr);
  }
}
</style>
