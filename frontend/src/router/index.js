import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
    },
    {
      path: "/login",
      name: "login",
      component: () => import('../views/Login.vue'),
      meta: { guestOnly: true }
    },
    {
      path: "/create",
      name: "create",
      component: () => import('../views/Create.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: "/article/:id",
      name: "article",
      component: () => import('../views/ArticleRead.vue'),
    }
  ],
})

router.beforeEach(async (to, from, next) => {
  const auth = useAuthStore()
  
  if (to.meta.guestOnly && auth.isAuthenticated) {
    next({ name: 'home' })
  } else if (to.meta.requiresAuth && !auth.isAuthenticated) {
    // Try to refresh token if not authenticated
    const refreshed = await auth.refreshToken()
    
    // Check authentication status again after refresh attempt
    if (refreshed && auth.isAuthenticated) {
      next() // Continue to requested route if refresh worked
    } else {
      next({ name: 'login' }) // Redirect to login if refresh failed
    }
  } else {
    next()
  }
})

export default router
