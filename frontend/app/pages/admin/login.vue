<script setup lang="ts">
definePageMeta({ layout: 'admin' })

const { adminUser, authChecked, fetchMe, login } = useAdminAuth()

if (!authChecked.value) {
  await fetchMe()
}
if (adminUser.value) {
  await navigateTo('/admin')
}

const loginValue = ref('')
const password = ref('')
const error = ref('')
const loading = ref(false)

async function onSubmit() {
  error.value = ''
  loading.value = true
  try {
    await login(loginValue.value, password.value)
    await navigateTo('/admin')
  } catch {
    error.value = 'Неверный логин или пароль'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="mx-auto flex min-h-screen max-w-sm flex-col justify-center px-4">
    <h1 class="mb-6 text-2xl font-semibold">Вход в админку</h1>
    <form class="flex flex-col gap-4" @submit.prevent="onSubmit">
      <input
        v-model="loginValue"
        type="text"
        placeholder="Логин"
        autocomplete="username"
        required
        class="rounded border border-gray-300 px-3 py-2"
      >
      <input
        v-model="password"
        type="password"
        placeholder="Пароль"
        autocomplete="current-password"
        required
        class="rounded border border-gray-300 px-3 py-2"
      >
      <p v-if="error" class="text-sm text-red-600">{{ error }}</p>
      <button
        type="submit"
        :disabled="loading"
        class="rounded bg-[var(--color-primary)] px-4 py-2 text-white disabled:opacity-50"
      >
        {{ loading ? 'Входим…' : 'Войти' }}
      </button>
    </form>
  </div>
</template>
