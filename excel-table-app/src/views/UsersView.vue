<template>
  <div>
    <el-card class="search-card">
      <div class="toolbar">
        <div class="toolbar-left">
          <el-input v-model="searchUsername" :placeholder="t('searchUsername')" clearable style="width:200px"
            @clear="applyFilter" @keyup.enter="applyFilter" />
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
        <el-table-column v-if="isAdmin" :label="t('operation')" width="140" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="openEdit(row)">{{ t('edit') }}</el-button>
            <el-button size="small" type="danger" @click="handleDelete(row)">{{ t('delete') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="isEdit ? t('editUser') : t('newUser')"
      width="520px" :close-on-click-modal="false">
      <el-form :model="form" label-width="120px">
        <el-form-item :label="t('username')" required>
          <el-input v-model="form.username" autocomplete="off" />
        </el-form-item>
        <el-form-item :label="t('password')" :required="!isEdit">
          <el-input v-model="form.password" type="password" show-password autocomplete="new-password"
            :placeholder="isEdit ? t('passwordOptionalHint') : ''" />
        </el-form-item>
        <el-form-item :label="t('userRole')" required>
          <el-select v-model="form.user_role" style="width:100%">
            <el-option v-for="opt in roleOptions" :key="opt.value" :label="t(opt.labelKey)" :value="opt.value" />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('canEditProducts')">
          <el-switch v-model="form.can_edit_products" />
        </el-form-item>
        <el-form-item v-if="needsRegions" :label="t('region')" required>
          <el-select v-model="form.regions" multiple style="width:100%" :placeholder="t('selectRegions')">
            <el-option v-for="r in regionOptions" :key="r.value" :label="t(r.labelKey)" :value="r.value" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ t('cancel') }}</el-button>
        <el-button type="primary" @click="submitForm" :loading="submitting">{{ t('save') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { request } from '@/utils/request'

const { t } = useI18n()

const user = JSON.parse(localStorage.getItem('user') || '{}')
const isAdmin = computed(() => user.user_role === 'admin')

const loading = ref(false)
const submitting = ref(false)
const dialogVisible = ref(false)
const isEdit = ref(false)
const editingId = ref<number | null>(null)

const searchUsername = ref('')
const list = ref<any[]>([])
const filterKeyword = ref('')

const regionOptions = [
  { value: '玻利维亚', labelKey: 'bolivia' },
  { value: '秘鲁', labelKey: 'peru' },
  { value: '智利', labelKey: 'chile' },
  { value: '西班牙', labelKey: 'spain' },
  { value: '美国', labelKey: 'usa' },
  { value: '南美', labelKey: 'southAmerica' },
]

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
  if (!form.value.username.trim()) {
    ElMessage.error(t('usernameRequired'))
    return
  }
  if (!isEdit.value && !form.value.password) {
    ElMessage.error(t('passwordRequired'))
    return
  }
  if (needsRegions.value && form.value.regions.length === 0) {
    ElMessage.error(t('regionsRequired'))
    return
  }

  submitting.value = true
  try {
    const payload: Record<string, unknown> = {
      username: form.value.username.trim(),
      user_role: form.value.user_role,
      can_edit_products: form.value.can_edit_products,
      regions: needsRegions.value ? form.value.regions : [],
    }
    if (form.value.password) {
      payload.password = form.value.password
    }

    const url = isEdit.value ? `/api/users/${editingId.value}` : '/api/users'
    const method = isEdit.value ? 'PUT' : 'POST'
    const res = await request(url, { method, body: JSON.stringify(payload) })
    if (res.ok) {
      ElNotification({ title: t('saveSuccess'), type: 'success', position: 'top-right', duration: 2500 })
      dialogVisible.value = false
      fetchUsers()
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
      ElNotification({ title: t('deleteSuccess'), type: 'success', position: 'top-right', duration: 2500 })
      fetchUsers()
    } else {
      const data = await res.json()
      ElMessage.error(data.error || t('saveFailed'))
    }
  } catch { }
}

onMounted(() => {
  if (isAdmin.value) {
    fetchUsers()
  }
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
