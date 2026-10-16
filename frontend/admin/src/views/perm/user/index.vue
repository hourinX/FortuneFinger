<template>
  <section class="user-layout-test">
    <aside class="user-layout-test__org">
      <div class="user-layout-test__org-header">
        <span class="user-layout-test__org-title">组织结构</span>
      </div>
      <a-input v-model:value="orgKeyword" allow-clear placeholder="搜索组织名称">
        <template #prefix><search-outlined /></template>
      </a-input>
      <div class="org-tree-wrap">
        <a-spin :spinning="orgLoading">
          <a-tree
            v-model:expanded-keys="expandedOrgKeys"
            v-model:selected-keys="selectedOrgKeys"
            block-node
            :show-line="{ showLeafIcon: false }"
            :tree-data="orgTreeData"
            @select="selectOrganization"
          />
        </a-spin>
      </div>
    </aside>
    <main class="user-layout-test__content">
      <a-form class="user-layout-test__toolbar" layout="inline" :model="query" @finish="search">
        <a-form-item>
          <a-input v-model:value="query.username" allow-clear placeholder="用户名">
            <template #prefix><search-outlined /></template>
          </a-input>
        </a-form-item>
        <a-form-item>
          <a-input v-model:value="query.nickname" allow-clear placeholder="昵称" />
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
          <a-space v-permission="'perm:user:list'" :size="8">
            <a-button type="primary" html-type="submit">查询</a-button>
            <a-button @click="reset">重置</a-button>
          </a-space>
        </a-form-item>
        <a-form-item class="user-layout-test__create">
          <a-button
            v-permission="'perm:user:add'"
            type="primary"
            @click="formRef?.openCreate()"
          >
            <plus-outlined />新增
          </a-button>
        </a-form-item>
      </a-form>

      <a-table
        class="user-layout-test__table"
        :columns="columns"
        :data-source="rows"
        :loading="loading"
        :pagination="pagination"
        row-key="id"
        :scroll="{ x: 1200 }"
        @change="changePage"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'user'">
            <div class="user-info">
              <button
                class="avatar-preview-trigger"
                type="button"
                :disabled="!record.avatar"
                :aria-label="record.avatar ? `预览${record.nickname}的头像` : undefined"
                @click="openAvatarPreview(record)"
              >
                <a-avatar :size="38" :src="record.avatar">
                  {{ record.nickname.slice(0, 1) }}
                </a-avatar>
              </button>
              <div class="user-info__identity">
                <strong>{{ record.nickname }}</strong>
                <span>@{{ record.username }}</span>
              </div>
            </div>
          </template>
          <template v-else-if="column.key === 'contact'">
            <div class="contact-info">
              <span><mail-outlined />{{ record.email || '未填写邮箱' }}</span>
              <span><phone-outlined />{{ record.mobile || '未填写手机号' }}</span>
            </div>
          </template>
          <template v-else-if="column.key === 'status'">
            <StatusPulse :status="record.status" />
          </template>
          <template v-else-if="column.key === 'department'">
            <span v-if="record.department_name" class="department-name">
              <apartment-outlined />{{ record.department_name }}
            </span>
            <span v-else-if="record.role_code === 'super_admin'" class="empty-role">无需分配</span>
            <span v-else class="empty-role">未分配</span>
          </template>
          <template v-else-if="column.key === 'role'">
            <a-tag v-if="record.role_name" color="blue">{{ record.role_name }}</a-tag>
            <a-tag v-else-if="record.role_code">角色已禁用/删除</a-tag>
            <span v-else class="empty-role">未分配</span>
          </template>
          <template v-else-if="column.key === 'createdAt'">
            <span class="created-time">{{ record.created_at }}</span>
          </template>
          <template v-else-if="column.key === 'action'">
            <a-space :size="14">
              <template v-if="canOperate(record)">
                <a
                  v-permission="'perm:user:role'"
                  class="action-link"
                  @click="assignRole(record)"
                >
                  <safety-outlined />分配角色
                </a>
                <a
                  v-permission="'perm:user:edit'"
                  class="action-link"
                  @click="formRef?.openEdit(record.id)"
                >
                  <edit-outlined />编辑
                </a>
                <a
                  v-permission="'perm:user:del'"
                  class="danger-link"
                  @click="removeUser(record)"
                >
                  <delete-outlined />删除
                </a>
              </template>
              <span v-else class="empty-role">—</span>
            </a-space>
          </template>
        </template>
        <template #emptyText><a-empty description="暂无用户数据" /></template>
      </a-table>
    </main>
  </section>

  <UserForm ref="formRef" :organizations="orgOptions" @success="loadUsers" />
  <AssignRole ref="assignRoleRef" @success="loadUsers" />
  <ImagePreview
    v-model:visible="avatarPreviewVisible"
    :src="avatarPreviewSrc"
    :alt="avatarPreviewAlt"
  />
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import {
  ApartmentOutlined,
  DeleteOutlined,
  EditOutlined,
  MailOutlined,
  PhoneOutlined,
  PlusOutlined,
  SafetyOutlined,
  SearchOutlined,
} from '@ant-design/icons-vue'
import { Modal, message } from 'ant-design-vue'
import ImagePreview from '@/components/ImagePreview.vue'
import StatusPulse from '@/components/StatusPulse.vue'
import { useDictionaryOptions } from '@/hooks/useDictionary'
import { getOrgOptions } from '@/api/org'
import { deleteUser, getUserList } from '@/api/user'
import { useAdminStore } from '@/store'
import UserForm from './components/form.vue'
import AssignRole from './components/assign-role.vue'

