<template>
  <div>
    <!-- 顶部操作栏 -->
    <div class="action-bar">
      <el-button @click="goBack">← {{ t('backToList') }}</el-button>
      <div class="action-right">
        <el-button type="primary" @click="saveInvoice" v-if="canEdit || canReview" :loading="saving">
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
              <!-- 新建/复制时可填 inv_no，否则只读 -->
              <el-input v-model="invoice.inv_no" :disabled="!isNew" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item :label="t('region')">
              <!-- 地区始终只读，由出口表单复制而来 -->
              <el-input v-model="invoice.region" disabled />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item :label="t('invoiceDate')">
              <el-date-picker v-model="invoice.invoice_date" type="date" value-format="YYYY-MM-DD"
                style="width:100%" />
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
    <el-card class="section-card">
      <template #header>{{ t('reviewInfo') }}</template>
      <el-form label-width="140px">
        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item :label="t('reviewStatus')">
              <el-switch v-model="invoice.reviewed" :disabled="!canReview"
                :active-text="t('reviewed')" :inactive-text="t('pending')" />
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

      <el-table :data="invoice.items" border style="width:100%"
        :row-class-name="getRowClass">
        <el-table-column label="#" type="index" width="50" fixed="left" />

        <!-- 出口共用字段（只读，从出口表单复制来的） -->
        <el-table-column :label="t('brand')" prop="brand" min-width="90" fixed="left">
          <template #default="{ row }">
            <el-input v-model="row.brand" size="small" :disabled="!canEdit" />
          </template>
        </el-table-column>
        <el-table-column :label="t('commodities')" prop="commodities" min-width="140">
          <template #default="{ row }">
            <el-input v-model="row.commodities" size="small" :disabled="!canEdit" />
          </template>
        </el-table-column>
        <el-table-column :label="t('modelNo')" prop="model_no" min-width="110">
          <template #default="{ row }">
            <el-input v-model="row.model_no" size="small" :disabled="!canEdit" />
          </template>
        </el-table-column>
        <el-table-column :label="t('descriptions')" prop="descriptions" min-width="140">
          <template #default="{ row }">
            <el-input v-model="row.descriptions" size="small" :disabled="!canEdit" />
          </template>
        </el-table-column>
        <el-table-column :label="t('cartonQty')" prop="carton_qty" width="80">
          <template #default="{ row }">
            <el-input-number v-model="row.carton_qty" size="small" :min="0" style="width:100%"
              :disabled="!canEdit" :controls="false" />
          </template>
        </el-table-column>
        <el-table-column :label="t('unitQty')" prop="unit_qty" width="80">
          <template #default="{ row }">
            <el-input-number v-model="row.unit_qty" size="small" :min="0" style="width:100%"
              :disabled="!canEdit" :controls="false" />
          </template>
        </el-table-column>
        <el-table-column :label="t('unit')" prop="unit" width="70">
          <template #default="{ row }">
            <el-input v-model="row.unit" size="small" :disabled="!canEdit" />
          </template>
        </el-table-column>
        <el-table-column :label="t('exportUnitPrice')" prop="export_unit_price" width="100">
          <template #default="{ row }">
            <el-input-number v-model="row.export_unit_price" size="small" :precision="4"
              style="width:100%" :disabled="!canEdit" :controls="false" />
          </template>
        </el-table-column>
        <el-table-column :label="t('totalAmount')" prop="total_amount" width="100">
          <template #default="{ row }">
            <el-input-number v-model="row.total_amount" size="small" :precision="2"
              style="width:100%" :disabled="!canEdit" :controls="false" />
          </template>
        </el-table-column>
        <el-table-column :label="t('grossWeight')" prop="gross_weight" width="85">
          <template #default="{ row }">
            <el-input-number v-model="row.gross_weight" size="small" :precision="2"
              style="width:100%" :disabled="!canEdit" :controls="false" />
          </template>
        </el-table-column>
        <el-table-column :label="t('netWeight')" prop="net_weight" width="85">
          <template #default="{ row }">
            <el-input-number v-model="row.net_weight" size="small" :precision="2"
              style="width:100%" :disabled="!canEdit" :controls="false" />
          </template>
        </el-table-column>
        <el-table-column :label="t('volume')" prop="volume" width="80">
          <template #default="{ row }">
            <el-input-number v-model="row.volume" size="small" :precision="2"
              style="width:100%" :disabled="!canEdit" :controls="false" />
          </template>
        </el-table-column>

        <!-- 进口特有字段（可编辑） -->
        <el-table-column :label="t('exchangeRate')" prop="exchange_rate" width="95">
          <template #default="{ row }">
            <el-input-number v-model="row.exchange_rate" size="small" :min="0" :precision="4"
              style="width:100%" :disabled="!canEdit" :controls="false" />
          </template>
        </el-table-column>
        <el-table-column :label="t('seaFreight')" prop="sea_freight" width="95">
          <template #default="{ row }">
            <el-input-number v-model="row.sea_freight" size="small" :min="0" :precision="4"
              style="width:100%" :disabled="!canEdit" :controls="false" />
          </template>
        </el-table-column>
        <el-table-column :label="t('cifPrice')" prop="cif_price" width="95">
          <template #default="{ row }">
            <el-input-number v-model="row.cif_price" size="small" :min="0" :precision="4"
              style="width:100%" :disabled="!canEdit" :controls="false" />
          </template>
        </el-table-column>
        <el-table-column :label="t('cifTotal')" prop="cif_total" width="95">
          <template #default="{ row }">
            <el-input-number v-model="row.cif_total" size="small" :min="0" :precision="2"
              style="width:100%" :disabled="!canEdit" :controls="false" />
          </template>
        </el-table-column>
        <el-table-column :label="t('cifBsTotal')" prop="cif_bs_total" width="95">
          <template #default="{ row }">
            <el-input-number v-model="row.cif_bs_total" size="small" :min="0" :precision="2"
              style="width:100%" :disabled="!canEdit" :controls="false" />
          </template>
        </el-table-column>
        <el-table-column :label="t('importDutyVat')" prop="import_duty_vat" width="95">
          <template #default="{ row }">
            <el-input-number v-model="row.import_duty_vat" size="small" :min="0" :precision="2"
              style="width:100%" :disabled="!canEdit" :controls="false" />
          </template>
        </el-table-column>
        <el-table-column :label="t('transportation')" prop="transportation" width="95">
          <template #default="{ row }">
            <el-input-number v-model="row.transportation" size="small" :min="0" :precision="2"
              style="width:100%" :disabled="!canEdit" :controls="false" />
          </template>
        </el-table-column>
        <el-table-column :label="t('othersCharge')" prop="others_charge" width="95">
          <template #default="{ row }">
            <el-input-number v-model="row.others_charge" size="small" :min="0" :precision="2"
              style="width:100%" :disabled="!canEdit" :controls="false" />
          </template>
        </el-table-column>
        <el-table-column :label="t('totalCost')" prop="total_cost" width="95">
          <template #default="{ row }">
            <el-input-number v-model="row.total_cost" size="small" :min="0" :precision="2"
              style="width:100%" :disabled="!canEdit" :controls="false" />
          </template>
        </el-table-column>
        <el-table-column :label="t('importUnitCost')" prop="import_unit_cost" width="105">
          <template #default="{ row }">
            <el-input-number v-model="row.import_unit_cost" size="small" :min="0" :precision="4"
              style="width:100%" :disabled="!canEdit" :controls="false" />
          </template>
        </el-table-column>
        <el-table-column :label="t('systemPrice')" prop="system_price" width="95">
          <template #default="{ row }">
            <el-input-number v-model="row.system_price" size="small" :min="0" :precision="4"
              style="width:100%" :disabled="!canEdit" :controls="false" />
          </template>
        </el-table-column>
        <el-table-column :label="t('shareRate')" prop="share_rate" width="85">
          <template #default="{ row }">
            <el-input-number v-model="row.share_rate" size="small" :min="0" :precision="4"
              style="width:100%" :disabled="!canEdit" :controls="false" />
          </template>
        </el-table-column>
        <el-table-column :label="t('priceAlert')" prop="price_alert" width="80" fixed="right">
          <template #default="{ row }">
            <el-tag v-if="row.price_alert === 'yellow'" type="warning" size="small">⚠ {{ t('alertYellow') }}</el-tag>
            <el-tag v-else-if="row.price_alert === 'red'" type="danger" size="small">🔴 {{ t('alertRed') }}</el-tag>
            <span v-else>—</span>
          </template>
        </el-table-column>
        <el-table-column v-if="canEdit" :label="t('deleteItem')" width="70" fixed="right">
          <template #default="{ $index }">
            <el-button type="danger" size="small" @click="removeItem($index)">
              {{ t('deleteItem') }}
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
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
//import { ElMessage, ElMessageBox, ElNotification } from 'element-plus' //会与unplugin-auto-import存在重复导入的冲突，导致提示框无法正确渲染，所以注释掉
import { useI18n } from 'vue-i18n'
import { request } from '@/utils/request'
import { useCtrlS } from '@/composables/useCtrlS'

