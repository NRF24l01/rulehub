<template>
    <div class="flex items-center justify-center min-h-screen bg-gray-100">
        <form
            class="bg-white p-8 rounded shadow-md w-full max-w-sm"
            @submit.prevent="handleLogin"
        >
            <h2 class="text-2xl font-bold mb-6 text-center">Войти</h2>
            <div class="mb-4">
                <label class="block mb-1 font-medium" for="nickname">Имя пользователя</label>
                <input
                    id="nickname"
                    v-model="nickname"
                    type="text"
                    class="w-full px-3 py-2 border rounded focus:outline-none focus:ring"
                    :class="{'border-red-500': nicknameError}"
                    autocomplete="username"
                />
                <p v-if="nicknameError" class="text-red-500 text-sm mt-1">{{ nicknameError }}</p>
            </div>
            <div class="mb-6">
                <label class="block mb-1 font-medium" for="password">Пароль</label>
                <input
                    id="password"
                    v-model="password"
                    type="password"
                    class="w-full px-3 py-2 border rounded focus:outline-none focus:ring"
                    :class="{'border-red-500': passwordError}"
                    autocomplete="current-password"
                />
                <p v-if="passwordError" class="text-red-500 text-sm mt-1">{{ passwordError }}</p>
            </div>
            <button
                type="submit"
                class="w-full bg-blue-600 text-white py-2 rounded hover:bg-blue-700 transition"
                :disabled="loading"
            >
                {{ loading ? 'Входимс...' : 'Войти' }}
            </button>
            <p v-if="loginError" class="text-red-500 text-center mt-4">{{ loginError }}</p>
        </form>
    </div>
</template>

<script setup>
import { ref } from 'vue'

const nickname = ref('')
const password = ref('')
const nicknameError = ref('')
const passwordError = ref('')
const loginError = ref('')
const loading = ref(false)

function validateNickname(value) {
    if (!value) return 'Ник нужен'
    if (value.length < 3 || value.length > 32)
        return 'Ник должен быть 3-32 символа'
    if (!/^[a-zA-Z0-9_]+$/.test(value))
        return 'Ник может содержать только буквы, цифры и подчеркивания'
    return ''
}

function validatePassword(value) {
    if (!value) return 'Пароль обязателен'
    if (value.length < 6 || value.length > 128)
        return 'Пароль должен быть 6-128 символов'
    if (!/[a-zA-Z]/.test(value) || !/[0-9]/.test(value))
        return 'Пароль должен содержать буквы и цифры'
    return ''
}

async function handleLogin() {
    nicknameError.value = validateNickname(nickname.value)
    passwordError.value = validatePassword(password.value)
    loginError.value = ''

    if (nicknameError.value || passwordError.value) return

    loading.value = true
    try {
        const res = await fetch(
            `${import.meta.env.VITE_BACKEND_URL}/auth/login`,
            {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    username: nickname.value,
                    password: password.value,
                }),
                credentials: 'include',
            }
        )
        if (res.status === 200) {
            // Successful login, handle as needed (e.g., redirect)
            // Optionally, get access_token from response if needed
            // const data = await res.json()
            // localStorage.setItem('access_token', data.access_token)
            window.location.reload()
        } else if (res.status === 401) {
            loginError.value = 'Неверные учетные данные'
        } else {
            loginError.value = 'Ошибка входа, попробуйте позже'
        }
    } catch (e) {
        loginError.value = 'Ошибка сети'
    } finally {
        loading.value = false
    }
}
</script>