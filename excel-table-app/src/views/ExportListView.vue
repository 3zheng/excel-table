<template>
  <div>
    <!-- 搜索栏 -->
    <div class="search-bar">
      <el-input
        v-model="searchInvNo"
        :placeholder="t('searchInvNo')"
        style="width: 200px"
        clearable
        @clear="fetchList"
        @keyup.enter="fetchList"
      />
      <el-select v-model="searchRegion" :placeholder="t('region')" clearable style="width:140px">
        <el-option :label="t('bolivia')" value="玻利维亚" />
        <el-option :label="t('peru')" value="秘鲁" />
        <el-option :label="t('chile')" value="智利" />
        <el-option :label="t('spain')" value="西班牙" />
        <el-option :label="t('usa')" value="美国" />
      </el-select>
      <el-date-picker
        v-model="dateRange"
        type="daterange"
        range-separator="~"
        :start-placeholder="t('startDate')"
        :end-placeholder="t('endDate')"
        value-format="YYYY-MM-DD"
        style="width: 260px"
      />
      <el-button type="primary" @click="fetchList">{{ t('search') }}</el-button>
      <el-button @click="resetSearch">{{ t('reset') }}</el-button>
      <el-button
        type="success"
        @click="goCreate"
        v-if="canCreate"
      >{{ t('newInvoice') }}</el-button>
    </div>

    <!-- 表格 -->
    <el-table :data="list" border stripe style="width:100%">
      <el-table-column label="INV. NO" prop="inv_no" width="160" />
      <el-table-column :label="t('region')" prop="region" width="100" />
      <el-table-column :label="t('invoiceDate')" prop="invoice_date" width="120" />
      <el-table-column :label="t('contractNo')" prop="contract_no" width="160" />
      <el-table-column :label="t('shippingLine')" prop="shipping_line" width="140" />
      <el-table-column :label="t('reviewStatus')" width="100">
        <template #default="{ row }">
          <el-tag :type="row.reviewed ? 'success' : 'warning'">
            {{ row.reviewed ? t('reviewed'): t('pending') }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column :label="t('reviewer')" prop="reviewer" width="120" />
      <el-table-column :label="t('createdAt')" prop="created_at" width="180" />
      <el-table-column :label="t('operation')" fixed="right" width="150">
        <template #default="{ row }">
          <el-button size="small" @click="goDetail(row.inv_no)">
            {{ t('view') }}
          </el-button>
          <el-button size="small" type="success" @click="goCopy(row.inv_no)">
            {{ t('copy') }}
          </el-button>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
//import { ElMessage } from 'element-plus' //会与unplugin-auto-import存在重复导入的冲突，导致提示框无法正确渲染，所以注释掉
import { request } from '@/utils/request'

const { t } = useI18n()

const router = useRouter()
const user = JSON.parse(localStorage.getItem('user') || '{}')
const role = user.user_role || ''

const canCreate = computed(() => ['admin', 'export_input'].includes(role))

const searchInvNo = ref('')
const searchRegion = ref('')
const dateRange = ref<string[]>([])
const list = ref([])

const fetchList = async () => {
  const params = new URLSearchParams()
  if (searchInvNo.value) params.append('inv_no', searchInvNo.value)
  if (searchRegion.value) params.append('region', searchRegion.value)
  if (dateRange.value && dateRange.value.length === 2) {
    params.append('start_date', dateRange.value[0] ?? '')
    params.append('end_date', dateRange.value[1] ?? '')
  }

  try {
    const res = await request(`/api/invoices/export?${params.toString()}`)
    const data = await res.json()
    list.value = data
  } catch {
    // request 已统一处理错误提示和 401 跳转
  }
}

const resetSearch = () => {
  searchInvNo.value = ''
  searchRegion.value = ''
  dateRange.value = []
  fetchList()
}

const goDetail = (invNo: string) => {
  router.push(`/export/${invNo}`)
}

const goCopy = (invNo: string) => {
  router.push(`/export/new?copy=${invNo}`)
}

const goCreate = () => {
  router.push('/export/new')
}

onMounted(() => {
  fetchList()
})
</script>

<style scoped>
.search-bar {
  display: flex;
  gap: 10px;
  margin-bottom: 16px;
  flex-wrap: wrap;
  align-items: center;
}
</style>