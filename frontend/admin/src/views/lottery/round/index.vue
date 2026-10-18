<template>
  <section class="round-page">
    <a-form class="round-toolbar" layout="inline" :model="query" @finish="search">
      <a-form-item>
        <a-input v-model:value="query.keyword" allow-clear placeholder="轮次或奖项名称">
          <template #prefix><search-outlined /></template>
        </a-input>
      </a-form-item>
      <a-form-item>
        <a-select
          v-model:value="query.activity"
          allow-clear
          placeholder="所属活动"
          :options="activityOptions"
        />
      </a-form-item>
      <a-form-item>
        <a-input v-model:value="query.date" type="date" />
      </a-form-item>
      <a-form-item>
        <a-select
          v-model:value="query.status"
          allow-clear
          placeholder="轮次状态"
          :options="statusOptions"
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

    <div class="schedule-summary">
      <div>
        <calendar-outlined />
        <span>共安排 <strong>{{ filteredRows.length }}</strong> 个抽奖轮次</span>
      </div>
      <span>所有轮次均使用完整日期和时间，不跨日复用时间</span>
    </div>

    <a-empty v-if="!groupedRows.length" description="暂无符合条件的抽奖轮次" />
    <section v-for="group in groupedRows" :key="group.date" class="schedule-group">
      <header class="schedule-group__header">
        <div class="schedule-date">
          <strong>{{ dateParts(group.date).day }}</strong>
          <span>{{ dateParts(group.date).yearMonth }}</span>
        </div>
        <div>
          <strong>{{ dateParts(group.date).weekday }}</strong>
          <span>{{ group.rows.length }} 个轮次</span>
        </div>
      </header>

      <div class="round-list">
        <article v-for="item in group.rows" :key="item.id" class="round-card">
          <div class="round-time">
            <strong>{{ item.time.slice(0, 5) }}</strong>
            <span>{{ item.time.slice(5) || ':00' }}</span>
          </div>
          <div class="round-sequence">
            <span>{{ String(item.sequence).padStart(2, '0') }}</span>
          </div>
          <div class="round-main">
            <div class="round-title">
              <strong>{{ item.name }}</strong>
              <a-tag :color="statusMeta[item.status].color">
                {{ statusMeta[item.status].label }}
              </a-tag>
              <a-tag v-if="item.autoStart" color="blue">自动开始</a-tag>
            </div>
            <span>{{ item.activity }}</span>
          </div>
          <div class="round-detail">
            <span>奖项</span>
            <strong>{{ item.prize }}</strong>
          </div>
          <div class="round-detail round-detail--count">
            <span>本轮名额</span>
            <strong>{{ item.winnerCount }} 人</strong>
          </div>
          <div class="round-detail round-detail--template">
            <span>大屏效果</span>
            <strong>{{ item.template }}</strong>
          </div>
          <a-space class="round-actions" :size="14">
            <a @click="openEdit(item)"><edit-outlined />编辑</a>
            <a class="danger-link" @click="removeRound(item)"><delete-outlined />删除</a>
          </a-space>
        </article>
      </div>
    </section>

    <a-modal
      v-model:open="formVisible"
      :title="editingId ? '编辑抽奖轮次' : '新增抽奖轮次'"
      width="680px"
      ok-text="确定"
      cancel-text="取消"
      @ok="saveRound"
    >
      <a-form class="round-form" layout="vertical" :model="form">
        <div class="round-form__grid">
          <a-form-item label="所属活动" required>
            <a-select v-model:value="form.activity" :options="activityOptions" />
          </a-form-item>
          <a-form-item label="轮次名称" required>
            <a-input v-model:value="form.name" placeholder="例如：第一轮抽奖" />
          </a-form-item>
          <a-form-item label="开奖日期" required>
            <a-input v-model:value="form.date" type="date" />
          </a-form-item>
          <a-form-item label="开奖时间" required>
            <a-input v-model:value="form.time" type="time" step="1" />
          </a-form-item>
          <a-form-item label="奖项" required>
            <a-input v-model:value="form.prize" placeholder="请选择或输入奖项" />
          </a-form-item>
          <a-form-item label="本轮中奖人数" required>
            <a-input-number v-model:value="form.winnerCount" :min="1" :max="999" />
          </a-form-item>
          <a-form-item label="大屏效果" required>
            <a-select v-model:value="form.template" :options="templateOptions" />
          </a-form-item>
          <a-form-item label="到点自动开始">
            <a-switch v-model:checked="form.autoStart" />
            <span class="switch-description">开启后按照设定日期和时间自动进入本轮抽奖</span>
          </a-form-item>
        </div>
      </a-form>
    </a-modal>
  </section>
