<template>
  <section class="activity-page">
    <a-form class="activity-toolbar" layout="inline" :model="query" @finish="search">
      <a-form-item>
        <a-input v-model:value="query.name" allow-clear placeholder="活动名称">
          <template #prefix><search-outlined /></template>
        </a-input>
      </a-form-item>
      <a-form-item>
        <a-select
          v-model:value="query.orgId"
          allow-clear
          placeholder="归属企业"
          :options="companyOptions"
        />
      </a-form-item>
      <a-form-item><a-date-picker v-model:value="query.date" placeholder="活动日期" /></a-form-item>
      <a-form-item>
        <a-select
          v-model:value="query.phase"
          allow-clear
          placeholder="活动状态"
          :options="phaseOptions"
        />
      </a-form-item>
      <a-form-item>
        <a-space v-permission="'lottery:activities:list'" :size="8">
          <a-button type="primary" html-type="submit">查询</a-button>
          <a-button @click="reset">重置</a-button>
        </a-space>
      </a-form-item>
      <a-form-item class="activity-toolbar__create">
        <a-button v-permission="'lottery:activities:add'" type="primary" @click="openCreate">
          <plus-outlined />新增活动
        </a-button>
      </a-form-item>
    </a-form>

    <a-table
      class="activity-table"
      :columns="columns"
      :data-source="rows"
      :loading="loading"
      :pagination="pagination"
      row-key="id"
      :scroll="{ x: 980 }"
      @change="changeTable"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'activity'">
          <div class="activity-identity">
            <div class="activity-identity__date">
              <strong>{{ dayOf(record.activity_start_time) }}</strong>
              <span>{{ monthOf(record.activity_start_time) }}</span>
            </div>
            <div>
              <strong>{{ record.name }}</strong
              ><span>{{ record.venue }}</span>
            </div>
          </div>
        </template>
        <template v-else-if="column.key === 'time'">
          <div class="time-range">
            <strong>
              {{ formatDate(record.activity_start_time) }} —
              {{ formatDate(record.activity_end_time) }}
            </strong>
            <span>
              {{ formatClock(record.activity_start_time) }} —
              {{ formatClock(record.activity_end_time) }}
            </span>
          </div>
        </template>
        <template v-else-if="column.key === 'checkin'">
          <div class="time-range">
            <strong>
              {{ formatDate(record.checkin_start_time) }} —
              {{ formatDate(record.checkin_end_time) }}
            </strong>
            <span>
              {{ formatClock(record.checkin_start_time) }} —
              {{ formatClock(record.checkin_end_time) }}
            </span>
          </div>
        </template>
        <template v-else-if="column.key === 'phase'">
          <a-tag v-if="record.status === 2" color="error">已禁用</a-tag>
          <a-tag v-else :style="activityStatusTagStyle(record.phase)">
            {{ phaseOptionMap.get(record.phase)?.label || record.phase }}
          </a-tag>
        </template>
        <template v-else-if="column.key === 'action'">
          <a-space :size="12">
            <a v-permission="'lottery:activities:template'" @click="openTemplate(record)">
              签到模板
            </a>
            <a :href="getCheckinScreenUrl(record.code)" target="_blank"> 签到大屏 </a>
            <a v-permission="'lottery:activities:edit'" @click="openEdit(record)">
              <edit-outlined />编辑
            </a>
            <a
              v-permission="'lottery:activities:del'"
              class="danger-link"
              @click="removeActivity(record)"
            >
              <delete-outlined />删除
            </a>
          </a-space>
        </template>
      </template>
      <template #emptyText><a-empty description="暂无活动数据" /></template>
    </a-table>

    <ActivityForm ref="formRef" :companies="companyOptions" @success="loadActivities" />
    <TemplateForm ref="templateFormRef" @confirm="confirmTemplate" />
  </section>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import dayjs from 'dayjs'
import { DeleteOutlined, EditOutlined, PlusOutlined, SearchOutlined } from '@ant-design/icons-vue'
import { Modal, message } from 'ant-design-vue'
import { deleteActivity, getActivityList, updateActivityCheckinTemplate } from '@/api/lottery'
import { getOrgOptions } from '@/api/org'
import { useDictionaryOptions } from '@/hooks/useDictionary'
import { dictionaryTagStyleByCode } from '@/utils/dictionary-tag'
import ActivityForm from './components/form.vue'
import TemplateForm from './components/template-form.vue'

const phaseOptions = useDictionaryOptions('activities_status')
const phaseOptionMap = computed(
  () => new Map(phaseOptions.value.map((option) => [option.value, option])),
)
function activityStatusTagStyle(phase) {
  return dictionaryTagStyleByCode('activities_status', phaseOptionMap.value.get(phase))
}
const columns = [
  { title: '活动信息', key: 'activity', width: 280 },
  { title: '归属企业', dataIndex: 'org_name', key: 'orgName', width: 160 },
  { title: '活动时间', key: 'time', width: 200 },
  { title: '签到时间', key: 'checkin', width: 200 },
  { title: '状态', key: 'phase', width: 90 },
  { title: '操作', key: 'action', width: 310, fixed: 'right' },
]
const rows = ref([])
const loading = ref(false)
const companyOptions = ref([])
const query = reactive({ name: '', orgId: undefined, date: undefined, phase: undefined })
const page = reactive({ current: 1, pageSize: 10, total: 0 })
const formRef = ref()
const templateFormRef = ref()
const pagination = computed(() => ({
  current: page.current,
  pageSize: page.pageSize,
  total: page.total,
  showSizeChanger: true,
  showTotal: (total) => `共 ${total} 条`,
}))

