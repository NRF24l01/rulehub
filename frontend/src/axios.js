import axios from 'axios'
import { useAuthStore } from '@/stores/auth'

const api = axios.create({
  baseURL: '/api',
  withCredentials: true
})

// Extract refresh token function to make it reusable
export async function refreshAccessToken() {
  const res = await axios.post(import.meta.env.VITE_BACKEND_URL + '/auth/refresh', {}, {
    withCredentials: true
  })
  return res.data.access_token
}

api.interceptors.request.use(config => {
  const auth = useAuthStore()
  if (auth.accessToken) {
    config.headers.Authorization = 'Bearer ' + auth.accessToken
  }
  return config
})

api.interceptors.response.use(
  response => response,
  async error => {
    const auth = useAuthStore()

    if (error.response && error.response.status === 401) {
      try {
        const newToken = await refreshAccessToken()
        auth.setToken(newToken)

        const config = error.config
        config.headers.Authorization = 'Bearer ' + newToken
        return api(config)
      } catch (refreshError) {
        auth.logout()
        return Promise.reject(refreshError)
      }
    }

    return Promise.reject(error)
  }
)

export default api
