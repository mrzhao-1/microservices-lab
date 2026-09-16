<script setup>
import { ref } from 'vue'

const username = ref('')
const email = ref('')
const result = ref(null)
const error = ref('')
const loading = ref(false)

async function register() {
  loading.value = true
  error.value = ''
  result.value = null
  try {
    const resp = await fetch('/api/user/register', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username: username.value, email: email.value }),
    })
    const data = await resp.json()
    if (!resp.ok) {
      error.value = data.error || '请求失败 HTTP ' + resp.status
    } else {
      result.value = data
    }
  } catch (e) {
    error.value = '调用失败：' + e.message
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <main class="wrap">
    <h1>用户注册</h1>
    <p class="hint">前端 → nginx → gateway(go-zero) → user-svc(go-kratos) → notification-svc(go-micro)</p>
    <form @submit.prevent="register">
      <label>用户名 <input v-model="username" placeholder="zhaohu" /></label>
      <label>邮箱 <input v-model="email" placeholder="zhaohu@example.com" /></label>
      <button type="submit" :disabled="loading">{{ loading ? '提交中…' : '注册' }}</button>
    </form>
    <p v-if="result" class="ok">注册成功：ID {{ result.id }}，用户名 {{ result.username }}</p>
    <p v-if="error" class="err">{{ error }}</p>
  </main>
</template>

<style>
body { font-family: system-ui, sans-serif; background: #f5f6f8; }
.wrap { max-width: 360px; margin: 80px auto; background: #fff; padding: 32px; border-radius: 12px; box-shadow: 0 2px 12px rgba(0,0,0,.06); }
.hint { color: #888; font-size: 13px; line-height: 1.6; }
form { display: flex; flex-direction: column; gap: 14px; margin-top: 20px; }
label { display: flex; flex-direction: column; gap: 4px; font-size: 14px; color: #333; }
input { padding: 10px; border: 1px solid #ddd; border-radius: 6px; font-size: 14px; }
button { padding: 10px; background: #1677ff; color: #fff; border: none; border-radius: 6px; font-size: 15px; cursor: pointer; }
button:disabled { opacity: .6; cursor: not-allowed; }
.ok { color: #16a34a; }
.err { color: #dc2626; word-break: break-all; }
</style>