function search() {
  page.current = 1
  loadActivities()
}

function reset() {
  Object.assign(query, { name: '', orgId: undefined, date: undefined, phase: undefined })
  search()
}

function openCreate() {
  formRef.value?.openCreate()
}

function openEdit(item) {
  formRef.value?.openEdit(item.id)
}

function openTemplate(item) {
  templateFormRef.value?.open(item, item.checkin_template)
}

function getCheckinScreenUrl(code) {
  const configuredUrl = import.meta.env.VITE_WEB_PREVIEW_URL
  const baseUrl =
    configuredUrl ||
    (import.meta.env.DEV ? `${location.protocol}//${location.hostname}:5174/checkin` : '/checkin')
  const url = new URL(baseUrl, location.origin)
  url.searchParams.set('template', code)
  return url.toString()
}

async function confirmTemplate({ activityId, template }) {
  try {
    await updateActivityCheckinTemplate(activityId, { checkin_template: template })
    message.success('签到模板保存成功')
    await loadActivities()
  } catch (error) {
    message.error(error.message)
  }
}

function removeActivity(item) {
  Modal.confirm({
    title: '确认删除活动？',
    content: `将删除“${item.name}”。已产生签到或抽奖数据的活动后续将不允许删除。`,
    okText: '删除',
    okType: 'danger',
    cancelText: '取消',
    async onOk() {
      try {
        await deleteActivity(item.id)
        message.success('活动删除成功')
        if (rows.value.length === 1 && page.current > 1) page.current -= 1
        await loadActivities()
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
  loadActivities()
}

async function loadActivities() {
  loading.value = true
  try {
    const resp = await getActivityList({
      page: page.current,
      pageSize: page.pageSize,
      name: query.name || undefined,
      org_id: query.orgId,
      date: query.date?.format('YYYY-MM-DD'),
      phase: query.phase,
    })
    rows.value = resp.list || []
    page.total = resp.total || 0
  } catch (error) {
    message.error(error.message)
  } finally {
    loading.value = false
  }
}

async function loadCompanyOptions() {
  try {
    const resp = await getOrgOptions()
    companyOptions.value = (resp.list || []).flatMap((group) => [
      { label: group.name, value: group.id },
      ...(group.children || []).map((company) => ({ label: company.name, value: company.id })),
    ])
  } catch (error) {
    message.error(error.message)
  }
}

function dayOf(value) {
  return dayjs(value).format('DD')
}
function monthOf(value) {
  return dayjs(value).format('MM月')
}
function formatDate(value) {
  return dayjs(value).format('YYYY-MM-DD')
}
function formatClock(value) {
  return dayjs(value).format('HH:mm')
}

onMounted(() => {
  loadActivities()
  loadCompanyOptions()
})
</script>

<style scoped lang="less">
.activity-page {
  min-width: 0;
}
.activity-toolbar {
  display: flex;
  min-height: 58px;
  align-items: center;
  padding: 0 0 18px;
}
.activity-toolbar :deep(.ant-form-item) {
  margin-bottom: 0;
}
.activity-toolbar :deep(.ant-input) {
  width: 190px;
}
.activity-toolbar :deep(.ant-select) {
  width: 150px;
}
.activity-toolbar__create {
  margin-right: 0 !important;
  margin-left: auto !important;
}
.activity-table :deep(.ant-table-thead > tr > th) {
  height: 48px;
  padding: 0 16px;
  color: #494949;
  font-size: 13px;
  font-weight: 600;
  background: #fafafa;
}
.activity-table :deep(.ant-table-tbody > tr > td) {
  height: 76px;
  padding: 10px 16px;
}
.activity-table :deep(.ant-tag) {
  margin: 0;
}
.activity-table a {
  color: #1677ff;
}
.activity-table .danger-link {
  color: #cf3f36;
}
.activity-table :deep(.anticon) {
  margin-right: 4px;
}
.activity-identity {
  display: flex;
  min-width: 0;
  align-items: center;
  gap: 12px;
}
.activity-identity > div:last-child {
  display: flex;
  min-width: 0;
  flex-direction: column;
  gap: 5px;
}
.activity-identity > div:last-child strong {
  overflow: hidden;
  color: #303030;
  font-size: 14px;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.activity-identity > div:last-child span {
  overflow: hidden;
  color: #999;
  font-size: 12px;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.activity-identity__date {
  display: flex;
  width: 44px;
  height: 46px;
  flex: 0 0 44px;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background: #fff1ef;
  border-radius: 7px;
}
.activity-identity__date strong {
  color: #db433c;
  font-size: 17px;
  line-height: 19px;
}
.activity-identity__date span {
  color: #b56c67;
  font-size: 9px;
}
.time-range {
  display: flex;
  flex-direction: column;
  gap: 5px;
}
.time-range strong {
  color: #464646;
  font-size: 12px;
  font-weight: 500;
}
.time-range span {
  color: #999;
  font-size: 11px;
}
</style>
