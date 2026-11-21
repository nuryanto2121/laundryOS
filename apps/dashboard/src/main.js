import { createApp } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'

import './assets/main.css'
import App from './App.vue'
import HomeView from './views/HomeView.vue'
import OrdersView from './views/OrdersView.vue'
import ReportsView from './views/ReportsView.vue'
import SettingsView from './views/SettingsView.vue'
import MobilePreviewView from './views/MobilePreviewView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
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

const app = createApp(App)
app.use(router)
app.mount('#app')
