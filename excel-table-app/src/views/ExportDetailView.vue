<template>
  <div>
    <!-- 顶部操作栏 -->
    <div class="action-bar">
      <el-button @click="goBack">← {{ t('backToList') }}</el-button>
      <div class="action-right">
        <el-button type="primary" @click="saveInvoice" v-if="canEdit" :loading="saving">
          {{ t('save') }}
        </el-button>
        <el-button type="danger" @click="deleteInvoice" v-if="isAdmin">
          {{ t('delete') }}
        </el-button>
      </div>
    </div>

    <!-- 表单头部信息 -->
    <el-card class="section-card">
      <template #header>{{ t('invoiceInfo') }}</template>
      <el-form :model="invoice" label-width="140px" :disabled="!canEdit">
        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item :label="t('invNo')">
              <el-input v-model="invoice.inv_no" :disabled="!isNew" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item :label="t('region')">
              <el-select v-model="invoice.region" :placeholder="t('region')" style="width:100%">
                <el-option :label="t('bolivia')" value="玻利维亚" />
                <el-option :label="t('peru')" value="秘鲁" />
                <el-option :label="t('chile')" value="智利" />
                <el-option :label="t('spain')" value="西班牙" />
                <el-option :label="t('usa')" value="美国" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item :label="t('invoiceDate')">
              <el-date-picker v-model="invoice.invoice_date" type="date" value-format="YYYY-MM-DD" style="width:100%"
                @change="onDateChange" @visible-change="(v: any) => console.log('[日期排查] 面板开关:', v)" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item :label="t('contractNo')">
              <el-input v-model="invoice.contract_no" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item :label="t('shippingLine')">
              <el-input v-model="invoice.shipping_line" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item :label="t('blNo')">
              <el-input v-model="invoice.bl_no" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item :label="t('containerNo')">
              <el-input v-model="invoice.container_no" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item :label="t('portOfLoading')">
              <el-input v-model="invoice.port_of_loading" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item :label="t('portOfDischarge')">
              <el-input v-model="invoice.port_of_discharge" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item :label="t('finalDestination')">
              <el-input v-model="invoice.final_destination" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item :label="t('buyerName')">
              <el-input v-model="invoice.buyer_name" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item :label="t('buyerTel')">
              <el-input v-model="invoice.buyer_tel" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item :label="t('sellerName')">
              <el-input v-model="invoice.seller_name" />
            </el-form-item>
          </el-col>
          <el-col :span="16">
            <el-form-item :label="t('buyerAddress')">
              <el-input v-model="invoice.buyer_address" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item :label="t('paymentTerm')">
              <el-input v-model="invoice.payment_term" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
    </el-card>

    <!-- 审核区域 -->
    <el-card class="section-card" v-if="!isNew">
      <template #header>{{ t('reviewInfo') }}</template>
      <el-form label-width="140px">
        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item :label="t('reviewStatus')">
              <el-switch v-model="invoice.reviewed" :disabled="!canReview" :active-text="t('reviewed')"
                :inactive-text="t('pending')" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item :label="t('reviewer')">
              <el-input :value="invoice.reviewer" disabled />
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item :label="t('reviewComment')">
              <el-input v-model="invoice.review_comment" type="textarea" :rows="2" :disabled="!canReview" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
    </el-card>

    <!-- 明细行 -->
    <el-card class="section-card">
      <template #header>
        <div class="card-header">
          <span>{{ t('itemDetails') }}</span>
          <div style="display:flex;gap:8px;align-items:center">
            <span class="legend">
              <span class="legend-dot yellow"></span>10%~50%
              <span class="legend-dot red" style="margin-left:8px"></span>&gt;50%
            </span>
            <el-button v-if="canEdit" type="primary" size="small" @click="addItem">
              {{ t('addItem') }}
            </el-button>
          </div>
        </div>
      </template>
      <el-table :data="invoice.items" border style="width:100%" :row-class-name="getRowClass">
        <el-table-column label="#" type="index" width="60" />
        <el-table-column :label="t('brand')" prop="brand" min-width="100">
          <template #default="{ row }">
            <el-input v-model="row.brand" size="small" :disabled="!canEdit" />
          </template>
        </el-table-column>
        <el-table-column :label="t('commodities')" prop="commodities" min-width="160">
          <template #default="{ row }">
            <el-input v-model="row.commodities" size="small" :disabled="!canEdit" />
          </template>
        </el-table-column>
        <el-table-column :label="t('modelNo')" prop="model_no" min-width="120">
          <template #default="{ row }">
            <el-input v-model="row.model_no" size="small" :disabled="!canEdit" />
          </template>
        </el-table-column>
        <el-table-column :label="t('descriptions')" prop="descriptions" min-width="160">
          <template #default="{ row }">
            <el-input v-model="row.descriptions" size="small" :disabled="!canEdit" />
          </template>
        </el-table-column>
        <!-- 箱数 -->
        <el-table-column :label="t('cartonQty')" prop="carton_qty" width="90">
          <template #default="{ row, $index }">
            <el-input-number v-model="row.carton_qty" size="small" :min="0" style="width:100%" :disabled="!canEdit"
              :controls="false" @change="onQtyChange(row, $index)" />
          </template>
        </el-table-column>
        <!-- 每箱货物数（新增） -->
        <el-table-column :label="t('qtyPerCarton')" prop="qty_per_carton" width="100">
          <template #default="{ row, $index }">
            <el-input-number v-model="row.qty_per_carton" size="small" :min="0" style="width:100%" :disabled="!canEdit"
              :controls="false" @change="onQtyChange(row, $index)" />
          </template>
        </el-table-column>
        <!-- 总货物量（受手动覆盖控制） -->
        <el-table-column :label="t('unitQty')" prop="unit_qty" width="100">
          <template #default="{ row, $index }">
            <el-input-number v-model="row.unit_qty" size="small" :min="0" style="width:100%"
              :disabled="!canEdit || !manualOverride[$index]" :controls="false" />
          </template>
        </el-table-column>
        <el-table-column :label="t('unit')" prop="unit" width="80">
          <template #default="{ row }">
            <el-input v-model="row.unit" size="small" :disabled="!canEdit" />
          </template>
        </el-table-column>
        <!-- 单价 -->
        <el-table-column :label="t('unitPrice')" prop="export_unit_price" width="100">
          <template #default="{ row, $index }">
            <el-input-number v-model="row.export_unit_price" size="small" :min="0" :precision="3" style="width:100%"
              :disabled="!canEdit" :controls="false" @change="onQtyChange(row, $index)" />
          </template>
        </el-table-column>
        <!-- 总金额（受手动覆盖控制） -->
        <el-table-column :label="t('totalAmount')" prop="total_amount" width="100">
          <template #default="{ row, $index }">
            <el-input-number v-model="row.total_amount" size="small" :min="0" :precision="2" style="width:100%"
              :disabled="!canEdit || !manualOverride[$index]" :controls="false" />
          </template>
        </el-table-column>
        <el-table-column :label="t('grossWeight')" prop="gross_weight" width="90">
          <template #default="{ row }">
            <el-input-number v-model="row.gross_weight" size="small" :min="0" :precision="2" style="width:100%"
              :disabled="!canEdit" :controls="false" />
          </template>
        </el-table-column>
        <el-table-column :label="t('netWeight')" prop="net_weight" width="90">
          <template #default="{ row }">
            <el-input-number v-model="row.net_weight" size="small" :min="0" :precision="2" style="width:100%"
              :disabled="!canEdit" :controls="false" />
          </template>
        </el-table-column>
        <el-table-column :label="t('volume')" prop="volume" width="90">
          <template #default="{ row }">
            <el-input-number v-model="row.volume" size="small" :min="0" :precision="2" style="width:100%"
              :disabled="!canEdit" :controls="false" />
          </template>
        </el-table-column>
        <el-table-column :label="t('priceAlert')" prop="price_alert" width="80" fixed="right">
          <template #default="{ row }">
            <el-tag v-if="row.price_alert === 'yellow'" type="warning" size="small">⚠ {{ t('alertYellow') }}</el-tag>
            <el-tag v-else-if="row.price_alert === 'red'" type="danger" size="small">🔴 {{ t('alertRed') }}</el-tag>
            <span v-else>—</span>
          </template>
        </el-table-column>
        <el-table-column :label="t('operation')" fixed="right" width="160" v-if="canEdit">
          <template #default="{ $index }">
            <el-button size="small" :type="manualOverride[$index] ? 'warning' : 'primary'" link
              @click="toggleOverride($index)">
              {{ manualOverride[$index] ? t('restoreAuto') : t('manualEdit') }}
            </el-button>
            <el-button size="small" type="danger" link @click="removeItem($index)">
              {{ t('delete') }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- FOB 合计 -->
      <div class="fob-total">
        FOB TOTAL: <strong>USD {{ fobTotal.toFixed(2) }}</strong>
      </div>
    </el-card>

    <!-- 底部操作栏 -->
    <div class="action-bar" v-if="canEdit || canReview">
      <span></span>
      <el-button type="primary" @click="saveInvoice" :loading="saving">
        {{ t('save') }}
      </el-button>
    </div>

  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
//import { ElMessage, ElMessageBox, ElNotification } from 'element-plus'//会与unplugin-auto-import存在重复导入的冲突，导致提示框无法正确渲染，所以注释掉
import { useI18n } from 'vue-i18n'
import { request } from '@/utils/request'
import { useCtrlS } from '@/composables/useCtrlS'

const { t } = useI18n()
const router = useRouter()
const route = useRoute()

const user = JSON.parse(localStorage.getItem('user') || '{}')
const role = user.user_role || ''

const isNew = computed(() => route.params.inv_no === 'new')
const copyFrom = route.query.copy as string | undefined  // 复制模式：源 inv_no
const isAdmin = computed(() => role === 'admin')
const canEdit = computed(() => ['admin', 'export_input'].includes(role))
const canReview = computed(() => ['admin', 'export_review'].includes(role))

const saving = ref(false)

const invoice = ref({
  inv_no: '',
  region: '',
  invoice_date: '',
  contract_no: '',
  shipping_line: '',
  bl_no: '',
  container_no: '',
  port_of_loading: '',
  port_of_discharge: '',
  final_destination: '',
  buyer_name: '',
  buyer_address: '',
  buyer_tel: '',
  seller_name: '',
  payment_term: '',
  fob_total: 0,
  reviewed: false,
  reviewer: '',
  review_comment: '',
  items: [] as any[],
})

function onDateChange(val: string | null) {
  console.log('[日期排查] date-picker @change 收到:', val, '类型:', typeof val)
  console.log('[日期排查] 当前 invoice.invoice_date:', invoice.value.invoice_date)
}

// 记录哪些行开启了手动覆盖（key 为行索引）
const manualOverride = ref<Record<number, boolean>>({})

const toggleOverride = (index: number) => {
  const next = !manualOverride.value[index]
  manualOverride.value[index] = next

  // 从手动切回自动时，立刻重新计算一次
  if (!next) {
    const row = invoice.value.items[index]
    if (row) calcItem(row)
  }
}

const calcItem = (row: any) => {
  const carton = Number(row.carton_qty) || 0
  const perCarton = Number(row.qty_per_carton) || 0
  const price = Number(row.export_unit_price) || 0

  row.unit_qty = carton * perCarton
  row.total_amount = +(row.unit_qty * price).toFixed(2)
}

const onQtyChange = (row: any, index: number) => {
  // 只有未开启手动覆盖时才自动算
  if (!manualOverride.value[index]) {
    calcItem(row)
  }
}



const fobTotal = computed(() =>
  invoice.value.items.reduce((sum, item) => sum + (item.total_amount || 0), 0)
)

const fetchDetail = async () => {
  const invNo = route.params.inv_no as string
  //原来的fetch请求改为用request
  //const res = await fetch(`/api/invoices/export/${invNo}`, {
  //  headers: { Authorization: `Bearer ${token}` }
  //})
  //if (res.ok) {
  //  const data = await res.json()
  //  invoice.value = data
  try {
    const res = await request(`/api/invoices/export/${invNo}`)
    const data = await res.json()
    invoice.value = data
  } catch {
    // request 已统一处理错误提示和 401 跳转
  }
}

const getRowClass = ({ row }: { row: any }) => {
  if (row.price_alert === 'red') return 'row-alert-red'
  if (row.price_alert === 'yellow') return 'row-alert-yellow'
  return ''
}

const addItem = () => {
  invoice.value.items.push({
    brand: '', commodities: '', model_no: '', descriptions: '',
    carton_qty: 0, qty_per_carton: 0, unit_qty: 0, unit: '',
    export_unit_price: 0, total_amount: 0,
    gross_weight: 0, net_weight: 0, volume: 0,
  })
}

const removeItem = (index: number) => {
  invoice.value.items.splice(index, 1)

  // 同步调整 manualOverride 的索引
  const newMap: Record<number, boolean> = {}
  Object.keys(manualOverride.value).forEach((k) => {
    const i = Number(k)
    const val = manualOverride.value[i]
    if (val === undefined) return // 跳过不存在的 key
    if (i < index) newMap[i] = val
    else if (i > index) newMap[i - 1] = val
  })
  manualOverride.value = newMap
}


const saveInvoice = async () => {
  saving.value = true
  const payload = {
    ...invoice.value,
    invoice_date: invoice.value.invoice_date
      ? invoice.value.invoice_date.slice(0, 10)
      : null
  }
  try {
    const method = isNew.value ? 'POST' : 'PUT'
    const url = isNew.value
      ? '/api/invoices/export'
      : `/api/invoices/export/${route.params.inv_no}`

    const res = await request(url, {
      method,
      body: JSON.stringify(payload)
    })
    const data = await res.json()
    if (res.ok) {
      ElNotification({
        title: t('saveSuccess'),
        type: 'success',
        position: 'top-right',
        duration: 2500,
        offset: 20,
        customClass: 'notification-large',
      })
      if (isNew.value) {
        router.replace(`/export/${invoice.value.inv_no}`)
      }
    } else {
      ElMessage.error(data.error || t('saveFailed'))
    }
  } finally {
    saving.value = false
  }
}

const deleteInvoice = async () => {
  await ElMessageBox.confirm(t('deleteConfirm'), t('warning'), { type: 'warning' })
  try {
    const res = await request(`/api/invoices/export/${route.params.inv_no}`, {
      method: 'DELETE',
    })
    if (res.ok) {
      ElNotification({
        title: t('deleteSuccess'),
        type: 'success',
        position: 'top-right',
        duration: 2500,
        offset: 20,
        customClass: 'notification-large',
      })
      router.push('/export')
    }
  } catch {
    // request 已统一处理错误提示
  }
}

const goBack = () => {
  router.push('/export')
}

// Ctrl+S 快捷键
useCtrlS(() => {
  if (canEdit.value || canReview.value) {
    saveInvoice()
  }
})

onMounted(async () => {
  if (copyFrom) {
    // 复制模式：拉取源表单数据，清空 inv_no
    try {
      const res = await request(`/api/invoices/export/${copyFrom}`)
      const data = await res.json()
      invoice.value = { ...data, inv_no: '' }
    } catch { }
  } else if (!isNew.value) {
    fetchDetail()
  }
})
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.legend {
  font-size: 12px;
  color: #606266;
  display: flex;
  align-items: center;
  gap: 4px;
}

.legend-dot {
  display: inline-block;
  width: 12px;
  height: 12px;
  border-radius: 50%;
}

.legend-dot.yellow {
  background-color: #e6a23c;
}

.legend-dot.red {
  background-color: #f56c6c;
}

:deep(.row-alert-yellow) {
  background-color: #fdf6ec !important;
}

:deep(.row-alert-red) {
  background-color: #fef0f0 !important;
}

.action-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.action-right {
  display: flex;
  gap: 8px;
}

.section-card {
  margin-bottom: 16px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.fob-total {
  text-align: right;
  margin-top: 12px;
  font-size: 15px;
  color: #303133;
}
</style>

<style>
/* ElNotification 放大样式，需要全局生效（不加 scoped） */
.notification-large {
  width: 340px !important;
  font-size: 16px !important;
}

.notification-large .el-notification__title {
  font-size: 18px !important;
  font-weight: 600 !important;
}
</style>