<template>
  <div>
    <!-- 强制选择地区（未选择时显示） -->
    <div v-if="!currentRegion" class="region-select-mask">
      <el-card class="region-select-card">
        <template #header>
          <div style="text-align:center;font-size:18px;font-weight:600">
            {{ t('region') }}
          </div>
        </template>
        <div style="text-align:center;margin-bottom:20px;color:#606266">
          {{ t('selectRegionFirst') }}
        </div>
        <el-select
          v-model="tempRegion"
          :placeholder="t('region')"
          size="large"
          style="width:100%"
          @change="confirmRegion"
        >
          <el-option v-for="r in REGIONS" :key="r.value" :label="t(r.labelKey)" :value="r.value" />
        </el-select>
      </el-card>
    </div>
    <!-- 已选择地区后显示正常内容 -->
    <template v-else>
      <!-- 顶部操作栏 -->
      <el-card class="search-card">
        <div class="toolbar">
          <div class="toolbar-left">
            <el-select
              v-model="currentRegion"
              :placeholder="t('region')"
              style="width:140px"
              @change="fetchProducts"
            >
              <el-option v-for="r in REGIONS" :key="r.value" :label="t(r.labelKey)" :value="r.value" />
            </el-select>
            <el-input
              v-model="searchModelNo"
              placeholder="搜索型号"
              clearable
              style="width:180px"
              @clear="fetchProducts"
              @keyup.enter="fetchProducts"
            />
            <el-button type="primary" @click="fetchProducts">{{ t('search') }}</el-button>
            <el-button @click="resetSearch">{{ t('reset') }}</el-button>
          </div>
          <div class="toolbar-right" v-if="canEdit">
            <el-button type="success" @click="openCreate">+ {{ t('newProduct') }}</el-button>
            <el-button type="warning" @click="openCopyDialog">
              {{ t('copyRegion') }}
            </el-button>
          </div>
        </div>
      </el-card>
      <!-- 产品列表 -->
      <el-card>
        <el-table :data="filteredList" border style="width:100%" v-loading="loading">
          <el-table-column :label="t('modelNo')" prop="model_no" width="140" fixed="left" />
          <el-table-column :label="t('brand')" prop="brand" width="110" />
          <el-table-column :label="t('commodities')" prop="commodities" min-width="130" />
          <el-table-column :label="t('descriptions')" prop="descriptions" min-width="160" />
          <el-table-column :label="t('unit')" prop="unit" width="70" />
          <el-table-column :label="t('exportRefPrice')" prop="export_ref_price" width="120">
            <template #default="{ row }">
              {{ row.export_ref_price != null ? row.export_ref_price : '—' }}
            </template>
          </el-table-column>
          <el-table-column :label="t('importRefCost')" prop="import_ref_cost" width="120">
            <template #default="{ row }">
              {{ row.import_ref_cost != null ? row.import_ref_cost : '—' }}
            </template>
          </el-table-column>
          <el-table-column :label="t('updatedAt')" prop="updated_at" width="160" />
          <el-table-column :label="t('updatedBy')" prop="updated_by" width="100" />
          <el-table-column v-if="canEdit" :label="t('operation')" width="140" fixed="right">
            <template #default="{ row }">
              <el-button size="small" @click="openEdit(row)">{{ t('edit') }}</el-button>
              <el-button size="small" type="danger" @click="handleDelete(row)">{{ t('delete') }}</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-card>
      <!-- 新增/编辑弹窗 -->
      <el-dialog
        v-model="dialogVisible"
        :title="isEdit ? t('editProduct') : t('newProduct')"
        width="600px"
        :close-on-click-modal="false"
      >
        <el-form :model="form" label-width="130px">
          <el-form-item :label="t('modelNo')">
            <el-input v-model="form.model_no" :disabled="isEdit" />
          </el-form-item>
          <el-form-item :label="t('brand')">
            <el-input v-model="form.brand" />
          </el-form-item>
          <el-form-item :label="t('commodities')">
            <el-input v-model="form.commodities" />
          </el-form-item>
          <el-form-item :label="t('descriptions')">
            <el-input v-model="form.descriptions" />
          </el-form-item>
          <el-form-item :label="t('unit')">
            <el-input v-model="form.unit" />
          </el-form-item>
          <el-divider>{{ t('priceInfo') }} - {{ regionLabel(currentRegion) }}</el-divider>
          <el-form-item :label="t('exportRefPrice')">
            <el-input-number
              v-model="form.export_ref_price"
              :precision="4"
              :min="0"
              style="width:100%"
              :controls="false"
            />
          </el-form-item>
          <el-form-item :label="t('importRefCost')">
            <el-input-number
              v-model="form.import_ref_cost"
              :precision="4"
              :min="0"
              style="width:100%"
              :controls="false"
            />
          </el-form-item>
        </el-form>
        <template #footer>
          <el-button @click="dialogVisible = false">{{ t('cancel') }}</el-button>
          <el-button type="primary" @click="submitForm" :loading="submitting">{{ t('save') }}</el-button>
        </template>
      </el-dialog>
      <!-- 一键复制地区弹窗 -->
      <el-dialog v-model="copyDialogVisible" :title="t('copyRegion')" width="400px">
        <el-form label-width="100px">
          <el-alert
            type="info"
            :closable="false"
            style="margin-bottom: 16px"
          >
            {{ t('currentlyManaging', { region: regionLabel(currentRegion) }) }}
          </el-alert>
          <el-form-item :label="t('fromRegion')">
            <el-select v-model="copyFrom" style="width:100%">
              <el-option v-for="r in REGIONS" :key="r.value" :label="t(r.labelKey)" :value="r.value" />
            </el-select>
          </el-form-item>
          <el-form-item :label="t('toRegion')">
            <el-select v-model="copyTo" style="width:100%">
              <el-option v-for="r in REGIONS" :key="r.value" :label="t(r.labelKey)" :value="r.value" />
            </el-select>
          </el-form-item>
          <el-alert type="info" :closable="false" style="margin-top:8px">
            {{ t('copyHint') }}
          </el-alert>
        </el-form>
        <template #footer>
          <el-button @click="copyDialogVisible = false">{{ t('cancel') }}</el-button>
          <el-button type="warning" @click="handleCopy" :loading="copying">{{ t('copyRegion') }}</el-button>
        </template>
      </el-dialog>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { request } from '@/utils/request'
