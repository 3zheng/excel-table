<template>
  <div class="login-container">
    <div class="login-box">
      <h2>{{ t('login') }}</h2>
      <el-form :model="form" @submit.prevent="handleLogin">
        <el-form-item>
          <el-input
            v-model="form.username"
            :placeholder="t('username')"
            prefix-icon="User"
          />
        </el-form-item>
        <el-form-item>
          <el-input
            v-model="form.password"
            :placeholder="t('password')"
            type="password"
            prefix-icon="Lock"
            show-password
          />
        </el-form-item>
        <el-form-item>
          <el-button
            type="primary"
            style="width:100%"
            :loading="loading"
            @click="handleLogin"
          >
            {{ t('login') }}
          </el-button>
        </el-form-item>
        <p v-if="errorMsg" class="error-msg">{{ errorMsg }}</p>
      </el-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const router = useRouter()

const form = ref({
  username: '',
  password: '',
})
const loading = ref(false)
const errorMsg = ref('')

const handleLogin = async () => {
  if (!form.value.username || !form.value.password) {
    errorMsg.value = t('loginEmptyError')
    return
  }
  loading.value = true
  errorMsg.value = ''
  try {
    const res = await fetch('/api/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(form.value)
    })
    const data = await res.json()
    if (res.ok) {
      localStorage.setItem('token', data.token)
      localStorage.setItem('user', JSON.stringify(data.user))
      router.push('/')
    } else {
      errorMsg.value = data.error
    }
  } catch {
    errorMsg.value = t('loginNetworkError')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container {
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f5f7fa;
}
.login-box {
  width: 360px;
  padding: 40px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0,0,0,0.1);
}
h2 {
  text-align: center;
  margin-bottom: 24px;
  font-weight: 500;
}
.error-msg {
  color: #f56c6c;
  font-size: 13px;
  text-align: center;
  margin-top: -8px;
}
</style>