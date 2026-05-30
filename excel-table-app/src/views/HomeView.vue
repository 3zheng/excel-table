<template>
  <el-container style="height: 100vh;">

    <!-- 侧边栏 -->
    <el-aside width="200px">
      <div class="logo">📊 货柜管理</div>
      <el-menu
        :default-active="$route.path"
        router
        background-color="#304156"
        text-color="#bfcbd9"
        active-text-color="#409eff"
      >
        <!-- 出口表单，A/B类和管理员可见 -->
        <el-menu-item
          index="/export"
          v-if="canSeeExport"
        >
          <el-icon><Ship /></el-icon>
          <span>{{ t('exportList') }}</span>
        </el-menu-item>

        <!-- 进口表单，C/D类和管理员可见 -->
        <el-menu-item
          index="/import"
          v-if="canSeeImport"
        >
          <el-icon><Box /></el-icon>
          <span>{{ t('importList') }}</span>
        </el-menu-item>

        <!-- 产品属性表，有编辑权限或管理员可见 -->
        <el-menu-item
          index="/products"
          v-if="canSeeProducts"
        >
          <el-icon><Goods /></el-icon>
          <span>{{ t('products') }}</span>
        </el-menu-item>

        <!-- 用户管理，仅管理员可见 -->
        <el-menu-item
          index="/users"
          v-if="isAdmin"
        >
          <el-icon><User /></el-icon>
          <span>{{ t('userManagement') }}</span>
        </el-menu-item>
      </el-menu>
    </el-aside>

    <!-- 主内容区 -->
    <el-container>
      <el-header>
        <div class="header-right">
          <span>{{ user.username }}</span>
          <el-button type="danger" size="small" @click="logout">
            {{ t('logout') }}
          </el-button>
        </div>
      </el-header>
      <el-main>
        <router-view />
      </el-main>
    </el-container>

  </el-container>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { Ship, Box, Goods, User } from '@element-plus/icons-vue'

const { t } = useI18n()
const router = useRouter()

const user = JSON.parse(localStorage.getItem('user') || '{}')
const role = user.user_role || ''

const isAdmin = computed(() => role === 'admin')
const canSeeExport = computed(() => ['admin', 'export_input', 'export_review'].includes(role))
const canSeeImport = computed(() => ['admin', 'import_input', 'import_review'].includes(role))
const canSeeProducts = computed(() => user.can_edit_products || role === 'admin')

const logout = () => {
  localStorage.removeItem('token')
  localStorage.removeItem('user')
  router.push('/login')
}
</script>

<style scoped>
.el-aside {
  background-color: #304156;
}
.logo {
  height: 60px;
  line-height: 60px;
  text-align: center;
  color: white;
  font-size: 16px;
  font-weight: 500;
  border-bottom: 1px solid #435163;
}
:deep(.el-menu) {
  border-right: none;
  width: 200px;
}
.el-header {
  background: white;
  border-bottom: 1px solid #e6e6e6;
  display: flex;
  align-items: center;
  justify-content: flex-end;
}
.header-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

</style>