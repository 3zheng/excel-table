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
              <el-date-picker v-model="invoice.invoice_date" type="date" value-format="YYYY-MM-DD" style="width:100%" />
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
        <el-table-column label="#" type="index" width="50" fixed="left" />

        <!-- 出口共用字段 -->
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
        <el-table-column :label="t('cartonQty')" prop="carton_qty" width="90">
          <template #default="{ row, $index }">
            <el-input-number v-model="row.carton_qty" size="small" :min="0" :value-on-clear="undefined"
              style="width:100%" :disabled="!canEdit" :controls="false" @change="onQtyChange(row, $index)" />
          </template>
        </el-table-column>
        <!-- 每箱货物数（新增） -->
        <el-table-column :label="t('qtyPerCarton')" prop="qty_per_carton" width="100">
          <template #default="{ row, $index }">
            <el-input-number v-model="row.qty_per_carton" size="small" :min="0" :value-on-clear="undefined"
              style="width:100%" :disabled="!canEdit" :controls="false" @change="onQtyChange(row, $index)" />
          </template>
        </el-table-column>
        <!-- 总货物量（受手动覆盖控制） -->
        <el-table-column :label="t('unitQty')" prop="unit_qty" width="100">
          <template #default="{ row, $index }">
            <el-input-number v-model="row.unit_qty" size="small" :min="0" :value-on-clear="undefined" style="width:100%"
              :disabled="!canEdit || !manualOverride[$index]" :controls="false" />
          </template>
        </el-table-column>
        <el-table-column :label="t('unit')" prop="unit" width="70">
          <template #default="{ row }">
            <el-input v-model="row.unit" size="small" :disabled="!canEdit" />
          </template>
        </el-table-column>
        <el-table-column :label="t('exportUnitPrice')" prop="export_unit_price" width="100">
          <template #default="{ row }">
            <el-input-number v-model="row.export_unit_price" size="small" :precision="4" :value-on-clear="undefined"
              style="width:100%" :disabled="!canEdit" :controls="false" />
          </template>
        </el-table-column>
        <el-table-column :label="t('totalAmount')" prop="total_amount" width="100">
          <template #default="{ row }">
            <el-input-number v-model="row.total_amount" size="small" :precision="2" :value-on-clear="undefined"
              style="width:100%" :disabled="!canEdit" :controls="false" />
          </template>
        </el-table-column>
        <el-table-column :label="t('grossWeight')" prop="gross_weight" width="85">
          <template #default="{ row }">
            <el-input-number v-model="row.gross_weight" size="small" :precision="2" :value-on-clear="undefined"
              style="width:100%" :disabled="!canEdit" :controls="false" />
          </template>
        </el-table-column>
        <el-table-column :label="t('netWeight')" prop="net_weight" width="85">
          <template #default="{ row }">
            <el-input-number v-model="row.net_weight" size="small" :precision="2" :value-on-clear="undefined"
              style="width:100%" :disabled="!canEdit" :controls="false" />
          </template>
        </el-table-column>
        <el-table-column :label="t('volume')" prop="volume" width="80">
          <template #default="{ row }">
            <el-input-number v-model="row.volume" size="small" :precision="2" :value-on-clear="undefined"
              style="width:100%" :disabled="!canEdit" :controls="false" />
          </template>
        </el-table-column>

        <!-- 进口特有字段（可编辑） -->
        <el-table-column :label="t('exchangeRate')" prop="exchange_rate" width="95">
          <template #default="{ row, $index }">
            <el-input-number v-model="row.exchange_rate" size="small" :min="0" :precision="4"
              :value-on-clear="undefined" style="width:100%" :disabled="!canEdit" :controls="false"
              @change="onCostChange(row, $index)" />
          </template>
        </el-table-column>

        <el-table-column :label="t('seaFreight')" prop="sea_freight" width="95">
          <template #default="{ row, $index }">
            <el-input-number v-model="row.sea_freight" size="small" :min="0" :precision="2" :value-on-clear="undefined"
              style="width:100%" :disabled="!canEdit" :controls="false" @change="onCostChange(row, $index)" />
          </template>
        </el-table-column>
        <el-table-column :label="t('unitQty')" prop="unit_qty" width="100">
          <template #default="{ row, $index }">
            <el-input-number v-model="row.unit_qty" size="small" :min="0" :value-on-clear="undefined" style="width:100%"
              :disabled="!canEdit || !manualOverride[$index]" :controls="false" />
          </template>
        </el-table-column>
        <el-table-column :label="t('cifTotal')" prop="cif_total" width="95">
          <template #default="{ row, $index }">
            <el-input-number v-model="row.cif_total" size="small" :min="0" :precision="2" :value-on-clear="undefined"
              style="width:100%" :disabled="!canEdit || !manualOverride[$index]" :controls="false" />
          </template>
        </el-table-column>
        <el-table-column :label="t('cifPrice')" prop="cif_price" width="95">
          <template #default="{ row, $index }">
            <el-input-number v-model="row.cif_price" size="small" :min="0" :precision="3" :value-on-clear="undefined"
              style="width:100%" :disabled="!canEdit || !manualOverride[$index]" :controls="false" />
          </template>
        </el-table-column>
        <el-table-column :label="t('cifBsTotal')" prop="cif_bs_total" width="95">
          <template #default="{ row, $index }">
            <el-input-number v-model="row.cif_bs_total" size="small" :min="0" :precision="2" :value-on-clear="undefined"
              style="width:100%" :disabled="!canEdit || !manualOverride[$index]" :controls="false" />
          </template>
        </el-table-column>
        <!-- 拆分后的关税列 -->
        <el-table-column :label="t('importDuty')" prop="import_duty" width="95">
          <template #default="{ row, $index }">
            <el-input-number v-model="row.import_duty" size="small" :min="0" :precision="2" :value-on-clear="undefined"
              style="width:100%" :disabled="!canEdit" :controls="false" @change="onCostChange(row, $index)" />
          </template>
        </el-table-column>
        <!-- 拆分后的增值税列 -->
        <el-table-column :label="t('vat')" prop="vat" width="95">
          <template #default="{ row, $index }">
            <el-input-number v-model="row.vat" size="small" :min="0" :precision="2" :value-on-clear="undefined"
              style="width:100%" :disabled="!canEdit" :controls="false" @change="onCostChange(row, $index)" />
          </template>
        </el-table-column>
        <el-table-column :label="t('transportation')" prop="transportation" width="95">
          <template #default="{ row, $index }">
            <el-input-number v-model="row.transportation" size="small" :min="0" :precision="2"
              :value-on-clear="undefined" style="width:100%" :disabled="!canEdit" :controls="false"
              @change="onCostChange(row, $index)" />
          </template>
        </el-table-column>
        <el-table-column :label="t('othersCharge')" prop="others_charge" width="95">
          <template #default="{ row, $index }">
            <el-input-number v-model="row.others_charge" size="small" :min="0" :precision="2"
              :value-on-clear="undefined" style="width:100%" :disabled="!canEdit" :controls="false"
              @change="onCostChange(row, $index)" />
          </template>
        </el-table-column>
        <el-table-column :label="t('totalCost')" prop="total_cost" width="95">
          <template #default="{ row, $index }">
            <el-input-number v-model="row.total_cost" size="small" :min="0" :precision="2" :value-on-clear="undefined"
              style="width:100%" :disabled="!canEdit || !manualOverride[$index]" :controls="false" />
          </template>
        </el-table-column>
        <el-table-column :label="t('importUnitCost')" prop="import_unit_cost" width="105">
          <template #default="{ row, $index }">
            <el-input-number v-model="row.import_unit_cost" size="small" :min="0" :precision="2"
              :value-on-clear="undefined" style="width:100%" :disabled="!canEdit" :controls="false"
              @change="onCostChange(row, $index)" />
          </template>
        </el-table-column>
        <el-table-column :label="t('pieceArea')" prop="piece_area" width="95">
          <template #default="{ row, $index }">
            <el-input-number v-model="row.piece_area" size="small" :min="0" :value-on-clear="undefined"
              style="width:100%" :disabled="!canEdit" :controls="false" @change="onCostChange(row, $index)" />
          </template>
        </el-table-column>
        <el-table-column :label="t('systemPrice')" prop="system_price" width="95">
          <template #default="{ row, $index }">
            <el-input-number v-model="row.system_price" size="small" :min="0" :precision="2" :value-on-clear="undefined"
              style="width:100%" :disabled="!canEdit || !manualOverride[$index]" :controls="false" />
          </template>
        </el-table-column>
        <el-table-column :label="t('shareRate')" prop="share_rate" width="85">
          <template #default="{ row, $index }">
            <el-input-number v-model="row.share_rate" size="small" :min="0" :precision="4" :value-on-clear="undefined"
              :formatter="(v: number | null) => v != null ? `${(v * 100).toFixed(2)}%` : ''"
              :parser="(v: string) => v.replace('%', '') ? Number(v.replace('%', '')) / 100 : 0"
              style="width:100%" :disabled="!canEdit || !manualOverride[$index]" :controls="false" />
          </template>
        </el-table-column>
        <el-table-column :label="t('priceAlert')" prop="price_alert" width="80" fixed="right">
          <template #default="{ row }">
            <el-tag v-if="row.price_alert === 'yellow'" type="warning" size="small">⚠ {{ t('alertYellow') }}</el-tag>
            <el-tag v-else-if="row.price_alert === 'red'" type="danger" size="small">🔴 {{ t('alertRed') }}</el-tag>
            <span v-else>—</span>
          </template>
        </el-table-column>
        <el-table-column v-if="canEdit" :label="t('deleteItem')" width="150" fixed="right">
          <template #default="{ $index }">
            <el-button size="small" :type="manualOverride[$index] ? 'warning' : 'primary'" link
              @click="toggleOverride($index)">
              {{ manualOverride[$index] ? t('restoreAuto') : t('manualEdit') }}
            </el-button>
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
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
//import { ElMessage, ElMessageBox, ElNotification } from 'element-plus' //会与unplugin-auto-import存在重复导入的冲突，导致提示框无法正确渲染，所以注释掉
import { useI18n } from 'vue-i18n'
import { request } from '@/utils/request'
import { useCtrlS } from '@/composables/useCtrlS'
import { REGIONS, BUYERS } from '@/constants'

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