import { REGIONS, BUYERS } from '@/constants'

const { t } = useI18n()

const user = JSON.parse(localStorage.getItem('user') || '{}')
const role = user.user_role || ''
const canEdit = computed(() => role === 'admin' || user.can_edit_products === true)

const loading = ref(false)
const submitting = ref(false)
const copying = ref(false)
const dialogVisible = ref(false)
const copyDialogVisible = ref(false)
const isEdit = ref(false)
const currentRegion = ref('') // 初始为空，强制用户选择
const tempRegion = ref('') // 选择弹窗临时值
const searchModelNo = ref('')
const list = ref<any[]>([])

const copyFrom = ref('')
const copyTo = ref('')

const form = ref({
  model_no: '',
  brand: '',
  commodities: '',
  descriptions: '',
  unit: '',
  export_ref_price: 0,
  import_ref_cost: 0,
})

// 地区名称国际化
/*
const regionKeys: Record<string, string> = {
  '玻利维亚': 'bolivia',
  '秘鲁': 'peru',
  '智利': 'chile',
  '西班牙': 'spain',
  '乌拉圭': 'uruguay',
  '委内瑞拉': 'venezuela',
  '阿根廷': 'argentina',
  '巴拉圭': 'paraguay',
  '巴拿马': 'panama',
  '日本': 'japan',
}
const regionLabel = (r: string) => t(regionKeys[r] || r)
*/
const regionLabel = (r: string) => {
  const target = REGIONS.find(item => item.value === r)
  return target ? t(target.labelKey) : r
}

