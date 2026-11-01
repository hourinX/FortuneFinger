<template>
  <section class="round-page">
    <a-form class="round-toolbar" layout="inline" :model="query" @finish="search">
      <a-form-item>
        <a-select
          v-model:value="query.activityId"
          allow-clear
          placeholder="所属活动"
          :options="activityOptions"
        />
      </a-form-item>
      <a-form-item>
        <a-select
          v-model:value="query.status"
          allow-clear
          placeholder="轮次状态"
          :options="roundStatusOptions"
        />
      </a-form-item>
      <a-form-item>
        <a-space :size="8">
          <a-button type="primary" html-type="submit">查询</a-button>
          <a-button @click="reset">重置</a-button>
        </a-space>
      </a-form-item>
      <a-form-item class="round-toolbar__create">
        <a-button type="primary" @click="openCreate"><plus-outlined />新增轮次</a-button>
      </a-form-item>
    </a-form>

    <a-table
      class="round-table"
      :columns="columns"
      :data-source="tableRows"
      :pagination="pagination"
      :loading="tableLoading"
      row-key="id"
      :scroll="{ x: 1305 }"
      :custom-row="customRow"
      @change="changeTable"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'activity'">
          <div class="activity-identity">
            <div class="activity-identity__content">
              <strong>{{ record.activityName }}</strong>
              <span class="activity-identity__date">
                <calendar-outlined />{{ record.activityDateRange }}
              </span>
              <small>{{ record.activityRoundCount }} 个抽奖轮次</small>
            </div>
          </div>
        </template>
        <template v-else-if="column.key === 'drawDate'">
          <div class="draw-time">
            <strong>{{ dateParts(record.date).day }}</strong>
            <div>
              <span>{{ dateParts(record.date).yearMonth }}</span>
              <small>{{ dateParts(record.date).weekday }}</small>
            </div>
          </div>
        </template>
        <template v-else-if="column.key === 'clock'">
          <strong class="clock-cell">{{ record.time.slice(0, 5) }}</strong>
        </template>
        <template v-else-if="column.key === 'round'">
          <div class="round-identity">
            <span>{{ String(record.sequence).padStart(2, '0') }}</span>
            <strong>{{ record.name }}</strong>
          </div>
        </template>
        <template v-else-if="column.key === 'prize'">
          <div class="prize-cell">
            <a-tag :style="prizeLevelTagStyle(record.level)">
              {{ prizeLevelOptionMap.get(record.level)?.label || record.level }}
            </a-tag>
            <strong>{{ record.prizeName }}</strong>
          </div>
        </template>
        <template v-else-if="column.key === 'winnerCount'">
          <strong class="winner-count">{{ record.winnerCount }} 人</strong>
        </template>
        <template v-else-if="column.key === 'status'">
          <a-tag :style="roundStatusTagStyle(record.status)">
            {{ roundStatusOptionMap.get(record.status)?.label || record.status }}
          </a-tag>
        </template>
        <template v-else-if="column.key === 'action'">
          <a-space :size="12">
            <a @click="openEdit(record)"><edit-outlined />编辑</a>
            <a class="danger-link" @click="removeRound(record)"> <delete-outlined />删除 </a>
          </a-space>
        </template>
      </template>
      <template #emptyText><a-empty description="暂无抽奖轮次" /></template>
    </a-table>

    <RoundForm ref="formRef" :activities="activityOptions" @success="loadRounds" />
  </section>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { CalendarOutlined, DeleteOutlined, EditOutlined, PlusOutlined } from '@ant-design/icons-vue'
import { Modal, message } from 'ant-design-vue'
import { useDictionaryOptions } from '@/hooks/useDictionary'
import { dictionaryTagStyleByCode } from '@/utils/dictionary-tag'
import { deleteRound, getActivityList, getRoundList } from '@/api/lottery'
import RoundForm from './components/form.vue'

