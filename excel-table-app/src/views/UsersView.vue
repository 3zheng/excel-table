<template>
  <div>
    <el-card class="search-card">
      <div class="toolbar">
        <div class="toolbar-left">
          <el-input
            v-model="searchUsername"
            :placeholder="t('searchUsername')"
            clearable
            style="width:200px"
            @clear="applyFilter"
            @keyup.enter="applyFilter"
          />
          <el-button type="primary" @click="applyFilter">{{ t('search') }}</el-button>
          <el-button @click="resetSearch">{{ t('reset') }}</el-button>
        </div>
        <div class="toolbar-right" v-if="isAdmin">
          <el-button type="success" @click="openCreate">+ {{ t('newUser') }}</el-button>
        </div>
      </div>
    </el-card>

    <el-card>
      <el-table :data="filteredList" border style="width:100%" v-loading="loading">
        <el-table-column :label="t('username')" prop="username" width="160" />
        <el-table-column :label="t('userRole')" prop="user_role" width="160">
          <template #default="{ row }">
            {{ roleLabel(row.user_role) }}
          </template>
        </el-table-column>
        <el-table-column :label="t('canEditProducts')" width="120" align="center">
          <template #default="{ row }">
            <el-tag :type="row.can_edit_products ? 'success' : 'info'" size="small">
              {{ row.can_edit_products ? t('yes') : t('no') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('region')" min-width="180">
          <template #default="{ row }">
            <template v-if="row.regions?.length">
              <el-tag v-for="r in row.regions" :key="r" size="small" style="margin:2px">
                {{ regionLabel(r) }}
              </el-tag>
            </template>
            <span v-else>—</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('createdAt')" prop="created_at" width="170" />
        <el-table-column :label="t('operation')" width="160" fixed="right">
          <template #default="{ row }">
            <!-- 可编辑的情况 -->
            <template v-if="canEditRow(row)">
              <el-button size="small" @click="openEdit(row)">{{ t('edit') }}</el-button>
              <el-button
                v-if="isAdmin && row.user_role !== 'admin'"
                size="small"
                type="danger"
                @click="handleDelete(row)"
              >
                {{ t('delete') }}
              </el-button>
            </template>
            <!-- 其他管理员只读 -->
            <el-tag v-else-if="row.user_role === 'admin'" type="info" size="small">{{ t('readonly') }}</el-tag>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 编辑 / 新建弹窗 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="520px"
      :close-on-click-modal="false"
    >
      <el-form :model="form" label-width="120px">
        <!-- 用户名 -->
        <el-form-item :label="t('username')" required>
          <el-input
            v-model="form.username"
            autocomplete="off"
            :disabled="usernameDisabled"
          />
        </el-form-item>

        <!-- 密码 -->
        <el-form-item :label="t('password')" :required="!isEdit || isSelfPasswordOnly">
          <el-input
            v-model="form.password"
            type="password"
            show-password
            autocomplete="new-password"
            :placeholder="isEdit ? t('passwordOptionalHint') : ''"
          />
        </el-form-item>

        <!-- 以下字段仅管理员编辑非自己时显示 -->
        <template v-if="showFullForm">
          <el-form-item :label="t('userRole')" required>
            <el-select v-model="form.user_role" style="width:100%" :disabled="isEditingSelf">
              <el-option
                v-for="opt in editableRoleOptions"
                :key="opt.value"
                :label="t(opt.labelKey)"
                :value="opt.value"
              />
            </el-select>
          </el-form-item>

          <el-form-item :label="t('canEditProducts')">
            <el-switch v-model="form.can_edit_products" :disabled="isEditingSelf" />
          </el-form-item>

          <el-form-item v-if="needsRegions" :label="t('region')" required>
            <el-select
              v-model="form.regions"
              multiple
              style="width:100%"
              :placeholder="t('selectRegions')"
              :disabled="isEditingSelf"
            >
              <el-option
                v-for="r in regionOptions"
                :key="r.value"
                :label="t(r.labelKey)"
                :value="r.value"
              />
            </el-select>
          </el-form-item>
        </template>
      </el-form>

      <template #footer>
        <el-button @click="dialogVisible = false">{{ t('cancel') }}</el-button>
        <el-button type="primary" @click="submitForm" :loading="submitting">
          {{ t('save') }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { request } from '@/utils/request'
import { REGIONS, BUYERS } from '@/constants'

const { t } = useI18n()

const user = JSON.parse(localStorage.getItem('user') || '{}')
const isAdmin = computed(() => user.user_role === 'admin')
const currentUserId = computed(() => user.id)

const loading = ref(false)
const submitting = ref(false)
const dialogVisible = ref(false)
const isEdit = ref(false)
const editingId = ref<number | null>(null)

const searchUsername = ref('')
const list = ref<any[]>([])
const filterKeyword = ref('')

const regionOptions = REGIONS //地区列表

const regionKeys: Record<string, string> = Object.fromEntries(
  regionOptions.map(r => [r.value, r.labelKey])
)
const regionLabel = (r: string) => t(regionKeys[r] || r)

const roleOptions = [
  { value: 'admin', labelKey: 'roleAdmin' },
  { value: 'export_input', labelKey: 'roleExportInput' },
  { value: 'export_review', labelKey: 'roleExportReview' },
  { value: 'import_input', labelKey: 'roleImportInput' },
  { value: 'import_review', labelKey: 'roleImportReview' },
]

// 新建/编辑别人时不出现 admin
const editableRoleOptions = computed(() =>
  roleOptions.filter(r => r.value !== 'admin')
)

const roleLabelKeys: Record<string, string> = Object.fromEntries(
  roleOptions.map(r => [r.value, r.labelKey])
)
const roleLabel = (role: string) => t(roleLabelKeys[role] || role)

const form = ref({
  username: '',
  password: '',
  user_role: 'export_input',
  can_edit_products: false,
  regions: [] as string[],
})

const needsRegions = computed(() =>
  form.value.user_role === 'import_input' || form.value.user_role === 'import_review'
)

// 是否正在编辑自己
const isEditingSelf = computed(() => isEdit.value && editingId.value === currentUserId.value)

// 非管理员只能改自己密码
const isSelfPasswordOnly = computed(() => !isAdmin.value)

// 用户名是否禁用
const usernameDisabled = computed(() => {
  // 非管理员：用户名只读
  if (!isAdmin.value) return true
  // 管理员编辑自己：可以改用户名
  return false
})

// 是否显示完整表单（角色、产品权限、地区）
const showFullForm = computed(() => {
  // 非管理员不显示
  if (!isAdmin.value) return false
  // 管理员编辑自己时也不显示（角色锁定）
  if (isEditingSelf.value) return false
  return true
})

const dialogTitle = computed(() => {
  if (!isEdit.value) return t('newUser')
  if (isSelfPasswordOnly.value) return t('editPassword') 
  return t('editUser')
})

// 当前行是否可编辑
const canEditRow = (row: any) => {
  if (isAdmin.value) {
    // 管理员：自己可编辑，其他管理员不可编辑
    if (row.user_role === 'admin') {
      return row.id === currentUserId.value
    }
    return true
  }
  // 非管理员：只能编辑自己
  return row.id === currentUserId.value
}

const filteredList = computed(() => {
  if (!filterKeyword.value) return list.value
  const kw = filterKeyword.value.toLowerCase()
  return list.value.filter(u => u.username.toLowerCase().includes(kw))
})

const fetchUsers = async () => {
  loading.value = true
  try {
    const res = await request('/api/users')
    if (res.ok) {
      list.value = await res.json()
    }
  } catch {
  } finally {
    loading.value = false
  }
}

const applyFilter = () => {
  filterKeyword.value = searchUsername.value.trim()
}

const resetSearch = () => {
  searchUsername.value = ''
  filterKeyword.value = ''
}

const openCreate = () => {
  isEdit.value = false
  editingId.value = null
  form.value = {
    username: '',
    password: '',
    user_role: 'export_input',
    can_edit_products: false,
    regions: [],
  }
  dialogVisible.value = true
}

const openEdit = (row: any) => {
  isEdit.value = true
  editingId.value = row.id
  form.value = {
    username: row.username,
    password: '',
    user_role: row.user_role,
    can_edit_products: !!row.can_edit_products,
    regions: [...(row.regions || [])],
  }
  dialogVisible.value = true
}

const submitForm = async () => {
  // 非管理员只校验密码
  if (!isAdmin.value) {
    if (!form.value.password) {
      ElMessage.error(t('passwordRequired') || '请输入新密码')
      return
    }
  } else {
    if (!form.value.username.trim()) {
      ElMessage.error(t('usernameRequired'))
      return
    }
    if (!isEdit.value && !form.value.password) {
      ElMessage.error(t('passwordRequired'))
      return
    }
    if (showFullForm.value && needsRegions.value && form.value.regions.length === 0) {
      ElMessage.error(t('regionsRequired'))
      return
    }
  }

  submitting.value = true
  try {
    const payload: Record<string, unknown> = {}

    if (isAdmin.value) {
      payload.username = form.value.username.trim()
      payload.user_role = isEditingSelf.value ? 'admin' : form.value.user_role
      payload.can_edit_products = isEditingSelf.value ? true : form.value.can_edit_products
      payload.regions = showFullForm.value && needsRegions.value ? form.value.regions : []
    }

    if (form.value.password) {
      payload.password = form.value.password
    }

    const url = isEdit.value ? `/api/users/${editingId.value}` : '/api/users'
    const method = isEdit.value ? 'PUT' : 'POST'
    const res = await request(url, { method, body: JSON.stringify(payload) })

    if (res.ok) {
      ElNotification({
        title: t('saveSuccess'),
        type: 'success',
        position: 'top-right',
        duration: 2500,
      })
      dialogVisible.value = false
      fetchUsers()

      // 如果改了自己的用户名，顺便更新 localStorage
      if (isEditingSelf.value && payload.username) {
        const u = JSON.parse(localStorage.getItem('user') || '{}')
        u.username = payload.username
        localStorage.setItem('user', JSON.stringify(u))
      }
    } else {
      const data = await res.json()
      ElMessage.error(data.error || t('saveFailed'))
    }
  } catch {
  } finally {
    submitting.value = false
  }
}

const handleDelete = async (row: any) => {
  try {
    await ElMessageBox.confirm(
      t('deleteUserConfirm', { name: row.username }),
      t('warning'),
      { type: 'warning', confirmButtonText: t('delete'), cancelButtonText: t('cancel') }
    )
    const res = await request(`/api/users/${row.id}`, { method: 'DELETE' })
    if (res.ok) {
      ElNotification({
        title: t('deleteSuccess'),
        type: 'success',
        position: 'top-right',
        duration: 2500,
      })
      fetchUsers()
    } else {
      const data = await res.json()
      ElMessage.error(data.error || t('saveFailed'))
    }
  } catch { }
}

onMounted(() => {
  fetchUsers()
})
</script>

<style scoped>
.search-card {
  margin-bottom: 16px;
}

.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 8px;
}

.toolbar-left {
  display: flex;
  gap: 8px;
  align-items: center;
  flex-wrap: wrap;
}

.toolbar-right {
  display: flex;
  gap: 8px;
}
</style>