const filteredList = computed(() => {
  if (!searchModelNo.value) return list.value
  return list.value.filter(p =>
    p.model_no.toLowerCase().includes(searchModelNo.value.toLowerCase())
  )
})
// 确认选择地区
const confirmRegion = () => {
  if (tempRegion.value) {
    currentRegion.value = tempRegion.value
    fetchProducts()
  }
}
// 打开复制弹窗，自动填入当前地区作为源地区
const openCopyDialog = () => {
  copyFrom.value = currentRegion.value
  copyTo.value = ''
  copyDialogVisible.value = true
}
const fetchProducts = async () => {
  if (!currentRegion.value) return
  loading.value = true
  try {
    const res = await request(`/api/products?region=${currentRegion.value}`)
    const data = await res.json()
    list.value = data
  } catch {
  } finally {
    loading.value = false
  }
}

const resetSearch = () => {
  searchModelNo.value = ''
  fetchProducts()
}

const openCreate = () => {
  isEdit.value = false
  form.value = {
    model_no: '', brand: '', commodities: '',
    descriptions: '', unit: '',
    export_ref_price: 0, import_ref_cost: 0,
  }
  dialogVisible.value = true
}

const openEdit = (row: any) => {
  isEdit.value = true
  form.value = {
    model_no: row.model_no,
    brand: row.brand || '',
    commodities: row.commodities || '',
    descriptions: row.descriptions || '',
    unit: row.unit || '',
    export_ref_price: row.export_ref_price ?? 0,
    import_ref_cost: row.import_ref_cost ?? 0,
  }
  dialogVisible.value = true
}

const submitForm = async () => {
  if (!form.value.model_no) {
    ElMessage.error('型号不能为空')
    return
  }
  submitting.value = true
  try {
    const payload = { ...form.value, region: currentRegion.value }
    const url = isEdit.value ? `/api/products/${form.value.model_no}` : '/api/products'
    const method = isEdit.value ? 'PUT' : 'POST'
    const res = await request(url, { method, body: JSON.stringify(payload) })
    if (res.ok) {
      ElNotification({ title: t('saveSuccess'), type: 'success', position: 'top-right', duration: 2500 })
      dialogVisible.value = false
      fetchProducts()
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
  // 先尝试删除，后端会检查引用
  try {
    const res = await request(`/api/products/${row.model_no}`, { method: 'DELETE' })
    if (res.status === 409) {
      const data = await res.json()
      // 有引用，弹出确认框
      await ElMessageBox.confirm(
        `该型号已在 ${data.ref_count} 条表单明细中使用，删除后表单数据保留但型号关联将断开，确认删除？`,
        t('warning'),
        { type: 'warning', confirmButtonText: '确认删除', cancelButtonText: t('cancel') }
      )
      // 用户确认后强制删除
      const forceRes = await request(`/api/products/${row.model_no}?force=true`, { method: 'DELETE' })
      if (forceRes.ok) {
        ElNotification({ title: t('deleteSuccess'), type: 'success', position: 'top-right', duration: 2500 })
        fetchProducts()
      }
    } else if (res.ok) {
      ElNotification({ title: t('deleteSuccess'), type: 'success', position: 'top-right', duration: 2500 })
      fetchProducts()
    }
  } catch { }
}

const handleCopy = async () => {
  if (!copyFrom.value || !copyTo.value) {
    ElMessage.error('请选择源地区和目标地区')
    return
  }
  if (copyFrom.value === copyTo.value) {
    ElMessage.warning('源地区和目标地区不能相同')
    return
  }
  copying.value = true
  try {
    const res = await request('/api/products/copy-prices', {
      method: 'POST',
      body: JSON.stringify({ from_region: copyFrom.value, to_region: copyTo.value })
    })
    if (res.ok) {
      const data = await res.json()
      ElNotification({
        title: `复制成功，共复制 ${data.copied} 条记录`,
        type: 'success',
        position: 'top-right',
        duration: 3000,
      })
      copyDialogVisible.value = false
      // 如果复制的目标就是当前地区，则刷新列表
      if (copyTo.value === currentRegion.value) {
        fetchProducts()
      }
    }
  } catch {
  } finally {
    copying.value = false
  }
}

onMounted(() => {
  // 不再自动加载，等用户选完地区再加载
})
</script>

<style scoped>
.region-select-mask {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 60vh;
}
.region-select-card {
  width: 380px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
}
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
