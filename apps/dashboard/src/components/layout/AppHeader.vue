<template>
  <header
    class="sticky top-0 z-30 bg-white/90 backdrop-blur-md border-b border-slate-100 shadow-xs"
  >
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 h-16 flex items-center justify-between">
      <!-- Left: Brand Logo & Title -->
      <div class="flex items-center gap-3">
        <div
          class="w-10 h-10 rounded-xl bg-gradient-to-tr from-brand-600 to-sky-400 flex items-center justify-center text-white shadow-md shadow-brand-500/20"
        >
          <Sparkles class="w-5 h-5" />
        </div>
        <div>
          <div class="flex items-center gap-2">
            <span class="font-extrabold text-lg tracking-tight text-slate-900">FAN Laundry</span>
            <span
              class="text-[10px] font-bold px-2 py-0.5 rounded-full bg-brand-50 text-brand-600 border border-brand-200 uppercase tracking-wider"
            >
              POS Station
            </span>
          </div>
          <p class="text-xs text-slate-500 hidden sm:block">
            Kalibaru Barat I, Cilincing — Jakarta Utara
          </p>
        </div>
      </div>

      <!-- Right: Outlet Switcher & Status Profile -->
      <div class="flex items-center gap-2 sm:gap-4">
        <!-- Outlet Switcher Trigger -->
        <button
          class="flex items-center gap-2 px-3 py-1.5 rounded-xl bg-slate-50 hover:bg-slate-100 border border-slate-200/80 transition-all text-xs font-semibold text-slate-700"
          @click="$emit('open-outlet-modal')"
        >
          <Store class="w-3.5 h-3.5 text-brand-600" />
          <span class="truncate max-w-[120px] sm:max-w-none">{{ activeOutletName }}</span>
          <ChevronDown class="w-3 h-3 text-slate-400" />
        </button>

        <!-- Live Sync Pill -->
        <div
          class="hidden sm:flex items-center gap-1.5 px-2.5 py-1 rounded-full bg-emerald-50 text-emerald-700 border border-emerald-200 text-xs font-medium"
        >
          <span class="w-2 h-2 rounded-full bg-emerald-500 animate-pulse"></span>
          <span>Online</span>
        </div>

        <!-- Cashier Avatar -->
        <div class="flex items-center gap-2 pl-2 border-l border-slate-200">
          <div
            class="w-8 h-8 rounded-full bg-slate-800 text-white flex items-center justify-center text-xs font-bold shadow-xs"
          >
            K1
          </div>
          <div class="hidden md:block text-left">
            <p class="text-xs font-bold text-slate-800 leading-none">Kasir Utama</p>
            <p class="text-[10px] text-slate-500">FAN Laundry</p>
          </div>
        </div>
      </div>
    </div>
  </header>
</template>

<script setup>
import { Sparkles, Store, ChevronDown } from 'lucide-vue-next'

defineProps({
  activeOutletName: {
    type: String,
    default: 'FAN Laundry Kalibaru',
  },
})

defineEmits(['open-outlet-modal'])
</script>
