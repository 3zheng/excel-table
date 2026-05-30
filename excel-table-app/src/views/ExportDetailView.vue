<template>
  <div>
    <!-- 顶部操作栏 -->
    <div class="action-bar">
      <el-button @click="goBack">← 返回列表</el-button>
      <div class="action-right">
        <el-button
          type="primary"
          @click="saveInvoice"
          v-if="canEdit"
          :loading="saving"
        >保存</el-button>
        <el-button
          type="danger"
          @click="deleteInvoice"
          v-if="isAdmin"
        >删除</el-button>
      </div>
    </div>

    <!-- 表单头部信息 -->
    <el-card class="section-card">
      <template #header>表单信息</template>
      <el-form :model="invoice" label-width="140px" :disabled="!canEdit">
        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item label="INV. NO">
              <el-input v-model="invoice.inv_no" :disabled="!isNew" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="日期">
              <el-date-picker
                v-model="invoice.invoice_date"
                type="date"
                value-format="YYYY-MM-DD"
                style="width:100%"
              />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="合同号">
              <el-input v-model="invoice.contract_no" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="船运公司">
              <el-input v-model="invoice.shipping_line" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="提单号 B/L NO">
              <el-input v-model="invoice.bl_no" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="柜号">
              <el-input v-model="invoice.container_no" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="装载港">
              <el-input v-model="invoice.port_of_loading" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="卸货港">
              <el-input v-model="invoice.port_of_discharge" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="目的地">
              <el-input v-model="invoice.final_destination" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="买方">
              <el-input v-model="invoice.buyer_name" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="买方电话">
              <el-input v-model="invoice.buyer_tel" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="卖方">
              <el-input v-model="invoice.seller_name" />
            </el-form-item>
          </el-col>
          <el-col :span="16">
            <el-form-item label="买方地址">
              <el-input v-model="invoice.buyer_address" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="付款条件">
              <el-input v-model="invoice.payment_term" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
    </el-card>

    <!-- 审核区域 -->
    <el-card class="section-card" v-if="!isNew">
      <template #header>审核信息</template>
      <el-form label-width="140px">
        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item label="审核状态">
              <el-switch
                v-model="invoice.reviewed"
                :disabled="!canReview"
                active-text="已审核"
                inactive-text="待审核"
              />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="审核人">
              <el-input :value="invoice.reviewer" disabled />
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item label="审核意见">
              <el-input
                v-model="invoice.review_comment"
                type="textarea"
                :rows="2"
                :disabled="!canReview"
              />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
    </el-card>

    <!-- 明细行 -->
    <el-card class="section-card">
      <template #header>
        <div class="card-header">
          <span>产品明细</span>
          <el-button
            type="primary"
            size="small"
            @click="addItem"
            v-if="canEdit"
          >新增行</el-button>
        </div>
      </template>
      <el-table :data="invoice.items" border style="width:100%">
        <el-table-column label="序号" type="index" width="60" />
        <el-table-column label="品牌" prop="brand" min-width="100">
          <template #default="{ row }">
            <el-input v-model="row.brand" size="small" :disabled="!canEdit" />
          </template>
        </el-table-column>
        <el-table-column label="产品名" prop="commodities" min-width="160">
          <template #default="{ row }">
            <el-input v-model="row.commodities" size="small" :disabled="!canEdit" />
          </template>
        </el-table-column>
        <el-table-column label="型号" prop="model_no" min-width="120">
          <template #default="{ row }">
            <el-input v-model="row.model_no" size="small" :disabled="!canEdit" />
          </template>
        </el-table-column>
        <el-table-column label="规格" prop="descriptions" min-width="160">
          <template #default="{ row }">
            <el-input v-model="row.descriptions" size="small" :disabled="!canEdit" />
          </template>
        </el-table-column>
        <el-table-column label="箱数" prop="carton_qty" width="90">
          <template #default="{ row }">
            <el-input-number v-model="row.carton_qty" size="small" :min="0" style="width:100%" :disabled="!canEdit" :controls="false"/>
          </template>
        </el-table-column>
        <el-table-column label="总数量" prop="unit_qty" width="90">
          <template #default="{ row }">
            <el-input-number v-model="row.unit_qty" size="small" :min="0" style="width:100%" :disabled="!canEdit" :controls="false"/>
          </template>
        </el-table-column>
        <el-table-column label="单位" prop="unit" width="80">
          <template #default="{ row }">
            <el-input v-model="row.unit" size="small" :disabled="!canEdit" />
          </template>
        </el-table-column>
        <el-table-column label="单价" prop="unit_price" width="100">
          <template #default="{ row }">
            <el-input-number v-model="row.unit_price" size="small" :min="0" :precision="3" style="width:100%" :disabled="!canEdit" :controls="false"/>
          </template>
        </el-table-column>
        <el-table-column label="总金额" prop="total_amount" width="100">
          <template #default="{ row }">
            <el-input-number v-model="row.total_amount" size="small" :min="0" :precision="2" style="width:100%" :disabled="!canEdit" :controls="false"/>
          </template>
        </el-table-column>
        <el-table-column label="毛重" prop="gross_weight" width="90">
          <template #default="{ row }">
            <el-input-number v-model="row.gross_weight" size="small" :min="0" :precision="2" style="width:100%" :disabled="!canEdit" :controls="false"/>
          </template>
        </el-table-column>
        <el-table-column label="净重" prop="net_weight" width="90">
          <template #default="{ row }">
            <el-input-number v-model="row.net_weight" size="small" :min="0" :precision="2" style="width:100%" :disabled="!canEdit" :controls="false"/>
          </template>
        </el-table-column>
        <el-table-column label="CBM" prop="volume" width="90">
          <template #default="{ row }">
            <el-input-number v-model="row.volume" size="small" :min="0" :precision="2" style="width:100%" :disabled="!canEdit" :controls="false"/>
          </template>
        </el-table-column>
        <el-table-column label="操作" fixed="right" width="80" v-if="canEdit">
          <template #default="{ $index }">
            <el-button size="small" type="danger" @click="removeItem($index)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- FOB 合计 -->
      <div class="fob-total">
        FOB TOTAL: <strong>USD {{ fobTotal.toFixed(2) }}</strong>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'

