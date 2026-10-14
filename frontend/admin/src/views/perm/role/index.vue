<template>
  <section class="role-page">
    <a-form class="toolbar" layout="inline" :model="query" @finish="search">
      <a-form-item>
        <a-input
          v-model:value="query.name"
          allow-clear
          placeholder="角色名称"
          @press-enter="search"
        >
          <template #prefix><search-outlined /></template>
        </a-input>
      </a-form-item>
      <a-form-item>
        <a-input
          v-model:value="query.code"
          allow-clear
          placeholder="角色编码"
          @press-enter="search"
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
        <a-button type="primary" @click="createRole"><plus-outlined />新增</a-button>
      </a-form-item>
    </a-form>

    <a-table
      class="role-table"
      :columns="columns"
      :data-source="visibleRows"
      :loading="loading"
      :pagination="false"
      row-key="id"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'name'">
          <div class="role-name" :style="{ paddingLeft: `${record.depth * 22}px` }">
            <team-outlined class="role-icon" />
            <span>{{ record.name }}</span>
            <button
              v-if="record.hasChildren"
              type="button"
              class="expand-button"
              :class="{ 'expand-button--open': !collapsedIds.has(record.id) }"
              @click="toggleRow(record.id)"
            >
              <right-outlined />
            </button>
          </div>
        </template>
        <template v-else-if="column.key === 'code'">
          <span class="code-text">{{ record.code }}</span>
        </template>
        <template v-else-if="column.key === 'status'">
          <StatusPulse :status="record.status" />
        </template>
        <template v-else-if="column.key === 'action'">
          <a-space :size="14">
            <a class="action-link" @click="configurePermission(record)">
              <safety-certificate-outlined />{{
                record.code === 'super_admin' ? '权限预览' : '配置权限'
              }}
            </a>
            <a class="action-link" @click="createChild(record)"><plus-outlined />添加</a>
            <template v-if="canOperate(record)">
              <a class="action-link" @click="editRole(record)"><edit-outlined />编辑</a>
              <a class="danger-link" @click="removeRole(record)"><delete-outlined />删除</a>
            </template>
          </a-space>
        </template>
      </template>
    </a-table>

    <RoleForm ref="formRef" @success="loadRoles" />
    <PermissionDrawer ref="permissionRef" />
  </section>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import {
  DeleteOutlined,
  EditOutlined,
  PlusOutlined,
  RightOutlined,
  SafetyCertificateOutlined,
  SearchOutlined,
  TeamOutlined,
} from '@ant-design/icons-vue'
import { Modal, message } from 'ant-design-vue'
import StatusPulse from '@/components/StatusPulse.vue'
import { useDictionaryOptions } from '@/hooks/useDictionary'
import { deleteRole, getRoleList } from '@/api/role'
import { useAdminStore } from '@/store'
import RoleForm from './components/form.vue'
import PermissionDrawer from './components/permission.vue'

const columns = [
  { title: '角色名称', dataIndex: 'name', key: 'name' },
  { title: '角色编码', dataIndex: 'code', key: 'code' },
  { title: '状态', dataIndex: 'status', key: 'status', width: 120 },
  { title: '创建时间', dataIndex: 'created_at', key: 'created_at', width: 180 },
  { title: '更新时间', dataIndex: 'updated_at', key: 'updated_at', width: 180 },
  { title: '操作', key: 'action', width: 340 },
]
const statusOptions = useDictionaryOptions('sys_status', true)
const query = reactive({ name: '', code: '', status: undefined })
const admin = useAdminStore()
const rows = ref([])
const loading = ref(false)
const formRef = ref()
const permissionRef = ref()
const collapsedIds = ref(new Set())
const initialized = ref(false)
const rowMap = computed(() => new Map(rows.value.map((row) => [row.id, row])))
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

function flattenRoles(roles, depth = 0) {
  return roles.flatMap((role) => {
    const children = role.children || []
    const row = {
      ...role,
      parentId: role.parent_id,
      depth,
      hasChildren: children.length > 0,
    }
    delete row.children
    return [row, ...flattenRoles(children, depth + 1)]
  })
}

function toggleRow(id) {
  const next = new Set(collapsedIds.value)
  if (next.has(id)) next.delete(id)
  else next.add(id)
  collapsedIds.value = next
}

function search() {
  loadRoles()
}

function reset() {
  Object.assign(query, { name: '', code: '', status: undefined })
  loadRoles()
}

function createRole() {
  formRef.value?.openCreate()
}

function configurePermission(record) {
  permissionRef.value?.open(record)
}

function createChild(record) {
  formRef.value?.openCreate(record)
}

function editRole(record) {
  formRef.value?.openEdit(record.id)
}

function canOperate(record) {
  return record.code !== 'super_admin' || admin.profile?.role_code === 'super_admin'
}

async function loadRoles() {
  loading.value = true
  try {
    const data = await getRoleList({
      name: query.name || undefined,
      code: query.code || undefined,
      status: query.status,
    })
    rows.value = flattenRoles(data?.list || [])
    if (!initialized.value) {
      collapsedIds.value = new Set(rows.value.filter((row) => row.hasChildren).map((row) => row.id))
      initialized.value = true
    }
  } catch (error) {
    message.error(error.message)
  } finally {
    loading.value = false
  }
}

function removeRole(record) {
  Modal.confirm({
    title: '确认删除角色？',
    content: `删除“${record.name}”后无法恢复。`,
    okText: '删除',
    okType: 'danger',
    cancelText: '取消',
    async onOk() {
      try {
        await deleteRole(record.id)
        message.success('角色删除成功')
        await loadRoles()
      } catch (error) {
        message.error(error.message)
        return Promise.reject(error)
      }
    },
  })
}

onMounted(loadRoles)
</script>

<style scoped lang="less">
.role-page {
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

.toolbar :deep(.ant-input) {
  width: 190px;
}

.toolbar :deep(.ant-select) {
  width: 132px;
}

.toolbar :deep(.create-action) {
  margin-right: 0;
  margin-left: auto;
}

.role-table :deep(.ant-table-thead > tr > th) {
  height: 48px;
  padding: 0 16px;
  color: #494949;
  font-size: 13px;
  font-weight: 600;
  background: #fafafa;
  border-bottom: 1px solid #e8e8e8;
}

.role-table :deep(.ant-table-tbody > tr > td) {
  height: 52px;
  padding: 10px 16px;
  color: #3b3b3b;
  border-bottom: 1px solid #f0f0f0;
}

.role-table :deep(.ant-table-tbody > tr:hover > td) {
  background: #fafafa;
}

.role-name {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 0;
  font-weight: 500;
}

.role-icon {
  color: #7b8493;
  font-size: 16px;
}

.expand-button {
  display: inline-flex;
  width: 22px;
  height: 22px;
  align-items: center;
  justify-content: center;
  padding: 0;
  color: #8c8c8c;
  font-size: 10px;
  background: transparent;
  border: 0;
  cursor: pointer;
}

.expand-button :deep(.anticon) {
  transition: transform 160ms ease;
}

.expand-button--open :deep(.anticon) {
  transform: rotate(90deg);
}

.code-text {
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
