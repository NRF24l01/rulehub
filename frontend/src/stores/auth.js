import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useAuthStore = defineStore('auth', () => {
  const accessToken = ref(null)

  const isAuthenticated = computed(() => !!accessToken.value)
  const authHeader = computed(() =>
    accessToken.value ? { Authorization: 'Bearer ' + accessToken.value } : {}
  )

  function setToken(token) {
    accessToken.value = token
  }

  function logout() {
    accessToken.value = null
  }

  return {
    accessToken,
    isAuthenticated,
    authHeader,
    setToken,
    logout,
  }
}, {
  persist: {
    enabled: true,
    strategies: [
      { storage: sessionStorage }
    ]
  }
})
