<template>
  <section class="prize-page">
    <a-form class="prize-toolbar" layout="inline" :model="query" @finish="search">
      <a-form-item>
        <a-input v-model:value="query.name" allow-clear placeholder="奖品名称">
          <template #prefix><search-outlined /></template>
        </a-input>
      </a-form-item>
      <a-form-item>
        <a-select
          v-model:value="query.level"
          allow-clear
          placeholder="奖项等级"
          :options="levelOptions"
        />
      </a-form-item>
      <a-form-item>
        <a-select
          v-model:value="query.types"
          allow-clear
          placeholder="奖品类型"
          :options="kindOptions"
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
        <a-space v-permission="'prize:pool:list'" :size="8">
          <a-button type="primary" html-type="submit">查询</a-button>
          <a-button @click="reset">重置</a-button>
        </a-space>
      </a-form-item>
      <a-form-item class="prize-toolbar__create">
        <a-button v-permission="'prize:pool:add'" type="primary" @click="openCreate">
          <plus-outlined />新增奖品
        </a-button>
      </a-form-item>
    </a-form>

    <a-table
      class="prize-table"
      :columns="columns"
      :data-source="rows"
      :loading="loading"
      :pagination="pagination"
      row-key="id"
      :scroll="{ x: 850 }"
      @change="changeTable"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'prize'">
          <div class="prize-identity">
            <div
              class="prize-cover"
              :class="{ 'prize-cover--previewable': record.image_url }"
              @click="previewPrizeImage(record.image_url)"
            >
              <img v-if="record.image_url" :src="record.image_url" :alt="record.name" />
              <gift-outlined v-else />
            </div>
            <div>
              <strong>{{ record.name }}</strong>
              <span>{{ kindOptionMap.get(record.types)?.label || record.types }}</span>
            </div>
          </div>
        </template>
        <template v-else-if="column.key === 'level'">
          <span v-if="!levelOptionMap.get(record.level)?.color">
            {{ levelOptionMap.get(record.level)?.label || record.level }}
          </span>
          <a-tag v-else :style="prizeLevelTagStyle(record.level)">
            {{ levelOptionMap.get(record.level)?.label || record.level }}
          </a-tag>
        </template>
        <template v-else-if="column.key === 'quantity'">
          <span>{{ record.quantity }} 份</span>
        </template>
        <template v-else-if="column.key === 'status'">
          <StatusPulse :status="record.status" />
        </template>
        <template v-else-if="column.key === 'action'">
          <a-space :size="14">
            <a v-permission="'prize:pool:edit'" @click="openEdit(record)">
              <edit-outlined />编辑
            </a>
            <a
              v-permission="'prize:pool:del'"
              class="danger-link"
              @click="removePrize(record)"
            >
              <delete-outlined />删除
            </a>
          </a-space>
        </template>
      </template>
      <template #emptyText><a-empty description="暂无奖品数据" /></template>
    </a-table>
    <PoolForm ref="formRef" @success="loadPrizes" />
    <ImagePreview v-model:visible="imagePreviewVisible" :src="imagePreviewSrc" alt="奖品图片预览" />
  </section>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import {
  DeleteOutlined,
  EditOutlined,
  GiftOutlined,
  PlusOutlined,
  SearchOutlined,
} from '@ant-design/icons-vue'
import { Modal, message } from 'ant-design-vue'
import { deletePrize, getPrizeList } from '@/api/prize'
import ImagePreview from '@/components/ImagePreview.vue'
import StatusPulse from '@/components/StatusPulse.vue'
import { useDictionaryOptions } from '@/hooks/useDictionary'
import { dictionaryTagStyleByCode } from '@/utils/dictionary-tag'
import PoolForm from './components/form.vue'