const columns = [
  { title: '用户信息', key: 'user', width: 200 },
  { title: '联系方式', key: 'contact', width: 230 },
  { title: '所属部门', key: 'department', width: 150 },
  { title: '所属角色', key: 'role', width: 130 },
  { title: '状态', dataIndex: 'status', key: 'status', width: 90 },
  { title: '创建时间', key: 'createdAt', width: 170 },
  { title: '操作', key: 'action', width: 260 },
]
const statusOptions = useDictionaryOptions('sys_status', true)
const query = reactive({ username: '', nickname: '', status: undefined })
const admin = useAdminStore()
const orgKeyword = ref('')
const selectedOrgKeys = ref(['all'])
const expandedOrgKeys = ref([])
const orgOptions = ref([])
const orgLoading = ref(false)
const orgTreeData = computed(() => [
  {
    key: 'all',
    title: '全部组织',
    children: filterOrgOptions(orgOptions.value, orgKeyword.value),
  },
])
const page = reactive({ current: 1, pageSize: 10 })
const formRef = ref()
const assignRoleRef = ref()
const avatarPreviewVisible = ref(false)
const avatarPreviewSrc = ref('')
const avatarPreviewAlt = ref('')
const rows = ref([])
const total = ref(0)
const loading = ref(false)
const pagination = computed(() => ({
  current: page.current,
  pageSize: page.pageSize,
  total: total.value,
  showSizeChanger: true,
  showQuickJumper: true,
  showTotal: (total) => `共 ${total} 条`,
  locale: { items_per_page: '条/页', jump_to: '跳转至' },
}))

async function search() {
  page.current = 1
  await loadUsers()
}

function reset() {
  Object.assign(query, { username: '', nickname: '', status: undefined })
  search()
}

function changePage(next) {
  page.current = next.current
  page.pageSize = next.pageSize
  loadUsers()
}

function assignRole(record) {
  assignRoleRef.value?.open(record)
}

function selectOrganization(keys) {
  selectedOrgKeys.value = keys.length ? keys : ['all']
  page.current = 1
  loadUsers()
}

function canOperate(record) {
  return record.role_code !== 'super_admin' || admin.profile?.role_code === 'super_admin'
}

