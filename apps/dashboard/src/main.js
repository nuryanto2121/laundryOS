import { createApp } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'

import './assets/main.css'
import App from './App.vue'
import HomeView from './views/HomeView.vue'
import OrdersView from './views/OrdersView.vue'
import ReportsView from './views/ReportsView.vue'
import SettingsView from './views/SettingsView.vue'
import MobilePreviewView from './views/MobilePreviewView.vue'
import LoginView from './views/LoginView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: LoginView,
      meta: { public: true },
    },
    {
      path: '/',
      name: 'home',
      component: HomeView,
    },
    {
      path: '/orders',
      name: 'orders',
      component: OrdersView,
    },
    {
      path: '/reports',
      name: 'reports',
      component: ReportsView,
    },
    {
      path: '/settings',
      name: 'settings',
      component: SettingsView,
    },
    {
      path: '/mobile',
      name: 'mobile',
      component: MobilePreviewView,
    },
    {
      path: '/mobile-preview',
      redirect: '/mobile',
    },
    {
      path: '/mobile-preview.html',
      redirect: '/mobile',
    },
    {
      path: '/:pathMatch(.*)*',
      redirect: '/',
    },
  ],
})

// Navigation Guard for Authentication and Settings RBAC
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('laundryos_auth_token')
  const rawUser = localStorage.getItem('laundryos_auth_user')
  const isAuthenticated = Boolean(token && rawUser)

  let user = null
  if (rawUser) {
    try {
      user = JSON.parse(rawUser)
    } catch (e) {
      user = null
    }
  }

  // Allow public routes
  if (to.meta.public) {
    if (isAuthenticated) {
      return next('/')
    }
    return next()
  }

  // Redirect to login if unauthenticated
  if (!isAuthenticated) {
    return next({ path: '/login', query: { redirect: to.fullPath } })
  }

  // Restrict Kasir from accessing Settings
  if (to.path === '/settings' && user?.role === 'KASIR') {
    alert('Akses Ditolak: Karyawan tidak diizinkan mengakses menu Pengaturan.')
    return next('/')
  }

  next()
})

const app = createApp(App)
app.use(router)
app.mount('#app')
