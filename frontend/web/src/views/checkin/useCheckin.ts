
import { computed, onMounted, onUnmounted, reactive, ref, toRef, type Reactive } from 'vue'

export interface CheckinSource {
  qrCode: string
  checkedIn: number
  total: number
  deadline: number
}

export interface CheckinContent {
  company: string
  title: string
  description: string
}

export function useCheckin(source: Reactive<CheckinSource>) {
  const arrival = ref(0)
  const refreshPeriod = 30_000
  const now = ref(Date.now())
  const nextRefresh = ref(now.value + refreshPeriod)
  const refreshing = ref(false)
  const refreshVersion = ref(0)
  let clock: ReturnType<typeof setInterval> | undefined
  let feedback: ReturnType<typeof setTimeout> | undefined

  const remaining = computed(() => Math.max(0, Math.ceil((source.deadline - now.value) / 1000)))
  const closed = computed(() => remaining.value === 0)
  const countdown = computed(() => [
    Math.floor(remaining.value / 3600),
    Math.floor((remaining.value % 3600) / 60),
    remaining.value % 60,
  ].map(value => String(value).padStart(2, '0')))
  const refreshSeconds = computed(() => Math.max(0, Math.ceil((nextRefresh.value - now.value) / 1000)))
  const refreshProgress = computed(() => Math.max(0, Math.min(100, ((nextRefresh.value - now.value) / refreshPeriod) * 100)))
  const attendanceProgress = computed(() => (
    source.total > 0 ? Math.min(100, Math.round((source.checkedIn / source.total) * 100)) : 0
  ))
  const deadlineLabel = computed(() => new Intl.DateTimeFormat('zh-CN', {
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    hour12: false,
  }).format(source.deadline))

  const qr = reactive({
    get value() {
      return source.qrCode
    },
    refreshing,
    refreshVersion,
    refreshSeconds,
    refreshProgress,
  })
  const attendance = reactive({
    get checkedIn() {
      return source.checkedIn
    },
    get total() {
      return source.total
    },
    progress: attendanceProgress,
    arrival,
  })
  const deadline = reactive({
    get value() {
      return source.deadline
    },
    closed,
    countdown,
    label: deadlineLabel,
  })

  function simulateArrival() {
    if (deadline.closed || source.checkedIn >= source.total) return
    source.checkedIn += 1
    arrival.value += 1
  }

  function tick() {
    now.value = Date.now()
    if (deadline.closed) {
      clearInterval(clock)
      return
    }
    if (now.value >= nextRefresh.value) {
      nextRefresh.value = now.value + refreshPeriod
      refreshVersion.value += 1
      refreshing.value = true
      clearTimeout(feedback)
      feedback = setTimeout(() => {
        refreshing.value = false
      }, 900)
    }
  }

  onMounted(() => {
    clock = setInterval(tick, 250)
  })
  onUnmounted(() => {
    clearInterval(clock)
    clearTimeout(feedback)
  })

  return { qr, attendance, deadline, simulateArrival }
}

export type CheckinViewState = ReturnType<typeof useCheckin>

export interface CheckinTemplateProps {
  content: CheckinContent
  qr: CheckinViewState['qr']
  attendance: CheckinViewState['attendance']
  deadline: CheckinViewState['deadline']
}

export function createCheckinTemplateState(props: CheckinTemplateProps) {
  const event = reactive({
    get company() {
      return props.content.company
    },
    get title() {
      return props.content.title
    },
    get description() {
      return props.content.description
    },
    get deadline() {
      return props.deadline.value
    },
    get checkedIn() {
      return props.attendance.checkedIn
    },
    get total() {
      return props.attendance.total
    },
  })

  return {
    event,
    refreshing: toRef(props.qr, 'refreshing'),
    refreshVersion: toRef(props.qr, 'refreshVersion'),
    refreshSeconds: toRef(props.qr, 'refreshSeconds'),
    refreshProgress: toRef(props.qr, 'refreshProgress'),
    attendance: toRef(props.attendance, 'progress'),
    arrival: toRef(props.attendance, 'arrival'),
    closed: toRef(props.deadline, 'closed'),
    countdown: toRef(props.deadline, 'countdown'),
    deadlineLabel: toRef(props.deadline, 'label'),
  }
}

export type CheckinTemplateState = ReturnType<typeof createCheckinTemplateState>