function openAvatarPreview(record) {
  if (!record.avatar) return
  avatarPreviewSrc.value = record.avatar
  avatarPreviewAlt.value = `${record.nickname}的头像`
  avatarPreviewVisible.value = true
}

async function loadUsers() {
  loading.value = true
  try {
    const data = await getUserList({
      page: page.current,
      pageSize: page.pageSize,
      username: query.username || undefined,
      nickname: query.nickname || undefined,
      status: query.status,
      org_id: selectedOrgKeys.value[0] === 'all' ? undefined : selectedOrgKeys.value[0],
    })
    rows.value = data?.list || []
    total.value = data?.total || 0
  } catch (error) {
    message.error(error.message)
  } finally {
    loading.value = false
  }
}

function mapOrgOptions(options) {
  return options.map((option) => ({
    key: option.id,
    title: option.name,
    children: mapOrgOptions(option.children || []),
  }))
}

function filterOrgOptions(options, keyword) {
  if (!keyword) return options
  const normalizedKeyword = keyword.toLocaleLowerCase()
  return options.flatMap((option) => {
    const children = filterOrgOptions(option.children || [], keyword)
    const matched = option.title.toLocaleLowerCase().includes(normalizedKeyword)
    if (!matched && !children.length) return []
    return [{ ...option, children }]
  })
}

function collectOrgKeys(options, keys = ['all']) {
  options.forEach((option) => {
    if (option.children?.length) keys.push(option.key)
    collectOrgKeys(option.children || [], keys)
  })
  return keys
}

async function loadOrgOptions() {
  orgLoading.value = true
  try {
    const data = await getOrgOptions()
    orgOptions.value = mapOrgOptions(data?.list || [])
    expandedOrgKeys.value = collectOrgKeys(orgOptions.value)
  } catch (error) {
    message.error(error.message)
  } finally {
    orgLoading.value = false
  }
}

function removeUser(record) {
  Modal.confirm({
    title: '确认删除用户？',
    content: `删除“${record.nickname}”后无法恢复。`,
    okText: '删除',
    okType: 'danger',
    cancelText: '取消',
    async onOk() {
      try {
        await deleteUser(record.id)
        message.success('用户删除成功')
        if (rows.value.length === 1 && page.current > 1) page.current -= 1
        await loadUsers()
      } catch (error) {
        message.error(error.message)
        return Promise.reject(error)
      }
    },
  })
}
onMounted(() => {
  loadOrgOptions()
  loadUsers()
})
</script>

<style scoped lang="less">
.user-layout-test {
  display: grid;
  grid-template-columns: 236px minmax(0, 1fr);
  width: 100%;
  max-width: 100%;
  min-width: 0;
  min-height: calc(100vh - 168px);
  gap: 18px;
  overflow: hidden;
}

.user-layout-test__org {
  display: flex;
  flex-direction: column;
  width: 100%;
  height: 100%;
  min-width: 0;
  padding: 16px 12px;
  overflow: hidden;
  background: #f8f8f8;
  border: 1px solid #e9e9e9;
  border-radius: 8px;
}

.user-layout-test__org-header {
  display: flex;
  height: 34px;
  align-items: flex-start;
  justify-content: space-between;
  padding: 0 2px;
}

.user-layout-test__org-title {
  color: #2f2f2f;
  font-size: 15px;
  font-weight: 600;
}

.user-layout-test__org :deep(.ant-input-affix-wrapper) {
  flex: 0 0 auto;
}

.user-layout-test__org .org-tree-wrap {
  min-height: 0;
  flex: 1 1 auto;
  margin-top: 16px;
  overflow: auto;
  scrollbar-width: none;
  -ms-overflow-style: none;
}

.user-layout-test__org .org-tree-wrap::-webkit-scrollbar {
  display: none;
}

.user-layout-test__org :deep(.ant-tree),
.user-layout-test__org :deep(.ant-tree-list),
.user-layout-test__org :deep(.ant-tree-list-holder),
.user-layout-test__org :deep(.ant-tree-list-holder-inner) {
  background: transparent;
}