const router = useRouter()
const route = useRoute()

const user = JSON.parse(localStorage.getItem('user') || '{}')
const role = user.user_role || ''
const token = localStorage.getItem('token') || ''

const isNew = computed(() => route.params.inv_no === 'new')
const isAdmin = computed(() => role === 'admin')
const canEdit = computed(() => ['admin', 'export_input'].includes(role))
const canReview = computed(() => ['admin', 'export_review'].includes(role))

const saving = ref(false)

console.log('route.params:', route.params)
console.log('isNew:', isNew.value)

const invoice = ref({
  inv_no: '',
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

const fetchDetail = async () => {
  const invNo = route.params.inv_no as string
  const res = await fetch(`/api/invoices/export/${invNo}`, {
    headers: { Authorization: `Bearer ${token}` }
  })
  if (res.ok) {
    const data = await res.json()
    invoice.value = data
  }
}

const addItem = () => {
  invoice.value.items.push({
    brand: '',
    commodities: '',
    model_no: '',
    descriptions: '',
    carton_qty: 0,
    unit_qty: 0,
    unit: '',
    unit_price: 0,
    total_amount: 0,
    gross_weight: 0,
    net_weight: 0,
    volume: 0,
  })
}

const removeItem = (index: number) => {
  invoice.value.items.splice(index, 1)
}

const saveInvoice = async () => {
  saving.value = true
  // 处理日期格式，只保留 YYYY-MM-DD 部分
  const payload = {
    ...invoice.value,
    invoice_date: invoice.value.invoice_date 
      ? invoice.value.invoice_date.slice(0, 10) 
      : null
  }
  console.log('saveInvoice payload:', JSON.stringify(payload))
  console.log('invoice.review_comment:', invoice.value.review_comment)
  
  try {
    const method = isNew.value ? 'POST' : 'PUT'
    const url = isNew.value
      ? '/api/invoices/export'
      : `/api/invoices/export/${route.params.inv_no}`

    const res = await fetch(url, {
      method,
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`
      },
      body: JSON.stringify(payload)
    })
    const data = await res.json()
    if (res.ok) {
      ElMessage.success('保存成功')
      if (isNew.value) {
        router.replace(`/export/${invoice.value.inv_no}`)
      }
    } else {
      ElMessage.error(data.error || '保存失败')
    }
  } finally {
    saving.value = false
  }
}

const deleteInvoice = async () => {
  await ElMessageBox.confirm('确定要删除这张表单吗？', '警告', {
    type: 'warning'
  })
  const res = await fetch(`/api/invoices/export/${route.params.inv_no}`, {
    method: 'DELETE',
    headers: { Authorization: `Bearer ${token}` }
  })
  if (res.ok) {
    ElMessage.success('删除成功')
    router.push('/export')
  }
}

const goBack = () => {
  router.push('/export')
}

onMounted(() => {
  if (!isNew.value) {
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
</style>