// 1. 单条进口明细行类型
export interface InvoiceItem {
  inv_no?: string
  invoice_type?: string
  item_no?: number | null

  // 基础与出口字段
  brand?: string
  commodities?: string
  model_no?: string
  descriptions?: string
  carton_qty?: number
  qty_per_carton?: number
  unit_qty?: number
  unit?: string
  export_unit_price?: number
  total_amount?: number
  gross_weight?: number
  net_weight?: number
  volume?: number

  // 进口特有与费用字段
  exchange_rate?: number
  sea_freight?: number
  cif_price?: number
  cif_total?: number
  cif_bs_total?: number
  import_duty?: number
  vat?: number
  transportation?: number
  others_charge?: number
  total_cost?: number
  import_unit_cost?: number
  piece_area?: number
  system_price?: number
  share_rate?: number

  // 状态标识
  price_alert?: 'normal' | 'yellow' | 'red' | string
}

// 2. 完整发票表单类型
export interface Invoice {
  inv_no: string
  region: string
  invoice_date: string
  contract_no: string
  shipping_line: string
  bl_no: string
  container_no: string
  port_of_loading: string
  port_of_discharge: string
  final_destination: string
  buyer_name: string
  buyer_address: string
  buyer_tel: string
  seller_name: string
  payment_term: string
  fob_total: number
  reviewed: boolean
  reviewer: string
  review_comment: string
  items: InvoiceItem[]
}