const { t } = useI18n()
const router = useRouter()
const route = useRoute()

const user = JSON.parse(localStorage.getItem('user') || '{}')
const role = user.user_role || ''

const isNew = computed(() => route.params.inv_no === 'new')
const copyFrom = route.query.copy as string | undefined
const isAdmin = computed(() => role === 'admin')
// 进口录入员和管理员可以编辑进口特有字段
const canEdit = computed(() => ['admin', 'import_input'].includes(role))
// 进口审核员和管理员可以操作审核字段
const canReview = computed(() => ['admin', 'import_review'].includes(role))

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

const fobTotal = computed(() =>
  invoice.value.items.reduce((sum, item) => sum + (item.total_amount || 0), 0)
)

// 行样式：根据 price_alert 高亮整行
const getRowClass = ({ row }: { row: any }) => {
  if (row.price_alert === 'red') return 'row-alert-red'
  if (row.price_alert === 'yellow') return 'row-alert-yellow'
  return ''
}

const addItem = () => {
  invoice.value.items.push({
    inv_no: invoice.value.inv_no,
    invoice_type: 'import',
    item_no: null,
    brand: '', commodities: '', model_no: '', descriptions: '',
    carton_qty: 0, unit_qty: 0, unit: '',
    export_unit_price: 0, total_amount: 0,
    gross_weight: 0, net_weight: 0, volume: 0,
    share_rate: 0, sea_freight: 0, cif_price: 0,
    cif_total: 0, cif_bs_total: 0, import_duty_vat: 0,
    transportation: 0, others_charge: 0, total_cost: 0,
    import_unit_cost: 0, system_price: 0, exchange_rate: 0,
    price_alert: 'normal',
  })
}

