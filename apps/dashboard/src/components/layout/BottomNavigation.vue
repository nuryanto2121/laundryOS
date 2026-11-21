<template>
  <nav
    class="md:hidden fixed bottom-0 left-0 right-0 z-40 bg-white/95 backdrop-blur-lg border-t border-slate-200/80 px-2 py-1 shadow-lg"
  >
    <div
      class="grid items-center"
      :class="auth.userRole.value === 'KASIR' ? 'grid-cols-3' : 'grid-cols-4'"
    >
      <!-- 1. Beranda -->
      <RouterLink
        to="/"
        class="flex flex-col items-center justify-center py-2 px-1 rounded-xl transition-all relative"
        :class="
          route.path === '/'
            ? 'text-brand-600 font-bold'
            : 'text-slate-600 hover:text-slate-800 font-medium'
        "
      >
        <div class="relative p-1 rounded-lg" :class="route.path === '/' ? 'bg-brand-50' : ''">
          <Home class="w-5 h-5" />
        </div>
        <span class="text-[11px] mt-0.5">Beranda</span>
        <span v-if="route.path === '/'" class="w-4 h-1 bg-brand-600 rounded-full mt-0.5"></span>
      </RouterLink>

      <!-- 2. Transaksi / Antrian -->
      <RouterLink
        to="/orders"
        class="flex flex-col items-center justify-center py-2 px-1 rounded-xl transition-all relative"
        :class="
          route.path.startsWith('/orders')
            ? 'text-brand-600 font-bold'
            : 'text-slate-600 hover:text-slate-800 font-medium'
        "
      >
        <div
          class="relative p-1 rounded-lg"
          :class="route.path.startsWith('/orders') ? 'bg-brand-50' : ''"
        >
          <ShoppingBag class="w-5 h-5" />
          <!-- Active Badge Counter if orders in queue -->
          <span
            v-if="queueCount > 0"
            class="absolute -top-1 -right-2 px-1.5 py-0.2 rounded-full bg-rose-500 text-white text-[9px] font-extrabold shadow-xs animate-pulse"
          >
            {{ queueCount }}
          </span>
        </div>
        <span class="text-[11px] mt-0.5">Transaksi</span>
        <span
          v-if="route.path.startsWith('/orders')"
          class="w-4 h-1 bg-brand-600 rounded-full mt-0.5"
        ></span>
      </RouterLink>

      <!-- 3. Laporan -->
      <RouterLink
        to="/reports"
        class="flex flex-col items-center justify-center py-2 px-1 rounded-xl transition-all relative"
        :class="
          route.path.startsWith('/reports')
            ? 'text-brand-600 font-bold'
            : 'text-slate-600 hover:text-slate-800 font-medium'
        "
      >
        <div
          class="relative p-1 rounded-lg"
          :class="route.path.startsWith('/reports') ? 'bg-brand-50' : ''"
        >
          <BarChart3 class="w-5 h-5" />
        </div>
        <span class="text-[11px] mt-0.5">Laporan</span>
        <span
          v-if="route.path.startsWith('/reports')"
          class="w-4 h-1 bg-brand-600 rounded-full mt-0.5"
        ></span>
      </RouterLink>

      <!-- 4. Pengaturan -->
      <RouterLink
        v-if="auth.userRole.value !== 'KASIR'"
        to="/settings"
        class="flex flex-col items-center justify-center py-2 px-1 rounded-xl transition-all relative"
        :class="
          route.path.startsWith('/settings')
            ? 'text-brand-600 font-bold'
            : 'text-slate-600 hover:text-slate-800 font-medium'
        "
      >
        <div
          class="relative p-1 rounded-lg"
          :class="route.path.startsWith('/settings') ? 'bg-brand-50' : ''"
        >
          <Settings class="w-5 h-5" />
        </div>
        <span class="text-[11px] mt-0.5">Setting</span>
        <span
          v-if="route.path.startsWith('/settings')"
          class="w-4 h-1 bg-brand-600 rounded-full mt-0.5"
        ></span>
      </RouterLink>
    </div>
  </nav>
</template>

<script setup>
import { useRoute } from 'vue-router'
import { Home, ShoppingBag, BarChart3, Settings } from 'lucide-vue-next'
import { useAuth } from '../../composables/useAuth'

const route = useRoute()
const auth = useAuth()

defineProps({
  queueCount: {
    type: Number,
    default: 0,
  },
})
</script>
