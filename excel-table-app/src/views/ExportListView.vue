<template>
  <div>
    <!-- 搜索栏 -->
    <div class="search-bar">
      <el-input
        v-model="searchInvNo"
        placeholder="搜索 INV. NO"
        style="width: 200px"
        clearable
        @clear="fetchList"
        @keyup.enter="fetchList"
      />
      <el-date-picker
        v-model="dateRange"
        type="daterange"
        range-separator="~"
        start-placeholder="开始日期"
        end-placeholder="结束日期"
        value-format="YYYY-MM-DD"
        style="width: 260px"
      />
      <el-button type="primary" @click="fetchList">搜索</el-button>
      <el-button @click="resetSearch">重置</el-button>
      <el-button
        type="success"
        @click="goCreate"
        v-if="canCreate"
      >新建表单</el-button>
    </div>

    <!-- 表格 -->
    <el-table :data="list" border stripe style="width:100%">
      <el-table-column label="INV. NO" prop="inv_no" width="160" />
      <el-table-column label="日期" prop="invoice_date" width="120" />
      <el-table-column label="合同号" prop="contract_no" width="160" />
      <el-table-column label="船运公司" prop="shipping_line" width="140" />
      <el-table-column label="审核状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.reviewed ? 'success' : 'warning'">
            {{ row.reviewed ? '已审核' : '待审核' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="审核人" prop="reviewer" width="120" />
      <el-table-column label="创建时间" prop="created_at" width="180" />
      <el-table-column label="操作" fixed="right" width="100">
        <template #default="{ row }">
          <el-button size="small" @click="goDetail(row.inv_no)">
            查看
          </el-button>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const user = JSON.parse(localStorage.getItem('user') || '{}')
const role = user.user_role || ''

const canCreate = computed(() => ['admin', 'export_input'].includes(role))

const searchInvNo = ref('')
const dateRange = ref<string[]>([])
const list = ref([])

const getToken = () => localStorage.getItem('token') || ''

const fetchList = async () => {
  const params = new URLSearchParams()
  if (searchInvNo.value) params.append('inv_no', searchInvNo.value)
  if (dateRange.value && dateRange.value.length === 2) {
    params.append('start_date', dateRange.value[0] ?? '')
    params.append('end_date', dateRange.value[1] ?? '')
  }

  const res = await fetch(`/api/invoices/export?${params.toString()}`, {
    headers: { Authorization: `Bearer ${getToken()}` }
  })
  if (res.status === 401) {
    router.push('/login')
    return
  }
  const data = await res.json()
  list.value = data
}

const resetSearch = () => {
  searchInvNo.value = ''
  dateRange.value = []
  fetchList()
}

const goDetail = (invNo: string) => {
  router.push(`/export/${invNo}`)
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