const activityOptions = ref([])
const roundStatusOptions = useDictionaryOptions('activities_status')
const roundStatusOptionMap = computed(
  () => new Map(roundStatusOptions.value.map((option) => [option.value, option])),
)
function roundStatusTagStyle(status) {
  return dictionaryTagStyleByCode('activities_status', roundStatusOptionMap.value.get(status))
}
const prizeLevelOptions = useDictionaryOptions('prize_level')
const prizeLevelOptionMap = computed(
  () => new Map(prizeLevelOptions.value.map((option) => [option.value, option])),
)
function prizeLevelTagStyle(level) {
  return dictionaryTagStyleByCode('prize_level', prizeLevelOptionMap.value.get(level))
}
const columns = [
  {
    title: '活动信息',
    key: 'activity',
    width: 280,
    customCell: (record) => ({
      rowSpan: record.activityRowSpan,
      class: record.activityPageEnd ? 'activity-column activity-page-end' : 'activity-column',
    }),
  },
  {
    title: '开奖日期',
    key: 'drawDate',
    width: 135,
    customCell: (record) => ({ rowSpan: record.dateRowSpan, class: 'draw-date-column' }),
  },
  { title: '开奖时间', key: 'clock', width: 95 },
  { title: '轮次信息', key: 'round', width: 220 },
  { title: '奖项与奖品', key: 'prize', width: 260 },
  { title: '本轮名额', key: 'winnerCount', width: 105 },
  { title: '状态', key: 'status', width: 90 },
  { title: '操作', key: 'action', width: 150, fixed: 'right' },
]
const query = reactive({ activityId: undefined, status: undefined })
const rows = ref([])
const page = reactive({ current: 1, pageSize: 10, total: 0 })
const tableLoading = ref(false)
const formRef = ref()
const pagination = computed(() => ({
  current: page.current,
  pageSize: page.pageSize,
  total: page.total,
  showSizeChanger: true,
  showTotal: (total) => `共 ${total} 条`,
}))
const tableRows = computed(() => {
  const activities = new Map()
  rows.value.forEach((item) => {
    if (!activities.has(item.activityId)) activities.set(item.activityId, [])
    activities.get(item.activityId).push(item)
  })
  const result = [...activities.values()].flatMap((activityRows) =>
    activityRows.slice().sort((a, b) => `${a.date} ${a.time}`.localeCompare(`${b.date} ${b.time}`)),
  )
  const activityCounts = new Map()
  rows.value.forEach((item) => {
    activityCounts.set(item.activityId, (activityCounts.get(item.activityId) || 0) + 1)
  })
  const mappedRows = result.map((item, index) => {
    const activity = activityOptions.value.find((option) => option.value === item.activityId)
    return {
      ...item,
      activityDateRange:
        item.activityDateRange ||
        (activity ? `${activity.startDate} — ${activity.endDate}` : item.date),
      activityRoundCount: item.activityRoundCount ?? activityCounts.get(item.activityId),
      activityRowSpan: 1,
      dateRowSpan: 1,
      activityPageEnd: false,
      activityGroupStart: index > 0 && result[index - 1].activityId !== item.activityId,
    }
  })
  for (let start = 0; start < mappedRows.length; start += page.pageSize) {
    const pageRows = mappedRows.slice(start, start + page.pageSize)
    setRowSpan(pageRows, 'activityRowSpan', (left, right) => left.activityId === right.activityId)
    setRowSpan(
      pageRows,
      'dateRowSpan',
      (left, right) => left.activityId === right.activityId && left.date === right.date,
    )
    let lastActivityStart = pageRows.length - 1
    while (lastActivityStart >= 0 && pageRows[lastActivityStart].activityRowSpan === 0) {
      lastActivityStart -= 1
    }
    if (lastActivityStart >= 0) pageRows[lastActivityStart].activityPageEnd = true
  }
  return mappedRows
})

function setRowSpan(list, field, matches) {
  for (let start = 0; start < list.length;) {
    let end = start + 1
    while (end < list.length && matches(list[start], list[end])) end += 1
    list[start][field] = end - start
    for (let index = start + 1; index < end; index += 1) list[index][field] = 0
    start = end
  }
}

function dateParts(date) {
  const value = new Date(`${date}T00:00:00`)
  return {
    day: date.slice(8),
    yearMonth: date.slice(0, 7).replace('-', ' / '),
    weekday: new Intl.DateTimeFormat('zh-CN', { weekday: 'long' }).format(value),
  }
}

async function search() {
  page.current = 1
  await loadRounds()
}

async function reset() {
  Object.assign(query, { activityId: undefined, status: undefined })
  page.current = 1
  await loadRounds()
}

async function changeTable(paginationState) {
  page.current = paginationState.current
  page.pageSize = paginationState.pageSize
  await loadRounds()
}

function customRow(record) {
  return { class: record.activityGroupStart ? 'activity-group-start' : '' }
}

function openCreate() {
  formRef.value?.openCreate()
}

function openEdit(item) {
  formRef.value?.openEdit(item.id)
}

function removeRound(item) {
  Modal.confirm({
    title: '确认删除抽奖轮次？',
    content: `将删除“${item.name}（${item.date} ${item.time}）”。`,
    okText: '删除',
    okType: 'danger',
    cancelText: '取消',
    async onOk() {
      await deleteRound(item.id)
      message.success('轮次删除成功')
      await loadRounds()
    },
  })
}

async function loadRounds() {
  tableLoading.value = true
  try {
    const resp = await getRoundList({
      page: page.current,
      pageSize: page.pageSize,
      activity_id: query.activityId,
      phase: query.status,
    })
    rows.value = flattenRoundActivities(resp?.list || [])
    page.total = resp?.total || 0
  } catch (error) {
    message.error(error.message)
  } finally {
    tableLoading.value = false
  }
}

