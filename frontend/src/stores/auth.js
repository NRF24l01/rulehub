import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { refreshAccessToken } from '@/axios'

function isTokenExpired(token) {
  if (!token) return true
  try {
    const payload = JSON.parse(atob(token.split('.')[1]))
    if (!payload.exp) return false
    return Date.now() >= payload.exp * 1000
  } catch (e) {
    return true
  }
}

export const useAuthStore = defineStore('auth', () => {
  const accessToken = ref(null)

  const isAuthenticated = computed(() =>
    !!accessToken.value && !isTokenExpired(accessToken.value)
  )
  const authHeader = computed(() =>
    accessToken.value ? { Authorization: 'Bearer ' + accessToken.value } : {}
  )

  function setToken(token) {
    accessToken.value = token
  }

  function logout() {
    accessToken.value = null
  }

  // Add method to refresh token
  async function refreshToken() {
    try {
      const newToken = await refreshAccessToken()
      setToken(newToken)
      return true
    } catch (error) {
      logout()
      return false
    }
  }

  return {
    accessToken,
    isAuthenticated,
    authHeader,
    setToken,
    logout,
    refreshToken
  }
}, {
  persist: {
    enabled: true,
    strategies: [
      { storage: sessionStorage }
    ]
  }
})
