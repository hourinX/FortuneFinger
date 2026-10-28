
import { computed, onMounted, onUnmounted, reactive, ref, toRef, type Reactive } from 'vue'

export interface CheckinSource {
  qrCode: string
  status: 'pending' | 'active' | 'ended'
  checkedIn: number
  total: number | null
  deadline: number
  serverTimeOffset: number
}

export interface CheckinContent {
  company: string
  title: string
  description: string
}

export function useCheckin(source: Reactive<CheckinSource>) {
  const arrival = ref(0)
  const refreshPeriod = 30_000
  const now = ref(Date.now() + source.serverTimeOffset)
  const nextRefresh = ref(now.value + refreshPeriod)
  const refreshing = ref(false)
  const refreshVersion = ref(0)
  let clock: ReturnType<typeof setInterval> | undefined
  let feedback: ReturnType<typeof setTimeout> | undefined

  const remaining = computed(() => Math.max(0, Math.ceil((source.deadline - now.value) / 1000)))
  const closed = computed(() => source.status === 'ended' || remaining.value === 0)
  const countdown = computed(() => [
    Math.floor(remaining.value / 3600),
    Math.floor((remaining.value % 3600) / 60),
    remaining.value % 60,
  ].map(value => String(value).padStart(2, '0')))
  const refreshSeconds = computed(() => Math.max(0, Math.ceil((nextRefresh.value - now.value) / 1000)))
  const refreshProgress = computed(() => Math.max(0, Math.min(100, ((nextRefresh.value - now.value) / refreshPeriod) * 100)))
  const attendanceProgress = computed(() => (
    source.total !== null && source.total > 0
      ? Math.min(100, Math.round((source.checkedIn / source.total) * 100))
      : 0
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
    get active() {
      return source.status === 'active'
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
    get hasTotal() {
      return source.total !== null
    },
    progress: attendanceProgress,
    arrival,
  })
  const deadline = reactive({
    get value() {
      return source.deadline
    },
    get status() {
      return source.status
    },
    closed,
    countdown,
    label: deadlineLabel,
  })

  function simulateArrival() {
    if (source.status !== 'active' || (source.total !== null && source.checkedIn >= source.total)) return
    source.checkedIn += 1
    arrival.value += 1
  }

  function setAttendance(checkedIn: number, total: number | null) {
    if (checkedIn > source.checkedIn) {
      arrival.value += checkedIn - source.checkedIn
    }
    source.checkedIn = checkedIn
    source.total = total
  }

  function setQRCode(value: string, expiresAt?: string) {
    source.qrCode = value
    nextRefresh.value = expiresAt ? Date.parse(expiresAt) : now.value + refreshPeriod
  }

  function tick() {
    now.value = Date.now() + source.serverTimeOffset
    if (source.status !== 'active') {
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

  return { qr, attendance, deadline, setAttendance, setQRCode, simulateArrival }
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
    qrActive: toRef(props.qr, 'active'),
    attendance: toRef(props.attendance, 'progress'),
    hasTotal: toRef(props.attendance, 'hasTotal'),
    arrival: toRef(props.attendance, 'arrival'),
    closed: toRef(props.deadline, 'closed'),
    countdown: toRef(props.deadline, 'countdown'),
    deadlineLabel: toRef(props.deadline, 'label'),
    status: toRef(props.deadline, 'status'),
  }
}

export type CheckinTemplateState = ReturnType<typeof createCheckinTemplateState>