const removeItem = (index: number) => {
  invoice.value.items.splice(index, 1)
}

const fetchDetail = async () => {
  const invNo = route.params.inv_no as string
  try {
    const res = await request(`/api/invoices/import/${invNo}`)
    const data = await res.json()
    invoice.value = data
  } catch {
    // request 已统一处理错误提示和 401 跳转，403 会提示权限不足
  }
}

const saveInvoice = async () => {
  saving.value = true
  //console.log('saveInvoice called, will use ElNotification')
  const payload = {
    ...invoice.value,
    invoice_date: invoice.value.invoice_date
      ? invoice.value.invoice_date.slice(0, 10)
      : null
  }
  try {
    const method = isNew.value ? 'POST' : 'PUT'
    const url = isNew.value
      ? '/api/invoices/import'
      : `/api/invoices/import/${route.params.inv_no}`

    const res = await request(url, {
      method,
      body: JSON.stringify(payload)
    })
    const data = await res.json()
    if (res.ok) {
      // 保存后重新拉取，price_alert 由后端计算后返回
      //console.log('about to show notification')
      ElNotification({
      	//title: '保存修改语句测试',
        title: t('saveSuccess'),
        type: 'success',
        position: 'top-right',
        duration: 2500,
        offset: 20,
        customClass: 'notification-large',
      })
      if (isNew.value) {
        router.replace(`/import/${invoice.value.inv_no}`)
      } else {
        await fetchDetail()
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
    const res = await request(`/api/invoices/import/${route.params.inv_no}`, {
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
      router.push('/import')
    }
  } catch {
    // request 已统一处理错误提示
  }
}

const goBack = () => {
  router.push('/import')
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
      const res = await request(`/api/invoices/import/${copyFrom}`)
      const data = await res.json()
      invoice.value = { ...data, inv_no: '' }
    } catch { }
  } else if (!isNew.value) {
    fetchDetail()
  }
})
</script>

<style scoped>
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

/* 价格预警行高亮 */
:deep(.row-alert-yellow) {
  background-color: #fdf6ec !important;
}

:deep(.row-alert-red) {
  background-color: #fef0f0 !important;
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
