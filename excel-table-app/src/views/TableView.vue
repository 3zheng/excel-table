<template>
  <div class="container">
    <h2>📊 {{ t('total') }}</h2>

    <div class="toolbar">
      <el-button type="primary" @click="addRow">{{ t('addRow') }}</el-button>
      <el-button type="danger" @click="deleteSelected" :disabled="!selectedRows.length">
        {{ t('deleteSelected') }} ({{ selectedRows.length }})
      </el-button>
    </div>

    <el-table
      :data="tableData"
      border
      stripe
      @selection-change="handleSelect"
      style="width: 100%"
    >
      <el-table-column type="selection" width="50" />
      <el-table-column :label="t('index')" type="index" width="60" />

      <el-table-column :label="t('name')" prop="name" min-width="120">
        <template #default="{ row }">
          <el-input v-model="row.name" size="small" @change="() => updateRow(row)" />
        </template>
      </el-table-column>

      <el-table-column :label="t('dept')" prop="dept" min-width="120">
        <template #default="{ row }">
          <el-select v-model="row.dept" size="small" style="width:100%" @change="() => updateRow(row)">
            <el-option :label="t('depts.rd')" value="研发" />
            <el-option :label="t('depts.product')" value="产品" />
            <el-option :label="t('depts.operation')" value="运营" />
            <el-option :label="t('depts.marketing')" value="市场" />
          </el-select>
        </template>
      </el-table-column>

      <el-table-column :label="t('salary')" prop="salary" min-width="120">
        <template #default="{ row }">
          <el-input-number v-model="row.salary" :min="0" size="small" style="width:100%" @change="() => updateRow(row)" />
        </template>
      </el-table-column>

      <el-table-column :label="t('date')" prop="date" min-width="160">
        <template #default="{ row }">
          <el-date-picker
            v-model="row.date"
            type="date"
            size="small"
            value-format="YYYY-MM-DD"
            style="width:100%"
            @change="() => updateRow(row)"
          />
        </template>
      </el-table-column>
    </el-table>

    <div class="summary">
      {{ t('total') }}：<strong>¥{{ totalSalary.toLocaleString() }}</strong>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

interface Row {
  id: number
  name: string
  dept: string
  salary: number
  date: string
}

const tableData = ref<Row[]>([])
const selectedRows = ref<Row[]>([])

const totalSalary = computed(() =>
  tableData.value.reduce((sum, r) => sum + (r.salary || 0), 0)
)

const updateRow = async (row: Row) => {
  const res = await fetch(`/api/rows/${row.id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(row)
  })
  console.log('updateRow response', res.status)
}
// 获取所有数据
const fetchRows = async () => {
  const res = await fetch('/api/rows')
  tableData.value = await res.json()
}

// 新增行
const addRow = async () => {
  const res = await fetch('/api/rows', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ name: '', dept: '研发', salary: 0, date: '' })
  })
  const newRow = await res.json()
  tableData.value.push(newRow)
}

// 删除选中
const deleteSelected = async () => {
  await Promise.all(
    selectedRows.value.map(row =>
      fetch(`/api/rows/${row.id}`, { method: 'DELETE' })
    )
  )
  tableData.value = tableData.value.filter(
    r => !selectedRows.value.includes(r)
  )
  selectedRows.value = []
}

const handleSelect = (rows: Row[]) => {
  selectedRows.value = rows
}

onMounted(() => {
  fetchRows()
})
</script>

<style>
.container { max-width: 900px; margin: 40px auto; padding: 0 20px; }
h2 { margin-bottom: 16px; font-weight: 500; }
.toolbar { margin-bottom: 12px; display: flex; gap: 8px; }
.summary { margin-top: 12px; text-align: right; font-size: 14px; color: #606266; }
</style>