const levelOptions = useDictionaryOptions('prize_level')
const kindOptions = useDictionaryOptions('prize_types', true)
const statusOptions = useDictionaryOptions('sys_status', true)
const levelOptionMap = computed(
  () => new Map(levelOptions.value.map((option) => [option.value, option])),
)
const kindOptionMap = computed(
  () => new Map(kindOptions.value.map((option) => [option.value, option])),
)
function prizeLevelTagStyle(level) {
  return dictionaryTagStyleByCode('prize_level', levelOptionMap.value.get(level))
}
const columns = [
  { title: '奖品信息', key: 'prize', width: 270 },
  { title: '归属企业', dataIndex: 'org_name', key: 'orgName', width: 180 },
  { title: '奖项等级', key: 'level', width: 110 },
  { title: '奖品数量', key: 'quantity', width: 110 },
  { title: '状态', key: 'status', width: 95 },
  { title: '操作', key: 'action', width: 170 },
]
const rows = ref([])
const loading = ref(false)
const query = reactive({ name: '', level: undefined, types: undefined, status: undefined })
const page = reactive({ current: 1, pageSize: 10, total: 0 })
const formRef = ref()
const imagePreviewVisible = ref(false)
const imagePreviewSrc = ref('')
const pagination = computed(() => ({
  current: page.current,
  pageSize: page.pageSize,
  total: page.total,
  showSizeChanger: true,
  showTotal: (total) => `共 ${total} 条`,
}))

function search() {
  page.current = 1
  loadPrizes()
}

function reset() {
  Object.assign(query, { name: '', level: undefined, types: undefined, status: undefined })
  search()
}

function openCreate() {
  formRef.value?.openCreate()
}

function openEdit(item) {
  formRef.value?.openEdit(item.id)
}

function previewPrizeImage(src) {
  if (!src) return
  imagePreviewSrc.value = src
  imagePreviewVisible.value = true
}

function removePrize(item) {
  Modal.confirm({
    title: '确认删除奖品？',
    content: `删除“${item.name}”后无法恢复。`,
    okText: '删除',
    okType: 'danger',
    cancelText: '取消',
    async onOk() {
      try {
        await deletePrize(item.id)
        message.success('奖品删除成功')
        if (rows.value.length === 1 && page.current > 1) page.current -= 1
        await loadPrizes()
      } catch (error) {
        message.error(error.message)
        return Promise.reject(error)
      }
    },
  })
}

function changeTable(paginationState) {
  page.current = paginationState.current
  page.pageSize = paginationState.pageSize
  loadPrizes()
}

async function loadPrizes() {
  loading.value = true
  try {
    const resp = await getPrizeList({
      page: page.current,
      pageSize: page.pageSize,
      name: query.name || undefined,
      level: query.level,
      types: query.types,
      status: query.status,
    })
    rows.value = resp.list || []
    page.total = resp.total || 0
  } catch (error) {
    message.error(error.message)
  } finally {
    loading.value = false
  }
}

onMounted(loadPrizes)
</script>

<style scoped lang="less">
.prize-page {
  min-width: 0;
}

.prize-toolbar {
  display: flex;
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

.prize-toolbar__create {
  margin-right: 0 !important;
  margin-left: auto !important;
}

.prize-table {
  :deep(.ant-table-thead > tr > th) {
    height: 48px;
    padding: 0 16px;
    color: #494949;
    font-size: 13px;
    font-weight: 600;
    background: #fafafa;
  }

  :deep(.ant-table-tbody > tr > td) {
    height: 72px;
    padding: 10px 16px;
  }

  :deep(.ant-tag) {
    margin: 0;
  }

  a {
    color: #1677ff;
  }

  .danger-link {
    color: #cf3f36;
  }

  :deep(.anticon) {
    margin-right: 4px;
  }
}

.prize-identity {
  display: flex;
  align-items: center;
  gap: 12px;

  > div:last-child {
    display: flex;
    min-width: 0;
    flex-direction: column;
    gap: 4px;
  }

  strong {
    overflow: hidden;
    color: #303030;
    font-size: 14px;
    white-space: nowrap;
    text-overflow: ellipsis;
  }

  span {
    color: #999;
    font-size: 12px;
  }
}

.prize-cover {
  display: flex;
  width: 44px;
  height: 44px;
  flex: 0 0 auto;
  align-items: center;
  justify-content: center;
  color: #d94a43;
  font-size: 20px;
  background: #fff0ef;
  border-radius: 8px;
  overflow: hidden;

  img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
}

.prize-cover--previewable {
  cursor: zoom-in;
}

</style>