</template>

<script setup>
import { computed, reactive, ref } from 'vue'
import {
  CalendarOutlined,
  DeleteOutlined,
  EditOutlined,
  PlusOutlined,
  SearchOutlined,
} from '@ant-design/icons-vue'
import { Modal, message } from 'ant-design-vue'

const activityOptions = [
  { label: '2026 FortuneFinger 年会', value: '2026 FortuneFinger 年会' },
  { label: '2027 新年庆典', value: '2027 新年庆典' },
]
const statusOptions = [
  { label: '待开始', value: 'pending' },
  { label: '进行中', value: 'drawing' },
  { label: '已结束', value: 'finished' },
]
const templateOptions = [
  { label: '透明球机', value: '透明球机' },
  { label: '金色礼签', value: '金色礼签' },
]
const statusMeta = {
  pending: { label: '待开始', color: 'default' },
  drawing: { label: '进行中', color: 'processing' },
  finished: { label: '已结束', color: 'success' },
}
const initialRows = [
  {
    id: '1',
    sequence: 1,
    activity: '2026 FortuneFinger 年会',
    name: '第一轮抽奖',
    date: '2026-12-31',
    time: '18:30:00',
    prize: '三等奖 · 智能手环',
    winnerCount: 20,
    template: '金色礼签',
    autoStart: true,
    status: 'pending',
  },
  {
    id: '2',
    sequence: 2,
    activity: '2026 FortuneFinger 年会',
    name: '第二轮抽奖',
    date: '2026-12-31',
    time: '19:20:00',
    prize: '二等奖 · 平板电脑',
    winnerCount: 6,
    template: '金色礼签',
    autoStart: true,
    status: 'pending',
  },
  {
    id: '3',
    sequence: 3,
    activity: '2026 FortuneFinger 年会',
    name: '跨年特别轮',
    date: '2027-01-01',
    time: '00:05:00',
    prize: '一等奖 · 智能科技礼盒',
    winnerCount: 3,
    template: '透明球机',
    autoStart: true,
    status: 'pending',
  },
]

const query = reactive({ keyword: '', activity: undefined, date: '', status: undefined })
const activeQuery = reactive({ ...query })
const rows = ref(initialRows)
const formVisible = ref(false)
const editingId = ref()
const form = reactive(initialForm())
const filteredRows = computed(() =>
  rows.value.filter((item) => {
    const keyword = activeQuery.keyword.toLocaleLowerCase()
    return (
      (!keyword || `${item.name}${item.prize}`.toLocaleLowerCase().includes(keyword)) &&
      (!activeQuery.activity || item.activity === activeQuery.activity) &&
      (!activeQuery.date || item.date === activeQuery.date) &&
      (!activeQuery.status || item.status === activeQuery.status)
    )
  }),
)
const groupedRows = computed(() => {
  const groups = new Map()
  filteredRows.value
    .slice()
    .sort((a, b) => `${a.date} ${a.time}`.localeCompare(`${b.date} ${b.time}`))
    .forEach((item) => {
      if (!groups.has(item.date)) groups.set(item.date, [])
      groups.get(item.date).push(item)
    })
  return [...groups].map(([date, groupRows]) => ({ date, rows: groupRows }))
})