const invoice = ref<Invoice>({
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
  items: [],
})

const fobTotal = computed(() =>
  invoice.value.items.reduce((sum, item) => sum + (item.total_amount || 0), 0)
)

watch(fobTotal, (newVal) => {
  invoice.value.fob_total = newVal // 副作用：更新外部状态
})

// 行样式：根据 price_alert 高亮整行
const getRowClass = ({ row }: { row: any }) => {
  if (row.price_alert === 'red') return 'row-alert-red'
  if (row.price_alert === 'yellow') return 'row-alert-yellow'
  return ''
}

const onBuyerChange = (name: string) => {
  const buyer = BUYERS.find(b => b.name === name)
  if (buyer) {
    invoice.value.buyer_address = buyer.address
    invoice.value.buyer_tel = buyer.tel
  }
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

const calcItem = (row: InvoiceItem) => {
  const carton = row.carton_qty
  const perCarton = row.qty_per_carton
  const price = row.export_unit_price

  // 只有箱数和每箱数量都有值（且非 null/undefined）时才自动计算总数量
  if (carton !== undefined && carton !== null && perCarton !== undefined && perCarton !== null) {
    row.unit_qty = Number(carton) * Number(perCarton)
  } else {
    row.unit_qty = undefined
  }

  // 只有总数量和单价都有值时才计算总金额
  if (row.unit_qty !== undefined && row.unit_qty !== null && price !== undefined && price !== null) {
    row.total_amount = +(Number(row.unit_qty) * Number(price)).toFixed(2)
  } else {
    row.total_amount = undefined
  }
  /*
  CIF总额 cifTotal = 总金额+海运费
  CIF单价 cifPrice =（总金额+海运费）/ 总数量 
  CIF本位币总额 cifBsTotal = CIF总额 * 汇率
  总成本 totalCost = CIF本位币总额 + 关税 - 增值税 + 本地运费 + 其他费用
  进口单位成本 importUnitCost = 总成本 / 总数量
  分摊比例 shareRate = 单行的总金额 / 所有行累加的总金额
  如果单片面积有填值，那么
  系统价格systemPrice = 进口单位成本 * 单片面积
  否则
  系统价格systemPrice = 进口单位成本
  */
  const exchangeRate = row.exchange_rate
  const seaFreight = row.sea_freight
  const importDuty = row.import_duty
  const vat = row.vat
  const transportation = row.transportation
  const othersCharge = row.others_charge

  if (seaFreight !== undefined && seaFreight !== null && row.total_amount !== undefined && row.total_amount !== null) {
    row.cif_total = +(Number(seaFreight) + Number(row.total_amount)).toFixed(2)
  } else {
    row.cif_total = undefined
  }
  if (row.cif_total !== undefined && row.cif_total !== null && row.unit_qty !== undefined && row.unit_qty !== null) {
    row.cif_price = +(Number(row.cif_total) / Number(row.unit_qty)).toFixed(2)
  } else {
    row.cif_price = undefined
  }
  if (row.cif_total !== undefined && row.cif_total !== null && exchangeRate !== undefined && exchangeRate !== null) {
    row.cif_bs_total = +(Number(row.cif_total) * Number(exchangeRate)).toFixed(2)
  } else {
    row.cif_bs_total = undefined
  }
  if (row.cif_bs_total !== undefined && row.cif_bs_total !== null && importDuty !== undefined && importDuty !== null
    && vat !== undefined && vat !== null && transportation !== undefined && transportation !== null && othersCharge !== undefined && othersCharge !== null 
  ) {
    row.total_cost = +(Number(row.cif_bs_total) + Number(importDuty) + Number(importDuty) - Number(vat) + Number(transportation) + Number(othersCharge)).toFixed(2)
  } else {
    row.total_cost = undefined
  }
  if (row.total_cost !== undefined && row.total_cost !== null && row.unit_qty !== undefined && row.unit_qty !== null) {
    row.import_unit_cost = +(Number(row.total_cost) / Number(row.unit_qty)).toFixed(2)
  } else {
    row.import_unit_cost = undefined
  }
  if (row.total_amount !== undefined && row.total_amount !== null ) {
    row.share_rate = +(Number(row.total_amount) / Number(invoice.value.fob_total)).toFixed(4)
  } else {
    row.share_rate = undefined
  }
  if (row.import_unit_cost !== undefined && row.import_unit_cost !== null ) {
    if (row.piece_area !== undefined && row.piece_area !== null ){
      row.system_price = +(Number(row.import_unit_cost) * Number(row.piece_area)).toFixed(2)
    }else{
      row.system_price = row.import_unit_cost
    }    
  } else {
    row.system_price = undefined
  }
}

const onQtyChange = (row: InvoiceItem, index: number) => {
  // 只有未开启手动覆盖时才自动算
  if (!manualOverride.value[index]) {
    calcItem(row)
  }
}

const onCostChange = (row: InvoiceItem, index: number) => {
  if (!manualOverride.value[index]) {
    calcItem(row)
  }
}


const addItem = () => {
  invoice.value.items.push({
    inv_no: invoice.value.inv_no,
    invoice_type: 'import',
    item_no: null,
    brand: '', commodities: '', model_no: '', descriptions: '',
    carton_qty: undefined, unit_qty: undefined, unit: '',
    export_unit_price: undefined, total_amount: undefined,
    gross_weight: undefined, net_weight: undefined, volume: undefined,
    share_rate: undefined, sea_freight: undefined, cif_price: undefined,
    cif_total: undefined, cif_bs_total: undefined,
    import_duty: undefined, vat: undefined,// 替换 import_duty_vat: undefined
    transportation: undefined, others_charge: undefined, total_cost: undefined,
    import_unit_cost: undefined, piece_area: undefined, system_price: undefined,
    exchange_rate: undefined, price_alert: 'normal',
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
      : null,
    // 遍历明细行，把所有可能为 undefined 的数值字段统一兜底为 0
    items: invoice.value.items.map(item => ({
      ...item,
      // 出口/基础数值字段
      carton_qty: item.carton_qty ?? 0,
      unit_qty: item.unit_qty ?? 0,
      export_unit_price: item.export_unit_price ?? 0,
      total_amount: item.total_amount ?? 0,
      gross_weight: item.gross_weight ?? 0,
      net_weight: item.net_weight ?? 0,
      volume: item.volume ?? 0,

      // 进口特有数值字段
      exchange_rate: item.exchange_rate ?? 0,
      sea_freight: item.sea_freight ?? 0,
      cif_price: item.cif_price ?? 0,
      cif_total: item.cif_total ?? 0,
      cif_bs_total: item.cif_bs_total ?? 0,
      // 替换 import_duty_vat: item.import_duty_vat ?? 0
      import_duty: item.import_duty ?? 0,
      vat: item.vat ?? 0,
      transportation: item.transportation ?? 0,
      others_charge: item.others_charge ?? 0,
      total_cost: item.total_cost ?? 0,
      import_unit_cost: item.import_unit_cost ?? 0,
      piece_area: item.piece_area ?? 0,
      system_price: item.system_price ?? 0,
      share_rate: item.share_rate ?? 0,
    })),
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
