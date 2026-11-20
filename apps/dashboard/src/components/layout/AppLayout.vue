<template>
  <div class="min-h-screen bg-slate-50 flex flex-col text-slate-800 overflow-x-hidden">
    <!-- Top Header -->
    <AppHeader :active-outlet-name="activeOutlet" @open-outlet-modal="showOutletModal = true" />

    <!-- Desktop Navigation Tab Bar (Visible on md and up) -->
    <div class="hidden md:block bg-white border-b border-slate-200/80 sticky top-16 z-20 shadow-xs">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <nav class="flex space-x-6">
          <RouterLink
            v-for="item in navItems"
            :key="item.path"
            :to="item.path"
            class="flex items-center gap-2 py-3.5 px-3 border-b-2 text-sm font-semibold transition-all relative"
            :class="
              isCurrentPath(item.path)
                ? 'border-brand-600 text-brand-600'
                : 'border-transparent text-slate-500 hover:text-slate-800 hover:border-slate-300'
            "
          >
            <component :is="item.icon" class="w-4 h-4" />
            <span>{{ item.label }}</span>
            <span
              v-if="item.badge"
              class="px-1.5 py-0.2 rounded-full bg-rose-500 text-white text-[10px] font-bold"
            >
              {{ item.badge }}
            </span>
          </RouterLink>
        </nav>
      </div>
    </div>

    <!-- Main Page Content Container -->
    <main class="flex-1 max-w-7xl w-full mx-auto px-4 sm:px-6 lg:px-8 py-4 sm:py-6 pb-24 md:pb-10">
      <slot
        :open-expense-modal="() => (showExpenseModal = true)"
        :open-search-modal="() => (showSearchModal = true)"
      />
    </main>

    <!-- Mobile Bottom Navigation (Sticky bottom on phone/tablet) -->
    <BottomNavigation :queue-count="queueCount" />

    <!-- Global Floating Toast Notification Container -->
    <ToastContainer />

    <!-- Outlet Switcher Modal -->
    <div
      v-if="showOutletModal"
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/40 backdrop-blur-xs"
    >
      <div class="bg-white rounded-2xl max-w-md w-full p-5 shadow-2xl border border-slate-100">
        <div class="flex items-center justify-between pb-3 border-b border-slate-100">
          <div class="flex items-center gap-2">
            <Store class="w-5 h-5 text-brand-600" />
            <h3 class="font-bold text-slate-800">Pilih Outlet Aktif</h3>
          </div>
          <button
            class="p-1.5 rounded-lg text-slate-400 hover:text-slate-600"
            @click="showOutletModal = false"
          >
            <X class="w-5 h-5" />
          </button>
        </div>
        <div class="py-4 space-y-2">
          <div
            class="p-3 rounded-xl border-2 flex items-center justify-between cursor-pointer transition-all"
            :class="
              activeOutlet === 'FAN Laundry Kalibaru'
                ? 'border-brand-500 bg-brand-50/50'
                : 'border-slate-200 hover:border-slate-300'
            "
            @click="selectOutlet('FAN Laundry Kalibaru')"
          >
            <div>
              <p class="font-bold text-sm text-slate-900">FAN Laundry Kalibaru</p>
              <p class="text-xs text-slate-500">
                Jl. Kalibaru Barat I Rt 007 Rw 008, Cilincing, Jakarta Utara
              </p>
            </div>
            <CheckCircle2
              v-if="activeOutlet === 'FAN Laundry Kalibaru'"
              class="w-5 h-5 text-brand-600"
            />
          </div>
        </div>
        <button
          class="w-full py-2.5 rounded-xl bg-slate-900 text-white font-semibold text-sm hover:bg-slate-800"
          @click="showOutletModal = false"
        >
          Tutup
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRoute } from 'vue-router'
import { Home, ShoppingBag, BarChart3, Settings, Store, X, CheckCircle2 } from 'lucide-vue-next'
import AppHeader from './AppHeader.vue'
import BottomNavigation from './BottomNavigation.vue'
import ToastContainer from '../ui/ToastContainer.vue'

const route = useRoute()
const activeOutlet = ref('FAN Laundry Kalibaru')
const showOutletModal = ref(false)
const showExpenseModal = ref(false)
const showSearchModal = ref(false)
const queueCount = ref(0)

const isCurrentPath = (path) => {
  if (path === '/') return route.path === '/'
  return route.path.startsWith(path)
}

const navItems = computed(() => [
  { path: '/', label: 'Beranda', icon: Home },
  {
    path: '/orders',
    label: 'Transaksi & Antrian',
    icon: ShoppingBag,
    badge: queueCount.value > 0 ? queueCount.value : null,
  },
  { path: '/reports', label: 'Laporan Keuangan & CRM', icon: BarChart3 },
  { path: '/settings', label: 'Pengaturan & Master Data', icon: Settings },
])

const selectOutlet = (name) => {
  activeOutlet.value = name
  showOutletModal.value = false
}
</script>