function initialForm() {
  return {
    activity: activityOptions[0].value,
    name: '',
    date: '',
    time: '18:30:00',
    prize: '',
    winnerCount: 1,
    template: templateOptions[0].value,
    autoStart: true,
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

function search() {
  Object.assign(activeQuery, query)
}

function reset() {
  Object.assign(query, { keyword: '', activity: undefined, date: '', status: undefined })
  search()
}

function openCreate() {
  editingId.value = undefined
  Object.assign(form, initialForm())
  formVisible.value = true
}

function openEdit(item) {
  editingId.value = item.id
  Object.assign(form, item)
  formVisible.value = true
}

function saveRound() {
  if (!form.name || !form.date || !form.time || !form.prize) {
    message.warning('请填写完整的轮次信息')
    return
  }
  if (editingId.value) {
    rows.value = rows.value.map((item) =>
      item.id === editingId.value ? { ...item, ...form } : item,
    )
  } else {
    rows.value.push({
      ...form,
      id: String(Date.now()),
      sequence: rows.value.length + 1,
      status: 'pending',
    })
  }
  formVisible.value = false
  message.success(editingId.value ? '轮次修改成功' : '轮次新增成功')
}

function removeRound(item) {
  Modal.confirm({
    title: '确认删除抽奖轮次？',
    content: `将删除“${item.name}（${item.date} ${item.time}）”。`,
    okText: '删除',
    okType: 'danger',
    cancelText: '取消',
    onOk() {
      rows.value = rows.value.filter((row) => row.id !== item.id)
      message.success('轮次删除成功')
    },
  })
}
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

.schedule-summary {
  display: flex;
  min-height: 42px;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 18px;
  padding: 0 14px;
  color: #777;
  font-size: 13px;
  background: #fafafa;
  border: 1px solid #eeeeee;
  border-radius: 6px;

  > div {
    display: flex;
    align-items: center;
    gap: 9px;
    color: #4a4a4a;
  }

  strong {
    color: var(--color-primary);
  }
}

.schedule-group {
  display: grid;
  grid-template-columns: 112px minmax(0, 1fr);
  border-top: 1px solid #ededed;
}

.schedule-group:last-of-type {
  border-bottom: 1px solid #ededed;
}

.schedule-group__header {
  padding: 20px 18px 20px 6px;
  border-right: 1px solid #ededed;

  > div:last-child {
    display: flex;
    flex-direction: column;
    gap: 4px;
    margin-top: 12px;
    color: #4f4f4f;

    span {
      color: #999;
      font-size: 12px;
    }
  }
}

.schedule-date {
  display: flex;
  align-items: flex-end;
  gap: 7px;

  strong {
    color: #292929;
    font-size: 30px;
    font-weight: 600;
    line-height: 1;
  }

  span {
    color: #999;
    font-size: 11px;
  }
}

.round-list {
  min-width: 0;
}

.round-card {
  display: grid;
  min-height: 88px;
  grid-template-columns: 82px 38px minmax(180px, 1.2fr) minmax(170px, 1fr) 100px 130px 130px;
  align-items: center;
  padding: 12px 16px 12px 20px;
  border-bottom: 1px solid #f0f0f0;
  transition: background 160ms ease;

  &:last-child {
    border-bottom: 0;
  }

  &:hover {
    background: #fafafa;
  }
}

.round-time {
  display: flex;
  align-items: baseline;
  color: #2f2f2f;

  strong {
    font-size: 19px;
    font-weight: 600;
  }

  span {
    color: #aaa;
    font-size: 11px;
  }
}

.round-sequence {
  position: relative;
  display: flex;
  width: 24px;
  height: 24px;
  align-items: center;
  justify-content: center;
  color: #888;
  font-size: 11px;
  background: #f3f3f3;
  border-radius: 50%;
}

.round-main,
.round-detail {
  display: flex;
  min-width: 0;
  flex-direction: column;
  gap: 5px;

  > span {
    overflow: hidden;
    color: #999;
    font-size: 12px;
    white-space: nowrap;
    text-overflow: ellipsis;
  }

  > strong {
    overflow: hidden;
    color: #444;
    font-size: 13px;
    font-weight: 500;
    white-space: nowrap;
    text-overflow: ellipsis;
  }
}

.round-title {
  display: flex;
  min-width: 0;
  align-items: center;
  gap: 7px;

  strong {
    overflow: hidden;
    color: #2f2f2f;
    font-size: 15px;
    white-space: nowrap;
    text-overflow: ellipsis;
  }

  :deep(.ant-tag) {
    margin: 0;
  }
}

.round-detail--count strong {
  color: var(--color-primary);
}

.round-actions {
  justify-self: end;

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

.round-form {
  padding-top: 12px;
}

.round-form__grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 0 20px;

  :deep(.ant-input-number),
  :deep(.ant-select) {
    width: 100%;
  }
}

.switch-description {
  margin-left: 10px;
  color: #999;
  font-size: 12px;
}

@media (max-width: 1280px) {
  .round-card {
    grid-template-columns: 75px 34px minmax(170px, 1fr) minmax(150px, 1fr) 90px 120px;
  }

  .round-detail--template {
    display: none;
  }
}
</style>
