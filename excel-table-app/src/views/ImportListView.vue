<template>
  <div>
    <!-- 搜索栏 -->
    <el-card class="search-card">
      <el-form inline>
        <el-form-item :label="t('invNo')">
          <el-input v-model="filters.inv_no" :placeholder="t('search')" clearable style="width:180px" />
        </el-form-item>
        <el-form-item :label="t('region')">
          <el-select v-model="filters.region" :placeholder="t('all')" clearable style="width:140px">
            <el-option v-for="r in availableRegions" :key="r" :label="r" :value="r" />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('invoiceDate')">
          <el-date-picker v-model="filters.start_date" type="date" value-format="YYYY-MM-DD"
            :placeholder="t('startDate')" style="width:150px" />
        </el-form-item>
        <el-form-item label="~">
          <el-date-picker v-model="filters.end_date" type="date" value-format="YYYY-MM-DD"
            :placeholder="t('endDate')" style="width:150px" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="fetchList">{{ t('search') }}</el-button>
          <el-button @click="resetFilters">{{ t('reset') }}</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 列表 -->
    <el-card>
      <el-table :data="list" border style="width:100%" v-loading="loading"
        @row-click="(row: any) => goDetail(row.inv_no)">
        <el-table-column :label="t('invNo')" prop="inv_no" min-width="150" />
        <el-table-column :label="t('region')" prop="region" width="120" />
        <el-table-column :label="t('invoiceDate')" prop="invoice_date" width="120" />
        <el-table-column :label="t('contractNo')" prop="contract_no" min-width="150" />
        <el-table-column :label="t('shippingLine')" prop="shipping_line" min-width="120" />
        <el-table-column :label="t('reviewStatus')" prop="reviewed" width="100">
          <template #default="{ row }">
            <el-tag :type="row.reviewed ? 'success' : 'warning'" size="small">
              {{ row.reviewed ? t('reviewed') : t('pending') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('reviewer')" prop="reviewer" width="100" />
        <el-table-column :label="t('createdAt')" prop="created_at" width="160" />
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { request } from '@/utils/request'

const { t } = useI18n()
const router = useRouter()

const loading = ref(false)
const list = ref<any[]>([])

// 管理员可以看所有地区，非管理员从接口返回数据里动态提取
const availableRegions = ref<string[]>([])

const filters = ref({
  inv_no: '',
  region: '',
  start_date: '',
  end_date: '',
})

const fetchList = async () => {
  loading.value = true
  try {
    const params = new URLSearchParams()
    if (filters.value.inv_no) params.append('inv_no', filters.value.inv_no)
    if (filters.value.region) params.append('region', filters.value.region)
    if (filters.value.start_date) params.append('start_date', filters.value.start_date)
    if (filters.value.end_date) params.append('end_date', filters.value.end_date)

    const res = await request(`/api/invoices/import?${params.toString()}`)
    const data = await res.json()
    list.value = data
    const regions = [...new Set(data.map((item: any) => item.region).filter(Boolean))] as string[]
    availableRegions.value = regions
  } catch {
    // request 已统一处理错误提示和 401 跳转
  } finally {
    loading.value = false
  }
}

const resetFilters = () => {
  filters.value = { inv_no: '', region: '', start_date: '', end_date: '' }
  fetchList()
}

const goDetail = (invNo: string) => {
  router.push(`/import/${invNo}`)
}

onMounted(() => {
  fetchList()
})
</script>

<style scoped>
.search-card {
  margin-bottom: 16px;
}

:deep(.el-table__row) {
  cursor: pointer;
}
</style>