.user-layout-test__org :deep(.ant-tree-node-content-wrapper) {
  background: transparent;
}

.user-layout-test__org :deep(.ant-tree-node-content-wrapper:hover) {
  background: rgba(228, 61, 54, 0.06);
}

.user-layout-test__org :deep(.ant-tree-node-selected) {
  color: var(--color-primary);
  background: rgba(228, 61, 54, 0.1) !important;
}

.user-layout-test__content {
  width: 100%;
  max-width: 100%;
  min-width: 0;
  overflow: hidden;
}

.user-layout-test__toolbar {
  display: flex;
  width: 100%;
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

.user-layout-test__create {
  margin-right: 0 !important;
  margin-left: auto !important;
}

.user-layout-test__table {
  width: 100%;
  max-width: 100%;
  min-width: 0;

  :deep(.ant-spin-nested-loading),
  :deep(.ant-spin-container),
  :deep(.ant-table),
  :deep(.ant-table-container) {
    width: 100%;
    max-width: 100%;
    min-width: 0;
  }

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
    height: 68px;
    padding: 10px 16px;
    color: #3b3b3b;
    border-bottom: 1px solid #f0f0f0;
  }

  :deep(.ant-table-tbody > tr:hover > td) {
    background: #fafafa;
  }
}

.user-layout-test__content .user-info {
  display: flex;
  align-items: center;
  gap: 12px;
  min-width: 0;

  :deep(.ant-avatar) {
    flex: 0 0 auto;
    color: var(--color-primary);
    background: var(--color-primary-light);
    box-shadow: 0 0 0 3px #fafafa;
  }

  .user-info__identity {
    display: flex;
    min-width: 0;
    flex-direction: column;
    gap: 3px;
  }

  strong {
    overflow: hidden;
    color: #2f2f2f;
    font-size: 14px;
    font-weight: 600;
    white-space: nowrap;
    text-overflow: ellipsis;
  }

  span {
    overflow: hidden;
    color: #929292;
    font-size: 12px;
    white-space: nowrap;
    text-overflow: ellipsis;
  }
}

.avatar-preview-trigger {
  display: inline-flex;
  flex: 0 0 auto;
  padding: 0;
  background: transparent;
  border: 0;
  border-radius: 50%;
  cursor: pointer;

  &:disabled {
    cursor: default;
  }
}

.contact-info {
  display: flex;
  min-width: 0;
  flex-direction: column;
  align-items: flex-start;
}

.contact-info {
  gap: 6px;

  span {
    display: flex;
    max-width: 100%;
    align-items: center;
    gap: 7px;
    overflow: hidden;
    color: #555;
    font-size: 13px;
    line-height: 17px;
    white-space: nowrap;
    text-overflow: ellipsis;
  }

  :deep(.anticon) {
    flex: 0 0 auto;
    color: #a0a0a0;
    font-size: 13px;
  }
}

.department-name {
  display: flex;
  max-width: 100%;
  align-items: center;
  gap: 7px;
  overflow: hidden;
  color: #4a4a4a;
  font-size: 13px;
  white-space: nowrap;
  text-overflow: ellipsis;

  :deep(.anticon) {
    flex: 0 0 auto;
    color: #8a94a3;
  }
}

.user-layout-test__content :deep(.ant-tag) {
  margin: 0;
  line-height: 20px;
}

.created-time {
  color: #6f6f6f;
  font-size: 13px;
  white-space: nowrap;
}

.user-layout-test__content .action-link {
  color: #1677ff;
}

.user-layout-test__content .danger-link {
  color: #cf3f36;
}

.user-layout-test__content .empty-role {
  color: #a6a6a6;
}

.user-layout-test__content .action-link :deep(.anticon),
.user-layout-test__content .danger-link :deep(.anticon) {
  margin-right: 4px;
}
</style>
