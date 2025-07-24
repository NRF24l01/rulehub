<template>
    <div class="flex items-center justify-center min-h-screen bg-gray-100">
        <div class="bg-white p-8 rounded-lg shadow-lg w-full max-w-md">
            <h2 class="text-2xl font-bold mb-6 text-center text-gray-800">Войти</h2>
            
            <form @submit.prevent="handleLogin">
                <!-- Username field -->
                <div class="mb-4">
                    <label class="block mb-2 font-medium text-gray-700" for="nickname">
                        Имя пользователя
                    </label>
                    <input
                        id="nickname"
                        v-model="nickname"
                        type="text"
                        class="w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 transition"
                        :class="{'border-red-500': nicknameError}"
                        placeholder="Введите имя пользователя"
                        autocomplete="username"
                    />
                    <p v-if="nicknameError" class="text-red-500 text-sm mt-1">{{ nicknameError }}</p>
                </div>
                
                <!-- Password field -->
                <div class="mb-6">
                    <label class="block mb-2 font-medium text-gray-700" for="password">
                        Пароль
                    </label>
                    <div class="relative">
                        <input
                            id="password"
                            v-model="password"
                            :type="showPassword ? 'text' : 'password'"
                            class="w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 transition"
                            :class="{'border-red-500': passwordError}"
                            placeholder="Введите пароль"
                            autocomplete="current-password"
                        />
                        <button 
                            type="button" 
                            class="absolute right-3 top-2.5 text-gray-500 hover:text-gray-700"
                            @click="showPassword = !showPassword"
                        >
                            {{ showPassword ? '🙈' : '👁️' }}
                        </button>
                    </div>
                    <p v-if="passwordError" class="text-red-500 text-sm mt-1">{{ passwordError }}</p>
                </div>

                <!-- Submit button -->
                <button
                    type="submit"
                    class="w-full bg-blue-600 text-white py-3 rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 transition"
                    :disabled="loading || isFormInvalid"
                    :class="{'opacity-70 cursor-not-allowed': loading || isFormInvalid}"
                >
                    <span v-if="loading" class="inline-flex items-center">
                        <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                        </svg>
                        Входимс...
                    </span>
                    <span v-else>Войти</span>
                </button>
            </form>

            <!-- Error message -->
            <div v-if="loginError" class="mt-4 p-3 bg-red-100 border border-red-400 text-red-700 rounded">
                {{ loginError }}
            </div>

            <!-- Register link -->
            <p class="mt-6 text-center text-gray-600">
                Нет аккаунта? 
                <router-link to="/register" class="text-blue-600 hover:underline">Не сможете зарегистрироваться</router-link>
            </p>
        </div>
    </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'

const auth = useAuthStore()
const router = useRouter()

const nickname = ref('')
const password = ref('')
const nicknameError = ref('')
const passwordError = ref('')
const loginError = ref('')
const loading = ref(false)
const showPassword = ref(false)

// Новые computed для проверки валидности
const isNicknameValid = computed(() => {
    return (
        nickname.value &&
        nickname.value.length >= 3 &&
        nickname.value.length <= 32 &&
        /^[a-zA-Z0-9_]+$/.test(nickname.value)
    )
})

const isPasswordValid = computed(() => {
    return (
        password.value &&
        password.value.length >= 6 &&
        password.value.length <= 128 &&
        /[a-zA-Z]/.test(password.value) &&
        /[0-9]/.test(password.value)
    )
})

const isFormInvalid = computed(() => {
    return !isNicknameValid.value || !isPasswordValid.value || loading.value
})

// Validation functions (только для ошибок)
function validateNickname() {
    if (!nickname.value) {
        nicknameError.value = 'Ник нужен'
        return false
    }
    if (nickname.value.length < 3 || nickname.value.length > 32) {
        nicknameError.value = 'Ник должен быть 3-32 символа'
        return false
    }
    if (!/^[a-zA-Z0-9_]+$/.test(nickname.value)) {
        nicknameError.value = 'Ник может содержать только буквы, цифры и подчеркивания'
        return false
    }
    nicknameError.value = ''
    return true
}

function validatePassword() {
    if (!password.value) {
        passwordError.value = 'Пароль обязателен'
        return false
    }
    if (password.value.length < 6 || password.value.length > 128) {
        passwordError.value = 'Пароль должен быть 6-128 символов'
        return false
    }
    if (!/[a-zA-Z]/.test(password.value) || !/[0-9]/.test(password.value)) {
        passwordError.value = 'Пароль должен содержать буквы и цифры'
        return false
    }
    passwordError.value = ''
    return true
}

// Показывать ошибки сразу при вводе
watch(nickname, () => {
    validateNickname()
})

watch(password, () => {
    validatePassword()
})

async function handleLogin() {
  loginError.value = ''

  const isNicknameValid = validateNickname()
  const isPasswordValid = validatePassword()
  if (!isNicknameValid || !isPasswordValid) return

  loading.value = true

  try {
    const res = await fetch(
      `${import.meta.env.VITE_BACKEND_URL}/auth/login`,
      {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          username: nickname.value.trim(),
          password: password.value,
        }),
        credentials: 'include', // Важно, чтобы куки (refresh token) отправились
      }
    )

    if (res.ok) {
      const data = await res.json()
      if (data.access_token) {
        auth.setToken(data.access_token) // Сохраняем access_token в Pinia
        router.push('/dashboard') // Переходим в SPA на dashboard без перезагрузки
      } else {
        loginError.value = 'Сервер не вернул access_token'
      }
    } else if (res.status === 401) {
      loginError.value = 'Неверные учетные данные'
    } else {
      const errorData = await res.json().catch(() => null)
      loginError.value = errorData?.message || 'Ошибка входа, попробуйте позже'
    }
  } catch (e) {
    console.error('Login error:', e)
    loginError.value = 'Ошибка сети. Проверьте подключение к интернету.'
  } finally {
    loading.value = false
  }
}
</script>