async function loadActivityOptions() {
  try {
    const resp = await getActivityList({ page: 1, pageSize: 100, status: 1 })
    activityOptions.value = (resp.list || []).map((item) => ({
      label: item.name,
      value: item.id,
      orgId: item.org_id,
      startDate: item.activity_start_time.slice(0, 10),
      endDate: item.activity_end_time.slice(0, 10),
    }))
  } catch (error) {
    message.error(error.message)
  }
}

function flattenRoundActivities(activities) {
  return activities.flatMap((activity) =>
    (activity.dates || []).flatMap((date) =>
      (date.rounds || []).map((round) => ({
        id: round.id,
        sequence: round.sequence,
        activityId: activity.activity_id,
        activityName: activity.activity_name,
        activityDateRange: `${activity.activity_start_time.slice(0, 10)} — ${activity.activity_end_time.slice(0, 10)}`,
        activityRoundCount: activity.round_count,
        name: round.name,
        date: date.draw_date,
        time: round.draw_time,
        level: round.level,
        prizeId: round.prize_id,
        prizeName: round.prize_name,
        winnerCount: round.winner_count,
        status: round.phase,
      })),
    ),
  )
}

onMounted(() => {
  loadActivityOptions()
  loadRounds()
})
</script>

<style scoped lang="less">
.round-page {
  min-width: 0;
}

.round-toolbar {
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
    width: 170px;
  }
}

.round-toolbar__create {
  margin-right: 0 !important;
  margin-left: auto !important;
}

.round-table {
  :deep(.ant-table-thead > tr > th) {
    height: 48px;
    padding: 0 16px;
    color: #494949;
    font-size: 13px;
    font-weight: 600;
    background: #fafafa;
  }

  :deep(.ant-table-tbody > tr > td) {
    height: 76px;
    padding: 10px 16px;
  }

  :deep(.ant-table-tbody .activity-column) {
    border-right: 1px solid #ededed !important;
  }

  :deep(.ant-table-tbody .draw-date-column) {
    border-right: 1px solid #ededed !important;
  }

  :deep(.ant-table-tbody .activity-page-end) {
    border-bottom: 1px solid #f0f0f0 !important;
  }

  :deep(.ant-table-tbody > tr.activity-group-start > td) {
    border-top: 2px solid #e6e6e6 !important;
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

.activity-identity {
  position: relative;
  display: flex;
  min-width: 0;
  align-items: stretch;
  padding-left: 14px;

  &::before {
    position: absolute;
    top: 2px;
    bottom: 2px;
    left: 0;
    width: 3px;
    background: linear-gradient(180deg, #e95a52 0%, #cf3f36 100%);
    border-radius: 3px;
    content: '';
  }
}

.activity-identity__content {
  display: flex;
  min-width: 0;
  flex-direction: column;
  gap: 7px;

  strong,
  span {
    overflow: hidden;
    white-space: nowrap;
    text-overflow: ellipsis;
  }

  strong {
    color: #303030;
    font-size: 14px;
    font-weight: 600;
    line-height: 21px;
  }

  span {
    color: #999;
    font-size: 11px;
  }

  small {
    width: fit-content;
    padding: 2px 8px;
    color: #c94a43;
    font-size: 10px;
    line-height: 17px;
    background: #fff1ef;
    border-radius: 10px;
  }
}

.activity-identity__date {
  display: flex;
  align-items: center;
  gap: 1px;
}

.round-identity {
  display: flex;
  min-width: 0;
  align-items: center;
  gap: 10px;

  > span {
    display: inline-flex;
    width: 25px;
    height: 25px;
    flex: 0 0 25px;
    align-items: center;
    justify-content: center;
    color: #888;
    font-size: 10px;
    background: #f3f3f3;
    border-radius: 50%;
  }

  strong {
    overflow: hidden;
    color: #303030;
    font-size: 14px;
    white-space: nowrap;
    text-overflow: ellipsis;
  }
}

.draw-time {
  display: flex;
  align-items: center;
  gap: 9px;

  > strong {
    color: #303030;
    font-size: 25px;
    font-weight: 600;
    line-height: 1;
  }

  > div {
    display: flex;
    flex-direction: column;
    gap: 3px;
    white-space: nowrap;

    span {
      color: #666;
      font-size: 11px;
    }

    small {
      color: #999;
      font-size: 10px;
    }
  }
}

.prize-cell {
  display: flex;
  min-width: 0;
  align-items: center;
  gap: 9px;

  :deep(.ant-tag) {
    flex: 0 0 auto;
  }

  strong {
    overflow: hidden;
    color: #3f3f3f;
    font-size: 13px;
    font-weight: 500;
    white-space: nowrap;
    text-overflow: ellipsis;
  }
}

.clock-cell {
  color: #464646;
  font-size: 13px;
  font-weight: 500;
}

.winner-count {
  color: var(--color-primary);
  font-weight: 500;
}
</style>
